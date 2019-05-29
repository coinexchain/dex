package market

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/transient"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
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

func newKeeperForTest() OrderKeeper {
	return NewOrderKeeper(transient.NewStore(), "CET/USDT", msgCdc)
}

func simpleAddr(s string) (sdk.AccAddress, error) {
	return sdk.AccAddressFromHex("01234567890123456789012345678901234" + s)
}

func newTO(sender string, seq uint64, price int64, qua int64, side byte, tif int, h int64) *Order {
	addr, _ := simpleAddr(sender)
	decPrice := sdk.NewDec(price).QuoInt(sdk.NewInt(10000))
	return &Order{
		Sender:      addr,
		Sequence:    seq,
		Symbol:      "CET/USDT",
		OrderType:   LIMIT,
		Price:       decPrice,
		Quantity:    qua,
		Side:        side,
		TimeInForce: tif,
		Height:      h,
	}
}

func sameTO(a, b *Order) bool {
	return !bytes.Equal(a.Sender, b.Sender) && a.Sequence == b.Sequence &&
		a.Symbol == b.Symbol && a.OrderType == b.OrderType && a.Price.Equal(b.Price) &&
		a.Quantity == b.Quantity && a.Side == b.Side && a.TimeInForce == b.TimeInForce &&
		a.Height == b.Height
}

func createTO1() []*Order {
	return []*Order{
		//sender seq   price quantity       height
		newTO("00001", 1, 11051, 50, Buy, GTE, 998),   //0
		newTO("00002", 2, 11080, 50, Buy, GTE, 998),   //1 good
		newTO("00002", 3, 10900, 50, Buy, GTE, 992),   //2
		newTO("00003", 2, 11010, 100, Sell, IOC, 997), //3 good
		newTO("00004", 4, 11032, 60, Sell, GTE, 990),  //4
		newTO("00005", 5, 12039, 120, Sell, GTE, 996), //5
	}
}

func createTO3() []*Order {
	return []*Order{
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
	keeper := newKeeperForTest()
	for _, order := range orders {
		keeper.Add(order)
		fmt.Printf("0: %s %d\n", order.OrderID(), order.Height)
	}
	newOrder := newTO("00005", 6, 11030, 20, Sell, GTE, 993)
	if keeper.Remove(newOrder) == nil {
		t.Errorf("Error in Remove")
	}
	orders1 := keeper.GetOlderThan(997)
	if !(sameTO(orders1[0], orders[5]) && sameTO(orders1[1], orders[2]) && sameTO(orders1[2], orders[4])) {
		t.Errorf("Error in GetOlderThan")
	}
	orders2 := keeper.GetOrdersAtHeight(998)
	if !(sameTO(orders2[0], orders[0]) && sameTO(orders2[1], orders[1])) {
		t.Errorf("Error in GetOlderThan")
	}
	addr, _ := simpleAddr("00002")
	orderList := keeper.GetOrdersFromUser(addr.String())
	refOrderList := []string{addr.String() + "-3", addr.String() + "-2"}
	if orderList[0] != refOrderList[0] || orderList[1] != refOrderList[1] {
		t.Errorf("Error in GetOrdersFromUser")
	}
	for _, order := range keeper.GetMatchingCandidates() {
		fmt.Printf("orderID %s %s\n", order.OrderID(), order.Price.String())
	}
	for _, order := range orders {
		if !keeper.Exists(order.OrderID()) {
			t.Errorf("Can not find added orders!")
			continue
		}
		qorder := keeper.QueryOrder(order.OrderID())
		if !sameTO(order, qorder) {
			t.Errorf("Order's content is changed!")
		}
		keeper.Remove(order)
		if keeper.Exists(order.OrderID()) {
			t.Errorf("Can find a deleted order!")
			continue
		}
	}
}

func TestOrderBook2(t *testing.T) {
	orders := createTO1()
	keeper1 := newKeeperForTest()
	keeper2 := newKeeperForTest()
	for _, order := range orders {
		if order.Side == Buy {
			keeper1.Add(order)
		}
		if order.Side == Sell {
			keeper2.Add(order)
		}
	}
	if len(keeper1.GetMatchingCandidates()) != 0 {
		t.Errorf("Matching result must be nil!")
	}
	if len(keeper2.GetMatchingCandidates()) != 0 {
		t.Errorf("Matching result must be nil!")
	}
}

func TestOrderBook3(t *testing.T) {
	orders := createTO3()
	keeper := newKeeperForTest()
	for _, order := range orders {
		keeper.Add(order)
	}
	if len(keeper.GetMatchingCandidates()) != 0 {
		t.Errorf("Matching result must be nil!")
	}
}
