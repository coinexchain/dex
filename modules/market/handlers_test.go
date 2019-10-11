package market

import (
	"bytes"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

type testInput struct {
	ctx     sdk.Context
	mk      keepers.Keeper
	handler sdk.Handler
	akp     auth.AccountKeeper
	keys    storeKeys
	cdc     *codec.Codec // mk.cdc
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
	haveCetAddress      sdk.AccAddress
	notHaveCetAddress   sdk.AccAddress
	forbidAddr          sdk.AccAddress
	stock                     = "tusdt"
	money                     = "teos"
	OriginHaveCetAmount int64 = 1e13
	issueAmount         int64 = 210000000000
	Bech32MainPrefix          = "coinex"
)

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
}

type mockBancorKeeper struct{}

func (mbk mockBancorKeeper) IsBancorExist(ctx sdk.Context, stock string) bool {
	return false
}

func initAddress() {
	haveCetAddress, _ = simpleAddr("00001")
	notHaveCetAddress, _ = simpleAddr("00002")
	forbidAddr, _ = simpleAddr("00003")
}

func prepareAssetKeeper(t *testing.T, keys storeKeys, cdc *codec.Codec, ctx sdk.Context, addrForbid, tokenForbid bool) (types.ExpectedAssetStatusKeeper, auth.AccountKeeper) {
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
		sdk.CodespaceRoot, map[string]bool{},
	)

	// account permissions
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		authx.ModuleName:          nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		types.ModuleName:          nil,
		asset.ModuleName:          {supply.Minter},
	}
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
		msgqueue.NewProducer(nil),
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
	cetacc.SetCoins(coins)
	ak.SetAccount(ctx, cetacc)
	usdtacc := ak.NewAccountWithAddress(ctx, forbidAddr)
	usdtacc.SetCoins(sdk.NewCoins(sdk.NewCoin(stock, sdk.NewInt(issueAmount)),
		sdk.NewCoin(dex.CET, sdk.NewInt(issueAmount))))
	ak.SetAccount(ctx, usdtacc)
	onlyIssueToken := ak.NewAccountWithAddress(ctx, notHaveCetAddress)
	onlyIssueToken.SetCoins(dex.NewCetCoins(asset.IssueTokenFee))
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

	return tk, ak
}

func prepareBankxKeeper(keys storeKeys, cdc *codec.Codec, ctx sdk.Context) types.ExpectedBankxKeeper {
	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace)
	producer := msgqueue.NewProducer(nil)
	ak := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)

	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot, map[string]bool{})
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		authx.ModuleName:          nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		types.ModuleName:          nil,
		asset.ModuleName:          {supply.Minter},
	}
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

func prepareMockInput(t *testing.T, addrForbid, tokenForbid bool) testInput {
	cdc := codec.New()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	initAddress()

	keys := storeKeys{}
	keys.marketKey = sdk.NewKVStoreKey(types.StoreKey)
	keys.assetCapKey = sdk.NewKVStoreKey(asset.StoreKey)
	keys.authCapKey = sdk.NewKVStoreKey(auth.StoreKey)
	keys.authxCapKey = sdk.NewKVStoreKey(authx.StoreKey)
	keys.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	keys.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	keys.authxKey = sdk.NewKVStoreKey(authx.StoreKey)
	keys.keyStaking = sdk.NewKVStoreKey(staking.StoreKey)
	keys.tkeyStaking = sdk.NewTransientStoreKey(staking.TStoreKey)
	keys.keySupply = sdk.NewKVStoreKey(supply.StoreKey)

	ms.MountStoreWithDB(keys.assetCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keys.marketKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.authxKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys.keySupply, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	ak, akp := prepareAssetKeeper(t, keys, cdc, ctx, addrForbid, tokenForbid)
	bk := prepareBankxKeeper(keys, cdc, ctx)
	paramsKeeper := params.NewKeeper(cdc, keys.keyParams, keys.tkeyParams, params.DefaultCodespace)
	// akp := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	mk := keepers.NewKeeper(keys.marketKey, ak, bk, cdc,
		msgqueue.NewProducer(nil), paramsKeeper.Subspace(types.StoreKey), mockBancorKeeper{}, akp)
	types.RegisterCodec(cdc)

	// akp := auth.NewAccountKeeper(cdc, keys.authCapKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	// subspace := paramsKeeper.Subspace(StoreKey)
	// keeper := NewKeeper(keys.marketKey, ak, bk, mockFeeKeeper{}, msgCdc, msgqueue.NewProducer(), subspace)
	parameters := types.DefaultParams()
	mk.SetParams(ctx, parameters)

	return testInput{ctx: ctx, mk: mk, handler: NewHandler(mk), akp: akp, keys: keys, cdc: cdc}
}

