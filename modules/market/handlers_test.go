package market

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market/match"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/types"
)

type mockFeeKeeper struct {
	storeCoins sdk.Coins
}

func (k mockFeeKeeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	k.storeCoins.Add(coins)
	fmt.Println(coins.String())

	return k.storeCoins
}

type testInput struct {
	ctx     sdk.Context
	mk      Keeper
	handler sdk.Handler
	akp     auth.AccountKeeper
}

func (t testInput) getCoinFromAddr(addr sdk.AccAddress, denom string) (cetCoin sdk.Coin) {
	coins := t.akp.GetAccount(t.ctx, addr).GetCoins()
	for _, coin := range coins {
		if coin.Denom == denom {
			cetCoin = coin
			return
		}
	}
	return
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
	stock                     = "tusdt"
	money                     = "teos"
	OriginHaveCetAmount int64 = 1E13
	issueAmount         int64 = 210000000000
)

type storeKeys struct {
	assetCapKey *sdk.KVStoreKey
	authCapKey  *sdk.KVStoreKey
	authxCapKey *sdk.KVStoreKey
	fckCapKey   *sdk.KVStoreKey
	keyParams   *sdk.KVStoreKey
	tkeyParams  *sdk.TransientStoreKey
	marketKey   *sdk.KVStoreKey
	authxKey    *sdk.KVStoreKey
	keyStaking  *sdk.KVStoreKey
	tkeyStaking *sdk.TransientStoreKey
}

func prepareAssetKeeper(t *testing.T, keys storeKeys, cdc *codec.Codec, ctx sdk.Context, addrForbid, tokenForbid bool) ExpectedAssetStatusKeeper {
	asset.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	//create auth, asset keeper
	ak := auth.NewAccountKeeper(
		cdc,
		keys.authCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams).Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	axk := authx.NewKeeper(
		cdc,
		keys.authxCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams).Subspace(authx.DefaultParamspace),
	)
	bk := bank.NewBaseKeeper(
		ak,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams).Subspace(bank.DefaultParamspace),
		sdk.CodespaceRoot,
	)
	fck := auth.NewFeeCollectionKeeper(
		cdc,
		keys.fckCapKey,
	)
	ask := asset.NewBaseTokenKeeper(
		cdc,
		keys.assetCapKey,
	)
	bkx := bankx.NewKeeper(
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams).Subspace(bankx.DefaultParamspace),
		axk, bk, ak, fck, ask,
		msgqueue.NewProducer(),
	)

	sk := staking.NewKeeper(cdc, keys.keyStaking, keys.tkeyStaking, bk,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams).Subspace(staking.DefaultParamspace),
		stakingtypes.DefaultCodespace)

	tk := asset.NewBaseKeeper(
		cdc,
		keys.assetCapKey,
		params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams).Subspace(asset.DefaultParamspace),
		bkx,
		&sk,
	)
	tk.SetParams(ctx, asset.DefaultParams())

	// create an account by auth keeper
	cetacc := ak.NewAccountWithAddress(ctx, haveCetAddress)
	coins := types.NewCetCoins(OriginHaveCetAmount).
		Add(sdk.NewCoins(sdk.NewCoin(stock, sdk.NewInt(issueAmount))))
	cetacc.SetCoins(coins)
	ak.SetAccount(ctx, cetacc)
	usdtacc := ak.NewAccountWithAddress(ctx, forbidAddr)
	usdtacc.SetCoins(sdk.NewCoins(sdk.NewCoin(stock, sdk.NewInt(issueAmount)),
		sdk.NewCoin(types.CET, sdk.NewInt(issueAmount))))
	ak.SetAccount(ctx, usdtacc)
	onlyIssueToken := ak.NewAccountWithAddress(ctx, notHaveCetAddress)
	onlyIssueToken.SetCoins(types.NewCetCoins(asset.IssueTokenFee))
	ak.SetAccount(ctx, onlyIssueToken)

	// issue tokens
	msgStock := asset.NewMsgIssueToken(stock, stock, issueAmount, haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "")
	msgMoney := asset.NewMsgIssueToken(money, money, issueAmount, notHaveCetAddress,
		false, false, addrForbid, tokenForbid, "", "")
	msgCet := asset.NewMsgIssueToken("cet", "cet", issueAmount, haveCetAddress,
		false, false, addrForbid, tokenForbid, "", "")
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

