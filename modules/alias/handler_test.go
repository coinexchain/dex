package alias

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	dex "github.com/coinexchain/dex/types"
)

var logStr string

type mocBankxKeeper struct {
	maxAmount sdk.Int
}

func (k *mocBankxKeeper) DeductInt64CetFee(ctx sdk.Context, addr sdk.AccAddress, amt int64) sdk.Error {
	return k.DeductFee(ctx, addr, dex.NewCetCoins(amt))
}

func (k *mocBankxKeeper) DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	coinStrList := make([]string, len(amt))
	for i, coin := range amt {
		if coin.Amount.GT(k.maxAmount) {
			return sdk.NewError(types.CodeSpaceAlias, 1199, "Not enough coins")
		}
		coinStrList[i] = coin.Amount.String() + coin.Denom
	}
	logStr = "Deduct " + strings.Join(coinStrList, ",") + " from " + addr.String()
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
			maxAmount: sdk.NewInt(913000000000),
		},
		&mocAssetKeeper{
			tokenIssuer: map[string]sdk.AccAddress{"cet": simpleAddr("00000")},
		},
		paramsKeeper.Subspace(types.StoreKey),
	)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())
	parameters := types.DefaultParams()
	keeper.SetParams(ctx, parameters)

	return ctx, &keeper
}

func simpleAddr(s string) sdk.AccAddress {
	a, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
	return a
}

