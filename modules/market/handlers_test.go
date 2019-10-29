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

var bancorExist bool

type mockBancorKeeper struct{}

func (mbk mockBancorKeeper) IsBancorExist(ctx sdk.Context, stock string) bool {
	return bancorExist
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
	onlyIssueToken.SetCoins(dex.NewCetCoins(asset.DefaultIssueTokenFee))
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
	}
	if addrForbid {
		msgForbidAddr := asset.MsgForbidAddr{
			Symbol:    stock,
			OwnerAddr: haveCetAddress,
			Addresses: []sdk.AccAddress{forbidAddr},
		}
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
	remainCoin := dex.NewCetCoin(OriginHaveCetAmount + issueAmount - asset.DefaultIssueTokenFee*2)
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

	//// failed by price precision
	//failedPricePrecision := msgMarket
	//failedPricePrecision.Money = "cet"
	//failedPricePrecision.PricePrecision = 20
	//ret = input.handler(input.ctx, failedPricePrecision)
	//require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "create market info should failed")
	//require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

	//failedPricePrecision.PricePrecision = 19
	//ret = input.handler(input.ctx, failedPricePrecision)
	//require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "create market info should failed")
	//require.Equal(t, true, input.hasCoins(haveCetAddress, sdk.Coins{remainCoin}), "The amount is error")

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

	//failedSymbolOrder := msgOrder
	//failedSymbolOrder.TradingPair = GetSymbol(stock, "no exsit")
	//oldCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
	//ret = input.handler(input.ctx, failedSymbolOrder)
	newCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
	//require.Equal(t, types.CodeInvalidSymbol, ret.Code, "create GTE order should failed by invalid symbol")
	//require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedPricePrecisionOrder := msgOrder
	failedPricePrecisionOrder.PricePrecision = 9
	ret = input.handler(input.ctx, failedPricePrecisionOrder)
	oldCetCoin := input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "create GTE order should failed by invalid price precision")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	failedInsufficientCoinOrder := msgOrder
	failedInsufficientCoinOrder.Quantity = issueAmount * 10
	ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	require.Equal(t, types.CodeInsufficientCoin, ret.Code, "create GTE order should failed by insufficient coin")
	require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	//failedInsufficientCoinOrder = msgOrder
	//failedInsufficientCoinOrder.Quantity = 0
	//ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	//oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	//require.Equal(t, types.CodeInvalidOrderAmount, ret.Code)
	//require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

	//failedInsufficientCoinOrder = msgOrder
	//failedInsufficientCoinOrder.Quantity = 0
	//failedInsufficientCoinOrder.Side = BUY
	//ret = input.handler(input.ctx, failedInsufficientCoinOrder)
	//oldCetCoin = input.getCoinFromAddr(haveCetAddress, dex.CET)
	//require.Equal(t, types.CodeInvalidOrderAmount, ret.Code)
	//require.Equal(t, true, IsEqual(oldCetCoin, newCetCoin, zeroCet), "The amount is error")

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

	now := time.Now()
	msgCancelMarket := types.MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   GetSymbol(stock, "cet"),
		EffectiveTime: now.UnixNano() + int64(types.DefaultMarketMinExpiredTime),
	}

	header := abci.Header{Time: now, Height: 10}
	input.ctx = input.ctx.WithBlockHeader(header)
	failedTime := msgCancelMarket
	failedTime.EffectiveTime = 10
	ret := input.handler(input.ctx, failedTime)
	require.Equal(t, types.CodeInvalidCancelTime, ret.Code, "cancel order should failed by invalid cancel time")

	//failedSymbol := msgCancelMarket
	//failedSymbol.TradingPair = GetSymbol(stock, "not exist")
	//ret = input.handler(input.ctx, failedSymbol)
	//require.Equal(t, types.CodeInvalidSymbol, ret.Code, "cancel order should failed by invalid symbol")

	failedSender := msgCancelMarket
	failedSender.Sender = notHaveCetAddress
	ret = input.handler(input.ctx, failedSender)
	require.Equal(t, types.CodeNotMatchSender, ret.Code, "cancel order should failed by not match sender")

	failedByNotForbidden := msgCancelMarket
	ret = input.handler(input.ctx, failedByNotForbidden)
	require.EqualValues(t, types.CodeDelistNotAllowed, ret.Code)

}

