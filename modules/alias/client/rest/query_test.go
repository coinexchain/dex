package rest

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/alias/internal/keepers"
)

var ResultParam *keepers.QueryAliasInfoParam
var ResultPath string

func RestQueryForTest(cdc *codec.Codec, cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request,
	query string, param interface{}, defaultRes []byte) {
	ResultParam = param.(*keepers.QueryAliasInfoParam)
	ResultPath = query
}

func TestQuery(t *testing.T) {
	restutil.RestQuery = RestQueryForTest
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")

	router := mux.NewRouter()
	registerQueryRoutes(context.NewCLIContextWithFrom(""), router, nil)
	respWr := restutil.NewResponseWriter4UT()

	req, _ := http.NewRequest("GET", "http://example.com/alias/address-of-alias/super_boy", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/alias/alias-info", ResultPath)
	assert.Equal(t, &keepers.QueryAliasInfoParam{
		Alias:   "super_boy",
		QueryOp: keepers.GetAddressFromAlias,
	}, ResultParam)

	req, _ = http.NewRequest("GET", "http://example.com/alias/aliases-of-address/coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a", nil)
	router.ServeHTTP(respWr, req)
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	assert.Equal(t, "custom/alias/alias-info", ResultPath)
	assert.Equal(t, &keepers.QueryAliasInfoParam{
		Owner:   addr,
		QueryOp: keepers.ListAliasOfAccount,
	}, ResultParam)

	req, _ = http.NewRequest("GET", "http://example.com/alias/aliases-of-address/coinex1px8alypku5j84qlwzdpynhn4ny", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/alias/alias-info", ResultPath)
	correct := `{"error":"decoding bech32 failed: checksum failed. Expected ffw0y2, got nhn4ny."}`
	assert.Equal(t, correct, string(respWr.GetBody()))
}
