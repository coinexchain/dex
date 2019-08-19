package keepers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/alias/internal/types"
)

var (
	AliasToAccountKey    = []byte{0x10}
	AliasToAccountKeyEnd = []byte{0x11}
	AccountToAliasKey    = []byte{0x12}
)

type AliasEntry struct {
	Alias     string         `json:"alias"`
	Addr      sdk.AccAddress `json:"addr"`
	AsDefault bool           `json:"is_default"`
}

type AliasKeeper struct {
	aliasKey sdk.StoreKey
}

func NewAliasKeeper(key sdk.StoreKey) *AliasKeeper {
	return &AliasKeeper{
		aliasKey: key,
	}
}

func (keeper *AliasKeeper) GetAddressFromAlias(ctx sdk.Context, alias string) ([]byte, bool) {
	store := ctx.KVStore(keeper.aliasKey)
	key := append(AliasToAccountKey, []byte(alias)...)
	value := store.Get(key)
	if len(value) == 0 {
		return nil, false
	}
	return value[1:], value[0] != 0
}

func getAccountToAliasKey(addr sdk.AccAddress, alias []byte) []byte {
	return append(append(append(AccountToAliasKey, []byte{byte(len(addr))}...), addr...), alias...)
}

func (keeper *AliasKeeper) GetAliasListOfAccount(ctx sdk.Context, addr sdk.AccAddress) []string {
	store := ctx.KVStore(keeper.aliasKey)
	keyStart := getAccountToAliasKey(addr, []byte{0})
	keyEnd := getAccountToAliasKey(addr, []byte{255})
	defaultAlias := ""
	aliasList := make([]string, 1, 5)
	start := len(keyStart) - 1
	iter := store.Iterator(keyStart, keyEnd)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		alias := iter.Key()[start:]
		asDefault := iter.Value()
		if asDefault[0] == 0 {
			aliasList = append(aliasList, string(alias))
		} else {
			defaultAlias = string(alias)
		}
	}
	aliasList[0] = defaultAlias
	return aliasList
}

func (keeper *AliasKeeper) GetAllAlias(ctx sdk.Context) []AliasEntry {
	store := ctx.KVStore(keeper.aliasKey)
	iter := store.Iterator(AliasToAccountKey, AliasToAccountKeyEnd)
	defer iter.Close()
	res := make([]AliasEntry, 0, 1000)
	for ; iter.Valid(); iter.Next() {
		res = append(res, AliasEntry{
			Alias:     string(iter.Key()[1:]),
			Addr:      sdk.AccAddress(iter.Value()[1:]),
			AsDefault: iter.Value()[0] != 0,
		})
	}
	return res
}

func (keeper *AliasKeeper) AddAlias(ctx sdk.Context, alias string, addr sdk.AccAddress, asDefault bool, maxCount int) (bool, bool) {
	aliasList := keeper.GetAliasListOfAccount(ctx, addr)
	hasDefault := aliasList[0] != ""
	if !hasDefault {
		aliasList = aliasList[1:]
	}
	addNewAlias := true
	for _, a := range aliasList {
		if a == alias {
			addNewAlias = false
		}
	}
	if addNewAlias && len(aliasList) >= maxCount && maxCount > 0 {
		return false, addNewAlias
	}
	if asDefault && hasDefault {
		keeper.setAlias(ctx, aliasList[0], addr, false)
	}
	keeper.setAlias(ctx, alias, addr, asDefault)
	return true, addNewAlias
}

func (keeper *AliasKeeper) setAlias(ctx sdk.Context, alias string, addr sdk.AccAddress, asDefault bool) {
	store := ctx.KVStore(keeper.aliasKey)
	key := append(AliasToAccountKey, []byte(alias)...)
	asDefaultAliasByte := []byte{0}
	if asDefault {
		asDefaultAliasByte = []byte{1}
	}
	store.Set(key, append(asDefaultAliasByte, addr...))
	key = getAccountToAliasKey(addr, []byte(alias))
	store.Set(key, asDefaultAliasByte)
}

func (keeper *AliasKeeper) RemoveAlias(ctx sdk.Context, alias string, addr sdk.AccAddress) {
	store := ctx.KVStore(keeper.aliasKey)
	key := append(AliasToAccountKey, []byte(alias)...)
	store.Delete(key)
	key = getAccountToAliasKey(addr, []byte(alias))
	store.Delete(key)
}

//============================================================================

type Keeper struct {
	paramSubspace params.Subspace
	AliasKeeper   *AliasKeeper
	BankKeeper    types.ExpectedBankxKeeper
	AssetKeeper   types.ExpectedAssetStatusKeeper
}

func NewKeeper(key sdk.StoreKey,
	bankKeeper types.ExpectedBankxKeeper,
	assetKeeper types.ExpectedAssetStatusKeeper,
	paramstore params.Subspace) Keeper {

	return Keeper{
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
		AliasKeeper:   NewAliasKeeper(key),
		BankKeeper:    bankKeeper,
		AssetKeeper:   assetKeeper,
	}
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}