func TestMarketInfoSetFailed(t *testing.T) {
	input := prepareMockInput(t, false, true)
	remainCoin := dex.NewCetCoin(OriginHaveCetAmount + issueAmount - asset.IssueTokenFee*2)
	msgMarket := types.MsgCreateTradingPair{
		Stock:          stock,
		Money:          money,
		Creator:        haveCetAddress,
		PricePrecision: 8,
	}

	// failed by token not exist
	failedToken := msgMarket
	failedToken.Money = "tbtc"
	ret := input.handler(input.ctx, failedToken)
	require.Equal(t, types.CodeInvalidToken, ret.Code, "create market info should failed by token not exist")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedToken.Stock = "tiota"
	failedToken.Money = money
	ret = input.handler(input.ctx, failedToken)
	require.Equal(t, types.CodeInvalidToken, ret.Code, "create market info should failed by token not exist")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by not token issuer
	failedTokenIssuer := msgMarket
	addr, _ := simpleAddr("00008")
	failedTokenIssuer.Creator = addr
	ret = input.handler(input.ctx, failedTokenIssuer)
	require.Equal(t, types.CodeInvalidTokenIssuer, ret.Code, "create market info should failed by not token issuer")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by price precision
	failedPricePrecision := msgMarket
	failedPricePrecision.Money = "cet"
	failedPricePrecision.PricePrecision = 20
	ret = input.handler(input.ctx, failedPricePrecision)
	require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	failedPricePrecision.PricePrecision = 19
	ret = input.handler(input.ctx, failedPricePrecision)
	require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by not have sufficient cet
	failedInsufficient := msgMarket
	failedInsufficient.Creator = notHaveCetAddress
	failedInsufficient.Money = "cet"
	failedInsufficient.Stock = money
	ret = input.handler(input.ctx, failedInsufficient)
	require.Equal(t, types.CodeInsufficientCoin, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	// failed by not have cet trade
	failedNotHaveCetTrade := msgMarket
	ret = input.handler(input.ctx, failedNotHaveCetTrade)
	require.Equal(t, types.CodeNotListedAgainstCet, ret.Code, "create market info should failed")
	require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")
}

func createMarket(input testInput) sdk.Result {
	return createImpMarket(input, stock, money, 0)
}

func createImpMarket(input testInput, stock, money string, orderPrecision byte) sdk.Result {
	msgMarketInfo := types.MsgCreateTradingPair{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8, OrderPrecision: orderPrecision}
	return input.handler(input.ctx, msgMarketInfo)
}

func createCetMarket(input testInput, stock string, orderPrecision byte) sdk.Result {
	return createImpMarket(input, stock, dex.CET, orderPrecision)
}

func IsEqual(old, new sdk.Coin, diff sdk.Coin) bool {

	return old.IsEqual(new.Add(diff))
}

func TestMarketInfoSetSuccess(t *testing.T) {
	for i := 0; i <= 10; i++ {
		input := prepareMockInput(t, true, true)
		oldCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
		params := input.mk.GetParams(input.ctx)

		ret := createCetMarket(input, stock, byte(i))
		newCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
		require.Equal(t, true, ret.IsOK(), "create market info should succeed")
		require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, dex.NewCetCoin(params.CreateMarketFee)), "The amount is error")
		info, err := input.mk.GetMarketInfo(input.ctx, GetSymbol(stock, dex.CET))
		require.Nil(t, err)
		if i <= int(types.MaxOrderPrecision) {
			require.EqualValues(t, i, info.OrderPrecision)
		} else {
			require.EqualValues(t, 0, info.OrderPrecision)
		}

		for i := 0; i <= 9; i++ {
			ret = createCetMarket(input, stock, byte(i))
			require.Equal(t, types.CodeRepeatTradingPair, ret.Code)
			require.Equal(t, false, ret.IsOK(), "repeatedly creating market would fail")
		}
	}
}

