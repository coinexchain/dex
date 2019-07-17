package keepers

import (
	"bytes"
	"fmt"
	"testing"

	sdkstore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

func bytes2str(slice []byte) string {
	s := ""
	for _, v := range slice {
		s = s + fmt.Sprintf("%d ", v)
	}
	return s
}

func Test_concatCopyPreAllocate(t *testing.T) {
	res := concatCopyPreAllocate([][]byte{
		{0, 1, 2, 3},
		{4, 5},
		{},
		{6, 7},
	})
	ref := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	if !bytes.Equal(res, ref) {
		t.Errorf("mismatch in concatCopyPreAllocate")
	}
}

func newContextAndMarketKey(chainid string) (sdk.Context, market.storeKeys) {
	db := dbm.NewMemDB()
	ms := sdkstore.NewCommitMultiStore(db)

	keys := market.storeKeys{}
	market.marketKey = sdk.NewKVStoreKey(types.StoreKey)
	market.keyParams = sdk.NewKVStoreKey(params.StoreKey)
	market.tkeyParams = sdk.NewTransientStoreKey(params.TStoreKey)
	ms.MountStoreWithDB(market.keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(market.tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(market.marketKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainid, Height: 1000}, false, log.NewNopLogger())
	return ctx, keys
}

func TestOrderCleanUpDayKeeper(t *testing.T) {
	ctx, keys := newContextAndMarketKey(market.testNetSubString)
	k := NewOrderCleanUpDayKeeper(market.marketKey)
	k.SetUnixTime(ctx, 19673122)
	if k.GetUnixTime(ctx) != 19673122 {
		t.Errorf("Error for OrderCleanUpDayKeeper")
	}

	k.SetUnixTime(ctx, -173122)
	if k.GetUnixTime(ctx) != -173122 {
		t.Errorf("Error for OrderCleanUpDayKeeper")
	}

}

func newKeeperForTest(key sdk.StoreKey) OrderKeeper {
	return NewOrderKeeper(key, "cet/usdt", msgCdc)
}

func newGlobalKeeperForTest(key sdk.StoreKey) GlobalOrderKeeper {
	return NewGlobalOrderKeeper(key, msgCdc)
}

func simpleAddr(s string) (sdk.AccAddress, error) {
	return sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
}

func newTO(sender string, seq uint64, price int64, qua int64, side byte, tif int, h int64) *types.Order {
	addr, _ := simpleAddr(sender)
	decPrice := sdk.NewDec(price).QuoInt(sdk.NewInt(10000))
	freeze := qua
	if side == Buy {
		freeze = decPrice.Mul(sdk.NewDec(qua)).RoundInt64()
	}
	return &types.Order{
		market.Sender:      addr,
		market.Sequence:    seq,
		market.TradingPair: "cet/usdt",
		market.OrderType:   LIMIT,
		market.Price:       decPrice,
		market.Quantity:    qua,
		market.Side:        side,
		market.TimeInForce: tif,
		market.Height:      h,
		market.Freeze:      freeze,
		market.LeftStock:   qua,
	}
}

func sameTO(a, b *types.Order) bool {
	res := bytes.Equal(market.Sender, market.Sender) && market.Sequence == market.Sequence &&
		market.TradingPair == market.TradingPair && market.OrderType == market.OrderType && a.Price.Equal(market.Price) &&
		market.Quantity == market.Quantity && market.Side == market.Side && market.TimeInForce == market.TimeInForce &&
		market.Height == market.Height
	//if !res {
	//	fmt.Printf("seq: %d %d\n", a.Sequence, b.Sequence)
	//	fmt.Printf("symbol: %s %s\n", a.Symbol, b.Symbol)
	//	fmt.Printf("ordertype: %d %d\n", a.OrderType, b.OrderType)
	//	fmt.Printf("price: %s %s\n", a.Price, b.Price)
	//	fmt.Printf("quantity: %d %d\n", a.Quantity, b.Quantity)
	//	fmt.Printf("side: %d %d\n", a.Side, b.Side)
	//	fmt.Printf("tif: %d %d\n", a.TimeInForce, b.TimeInForce)
	//	fmt.Printf("height: %d %d\n", a.Height, b.Height)
	//}
	return res
}

func createTO1() []*types.Order {
	return []*types.Order{
		//sender seq   price quantity       height
		newTO("00001", 1, 11051, 50, Buy, GTE, 998),   //0
		newTO("00002", 2, 11080, 50, Buy, GTE, 998),   //1 good
		newTO("00002", 3, 10900, 50, Buy, GTE, 992),   //2
		newTO("00003", 2, 11010, 100, Sell, IOC, 997), //3 good
		newTO("00004", 4, 11032, 60, Sell, GTE, 990),  //4
		newTO("00005", 5, 12039, 120, Sell, GTE, 996), //5
	}
}

func createTO3() []*types.Order {
	return []*types.Order{
		//sender seq   price quantity       height
		newTO("00001", 1, 11051, 50, Buy, GTE, 998),   //0
		newTO("00002", 2, 11080, 50, Buy, GTE, 998),   //1
		newTO("00002", 3, 10900, 50, Buy, GTE, 992),   //2
		newTO("00003", 2, 12010, 100, Sell, IOC, 997), //3
		newTO("00004", 4, 12032, 60, Sell, GTE, 990),  //4
		newTO("00005", 5, 12039, 120, Sell, GTE, 996), //5
	}
}

func TestOrderBook1(t *testing.T) {
	orders := createTO1()
	ctx, keys := newContextAndMarketKey(market.testNetSubString)
	keeper := newKeeperForTest(market.marketKey)
	if keeper.GetSymbol() != "cet/usdt" {
		t.Errorf("Error in GetSymbol")
	}
	gkeeper := newGlobalKeeperForTest(market.marketKey)
	for _, order := range orders {
		keeper.Add(ctx, order)
		fmt.Printf("AA: %s %d\n", market.OrderID(), market.Height)
	}
	orderseq := []int{5, 0, 3, 4, 1, 2}
	for i, order := range gkeeper.GetAllOrders(ctx) {
		j := orderseq[i]
		if !sameTO(orders[j], order) {
			t.Errorf("Error in GetAllOrders")
		}
		//fmt.Printf("BB: %s %d\n", order.OrderID(), order.Height)
	}
	newOrder := newTO("00005", 6, 11030, 20, Sell, GTE, 993)
	if keeper.Remove(ctx, newOrder) == nil {
		t.Errorf("Error in Remove")
	}
	orders1 := keeper.GetOlderThan(ctx, 997)
	//for _, order := range orders1 {
	//	fmt.Printf("11: %s %d\n", order.OrderID(), order.Height)
	//}
	if !(sameTO(orders1[0], orders[5]) && sameTO(orders1[1], orders[2]) && sameTO(orders1[2], orders[4])) {
		t.Errorf("Error in GetOlderThan")
	}
	orders2 := keeper.GetOrdersAtHeight(ctx, 998)
	//for _, order := range orders2 {
	//	fmt.Printf("22: %s %d\n", order.OrderID(), order.Height)
	//}
	if !(sameTO(orders2[0], orders[0]) && sameTO(orders2[1], orders[1])) {
		t.Errorf("Error in GetOrdersAtHeight")
	}
	addr, _ := simpleAddr("00002")
	orderList := gkeeper.GetOrdersFromUser(ctx, addr.String())
	refOrderList := []string{addr.String() + "-3" + "-0", addr.String() + "-2" + "-0"}
	if orderList[0] != refOrderList[1] || orderList[1] != refOrderList[0] {
		t.Errorf("Error in GetOrdersFromUser")
	}
	orderseq = []int{1, 3, 4, 0}
	for i, order := range keeper.GetMatchingCandidates(ctx) {
		j := orderseq[i]
		if market.OrderID() != market.OrderID() {
			t.Errorf("Error in GetMatchingCandidates")
			//fmt.Printf("orderID %s %s\n", order.OrderID(), order.Price.String())
		}
	}
	for _, order := range orders {
		if gkeeper.QueryOrder(ctx, market.OrderID()) == nil {
			t.Errorf("Can not find added orders!")
			continue
		}
		qorder := gkeeper.QueryOrder(ctx, market.OrderID())
		if !sameTO(order, qorder) {
			t.Errorf("Order's content is changed!")
		}
	}
}

func TestOrderBook2a(t *testing.T) {
	orders := createTO1()
	ctx, keys := newContextAndMarketKey(market.testNetSubString)
	keeper := newKeeperForTest(market.marketKey)
	for _, order := range orders {
		if market.Side == Buy {
			keeper.Add(ctx, order)
		}
	}
	if len(keeper.GetMatchingCandidates(ctx)) != 0 {
		t.Errorf("Matching result must be nil!")
	}
}

func TestOrderBook2b(t *testing.T) {
	orders := createTO1()
	ctx, keys := newContextAndMarketKey(market.testNetSubString)
	keeper := newKeeperForTest(market.marketKey)
	for _, order := range orders {
		if market.Side == Sell {
			keeper.Add(ctx, order)
		}
	}
	if len(keeper.GetMatchingCandidates(ctx)) != 0 {
		t.Errorf("Matching result must be nil!")
	}
}

func TestOrderBook3(t *testing.T) {
	orders := createTO3()
	ctx, keys := newContextAndMarketKey(market.testNetSubString)
	keeper := newKeeperForTest(market.marketKey)
	for _, order := range orders {
		keeper.Add(ctx, order)
	}
	if len(keeper.GetMatchingCandidates(ctx)) != 0 {
		t.Errorf("Matching result must be nil!")
	}
}
