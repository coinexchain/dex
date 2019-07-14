package types

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type ModuleClient interface {
	// register rest routes
	RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router)

	// get the root tx command of this module
	GetTxCmd(cdc *codec.Codec) *cobra.Command

	// get the root query command of this module
	GetQueryCmd(cdc *codec.Codec) *cobra.Command
}