func TestCreateOrderFailed(t *testing.T) {
	input := prepareMockInput(t, false, true)
	msgOrder := types.MsgCreateOrder{
		Sender:         haveCetAddress,
		TradingPair:    GetSymbol(stock, money),
		OrderType:      types.LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           types.SELL,
		TimeInForce:    types.GTE,
	}
	ret := createCetMarket(input, stock, 1)
	require.Equal(t, true, ret.IsOK(), "create market trade should success")
	ret = createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market trade should success")
	zeroCet := sdk.NewCoin("cet", sdk.NewInt(0))

	failedSymbolOrder := msgOrder
	failedSymbolOrder.TradingPair = GetSymbol(stock, "no exsit")
	oldCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
	ret = input.handler(input.ctx, failedSymbolOrder)
	newCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeInvalidSymbol, ret.Code, "create GTE order should failed by invalid symbol")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedPricePrecisionOrder := msgOrder
	failedPricePrecisionOrder.PricePrecision = 9
	ret = input.handler(input.ctx, failedPricePrecisionOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "create GTE order should failed by invalid price precision")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedInsufficientCoinOrder := msgOrder
	failedInsufficientCoinOrder.Quantity = issueAmount * 10
	ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeInsufficientCoin, ret.Code, "create GTE order should failed by insufficient coin")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedInsufficientCoinOrder = msgOrder
	failedInsufficientCoinOrder.Quantity = 0
	ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeOrderQuantityTooSmall, ret.Code, "create GTE order should failed by too small commission coin")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedInsufficientCoinOrder = msgOrder
	failedInsufficientCoinOrder.Quantity = 0
	failedInsufficientCoinOrder.Side = BUY
	ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeOrderQuantityTooSmall, ret.Code, "create GTE order should failed by too small commission coin")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedTokenForbidOrder := msgOrder
	ret = input.handler(input.ctx, failedTokenForbidOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeTokenForbidByIssuer, ret.Code, "create GTE order should failed by token forbidden by issuer")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	input = prepareMockInput(t, true, false)
	ret = createCetMarket(input, stock, 0)
	require.Equal(t, true, ret.IsOK(), "create market failed")
	ret = createMarket(input)
	require.Equal(t, true, ret.IsOK(), "create market failed")

	failedAddrForbidOrder := msgOrder
	failedAddrForbidOrder.Sender = forbidAddr
	newCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	ret = input.handler(input.ctx, failedAddrForbidOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeAddressForbidByIssuer, ret.Code, "create GTE order should failed by token forbidden by issuer")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedMaxAmount := msgOrder
	failedMaxAmount.Side = SELL
	failedMaxAmount.Quantity = 1e18 * 5
	ret = input.handler(input.ctx, failedMaxAmount)
	require.Equal(t, types.CodeInvalidOrderAmount, ret.Code, "create GTE order should failed by token forbidden by issuer")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	ret = input.handler(input.ctx, msgOrder)
	require.Equal(t, true, ret.IsOK(), "create order should succeed")

	failedOrderHaveExist := msgOrder
	newCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	ret = input.handler(input.ctx, failedOrderHaveExist)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeOrderAlreadyExist, ret.Code, "create order should failed by order exist")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")
}

func TestCalculateAmount(t *testing.T) {
	// price quantity price-precision
	items := [][]int64{{100, 10000, 2}, {300, 2000, 3}, {500, 4500, 2}}
	results := []int64{10000, 600, 22500}
	for i, item := range items {
		ret, _ := calculateAmount(item[0], item[1], byte(item[2]))
		if ret.RoundInt64() != results[i] {
			t.Errorf("amount is error, actual : %d, expect : %d", ret.RoundInt64(), results[i])
		}
	}

	for i := 2; i <= 5; i++ {
		_, err := calculateAmount(math.MaxInt64, int64(i), 0)
		require.NotNil(t, err)
	}
}

