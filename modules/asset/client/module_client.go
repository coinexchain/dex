package client

import (
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/asset/client/cli"
	"github.com/coinexchain/dex/modules/asset/client/rest"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func NewAssetModuleClient() types.ModuleClient {
	return AssetModuleClient{}
}

type AssetModuleClient struct{}

func (AssetModuleClient) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, asset.ModuleCdc, asset.StoreKey)
}

func (AssetModuleClient) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

func (AssetModuleClient) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}
