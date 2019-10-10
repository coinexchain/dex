package keepers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/alias/internal/types"
)

func newContextAndKeeper(chainid string) (sdk.Context, *AliasKeeper) {
	db := dbm.NewMemDB()
	ms := sdkstore.NewCommitMultiStore(db)

	key := sdk.NewKVStoreKey(types.StoreKey)

	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	keeper := NewAliasKeeper(key)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())

	return ctx, keeper
}

func simpleAddr(s string) sdk.AccAddress {
	a, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
	return a
}

func TestGetAccountToAliasKey(t *testing.T) {
	key := getAccountToAliasKey([]byte("addr"), []byte("alias"))
	require.Equal(t, "\x12\x04addralias", string(key))
}

func Test1(t *testing.T) {
	tom := simpleAddr("00001")
	bob := simpleAddr("00002")
	alice := simpleAddr("00003")
	fmt.Printf("tom: %s\n", tom.String())
	fmt.Printf("bob: %s\n", bob.String())
	fmt.Printf("alice: %s\n", alice.String())
	ctx, ak := newContextAndKeeper("test-1")

	maxCount := 5
	ak.AddAlias(ctx, "alice", alice, false, maxCount)
	ak.AddAlias(ctx, "tom", tom, false, maxCount)
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
	addr, isDefault = ak.GetAddressFromAlias(ctx, "hehe")
	require.Equal(t, []byte(nil), addr)
	require.Equal(t, false, isDefault)

	aliasList := ak.GetAliasListOfAccount(ctx, alice)
	require.Equal(t, []string{"", "alice", "goodgirl"}, aliasList)
	aliasList = ak.GetAliasListOfAccount(ctx, tom)
	require.Equal(t, []string{"", "tom", "tom@gmail.com"}, aliasList)
	aliasList = ak.GetAliasListOfAccount(ctx, bob)
	require.Equal(t, []string{"", "bob"}, aliasList)

	aliasEntryList := ak.GetAllAlias(ctx)
	refList := []AliasEntry{
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

	ak.AddAlias(ctx, "bob", bob, true, maxCount)
	ak.AddAlias(ctx, "goodgirl", alice, true, maxCount)
	addr, isDefault = ak.GetAddressFromAlias(ctx, "bob")
	require.Equal(t, bob, sdk.AccAddress(addr))
	require.Equal(t, true, isDefault)
	addr, isDefault = ak.GetAddressFromAlias(ctx, "goodgirl")
	require.Equal(t, alice, sdk.AccAddress(addr))
	require.Equal(t, true, isDefault)
	aliasList = ak.GetAliasListOfAccount(ctx, alice)
	require.Equal(t, []string{"goodgirl", "alice"}, aliasList)
	aliasList = ak.GetAliasListOfAccount(ctx, bob)
	require.Equal(t, []string{"bob"}, aliasList)
	ak.AddAlias(ctx, "alice", alice, true, maxCount)
	aliasList = ak.GetAliasListOfAccount(ctx, alice)
	require.Equal(t, []string{"alice", "goodgirl"}, aliasList)
	ok, addNewAlias := ak.AddAlias(ctx, "supergirl", alice, true, 2)
	require.Equal(t, false, ok)
	require.Equal(t, true, addNewAlias)
}
