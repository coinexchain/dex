package match

import (
	"bytes"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BID = 1
const BUY = 1

const ASK = 2
const SELL = 2

type Account interface {
	String() string
}

type OrderForTrade interface {
	GetPrice() sdk.Dec
	GetAmount() int64
	GetHeight() int64
	GetHash() []byte
	GetSide() int
	GetOwner() Account
	Deal(otherSide OrderForTrade, amount int64, price sdk.Dec)
	String() string
}

// match bid order list against ask order list
func Match(highPrice, midPrice, lowPrice sdk.Dec, bidList []OrderForTrade, askList []OrderForTrade) {
	sort.Slice(bidList, func(i, j int) bool {
		return precede(bidList[i], bidList[j])
	})
	sort.Slice(askList, func(i, j int) bool {
		return precede(askList[i], askList[j])
	})
	//for _, order := range bidList {
	//	fmt.Printf("bid %s\n", order.String())
	//}
	//for _, order := range askList {
	//	fmt.Printf("ask %s\n", order.String())
	//}
	for len(bidList) != 0 && len(askList) != 0 && askList[0].GetPrice().LTE(bidList[0].GetPrice()) {
		price := GetExecutionPrice(highPrice, midPrice, lowPrice, append(bidList, askList...))
		//fmt.Printf("Now price is %s\n", price)
		bidList, askList = ExecuteOrderList(price, bidList, askList)
		//if len(bidList) != 0 && len(askList) != 0 {
		//	fmt.Printf("bidList len:%d p:%s askList len:%d p:%s\n",
		//		len(bidList), bidList[0].GetPrice(), len(askList), askList[0].GetPrice())
		//}
	}
}

// return true if a should precede b in a sorted list, i.e. index of a is smaller
func precede(a, b OrderForTrade) bool {
	if (a.GetSide() == ASK && a.GetPrice().LT(b.GetPrice())) || //for ask, lower price has priority
		(a.GetSide() == BID && a.GetPrice().GT(b.GetPrice())) { //for bid, higher price has priority
		return true
	} else if (a.GetSide() == ASK && a.GetPrice().GT(b.GetPrice())) ||
		(a.GetSide() == BID && a.GetPrice().LT(b.GetPrice())) {
		return false
	} else if a.GetHeight() < b.GetHeight() { //lower height has priority
		return true
	} else if a.GetHeight() > b.GetHeight() {
		return false
	} else {
		return bytes.Compare(a.GetHash(), b.GetHash()) < 0
	}
}

// Given price, execute the orders in bidList and askList
func ExecuteOrderList(price sdk.Dec, bidList []OrderForTrade, askList []OrderForTrade) (newBidList []OrderForTrade, newAskList []OrderForTrade) {
	for {
		if len(askList) == 0 || len(bidList) == 0 ||
			bidList[0].GetPrice().LT(price) || askList[0].GetPrice().GT(price) {
			break
		}
		//fmt.Printf("ask0 %s len:%d p:%s a:%d | bid0 %s len:%d p:%s a:%d refP:%s\n",
		//askList[0].GetOwner().String(), len(askList), askList[0].GetPrice().String(), askList[0].GetAmount(),
		//bidList[0].GetOwner().String(), len(bidList), bidList[0].GetPrice().String(), bidList[0].GetAmount(), price.String())
		askList = ExecuteOrder(price, bidList[0], askList)
		if bidList[0].GetAmount() == 0 {
			bidList = bidList[1:]
		}
		if len(askList) == 0 || len(bidList) == 0 ||
			bidList[0].GetPrice().LT(price) || askList[0].GetPrice().GT(price) {
			break
		}
		bidList = ExecuteOrder(price, askList[0], bidList)
		if askList[0].GetAmount() == 0 {
			askList = askList[1:]
		}
	}
	return bidList, askList
}

// Given price, execute currOrder against orders in orderList
func ExecuteOrder(price sdk.Dec, currOrder OrderForTrade, orderList []OrderForTrade) []OrderForTrade {
	firstNonZeroIndex := 0
	for _, otherSide := range orderList {
		if otherSide.GetSide() == BUY {
			if otherSide.GetPrice().LT(price) {
				break
			}
		} else {
			if otherSide.GetPrice().GT(price) {
				break
			}
		}
		minAmount := otherSide.GetAmount()
		if currOrder.GetAmount() < otherSide.GetAmount() {
			minAmount = currOrder.GetAmount()
		}
		currOrder.Deal(otherSide, minAmount, price)
		if otherSide.GetAmount() == 0 {
			firstNonZeroIndex++
		}
		if currOrder.GetAmount() == 0 {
			break
		}
	}

	if firstNonZeroIndex < len(orderList) {
		return orderList[firstNonZeroIndex:]
	}

	return nil
}

type PricePoint struct {
	price                sdk.Dec
	accumulatedAskAmount int64
	askAmount            int64
	accumulatedBidAmount int64
	bidAmount            int64
	executionAmount      int64
	imbalance            int64
	absImbalance         int64
}

func (pp *PricePoint) String() string {
	return fmt.Sprintf("aS:%d\tS:%d\t%s\tB:%d\taB:%d\te:%d\ti:%d", pp.accumulatedAskAmount, pp.askAmount,
		pp.price.String(), pp.bidAmount, pp.accumulatedBidAmount, pp.executionAmount, pp.imbalance)
}

func GetExecutionPrice(highPrice, midPrice, lowPrice sdk.Dec, orders []OrderForTrade) sdk.Dec {
	//for _,order := range orders {
	//	fmt.Printf("%s\n",order.String())
	//}
	ppList := createPricePointList(orders)
	accumulateForPricePointList(ppList)
	//for _, pp := range ppList {
	//	fmt.Printf("%s\n", pp.String())
	//}
	return calculateExecutionPrice(highPrice, midPrice, lowPrice, ppList)
}

// create a slice of PricePoint from orders, and fill three fields: price, askAmount and bidAmount
func createPricePointList(orders []OrderForTrade) []PricePoint {
	ppList := make([]PricePoint, 0, 100)
	ppMap := make(map[string]int)
	for _, order := range orders {
		s := order.GetPrice().String()
		offset, ok := ppMap[s]
		if !ok {
			offset = len(ppList)
			ppMap[s] = offset
			ppList = append(ppList, PricePoint{
				price:                order.GetPrice(),
				accumulatedAskAmount: 0,
				askAmount:            0,
				accumulatedBidAmount: 0,
				bidAmount:            0,
			})
		}
		//fmt.Printf("ppList[%d]: %s\n", offset, ppList[offset].String())
		if order.GetSide() == ASK {
			ppList[offset].askAmount += order.GetAmount()
		} else if order.GetSide() == BID {
			ppList[offset].bidAmount += order.GetAmount()
		}
	}
	return ppList
}

// sort the slice of PricePoint, in price-descending order
// then fill four fields: accumulatedBidAmount, accumulatedAskAmount
func accumulateForPricePointList(ppList []PricePoint) {
	sort.Slice(ppList, func(i, j int) bool {
		// order with higher price come first
		return ppList[i].price.GT(ppList[j].price)
	})
	accumulatedBidAmount := int64(0)
	for i := 0; i < len(ppList); i++ { // buy, scan from top to bottom
		accumulatedBidAmount += ppList[i].bidAmount
		ppList[i].accumulatedBidAmount = accumulatedBidAmount
	}
	accumulatedAskAmount := int64(0)
	for i := len(ppList) - 1; i >= 0; i-- { //sell, scan from bottom to top
		accumulatedAskAmount += ppList[i].askAmount
		ppList[i].accumulatedAskAmount = accumulatedAskAmount
	}
	for i := 0; i < len(ppList); i++ {
		ppList[i].executionAmount = ppList[i].accumulatedAskAmount
		if ppList[i].accumulatedBidAmount < ppList[i].accumulatedAskAmount {
			ppList[i].executionAmount = ppList[i].accumulatedBidAmount
		}
		ppList[i].imbalance = ppList[i].accumulatedBidAmount - ppList[i].accumulatedAskAmount
		ppList[i].absImbalance = ppList[i].imbalance
		if ppList[i].absImbalance < 0 {
			ppList[i].absImbalance = -ppList[i].imbalance
		}
	}
}

// return true if a should precede b in a sorted list, i.e. index of a is smaller
func (pp PricePoint) precede(b PricePoint) bool {
	if pp.executionAmount > b.executionAmount {
		return true
	} else if pp.executionAmount < b.executionAmount {
		return false
	} else if pp.absImbalance < b.absImbalance {
		return true
	} else if pp.absImbalance > b.absImbalance {
		return false
	} else {
		return pp.price.GT(b.price)
	}
}

// calculate execution price from ppList, with respect to highPrice/midPrice/lowPrice
func calculateExecutionPrice(highPrice, midPrice, lowPrice sdk.Dec, ppList []PricePoint) sdk.Dec {
	sort.Slice(ppList, func(i, j int) bool {
		return ppList[i].precede(ppList[j])
	})
	// create a PricePoint list whose every member has the same largest executionAmount
	ppListSameEA := []*PricePoint{&ppList[0]}
	for i := 1; i < len(ppList); i++ {
		if ppList[i].executionAmount == ppList[0].executionAmount {
			ppListSameEA = append(ppListSameEA, &ppList[i])
		} else {
			break
		}
	}
	// if only one price has the largest executionAmount, then use it
	if len(ppListSameEA) == 1 {
		return ppListSameEA[0].price
	}

	sort.Slice(ppListSameEA, func(i, j int) bool {
		return ppListSameEA[i].absImbalance < ppListSameEA[j].absImbalance
	})

	// create a PricePoint list whose every member has the same smallest absImbalance
	ppListSameImbalance := []*PricePoint{ppListSameEA[0]}
	for i := 1; i < len(ppListSameEA); i++ {
		if ppListSameEA[i].absImbalance == ppListSameEA[0].absImbalance {
			ppListSameImbalance = append(ppListSameImbalance, ppListSameEA[i])
		} else {
			break
		}
	}

	// if only one price has the smallest absImbalance, then use it
	if len(ppListSameImbalance) == 1 {
		return ppListSameImbalance[0].price
	}
	return calculateExecutionPriceWithRef(highPrice, midPrice, lowPrice, ppListSameImbalance)
}

// handle the special case: when more than one price has the same largest executionAmount and the
// same smallest absImbalance, we must consider the market pressure
func calculateExecutionPriceWithRef(highPrice, midPrice, lowPrice sdk.Dec, ppListSameImbalance []*PricePoint) sdk.Dec {
	allImbalanceIsNegative := true
	allImbalanceIsPositive := true
	ppWithHighestPrice := ppListSameImbalance[0]
	ppWithLowestPrice := ppListSameImbalance[len(ppListSameImbalance)-1]
	ppWithMiddlePrice := ppListSameImbalance[len(ppListSameImbalance)/2]
	midPriceIsZero := midPrice.Equal(sdk.ZeroDec())
	allPriceLargerThanHigh := ppWithLowestPrice.price.GT(highPrice) && !midPriceIsZero
	allPriceSmallerThanHigh := ppWithHighestPrice.price.LT(highPrice) && !midPriceIsZero
	allPriceLargerThanLow := ppWithLowestPrice.price.GT(lowPrice) && !midPriceIsZero
	allPriceSmallerThanLow := ppWithHighestPrice.price.LT(lowPrice) && !midPriceIsZero
	if midPriceIsZero {
		return ppWithMiddlePrice.price
	}
	for _, pp := range ppListSameImbalance {
		if pp.imbalance < 0 {
			allImbalanceIsPositive = false
		}
		if pp.imbalance > 0 {
			allImbalanceIsNegative = false
		}
	}
	if allImbalanceIsPositive { // with more buyer, we want higher price
		/*
			For scenarios that all the the equivalent surplus amounts are positive, if all the prices are below the reference price plus an upper limit percentage (e.g. 5%), then algorithm uses the highest of the potential equilibrium prices. If all the prices are above the reference price plus an upper limit, use the lowest price; for other cases, use the reference price plus the upper limit.
		*/
		//fmt.Println("allImbalanceIsPositive")
		if allPriceSmallerThanHigh {
			return ppWithHighestPrice.price
		} else if allPriceLargerThanHigh {
			return ppWithLowestPrice.price
		} else {
			return highPrice
		}
	} else if allImbalanceIsNegative { // with more seller, we want lower price
		/*
			Conversely, if market pressure is on the sell side, if all prices are above the reference price minus a lower percentage limit, then the algorithm uses the lowest of the potential prices. If all the price are below the reference price minus the lower percentage limit, use the highest price, otherwise use the reference price minus the lower percentage limit.
		*/
		//fmt.Println("allImbalanceIsNegative")
		if allPriceSmallerThanLow {
			return ppWithHighestPrice.price
		} else if allPriceLargerThanLow {
			return ppWithLowestPrice.price
		} else {
			return lowPrice
		}
	} else {
		/*
			When both positive and negative surplus amounts exists at the lowest, if the reference price falls at / into these prices, the reference price should be chose, otherwise the price closest to the reference price would be chosen.
		*/
		//fmt.Println("Imbalance : Negative and Positive")
		if ppWithHighestPrice.price.LT(midPrice) {
			return ppWithHighestPrice.price
		} else if ppWithLowestPrice.price.GT(midPrice) {
			return ppWithLowestPrice.price
		} else {
			return midPrice
		}
	}
}