func TestCancelMarketSuccess(t *testing.T) {
	input := prepareMockInput(t, false, true)
	createCetMarket(input, stock, 0)

	msgCancelMarket := types.MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   GetSymbol(stock, "cet"),
		EffectiveTime: int64(types.DefaultMarketMinExpiredTime + 10),
	}

	ret := input.handler(input.ctx, msgCancelMarket)
	require.Equal(t, true, ret.IsOK(), "cancel market should success")

	msgCancelMarket = types.MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   GetSymbol(stock, "cet"),
		EffectiveTime: int64(types.DefaultMarketMinExpiredTime + 10),
	}

	ret = input.handler(input.ctx, msgCancelMarket)
	require.Equal(t, false, ret.IsOK(), "repeatedly cancel market will fail")
	require.EqualValues(t, types.CodeDelistRequestExist, ret.Code)

	dlk := keepers.NewDelistKeeper(input.keys.marketKey)
	delSymbol := dlk.GetDelistSymbolsBeforeTime(input.ctx, int64(types.DefaultMarketMinExpiredTime+10+1))[0]
	require.EqualValues(t, delSymbol, GetSymbol(stock, dex.CET))
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

	//msgFailedByPricePrecision := msg
	//msgFailedByPricePrecision.PricePrecision = 19
	//ret = input.handler(input.ctx, msgFailedByPricePrecision)
	//require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "the tx should failed by dis match sender")

	//msgFailedByPricePrecision.PricePrecision = 2
	//ret = input.handler(input.ctx, msgFailedByPricePrecision)
	//require.Equal(t, types.CodeInvalidPricePrecision, ret.Code, "the tx should failed, the price precision can only be increased")

	//msgFailedByInvalidSymbol := msg
	//msgFailedByInvalidSymbol.TradingPair = GetSymbol(stock, "not find")
	//ret = input.handler(input.ctx, msgFailedByInvalidSymbol)
	//require.Equal(t, types.CodeInvalidSymbol, ret.Code, "the tx should failed by dis match sender")
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
	require.Equal(t, int64(0), fee)

	msg.ExistBlocks = 18000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(0), fee)

	msg.ExistBlocks = 20000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(1), fee)

	msg.ExistBlocks = 20001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(1), fee)

	msg.ExistBlocks = 28000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(1), fee)

	msg.ExistBlocks = 30000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(2), fee)

	msg.ExistBlocks = 30001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(2), fee)
	//
	params = types.Params{
		GTEOrderLifetime:           10000,
		GTEOrderFeatureFeeByBlocks: 10,
	}
	msg.ExistBlocks = 8000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(0), fee)

	msg.ExistBlocks = 10000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(0), fee)

	msg.ExistBlocks = 10001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(0), fee)

	msg.ExistBlocks = 18000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(8), fee)

	msg.ExistBlocks = 20000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(10), fee)

	msg.ExistBlocks = 20001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(10), fee)

	msg.ExistBlocks = 28000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(18), fee)

	msg.ExistBlocks = 30000
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(20), fee)

	msg.ExistBlocks = 30001
	fee = calFeatureFeeForExistBlocks(msg, params)
	require.Equal(t, int64(20), fee)
}

