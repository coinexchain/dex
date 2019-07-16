package client
//
//import (
//	"github.com/coinexchain/dex/modules/asset"
//	"github.com/coinexchain/dex/modules/asset/client/rest"
//	"github.com/coinexchain/dex/types"
//	"github.com/cosmos/cosmos-sdk/client"
//	"github.com/cosmos/cosmos-sdk/client/context"
//	"github.com/cosmos/cosmos-sdk/codec"
//	"github.com/gorilla/mux"
//	"github.com/spf13/cobra"
//
//	assCli "github.com/coinexchain/dex/modules/asset/client/cli"
//)
//
//type AssetModuleClient struct {
//}
//
//func (mc AssetModuleClient) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
//	rest.RegisterRoutes(ctx, rtr, asset.ModuleCdc, asset.StoreKey)
//}
//
//// get the root tx command of this module
//func (mc AssetModuleClient) GetTxCmd(cdc *codec.Codec) *cobra.Command {
//	assTxCmd := &cobra.Command{
//		Use:   asset.ModuleName,
//		Short: "Asset transactions subcommands",
//	}
//
//	assTxCmd.AddCommand(client.PostCommands(
//		assCli.IssueTokenCmd(asset.QuerierRoute, cdc),
//		assCli.TransferOwnershipCmd(cdc),
//		assCli.MintTokenCmd(cdc),
//		assCli.BurnTokenCmd(cdc),
//		assCli.ForbidTokenCmd(cdc),
//		assCli.UnForbidTokenCmd(cdc),
//		assCli.AddTokenWhitelistCmd(cdc),
//		assCli.RemoveTokenWhitelistCmd(cdc),
//		assCli.ForbidAddrCmd(cdc),
//		assCli.UnForbidAddrCmd(cdc),
//		assCli.ModifyTokenURLCmd(cdc),
//		assCli.ModifyTokenDescriptionCmd(cdc),
//	)...)
//
//	return assTxCmd
//}
//
//// get the root query command of this module
//func (mc AssetModuleClient) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
//	// Group asset queries under a subcommand
//	assQueryCmd := &cobra.Command{
//		Use:   asset.ModuleName,
//		Short: "Querying commands for the asset module",
//	}
//
//	assQueryCmd.AddCommand(client.GetCommands(
//		assCli.GetTokenCmd(asset.QuerierRoute, cdc),
//		assCli.GetTokenListCmd(asset.QuerierRoute, cdc),
//		assCli.GetWhitelistCmd(asset.QuerierRoute, cdc),
//		assCli.GetForbiddenAddrCmd(asset.QuerierRoute, cdc),
//		assCli.GetReservedSymbolsCmd(asset.QuerierRoute, cdc),
//	)...)
//
//	return assQueryCmd
//}
//
//func NewAssetModuleClient() types.ModuleClient {
//	return AssetModuleClient{}
//}
