package dev

import (
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/stakingx"
)

func DefaultParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default-params",
		Short: "Print default params",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			printDefaultParams()
			return nil
		},
	}
	return cmd
}

func printDefaultParams() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module", "Key", "Value", "Type"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
	})

	fillParamsTable(table, "authx", authx.DefaultParams())
	fillParamsTable(table, "bankx", bankx.DefaultParams())
	fillParamsTable(table, "stakingx", stakingx.DefaultParams())
	fillParamsTable(table, "asset", asset.DefaultParams())
	fillParamsTable(table, "market", market.DefaultParams())

	table.Render()
}

func fillParamsTable(table *tablewriter.Table, modname string, obj interface{}) {
	ps := castToParamSet(obj)
	for _, pair := range ps.ParamSetPairs() {
		t := reflect.Indirect(reflect.ValueOf(pair.Value)).Type().Name()
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
		table.Append([]string{modname, string(pair.Key), fmt.Sprintf("%v", v), t})
	}
}

func castToParamSet(obj interface{}) params.ParamSet {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	vpi := vp.Interface()
	return vpi.(params.ParamSet)
}