func prepareBankxKeeper(keys storeKeys, cdc *codec.Codec, ctx sdk.Context) ExpectedBankxKeeper {

	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams)
	producer := msgqueue.NewProducer()
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	fck := auth.NewFeeCollectionKeeper(cdc, keys.fckCapKey)
	axk := authx.NewKeeper(cdc, keys.authxKey, paramsKeeper.Subspace(authx.DefaultParamspace))
	ask := asset.NewBaseTokenKeeper(cdc, keys.assetCapKey)
	bxkKeeper := bankx.NewKeeper(paramsKeeper.Subspace("bankx"), axk, bk, ak, fck, ask, producer)
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
	keys.authxCapKey = sdk.NewKVStoreKey(authx.StoreKey)
	keys.fckCapKey = sdk.NewKVStoreKey(auth.FeeStoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	keys.authxKey = sdk.NewKVStoreKey(authx.StoreKey)
	keys.keyStaking = sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keys.tkeyStaking = sdk.NewTransientStoreKey(stakingtypes.TStoreKey)

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

	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams)
	mk := NewKeeper(keys.marketKey, ak, bk, mockFeeKeeper{}, cdc,
		msgqueue.NewProducer(), paramsKeeper.Subspace(StoreKey))
	RegisterCodec(mk.cdc)

	akp := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	// subspace := paramsKeeper.Subspace(StoreKey)
	// keeper := NewKeeper(keys.marketKey, ak, bk, mockFeeKeeper{}, msgCdc, msgqueue.NewProducer(), subspace)
	parameters := DefaultParams()
	mk.SetParams(ctx, parameters)

	return testInput{ctx: ctx, mk: mk, handler: NewHandler(mk), akp: akp}
}

func TestMarketInfoSetFailed(t *testing.T) {
	input := prepareMockInput(t, false, true)
	remainCoin := types.NewCetCoin(OriginHaveCetAmount + issueAmount - asset.IssueTokenFee*2)
	msgMarket := MsgCreateTradingPair{
		Stock:          stock,
		Money:          money,
		Creator:        haveCetAddress,
		PricePrecision: 8,
	}

	// failed by token not exist
	failedToken := msgMarket
	failedToken.Money = "tbtc"
	ret := input.handler(input.ctx, failedToken)
	require.Equal(t, CodeInvalidToken, ret.Code, "create market info should failed by token not exist")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedToken.Stock = "tiota"
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
	failedPricePrecision.Money = "cet"
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
	failedInsufficient.Money = "cet"
	failedInsufficient.Stock = money
	ret = input.handler(input.ctx, failedInsufficient)
	require.Equal(t, CodeInsufficientCoin, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by not have cet trade
	failedNotHaveCetTrade := msgMarket
	ret = input.handler(input.ctx, failedNotHaveCetTrade)
	require.Equal(t, CodeStockNoHaveCetTrade, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

}

func createMarket(input testInput) sdk.Result {
	return createImpMarket(input, stock, money)
}

func createImpMarket(input testInput, stock, money string) sdk.Result {
	msgMarketInfo := MsgCreateTradingPair{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8}
	return input.handler(input.ctx, msgMarketInfo)
}

func createCetMarket(input testInput, stock string) sdk.Result {
	return createImpMarket(input, stock, types.CET)
}

func IsEqual(old, new sdk.Coin, diff sdk.Coin) bool {

	return old.IsEqual(new.Add(diff))
}

func TestMarketInfoSetSuccess(t *testing.T) {
	input := prepareMockInput(t, true, true)
	oldCetCoin := input.getCoinFromAddr(haveCetAddress, types.CET)
	params := input.mk.GetParams(input.ctx)

	ret := createCetMarket(input, stock)
	newCetCoin := input.getCoinFromAddr(haveCetAddress, types.CET)
	require.Equal(t, true, ret.IsOK(), "create market info should succeed")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, types.NewCetCoin(params.CreateMarketFee)), "The amount is error")

}

