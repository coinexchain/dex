package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/client/restutil"
)

func TestTx(t *testing.T) {
	//testTx(t, "/asset/tokens", "*rest.issueReq")
	testTx(t, "/asset/tokens/abc/ownerships", "*rest.transferOwnerReq")
	testTx(t, "/asset/tokens/abc/mints", "*rest.mintTokenReq")
	testTx(t, "/asset/tokens/abc/burns", "*rest.burnTokenReq")
	testTx(t, "/asset/tokens/abc/forbids", "*rest.forbidTokenReq")
	testTx(t, "/asset/tokens/abc/unforbids", "*rest.unForbidTokenReq")
	testTx(t, "/asset/tokens/abc/forbidden/whitelist", "*rest.addWhiteListReq")
	testTx(t, "/asset/tokens/abc/unforbidden/whitelist", "*rest.removeWhiteListReq")
	testTx(t, "/asset/tokens/abc/forbidden/addresses", "*rest.forbidAddrReq")
	testTx(t, "/asset/tokens/abc/unforbidden/addresses", "*rest.unforbidAddrReq")
	testTx(t, "/asset/tokens/abc/infos", "*rest.modifyTokenInfoReq")
}

func testTx(t *testing.T, restPath string, expectedReqType string) {
	oldFactory := restutil.NewRestHandler
	defer func() {
		restutil.NewRestHandler = oldFactory
	}()

	executed := false
	restutil.NewRestHandler = func(cdc *codec.Codec, cliCtx context.CLIContext, req restutil.RestReq) http.HandlerFunc {
		return func(http.ResponseWriter, *http.Request) {
			executed = true
			require.Equal(t, expectedReqType, fmt.Sprintf("%T", req))
		}
	}

	ctx := context.CLIContext{}
	r := mux.NewRouter()
	registerTXRoutes(ctx, r, nil)

	url, err := url.Parse(restPath)
	require.NoError(t, err)

	req := &http.Request{Method: "POST", URL: url}
	req = mux.SetURLVars(req, map[string]string{"symbol": testSymbol})
	var match mux.RouteMatch

	ok := r.Match(req, &match)
	require.True(t, ok)
	require.NotNil(t, match.Handler)

	match.Handler.ServeHTTP(nil, req)
	require.True(t, executed)
}
