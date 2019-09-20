package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/alias/internal/types"
)

type AliasUpdateReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Alias     string       `json:"alias"`
	IsAdd     bool         `json:"is_add"`
	AsDefault bool         `json:"as_default"`
}

var _ restutil.RestReq = &AliasUpdateReq{}

func (req *AliasUpdateReq) New() restutil.RestReq {
	return new(AliasUpdateReq)
}
func (req *AliasUpdateReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *AliasUpdateReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	return &types.MsgAliasUpdate{
		Owner:     sender,
		Alias:     req.Alias,
		IsAdd:     req.IsAdd,
		AsDefault: req.AsDefault,
	}
}

func aliasUpdateHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(AliasUpdateReq))
}
