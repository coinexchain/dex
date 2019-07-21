package rest

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

type BancorInitReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Token     string       `json:"token"`
	MaxSupply string       `json:"max_supply"`
	MaxPrice  string       `json:"max_price"`
}

func bancorInitHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BancorInitReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sequence := req.BaseReq.Sequence
		if sequence == 0 {
			_, sequence, err = auth.NewAccountRetriever(cliCtx).GetAccountNumberSequence(sender)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Can not get sequence from blockchain.")
				return
			}
		}
		req.BaseReq.Sequence = sequence

		var maxPrice sdk.Dec
		types.FillDec(req.MaxPrice, &maxPrice)
		if maxPrice.IsZero() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Max Price is Invalid or Zero")
			return
		}
		maxSupply, ok := sdk.NewIntFromString(req.MaxSupply)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Max Supply is Invalid")
			return
		}
		msg := &types.MsgBancorInit{
			Owner:     sender,
			Token:     req.Token,
			MaxSupply: maxSupply,
			MaxPrice:  maxPrice,
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

type BancorTradeReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Token      string       `json:"token"`
	Amount     string       `json:"amount"`
	IsBuy      bool         `json:"is_buy"`
	MoneyLimit string       `json:"money_limit"`
}

func bancorTradeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BancorTradeReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sequence := req.BaseReq.Sequence
		if sequence == 0 {
			_, sequence, err = auth.NewAccountRetriever(cliCtx).GetAccountNumberSequence(sender)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Can not get sequence from blockchain.")
				return
			}
		}
		req.BaseReq.Sequence = sequence

		amount, err := strconv.ParseInt(req.Amount, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Amount.")
			return
		}

		moneyLimit, err := strconv.ParseInt(req.MoneyLimit, 10, 63)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Money Limit.")
			return
		}

		msg := &types.MsgBancorTrade{
			Sender:     sender,
			Token:      req.Token,
			Amount:     amount,
			IsBuy:      req.IsBuy,
			MoneyLimit: moneyLimit,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
