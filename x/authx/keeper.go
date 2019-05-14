package authx

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// StoreKey is string representation of the store key for authx
	StoreKey = "accx"
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
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) AccountXKeeper {
	return AccountXKeeper{
		key: key,
		cdc: cdc,
	}
}

func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
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

// -----------------------------------------------------------------------------
// Codec

func (axk AccountXKeeper) decodeAccountX(bz []byte) (ax AccountX) {
	err := axk.cdc.UnmarshalBinaryBare(bz, &ax)
	if err != nil {
		panic(err)
	}
	return
}
