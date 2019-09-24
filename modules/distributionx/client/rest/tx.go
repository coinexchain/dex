package rest

import (
	"net/http"

	"github.com/coinexchain/dex/client/restutil"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/modules/distributionx/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/distribution/{address}/donates", DonateTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Amount  sdk.Coins    `json:"amount"`
}

func (sr SendReq) New() restutil.RestReq {
	return new(SendReq)
}
func (sr SendReq) GetBaseReq() *rest.BaseReq {
	return &sr.BaseReq
}
func (sr SendReq) GetMsg(r *http.Request, sender sdk.AccAddress) (sdk.Msg, error) {

	from, err := sdk.AccAddressFromBech32(sr.BaseReq.From)
	if err != nil {
		return types.MsgDonateToCommunityPool{}, err
	}
	return types.NewMsgDonateToCommunityPool(from, sr.Amount), nil
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func DonateTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(SendReq))
}
