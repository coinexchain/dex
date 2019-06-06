package market

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market/match"
	"github.com/coinexchain/dex/types"
)

type testInput struct {
	ctx     sdk.Context
	mk      Keeper
	handler sdk.Handler
	akp     auth.AccountKeeper
}

func (t testInput) hasCoins(addr sdk.AccAddress, coins sdk.Coins) bool {

	coinsStore := t.akp.GetAccount(t.ctx, addr).GetCoins()
	if len(coinsStore) < len(coins) {
		return false
	}

	for _, coin := range coins {
		find := false
		for _, coinC := range coinsStore {
			if coinC.Denom == coin.Denom {
				find = true
				if coinC.IsEqual(coin) {
					break
				} else {
					return false
				}
			}
		}
		if !find {
			return false
		}
	}

	return true
}

var (
	haveCetAddress            = getAddr("000001")
	notHaveCetAddress         = getAddr("000002")
	forbidAddr                = getAddr("000003")
	stock                     = "usdt"
	money                     = "eos"
	OriginHaveCetAmount int64 = 1E13
	issueAmount         int64 = 210000000000
)

type storeKeys struct {
	assetCapKey *sdk.KVStoreKey
	authCapKey  *sdk.KVStoreKey
	fckCapKey   *sdk.KVStoreKey
	keyParams   *sdk.KVStoreKey
	tkeyParams  *sdk.TransientStoreKey
	marketKey   *sdk.KVStoreKey
	authxKey    *sdk.KVStoreKey
}

func prepareAssetKeeper(t *testing.T, keys storeKeys, cdc *codec.Codec, ctx sdk.Context, addrForbid, tokenForbid bool) ExpectedAssertStatusKeeper {
	asset.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	//create auth, asset keeper
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, params.NewKeeper(cdc, keys.keyParams,
		keys.tkeyParams).Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)

	fck := auth.NewFeeCollectionKeeper(cdc, keys.fckCapKey)
	tk := asset.NewKeeper(cdc, keys.assetCapKey, params.NewKeeper(cdc, keys.keyParams,
		keys.tkeyParams).Subspace(asset.DefaultParamspace), ak, fck)
	tk.SetParams(ctx, asset.DefaultParams())

	// create an account by auth keeper
	cetacc := ak.NewAccountWithAddress(ctx, haveCetAddress)
	cetacc.SetCoins(types.NewCetCoins(OriginHaveCetAmount))
	ak.SetAccount(ctx, cetacc)
	usdtacc := ak.NewAccountWithAddress(ctx, forbidAddr)
	usdtacc.SetCoins(sdk.Coins{sdk.NewCoin(stock, sdk.NewInt(issueAmount))})
	ak.SetAccount(ctx, usdtacc)
	onlyIssueToken := ak.NewAccountWithAddress(ctx, notHaveCetAddress)
	onlyIssueToken.SetCoins(types.NewCetCoins(asset.IssueTokenFee))
	ak.SetAccount(ctx, onlyIssueToken)

	// issue tokens
	msgStock := asset.NewMsgIssueToken(stock, stock, issueAmount, haveCetAddress,
		false, false, addrForbid, tokenForbid)
	msgMoney := asset.NewMsgIssueToken(money, money, issueAmount, notHaveCetAddress,
		false, false, addrForbid, tokenForbid)
	handler := asset.NewHandler(tk)
	ret := handler(ctx, msgStock)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed", ret)
	ret = handler(ctx, msgMoney)
	require.Equal(t, true, ret.IsOK(), "issue token should succeed", ret)

	if tokenForbid {
		msgForbidToken := asset.MsgForbidToken{
			Symbol:       stock,
			OwnerAddress: haveCetAddress,
		}
		tk.ForbidToken(ctx, msgForbidToken)
		msgForbidToken.Symbol = money
		tk.ForbidToken(ctx, msgForbidToken)
	}
	if addrForbid {
		msgForbidAddr := asset.MsgForbidAddr{
			Symbol:     money,
			OwnerAddr:  haveCetAddress,
			ForbidAddr: []sdk.AccAddress{forbidAddr},
		}
		tk.ForbidAddress(ctx, msgForbidAddr)
		msgForbidAddr.Symbol = stock
		tk.ForbidAddress(ctx, msgForbidAddr)
	}

	return tk
}

