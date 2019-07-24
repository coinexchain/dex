package alias

//import (
//	"errors"
//	"fmt"
//	"github.com/stretchr/testify/require"
//	"testing"
//
//	"github.com/coinexchain/dex/modules/alias/internal/keepers"
//	"github.com/coinexchain/dex/modules/alias/internal/types"
//
//	sdkstore "github.com/cosmos/cosmos-sdk/store"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	abci "github.com/tendermint/tendermint/abci/types"
//	dbm "github.com/tendermint/tendermint/libs/db"
//	"github.com/tendermint/tendermint/libs/log"
//)
//
//func newContextAndKeeper(chainid string) (sdk.Context, *Keeper) {
//	db := dbm.NewMemDB()
//	ms := sdkstore.NewCommitMultiStore(db)
//
//	key := sdk.NewKVStoreKey(StoreKey)
//	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
//	ms.LoadLatestVersion()
//
//	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())
//	k := keepers.NewKeeper(key)
//	return ctx, k
//}
//
//func simpleAddr(s string) sdk.AccAddress {
//	a, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
//	return a
//}
//
//func Test1(t *testing.T) {
//	tom := simpleAddr("00001")
//	bob := simpleAddr("00002")
//	alice := simpleAddr("00003")
//	fmt.Printf("Tom: %s\n", tom.String())
//	fmt.Printf("Bob: %s\n", bob.String())
//	fmt.Printf("Alice: %s\n", alice.String())
//	ctx, keeper := newContextAndKeeper("test-1")
//	InitGenesis(ctx, *keeper, DefaultGenesisState())
//	aliasMap := keeper.GetAllAlias(ctx)
//	require.Equal(t, len(aliasMap), 0)
//
//	genS := NewGenesisState(map[string]sdk.AccAddress{
//		"Alice": alice,
//		"Tom":   tom,
//	})
//	InitGenesis(ctx, *keeper, genS)
//
//	keeper.AddAlias(ctx, "GoodGirl", alice)
//	keeper.AddAlias(ctx, "汤姆", tom)
//	keeper.AddAlias(ctx, "Bob", bob)
//	addr := sdk.AccAddress(keeper.GetAddressFromAlias(ctx, "Alice"))
//	require.Equal(t, addr, alice)
//	addr = sdk.AccAddress(keeper.GetAddressFromAlias(ctx, "GoodGirl"))
//	require.Equal(t, addr, alice)
//	addr = sdk.AccAddress(keeper.GetAddressFromAlias(ctx, "Tom"))
//	require.Equal(t, addr, tom)
//	addr = sdk.AccAddress(keeper.GetAddressFromAlias(ctx, "汤姆"))
//	require.Equal(t, addr, tom)
//	addr = sdk.AccAddress(keeper.GetAddressFromAlias(ctx, "Bob"))
//	require.Equal(t, addr, bob)
//
//	aliasList := keeper.GetAliasListOfAccount(ctx, alice)
//	require.Equal(t, aliasList, []string{"Alice", "GoodGirl"})
//	aliasList = keeper.GetAliasListOfAccount(ctx, tom)
//	require.Equal(t, aliasList, []string{"Tom", "汤姆"})
//	aliasList = keeper.GetAliasListOfAccount(ctx, bob)
//	require.Equal(t, aliasList, []string{"Bob"})
//
//	aliasMap = keeper.GetAllAlias(ctx)
//	refMap := map[string]sdk.AccAddress{
//		"Alice":    alice,
//		"GoodGirl": alice,
//		"Tom":      tom,
//		"汤姆":       tom,
//		"Bob":      bob,
//	}
//	require.Equal(t, aliasMap, refMap)
//
//	keeper.RemoveAlias(ctx, "GoodGirl", alice)
//	keeper.RemoveAlias(ctx, "Bob", bob)
//	aliasList = keeper.GetAliasListOfAccount(ctx, alice)
//	require.Equal(t, aliasList, []string{"Alice"})
//	aliasList = keeper.GetAliasListOfAccount(ctx, bob)
//	require.Equal(t, 0, len(aliasList))
//
//	aliasMap = ExportGenesis(ctx, *keeper).AliasInfoMap
//	refMap = map[string]sdk.AccAddress{
//		"Alice": alice,
//		"Tom":   tom,
//		"汤姆":    tom,
//	}
//	require.Equal(t, aliasMap, refMap)
//
//	err := NewGenesisState(map[string]sdk.AccAddress{
//		"0lice": alice,
//	}).Validate()
//	refErr := errors.New("Invalid Alias")
//	require.Equal(t, err, refErr)
//	err = NewGenesisState(map[string]sdk.AccAddress{
//		string([]byte{255, 255}): alice,
//	}).Validate()
//	require.Equal(t, err, refErr)
//	err = NewGenesisState(map[string]sdk.AccAddress{
//		"": alice,
//	}).Validate()
//	require.Equal(t, err, refErr)
//	err = NewGenesisState(map[string]sdk.AccAddress{
//		"Love You": alice,
//	}).Validate()
//	require.Equal(t, err, refErr)
//	err = NewGenesisState(map[string]sdk.AccAddress{
//		"ABCDEFGHIJKLMNOPQRSTUVWXYZ": alice,
//	}).Validate()
//	require.Equal(t, err, refErr)
//
//	handlerFunc := NewHandler(*keeper)
//
//	handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "SuperGirl", IsAdd: true})
//	handlerFunc(ctx, types.MsgAliasUpdate{Owner: bob, Alias: "Superman", IsAdd: true})
//	handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: "Tom", IsAdd: false})
//
//	aliasMap = ExportGenesis(ctx, *keeper).AliasInfoMap
//	refMap = map[string]sdk.AccAddress{
//		"Alice":     alice,
//		"SuperGirl": alice,
//		"Superman":  bob,
//		"汤姆":        tom,
//	}
//	require.Equal(t, aliasMap, refMap)
//
//	res := handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: "Tom", IsAdd: false})
//	require.Equal(t, res, types.ErrNoSuchAlias().Result())
//	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "Superman", IsAdd: false})
//	require.Equal(t, res, types.ErrNoSuchAlias().Result())
//	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: "Superman", IsAdd: true})
//	require.Equal(t, res, types.ErrAliasAlreadyExists().Result())
//
//	msg := types.MsgAliasUpdate{Owner: []byte{}, Alias: "SuperGirl", IsAdd: true}
//	err = msg.ValidateBasic()
//	require.Equal(t, err, sdk.ErrInvalidAddress("missing owner address"))
//	msg = types.MsgAliasUpdate{Owner: tom, Alias: "", IsAdd: true}
//	err = msg.ValidateBasic()
//	require.Equal(t, err, types.ErrEmptyAlias())
//	msg = types.MsgAliasUpdate{Owner: tom, Alias: "I Love U", IsAdd: true}
//	err = msg.ValidateBasic()
//	require.Equal(t, err, types.ErrInvalidAlias())
//}
