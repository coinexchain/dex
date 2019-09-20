package restutil

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type RestReq interface {
	GetBaseReq() *rest.BaseReq
	GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg
}

type RestHandlerBuilder struct {
	cdc     *codec.Codec
	cliCtx  context.CLIContext
	restReq RestReq
}

func NewRestHandlerBuilder(cdc *codec.Codec, cliCtx context.CLIContext, req RestReq) *RestHandlerBuilder {
	return &RestHandlerBuilder{
		cdc:     cdc,
		cliCtx:  cliCtx,
		restReq: req,
	}
}

func (rhb *RestHandlerBuilder) preProc(w http.ResponseWriter, r *http.Request) (sdk.AccAddress, bool) {
	if !rest.ReadRESTReq(w, r, rhb.cdc, rhb.restReq) {
		return nil, false
	}

	baseReq := rhb.restReq.GetBaseReq()
	*baseReq = baseReq.Sanitize()
	if !baseReq.ValidateBasic(w) {
		return nil, false
	}

	sender, err := sdk.AccAddressFromBech32(baseReq.From)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return nil, false
	}

	sequence := baseReq.Sequence
	if sequence == 0 {
		_, sequence, err = auth.NewAccountRetriever(rhb.cliCtx).GetAccountNumberSequence(sender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Can not get sequence from blockchain.")
			return nil, false
		}
	}
	baseReq.Sequence = sequence

	return sender, true
}

func (rhb *RestHandlerBuilder) Build() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sender, ok := rhb.preProc(w, r)
		if !ok {
			return
		}
		msg := rhb.restReq.GetMsg(w, sender)
		if msg == nil {
			return
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, rhb.cliCtx, *rhb.restReq.GetBaseReq(), []sdk.Msg{msg})
	}
}
