package client

import (
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/authx/client/cli"
	"github.com/coinexchain/dex/modules/authx/client/rest"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func NewAuthXModuleClient() types.ModuleClient {
	return AuthXModuleClient{}
}

type AuthXModuleClient struct{}

func (AuthXModuleClient) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, authx.ModuleCdc)
}

func (AuthXModuleClient) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

func (AuthXModuleClient) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetAccountXCmd(cdc)
}
