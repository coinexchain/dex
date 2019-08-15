package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//kafka msg
type MsgBancorCreateForKafka struct {
	Owner              sdk.AccAddress `json:"owner"`
	Stock              string         `json:"stock"`
	Money              string         `json:"money"`
	InitPrice          sdk.Dec        `json:"init_price"`
	MaxSupply          sdk.Int        `json:"max_supply"`
	MaxPrice           sdk.Dec        `json:"max_price"`
	EarliestCancelTime int64          `json:"earliest_cancel_time"`
	BlockHeight        int64          `json:"block_height"`
}

type MsgBancorInfoForKafka struct {
	Owner              sdk.AccAddress `json:"sender"`
	Stock              string         `json:"stock"`
	Money              string         `json:"money"`
	InitPrice          sdk.Dec        `json:"init_price"`
	MaxSupply          sdk.Int        `json:"max_supply"`
	MaxPrice           sdk.Dec        `json:"max_price"`
	Price              sdk.Dec        `json:"price"`
	StockInPool        sdk.Int        `json:"stock_in_pool"`
	MoneyInPool        sdk.Int        `json:"money_in_pool"`
	EarliestCancelTime int64          `json:"earliest_cancel_time"`
	BlockHeight        int64          `json:"block_height"`
}

type MsgBancorTradeInfoForKafka struct {
	Sender      sdk.AccAddress `json:"sender"`
	Stock       string         `json:"stock"`
	Money       string         `json:"money"`
	Amount      int64          `json:"amount"`
	Side        byte           `json:"side"`
	MoneyLimit  int64          `json:"money_limit"`
	TxPrice     sdk.Dec        `json:"transaction_price"`
	BlockHeight int64          `json:"block_height"`
}

type MsgBancorCancelForKafka struct {
	Owner       sdk.AccAddress `json:"owner"`
	Stock       string         `json:"stock"`
	Money       string         `json:"money"`
	BlockHeight int64          `json:"block_height"`
}
