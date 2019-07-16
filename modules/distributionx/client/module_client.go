package client

import (
	"github.com/coinexchain/dex/modules/distributionx"
	"github.com/coinexchain/dex/modules/distributionx/client/cli"
	"github.com/coinexchain/dex/modules/distributionx/client/rest"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func NewDistrXModuleClient() types.ModuleClient {
	return DistrXModuleClient{}
}

type DistrXModuleClient struct{}

func (DistrXModuleClient) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, distributionx.ModuleCdc)
}

func (DistrXModuleClient) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.DonateTxCmd(cdc)
}

func (DistrXModuleClient) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}
