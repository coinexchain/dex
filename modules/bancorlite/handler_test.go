package bancorlite_test

import (
	"github.com/coinexchain/dex/modules/bancorlite"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/testapp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/coinexchain/dex/modules/bankx"
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
	forbidAddr                = getAddr("000003")
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

func prepareAsset(t *testing.T, app *testapp.TestApp, ctx sdk.Context, addrForbid, tokenForbid bool) {
	ak := app.AccountKeeper
	sk := app.SupplyKeeper
	tk := app.AssetKeeper
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	ak.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))
	sk.SetSupply(ctx, supply.Supply{Total: sdk.Coins{}})
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
}

func prepareMockInput(t *testing.T, addrForbid, tokenForbid bool) testInput {
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	keeper := testApp.BancorKeeper
	keeper.Bik.SetParams(ctx, bancorlite.DefaultParams())
	prepareAsset(t, testApp, ctx, addrForbid, tokenForbid)
	testApp.BankKeeper.SetSendEnabled(ctx, true)
	testApp.BankxKeeper.SetParams(ctx, bankx.DefaultParams())
	_ = testApp.MarketKeeper.SetMarket(ctx, market.MarketInfo{Stock: stock, Money: "cet"})
	return testInput{ctx: ctx, bik: keeper, handler: bancorlite.NewHandler(keeper), akp: testApp.AccountKeeper, cdc: testApp.Cdc}
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
			if got := input.handler(tt.args.ctx, tt.args.msg); !reflect.DeepEqual(got, tt.want) {
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
func Test_handleMsgBancorTradeAfterInit(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		k        bancorlite.Keeper
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
	initRes := input.handler(input.ctx, msgInit)
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
			if got := input.handler(tt.args.ctx, tt.args.msgTrade); !reflect.DeepEqual(got, tt.want) {
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

	msgInit := types.MsgBancorInit{
		Owner:              haveCetAddress,
		Stock:              stock,
		Money:              money,
		InitPrice:          sdk.NewDec(0),
		MaxSupply:          sdk.NewInt(100),
		MaxPrice:           sdk.NewDec(10),
		EarliestCancelTime: 0,
	}
	initRes := input.handler(input.ctx, msgInit)
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
					Owner: notHaveCetAddress,
					Stock: stock,
					Money: money,
				},
			},
			want: types.ErrNotBancorOwner().Result(),
		},
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
			want: sdk.Result{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := input.handler(tt.args.ctx, tt.args.msgCancel); !reflect.DeepEqual(got, tt.want) {
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
