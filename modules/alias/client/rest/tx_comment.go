package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"net/http"

	"github.com/coinexchain/dex/modules/alias/internal/types"
)

type AliasUpdateReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Alias     string       `json:"alias"`
	IsAdd     bool         `json:"is_add"`
	AsDefault bool         `json:"as_default"`
}

func aliasUpdateHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AliasUpdateReq
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

		msg := &types.MsgAliasUpdate{
			Owner:     sender,
			Alias:     req.Alias,
			IsAdd:     req.IsAdd,
			AsDefault: req.AsDefault,
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
