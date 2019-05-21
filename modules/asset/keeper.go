package asset

import (
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
	TokenStoreKeyPrefix = []byte{0x00}
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
	bz := store.Get([]byte(symbol))
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
func (tk TokenKeeper) SetToken(ctx sdk.Context, token Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(tk.key)
	bz, err := tk.cdc.MarshalBinaryBare(token)
	if err != nil {
		panic(err)
	}
	store.Set([]byte(symbol), bz)

}

//IssueToken - new token and store
func (tk TokenKeeper) IssueToken(ctx sdk.Context, msg MsgIssueToken) (err sdk.Error) {

	token, err := NewToken(msg.Name, msg.Symbol, msg.TotalSupply, msg.Owner,
		msg.Mintable, msg.Burnable, msg.AddrFreezeable, msg.TokenFreezeable)

	if err != nil {
		return err
	}
	tk.SetToken(ctx, token)

	return nil
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