func prepareBankxKeeper(keys storeKeys, cdc *codec.Codec, ctx sdk.Context) ExpectedBankxKeeper {

	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams)
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	fck := auth.NewFeeCollectionKeeper(cdc, keys.fckCapKey)
	axk := authx.NewKeeper(cdc, keys.authxKey, paramsKeeper.Subspace(authx.DefaultParamspace))
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, fck)
	bk.SetSendEnabled(ctx, true)
	bxkKeeper.SetParam(ctx, bankx.DefaultParams())

	return bxkKeeper
}

func prepareMockInput(t *testing.T, addrForbid, tokenForbid bool) testInput {
	cdc := codec.New()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	keys := storeKeys{}
	keys.marketKey = sdk.NewKVStoreKey(StoreKey)
	keys.assetCapKey = sdk.NewKVStoreKey(asset.StoreKey)
	keys.authCapKey = sdk.NewKVStoreKey(auth.StoreKey)
	keys.fckCapKey = sdk.NewKVStoreKey(auth.FeeStoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	keys.authxKey = sdk.NewKVStoreKey(authx.StoreKey)
	ms.MountStoreWithDB(keys.assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keys.marketKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authxKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak := prepareAssetKeeper(t, keys, cdc, ctx, addrForbid, tokenForbid)
	bk := prepareBankxKeeper(keys, cdc, ctx)

	mk := NewKeeper(keys.marketKey, ak, bk, cdc, params.NewKeeper(
		cdc, keys.keyParams, keys.tkeyParams).Subspace(StoreKey))
	RegisterCodec(mk.cdc)
	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams)
	akp := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	return testInput{ctx: ctx, mk: mk, handler: NewHandler(mk), akp: akp}
}

func TestMarketInfoSetFailed(t *testing.T) {
	input := prepareMockInput(t, false, true)
	remainCoin := types.NewCetCoin(OriginHaveCetAmount - asset.IssueTokenFee)
	msgMarket := MsgCreateMarketInfo{
		Stock:          stock,
		Money:          money,
		Creator:        haveCetAddress,
		PricePrecision: 8,
	}

	// failed by token not exist
	failedToken := msgMarket
	failedToken.Money = "btc"
	ret := input.handler(input.ctx, failedToken)
	require.Equal(t, CodeInvalidToken, ret.Code, "create market info should failed by token not exist")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedToken.Stock = "iota"
	failedToken.Money = money
	ret = input.handler(input.ctx, failedToken)
	require.Equal(t, CodeInvalidToken, ret.Code, "create market info should failed by token not exist")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by not token issuer
	failedTokenIssuer := msgMarket
	addr, _ := simpleAddr("00008")
	failedTokenIssuer.Creator = addr
	ret = input.handler(input.ctx, failedTokenIssuer)
	require.Equal(t, CodeInvalidTokenIssuer, ret.Code, "create market info should failed by not token issuer")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by price precision
	failedPricePrecision := msgMarket
	failedPricePrecision.PricePrecision = 6
	ret = input.handler(input.ctx, failedPricePrecision)
	require.Equal(t, CodeInvalidPricePrecision, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedPricePrecision.PricePrecision = 19
	ret = input.handler(input.ctx, failedPricePrecision)
	require.Equal(t, CodeInvalidPricePrecision, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by not have sufficient cet
	failedInsufficient := msgMarket
	failedInsufficient.Creator = notHaveCetAddress
	ret = input.handler(input.ctx, failedPricePrecision)
	require.Equal(t, CodeInvalidPricePrecision, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

}

func createMarket(input testInput) sdk.Result {
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8}
	return input.handler(input.ctx, msgMarketInfo)
}

func TestMarketInfoSetSuccess(t *testing.T) {
	input := prepareMockInput(t, true, true)

	//TODO. Need to determine where the deductions are incurred
	remainCoin := types.NewCetCoin(OriginHaveCetAmount - asset.IssueTokenFee)

	ret := createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market info should succeed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")
}

func TestCreateOrderFailed(t *testing.T) {
	input := prepareMockInput(t, false, true)
	msgOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		Symbol:         stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.SELL,
		TimeInForce:    GTE,
	}
	createMarket(input)

	failedSymbolOrder := msgOrder
	failedSymbolOrder.Symbol = stock + SymbolSeparator + "no exsit"
	ret := input.handler(input.ctx, failedSymbolOrder)
	remainCoin := types.NewCetCoin(OriginHaveCetAmount - asset.IssueTokenFee)
	require.Equal(t, CodeInvalidSymbol, ret.Code, "create GTE order should failed by invalid symbol")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedPricePrecisionOrder := msgOrder
	failedPricePrecisionOrder.PricePrecision = 9
	failedPricePrecisionOrder.Symbol = stock + SymbolSeparator + money
	ret = input.handler(input.ctx, failedPricePrecisionOrder)
	require.Equal(t, CodeInvalidPricePrecision, ret.Code, "create GTE order should failed by invalid price precision")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedInsufficientCoinOrder := msgOrder
	failedInsufficientCoinOrder.Quantity = issueAmount * 10
	ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	require.Equal(t, CodeInsufficientCoin, ret.Code, "create GTE order should failed by insufficient coin")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedTokenForbidOrder := msgOrder
	ret = input.handler(input.ctx, failedTokenForbidOrder)
	require.Equal(t, CodeTokenForbidByIssuer, ret.Code, "create GTE order should failed by token forbidden by issuer")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	input = prepareMockInput(t, true, false)
	createMarket(input)
	failedAddrForbidOrder := msgOrder
	failedAddrForbidOrder.Sender = forbidAddr
	ret = input.handler(input.ctx, failedAddrForbidOrder)
	//should be replace, when the forbidden addr is implement.
	//require.Equal(t, CodeTokenForbidByIssuer, ret.Code, "create GTE order should failed by token forbidden by issuer")
	//require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

}

func TestCreateOrderSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	msgGteOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		Symbol:         stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.SELL,
		TimeInForce:    GTE,
	}

	ret := createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market should succeed")

	ret = input.handler(input.ctx, msgGteOrder)
	remainCoin := sdk.NewCoin(stock, sdk.NewInt(issueAmount-10000000))
	require.Equal(t, true, ret.IsOK(), "create GTE order should succeed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	glk := NewGlobalOrderKeeper(input.mk.marketKey, input.mk.cdc)
	order := glk.QueryOrder(input.ctx, assemblyOrderID(haveCetAddress, 1))
	require.Equal(t, true, isSameOrderAndMsg(order, msgGteOrder), "order should equal msg")

	msgIOCOrder := MsgCreateOrder{
		Sender:         notHaveCetAddress,
		Sequence:       2,
		Symbol:         stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       68293762,
		Side:           Buy,
		TimeInForce:    IOC,
	}

	ret = input.handler(input.ctx, msgIOCOrder)
	remainIocCoin := sdk.NewCoin(money, sdk.NewInt(issueAmount).Sub(calculateAmount(300, 68293762, 8).RoundInt()))
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, input.hasCoins(notHaveCetAddress, sdk.Coins{remainIocCoin}), "The amount is error")

	order = glk.QueryOrder(input.ctx, assemblyOrderID(notHaveCetAddress, 2))
	require.Equal(t, true, isSameOrderAndMsg(order, msgIOCOrder), "order should equal msg")
}

