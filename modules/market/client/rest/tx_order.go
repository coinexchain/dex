package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

// SendReq defines the properties of a send request's body.
type createOrderReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	OrderType      int          `json:"order_type"`
	TradingPair    string       `json:"trading_pair"`
	Identify       int          `json:"identify"`
	PricePrecision int          `json:"price_precision"`
	Price          int64        `json:"price"`
	Quantity       int64        `json:"quantity"`
	Side           int          `json:"side"`
	ExistBlocks    int          `json:"exist_blocks"`
	TimeInForce    int          `json:"time_in_force"`
}

func (req *createOrderReq) New() restutil.RestReq {
	return new(createOrderReq)
}
func (req *createOrderReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *createOrderReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
	msg := types.MsgCreateOrder{
		Sender:         sender,
		TradingPair:    req.TradingPair,
		Identify:       byte(req.Identify),
		OrderType:      byte(req.OrderType),
		PricePrecision: byte(req.PricePrecision),
		Price:          req.Price,
		Quantity:       req.Quantity,
		Side:           byte(req.Side),
		TimeInForce:    types.IOC,
		ExistBlocks:    req.ExistBlocks,
	}
	if r.URL.Path == "/market/gte-orders" {
		msg.TimeInForce = types.GTE
	}
	return msg, nil
}

type cancelOrderReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	OrderID string       `json:"order_id"`
}

func (req *cancelOrderReq) New() restutil.RestReq {
	return new(cancelOrderReq)
}
func (req *cancelOrderReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}
func (req *cancelOrderReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
	msg := &types.MsgCancelOrder{
		OrderID: req.OrderID,
		Sender:  sender,
	}
	return msg, nil
}

func createGTEOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return createOrderAndBroadCast(cdc, cliCtx)
}

func createIOCOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return createOrderAndBroadCast(cdc, cliCtx)
}

func cancelOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req cancelOrderReq
	return restutil.NewRestHandler(cdc, cliCtx, &req)
}

func createOrderAndBroadCast(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req createOrderReq
	return restutil.NewRestHandler(cdc, cliCtx, &req)
}
