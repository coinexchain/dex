package bancorlite_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bancorlite"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/testapp"
	dex "github.com/coinexchain/dex/types"
)

type testInput struct {
	ctx     sdk.Context
	bik     keepers.Keeper
	handler sdk.Handler
	akp     auth.AccountKeeper
	cdc     *codec.Codec // mk.cdc
}

var (
	haveCetAddress            = getAddr("000001")
	notHaveCetAddress         = getAddr("000002")
	tradeAddr                 = getAddr("000003")
	stock                     = "tusdt"
	money                     = "teos"
	OriginHaveCetAmount int64 = 1e13
	issueAmount         int64 = 210000000000
)

func getAddr(input string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromHex(input)
	if err != nil {
		panic(err)
	}
	return addr
}

func prepareApp() (testApp *testapp.TestApp, ctx sdk.Context) {
	testApp = testapp.NewTestApp()
	ctx = testApp.NewCtx()
	return
}

func prepareSupply(ctx sdk.Context, sk supply.Keeper) {
	sk.SetSupply(ctx, supply.Supply{Total: sdk.Coins{}})
}

func prepareBancor(ctx sdk.Context, k keepers.Keeper) {
	k.SetParams(ctx, bancorlite.DefaultParams())
}
func prepareAsset(t *testing.T, tk asset.Keeper, ctx sdk.Context, addrForbid, tokenForbid bool) {
	tk.SetParams(ctx, asset.DefaultParams())

	// issue tokens
	msgStock := asset.NewMsgIssueToken(stock, stock, sdk.NewInt(issueAmount), haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", asset.TestIdentityString)
	msgMoney := asset.NewMsgIssueToken(money, money, sdk.NewInt(issueAmount), notHaveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", asset.TestIdentityString)
	msgCet := asset.NewMsgIssueToken("cet", "cet", sdk.NewInt(issueAmount), haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "", asset.TestIdentityString)
	handler := asset.NewHandler(tk)
	ret := handler(ctx, msgStock)
	require.Equal(t, true, ret.IsOK(), "issue stock should succeed", ret)
	ret = handler(ctx, msgMoney)
	require.Equal(t, true, ret.IsOK(), "issue money should succeed", ret)
	ret = handler(ctx, msgCet)
	require.Equal(t, true, ret.IsOK(), "issue cet should succeed", ret)
}

func prepareAccounts(ctx sdk.Context, ak auth.AccountKeeper) {
	// create an account by auth keeper
	cetacc := ak.NewAccountWithAddress(ctx, haveCetAddress)
	coins := dex.NewCetCoins(OriginHaveCetAmount).
		Add(sdk.NewCoins(sdk.NewCoin(stock, sdk.NewInt(issueAmount))))
	_ = cetacc.SetCoins(coins)
	ak.SetAccount(ctx, cetacc)
	eosacc := ak.NewAccountWithAddress(ctx, tradeAddr)
	_ = eosacc.SetCoins(sdk.NewCoins(sdk.NewCoin(money, sdk.NewInt(issueAmount)),
		sdk.NewCoin(dex.CET, sdk.NewInt(issueAmount))))
	ak.SetAccount(ctx, eosacc)
	onlyIssueToken := ak.NewAccountWithAddress(ctx, notHaveCetAddress)
	_ = onlyIssueToken.SetCoins(dex.NewCetCoins(asset.DefaultIssue3CharTokenFee))
	ak.SetAccount(ctx, onlyIssueToken)

	//set module account
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))
}

func prepareBank(ctx sdk.Context, keeper bank.Keeper) {
	keeper.SetSendEnabled(ctx, true)
}

func prepareBankx(ctx sdk.Context, keeper bankx.Keeper) {
	keeper.SetParams(ctx, bankx.DefaultParams())
}

func prepareMarket(ctx sdk.Context, keeper market.Keeper) {
	keeper.SetParams(ctx, market.DefaultParams())
	_ = keeper.SetMarket(ctx, market.MarketInfo{Stock: stock, Money: "cet", LastExecutedPrice: sdk.NewDec(1e9)})
}

func prepareMockInput(t *testing.T, addrForbid, tokenForbid bool) testInput {
	testApp, ctx := prepareApp()

	prepareSupply(ctx, testApp.SupplyKeeper)
	prepareAccounts(ctx, testApp.AccountKeeper)
	prepareBancor(ctx, testApp.BancorKeeper)
	prepareAsset(t, testApp.AssetKeeper, ctx, addrForbid, tokenForbid)
	prepareBank(ctx, testApp.BankKeeper)
	prepareBankx(ctx, testApp.BankxKeeper)
	prepareMarket(ctx, testApp.MarketKeeper)

	return testInput{ctx: ctx, bik: testApp.BancorKeeper, handler: bancorlite.NewHandler(testApp.BancorKeeper), akp: testApp.AccountKeeper, cdc: testApp.Cdc}
}