func TestCreateOrderFailed(t *testing.T) {
	input := prepareMockInput(t, false, true)
	msgOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		TradingPair:    stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.SELL,
		TimeInForce:    GTE,
	}
	ret := createCetMarket(input, stock)
	require.Equal(t, true, ret.IsOK(), "create market trade should success")
	ret = createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market trade should success")
	zeroCet := sdk.NewCoin("cet", sdk.NewInt(0))

	failedSymbolOrder := msgOrder
	failedSymbolOrder.TradingPair = stock + SymbolSeparator + "no exsit"
	oldCetCoin := input.getCoinFromAddr(haveCetAddress, types.CET)
	ret = input.handler(input.ctx, failedSymbolOrder)
	newCetCoin := input.getCoinFromAddr(haveCetAddress, types.CET)
	require.Equal(t, CodeInvalidSymbol, ret.Code, "create GTE order should failed by invalid symbol")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedPricePrecisionOrder := msgOrder
	failedPricePrecisionOrder.PricePrecision = 9
	ret = input.handler(input.ctx, failedPricePrecisionOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, types.CET)
	require.Equal(t, CodeInvalidPricePrecision, ret.Code, "create GTE order should failed by invalid price precision")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedInsufficientCoinOrder := msgOrder
	failedInsufficientCoinOrder.Quantity = issueAmount * 10
	ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, types.CET)
	require.Equal(t, CodeInsufficientCoin, ret.Code, "create GTE order should failed by insufficient coin")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedTokenForbidOrder := msgOrder
	ret = input.handler(input.ctx, failedTokenForbidOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, types.CET)
	require.Equal(t, CodeTokenForbidByIssuer, ret.Code, "create GTE order should failed by token forbidden by issuer")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	input = prepareMockInput(t, true, false)
	ret = createCetMarket(input, stock)
	require.Equal(t, true, ret.IsOK(), "create market failed")
	ret = createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market failed")

	failedAddrForbidOrder := msgOrder
	failedAddrForbidOrder.Sender = forbidAddr
	newCetCoin = input.getCoinFromAddr(haveCetAddress, types.CET)
	ret = input.handler(input.ctx, failedAddrForbidOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, types.CET)
	require.Equal(t, CodeAddressForbidByIssuer, ret.Code, "create GTE order should failed by token forbidden by issuer")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

}

func TestCreateOrderSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	msgGteOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		TradingPair:    stock + SymbolSeparator + "cet",
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.SELL,
		TimeInForce:    GTE,
	}

	param := input.mk.GetParams(input.ctx)

	ret := createCetMarket(input, stock)
	require.Equal(t, true, ret.IsOK(), "create market should succeed")

	oldCoin := input.getCoinFromAddr(haveCetAddress, stock)
	ret = input.handler(input.ctx, msgGteOrder)
	newCoin := input.getCoinFromAddr(haveCetAddress, stock)
	frozenMoney := sdk.NewCoin(stock, sdk.NewInt(msgGteOrder.Quantity))
	require.Equal(t, true, ret.IsOK(), "create GTE order should succeed")
	require.Equal(t, true, IsEqual(oldCoin, newCoin, frozenMoney), "The amount is error")

	glk := NewGlobalOrderKeeper(input.mk.marketKey, input.mk.cdc)
	order := glk.QueryOrder(input.ctx, assemblyOrderID(haveCetAddress, 1))
	require.Equal(t, true, isSameOrderAndMsg(order, msgGteOrder), "order should equal msg")

	msgIOCOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       2,
		TradingPair:    stock + SymbolSeparator + "cet",
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       68293762,
		Side:           Buy,
		TimeInForce:    IOC,
	}

	oldCoin = input.getCoinFromAddr(haveCetAddress, types.CET)
	ret = input.handler(input.ctx, msgIOCOrder)
	newCoin = input.getCoinFromAddr(haveCetAddress, types.CET)
	frozenMoney = sdk.NewCoin(types.CET, calculateAmount(msgIOCOrder.Price, msgIOCOrder.Quantity, msgIOCOrder.PricePrecision).RoundInt())
	frozenFee := sdk.NewCoin(types.CET, sdk.NewInt(param.FixedTradeFee))
	totalFrozen := frozenMoney.Add(frozenFee)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCoin, newCoin, totalFrozen), "The amount is error")

	order = glk.QueryOrder(input.ctx, assemblyOrderID(haveCetAddress, 2))
	require.Equal(t, true, isSameOrderAndMsg(order, msgIOCOrder), "order should equal msg")
}

