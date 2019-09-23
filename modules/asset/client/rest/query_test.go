package rest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

const testSymbol = "abc"

func TestQuery(t *testing.T) {
	testQuery(t, "/asset/tokens/abc", "custom/asset/token-info", types.NewQueryAssetParams(testSymbol))
	testQuery(t, "/asset/tokens", "custom/asset/token-list", nil)
	testQuery(t, "/asset/tokens/abc/forbidden/whitelist", "custom/asset/token-whitelist", types.NewQueryWhitelistParams(testSymbol))
	testQuery(t, "/asset/tokens/abc/forbidden/addresses", "custom/asset/addr-forbidden", types.NewQueryForbiddenAddrParams(testSymbol))
	testQuery(t, "/asset/tokens/reserved/symbols", "custom/asset/reserved-symbols", nil)
	testQuery(t, "/asset/parameters", "custom/asset/parameters", nil)
}

func testQuery(t *testing.T, restPath string,
	expectedQueryPath string, expectedParam interface{}) {

	oldRestQuery := restutil.RestQuery
	defer func() {
		restutil.RestQuery = oldRestQuery
	}()

	executed := false
	restutil.RestQuery = func(cdc *codec.Codec, cliCtx context.CLIContext,
		w http.ResponseWriter, r *http.Request, path string, param interface{}, defaultRes []byte) {

		executed = true
		require.Equal(t, expectedQueryPath, path)
		require.Equal(t, expectedParam, param)
	}

	ctx := context.CLIContext{}
	r := mux.NewRouter()
	registerQueryRoutes(ctx, r, nil, "asset")

	url, err := url.Parse(restPath)
	require.NoError(t, err)

	req := &http.Request{Method: "GET", URL: url}
	req = mux.SetURLVars(req, map[string]string{"symbol": testSymbol})
	var match mux.RouteMatch

	ok := r.Match(req, &match)
	require.True(t, ok)
	require.NotNil(t, match.Handler)

	match.Handler.ServeHTTP(nil, req)
	require.True(t, executed)
}
