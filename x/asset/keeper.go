package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// StoreKey is string representation of the store key for asset
	StoreKey = "asset"

	// QuerierRoute is the querier route for asset
	QuerierRoute = StoreKey
)

var (
	// TokenStoreKeyPrefix prefix for asset-by-TokenSymbol store
	TokenStoreKeyPrefix = []byte{0x00}
)

// Keeper encodes/decodes tokens using the go-amino (binary)
// encoding/decoding library.
type Keeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical token constructor.
	proto func() Token

	// The codec codec for binary encoding/decoding of token.
	cdc *codec.Codec

	paramSubspace params.Subspace
}

// NewKeeper returns a new Keeper that uses go-amino to
// (binary) encode and decode concrete Token.
// nolint
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore params.Subspace, proto func() Token,
) Keeper {

	return Keeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// GetToken implements token Keeper.
func (keeper Keeper) GetToken(ctx sdk.Context, symbol string) Token {
	store := ctx.KVStore(keeper.key)
	bz := store.Get(TokenStoreKey(symbol))
	if bz == nil {
		return nil
	}
	token := keeper.decodeToken(bz)
	return token
}

// GetAllTokens returns all tokens in the token Keeper.
func (keeper Keeper) GetAllTokens(ctx sdk.Context) []Token {
	var tokens []Token
	appendToken := func(token Token) (stop bool) {
		tokens = append(tokens, token)
		return false
	}
	keeper.IterateAccounts(ctx, appendToken)
	return tokens
}

// RemoveToken removes an token for the asset mapper store.
func (keeper Keeper) RemoveToken(ctx sdk.Context, token Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(keeper.key)
	store.Delete(TokenStoreKey(symbol))
}

// IterateToken implements token Keeper
func (keeper Keeper) IterateAccounts(ctx sdk.Context, process func(Token) (stop bool)) {
	store := ctx.KVStore(keeper.key)
	iter := sdk.KVStorePrefixIterator(store, TokenStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := keeper.decodeToken(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// SetToken  implements token Keeper.
func (keeper Keeper) SetToken(ctx sdk.Context, token Token) {
	symbol := token.GetSymbol()
	store := ctx.KVStore(keeper.key)
	bz, err := keeper.cdc.MarshalBinaryBare(token)
	if err != nil {
		panic(err)
	}
	store.Set(TokenStoreKey(symbol), bz)

}

func (keeper Keeper) IssueToken(ctx sdk.Context, token MsgIssueToken) (tags sdk.Tags, err sdk.Error) {
	//TODO:
	//deduct the fee from issuer’s account
	//New token info is saved on the CoinEx Chain
	return
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (keeper Keeper) SetParams(ctx sdk.Context, params Params) {
	keeper.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (keeper Keeper) GetParams(ctx sdk.Context) (params Params) {
	keeper.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Misc.

// TokenStoreKey turn an token symbol to key used to get it from the asset store
func TokenStoreKey(symbol string) []byte {
	return append(TokenStoreKeyPrefix, []byte(symbol)...)
}

func (keeper Keeper) decodeToken(bz []byte) (token Token) {
	err := keeper.cdc.UnmarshalBinaryBare(bz, &token)
	if err != nil {
		panic(err)
	}
	return
}

