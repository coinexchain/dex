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
	TokenStoreKeyPrefix = []byte{0x01}
	WhitelistKeyPrefix  = []byte{0x02}
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
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramstore params.Subspace,
	ak auth.AccountKeeper, fck auth.FeeCollectionKeeper) TokenKeeper {

	return TokenKeeper{
		key:           key,
		ak:            ak,
		fck:           fck,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// setToken  implements token Keeper.
func (tk TokenKeeper) setToken(ctx sdk.Context, token Token) sdk.Error {
	symbol := token.GetSymbol()
	store := ctx.KVStore(tk.key)

	bz, err := tk.cdc.MarshalBinaryBare(token)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(TokenStoreKey(symbol), bz)
	return nil
}

// GetToken implements token Keeper.
func (tk TokenKeeper) GetToken(ctx sdk.Context, symbol string) Token {
	store := ctx.KVStore(tk.key)
	bz := store.Get(TokenStoreKey(symbol))
	if bz == nil {
		return nil
	}
	return tk.decodeToken(bz)
}

// GetAllTokens returns all tokens in the token Keeper.
func (tk TokenKeeper) GetAllTokens(ctx sdk.Context) []Token {
	tokens := make([]Token, 0)

	tk.IterateToken(ctx, func(token Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	return tokens
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
		acc := tk.decodeToken(iter.Value())
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// removeToken remove an token for the asset mapper store.
func (tk TokenKeeper) removeToken(ctx sdk.Context, token Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(tk.key)
	store.Delete(TokenStoreKey(symbol))
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

	return tk.setToken(ctx, token)
}

func (tk TokenKeeper) checkPrecondition(ctx sdk.Context, msg sdk.Msg, symbol string, owner sdk.AccAddress) (Token, sdk.Error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	token := tk.GetToken(ctx, symbol)
	if token == nil {
		return nil, ErrorTokenNotFound(fmt.Sprintf("token %s not found", symbol))
	}

	if !token.GetOwner().Equals(owner) {
		return nil, ErrorInvalidTokenOwner("Only token owner can do this action")
	}

	return token, nil
}

//TransferOwnership - transfer token owner
func (tk TokenKeeper) TransferOwnership(ctx sdk.Context, msg MsgTransferOwnership) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OriginalOwner)
	if err != nil {
		return err
	}

	if err := token.SetOwner(msg.NewOwner); err != nil {
		return ErrorInvalidTokenOwner("token new owner is invalid")
	}

	return tk.setToken(ctx, token)
}

//MintToken - mint token
func (tk TokenKeeper) MintToken(ctx sdk.Context, msg MsgMintToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetMintable() {
		return ErrorInvalidTokenMint(fmt.Sprintf("token %s do not support mint", msg.Symbol))
	}

	if err := token.SetTotalMint(msg.Amount + token.GetTotalMint()); err != nil {
		return ErrorInvalidTokenMint(err.Error())
	}

	if err := token.SetTotalSupply(msg.Amount + token.GetTotalSupply()); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	return tk.setToken(ctx, token)
}

//BurnToken - burn token
func (tk TokenKeeper) BurnToken(ctx sdk.Context, msg MsgBurnToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetBurnable() {
		return ErrorInvalidTokenBurn(fmt.Sprintf("token %s do not support burn", msg.Symbol))
	}

	if err := token.SetTotalBurn(msg.Amount + token.GetTotalBurn()); err != nil {
		return ErrorInvalidTokenBurn(err.Error())
	}

	if err := token.SetTotalSupply(token.GetTotalSupply() - msg.Amount); err != nil {
		return ErrorInvalidTokenSupply(err.Error())
	}

	return tk.setToken(ctx, token)
}

func (tk TokenKeeper) addWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(tk.key)
	for _, acc := range whitelist {
		store.Set(WhitelistKey(symbol, acc), nil)
	}

	return nil
}

func (tk TokenKeeper) removeWhitelist(ctx sdk.Context, symbol string, whitelist []sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(tk.key)
	for _, acc := range whitelist {
		store.Delete(WhitelistKey(symbol, acc))
	}

	return nil
}

//ForbidToken - forbid token
func (tk TokenKeeper) ForbidToken(ctx sdk.Context, msg MsgForbidToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid", msg.Symbol))
	}
	if token.GetIsForbidden() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s has been forbidden", msg.Symbol))
	}
	token.SetIsForbidden(true)

	return tk.setToken(ctx, token)
}

//UnForbidToken - unforbid token
func (tk TokenKeeper) UnForbidToken(ctx sdk.Context, msg MsgUnForbidToken) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support unforbid", msg.Symbol))
	}
	if !token.GetIsForbidden() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s has not been forbidden", msg.Symbol))
	}
	token.SetIsForbidden(false)

	return tk.setToken(ctx, token)
}

func (tk TokenKeeper) AddTokenForbidWhitelist(ctx sdk.Context, msg MsgAddForbidWhitelist) sdk.Error {
	token, err := tk.checkPrecondition(ctx, msg, msg.Symbol, msg.OwnerAddress)
	if err != nil {
		return err
	}

	if !token.GetTokenForbiddable() {
		return ErrorInvalidTokenForbidden(fmt.Sprintf("token %s do not support forbid and add whitelist", msg.Symbol))
	}
	if err = tk.addWhitelist(ctx, msg.Symbol, msg.Whitelist); err != nil {
		return ErrorInvalidTokenWhitelist(fmt.Sprintf("token whitelist is invalid"))
	}
	return tk.setToken(ctx, token)
}

// -----------------------------------------------------------------------------
// ExpectedAssertStatusKeeper

//IsTokenForbidden - check whether the coin's owner has forbidden "denom", forbiding transmission and exchange.
func (tk TokenKeeper) IsTokenFrozen(ctx sdk.Context, denom string) bool {
	token := tk.GetToken(ctx, denom)
	if token != nil {
		return token.GetIsForbidden()
	}

	return true
}

// IsTokenExists - check whether there is a coin named "denom"
func (tk TokenKeeper) IsTokenExists(ctx sdk.Context, denom string) bool {
	return tk.GetToken(ctx, denom) != nil
}

// IsTokenIssuer - check whether addr is a token issuer
func (tk TokenKeeper) IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	if addr.Empty() {
		return false
	}

	token := tk.GetToken(ctx, denom)
	return token != nil && token.GetOwner().Equals(addr)
}

func (tk TokenKeeper) IsForbiddenByTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	//TODO:
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

// WhitelistKey - return WhitelistKeyPrefix-Symbol-AccAddress KEY
func WhitelistKey(symbol string, acc sdk.AccAddress) []byte {
	return append(append(WhitelistKeyPrefix, []byte(symbol)...), acc.Bytes()...)
}

func (tk TokenKeeper) decodeToken(bz []byte) (token Token) {
	if err := tk.cdc.UnmarshalBinaryBare(bz, &token); err != nil {
		panic(err)
	}
	return
}
