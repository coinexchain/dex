package bancorlite

import (
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	dex "github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	ctx     sdk.Context
	bik     keepers.Keeper
	handler sdk.Handler
	akp     auth.AccountKeeper
	keys    storeKeys
	cdc     *codec.Codec // mk.cdc
}

type storeKeys struct {
	assetCapKey *sdk.KVStoreKey
	authCapKey  *sdk.KVStoreKey
	authxCapKey *sdk.KVStoreKey
	keyParams   *sdk.KVStoreKey
	tkeyParams  *sdk.TransientStoreKey
	marketKey   *sdk.KVStoreKey
	authxKey    *sdk.KVStoreKey
	keyStaking  *sdk.KVStoreKey
	tkeyStaking *sdk.TransientStoreKey
	keySupply   *sdk.KVStoreKey
	keyBancor   *sdk.KVStoreKey
}

var (
	haveCetAddress            = getAddr("000001")
	notHaveCetAddress         = getAddr("000002")
	forbidAddr                = getAddr("000003")
	stock                     = "tusdt"
	money                     = "teos"
	OriginHaveCetAmount int64 = 1E13
	issueAmount         int64 = 210000000000
)

func getAddr(input string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromHex(input)
	if err != nil {
		panic(err)
	}
	return addr
}

func prepareAssetKeeper(t *testing.T, keys storeKeys, cdc *codec.Codec, ctx sdk.Context, addrForbid, tokenForbid bool) types.ExpectedAssetStatusKeeper {
	asset.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)

	//create auth, asset keeper
	ak := auth.NewAccountKeeper(
		cdc,
		keys.authCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	bk := bank.NewBaseKeeper(
		ak,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(bank.DefaultParamspace),
		sdk.CodespaceRoot,
	)

	// account permissions
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     {supply.Basic},
		authx.ModuleName:          {supply.Basic},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		asset.ModuleName:          {supply.Minter},
		ModuleName:                {supply.Basic},
	}
	sk := supply.NewKeeper(cdc, keys.keySupply, ak, bk, supply.DefaultCodespace, maccPerms)
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))
	sk.SetSupply(ctx, supply.Supply{Total: sdk.Coins{}})
	axk := authx.NewKeeper(
		cdc,
		keys.authxCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(authx.DefaultParamspace),
		sk,
		ak,
	)

	ask := asset.NewBaseTokenKeeper(
		cdc,
		keys.assetCapKey,
	)
	bkx := bankx.NewKeeper(
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(bankx.DefaultParamspace),
		axk, bk, ak, ask,
		sk,
		msgqueue.NewProducer(),
	)
	tk := asset.NewBaseKeeper(
		cdc,
		keys.assetCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(asset.DefaultParamspace),
		bkx,
		sk,
	)
	tk.SetParams(ctx, asset.DefaultParams())

	// create an account by auth keeper
	cetacc := ak.NewAccountWithAddress(ctx, haveCetAddress)
	coins := dex.NewCetCoins(OriginHaveCetAmount).
		Add(sdk.NewCoins(sdk.NewCoin(stock, sdk.NewInt(issueAmount))))
	_ = cetacc.SetCoins(coins)
	ak.SetAccount(ctx, cetacc)
	usdtacc := ak.NewAccountWithAddress(ctx, forbidAddr)
	_ = usdtacc.SetCoins(sdk.NewCoins(sdk.NewCoin(stock, sdk.NewInt(issueAmount)),
		sdk.NewCoin(dex.CET, sdk.NewInt(issueAmount))))
	ak.SetAccount(ctx, usdtacc)
	onlyIssueToken := ak.NewAccountWithAddress(ctx, notHaveCetAddress)
	_ = onlyIssueToken.SetCoins(dex.NewCetCoins(asset.IssueTokenFee))
	ak.SetAccount(ctx, onlyIssueToken)

	// issue tokens
	msgStock := asset.NewMsgIssueToken(stock, stock, issueAmount, haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", "")
	msgMoney := asset.NewMsgIssueToken(money, money, issueAmount, notHaveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", "")
	msgCet := asset.NewMsgIssueToken("cet", "cet", issueAmount, haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", "")
	handler := asset.NewHandler(tk)
	ret := handler(ctx, msgStock)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed", ret)
	ret = handler(ctx, msgMoney)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed", ret)
	ret = handler(ctx, msgCet)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed", ret)

	if tokenForbid {
		msgForbidToken := asset.MsgForbidToken{
			Symbol:       stock,
			OwnerAddress: haveCetAddress,
		}
		tk.ForbidToken(ctx, msgForbidToken.Symbol, msgForbidToken.OwnerAddress)
		msgForbidToken.Symbol = money
		tk.ForbidToken(ctx, msgForbidToken.Symbol, msgForbidToken.OwnerAddress)
	}
	if addrForbid {
		msgForbidAddr := asset.MsgForbidAddr{
			Symbol:    money,
			OwnerAddr: haveCetAddress,
			Addresses: []sdk.AccAddress{forbidAddr},
		}
		tk.ForbidAddress(ctx, msgForbidAddr.Symbol, msgForbidAddr.OwnerAddr, msgForbidAddr.Addresses)
		msgForbidAddr.Symbol = stock
		tk.ForbidAddress(ctx, msgForbidAddr.Symbol, msgForbidAddr.OwnerAddr, msgForbidAddr.Addresses)
	}

	return tk
}

func prepareBankxKeeper(keys storeKeys, cdc *codec.Codec, ctx sdk.Context) types.ExpectedBankxKeeper {
	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace)
	producer := msgqueue.NewProducer()
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)

	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     {supply.Basic},
		authx.ModuleName:          {supply.Basic},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		types.ModuleName:          {supply.Basic},
		asset.ModuleName:          {supply.Minter},
	}
	sk := supply.NewKeeper(cdc, keys.keySupply, ak, bk, supply.DefaultCodespace, maccPerms)
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))

	axk := authx.NewKeeper(cdc, keys.authxKey, paramsKeeper.Subspace(authx.DefaultParamspace), sk, ak)
	ask := asset.NewBaseTokenKeeper(cdc, keys.assetCapKey)
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, ask, sk, producer)
	bk.SetSendEnabled(ctx, true)
	bxkKeeper.SetParam(ctx, bankx.DefaultParams())

	return bxkKeeper
}