func assemblyOrderID(addr sdk.AccAddress, seq uint64) string {
	return fmt.Sprintf("%s-%d", addr, seq)
}

func isSameOrderAndMsg(order *Order, msg MsgCreateOrder) bool {
	return bytes.Equal(order.Sender, msg.Sender) && order.Sequence == msg.Sequence &&
		order.Symbol == msg.Symbol && order.OrderType == msg.OrderType && order.Price.Equal(sdk.NewDec(msg.Price)) &&
		order.Quantity == msg.Quantity && order.Side == msg.Side && order.TimeInForce == msg.TimeInForce
}

func getAddr(input string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromHex(input)
	if err != nil {
		panic(err)
	}
	return addr
}

func TestCancelOrderFailed(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createMarket(input)

	cancelOrder := MsgCancelOrder{
		Sender:  haveCetAddress,
		OrderID: assemblyOrderID(haveCetAddress, 1),
	}

	failedOrderNotExist := cancelOrder
	ret := input.handler(input.ctx, failedOrderNotExist)
	require.Equal(t, CodeNotFindOrder, ret.Code, "cancel order should failed by not exist ")

	// create order
	msgIOCOrder := MsgCreateOrder{
		Sender:         notHaveCetAddress,
		Sequence:       2,
		Symbol:         stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       68293762,
		Side:           Buy,
		TimeInForce:    IOC,
	}
	ret = input.handler(input.ctx, msgIOCOrder)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)

	failedNotOrderSender := cancelOrder
	failedNotOrderSender.OrderID = assemblyOrderID(notHaveCetAddress, 2)
	ret = input.handler(input.ctx, failedNotOrderSender)
	require.Equal(t, CodeNotMatchSender, ret.Code, "cancel order should failed by not match order sender")

}

func TestCancelOrderSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createMarket(input)

	// create order
	msgIOCOrder := MsgCreateOrder{
		Sender:         notHaveCetAddress,
		Sequence:       2,
		Symbol:         stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       68293762,
		Side:           Buy,
		TimeInForce:    IOC,
	}
	ret := input.handler(input.ctx, msgIOCOrder)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)

	cancelOrder := MsgCancelOrder{
		Sender:  notHaveCetAddress,
		OrderID: assemblyOrderID(notHaveCetAddress, 2),
	}
	ret = input.handler(input.ctx, cancelOrder)
	require.Equal(t, true, ret.IsOK(), "cancel order should succeed ; ", ret.Log)

	remainCoin := sdk.NewCoin(money, sdk.NewInt(issueAmount))
	require.Equal(t, true, input.hasCoins(notHaveCetAddress, sdk.Coins{remainCoin}), "The amount is error ")
}

func TestCancelMarketFailed(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createMarket(input)

	msgCancelMarket := MsgCancelMarket{
		Sender:          haveCetAddress,
		Symbol:          stock + SymbolSeparator + money,
		EffectiveHeight: MinEffectHeight + 10,
	}

	failedHeight := msgCancelMarket
	failedHeight.EffectiveHeight = 10
	ret := input.handler(input.ctx, failedHeight)
	require.Equal(t, CodeInvalidHeight, ret.Code, "cancel order should failed by invalid cancel height")

	failedSymbol := msgCancelMarket
	failedSymbol.Symbol = stock + SymbolSeparator + "not exist"
	ret = input.handler(input.ctx, failedSymbol)
	require.Equal(t, CodeInvalidSymbol, ret.Code, "cancel order should failed by invalid symbol")

	failedSender := msgCancelMarket
	failedSender.Sender = notHaveCetAddress
	ret = input.handler(input.ctx, failedSender)
	require.Equal(t, CodeNotMatchSender, ret.Code, "cancel order should failed by not match sender")

}

func TestCancelMarketSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createMarket(input)

	msgCancelMarket := MsgCancelMarket{
		Sender:          haveCetAddress,
		Symbol:          stock + SymbolSeparator + money,
		EffectiveHeight: MinEffectHeight + 10,
	}

	ret := input.handler(input.ctx, msgCancelMarket)
	require.Equal(t, true, ret.IsOK(), "cancel market should success")

	dlk := NewDelistKeeper(input.mk.marketKey)
	delSymbol := dlk.GetDelistSymbolsAtHeight(input.ctx, MinEffectHeight+10)[0]
	if delSymbol != stock+SymbolSeparator+money {
		t.Error("Not find del market in store")
	}

}
