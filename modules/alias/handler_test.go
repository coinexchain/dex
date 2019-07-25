package alias

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

//// Bankx Keeper will implement the interface
//type ExpectedBankxKeeper interface {
//	DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
//}
//// Asset Keeper will implement the interface
//type ExpectedAssetStatusKeeper interface {
//	IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
//}

var logStrList = make([]string, 0, 100)

func logStrClear() {
	logStrList = logStrList[:0]
}

func logStrAppend(s string) {
	logStrList = append(logStrList, s)
}

type mocBankxKeeper struct {
	maxAmount sdk.Int
}

func (k *mocBankxKeeper) DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	coinStrList := make([]string, len(amt))
	for i, coin := range amt {
		if coin.Amount.GT(k.maxAmount) {
			return sdk.NewError(types.CodeSpaceAlias, 1199, "Not enough coins")
		}
		coinStrList[i] = coin.Amount.String() + coin.Denom
	}
	s := "Deduct " + strings.Join(coinStrList, ",") + " from " + addr.String()
	logStrAppend(s)
	return nil
}

type mocAssetKeeper struct {
	tokenIssuer map[string]sdk.AccAddress
}

func (k *mocAssetKeeper) IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	a, ok := k.tokenIssuer[denom]
	return ok && bytes.Equal(addr, a)
}

func newContextAndKeeper(chainid string) (sdk.Context, *Keeper) {
	db := dbm.NewMemDB()
	ms := sdkstore.NewCommitMultiStore(db)

	key := sdk.NewKVStoreKey(types.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	paramsKeeper := params.NewKeeper(types.ModuleCdc, keyParams, tkeyParams, params.DefaultCodespace)
	keeper := keepers.NewKeeper(key,
		&mocBankxKeeper{
			maxAmount: sdk.NewInt(1000),
		},
		&mocAssetKeeper{
			tokenIssuer: map[string]sdk.AccAddress{"cet": simpleAddr("00000")},
		},
		paramsKeeper.Subspace(types.StoreKey),
	)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())
	parameters := keepers.DefaultParams()
	keeper.SetParams(ctx, parameters)

	return ctx, &keeper
}

func simpleAddr(s string) sdk.AccAddress {
	a, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
	return a
}

