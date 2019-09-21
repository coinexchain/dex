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
type createMarketReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	Stock          string       `json:"stock"`
	Money          string       `json:"money"`
	PricePrecision int          `json:"price_precision"`
}

func (req *createMarketReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *createMarketReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	msg := types.NewMsgCreateTradingPair(req.Stock, req.Money, sender, byte(req.PricePrecision))
	return msg
}

type cancelMarketReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	TradingPair string       `json:"trading_pair"`
	Time        int64        `json:"time"`
}

func (req *cancelMarketReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *cancelMarketReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	msg := types.MsgCancelTradingPair{
		Sender:        sender,
		TradingPair:   req.TradingPair,
		EffectiveTime: req.Time,
	}
	return msg
}

type modifyPricePrecision struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	TradingPair    string       `json:"trading_pair"`
	PricePrecision int          `json:"price_precision"`
}

func (req *modifyPricePrecision) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *modifyPricePrecision) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	msg := types.MsgModifyPricePrecision{
		Sender:         sender,
		TradingPair:    req.TradingPair,
		PricePrecision: byte(req.PricePrecision),
	}
	return msg
}

func createMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req createMarketReq
	builder := restutil.NewRestHandlerBuilder(cdc, cliCtx, &req)
	return builder.Build()
}

func cancelMarketHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req cancelMarketReq
	builder := restutil.NewRestHandlerBuilder(cdc, cliCtx, &req)
	return builder.Build()
}

func modifyTradingPairPricePrecision(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req modifyPricePrecision
	builder := restutil.NewRestHandlerBuilder(cdc, cliCtx, &req)
	return builder.Build()
}
