package main

import (
	"fmt"
	"math/rand"
	"strings"
	"encoding/json"
	"encoding/binary"

	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/match"
	"github.com/emirpasic/gods/maps/treemap"
	"golang.org/x/crypto/blake2b"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var LastPrice sdk.Dec
var Keeper *OrderKeeper
var Height int64
var OrderCount int64
var DealCount int64

type Account struct {
	ID int64
}

func (acc Account) String() string {
	return fmt.Sprintf("%d", acc.ID)
}

type Order struct {
	Price sdk.Dec `json:"price"`
	Amount int64 `json:"amount"`
	Height int64 `json:"height"`
	ID int64 `json:"id"`
	Side int `json:"side"`
}

func RandOrder(r *rand.Rand, priceRange int64, amountRange int64) *Order {
	OrderCount++
	side := market.SELL
	if r.Int31n(2) == 1 {
		side = market.BUY
	}
	return &Order{
		Price: sdk.NewDec(r.Int63n(priceRange)),
		Amount: r.Int63n(amountRange),
		Height: Height,
		ID: OrderCount,
		Side: side,
	}
}

func (order *Order) GetPrice() sdk.Dec {
	return order.Price
}
func (order *Order) GetAmount() int64 {
	return order.Amount
}
func (order *Order) GetHeight() int64 {
	return order.Height
}
func (order *Order) GetHash() []byte {
	bz, err := json.Marshal(order)
	if err!=nil {
		panic(err.Error())
	}
	res := blake2b.Sum256(bz)
	return res[:]
}
func (order *Order) GetSide() int {
	return order.Side
}
func (order *Order) GetOwner() match.Account {
	return &Account{ID: order.ID%10000}
}
func (order *Order) String() string {
	return fmt.Sprintf("%d", order.ID)
}
func (order *Order) Key() string {
	priceSlice := market.DecToBigEndianBytes(order.Price)
	idSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(idSlice, uint64(order.ID))
	return string(append(priceSlice, idSlice...))
}

type OrderKeeper struct {
	sellMap   *treemap.Map //map[string]*Order
	buyMap   *treemap.Map //map[string]*Order
}

func removeEntriesRandomly(m *treemap.Map, r *rand.Rand, step int32) {
	keysToRemove := make([]string, 0, 1000)
	iter := m.Iterator()
	iter.Begin()
	for {
		if ok := iter.Next(); !ok {
			break
		}
		if r.Int31n(step)==0 {
			keysToRemove = append(keysToRemove, iter.Key().(string))
		}
	}
	for _, key := range keysToRemove {
		m.Remove(key)
	}
}

func (keeper *OrderKeeper) Size() int {
	return keeper.sellMap.Size() + keeper.buyMap.Size()
}

func (keeper *OrderKeeper) AddOrder(order *Order) {
	if order.Side==market.SELL {
		keeper.sellMap.Put(order.Key(), order)
	} else {
		keeper.buyMap.Put(order.Key(), order)
	}
}

func (keeper *OrderKeeper) RemoveOrder(order *Order) {
	if order.Side==market.SELL {
		ptr, ok := keeper.sellMap.Get(order.Key())
		if !ok || ptr!=order {
			panic("Order not exist")
		}
		keeper.sellMap.Remove(order.Key())
	} else {
		ptr, ok := keeper.buyMap.Get(order.Key())
		if !ok || ptr!=order {
			panic("Order not exist")
		}
		keeper.buyMap.Remove(order.Key())
	}
}

func (keeper *OrderKeeper) GetLowestSell() *Order {
	iter := keeper.sellMap.Iterator()
	iter.Begin()
	ok := iter.Next()
	if !ok {
		panic("Empty Map")
	}
	return iter.Value().(*Order)
}

func (keeper *OrderKeeper) GetHighestBuy() *Order {
	iter := keeper.buyMap.Iterator()
	iter.End()
	ok := iter.Prev()
	if !ok {
		panic("Empty Map")
	}
	return iter.Value().(*Order)
}

func (keeper *OrderKeeper) GetLowestSellUntil(until string) []match.OrderForTrade {
	tail := string([]byte{0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF,0xFF})
	until = until + tail

	res := make([]match.OrderForTrade, 0, 1000)
	iter := keeper.sellMap.Iterator()
	iter.Begin()
	for {
		if ok := iter.Next(); !ok {
			break
		}
		if strings.Compare(iter.Key().(string), until) > 0 {
			break
		}
		res = append(res, match.OrderForTrade(iter.Value().(*Order)))
	}
	return res
}

func (keeper *OrderKeeper) GetHighestBuyUntil(until string) []match.OrderForTrade {
	res := make([]match.OrderForTrade, 0, 1000)
	iter := keeper.buyMap.Iterator()
	iter.End()
	for {
		if ok := iter.Prev(); !ok {
			break
		}
		if strings.Compare(iter.Key().(string), until) < 0 {
			break
		}
		res = append(res, match.OrderForTrade(iter.Value().(*Order)))
	}
	return res
}

func (order *Order) checkPrice(price sdk.Dec) {
	if order.Side==market.SELL && price.LT(order.Price) {
		panic("Can not sell at a lower price")
	}
	if order.Side==market.BUY && price.GT(order.Price) {
		panic("Can not buy at a higher price")
	}
}

func (order *Order) Deal(otherSide match.OrderForTrade, amount int64, price sdk.Dec) {
	otherOrder := otherSide.(*Order)
	order.checkPrice(price)
	otherOrder.checkPrice(price)
	if otherOrder.Amount<amount {
		panic("amount not enough")
	}
	if order.Amount<amount {
		panic("amount not enough")
	}
	otherOrder.Amount -= amount
	order.Amount -= amount
	if otherOrder.Amount == 0 {
		Keeper.RemoveOrder(otherOrder)
	}
	if order.Amount == 0 {
		Keeper.RemoveOrder(order)
	}
	LastPrice = price
	DealCount++
}

func runTest(seed int64, priceRange int64, amountRange int64, delStep int32, liveOrderUpper, liveOrderLower int, heightLimit int) {
	DealCount = 0
	LastPrice = sdk.ZeroDec()
	Keeper = &OrderKeeper{
		sellMap: treemap.NewWithStringComparator(),
		buyMap: treemap.NewWithStringComparator(),
	}
	Height = 0
	OrderCount = 0

	r := rand.New(rand.NewSource(seed))
	for Height:=0; Height<heightLimit; Height++ {
		//if height%1000==0 {
		//	fmt.Printf("Height: %d\n", Height)
		//}
		fmt.Printf("Height: %d OrderCount:%d DealCount:%d\n", Height, Keeper.Size(), DealCount)
		for Keeper.Size() < liveOrderUpper {
			randOrder := RandOrder(r, priceRange, amountRange)
			Keeper.AddOrder(randOrder)
		}
		highBuy := Keeper.GetHighestBuy().Price
		lowSell := Keeper.GetLowestSell().Price
		if highBuy.LT(lowSell) {
			continue
		}
		bidList := Keeper.GetHighestBuyUntil(string(market.DecToBigEndianBytes(lowSell)))
		askList := Keeper.GetLowestSellUntil(string(market.DecToBigEndianBytes(highBuy)))
		ratio := 10
		lowPrice := LastPrice.Mul(sdk.NewDec(int64(100 - ratio))).Quo(sdk.NewDec(100))
		highPrice := LastPrice.Mul(sdk.NewDec(int64(100 + ratio))).Quo(sdk.NewDec(100))

		match.Match(highPrice, LastPrice, lowPrice, bidList, askList)

		highBuy = Keeper.GetHighestBuy().Price
		lowSell = Keeper.GetLowestSell().Price
		if highBuy.GT(lowSell) {
			panic("Still can deal!")
		}

		for Keeper.Size() > liveOrderLower {
			removeEntriesRandomly(Keeper.buyMap, r, delStep)
			removeEntriesRandomly(Keeper.sellMap, r, delStep)
		}
	}
}

func main() {
	  //     seed, priceRange, amountRange, delStep, liveOrderUpper, liveOrderLower, heightLimit
	runTest(    0,    100,      1000,        3,       8000,           6000,           1000)
}
