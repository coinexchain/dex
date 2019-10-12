package market

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/market/internal/types"
)

func createOrdersAndMarkets(num int) (Order, []*Order, MarketInfo, []MarketInfo) {
	haveCetAddress, _ := simpleAddr("00001")
	mkInfo := MarketInfo{
		Stock:             stock,
		Money:             money,
		PricePrecision:    9,
		OrderPrecision:    19,
		LastExecutedPrice: sdk.NewDec(987),
	}
	orderInfo := Order{
		Sender: haveCetAddress,
		Price:  sdk.NewDec(19899),
	}

	mkInfos := make([]MarketInfo, num)
	for i := 0; i < num; i++ {
		mk := mkInfo
		mk.Stock = fmt.Sprintf("%s%.3d", stock, i)
		mkInfos[i] = mk
	}

	orderInfos := make([]*Order, num)
	for i := 0; i < num; i++ {
		in := orderInfo
		in.Sequence = uint64(i)
		in.Identify = byte(i)
		orderInfos[i] = &in
	}
	return orderInfo, orderInfos, mkInfo, mkInfos
}

func TestNewGenesis(t *testing.T) {

	orderInfo, orderInfos, _, mkInfos := createOrdersAndMarkets(9)

	state := NewGenesisState(types.DefaultParams(), orderInfos, mkInfos, 876738)
	require.Nil(t, state.Validate())

	require.EqualValues(t, 876738, state.OrderCleanTime)
	require.EqualValues(t, mkInfos, state.MarketInfos)
	for i := 0; i < len(state.Orders); i++ {
		in := orderInfo
		in.Sequence = uint64(i)
		in.Identify = byte(i)
		fmt.Println(orderInfo)
		require.EqualValues(t, in, *state.Orders[i])
	}
}

func TestExportGenesis(t *testing.T) {
	input := prepareMockInput(t, false, false)
	_, orderInfos, _, mkInfos := createOrdersAndMarkets(9)
	state := NewGenesisState(types.DefaultParams(), orderInfos, mkInfos, 876738)
	require.Nil(t, state.Validate())
	InitGenesis(input.ctx, input.mk, state)
	orders := make(map[string]Order)
	for _, order := range orderInfos {
		orders[order.OrderID()] = *order
	}

	exportState := ExportGenesis(input.ctx, input.mk)
	require.Nil(t, exportState.Validate())
	require.EqualValues(t, state.OrderCleanTime, exportState.OrderCleanTime)
	require.EqualValues(t, state.Params, exportState.Params)
	for _, exOrder := range exportState.Orders {
		require.EqualValues(t, orders[exOrder.OrderID()], *exOrder)
	}
	for i, exMarket := range exportState.MarketInfos {
		require.EqualValues(t, mkInfos[i], exMarket)
	}
}

func TestValidateGenesis(t *testing.T) {
	orderInfo, orderInfos, mkInfo, mkInfos := createOrdersAndMarkets(9)

	mkInfo.Stock = fmt.Sprintf("%s%.3d", stock, 1)
	mkInfos = append(mkInfos, mkInfo)
	state := NewGenesisState(types.DefaultParams(), orderInfos, mkInfos, 876738)
	err := state.Validate()
	require.NotNil(t, err)
	require.EqualValues(t, "duplicate market found during market ValidateGenesis", err.Error())

	mkInfos = mkInfos[0 : len(mkInfos)-1]
	orderInfos = append(orderInfos, &orderInfo)
	state = NewGenesisState(types.DefaultParams(), orderInfos, mkInfos, 876738)
	err = state.Validate()
	require.NotNil(t, err)
	require.EqualValues(t, "duplicate order found during market ValidateGenesis", err.Error())

}
