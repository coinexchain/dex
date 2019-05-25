package match

import (
	"crypto/sha256"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

type dealRecord struct {
	self   string
	other  string
	amount int64
	price  sdk.Dec
}

func newDR(self string, other string, amount int64, price int64) dealRecord {
	return dealRecord{
		self:   self,
		other:  other,
		amount: amount,
		price:  sdk.NewDec(price),
	}
}

func newDRDec(self string, other string, amount int64, price sdk.Dec) dealRecord {
	return dealRecord{
		self:   self,
		other:  other,
		amount: amount,
		price:  price,
	}
}

type mocAccount struct {
	name string
}

func (ma mocAccount) String() string {
	return ma.name
}

var currDealRecordList []dealRecord
var currDealRecordIndex int
var testHandler *testing.T

type mocOrder struct {
	price        sdk.Dec
	height       int64
	totalAmount  int64
	remainAmount int64
	orderType    int
	owner        mocAccount
}

var _ OrderForTrade = &mocOrder{}
var _ Account = &mocAccount{}

func (order *mocOrder) GetPrice() sdk.Dec {
	return order.price
}

func (order *mocOrder) GetAmount() int64 {
	return order.remainAmount
}

func (order *mocOrder) GetHeight() int64 {
	return order.height
}

func (order *mocOrder) GetHash() []byte {
	res := sha256.Sum256([]byte(order.owner.name))
	return res[:]
}

func (order *mocOrder) GetType() int {
	return order.orderType
}

func (order *mocOrder) GetOwner() Account {
	return &order.owner
}

func (order *mocOrder) Deal(otherSide OrderForTrade, amount int64, price sdk.Dec) {
	other := otherSide.(*mocOrder)
	fmt.Printf("Deal: %s|%d-%s|%d %d price:%s\n", order.GetOwner(), order.GetAmount(), other.GetOwner(), other.GetAmount(), amount, price.String())
	order.remainAmount -= amount
	other.remainAmount -= amount
	if len(currDealRecordList) == 0 {
		return
	}
	if currDealRecordIndex >= len(currDealRecordList) {
		fmt.Printf("More deals than expected! %d>=%d", currDealRecordIndex, len(currDealRecordList))
		testHandler.Errorf("More deals than expected!")
	} else {
		dr := currDealRecordList[currDealRecordIndex]
		currDealRecordIndex++
		pass := dr.self == order.GetOwner().String() && dr.other == other.GetOwner().String() &&
			dr.amount == amount && dr.price.Equal(price)
		if !pass {
			testHandler.Errorf("incorrect deal! i:%d this:%s that:%s a:%d p:%s\n", currDealRecordIndex, dr.self,
				dr.other, dr.amount, dr.price.String())
			fmt.Printf("incorrect deal! i:%d this:%s that:%s a:%d p:%s\n", currDealRecordIndex, dr.self,
				dr.other, dr.amount, dr.price.String())
		}
	}
}

func (order *mocOrder) String() string {
	s := "sell"
	if order.GetType() == BUY {
		s = "buy"
	}
	return fmt.Sprintf("%s %s %d at %s (%d)", order.GetOwner().String(), s, order.GetAmount(),
		order.GetPrice().String(), order.GetHeight())
}

func newMocOrder(price int64, height int64, totalAmount int64, orderType int, owner string) OrderForTrade {
	return &mocOrder{
		price:        sdk.NewDec(price),
		height:       height,
		totalAmount:  totalAmount,
		remainAmount: totalAmount,
		orderType:    orderType,
		owner:        mocAccount{name: owner},
	}
}

func createOrders1() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(100, 1, 150, BUY, "buyer1"),
		newMocOrder(98, 1, 150, BUY, "buyer2"),
		newMocOrder(98, 1, 250, SELL, "seller1"),
		newMocOrder(97, 1, 50, SELL, "seller2"),
	}
}

