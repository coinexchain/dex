package rest

import (
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/client/cli"
	"github.com/coinexchain/dex/modules/market/match"
)

// SendReq defines the properties of a send request's body.
type createOrderReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	OrderType      int          `json:"order_type"`
	Symbol         string       `json:"symbol"`
	PricePrecision int          `json:"price_precision"`
	Price          int64        `json:"price"`
	Quantity       int64        `json:"quantity"`
	Side           int          `json:"side"`
	ExistBlocks    int          `json:"exist_blocks"`
}

type cancelOrderReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	OrderID string       `json:"order_id"`
}

func createGTEOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createOrderAndBroadCast(w, r, cdc, cliCtx, true)
	}
}

func createIOCOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createOrderAndBroadCast(w, r, cdc, cliCtx, false)
	}
}

func cancelOrderHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req cancelOrderReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, _ := sdk.AccAddressFromBech32(req.BaseReq.From)
		msg, err := cli.CheckSenderAndOrderID(sender, req.OrderID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func createOrderAndBroadCast(w http.ResponseWriter, r *http.Request, cdc *codec.Codec, cliCtx context.CLIContext, isGTE bool) {
	var req createOrderReq
	if !rest.ReadRESTReq(w, r, cdc, &req) {
		return
	}

	req.BaseReq = req.BaseReq.Sanitize()
	if !req.BaseReq.ValidateBasic(w) {
		return
	}

	if req.ExistBlocks <= 0 {
		req.ExistBlocks = market.DefaultGTEOrderLifetime
	}

	creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := queryMarketInfo(cdc, cliCtx, req.Symbol); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	force := market.GTE
	if !isGTE {
		force = market.IOC
	}

	sequence, err := cliCtx.GetAccountSequence(creator)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Don't get sequence from blockchain.")
		return
	}

	msg := market.MsgCreateOrder{
		Sender:         creator,
		Sequence:       sequence,
		Symbol:         req.Symbol,
		OrderType:      byte(req.OrderType),
		PricePrecision: byte(req.PricePrecision),
		Price:          req.Price,
		Quantity:       req.Quantity,
		Side:           byte(req.Side),
		TimeInForce:    force,
		ExistBlocks:    req.ExistBlocks,
	}

	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	symbols := strings.Split(msg.Symbol, market.SymbolSeparator)
	userToken := symbols[0]
	if msg.Side == match.BUY {
		userToken = symbols[1]
	}

	account, err := cliCtx.GetAccount(creator)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "No have insufficient cet to create market in blockchain")
		return
	}
	if !account.GetCoins().IsAllGTE(sdk.Coins{sdk.NewCoin(userToken, sdk.NewInt(msg.Quantity))}) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "No have insufficient cet to create market in blockchain")
		return
	}

	clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
}
