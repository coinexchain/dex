package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/client/restutil"

	"github.com/cosmos/cosmos-sdk/codec"
)

var testAddr = "coinex12kcupm2x8fw0gglgcz8850kw0k2kx0ff8sr3rn"

func TestDonateTxRequestHandlerFn(t *testing.T) {
	oldFactory := restutil.NewRestHandler
	defer func() {
		restutil.NewRestHandler = oldFactory
	}()

	executed := false
	restutil.NewRestHandler = func(cdc *codec.Codec, cliCtx context.CLIContext, req restutil.RestReq) http.HandlerFunc {
		return func(http.ResponseWriter, *http.Request) {
			executed = true
			require.Equal(t, "*rest.SendReq", fmt.Sprintf("%T", req))
		}
	}

	ctx := context.CLIContext{}
	r := mux.NewRouter()
	RegisterRoutes(ctx, r, nil)

	url, err := url.Parse(fmt.Sprintf("/distribution/%s/donates", testAddr))
	require.NoError(t, err)

	req := &http.Request{Method: "POST", URL: url}
	req = mux.SetURLVars(req, map[string]string{"address": testAddr})
	var match mux.RouteMatch

	ok := r.Match(req, &match)
	require.True(t, ok)
	require.NotNil(t, match.Handler)

	match.Handler.ServeHTTP(nil, req)
	require.True(t, executed)
}
