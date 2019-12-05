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
	cmd := GetQueryCmd(nil)
	args := []string{
		"trading-pairs",
	}
	cmd.SetArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/market-list", ResultPath)

	args = []string{
		"trading-pair",
		"eth/cet",
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/market-info", ResultPath)
	assert.Equal(t, keepers.QueryMarketParam{TradingPair: "eth/cet"}, ResultParam)

	args = []string{
		"trading-pair",
		"eth-cet",
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, "trading-pair illegal : eth-cet, For example : eth/cet.", err.Error())
	assert.Equal(t, "custom/market/market-info", ResultPath)

	user := "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a"
	orderID := user + "-1025"
	args = []string{
		"order-info",
		orderID,
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/order-info", ResultPath)
	assert.Equal(t, keepers.QueryOrderParam{OrderID: orderID}, ResultParam)

	orderID = user + "=1025"
	args = []string{
		"order-info",
		orderID,
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, "order-id is incorrect", err.Error())
	assert.Equal(t, "custom/market/order-info", ResultPath)

	args = []string{
		"orderbook",
		"eth/cet",
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/orders-in-market", ResultPath)
	assert.Equal(t, keepers.QueryMarketParam{TradingPair: "eth/cet"}, ResultParam)

	args = []string{
		"order-list",
		user,
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, "custom/market/user-order-list", ResultPath)
	assert.Equal(t, keepers.QueryUserOrderList{User: user}, ResultParam)

	args = []string{
		"order-list",
		"coinex1px8alypku5j84qlwzdpy",
	}
	cmd.SetArgs(args)
	err = cmd.Execute()
	assert.Equal(t, "decoding bech32 failed: checksum failed. Expected 026624, got lwzdpy.", err.Error())
	assert.Equal(t, "custom/market/user-order-list", ResultPath)

}
