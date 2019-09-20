package rest

import (
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
)

type BancorInitReq struct {
	BaseReq            rest.BaseReq `json:"base_req"`
	Stock              string       `json:"stock"`
	Money              string       `json:"money"`
	InitPrice          string       `json:"init_price"`
	MaxSupply          string       `json:"max_supply"`
	MaxPrice           string       `json:"max_price"`
	EarliestCancelTime string       `json:"earliest_cancel_time"`
}

var _ restutil.RestReq = &BancorInitReq{}

func (req *BancorInitReq) New() restutil.RestReq {
	return new(BancorInitReq)
}
func (req *BancorInitReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *BancorInitReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	maxPrice, err := sdk.NewDecFromStr(req.MaxPrice)
	if err != nil || maxPrice.IsZero() {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Max Price is Invalid or Zero")
		return nil
	}
	initPrice, err := sdk.NewDecFromStr(req.InitPrice)
	if err != nil || initPrice.IsNegative() {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Negative init price")
	}
	maxSupply, ok := sdk.NewIntFromString(req.MaxSupply)
	if !ok {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Max Supply is Invalid")
		return nil
	}
	time, converr := strconv.ParseInt(req.EarliestCancelTime, 10, 64)
	if converr != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid enable cancel time")
		return nil
	}

	return &types.MsgBancorInit{
		Owner:              sender,
		Stock:              req.Stock,
		Money:              req.Money,
		InitPrice:          initPrice,
		MaxSupply:          maxSupply,
		MaxPrice:           maxPrice,
		EarliestCancelTime: time,
	}
}

func bancorInitHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(BancorInitReq))
}

type BancorTradeReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Stock      string       `json:"stock"`
	Money      string       `json:"money"`
	Amount     string       `json:"amount"`
	IsBuy      bool         `json:"is_buy"`
	MoneyLimit string       `json:"money_limit"`
}

var _ restutil.RestReq = &BancorTradeReq{}

func (req *BancorTradeReq) New() restutil.RestReq {
	return new(BancorTradeReq)
}
func (req *BancorTradeReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *BancorTradeReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	amount, err := strconv.ParseInt(req.Amount, 10, 64)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Amount.")
		return nil
	}

	moneyLimit, err := strconv.ParseInt(req.MoneyLimit, 10, 64)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Money Limit.")
		return nil
	}

	return &types.MsgBancorTrade{
		Sender:     sender,
		Stock:      req.Stock,
		Money:      req.Money,
		Amount:     amount,
		IsBuy:      req.IsBuy,
		MoneyLimit: moneyLimit,
	}
}

func bancorTradeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(BancorTradeReq))
}

type BancorCancelReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Stock   string       `json:"stock"`
	Money   string       `json:"money"`
}

var _ restutil.RestReq = &BancorCancelReq{}

func (req *BancorCancelReq) New() restutil.RestReq {
	return new(BancorCancelReq)
}
func (req *BancorCancelReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *BancorCancelReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	return &types.MsgBancorCancel{
		Owner: sender,
		Stock: req.Stock,
		Money: req.Money,
	}
}

func bancorCancelHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(BancorCancelReq))
}