func TestCreateOrderFiledByOrderPrecision(t *testing.T) {
	for i := 1; i <= 8; i++ {
		input := prepareMockInput(t, false, false)
		msgGteOrder := types.MsgCreateOrder{
			Sender:         haveCetAddress,
			Identify:       1,
			TradingPair:    stock + types.SymbolSeparator + "cet",
			OrderType:      types.LimitOrder,
			PricePrecision: 8,
			Price:          100,
			Quantity:       10000000,
			Side:           types.SELL,
			TimeInForce:    types.GTE,
		}

		ret := createCetMarket(input, stock, byte(i))
		require.Equal(t, true, ret.IsOK(), "create market should succeed")
		failedorderPrecision := msgGteOrder
		for j := 1; j <= 8; j++ {
			failedorderPrecision.Quantity = int64(rand.Intn(int(math.Pow10(i)) - 1))
			if failedorderPrecision.Quantity == 0 {
				failedorderPrecision.Quantity = 1
			}
			failedorderPrecision.TradingPair = stock + types.SymbolSeparator + dex.CET
			ret = input.handler(input.ctx, failedorderPrecision)
			require.Equal(t, false, ret.IsOK(), "create GTE order should failed")
			require.Equal(t, types.CodeInvalidOrderAmount, ret.Code, "invalid order amount, must be a multiple of granularity ")
		}
	}

}

func TestCreateOrderSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	msgGteOrder := types.MsgCreateOrder{
		Sender:         haveCetAddress,
		Identify:       1,
		TradingPair:    GetSymbol(stock, "cet"),
		OrderType:      types.LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           types.SELL,
		TimeInForce:    types.GTE,
	}

	param := input.mk.GetParams(input.ctx)

	ret := createCetMarket(input, stock, 10)
	require.Equal(t, true, ret.IsOK(), "create market should succeed")

	seq, err := input.mk.QuerySeqWithAddr(input.ctx, msgGteOrder.Sender)
	require.Equal(t, nil, err)
	oldCoin := input.getCoinFromAddr(haveCetAddress, stock)
	ret = input.handler(input.ctx, msgGteOrder)
	newCoin := input.getCoinFromAddr(haveCetAddress, stock)
	frozenMoney := sdk.NewCoin(stock, sdk.NewInt(msgGteOrder.Quantity))
	require.Equal(t, true, ret.IsOK(), "create GTE order should succeed")
	require.Equal(t, true, IsEqual(oldCoin, newCoin, frozenMoney), "The amount is error")

	glk := keepers.NewGlobalOrderKeeper(input.keys.marketKey, input.cdc)
	orderID, err2 := types.AssemblyOrderID(msgGteOrder.Sender.String(), seq, msgGteOrder.Identify)
	require.Equal(t, nil, err2)
	order := glk.QueryOrder(input.ctx, orderID)
	require.Equal(t, true, isSameOrderAndMsg(order, msgGteOrder), "order should equal msg")

	msgIOCOrder := types.MsgCreateOrder{
		Sender:         haveCetAddress,
		Identify:       2,
		TradingPair:    GetSymbol(stock, "cet"),
		OrderType:      types.LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       68293762,
		Side:           types.BUY,
		TimeInForce:    types.IOC,
	}

	seq, err = input.mk.QuerySeqWithAddr(input.ctx, msgGteOrder.Sender)
	require.Equal(t, nil, err)
	oldCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	ret = input.handler(input.ctx, msgIOCOrder)
	newCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	frozen, _ := calculateAmount(msgIOCOrder.Price, msgIOCOrder.Quantity, msgIOCOrder.PricePrecision)
	frozenMoney = sdk.NewCoin(dex.CET, frozen.RoundInt())
	frozenFee := sdk.NewCoin(dex.CET, sdk.NewInt(param.FixedTradeFee))
	totalFrozen := frozenMoney.Add(frozenFee)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCoin, newCoin, totalFrozen), "The amount is error")

	orderID, err2 = types.AssemblyOrderID(msgIOCOrder.Sender.String(), seq, msgIOCOrder.Identify)
	require.Equal(t, nil, err2)
	order = glk.QueryOrder(input.ctx, orderID)
	require.Equal(t, true, isSameOrderAndMsg(order, msgIOCOrder), "order should equal msg")
}

