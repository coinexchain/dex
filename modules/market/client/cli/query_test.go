package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
)

var ResultParam interface{}
var ResultPath string

func CliQueryForTest(cdc *codec.Codec, path string, param interface{}) error {
	ResultParam = param
	ResultPath = path
	return nil
}

func TestQuery(t *testing.T) {
	cliutil.CliQuery = CliQueryForTest

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")

	cmd := QueryMarketListCmd(nil)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/market-list", ResultPath)

	args := []string{
		"eth/cet",
	}
	cmd = QueryMarketCmd(nil)
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/market-info", ResultPath)
	assert.Equal(t, keepers.QueryMarketParam{TradingPair: "eth/cet"}, ResultParam)

	user := "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a"
	orderID := user + "-1025"
	args = []string{
		orderID,
	}
	cmd = QueryOrderCmd(nil)
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/order-info", ResultPath)
	assert.Equal(t, keepers.QueryOrderParam{OrderID: orderID}, ResultParam)

	args = []string{
		user,
	}
	cmd = QueryUserOrderList(nil)
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/user-order-list", ResultPath)
	assert.Equal(t, keepers.QueryUserOrderList{User: user}, ResultParam)

}