func createOrders2() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(100, 1, 150, BUY, "buyer1"),
		newMocOrder(99, 1, 50, BUY, "buyer2"),
		newMocOrder(97, 1, 300, BUY, "buyer3"),
		newMocOrder(97, 1, 200, SELL, "seller1"),
		newMocOrder(96, 1, 100, SELL, "seller2"),
	}
}

func createOrders2_1() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(100, 3, 50, BUY, "buyer1a"),
		newMocOrder(100, 2, 50, BUY, "buyer1b"),
		newMocOrder(100, 1, 50, BUY, "buyer1c"),
		newMocOrder(99, 1, 50, BUY, "buyer2"),
		newMocOrder(97, 2, 50, BUY, "buyer3a"),
		newMocOrder(97, 2, 50, BUY, "buyer3b"),
		newMocOrder(97, 1, 100, BUY, "buyer3c"),
		newMocOrder(97, 3, 100, BUY, "buyer3d"),
		newMocOrder(97, 1, 200, SELL, "seller1"),
		newMocOrder(96, 1, 100, SELL, "seller2"),
	}
}

func createOrders3() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(102, 1, 300, BUY, "buyer1"),
		newMocOrder(100, 1, 100, BUY, "buyer2"),
		newMocOrder(99, 1, 200, BUY, "buyer3"),
		newMocOrder(98, 1, 300, BUY, "buyer4"),
		newMocOrder(98, 1, 250, SELL, "seller1"),
		newMocOrder(97, 1, 250, SELL, "seller2"),
		newMocOrder(96, 1, 1000, SELL, "seller3"),
	}
}

func createOrders4() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(102, 1, 30, BUY, "buyer1"),
		newMocOrder(101, 1, 10, BUY, "buyer2"),
		newMocOrder(99, 1, 50, BUY, "buyer3"),
		newMocOrder(96, 1, 15, BUY, "buyer4"),
		newMocOrder(98, 1, 10, SELL, "seller1"),
		newMocOrder(97, 1, 50, SELL, "seller2"),
		newMocOrder(95, 1, 50, SELL, "seller3"),
	}
}

func createOrders4_1() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(102, 1, 30, BUY, "buyer1"),
		newMocOrder(101, 1, 10, BUY, "buyer2"),
		newMocOrder(99, 1, 10, BUY, "buyer3"),
		newMocOrder(96, 1, 15, BUY, "buyer4"),
		newMocOrder(98, 1, 10, BUY, "buyer5"),
		newMocOrder(97, 1, 50, SELL, "seller2"),
		newMocOrder(95, 1, 50, SELL, "seller3"),
	}
}

func createOrders5_1() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(102, 1, 10, BUY, "buyer1"),
		newMocOrder(97, 1, 10, BUY, "buyer2"),
		newMocOrder(95, 1, 50, SELL, "seller1"),
	}
}

func createOrders5_2() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(99, 1, 10, BUY, "buyer1"),
		newMocOrder(94, 1, 10, BUY, "buyer2"),
		newMocOrder(92, 1, 50, SELL, "seller1"),
	}
}

func createOrders5_3() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(99, 1, 100, BUY, "buyer1"),
		newMocOrder(92, 1, 50, SELL, "seller1"),
	}
}

func createOrders5_4() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(101, 1, 10, BUY, "buyer1"),
		newMocOrder(96, 1, 10, BUY, "buyer2"),
		newMocOrder(94, 1, 50, SELL, "seller1"),
	}
}

func createOrders6() []OrderForTrade {
	//             price height totalAmount orderType owner
	return []OrderForTrade{
		newMocOrder(100, 1, 25, BUY, "buyer1"),
		newMocOrder(97, 1, 25, BUY, "buyer2"),
		newMocOrder(98, 1, 25, SELL, "seller1"),
		newMocOrder(95, 1, 25, SELL, "seller2"),
	}
}

