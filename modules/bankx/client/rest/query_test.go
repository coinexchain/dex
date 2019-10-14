package rest

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/restutil"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
)

var ResultParam *keeper.QueryAddrBalances
var ResultPath string

func TestQuery(t *testing.T) {
	restutil.RestQuery = func(cdc *codec.Codec, cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request,
		query string, param interface{}, defaultRes []byte) {
		ResultParam = param.(*keeper.QueryAddrBalances)
		ResultPath = query
	}
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")

	router := mux.NewRouter()
	RegisterRoutes(context.NewCLIContextWithFrom(""), router, nil)
	respWr := restutil.NewResponseWriter4UT()

	req, _ := http.NewRequest("GET", "http://example.com/bank/balances/coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a", nil)
	router.ServeHTTP(respWr, req)
	addr, _ := sdk.AccAddressFromBech32("coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	assert.Equal(t, "custom/bankx/balances", ResultPath)
	assert.Equal(t, &keeper.QueryAddrBalances{Addr: addr}, ResultParam)

	restutil.RestQuery = func(cdc *codec.Codec, cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request,
		query string, param interface{}, defaultRes []byte) {
		ResultParam = nil
		ResultPath = query
	}
	req, _ = http.NewRequest("GET", "http://example.com/bank/parameters", nil)
	router.ServeHTTP(respWr, req)
	assert.Equal(t, "custom/bankx/parameters", ResultPath)
}