func isSameOrderAndMsg(order *types.Order, msg types.MsgCreateOrder) bool {
	p := sdk.NewDec(msg.Price).Quo(sdk.NewDec(int64(math.Pow10(int(msg.PricePrecision)))))
	samePrice := order.Price.Equal(p)
	return bytes.Equal(order.Sender, msg.Sender) && order.TradingPair ==
		msg.TradingPair && order.OrderType == msg.OrderType && samePrice &&
		order.Quantity == msg.Quantity && order.Side == msg.Side &&
		order.TimeInForce == msg.TimeInForce
}

func TestCancelOrderFailed(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock, 0)
	cancelOrder := types.MsgCancelOrder{
		Sender: haveCetAddress,
	}

	failedInvalidOrderID := cancelOrder
	failedInvalidOrderID.OrderID, _ = types.AssemblyOrderID(haveCetAddress.String(), 1, 2)
	ret := input.handler(input.ctx, failedInvalidOrderID)
	require.Equal(t, types.CodeOrderNotFound, ret.Code, "cancel order should failed by not exist ")

	// create order
	msgIOCOrder := types.MsgCreateOrder{
		Sender:         haveCetAddress,
		Identify:       1,
		TradingPair:    GetSymbol(stock, "cet"),
		OrderType:      types.LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       68293762,
		Side:           types.BUY,
		TimeInForce:    types.IOC,
	}
	ret = input.handler(input.ctx, msgIOCOrder)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)

	seq, err := input.mk.QuerySeqWithAddr(input.ctx, msgIOCOrder.Sender)
	require.Equal(t, nil, err)
	failedNotOrderSender := cancelOrder
	failedNotOrderSender.OrderID, _ = types.AssemblyOrderID(msgIOCOrder.Sender.String(), seq, msgIOCOrder.Identify)
	failedNotOrderSender.Sender = notHaveCetAddress
	ret = input.handler(input.ctx, failedNotOrderSender)
	require.Equal(t, types.CodeNotMatchSender, ret.Code, "cancel order should failed by not match order sender")
}

func TestCancelOrderSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock, 0)

	// create order
	msgIOCOrder := types.MsgCreateOrder{
		Sender:         haveCetAddress,
		Identify:       2,
		TradingPair:    GetSymbol(stock, "cet"),
		OrderType:      types.LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       68293762,
		Side:           types.BUY,
		TimeInForce:    types.IOC,
	}
	seq, err := input.mk.QuerySeqWithAddr(input.ctx, msgIOCOrder.Sender)
	require.Equal(t, nil, err)
	ret := input.handler(input.ctx, msgIOCOrder)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)

	cancelOrder := types.MsgCancelOrder{
		Sender: haveCetAddress,
	}
	cancelOrder.OrderID, _ = types.AssemblyOrderID(msgIOCOrder.Sender.String(), seq, msgIOCOrder.Identify)
	ret = input.handler(input.ctx, cancelOrder)
	require.Equal(t, true, ret.IsOK(), "cancel order should succeed ; ", ret.Log)

	remainCoin := sdk.NewCoin(money, sdk.NewInt(issueAmount))
	require.Equal(t, true, input.hasCoins(notHaveCetAddress, sdk.Coins{remainCoin}), "The amount is error ")
}

func TestCancelMarketFailed(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock, 0)

	msgCancelMarket := types.MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   GetSymbol(stock, "cet"),
		EffectiveTime: time.Now().Unix() + types.DefaultMarketMinExpiredTime,
	}

	header := abci.Header{Time: time.Now(), Height: 10}
	input.ctx = input.ctx.WithBlockHeader(header)
	failedTime := msgCancelMarket
	failedTime.EffectiveTime = 10
	ret := input.handler(input.ctx, failedTime)
	require.Equal(t, types.CodeInvalidCancelTime, ret.Code, "cancel order should failed by invalid cancel time")

	failedSymbol := msgCancelMarket
	failedSymbol.TradingPair = GetSymbol(stock, "not exist")
	ret = input.handler(input.ctx, failedSymbol)
	require.Equal(t, types.CodeInvalidSymbol, ret.Code, "cancel order should failed by invalid symbol")

	failedSender := msgCancelMarket
	failedSender.Sender = notHaveCetAddress
	ret = input.handler(input.ctx, failedSender)
	require.Equal(t, types.CodeNotMatchSender, ret.Code, "cancel order should failed by not match sender")
}

