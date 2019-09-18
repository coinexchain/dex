package bancorlite

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
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
	OriginHaveCetAmount int64 = 1e13
	issueAmount         int64 = 210000000000
	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		authx.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		asset.ModuleName:          {supply.Minter},
		ModuleName:                nil,
	}
)

func getAddr(input string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromHex(input)
	if err != nil {
		panic(err)
	}
	return addr
}

func prepareAssetKeeper(t *testing.T, keys storeKeys, cdc *codec.Codec, ctx sdk.Context, addrForbid, tokenForbid bool) asset.Keeper {

	//create auth, asset keeper
	ak := auth.NewAccountKeeper(
		cdc,
		keys.authCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	bk := bank.NewBaseKeeper(
		ak,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(bank.DefaultParamspace),
		sdk.CodespaceRoot, map[string]bool{},
	)

	sk := supply.NewKeeper(cdc, keys.keySupply, ak, bk, maccPerms)
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))
	sk.SetSupply(ctx, supply.Supply{Total: sdk.Coins{}})
	axk := authx.NewKeeper(
		cdc,
		keys.authxCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace).Subspace(authx.DefaultParamspace),
		sk,
		ak,
		"",
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
	msgStock := asset.NewMsgIssueToken(stock, stock, sdk.NewInt(issueAmount), haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", asset.TestIdentityString)
	msgMoney := asset.NewMsgIssueToken(money, money, sdk.NewInt(issueAmount), notHaveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", asset.TestIdentityString)
	msgCet := asset.NewMsgIssueToken("cet", "cet", sdk.NewInt(issueAmount), haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", asset.TestIdentityString)
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

func prepareBankxKeeper(keys storeKeys, cdc *codec.Codec, ctx sdk.Context) bankx.Keeper {
	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace)
	producer := msgqueue.NewProducer()
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)

	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot, map[string]bool{})

	sk := supply.NewKeeper(cdc, keys.keySupply, ak, bk, maccPerms)
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))

	axk := authx.NewKeeper(cdc, keys.authxKey, paramsKeeper.Subspace(authx.DefaultParamspace), sk, ak, "")
	ask := asset.NewBaseTokenKeeper(cdc, keys.assetCapKey)
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, ask, sk, producer)
	bk.SetSendEnabled(ctx, true)
	bxkKeeper.SetParams(ctx, bankx.DefaultParams())

	return bxkKeeper
}

