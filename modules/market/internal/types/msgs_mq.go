package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CreateMarketInfo struct {
	Stock          string `json:"stock"`
	Money          string `json:"money"`
	PricePrecision byte   `json:"price_precision"`

	// create market info
	Creator      string `json:"creator"`
	CreateHeight int64  `json:"create_height"`
}

type CancelMarketInfo struct {
	Stock string `json:"stock"`
	Money string `json:"money"`

	// del market info
	Deleter string `json:"deleter"`
	DelTime int64  `json:"del_time"`
}

type CreateOrderInfo struct {
	OrderID          string  `json:"order_id"`
	Sender           string  `json:"sender"`
	TradingPair      string  `json:"trading_pair"`
	OrderType        byte    `json:"order_type"`
	Price            sdk.Dec `json:"price"`
	Quantity         int64   `json:"quantity"`
	Side             byte    `json:"side"`
	TimeInForce      int64   `json:"time_in_force"`
	Height           int64   `json:"height"`
	FrozenCommission int64   `json:"frozen_commission"`
	FrozenFeatureFee int64   `json:"frozen_feature_fee"`
	Freeze           int64   `json:"freeze"`
}

type FillOrderInfo struct {
	OrderID     string  `json:"order_id"`
	TradingPair string  `json:"trading_pair"`
	Height      int64   `json:"height"`
	Side        byte    `json:"side"`
	Price       sdk.Dec `json:"price"`

	// These fields will change when order was filled/canceled.
	LeftStock int64   `json:"left_stock"`
	Freeze    int64   `json:"freeze"`
	DealStock int64   `json:"deal_stock"`
	DealMoney int64   `json:"deal_money"`
	CurrStock int64   `json:"curr_stock"`
	CurrMoney int64   `json:"curr_money"`
	FillPrice sdk.Dec `json:"fill_price"`
}

type CancelOrderInfo struct {
	OrderID     string  `json:"order_id"`
	TradingPair string  `json:"trading_pair"`
	Height      int64   `json:"height"`
	Side        byte    `json:"side"`
	Price       sdk.Dec `json:"price"`

	// Del infos
	DelReason string `json:"del_reason"`

	// Fields of amount
	UsedCommission int64 `json:"used_commission"`
	UsedFeatureFee int64 `json:"used_feature_fee"`
	LeftStock      int64 `json:"left_stock"`
	RemainAmount   int64 `json:"remain_amount"`
	DealStock      int64 `json:"deal_stock"`
	DealMoney      int64 `json:"deal_money"`
}

type ModifyPricePrecisionInfo struct {
	Sender            string `json:"sender"`
	TradingPair       string `json:"trading_pair"`
	OldPricePrecision byte   `json:"old_price_precision"`
	NewPricePrecision byte   `json:"new_price_precision"`
}