func createDealRecord1() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller2", 50, 98),
		newDR("buyer1", "seller1", 100, 98),
		newDR("seller1", "buyer2", 150, 98),
	}
}
func createDealRecord2() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller2", 100, 97),
		newDR("buyer1", "seller1", 50, 97),
		newDR("seller1", "buyer2", 50, 97),
		newDR("seller1", "buyer3", 100, 97),
	}
}
func createDealRecord2_1() []dealRecord {
	return []dealRecord{
		newDR("buyer1c", "seller2", 50, 97),
		newDR("seller2", "buyer1b", 50, 97),
		newDR("buyer1a", "seller1", 50, 97),
		newDR("seller1", "buyer2", 50, 97),
		newDR("seller1", "buyer3c", 100, 97),
	}
}
func createDealRecord3() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller3", 300, 96),
		newDR("seller3", "buyer2", 100, 96),
		newDR("seller3", "buyer3", 200, 96),
		newDR("seller3", "buyer4", 300, 96),
	}
}
func createDealRecord4() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller3", 30, 97),
		newDR("seller3", "buyer2", 10, 97),
		newDR("seller3", "buyer3", 10, 97),
		newDR("buyer3", "seller2", 40, 97),
	}
}
func createDealRecord4_1() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller3", 30, 97),
		newDR("seller3", "buyer2", 10, 97),
		newDR("seller3", "buyer3", 10, 97),
		newDR("buyer5", "seller2", 10, 97),
	}
}
func createDealRecord5_1() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller1", 10, 95),
		newDR("seller1", "buyer2", 10, 95),
	}
}
func createDealRecord5_2() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller1", 10, 94),
		newDR("seller1", "buyer2", 10, 94),
	}
}
func createDealRecord5_3a() []dealRecord {
	return []dealRecord{
		newDRDec("buyer1", "seller1", 50, sdk.NewDec(945).QuoInt(sdk.NewInt(10))),
	}
}
func createDealRecord5_3b() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller1", 50, 99),
	}
}
func createDealRecord5_3c() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller1", 50, 92),
	}
}
func createDealRecord5_4() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller1", 10, 95),
		newDR("seller1", "buyer2", 10, 95),
	}
}
func createDealRecord6_1() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller2", 25, 99),
	}
}
func createDealRecord6_2() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller2", 25, 97),
	}
}
func createDealRecord6_3() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller2", 25, 95),
	}
}
func createDealRecord6_4() []dealRecord {
	return []dealRecord{
		newDR("buyer1", "seller2", 25, 100),
	}
}

func testGetExecutionPrice(mid int64, orders []OrderForTrade) sdk.Dec {
	midPrice := sdk.NewDec(mid)
	highPrice := midPrice.MulInt(sdk.NewInt(105)).QuoInt(sdk.NewInt(100))
	lowPrice := midPrice.MulInt(sdk.NewInt(95)).QuoInt(sdk.NewInt(100))
	return GetExecutionPrice(highPrice, midPrice, lowPrice, orders)
}

func testMatch(tag string, mid int64, orders []OrderForTrade, dealRecordList []dealRecord) {
	currDealRecordList = dealRecordList
	currDealRecordIndex = 0
	fmt.Printf("=======================%s===============================\n", tag)
	bidList := make([]OrderForTrade, 0, 10)
	askList := make([]OrderForTrade, 0, 10)
	for _, order := range orders {
		if order.GetType() == BID {
			bidList = append(bidList, order)
		}
		if order.GetType() == ASK {
			askList = append(askList, order)
		}
	}
	midPrice := sdk.NewDec(mid)
	highPrice := midPrice.MulInt(sdk.NewInt(105)).QuoInt(sdk.NewInt(100))
	lowPrice := midPrice.MulInt(sdk.NewInt(95)).QuoInt(sdk.NewInt(100))
	Match(highPrice, midPrice, lowPrice, bidList, askList)
	if currDealRecordIndex != len(currDealRecordList) {
		testHandler.Errorf("Missmatch in the count of deals")
	}
}

