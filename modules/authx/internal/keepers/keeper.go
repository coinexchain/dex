package keepers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/authx/internal/types"
)

var (
	// AddressStoreKeyPrefix prefix for accountx-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	PrefixUnlockedCoinsQueue = []byte("UnlockedCoinsQueue")
	KeyDelimiter             = []byte(";")
)

type AccountXKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace params.Subspace

	supplyKeeper SupplyKeeper

	ak ExpectedAccountKeeper

	EventTypeMsgQueue string
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSubspace params.Subspace, keeper SupplyKeeper, ak ExpectedAccountKeeper, eventTypeMsgQueue string) AccountXKeeper {
	// ensure authx module account is set
	if addr := keeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return AccountXKeeper{
		key:               key,
		cdc:               cdc,
		paramSubspace:     paramSubspace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:      keeper,
		ak:                ak,
		EventTypeMsgQueue: eventTypeMsgQueue,
	}
}

func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}

// -----------------------------------------------------------------------------
// AccountX

func (axk AccountXKeeper) GetOrCreateAccountX(ctx sdk.Context, addr sdk.AccAddress) types.AccountX {
	ax, ok := axk.GetAccountX(ctx, addr)
	if !ok {
		ax = types.AccountX{Address: addr}
		axk.SetAccountX(ctx, ax)
	}
	return ax
}

func (axk AccountXKeeper) GetAccountX(ctx sdk.Context, addr sdk.AccAddress) (ax types.AccountX, ok bool) {
	store := ctx.KVStore(axk.key)
	bz := store.Get(AddressStoreKey(addr))
	if bz == nil {
		return
	}

	acc := axk.decodeAccountX(bz)
	return acc, true
}

func (axk AccountXKeeper) SetAccountX(ctx sdk.Context, ax types.AccountX) {
	addr := ax.Address
	store := ctx.KVStore(axk.key)
	bz, err := axk.cdc.MarshalBinaryBare(ax)
	if err != nil {
		panic(err)
	}
	store.Set(AddressStoreKey(addr), bz)
}

func (axk AccountXKeeper) IterateAccounts(ctx sdk.Context, process func(types.AccountX) (stop bool)) {
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

func (axk AccountXKeeper) UnlockedCoinsQueueIterator(ctx sdk.Context, unlockedTime int64) sdk.Iterator {
	store := ctx.KVStore(axk.key)
	return store.Iterator(PrefixUnlockedCoinsQueue, sdk.PrefixEndBytes(PrefixUnlockedTimeQueueTime(unlockedTime)))
}

func (axk AccountXKeeper) InsertUnlockedCoinsQueue(ctx sdk.Context, unlockedTime int64, address sdk.AccAddress) {
	store := ctx.KVStore(axk.key)
	store.Set(KeyUnlockedCoinsQueue(unlockedTime, address), address)
}

func (axk AccountXKeeper) RemoveFromUnlockedCoinsQueue(ctx sdk.Context, unlockedTime int64, address sdk.AccAddress) {
	store := ctx.KVStore(axk.key)
	store.Delete(KeyUnlockedCoinsQueue(unlockedTime, address))
}

func (axk AccountXKeeper) RemoveFromUnlockedCoinsQueueByKey(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(axk.key)
	store.Delete(key)
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (axk AccountXKeeper) SetParams(ctx sdk.Context, params types.Params) {
	axk.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (axk AccountXKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	axk.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Codec

func (axk AccountXKeeper) decodeAccountX(bz []byte) (ax types.AccountX) {
	err := axk.cdc.UnmarshalBinaryBare(bz, &ax)

	if err != nil {
		panic(err)
	}
	return
}

// -----------------------------------------------------------------------------
// PreTotalSupply sets the Authx Module Account
func (axk AccountXKeeper) PreTotalSupply(ctx sdk.Context) {
	var expectedTotal sdk.Coins

	axk.IterateAccounts(ctx, func(acc types.AccountX) bool {
		expectedTotal = expectedTotal.Add(acc.GetAllCoins())
		return false
	})

	authxMacc := axk.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	_ = authxMacc.SetCoins(expectedTotal)
	axk.supplyKeeper.SetModuleAccount(ctx, authxMacc)
}

// -----------------------------------------------------------------------------
// Keys

func KeyUnlockedCoinsQueue(unlockedTime int64, address sdk.AccAddress) []byte {
	return bytes.Join([][]byte{
		PrefixUnlockedCoinsQueue,
		sdk.FormatTimeBytes(time.Unix(unlockedTime, 0)),
		address,
	}, KeyDelimiter)
}

func PrefixUnlockedTimeQueueTime(unlockedTime int64) []byte {
	return bytes.Join([][]byte{
		PrefixUnlockedCoinsQueue,
		sdk.FormatTimeBytes(time.Unix(unlockedTime, 0)),
	}, KeyDelimiter)
}
