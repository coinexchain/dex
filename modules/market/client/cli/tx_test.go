package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

var ResultMsg cliutil.MsgWithAccAddress

func CliRunCommandForTest(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
	ResultMsg = msg
	return nil
}

func CliQueryNull(cdc *codec.Codec, path string, param interface{}) error {
	return nil
}

func TestCmd(t *testing.T) {
	cliutil.CliRunCommand = CliRunCommandForTest
	cliutil.CliQuery = CliQueryNull

	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	cmd := GetTxCmd(nil)

	args := []string{
		"create-trading-pair",
		"--stock=eth",
		"--money=cet",
		"--price-precision=8",
		"--order-precision=8",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, &types.MsgCreateTradingPair{
		Stock:          "eth",
		Money:          "cet",
		PricePrecision: byte(8),
		OrderPrecision: byte(8),
	}, ResultMsg)

	args = []string{
		"create-trading-pair",
		"--stock=eth",
		"--money=cet",
		"--price-precision=800",
		"--order-precision=8",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)

	args = []string{
		"create-trading-pair",
		"--money=cet",
		"--price-precision=800",
		"--order-precision=8",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, "tx flag is error, please see help : $ cetcli tx market createmarket -h", err.Error())

	args = []string{
		"cancel-trading-pair",
		"--trading-pair=etc/cet",
		"--time=1234567",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, &types.MsgCancelTradingPair{
		EffectiveTime: 1234567,
		TradingPair:   "etc/cet",
	}, ResultMsg)

	args = []string{
		"modify-price-precision",
		"--trading-pair=etc/cet",
		"--price-precision=9",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, &types.MsgModifyPricePrecision{
		TradingPair:    "etc/cet",
		PricePrecision: byte(9),
	}, ResultMsg)

	args = []string{
		"create-gte-order",
		"--trading-pair=btc/cet",
		"--price-precision=9",
		"--order-type=2",
		"--price=520",
		"--quantity=12345678",
		"--side=1",
		"--price-precision=10",
		"--identify=0",
		"--blocks=40000",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, &types.MsgCreateOrder{
		Identify:       0,
		TradingPair:    "btc/cet",
		OrderType:      types.LIMIT,
		Side:           types.BUY,
		Price:          520,
		PricePrecision: 10,
		Quantity:       12345678,
		ExistBlocks:    40000,
		TimeInForce:    types.GTE,
	}, ResultMsg)

	args = []string{
		"create-ioc-order",
		"--trading-pair=btc/cet",
		"--price-precision=9",
		"--order-type=2",
		"--price=520",
		"--quantity=12345678",
		"--side=1",
		"--price-precision=10",
		"--identify=1",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, &types.MsgCreateOrder{
		Identify:       1,
		TradingPair:    "btc/cet",
		OrderType:      types.LIMIT,
		Side:           types.BUY,
		Price:          520,
		PricePrecision: 10,
		Quantity:       12345678,
		ExistBlocks:    0,
		TimeInForce:    types.IOC,
	}, ResultMsg)

	args = []string{
		"cancel-order",
		"--order-id=coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	assert.Equal(t, &types.MsgCancelOrder{
		OrderID: "coinex1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025",
	}, ResultMsg)
}