func TestCalFrozenFeeInOrder(t *testing.T) {
	input := prepareMockInput(t, false, false)
	mkInfo := MarketInfo{
		Stock:             "abc",
		Money:             "cet",
		LastExecutedPrice: sdk.NewDec(0),
		PricePrecision:    1,
		OrderPrecision:    3,
	}

	err := input.mk.SetMarket(input.ctx, mkInfo)
	require.Nil(t, err)

	checkInfo, failed := input.mk.GetMarketInfo(input.ctx, GetSymbol(mkInfo.Stock, mkInfo.Money))
	require.Nil(t, failed)
	require.EqualValues(t, mkInfo, checkInfo)

	param := types.Params{MarketFeeRate: 10, FixedTradeFee: 0}
	orderInfo := MsgCreateOrder{
		Price:       1,
		Quantity:    10000,
		TradingPair: GetSymbol(mkInfo.Stock, mkInfo.Money),
	}

	// LastExecutedPrice is zero, MarketFeeMin is zero, FixedTradeFee is zero
	cal, err := calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.Nil(t, err)
	require.EqualValues(t, 0, cal)

	// LastExecutedPrice is zero, MarketFeeMin is 100, FixedTradeFee is zero
	param.MarketFeeMin = 100
	_, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidOrderCommission, err.Code())

	// LastExecutedPrice is zero, MarketFeeMin is 100, FixedTradeFee is 10
	param.MarketFeeMin = 100
	param.FixedTradeFee = 10
	_, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidOrderCommission, err.Code())

	// LastExecutedPrice is zero, MarketFeeMin is 100, FixedTradeFee is 200
	param.MarketFeeMin = 100
	param.FixedTradeFee = 200
	cal, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.Nil(t, err)
	require.EqualValues(t, 200, cal)

	mkInfo.LastExecutedPrice = sdk.NewDec(10)
	err = input.mk.SetMarket(input.ctx, mkInfo)
	require.Nil(t, err)
	checkInfo, failed = input.mk.GetMarketInfo(input.ctx, GetSymbol(mkInfo.Stock, mkInfo.Money))
	require.Nil(t, failed)
	require.EqualValues(t, mkInfo, checkInfo)

	// LastExecutedPrice is 10, MarketFeeMin is 200; actual 10 * 10000 * (1 / 1000) = 100
	param.MarketFeeMin = 200
	_, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidOrderCommission, err.Code())

	// LastExecutedPrice is 10, MarketFeeMin is 100; actual 10 * 10000 * (1 / 1000) = 100
	param.MarketFeeMin = 100
	cal, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.Nil(t, err)
	require.EqualValues(t, 100, cal)

	// LastExecutedPrice is 10, MarketFeeMin is 100, MarketFeeRate is 10000; actual 10 * 1e18 = 1e19
	param.MarketFeeRate = 10000
	orderInfo.Quantity = types.MaxOrderAmount
	_, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidOrderAmount, err.Code())

	mkInfo.Stock = dex.CET
	mkInfo.Money = "abc"
	err = input.mk.SetMarket(input.ctx, mkInfo)
	require.Nil(t, err)
	checkInfo, failed = input.mk.GetMarketInfo(input.ctx, GetSymbol(mkInfo.Stock, mkInfo.Money))
	require.Nil(t, failed)
	require.EqualValues(t, mkInfo, checkInfo)

	// LastExecutedPrice is 10, MarketFeeMin is 100, MarketFeeRate is 10000; actual 1000
	orderInfo.Quantity = 1000
	orderInfo.TradingPair = GetSymbol(mkInfo.Stock, mkInfo.Money)
	cal, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.Nil(t, err)
	require.EqualValues(t, 1000, cal)

	// LastExecutedPrice is 10, MarketFeeMin is 10000, MarketFeeRate is 10000; actual 1000
	param.MarketFeeMin = 10000
	orderInfo.Quantity = 1000
	orderInfo.TradingPair = GetSymbol(mkInfo.Stock, mkInfo.Money)
	_, err = calFrozenFeeInOrder(input.ctx, param, input.mk, orderInfo)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidOrderCommission, err.Code())
}

