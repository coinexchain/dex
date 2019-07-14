package client

import (
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/asset/client/rest"
	asset_types "github.com/coinexchain/dex/modules/asset/types"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	assCli "github.com/coinexchain/dex/modules/asset/client/cli"
)

type AssetModuleClient struct {
}

func (mc AssetModuleClient) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, asset.ModuleCdc, asset_types.StoreKey)
}

// get the root tx command of this module
func (mc AssetModuleClient) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	assTxCmd := &cobra.Command{
		Use:   asset_types.ModuleName,
		Short: "Asset transactions subcommands",
	}

	assTxCmd.AddCommand(client.PostCommands(
		assCli.IssueTokenCmd(asset_types.QuerierRoute, cdc),
		assCli.TransferOwnershipCmd(cdc),
		assCli.MintTokenCmd(cdc),
		assCli.BurnTokenCmd(cdc),
		assCli.ForbidTokenCmd(cdc),
		assCli.UnForbidTokenCmd(cdc),
		assCli.AddTokenWhitelistCmd(cdc),
		assCli.RemoveTokenWhitelistCmd(cdc),
		assCli.ForbidAddrCmd(cdc),
		assCli.UnForbidAddrCmd(cdc),
		assCli.ModifyTokenURLCmd(cdc),
		assCli.ModifyTokenDescriptionCmd(cdc),
	)...)

	return assTxCmd
}

// get the root query command of this module
func (mc AssetModuleClient) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group asset queries under a subcommand
	assQueryCmd := &cobra.Command{
		Use:   asset_types.ModuleName,
		Short: "Querying commands for the asset module",
	}

	assQueryCmd.AddCommand(client.GetCommands(
		assCli.GetTokenCmd(asset_types.QuerierRoute, cdc),
		assCli.GetTokenListCmd(asset_types.QuerierRoute, cdc),
		assCli.GetWhitelistCmd(asset_types.QuerierRoute, cdc),
		assCli.GetForbiddenAddrCmd(asset_types.QuerierRoute, cdc),
		assCli.GetReservedSymbolsCmd(asset_types.QuerierRoute, cdc),
	)...)

	return assQueryCmd
}

func NewAssetModuleClient() types.ModuleClient {
	return AssetModuleClient{}
}
