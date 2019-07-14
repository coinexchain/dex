package client

import (
	asset_types "github.com/coinexchain/dex/modules/asset/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"

	assCli "github.com/coinexchain/dex/modules/asset/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
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

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
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