func TestCheckMsgCreateOrder(t *testing.T) {
	input := prepareMockInput(t, true, true)
	require.True(t, input.mk.IsTokenForbidden(input.ctx, stock))
	require.True(t, input.mk.IsForbiddenByTokenIssuer(input.ctx, stock, forbidAddr))

	// Insufficient coin
	msg := MsgCreateOrder{
		Sender:         haveCetAddress,
		Identify:       255,
		TradingPair:    GetSymbol(stock, dex.CET),
		OrderType:      LimitOrder,
		Side:           BUY,
		Price:          10,
		PricePrecision: 8,
		Quantity:       100,
		TimeInForce:    GTE,
		ExistBlocks:    10000,
	}
	err := checkMsgCreateOrder(input.ctx, input.mk, msg, OriginHaveCetAmount+1, 1, dex.CET, 1)
	require.EqualValues(t, err.Code(), types.CodeInsufficientCoin)

	err = checkMsgCreateOrder(input.ctx, input.mk, msg, issueAmount, OriginHaveCetAmount, dex.CET, 1)
	require.EqualValues(t, err.Code(), types.CodeInsufficientCoin)

	// Invalid market
	err = checkMsgCreateOrder(input.ctx, input.mk, msg, issueAmount, issueAmount, dex.CET, math.MaxUint64)
	require.EqualValues(t, err.Code(), types.CodeInvalidMarket)

	mkInfo := MarketInfo{
		Stock:             stock,
		Money:             dex.CET,
		PricePrecision:    6,
		OrderPrecision:    1,
		LastExecutedPrice: sdk.NewDec(0),
	}
	ret := input.mk.SetMarket(input.ctx, mkInfo)
	require.Nil(t, ret)

	// Invalid price precision
	err = checkMsgCreateOrder(input.ctx, input.mk, msg, issueAmount, issueAmount, dex.CET, math.MaxUint64)
	require.EqualValues(t, err.Code(), types.CodeInvalidPricePrecision)

	// Forbidden token
	msg.PricePrecision = 6
	err = checkMsgCreateOrder(input.ctx, input.mk, msg, issueAmount, issueAmount, dex.CET, math.MaxUint64)
	require.EqualValues(t, err.Code(), types.CodeTokenForbidByIssuer)

	mkInfo.Stock = money
	mkInfo.Money = dex.CET
	ret = input.mk.SetMarket(input.ctx, mkInfo)
	require.Nil(t, ret)

	// Invalid order quantity
	msg.Sender = haveCetAddress
	msg.Quantity = 2
	msg.TradingPair = GetSymbol(money, dex.CET)
	err = checkMsgCreateOrder(input.ctx, input.mk, msg, 1, 6, dex.CET, math.MaxUint64)
	require.EqualValues(t, types.CodeInvalidOrderAmount, err.Code())

	// Pass
	msg.Sender = forbidAddr
	msg.Quantity = 10
	msg.TradingPair = GetSymbol(money, dex.CET)
	err = checkMsgCreateOrder(input.ctx, input.mk, msg, 1, 60, dex.CET, math.MaxUint64)
	require.Nil(t, err)
}

func TestCheckMsgCreateTradingPair(t *testing.T) {
	input := prepareMockInput(t, false, false)

	msg := MsgCreateTradingPair{
		Creator:        forbidAddr,
		Stock:          stock,
		Money:          dex.CET,
		PricePrecision: 8,
		OrderPrecision: 8,
	}

	// Not exist token
	msg.Money = "test"
	err := checkMsgCreateTradingPair(input.ctx, msg, input.mk)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidToken, err.Code())

	msg.Money = dex.CET
	msg.Stock = "test"
	err = checkMsgCreateTradingPair(input.ctx, msg, input.mk)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidToken, err.Code())

	// Invalid token issuer
	msg.Stock = stock
	err = checkMsgCreateTradingPair(input.ctx, msg, input.mk)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidTokenIssuer, err.Code())

	// Stock/Cet trading pair not exist
	msg.Money = money
	msg.Creator = haveCetAddress
	err = checkMsgCreateTradingPair(input.ctx, msg, input.mk)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeNotListedAgainstCet, err.Code())

	// Insufficient coin
	input.mk.SetParams(input.ctx, types.Params{
		CreateMarketFee: OriginHaveCetAmount,
	})
	msg.Creator = haveCetAddress
	msg.Money = dex.CET
	msg.Stock = stock
	err = checkMsgCreateTradingPair(input.ctx, msg, input.mk)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInsufficientCoin, err.Code())

	// Success
	input.mk.SetParams(input.ctx, types.Params{
		CreateMarketFee: 100000,
	})
	err = checkMsgCreateTradingPair(input.ctx, msg, input.mk)
	require.Nil(t, err)

	err = input.mk.SetMarket(input.ctx, MarketInfo{
		Stock:             stock,
		Money:             dex.CET,
		PricePrecision:    8,
		OrderPrecision:    0,
		LastExecutedPrice: sdk.NewDec(0),
	})

	// Invalid Repeat market
	err = checkMsgCreateTradingPair(input.ctx, msg, input.mk)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeRepeatTradingPair, err.Code())
}

