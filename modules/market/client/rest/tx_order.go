package rest

import (
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/coinexchain/dex/modules/market/client/cli"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

// SendReq defines the properties of a send request's body.
type createOrderReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	OrderType      int          `json:"order_type"`
	TradingPair    string       `json:"trading_pair"`
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

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
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
		req.ExistBlocks = keepers.DefaultGTEOrderLifetime
	}

	creator, err := sdk.AccAddressFromBech32(req.BaseReq.From)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, _, err := queryMarketInfo(cdc, cliCtx, req.TradingPair); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	force := types.GTE
	if !isGTE {
		force = types.IOC
	}

	accRetriever := auth.NewAccountRetriever(cliCtx)
	sequence := req.BaseReq.Sequence
	if sequence == 0 {
		_, sequence, err = accRetriever.GetAccountNumberSequence(creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Don't get sequence from blockchain.")
			return
		}
	}

	msg := types.MsgCreateOrder{
		Sender:         creator,
		Sequence:       sequence,
		TradingPair:    req.TradingPair,
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

	symbols := strings.Split(msg.TradingPair, types.SymbolSeparator)
	userToken := symbols[0]
	if msg.Side == types.BUY {
		userToken = symbols[1]
	}

	account, err := accRetriever.GetAccount(creator)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "No have insufficient cet to create market in blockchain")
		return
	}
	if !account.GetCoins().IsAllGTE(sdk.Coins{sdk.NewCoin(userToken, sdk.NewInt(msg.Quantity))}) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "No have insufficient cet to create market in blockchain")
		return
	}

	utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
}