func TestPrice_1(t *testing.T) {
	p := testGetExecutionPrice(100, createOrders1())
	fmt.Printf("==================1: %s\n", p)
	if !p.Equal(sdk.NewDec(98)) {
		t.Errorf("::::1:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(100, createOrders2())
	fmt.Printf("==================2: %s\n", p)
	if !p.Equal(sdk.NewDec(97)) {
		t.Errorf("::::2:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(100, createOrders3())
	fmt.Printf("==================3: %s\n", p)
	if !p.Equal(sdk.NewDec(96)) {
		t.Errorf("::::3:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(100, createOrders4())
	fmt.Printf("==================4: %s\n", p)
	if !p.Equal(sdk.NewDec(97)) {
		t.Errorf("::::4:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(80, createOrders4_1())
	fmt.Printf("==================4_1: %s\n", p)
	if !p.Equal(sdk.NewDec(97)) {
		t.Errorf("::::4_1:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(80, createOrders5_1())
	fmt.Printf("==================5_1: %s\n", p)
	if !p.Equal(sdk.NewDec(95)) {
		t.Errorf("::::5_1:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(100, createOrders5_2())
	fmt.Printf("==================5_2: %s\n", p)
	if !p.Equal(sdk.NewDec(94)) {
		t.Errorf("::::5_2:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(90, createOrders5_3())
	fmt.Printf("==================5_3a: %s\n", p)
	if !p.Equal(sdk.NewDec(945).Quo(sdk.NewDec(10))) {
		t.Errorf("::::5_3a:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(97, createOrders5_3())
	fmt.Printf("==================5_3b: %s\n", p)
	if !p.Equal(sdk.NewDec(99)) {
		t.Errorf("::::5_3b:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(80, createOrders5_3())
	fmt.Printf("==================5_3c: %s\n", p)
	if !p.Equal(sdk.NewDec(92)) {
		t.Errorf("::::5_3c:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(100, createOrders5_4())
	fmt.Printf("==================5_4: %s\n", p)
	if !p.Equal(sdk.NewDec(95)) {
		t.Errorf("::::5_4:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(99, createOrders6())
	fmt.Printf("==================6_1: %s\n", p)
	if !p.Equal(sdk.NewDec(99)) {
		t.Errorf("::::6_1:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(97, createOrders6())
	fmt.Printf("==================6_2: %s\n", p)
	if !p.Equal(sdk.NewDec(97)) {
		t.Errorf("::::6_2:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(90, createOrders6())
	fmt.Printf("==================6_3: %s\n", p)
	if !p.Equal(sdk.NewDec(95)) {
		t.Errorf("::::6_3:::: Wrong:%s\n", p)
	}
	p = testGetExecutionPrice(110, createOrders6())
	fmt.Printf("==================6_4: %s\n", p)
	if !p.Equal(sdk.NewDec(100)) {
		t.Errorf("::::6_4:::: Wrong:%s\n", p)
	}
}

func TestMatch_1(t *testing.T) {
	testHandler = t
	testMatch("1", 100, createOrders1(), createDealRecord1())
	testMatch("2", 100, createOrders2(), createDealRecord2())
	testMatch("2_1", 100, createOrders2_1(), createDealRecord2_1())
	testMatch("3", 100, createOrders3(), createDealRecord3())
	testMatch("4", 100, createOrders4(), createDealRecord4())
	testMatch("4_1", 80, createOrders4_1(), createDealRecord4_1())
	testMatch("5_1", 80, createOrders5_1(), createDealRecord5_1())
	testMatch("5_2", 100, createOrders5_2(), createDealRecord5_2())
	testMatch("5_3a", 90, createOrders5_3(), createDealRecord5_3a())
	testMatch("5_3b", 97, createOrders5_3(), createDealRecord5_3b())
	testMatch("5_3c", 80, createOrders5_3(), createDealRecord5_3c())
	testMatch("5_4", 100, createOrders5_4(), createDealRecord5_4())
	testMatch("6_1", 99, createOrders6(), createDealRecord6_1())
	testMatch("6_2", 97, createOrders6(), createDealRecord6_2())
	testMatch("6_3", 90, createOrders6(), createDealRecord6_3())
	testMatch("6_4", 110, createOrders6(), createDealRecord6_4())
}