func TestGetDenomAndOrderAmount(t *testing.T) {
	msg := MsgCreateOrder{
		Sender:         haveCetAddress,
		Identify:       255,
		TradingPair:    GetSymbol(stock, dex.CET),
		OrderType:      LimitOrder,
		Side:           BUY,
		Price:          11,
		PricePrecision: 8,
		Quantity:       1e8,
		TimeInForce:    GTE,
		ExistBlocks:    10000,
	}

	// 1e8 * 11 / 10^8
	denom, amount, err := getDenomAndOrderAmount(msg)
	require.Nil(t, err)
	require.EqualValues(t, dex.CET, denom)
	require.EqualValues(t, 11, amount)

	// 10 * 11 / 10^8 ≈ 10^-6
	msg.Quantity = 10
	denom, amount, err = getDenomAndOrderAmount(msg)
	require.Nil(t, err)
	require.EqualValues(t, dex.CET, denom)
	require.EqualValues(t, 1, amount)

	msg.Quantity = types.MaxOrderAmount + 1
	msg.PricePrecision = 0
	msg.Price = 1
	_, _, err = getDenomAndOrderAmount(msg)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidOrderAmount, err.Code())

	msg.Side = SELL
	msg.Quantity = 100
	denom, amount, err = getDenomAndOrderAmount(msg)
	require.Nil(t, err)
	require.EqualValues(t, stock, denom)
	require.EqualValues(t, msg.Quantity, amount)

	msg.Quantity = types.MaxOrderAmount + 1
	msg.Side = SELL
	_, _, err = getDenomAndOrderAmount(msg)
	require.NotNil(t, err)
	require.EqualValues(t, types.CodeInvalidOrderAmount, err.Code())

}

func TestCheckMsgCancelOrder(t *testing.T) {
	input := prepareMockInput(t, false, false)

	orderID, err := types.AssemblyOrderID(haveCetAddress.String(), 1, 1)
	require.Nil(t, err)

	msg := MsgCancelOrder{
		OrderID: orderID,
		Sender:  haveCetAddress,
	}
	failed := checkMsgCancelOrder(input.ctx, msg, input.mk)
	require.NotNil(t, failed)
	require.EqualValues(t, types.CodeOrderNotFound, failed.Code())

	// Create order
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

	seq, err := input.mk.QuerySeqWithAddr(input.ctx, msgGteOrder.Sender)
	require.Nil(t, err)
	ret := createCetMarket(input, stock, 10)
	require.Equal(t, true, ret.IsOK(), "create market should succeed")
	ret = input.handler(input.ctx, msgGteOrder)
	require.Equal(t, true, ret.IsOK(), "create market should succeed")

	// Invalid order sender
	orderID, err = types.AssemblyOrderID(haveCetAddress.String(), seq, msgGteOrder.Identify)
	require.Nil(t, err)
	msg.OrderID = orderID
	msg.Sender = forbidAddr
	failed = checkMsgCancelOrder(input.ctx, msg, input.mk)
	require.NotNil(t, failed)
	require.EqualValues(t, types.CodeNotMatchSender, failed.Code())

}

