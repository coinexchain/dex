package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

var issueTokenFlags = []string{
	flagName,
	flagSymbol,
	flagTotalSupply,
	flagMintable,
	flagBurnable,
	flagAddrForbiddable,
	flagTokenForbiddable,
	flagTokenURL,
	flagTokenDescription,
	flagTokenIdentity,
}

// get the root tx command of this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	assTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Asset transactions subcommands",
	}

	assTxCmd.AddCommand(client.PostCommands(
		GetCmdIssueToken(types.QuerierRoute, cdc),
		GetCmdTransferOwnership(cdc),
		GetCmdMintToken(cdc),
		GetCmdBurnToken(cdc),
		GetCmdForbidToken(cdc),
		GetCmdUnForbidToken(cdc),
		GetCmdAddTokenWhitelist(cdc),
		GetCmdRemoveTokenWhitelist(cdc),
		GetCmdForbidAddr(cdc),
		GetCmdUnForbidAddr(cdc),
		GetCmdModifyTokenInfo(cdc),
	)...)

	return assTxCmd
}

// GetCmdIssueToken will create a issue token tx and sign.
func GetCmdIssueToken(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token",
		Short: "Create and sign a issue-token tx",
		Long: strings.TrimSpace(
			`Create and sign a issue-token tx, broadcast to nodes.

Example:
$ cetcli tx asset issue-token --name="ABC Token" \
	--symbol="abc" \
	--total-supply=2100000000000000 \
	--mintable=false \
	--burnable=true \
	--addr-forbiddable=false \
	--token-forbiddable=false \
	--url="www.abc.org" \
	--description="token abc is a example token" \
	--identity="552A83BA62F9B1F8" \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenOwner := cliCtx.GetFromAddress()

			msg, err := parseIssueFlags(tokenOwner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(types.NewQueryAssetParams(msg.Symbol))
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryToken)
			if res, _, _ := cliCtx.QueryWithData(route, bz); res != nil {
				return fmt.Errorf("token symbol already exists，please query tokens and issue another symbol")
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagName, "", "issue token name is limited to 32 unicode characters")
	cmd.Flags().String(flagSymbol, "", "issue token symbol is limited to [a-z][a-z0-9]{1,7}")
	cmd.Flags().String(flagTotalSupply, "0", "The amount before boosting should not exceed 90 billion.")
	cmd.Flags().Bool(flagMintable, false, "whether the token could be minted")
	cmd.Flags().Bool(flagBurnable, true, "whether the token could be burned")
	cmd.Flags().Bool(flagAddrForbiddable, false, "whether the token holder address can be forbidden by token owner")
	cmd.Flags().Bool(flagTokenForbiddable, false, "whether the token can be forbidden")
	cmd.Flags().String(flagTokenURL, "", "url of token website")
	cmd.Flags().String(flagTokenDescription, "", "description of token info")
	cmd.Flags().String(flagTokenIdentity, "", "identity of token")

	for _, flag := range issueTokenFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var transferOwnershipFlags = []string{
	flagSymbol,
	flagNewOwner,
}

// GetCmdTransferOwnership will create a transfer token  owner tx and sign.
func GetCmdTransferOwnership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "Create and sign a transfer-ownership tx",
		Long: strings.TrimSpace(
			`Create and sign a transfer-ownership tx, broadcast to nodes.

Example:
$ cetcli tx asset transfer-ownership --symbol="abc" \
	--new-owner=newkey \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseTransferOwnershipFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token`s ownership be transferred")
	cmd.Flags().String(flagNewOwner, "", "who do you want to transfer to ?")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range transferOwnershipFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var mintTokenFlags = []string{
	flagSymbol,
	flagAmount,
}

// GetCmdMintToken will create a mint token tx and sign.
func GetCmdMintToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-token",
		Short: "Create and sign a mint token tx",
		Long: strings.TrimSpace(
			`Create and sign a mint token tx, broadcast to nodes.

Example:
$ cetcli tx asset mint-token --symbol="abc" \
	--amount=10000000000000000 \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseMintTokenFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token will be minted")
	cmd.Flags().String(flagAmount, "0", "the amount of mint")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range mintTokenFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var burnTokenFlags = []string{
	flagSymbol,
	flagAmount,
}

// GetCmdBurnToken will create a burn token tx and sign.
func GetCmdBurnToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-token",
		Short: "Create and sign a burn token tx",
		Long: strings.TrimSpace(
			`Create and sign a burn token tx, broadcast to nodes.

Example:
$ cetcli tx asset burn-token --symbol="abc" \
	--amount=10000000000000000 \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseBurnTokenFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token will be burned")
	cmd.Flags().String(flagAmount, "0", "the amount of burn")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range burnTokenFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var symbolFlags = []string{
	flagSymbol,
}

// GetCmdForbidToken will create a Forbid token tx and sign.
func GetCmdForbidToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forbid-token",
		Short: "Create and sign a forbid token tx",
		Long: strings.TrimSpace(
			`Create and sign a forbid token tx, broadcast to nodes.

Example:
$ cetcli tx asset forbid-token --symbol="abc" \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseForbidTokenFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token will be forbidden")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range symbolFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

// GetCmdUnForbidToken will create a UnForbid token tx and sign.
func GetCmdUnForbidToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unforbid-token",
		Short: "Create and sign a unforbid token tx",
		Long: strings.TrimSpace(
			`Create and sign a unforbid token tx, broadcast to nodes.

Example:
$ cetcli tx asset unforbid-token --symbol="abc" \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseUnForbidTokenFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token will be un forbidden")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range symbolFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var whitelistFlags = []string{
	flagSymbol,
	flagWhitelist,
}

// GetCmdAddTokenWhitelist will create a add token whitelist tx and sign.
func GetCmdAddTokenWhitelist(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-whitelist",
		Short: "Create and sign a add-whitelist tx",
		Long: strings.TrimSpace(
			`Create and sign a add-whitelist tx, broadcast to nodes.
				Multiple addresses separated by commas.

Example:
$ cetcli tx asset add-whitelist --symbol="abc" \
	--whitelist=key,key,key \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseAddWhitelistFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token whitelist be added")
	cmd.Flags().String(flagWhitelist, "", "add token whitelist addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range whitelistFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

// GetCmdRemoveTokenWhitelist will create a remove token whitelist tx and sign.
func GetCmdRemoveTokenWhitelist(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-whitelist",
		Short: "Create and sign a remove-whitelist tx",
		Long: strings.TrimSpace(
			`Create and sign a remove-whitelist tx, broadcast to nodes.
				Multiple addresses separated by commas.

Example:
$ cetcli tx asset remove-whitelist --symbol="abc" \
	--whitelist=key,key,key \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseRemoveWhitelistFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token whitelist be remove")
	cmd.Flags().String(flagWhitelist, "", "remove token whitelist addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range whitelistFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var addressesFlags = []string{
	flagSymbol,
	flagAddresses,
}

// GetCmdForbidAddr will create forbid address tx and sign.
func GetCmdForbidAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "forbid-addr",
		Short: "Create and sign a forbid-addr tx",
		Long: strings.TrimSpace(
			`Create and sign a forbid-addr tx, broadcast to nodes.
				Multiple addresses separated by commas.

Example:
$ cetcli tx asset forbid-addr --symbol="abc" \
	--addresses=key,key,key \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseForbidAddrFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token address be forbidden")
	cmd.Flags().String(flagAddresses, "", "forbid addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range addressesFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

// GetCmdUnForbidAddr will create unforbid address tx and sign.
func GetCmdUnForbidAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unforbid-addr",
		Short: "Create and sign a unforbid-addr tx",
		Long: strings.TrimSpace(
			`Create and sign a unforbid-addr tx, broadcast to nodes.
				Multiple addresses separated by commas.

Example:
$ cetcli tx asset unforbid-addr --symbol="abc" \
	--addresses=key,key,key \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseUnForbidAddrFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token address be un-forbidden")
	cmd.Flags().String(flagAddresses, "", "unforbid addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range addressesFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var modifyTokenURLFlags = []string{
	flagSymbol,
	flagTokenURL,
	flagTokenDescription,
	flagTokenIdentity,
}

// GetCmdModifyTokenInfo will create a modify token info tx and sign.
func GetCmdModifyTokenInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-token-info",
		Short: "Modify token info",
		Long: strings.TrimSpace(
			`Create and sign a modify token info msg, broadcast to nodes.

Example:
$ cetcli tx asset modify-token-info --symbol="abc" \
	--url="www.abc.com" \
	--description="abc example description" \
	--identity="552A83BA62F9B1F8" \
	--from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, err := parseModifyTokenInfoFlags(nil)
			if err != nil {
				return err
			}
			return cliutil.CliRunCommand(cdc, msg)
		},
	}

	cmd.Flags().String(flagSymbol, "", "which token will be modify")
	cmd.Flags().String(flagTokenURL, types.DoNotModifyTokenInfo, "the url of token")
	cmd.Flags().String(flagTokenDescription, types.DoNotModifyTokenInfo, "the description of token")
	cmd.Flags().String(flagTokenIdentity, types.DoNotModifyTokenInfo, "the identity of token")

	_ = cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}