func prepareCodec() *codec.Codec {
	cdc := codec.New()
	types.RegisterCodec(cdc)
	asset.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	supply.RegisterCodec(cdc)

	return cdc
}
func prepareMockInput(t *testing.T, addrForbid, tokenForbid bool) testInput {

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	cdc := prepareCodec()

	keys := storeKeys{}
	keys.assetCapKey = sdk.NewKVStoreKey(asset.StoreKey)
	keys.authCapKey = sdk.NewKVStoreKey(auth.StoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	keys.authxKey = sdk.NewKVStoreKey(authx.StoreKey)
	keys.keyStaking = sdk.NewKVStoreKey(staking.StoreKey)
	keys.tkeyStaking = sdk.NewTransientStoreKey(staking.TStoreKey)
	keys.keySupply = sdk.NewKVStoreKey(supply.StoreKey)
	keys.marketKey = sdk.NewKVStoreKey(market.StoreKey)
	keys.keyBancor = sdk.NewKVStoreKey(StoreKey)

	ms.MountStoreWithDB(keys.assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keys.authxKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyBancor, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.marketKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak := prepareAssetKeeper(t, keys, cdc, ctx, addrForbid, tokenForbid)
	bk := prepareBankxKeeper(keys, cdc, ctx)
	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace)

	akp := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	mk := market.NewBaseKeeper(keys.marketKey, ak, bk, cdc,
		msgqueue.NewProducer(), paramsKeeper.Subspace(market.StoreKey), Keeper{}, akp)
	bik := keepers.NewBancorInfoKeeper(keys.keyBancor, cdc, paramsKeeper.Subspace(StoreKey))
	keeper := keepers.NewKeeper(bik, bk, ak, &mk, msgqueue.NewProducer())
	keeper.Bik.SetParams(ctx, DefaultParams())

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
			name: "not stock owner",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msg: types.MsgBancorInit{
					Owner:              notHaveCetAddress,
					Stock:              stock,
					Money:              money,
					InitPrice:          sdk.NewDec(0),
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           sdk.NewDec(10),
					EarliestCancelTime: 0,
				},
			},
			want: types.ErrNonOwnerIsProhibited().Result(),
		},
		{
			name: "positive",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msg: types.MsgBancorInit{
					Owner:              haveCetAddress,
					Stock:              stock,
					Money:              money,
					InitPrice:          sdk.NewDec(0),
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           sdk.NewDec(10),
					EarliestCancelTime: 0,
				},
			},
			want: sdk.Result{},
		},
		{
			name: "money is cet",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msg: types.MsgBancorInit{
					Owner:              haveCetAddress,
					Stock:              stock,
					Money:              dex.CET,
					InitPrice:          sdk.NewDec(0),
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           sdk.NewDec(10),
					EarliestCancelTime: 0,
				},
			},
			want: sdk.Result{},
		},
		{
			name: "stock not exist",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msg: types.MsgBancorInit{
					Owner:              haveCetAddress,
					Stock:              "abc",
					Money:              money,
					InitPrice:          sdk.NewDec(0),
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           sdk.NewDec(10),
					EarliestCancelTime: 0,
				},
			},
			want: types.ErrNoSuchToken().Result(),
		},
		{
			name: "trading pair already exist",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msg: types.MsgBancorInit{
					Owner:              haveCetAddress,
					Stock:              stock,
					Money:              money,
					InitPrice:          sdk.NewDec(0),
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           sdk.NewDec(10),
					EarliestCancelTime: 0,
				},
			},
			want: types.ErrBancorAlreadyExists().Result(),
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
					Stock:      stock,
					Money:      money,
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
func Test_handleMsgBancorTradeAfterInit(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		k        Keeper
		msgTrade types.MsgBancorTrade
	}
	input := prepareMockInput(t, false, false)

	msgInit := types.MsgBancorInit{
		Owner:              haveCetAddress,
		Stock:              stock,
		Money:              money,
		InitPrice:          sdk.NewDec(0),
		MaxSupply:          sdk.NewInt(100),
		MaxPrice:           sdk.NewDec(10),
		EarliestCancelTime: 0,
	}
	initRes := handleMsgBancorInit(input.ctx, input.bik, msgInit)
	require.True(t, initRes.IsOK())

	tests := []struct {
		name string
		args args
		want sdk.Result
	}{
		{
			name: "owner is prohibted from trading",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgTrade: types.MsgBancorTrade{
					Sender:     haveCetAddress,
					Stock:      stock,
					Money:      money,
					Amount:     10,
					IsBuy:      true,
					MoneyLimit: 100,
				},
			},
			want: types.ErrOwnerIsProhibited().Result(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgBancorTrade(tt.args.ctx, tt.args.k, tt.args.msgTrade); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBancorTrade() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_BancorCancel(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		k         Keeper
		msgCancel types.MsgBancorCancel
	}
	input := prepareMockInput(t, false, false)

	msgInit := types.MsgBancorInit{
		Owner:              haveCetAddress,
		Stock:              stock,
		Money:              money,
		InitPrice:          sdk.NewDec(0),
		MaxSupply:          sdk.NewInt(100),
		MaxPrice:           sdk.NewDec(10),
		EarliestCancelTime: 0,
	}
	initRes := handleMsgBancorInit(input.ctx, input.bik, msgInit)
	require.True(t, initRes.IsOK())

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
				msgCancel: types.MsgBancorCancel{
					Owner: haveCetAddress,
					Stock: stock,
					Money: money,
				},
			},
			want: types.ErrEarliestCancelTimeNotArrive().Result(),
		},
		{
			name: "negative token",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgCancel: types.MsgBancorCancel{
					Owner: notHaveCetAddress,
					Stock: stock,
					Money: money,
				},
			},
			want: types.ErrNotBancorOwner().Result(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMsgBancorCancel(tt.args.ctx, tt.args.k, tt.args.msgCancel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBancorTrade() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestKeeper(t *testing.T) {
	input := prepareMockInput(t, false, false)
	ctx := input.ctx
	k := input.bik
	e := k.IsBancorExist(ctx, "ccc")
	assert.False(t, e)

	k.Bik.Save(ctx, &keepers.BancorInfo{
		Stock: "ccc",
		Money: "cet",
	})
	e = k.IsBancorExist(ctx, "ccc")
	assert.True(t, e)

	e = k.IsBancorExist(ctx, "ccb")
	assert.False(t, e)

	bi := k.Bik.Load(ctx, "ccc/abc")
	assert.Nil(t, bi)

	bi = k.Bik.Load(ctx, "ccc/cet")
	assert.Equal(t, "ccc", bi.Stock)

	k.Bik.Remove(ctx, bi)
	e = k.IsBancorExist(ctx, "ccc")
	assert.False(t, e)
}