func TestCheckMsgCancelTradingPair(t *testing.T) {
	timeNow := time.Now()
	input := prepareMockInput(t, false, false)
	input.ctx = input.ctx.WithBlockTime(timeNow)
	param := input.mk.GetParams(input.ctx)

	msg := MsgCancelTradingPair{
		Sender:        haveCetAddress,
		TradingPair:   GetSymbol(stock, dex.CET),
		EffectiveTime: timeNow.UnixNano(),
	}

	// Invalid cancel time
	err := checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.EqualValues(t, types.CodeInvalidCancelTime, err.Code())

	msg.EffectiveTime = timeNow.UnixNano() + param.MarketMinExpiredTime - 1
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.EqualValues(t, types.CodeInvalidCancelTime, err.Code())

	// Invalid market
	msg.EffectiveTime = timeNow.UnixNano() + param.MarketMinExpiredTime
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.EqualValues(t, types.CodeInvalidMarket, err.Code())

	ret := createCetMarket(input, stock, 10)
	require.EqualValues(t, sdk.CodeOK, ret.Code)

	// Invalid sender
	msg.Sender = forbidAddr
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.EqualValues(t, types.CodeNotMatchSender, err.Code())

	// Token not forbidden when money = cet
	msg.Sender = haveCetAddress
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.EqualValues(t, types.CodeDelistNotAllowed, err.Code())

	// Token not forbidden when money != cet
	err = input.mk.SetMarket(input.ctx, MarketInfo{
		Stock: stock,
		Money: money,
	})
	require.Nil(t, err)

	msg.TradingPair = GetSymbol(stock, money)
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.Nil(t, err)

	// -----------------------

	input = prepareMockInput(t, true, true)
	input.ctx = input.ctx.WithBlockTime(timeNow)

	err = input.mk.SetMarket(input.ctx, MarketInfo{
		Stock: stock,
		Money: dex.CET,
	})
	require.Nil(t, err)

	// Token forbidden when money = cet
	msg.TradingPair = GetSymbol(stock, dex.CET)
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.Nil(t, err)

	// Bancor doesn't exist when money = cet
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.Nil(t, err)

	// Bancor exist when money = cet
	bancorExist = true
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.EqualValues(t, types.CodeDelistNotAllowed, err.Code())

	// Bancor exist when money != cet
	err = input.mk.SetMarket(input.ctx, MarketInfo{
		Stock: stock,
		Money: money,
	})
	require.Nil(t, err)

	msg.TradingPair = GetSymbol(stock, money)
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.Nil(t, err)

	// Bancor doesn't exist when money != cet
	bancorExist = false
	err = checkMsgCancelTradingPair(input.mk, msg, input.ctx)
	require.Nil(t, err)

}

func TestCheckMsgModifyPricePrecision(t *testing.T) {
	input := prepareMockInput(t, false, false)
	msg := MsgModifyPricePrecision{
		Sender:         haveCetAddress,
		TradingPair:    GetSymbol(stock, dex.CET),
		PricePrecision: 8,
	}

	// Invalid market
	err := checkMsgModifyPricePrecision(input.ctx, msg, input.mk)
	require.EqualValues(t, types.CodeInvalidMarket, err.Code())

	// Invalid price precision
	ret := createCetMarket(input, stock, 7)
	require.EqualValues(t, sdk.CodeOK, ret.Code)

	// Invalid tx sender
	msg.PricePrecision = 9
	msg.Sender = forbidAddr
	err = checkMsgModifyPricePrecision(input.ctx, msg, input.mk)
	require.EqualValues(t, types.CodeNotMatchSender, err.Code())

	msg.PricePrecision = 3
	msg.Sender = forbidAddr
	err = checkMsgModifyPricePrecision(input.ctx, msg, input.mk)
	require.EqualValues(t, types.CodeNotMatchSender, err.Code())
}
