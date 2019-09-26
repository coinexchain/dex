package rest

import (
	"net/http"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
)

var ResultParam interface{}
var ResultPath string

func RestQueryForTest(cdc *codec.Codec, cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request,
	query string, param interface{}, defaultRes []byte) {
	ResultParam = param
	ResultPath = query
}

func TestQuery(t *testing.T) {
	restutil.RestQuery = RestQueryForTest
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")

	router := mux.NewRouter()
	registerQueryRoutes(context.NewCLIContextWithFrom(""), router, nil)
	respWr := restutil.NewResponseWriter4UT()

	req, _ := http.NewRequest("GET", "http://example.com/market/trading-pairs/etc/cet", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/market/market-info", ResultPath)
	assert.Equal(t, keepers.QueryMarketParam{
		TradingPair: "etc/cet",
	}, ResultParam)

	req, _ = http.NewRequest("GET", "http://example.com/market/exist-trading-pairs", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/market/market-list", ResultPath)

	req, _ = http.NewRequest("GET", "http://example.com/market/orders/coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/market/order-info", ResultPath)
	assert.Equal(t, keepers.QueryOrderParam{
		OrderID: "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025",
	}, ResultParam)

	req, _ = http.NewRequest("GET", "http://example.com/market/orders/account/coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/market/user-order-list", ResultPath)
	assert.Equal(t, keepers.QueryUserOrderList{
		User: "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a",
	}, ResultParam)

	req, _ = http.NewRequest("GET", "http://example.com/market/parameters", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/market/parameters", ResultPath)
}
