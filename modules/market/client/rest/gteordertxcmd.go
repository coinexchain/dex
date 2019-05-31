package rest

import (
	"github.com/coinexchain/dex/modules/market"
	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"net/http"
)

// SendReq defines the properties of a send request's body.
type createGteOrederReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	OrderType      int          `json:"order_type"`
	Symbol         string       `json:"symbol"`
	PricePrecision int          `json:"price_precision"`
	Price          int64        `json:"price"`
	Quantity       int64        `json:"quantity"`
	Side           int          `json:"side"`
	TimeInForce    int          `json:"time_in_force"`
}

func createGTEOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGteOrederReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := market.MsgCreateOrder{
			Sender:         creator,
			Sequence:       req.BaseReq.Sequence,
			Symbol:         req.Symbol,
			OrderType:      byte(req.OrderType),
			PricePrecision: byte(req.PricePrecision),
			Price:          req.Price,
			Quantity:       req.Quantity,
			Side:           byte(req.Side),
			TimeInForce:    req.TimeInForce,
		}

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{nil})
	}
}
