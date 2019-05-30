package asset

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "asset"

	// StoreKey is string representation of the store key for asset
	StoreKey = ModuleName

	// RouterKey is the message route for asset
	RouterKey = ModuleName

	// QuerierRoute is the querier route for asset
	QuerierRoute = ModuleName
)

var (
	// TokenStoreKeyPrefix prefix for asset-by-TokenSymbol store
	TokenStoreKeyPrefix = []byte{0x01}
)

// Keeper encodes/decodes tokens using the go-amino (binary)
// encoding/decoding library.
type TokenKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	ak  auth.AccountKeeper
	fck auth.FeeCollectionKeeper

	// The codec codec for binary encoding/decoding of token.
	cdc *codec.Codec

	paramSubspace params.Subspace
}

// NewKeeper returns a new Keeper that uses go-amino to
// (binary) encode and decode concrete Token.
// nolint
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore params.Subspace,
	ak auth.AccountKeeper, fck auth.FeeCollectionKeeper) TokenKeeper {

	return TokenKeeper{
		key:           key,
		ak:            ak,
		fck:           fck,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// GetToken implements token Keeper.
func (tk TokenKeeper) GetToken(ctx sdk.Context, symbol string) Token {
	store := ctx.KVStore(tk.key)
	bz := store.Get(TokenStoreKey(symbol))
	if bz == nil {
		return nil
	}
	token := tk.decodeToken(bz)
	return token
}

// GetAllTokens returns all tokens in the token Keeper.
func (tk TokenKeeper) GetAllTokens(ctx sdk.Context) []Token {
	var tokens []Token
	appendToken := func(token Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	}
	tk.IterateToken(ctx, appendToken)
	return tokens
}

// RemoveToken removes an token for the asset mapper store.
func (tk TokenKeeper) RemoveToken(ctx sdk.Context, token Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(tk.key)
	store.Delete(TokenStoreKey(symbol))
}

// IterateToken implements token Keeper
func (tk TokenKeeper) IterateToken(ctx sdk.Context, process func(Token) (stop bool)) {
	store := ctx.KVStore(tk.key)
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := tk.decodeToken(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// SetToken  implements token Keeper.
func (tk TokenKeeper) SetToken(ctx sdk.Context, token Token) sdk.Error {
	symbol := token.GetSymbol()
	store := ctx.KVStore(tk.key)

	bz, err := tk.cdc.MarshalBinaryBare(token)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(TokenStoreKey(symbol), bz)
	return nil
}

//IssueToken - new token and store
func (tk TokenKeeper) IssueToken(ctx sdk.Context, msg MsgIssueToken) sdk.Error {

	token, err := NewToken(msg.Name, msg.Symbol, msg.TotalSupply, msg.Owner,
		msg.Mintable, msg.Burnable, msg.AddrForbiddable, msg.TokenForbiddable)

	if err != nil {
		return err
	}

	if tk.IsTokenExists(ctx, token.Symbol) {
		return ErrorDuplicateTokenSymbol(fmt.Sprintf("token symbol already exists in store"))
	}

	err = tk.SetToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

//TransferOwnership - transfer token owner
func (tk TokenKeeper) TransferOwnership(ctx sdk.Context, msg MsgTransferOwnership) sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return ErrorInvalidTokenOwner(err.Error())
	}

	token := tk.GetToken(ctx, msg.Symbol)
	if token == nil {
		return ErrorNoTokenPersist("transfer invalid token`s ownership")
	}
	if !token.GetOwner().Equals(msg.OriginalOwner) {
		return ErrorInvalidTokenOwner("token original owner is invalid")
	}

	if err := token.SetOwner(msg.NewOwner); err != nil {
		return ErrorInvalidTokenOwner("token new owner is invalid")
	}
	if err := tk.SetToken(ctx, token); err != nil {
		return nil
	}
	return nil
}

//MintToken - mint token
func (tk TokenKeeper) MintToken(ctx sdk.Context, msg MsgMintToken) sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return ErrorInvalidTokenMint(err.Error())
	}

	token := tk.GetToken(ctx, msg.Symbol)
	if token == nil {
		return ErrorNoTokenPersist("mint invalid token")
	}
	if !token.GetOwner().Equals(msg.OwnerAddress) {
		return ErrorInvalidTokenOwner("only token owner can mint token")
	}
	if !token.GetMintable() {
		return ErrorInvalidTokenMint("token that cannot be minted")
	}

	amt := msg.Amount
	preMint := token.GetTotalMint()
	if err := token.SetTotalMint(amt + preMint); err != nil {
		return ErrorInvalidTokenMint(err.Error())
	}
	preSupply := token.GetTotalSupply()
	if err := token.SetTotalSupply(amt + preSupply); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	if err := tk.SetToken(ctx, token); err != nil {
		return nil
	}
	return nil
}

//BurnToken - burn token
func (tk TokenKeeper) BurnToken(ctx sdk.Context, msg MsgBurnToken) sdk.Error {
	if err := msg.ValidateBasic(); err != nil {
		return ErrorInvalidTokenMint(err.Error())
	}

	token := tk.GetToken(ctx, msg.Symbol)
	if token == nil {
		return ErrorNoTokenPersist("burn invalid token")
	}
	if !token.GetOwner().Equals(msg.OwnerAddress) {
		return ErrorInvalidTokenOwner("only token owner can burn token")
	}
	if !token.GetBurnable() {
		return ErrorInvalidTokenMint("token that cannot be burn")
	}

	amt := msg.Amount
	preMint := token.GetTotalBurn()
	if err := token.SetTotalBurn(amt + preMint); err != nil {
		return ErrorInvalidTokenMint(err.Error())
	}
	preSupply := token.GetTotalSupply()
	if err := token.SetTotalSupply(preSupply - amt); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	if err := tk.SetToken(ctx, token); err != nil {
		return nil
	}
	return nil
}

// -----------------------------------------------------------------------------
// ExpectedAssertStatusKeeper

//IsTokenFrozen - check whether the coin's owner has frozen "denom", forbiding transmission and exchange.
func (tk TokenKeeper) IsTokenFrozen(ctx sdk.Context, denom string) bool {
	token := tk.GetToken(ctx, denom)
	if token != nil {
		return token.GetIsFrozen()
	}

	return true
}

// IsTokenExists - check whether there is a coin named "denom"
func (tk TokenKeeper) IsTokenExists(ctx sdk.Context, denom string) bool {
	return tk.GetToken(ctx, denom) != nil
}

// IsTokenIssuer - check whether addr is a token issuer
func (tk TokenKeeper) IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	token := tk.GetToken(ctx, denom)
	if token != nil && token.GetOwner().Equals(addr) {
		return true
	}
	return false
}

func (tk TokenKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	//TODO: fzc
	return false
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (tk TokenKeeper) SetParams(ctx sdk.Context, params Params) {
	tk.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (tk TokenKeeper) GetParams(ctx sdk.Context) (params Params) {
	tk.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Misc.

// TokenStoreKey turn an token symbol to key used to get it from the asset store
func TokenStoreKey(symbol string) []byte {
	return append(TokenStoreKeyPrefix, []byte(symbol)...)
}

func (tk TokenKeeper) decodeToken(bz []byte) (token Token) {
	err := tk.cdc.UnmarshalBinaryBare(bz, &token)
	if err != nil {
		panic(err)
	}
	return
}