func TestCancelMarketSuccess(t *testing.T) {
	input := prepareMockInput(t, false, true)
	createCetMarket(input, stock, 0)

	msgCancelMarket := types.MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   GetSymbol(stock, "cet"),
		EffectiveTime: types.DefaultMarketMinExpiredTime + 10,
	}

	ret := input.handler(input.ctx, msgCancelMarket)
	require.Equal(t, true, ret.IsOK(), "cancel market should success")

	msgCancelMarket = types.MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   GetSymbol(stock, "cet"),
		EffectiveTime: types.DefaultMarketMinExpiredTime + 10,
	}

	ret = input.handler(input.ctx, msgCancelMarket)
	require.Equal(t, false, ret.IsOK(), "repeatedly cancel market will fail")

	dlk := keepers.NewDelistKeeper(input.keys.marketKey)
	delSymbol := dlk.GetDelistSymbolsBeforeTime(input.ctx, types.DefaultMarketMinExpiredTime+10+1)[0]
	if delSymbol != GetSymbol(stock, "cet") {
		t.Error("Not find del market in store")
	}
}

func TestChargeOrderFee(t *testing.T) {
	input := prepareMockInput(t, false, false)
	ret := createCetMarket(input, stock, 0)
	require.Equal(t, true, ret.IsOK(), "create market should success")
	param := input.mk.GetParams(input.ctx)

	msgOrder := types.MsgCreateOrder{
		Sender:         haveCetAddress,
		Identify:       1,
		TradingPair:    GetSymbol(stock, dex.CET),
		OrderType:      types.LimitOrder,
		PricePrecision: 8,
		Price:          300,
		Quantity:       100000000000,
		Side:           types.BUY,
		TimeInForce:    types.IOC,
	}

	// charge fix trade fee, because the stock/cet LastExecutedPrice is zero.
	oldCetCoin := input.getCoinFromAddr(msgOrder.Sender, dex.CET)
	ret = input.handler(input.ctx, msgOrder)
	newCetCoin := input.getCoinFromAddr(msgOrder.Sender, dex.CET)
	frozen, _ := calculateAmount(msgOrder.Price, msgOrder.Quantity, msgOrder.PricePrecision)
	frozeCoin := dex.NewCetCoin(frozen.RoundInt64())
	frozeFee := dex.NewCetCoin(param.FixedTradeFee)
	totalFreeze := frozeCoin.Add(frozeFee)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, totalFreeze), "The amount is error ")

	// If stock is cet symbol, Charge a percentage of the transaction fee,
	ret = createImpMarket(input, dex.CET, stock, 0)
	require.Equal(t, true, ret.IsOK(), "create market should success")
	stockIsCetOrder := msgOrder
	stockIsCetOrder.Identify = 2
	stockIsCetOrder.TradingPair = GetSymbol(dex.CET, stock)
	oldCetCoin = input.getCoinFromAddr(msgOrder.Sender, dex.CET)
	ret = input.handler(input.ctx, stockIsCetOrder)
	newCetCoin = input.getCoinFromAddr(msgOrder.Sender, dex.CET)
	rate := sdk.NewDec(param.MarketFeeRate).Quo(sdk.NewDec(int64(math.Pow10(types.MarketFeeRatePrecision))))
	frozeFee = dex.NewCetCoin(sdk.NewDec(stockIsCetOrder.Quantity).Mul(rate).RoundInt64())
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, frozeFee), "The amount is error ")

	marketInfo, err := input.mk.GetMarketInfo(input.ctx, msgOrder.TradingPair)
	require.Equal(t, nil, err, "get %s market failed", msgOrder.TradingPair)
	marketInfo.LastExecutedPrice = sdk.NewDec(12)
	err = input.mk.SetMarket(input.ctx, marketInfo)
	require.Equal(t, nil, err, "set %s market failed", msgOrder.TradingPair)

	// Freeze fee at market execution prices
	msgOrder.Identify = 3
	oldCetCoin = input.getCoinFromAddr(msgOrder.Sender, dex.CET)
	ret = input.handler(input.ctx, msgOrder)
	newCetCoin = input.getCoinFromAddr(msgOrder.Sender, dex.CET)
	frozeFee = dex.NewCetCoin(marketInfo.LastExecutedPrice.MulInt64(msgOrder.Quantity).Mul(rate).RoundInt64())
	totalFreeze = frozeFee.Add(frozeCoin)
	require.Equal(t, true, ret.IsOK(), "create Ioc order should succeed ; ", ret.Log)
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, totalFreeze), "The amount is error ")
}

