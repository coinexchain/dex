package rest

import (
	"errors"
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
	MaxMoney           string       `json:"max_money"`
	StockPrecision     string       `json:"stock_precision"`
	MaxPrice           string       `json:"max_price"`
	EarliestCancelTime string       `json:"earliest_cancel_time"`
}

var _ restutil.RestReq = (*BancorInitReq)(nil)

func (req *BancorInitReq) New() restutil.RestReq {
	return new(BancorInitReq)
}
func (req *BancorInitReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *BancorInitReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
	maxPrice, err := sdk.NewDecFromStr(req.MaxPrice)
	if err != nil || maxPrice.IsZero() {
		return nil, errors.New("Max Price is Invalid or Zero")
	}
	initPrice, err := sdk.NewDecFromStr(req.InitPrice)
	if err != nil || initPrice.IsNegative() {
		return nil, err
	}
	maxSupply, ok := sdk.NewIntFromString(req.MaxSupply)
	if !ok {
		return nil, errors.New("Max Supply is Invalid")
	}
	maxMoney, ok := sdk.NewIntFromString(req.MaxMoney)
	if !ok {
		return nil, errors.New("Max Money is Invalid")
	}
	time, convertErr := strconv.ParseInt(req.EarliestCancelTime, 10, 64)
	if convertErr != nil {
		return nil, errors.New("Invalid enable cancel time")
	}
	var precision int
	if req.StockPrecision == "" {
		precision = 0
	} else {
		precision, convertErr = strconv.Atoi(req.StockPrecision)
		if convertErr != nil {
			return nil, errors.New("Invalid stock precision")
		}
	}

	return &types.MsgBancorInit{
		Owner:              sender,
		Stock:              req.Stock,
		Money:              req.Money,
		InitPrice:          req.InitPrice,
		MaxSupply:          maxSupply,
		MaxMoney:           maxMoney,
		StockPrecision:     byte(precision),
		MaxPrice:           req.MaxPrice,
		EarliestCancelTime: time,
	}, nil
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

var _ restutil.RestReq = (*BancorTradeReq)(nil)

func (req *BancorTradeReq) New() restutil.RestReq {
	return new(BancorTradeReq)
}
func (req *BancorTradeReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *BancorTradeReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
	amount, err := strconv.ParseInt(req.Amount, 10, 64)
	if err != nil {
		return nil, errors.New("invalid amount")
	}

	moneyLimit, err := strconv.ParseInt(req.MoneyLimit, 10, 64)
	if err != nil {
		return nil, errors.New("invalid money limit")
	}

	return &types.MsgBancorTrade{
		Sender:     sender,
		Stock:      req.Stock,
		Money:      req.Money,
		Amount:     amount,
		IsBuy:      req.IsBuy,
		MoneyLimit: moneyLimit,
	}, nil
}

func bancorTradeHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(BancorTradeReq))
}

type BancorCancelReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Stock   string       `json:"stock"`
	Money   string       `json:"money"`
}

var _ restutil.RestReq = (*BancorCancelReq)(nil)

func (req *BancorCancelReq) New() restutil.RestReq {
	return new(BancorCancelReq)
}
func (req *BancorCancelReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *BancorCancelReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {
	return &types.MsgBancorCancel{
		Owner: sender,
		Stock: req.Stock,
		Money: req.Money,
	}, nil
}

func bancorCancelHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(BancorCancelReq))
}
