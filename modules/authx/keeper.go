package authx

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// StoreKey is string representation of the store key for authx
	StoreKey = "accx"
	// QuerierRoute is the querier route for accx
	QuerierRoute = StoreKey
)

var (
	// AddressStoreKeyPrefix prefix for accountx-by-address store
	AddressStoreKeyPrefix = []byte{0x01}
)

type AccountXKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSubspace params.Subspace) AccountXKeeper {
	return AccountXKeeper{
		key:           key,
		cdc:           cdc,
		paramSubspace: paramSubspace.WithKeyTable(ParamKeyTable()),
	}
}

func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}

// -----------------------------------------------------------------------------
// AccountX

func (axk AccountXKeeper) GetOrCreateAccountX(ctx sdk.Context, addr sdk.AccAddress) AccountX {
	ax, ok := axk.GetAccountX(ctx, addr)
	if !ok {
		ax = AccountX{Address: addr}
		axk.SetAccountX(ctx, ax)
	}
	return ax
}

func (axk AccountXKeeper) GetAccountX(ctx sdk.Context, addr sdk.AccAddress) (ax AccountX, ok bool) {
	store := ctx.KVStore(axk.key)
	bz := store.Get(AddressStoreKey(addr))
	if bz == nil {
		return
	}
	acc := axk.decodeAccountX(bz)
	return acc, true
}

func (axk AccountXKeeper) SetAccountX(ctx sdk.Context, ax AccountX) {
	addr := ax.Address
	store := ctx.KVStore(axk.key)
	bz, err := axk.cdc.MarshalBinaryBare(ax)
	if err != nil {
		panic(err)
	}
	store.Set(AddressStoreKey(addr), bz)
}

func (axk AccountXKeeper) IterateAccounts(ctx sdk.Context, process func(AccountX) (stop bool)) {
	store := ctx.KVStore(axk.key)
	iter := sdk.KVStorePrefixIterator(store, AddressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := axk.decodeAccountX(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (axk AccountXKeeper) SetParams(ctx sdk.Context, params Params) {
	axk.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (axk AccountXKeeper) GetParams(ctx sdk.Context) (params Params) {
	axk.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Codec

func (axk AccountXKeeper) decodeAccountX(bz []byte) (ax AccountX) {
	err := axk.cdc.UnmarshalBinaryBare(bz, &ax)
	if err != nil {
		panic(err)
	}
	return
}
