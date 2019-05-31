package dev

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/app"
)

func ExampleGenesisCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example-genesis",
		Short: "Print example genesis JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			genState := app.NewDefaultGenesisState()
			gsJSON, err := codec.MarshalJSONIndent(cdc, genState)
			if err != nil {
				return err
			}
			fmt.Println(string(gsJSON))
			return nil
		},
	}
	return cmd
}