func TestModifyPricePrecisionFaild(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock, 0)

	msg := types.MsgModifyPricePrecision{
		Sender:         haveCetAddress,
		TradingPair:    GetSymbol(stock, dex.CET),
		PricePrecision: 12,
	}

	msgFailedBySender := msg
	msgFailedBySender.Sender = notHaveCetAddress
	ret := input.handler(input.ctx, msgFailedBySender)
	require.Equal(t, types.CodeNotMatchSender, ret.Code, "the tx should failed by dis match sender")

	msgFailedByPricePrecision := msg
	msgFailedByPricePrecision.PricePrecision = 19
	ret = input.handler(input.ctx, msgFailedByPricePrecision)
	require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "the tx should failed by dis match sender")

	msgFailedByPricePrecision.PricePrecision = 2
	ret = input.handler(input.ctx, msgFailedByPricePrecision)
	require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "the tx should failed, the price precision can only be increased")

	msgFailedByInvalidSymbol := msg
	msgFailedByInvalidSymbol.TradingPair = GetSymbol(stock, "not find")
	ret = input.handler(input.ctx, msgFailedByInvalidSymbol)
	require.Equal(t, types.CodeInvalidSymbol, ret.Code, "the tx should failed by dis match sender")
}

func TestModifyPricePrecisionSuccess(t *testing.T) {
	input := prepareMockInput(t, false, false)
	createCetMarket(input, stock, 0)

	msg := types.MsgModifyPricePrecision{
		Sender:         haveCetAddress,
		TradingPair:    GetSymbol(stock, dex.CET),
		PricePrecision: 12,
	}

	oldCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
	ret := input.handler(input.ctx, msg)
	newCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, true, ret.IsOK(), "the tx should success")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, sdk.NewCoin(dex.CET, sdk.NewInt(0))), "the amount is error")
}

func TestGetGranularityOfOrder(t *testing.T) {
	var expectValue = []float64{math.Pow10(0), math.Pow10(1), math.Pow10(2),
		math.Pow10(3), math.Pow10(4), math.Pow10(5), math.Pow10(6),
		math.Pow10(7), math.Pow10(8), math.Pow10(0)}
	for i := 0; i <= 9; i++ {
		ret := types.GetGranularityOfOrder(byte(i))
		require.EqualValues(t, ret, expectValue[i])
	}
}

func TestCalFeatureFeeForExistBlocks(t *testing.T) {
	msg := types.MsgCreateOrder{
		ExistBlocks: 8000,
	}
	params := types.Params{
		GTEOrderLifetime:           10000,
		GTEOrderFeatureFeeByBlocks: 1,
	}
	fee := calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(0), fee)

	msg.ExistBlocks = 10000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(0), fee)

	msg.ExistBlocks = 10001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(1), fee)

	msg.ExistBlocks = 18000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(1), fee)

	msg.ExistBlocks = 20000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(1), fee)

	msg.ExistBlocks = 20001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(2), fee)

	msg.ExistBlocks = 28000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(2), fee)

	msg.ExistBlocks = 30000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(2), fee)

	msg.ExistBlocks = 30001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(3), fee)

}