func prepareMockInput(t *testing.T, addrForbid, tokenForbid bool) testInput {
	cdc := codec.New()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	keys := storeKeys{}
	keys.assetCapKey = sdk.NewKVStoreKey(asset.StoreKey)
	keys.authCapKey = sdk.NewKVStoreKey(auth.StoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	keys.authxKey = sdk.NewKVStoreKey(authx.StoreKey)
	keys.keyStaking = sdk.NewKVStoreKey(staking.StoreKey)
	keys.tkeyStaking = sdk.NewTransientStoreKey(staking.TStoreKey)
	keys.keySupply = sdk.NewKVStoreKey(supply.StoreKey)
	keys.keyBancor = sdk.NewKVStoreKey(StoreKey)

	ms.MountStoreWithDB(keys.assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keys.authxKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyBancor, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak := prepareAssetKeeper(t, keys, cdc, ctx, addrForbid, tokenForbid)
	bk := prepareBankxKeeper(keys, cdc, ctx)

	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace)
	types.RegisterCodec(cdc)
	bik := keepers.NewBancorInfoKeeper(keys.keyBancor, cdc)
	keeper := keepers.NewKeeper(bik, bk, ak)
	akp := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)

	return testInput{ctx: ctx, bik: keeper, handler: NewHandler(keeper), akp: akp, keys: keys, cdc: cdc}
}

func Test_handleMsgBancorInit(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		msg types.MsgBancorInit
	}
	input := prepareMockInput(t, false, false)
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		{
			name: "positive",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msg: types.MsgBancorInit{
					Owner:     haveCetAddress,
					Token:     stock,
					MaxSupply: sdk.NewInt(100),
					MaxPrice:  sdk.NewDec(10),
				},
			},
			want: sdk.Result{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgBancorInit(tt.args.ctx, tt.args.k, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBancorInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgBancorTrade(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		msg types.MsgBancorTrade
	}
	input := prepareMockInput(t, false, false)
	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		{
			name: "negative token",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msg: types.MsgBancorTrade{
					Sender:     haveCetAddress,
					Token:      stock,
					Amount:     10,
					IsBuy:      true,
					MoneyLimit: 100,
				},
			},
			want: types.ErrNoBancorExists().Result(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgBancorTrade(tt.args.ctx, tt.args.k, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBancorTrade() = %v, want %v", got, tt.want)
			}
		})
	}
}
