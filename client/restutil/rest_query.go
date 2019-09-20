package restutil

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

var RestQuery = func(cdc *codec.Codec, cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request,
	query string, param interface{}, defaultRes []byte) {

	var bz []byte
	var err error
	bz = nil
	if param != nil {
		bz, err = cdc.MarshalJSON(param)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
	if !ok {
		return
	}
	res, height, err := cliCtx.QueryWithData(query, bz)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cliCtx = cliCtx.WithHeight(height)
	if len(res) == 0 && len(defaultRes) > 0 {
		rest.PostProcessResponse(w, cliCtx, defaultRes)
	} else {
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