func Test1(t *testing.T) {
	codec.RunInitFuncList()
	tom := simpleAddr("00001")
	bob := simpleAddr("00002")
	alice := simpleAddr("00003")
	fmt.Printf("tom: %s\n", tom.String())
	fmt.Printf("bob: %s\n", bob.String())
	fmt.Printf("alice: %s\n", alice.String())
	ctx, keeper := newContextAndKeeper("test-1")
	InitGenesis(ctx, *keeper, DefaultGenesisState())
	aliasEntryList := keeper.GetAllAlias(ctx)
	require.Equal(t, 0, len(aliasEntryList))

	refList := []keepers.AliasEntry{
		{Alias: "alice", Addr: alice, AsDefault: false},
		{Alias: "tom", Addr: tom, AsDefault: false},
		{Alias: "tom@gmail.com", Addr: tom, AsDefault: false},
	}
	genS := NewGenesisState(types.DefaultParams(), refList)
	InitGenesis(ctx, *keeper, genS)

	aliasEntryList = ExportGenesis(ctx, *keeper).AliasEntryList
	require.Equal(t, refList, aliasEntryList)

	err := NewGenesisState(types.DefaultParams(), []keepers.AliasEntry{
		{Alias: "alice", Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, nil, err)
	err = NewGenesisState(types.DefaultParams(), []keepers.AliasEntry{
		{Alias: "爱丽丝", Addr: alice, AsDefault: false},
	}).Validate()
	refErr := errors.New("Invalid Alias")
	require.Equal(t, refErr, err)
	err = NewGenesisState(types.DefaultParams(), []keepers.AliasEntry{
		{Alias: string([]byte{255, 255}), Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)
	err = NewGenesisState(types.DefaultParams(), []keepers.AliasEntry{
		{Alias: "", Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)
	err = NewGenesisState(types.DefaultParams(), []keepers.AliasEntry{
		{Alias: "a", Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)
	err = NewGenesisState(types.DefaultParams(), []keepers.AliasEntry{
		{Alias: "abcdefghijklmnopqrstuvwxyabcdefghijklmnopqrstuvwxyabcabcdefghijklmnopqrstuvwxyabcdefghijklmnopqrstuvwxyabcdefghijklmnopqrstuvwxyzzzabcdefghijklmnopqrstuvwxyz", Addr: alice, AsDefault: false},
	}).Validate()
	require.Equal(t, refErr, err)

	params := types.DefaultParams()
	params.MaxAliasCount = 0
	err = NewGenesisState(params, []keepers.AliasEntry{
		{Alias: string([]byte{255, 255}), Addr: alice, AsDefault: false},
	}).Validate()
	refErr = fmt.Errorf("%s must be a positive number, is %d", types.KeyMaxAliasCount, 0)
	require.Equal(t, refErr, err)

	handlerFunc := NewHandler(*keeper)

	handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "supergirl", IsAdd: true, AsDefault: false})
	refLog := "Deduct 1000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqrctskts"
	require.Equal(t, refLog, logStr)
	handlerFunc(ctx, types.MsgAliasUpdate{Owner: bob, Alias: "superman", IsAdd: true, AsDefault: false})
	refLog = "Deduct 1000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz"
	require.Equal(t, refLog, logStr)
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
	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: "superman", IsAdd: true, AsDefault: true})
	require.Equal(t, types.ErrAliasAlreadyExists().Result(), res)
	aliasCet := "coinex-usdt.t"
	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: aliasCet, IsAdd: true})
	require.Equal(t, types.ErrCanOnlyBeUsedByCetOwner(aliasCet).Result(), res)

	cetOwner := simpleAddr("00000")
	refLog = "Deduct 1000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqqkc9q90"
	handlerFunc(ctx, types.MsgAliasUpdate{Owner: cetOwner, Alias: aliasCet, IsAdd: true, AsDefault: true})
	addr, isDefault := keeper.GetAddressFromAlias(ctx, aliasCet)
	require.Equal(t, cetOwner, sdk.AccAddress(addr))
	require.Equal(t, true, isDefault)
	require.Equal(t, refLog, logStr)

	handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "sup", IsAdd: true, AsDefault: false})
	refLog = "Deduct 500000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqrctskts"
	require.Equal(t, refLog, logStr)
	handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "supe", IsAdd: true, AsDefault: false})
	refLog = "Deduct 200000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqrctskts"
	require.Equal(t, refLog, logStr)
	handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "super", IsAdd: true, AsDefault: false})
	refLog = "Deduct 100000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqrctskts"
	require.Equal(t, refLog, logStr)
	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: alice, Alias: "superg", IsAdd: true, AsDefault: false})
	require.Equal(t, types.ErrMaxAliasCountReached().Result(), res)

	handlerFunc(ctx, types.MsgAliasUpdate{Owner: bob, Alias: "12345", IsAdd: true, AsDefault: false})
	refLog = "Deduct 100000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz"
	require.Equal(t, refLog, logStr)
	handlerFunc(ctx, types.MsgAliasUpdate{Owner: bob, Alias: "123456", IsAdd: true, AsDefault: false})
	refLog = "Deduct 10000000000cet from cosmos1qy352eufqy352eufqy352eufqy35qqqz9ayrkz"
	require.Equal(t, refLog, logStr)
	res = handlerFunc(ctx, types.MsgAliasUpdate{Owner: bob, Alias: "12", IsAdd: true, AsDefault: false})
	require.Equal(t, sdk.NewError(types.CodeSpaceAlias, 1199, "Not enough coins").Result(), res)

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

func TestReservedAliases(t *testing.T) {
	codec.RunInitFuncList()
	ctx, keeper := newContextAndKeeper("test-1")
	handlerFunc := NewHandler(*keeper)
	tom := simpleAddr("00001")

	reservedAliases := []string{
		"coinex",
		"cet",
		"viabtc",
		"cetdac",
		"coinex-usdt.t",
		"usdt.t-coinex",
		"usdt-coinex",
		"cet.coinex",
		"coinex.cet",
		"coinex__",
		"__coinex",
		"www.coinex.com",
		"www.coinex.org",
		"www.coinex.net",
		"coinex.com",
		"coinex.org",
		"coinex.net",
		"cetdac@coinex.org",
		"bob@coinex.com",
		"btc@coinex.com",
		"bob.mail@coinex.com",
	}

	for _, alias := range reservedAliases {
		res := handlerFunc(ctx, types.MsgAliasUpdate{Owner: tom, Alias: alias, IsAdd: true})
		require.Equal(t, types.ErrCanOnlyBeUsedByCetOwner(alias).Result(), res)
	}
}