func Test1(t *testing.T) {
	tom := simpleAddr("00001")
	bob := simpleAddr("00002")
	alice := simpleAddr("00003")
	fmt.Printf("tom: %s\n", tom.String())
	fmt.Printf("bob: %s\n", bob.String())
	fmt.Printf("alice: %s\n", alice.String())
	ctx, keeper := newContextAndKeeper("test-1")
	InitGenesis(ctx, *keeper, DefaultGenesisState())
	aliasEntryList := keeper.AliasKeeper.GetAllAlias(ctx)
	require.Equal(t, 0, len(aliasEntryList))

	genS := NewGenesisState(keepers.DefaultParams(), []keepers.AliasEntry{
		{Alias: "alice", Addr: alice, AsDefault: false},
		{Alias: "tom", Addr: tom, AsDefault: false},
	})
	InitGenesis(ctx, *keeper, genS)

	ak := keeper.AliasKeeper
	maxCount := 5
	ak.AddAlias(ctx, "goodgirl", alice, false, maxCount)
	ak.AddAlias(ctx, "tom@gmail.com", tom, false, maxCount)
	ak.AddAlias(ctx, "bob", bob, false, maxCount)
	addr, isDefault := ak.GetAddressFromAlias(ctx, "alice")
	require.Equal(t, alice, sdk.AccAddress(addr))
	require.Equal(t, false, isDefault)
	addr, isDefault = ak.GetAddressFromAlias(ctx, "goodgirl")
	require.Equal(t, alice, sdk.AccAddress(addr))
	require.Equal(t, false, isDefault)
	addr, isDefault = ak.GetAddressFromAlias(ctx, "tom")
	require.Equal(t, tom, sdk.AccAddress(addr))
	require.Equal(t, false, isDefault)
	addr, isDefault = ak.GetAddressFromAlias(ctx, "tom@gmail.com")
	require.Equal(t, tom, sdk.AccAddress(addr))
	require.Equal(t, false, isDefault)
	addr, isDefault = ak.GetAddressFromAlias(ctx, "bob")
	require.Equal(t, bob, sdk.AccAddress(addr))
	require.Equal(t, false, isDefault)

	aliasList := ak.GetAliasListOfAccount(ctx, alice)
	require.Equal(t, []string{"", "alice", "goodgirl"}, aliasList)
	aliasList = ak.GetAliasListOfAccount(ctx, tom)
	require.Equal(t, []string{"", "tom", "tom@gmail.com"}, aliasList)
	aliasList = ak.GetAliasListOfAccount(ctx, bob)
	require.Equal(t, []string{"", "bob"}, aliasList)

	aliasEntryList = ak.GetAllAlias(ctx)
	refList := []keepers.AliasEntry{
		{Alias: "alice", Addr: alice, AsDefault: false},
		{Alias: "bob", Addr: bob, AsDefault: false},
		{Alias: "goodgirl", Addr: alice, AsDefault: false},
		{Alias: "tom", Addr: tom, AsDefault: false},
		{Alias: "tom@gmail.com", Addr: tom, AsDefault: false},
	}
	require.Equal(t, refList, aliasEntryList)

	ak.RemoveAlias(ctx, "goodgirl", alice)
	ak.RemoveAlias(ctx, "bob", bob)
	aliasList = ak.GetAliasListOfAccount(ctx, alice)
	require.Equal(t, []string{"", "alice"}, aliasList)
	aliasList = ak.GetAliasListOfAccount(ctx, bob)
	require.Equal(t, []string{""}, aliasList)

	aliasEntryList = ExportGenesis(ctx, *keeper).AliasEntryList
	refList = []keepers.AliasEntry{
		{Alias: "alice", Addr: alice, AsDefault: false},
		{Alias: "tom", Addr: tom, AsDefault: false},
		{Alias: "tom@gmail.com", Addr: tom, AsDefault: false},
	}
	require.Equal(t, refList, aliasEntryList)

	err := NewGenesisState(keepers.DefaultParams(), []keepers.AliasEntry{
		{Alias: "爱丽丝", Addr: alice, AsDefault: false},
	}).Validate()
	refErr := errors.New("Invalid Alias")
	require.Equal(t, refErr, err)
	err = NewGenesisState(keepers.DefaultParams(), []keepers.AliasEntry{
		{Alias: string([]byte{255, 255}), Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)
	err = NewGenesisState(keepers.DefaultParams(), []keepers.AliasEntry{
		{Alias: "", Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)
	err = NewGenesisState(keepers.DefaultParams(), []keepers.AliasEntry{
		{Alias: "a", Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)
	err = NewGenesisState(keepers.DefaultParams(), []keepers.AliasEntry{
		{Alias: "abcdefghijklmnopqrstuvwxyabcdefghijklmnopqrstuvwxyabcabcdefghijklmnopqrstuvwxyabcdefghijklmnopqrstuvwxyabcdefghijklmnopqrstuvwxyzzzabcdefghijklmnopqrstuvwxyz", Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)

	handlerFunc := NewHandler(*keeper)

	handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "supergirl", IsAdd: true, AsDefault: false})
	handlerFunc(ctx, types.MsgAliasUpdate{Owner: bob, Alias: "superman", IsAdd: true, AsDefault: false})
	handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: "tom", IsAdd: false, AsDefault: false})

	aliasEntryList = ExportGenesis(ctx, *keeper).AliasEntryList
	refList = []keepers.AliasEntry{
		{Alias: "alice", Addr: alice, AsDefault: false},
		{Alias: "supergirl", Addr: alice, AsDefault: false},
		{Alias: "superman", Addr: bob, AsDefault: false},
		{Alias: "tom@gmail.com", Addr: tom, AsDefault: false},
	}
	require.Equal(t, refList, aliasEntryList)

	res := handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: "tom", IsAdd: false})
	require.Equal(t, types.ErrNoSuchAlias().Result(), res)
	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "superman", IsAdd: false})
	require.Equal(t, types.ErrNoSuchAlias().Result(), res)
	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: "superman", IsAdd: true})
	require.Equal(t, types.ErrAliasAlreadyExists().Result(), res)

	msg := types.MsgAliasUpdate{Owner: []byte{}, Alias: "supergirl", IsAdd: true}
	err = msg.ValidateBasic()
	require.Equal(t, sdk.ErrInvalidAddress("missing owner address"), err)
	msg = types.MsgAliasUpdate{Owner: tom, Alias: "", IsAdd: true}
	err = msg.ValidateBasic()
	require.Equal(t, types.ErrEmptyAlias(), err)
	msg = types.MsgAliasUpdate{Owner: tom, Alias: "I Love U", IsAdd: true}
	err = msg.ValidateBasic()
	require.Equal(t, types.ErrInvalidAlias(), err)
}