func Test_handleMsgBancorInit(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   bancorlite.Keeper
		msg types.MsgBancorInit
	}
	input := prepareMockInput(t, false, false)
	tests := []struct {
		name string
		args args
		want string
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
					InitPrice:          "0",
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           "10",
					EarliestCancelTime: 0,
				},
			},
			want: types.ErrNonOwnerIsProhibited().Result().Log,
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
					InitPrice:          "0",
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           "10",
					MaxMoney:           sdk.NewInt(900),
					EarliestCancelTime: 0,
				},
			},
			want: sdk.Result{}.Log,
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
					InitPrice:          "0",
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           "10",
					EarliestCancelTime: 0,
				},
			},
			want: sdk.Result{}.Log,
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
					InitPrice:          "0",
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           "10",
					EarliestCancelTime: 0,
				},
			},
			want: types.ErrNoSuchToken().Result().Log,
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
					InitPrice:          "0",
					MaxSupply:          sdk.NewInt(100),
					MaxPrice:           "10",
					EarliestCancelTime: 0,
				},
			},
			want: types.ErrBancorAlreadyExists().Result().Log,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := input.handler(tt.args.ctx, tt.args.msg).Log; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBancorInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgBancorTrade(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   bancorlite.Keeper
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
			if got := input.handler(tt.args.ctx, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBancorTrade() = %v, want %v", got, tt.want)
			}
		})
	}

}
func prepareBancorInit(input testInput) bool {
	msgInit := []types.MsgBancorInit{
		{
			Owner:              haveCetAddress,
			Stock:              stock,
			Money:              money,
			InitPrice:          "0",
			MaxSupply:          sdk.NewInt(100),
			MaxPrice:           "10",
			EarliestCancelTime: 0,
		}, {
			Owner:              haveCetAddress,
			Stock:              stock,
			Money:              "cet",
			InitPrice:          "10",
			MaxSupply:          sdk.NewInt(issueAmount / 2),
			MaxPrice:           "100",
			EarliestCancelTime: 0,
		},
	}
	for _, msg := range msgInit {
		initRes := input.handler(input.ctx, msg)
		if !initRes.IsOK() {
			return false
		}
	}
	return true
}
func Test_handleMsgBancorTradeAfterInit(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		k        bancorlite.Keeper
		msgTrade types.MsgBancorTrade
	}
	input := prepareMockInput(t, false, false)
	require.True(t, prepareBancorInit(input))

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
		}, {
			name: "trade succeed",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgTrade: types.MsgBancorTrade{
					Sender:     tradeAddr,
					Stock:      stock,
					Money:      money,
					Amount:     5,
					IsBuy:      true,
					MoneyLimit: 100,
				},
			},
			want: sdk.Result{},
		}, {
			name: "stock pool out of bond",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgTrade: types.MsgBancorTrade{
					Sender:     tradeAddr,
					Stock:      stock,
					Money:      money,
					Amount:     100,
					IsBuy:      true,
					MoneyLimit: 100,
				},
			},
			want: types.ErrStockInPoolOutofBound().Result(),
		},
		{
			name: "money cross limit",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgTrade: types.MsgBancorTrade{
					Sender:     tradeAddr,
					Stock:      stock,
					Money:      money,
					Amount:     5,
					IsBuy:      false,
					MoneyLimit: 200,
				},
			},
			want: types.ErrMoneyCrossLimit("less than").Result(),
		},
		{
			name: "Insufficient coins",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgTrade: types.MsgBancorTrade{
					Sender:     tradeAddr,
					Stock:      stock,
					Money:      "cet",
					Amount:     issueAmount / 2,
					IsBuy:      true,
					MoneyLimit: math.MaxInt64,
				},
			},
			want: sdk.ErrInsufficientCoins("insufficient account funds; 204220000000cet,209999999999teos,5tusdt < 5775000000000cet").Result(),
		},
		{
			name: "trade quantity too small",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgTrade: types.MsgBancorTrade{
					Sender:     tradeAddr,
					Stock:      stock,
					Money:      "cet",
					Amount:     1,
					IsBuy:      true,
					MoneyLimit: math.MaxInt64,
				},
			},
			want: types.ErrTradeQuantityTooSmall(0).Result(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := input.handler(tt.args.ctx, tt.args.msgTrade); !reflect.DeepEqual(got, tt.want) && !got.IsOK() {
				t.Errorf("handleMsgBancorTrade() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_BancorCancel(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		k         bancorlite.Keeper
		msgCancel types.MsgBancorCancel
	}
	input := prepareMockInput(t, false, false)
	require.True(t, prepareBancorInit(input))

	tests := []struct {
		name string
		args args
		want string
	}{
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
			want: types.ErrNotBancorOwner().Result().Log,
		},
		{
			name: "cancel succeed",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgCancel: types.MsgBancorCancel{
					Owner: haveCetAddress,
					Stock: stock,
					Money: money,
				},
			},
			want: sdk.Result{}.Log,
		},
		{
			name: "bancor does not exist",
			args: args{
				ctx: input.ctx,
				k:   input.bik,
				msgCancel: types.MsgBancorCancel{
					Owner: haveCetAddress,
					Stock: money,
					Money: stock,
				},
			},
			want: types.ErrNoBancorExists().Result().Log,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := input.handler(tt.args.ctx, tt.args.msgCancel).Log; !reflect.DeepEqual(got, tt.want) {
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

	k.Save(ctx, &keepers.BancorInfo{
		Stock: "ccc",
		Money: "cet",
	})
	e = k.IsBancorExist(ctx, "ccc")
	assert.True(t, e)

	e = k.IsBancorExist(ctx, "ccb")
	assert.False(t, e)

	bi := k.Load(ctx, "ccc/abc")
	assert.Nil(t, bi)

	bi = k.Load(ctx, "ccc/cet")
	assert.Equal(t, "ccc", bi.Stock)

	k.Remove(ctx, bi)
	e = k.IsBancorExist(ctx, "ccc")
	assert.False(t, e)
}