func assemblyOrderID(addr sdk.AccAddress, seq uint64) string {
	return fmt.Sprintf("%s-%d", addr, seq)
}

func isSameOrderAndMsg(order *Order, msg MsgCreateOrder) bool {
	p := sdk.NewDec(msg.Price).Quo(sdk.NewDec(int64(math.Pow10(int(msg.PricePrecision)))))
	samePrice := order.Price.Equal(p)
	return bytes.Equal(order.Sender, msg.Sender) && order.Sequence == msg.Sequence &&
		order.TradingPair == msg.TradingPair && order.OrderType == msg.OrderType && samePrice &&
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
	createCetMarket(input, stock)

	cancelOrder := MsgCancelOrder{
		Sender:  haveCetAddress,
		OrderID: assemblyOrderID(haveCetAddress, 1),
	}

	failedOrderNotExist := cancelOrder
	ret := input.handler(input.ctx, failedOrderNotExist)
	require.Equal(t, CodeNotFindOrder, ret.Code, "cancel order should failed by not exist ")

	// create order
	msgIOCOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       2,
		TradingPair:    stock + SymbolSeparator + "cet",
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
	require.Equal(t, CodeNotFindOrder, ret.Code, "cancel order should failed by not match order sender")

}

func TestCancelOrderSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock)

	// create order
	msgIOCOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       2,
		TradingPair:    stock + SymbolSeparator + "cet",
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
		Sender:  haveCetAddress,
		OrderID: assemblyOrderID(haveCetAddress, 2),
	}
	ret = input.handler(input.ctx, cancelOrder)
	require.Equal(t, true, ret.IsOK(), "cancel order should succeed ; ", ret.Log)

	remainCoin := sdk.NewCoin(money, sdk.NewInt(issueAmount))
	require.Equal(t, true, input.hasCoins(notHaveCetAddress, sdk.Coins{remainCoin}), "The amount is error ")
}

func TestCancelMarketFailed(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock)

	msgCancelMarket := MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   stock + SymbolSeparator + "cet",
		EffectiveTime: time.Now().Unix() + DefaultMarketMinExpiredTime,
	}

	header := abci.Header{Time: time.Now(), Height: 10}
	input.ctx = input.ctx.WithBlockHeader(header)
	failedTime := msgCancelMarket
	failedTime.EffectiveTime = 10
	ret := input.handler(input.ctx, failedTime)
	require.Equal(t, CodeInvalidTime, ret.Code, "cancel order should failed by invalid cancel time")

	failedSymbol := msgCancelMarket
	failedSymbol.TradingPair = stock + SymbolSeparator + "not exist"
	ret = input.handler(input.ctx, failedSymbol)
	require.Equal(t, CodeInvalidSymbol, ret.Code, "cancel order should failed by invalid symbol")

	failedSender := msgCancelMarket
	failedSender.Sender = notHaveCetAddress
	ret = input.handler(input.ctx, failedSender)
	require.Equal(t, CodeNotMatchSender, ret.Code, "cancel order should failed by not match sender")

}

func TestCancelMarketSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock)

	msgCancelMarket := MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   stock + SymbolSeparator + "cet",
		EffectiveTime: DefaultMarketMinExpiredTime + 10,
	}

	ret := input.handler(input.ctx, msgCancelMarket)
	require.Equal(t, true, ret.IsOK(), "cancel market should success")

	dlk := NewDelistKeeper(input.mk.marketKey)
	delSymbol := dlk.GetDelistSymbolsBeforeTime(input.ctx, DefaultMarketMinExpiredTime+10+1)[0]
	if delSymbol != stock+SymbolSeparator+"cet" {
		t.Error("Not find del market in store")
	}

}

func TestChargeOrderFee(t *testing.T) {
	input := prepareMockInput(t, false, false)
	ret := createCetMarket(input, stock)
	require.Equal(t, true, ret.IsOK(), "create market should success")
	param := input.mk.GetParams(input.ctx)

	msgOrder := MsgCreateOrder{
		Sender:         haveCetAddress,
		Sequence:       2,
		TradingPair:    stock + SymbolSeparator + types.CET,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       100000000000,
		Side:           Buy,
		TimeInForce:    IOC,
	}

	// charge fix trade fee, because the stock/cet LastExecutedPrice is zero.
	oldCetCoin := input.getCoinFromAddr(msgOrder.Sender, types.CET)
	ret = input.handler(input.ctx, msgOrder)
	newCetCoin := input.getCoinFromAddr(msgOrder.Sender, types.CET)
	frozeCoin := types.NewCetCoin(calculateAmount(msgOrder.Price, msgOrder.Quantity, msgOrder.PricePrecision).RoundInt64())
	frozeFee := types.NewCetCoin(param.FixedTradeFee)
	totalFreeze := frozeCoin.Add(frozeFee)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, totalFreeze), "The amount is error ")

	// If stock is cet symbol, Charge a percentage of the transaction fee,
	ret = createImpMarket(input, types.CET, stock)
	require.Equal(t, true, ret.IsOK(), "create market should success")
	stockIsCetOrder := msgOrder
	stockIsCetOrder.TradingPair = types.CET + SymbolSeparator + stock
	oldCetCoin = input.getCoinFromAddr(msgOrder.Sender, types.CET)
	ret = input.handler(input.ctx, stockIsCetOrder)
	newCetCoin = input.getCoinFromAddr(msgOrder.Sender, types.CET)
	rate := sdk.NewDec(param.MarketFeeRate).Quo(sdk.NewDec(int64(math.Pow10(MarketFeeRatePrecision))))
	frozeFee = types.NewCetCoin(sdk.NewDec(stockIsCetOrder.Quantity).Mul(rate).RoundInt64())
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, frozeFee), "The amount is error ")

	marketInfo, err := input.mk.GetMarketInfo(input.ctx, msgOrder.TradingPair)
	require.Equal(t, nil, err, "get %s market failed", msgOrder.TradingPair)
	marketInfo.LastExecutedPrice = sdk.NewDec(12)
	err = input.mk.SetMarket(input.ctx, marketInfo)
	require.Equal(t, nil, err, "set %s market failed", msgOrder.TradingPair)

	// Freeze fee at market execution prices
	oldCetCoin = input.getCoinFromAddr(msgOrder.Sender, types.CET)
	ret = input.handler(input.ctx, msgOrder)
	newCetCoin = input.getCoinFromAddr(msgOrder.Sender, types.CET)
	frozeFee = types.NewCetCoin(marketInfo.LastExecutedPrice.MulInt64(msgOrder.Quantity).Mul(rate).RoundInt64())
	totalFreeze = frozeFee.Add(frozeCoin)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, totalFreeze), "The amount is error ")

}
