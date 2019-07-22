package keepers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AliasToAccountKey    = []byte{0x10}
	AliasToAccountKeyEnd = []byte{0x11}
	AccountToAliasKey    = []byte{0x12}
)

type Keeper struct {
	aliasKey sdk.StoreKey
}

func NewKeeper(key sdk.StoreKey) *Keeper {
	return &Keeper{
		aliasKey: key,
	}
}

func (keeper *Keeper) GetAddressFromAlias(ctx sdk.Context, alias string) []byte {
	store := ctx.KVStore(keeper.aliasKey)
	key := append(AliasToAccountKey, []byte(alias)...)
	return store.Get(key)
}

func getAccountToAliasKey(addr sdk.AccAddress, alias []byte) []byte {
	return append(append(append(AccountToAliasKey, []byte{byte(len(addr))}...), addr...), alias...)
}

func (keeper *Keeper) GetAliasListOfAccount(ctx sdk.Context, addr sdk.AccAddress) []string {
	store := ctx.KVStore(keeper.aliasKey)
	keyStart := getAccountToAliasKey(addr, []byte{0})
	keyEnd := getAccountToAliasKey(addr, []byte{255})
	iter := store.Iterator(keyStart, keyEnd)
	defer iter.Close()
	res := make([]string, 0, 10)
	start := len(keyStart) - 1
	for ; iter.Valid(); iter.Next() {
		alias := iter.Key()[start:]
		res = append(res, string(alias))
	}
	return res
}

func (keeper *Keeper) GetAllAlias(ctx sdk.Context) map[string]sdk.AccAddress {
	store := ctx.KVStore(keeper.aliasKey)
	iter := store.Iterator(AliasToAccountKey, AliasToAccountKeyEnd)
	defer iter.Close()
	res := make(map[string]sdk.AccAddress)
	for ; iter.Valid(); iter.Next() {
		alias := string(iter.Key()[1:])
		res[alias] = sdk.AccAddress(iter.Value())
	}
	return res
}

func (keeper *Keeper) AddAlias(ctx sdk.Context, alias string, addr sdk.AccAddress) {
	store := ctx.KVStore(keeper.aliasKey)
	key := append(AliasToAccountKey, []byte(alias)...)
	store.Set(key, addr)
	key = getAccountToAliasKey(addr, []byte(alias))
	store.Set(key, []byte{})
}

func (keeper *Keeper) RemoveAlias(ctx sdk.Context, alias string, addr sdk.AccAddress) {
	store := ctx.KVStore(keeper.aliasKey)
	key := append(AliasToAccountKey, []byte(alias)...)
	store.Delete(key)
	key = getAccountToAliasKey(addr, []byte(alias))
	store.Delete(key)
}
