package cli

import (
	"fmt"
	"strings"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
)

const (
	FlagName             = "name"
	FlagSymbol           = "symbol"
	FlagTotalSupply      = "total-supply"
	FlagMintable         = "mintable"
	FlagBurnable         = "burnable"
	FlagAddrForbiddable  = "addr-forbiddable"
	FlagTokenForbiddable = "token-forbiddable"

	FlagNewOwner  = "new-owner"
	FlagAmount    = "amount"
	FlagWhitelist = "whitelist"
	FlagAddresses = "addresses"
)

var issueTokenFlags = []string{
	FlagName,
	FlagSymbol,
	FlagTotalSupply,
	FlagMintable,
	FlagBurnable,
	FlagAddrForbiddable,
	FlagTokenForbiddable,
}

// IssueTokenCmd will create a issue token tx and sign.
func IssueTokenCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
    --from mykey
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			tokenOwner := cliCtx.GetFromAddress()
			msg, err := parseIssueFlags(tokenOwner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(asset.NewQueryAssetParams(msg.Symbol))
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryToken)
			if res, _ := cliCtx.QueryWithData(route, bz); res != nil {
				return fmt.Errorf("token symbol already exists，please query tokens and issue another symbol")
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(tokenOwner)
			if err != nil {
				return err
			}

			issueFee := types.NewCetCoins(asset.IssueTokenFee)
			if !account.GetCoins().IsAllGTE(issueFee) {
				return fmt.Errorf("address %s doesn't have enough cet to issue token", tokenOwner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagName, "", "Issue token name limited to 32 unicode characters")
	cmd.Flags().String(FlagSymbol, "", "Issue token symbol limited to [a-z][a-z0-9]{1,7}")
	cmd.Flags().Int64(FlagTotalSupply, 0, "The total supply for token can have a maximum of "+
		"8 digits of decimal and is boosted by 1e8 in order to store as int64. "+
		"The amount before boosting should not exceed 90 billion.")
	cmd.Flags().Bool(FlagMintable, false, "Whether this token could be minted after the issuing")
	cmd.Flags().Bool(FlagBurnable, true, "Whether this token could be burned")
	cmd.Flags().Bool(FlagAddrForbiddable, false, " Whether the token holder address can be forbidden by token owner")
	cmd.Flags().Bool(FlagTokenForbiddable, false, "Whether the token can be forbidden")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range issueTokenFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var transferOwnershipFlags = []string{
	FlagSymbol,
	FlagNewOwner,
}

// TransferOwnershipCmd will create a transfer token  owner tx and sign.
func TransferOwnershipCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			originalOwner := cliCtx.GetFromAddress()
			msg, err := parseTransferOwnershipFlags(originalOwner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(originalOwner)
			if err != nil {
				return err
			}

			transferFee := types.NewCetCoins(asset.TransferOwnershipFee)
			if !account.GetCoins().IsAllGTE(transferFee) {
				return fmt.Errorf("address %s doesn't have enough cet to transfer ownership", originalOwner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token`s ownership be transferred")
	cmd.Flags().String(FlagNewOwner, "", "Who do you want to transfer to ?")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range transferOwnershipFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var mintTokenFlags = []string{
	FlagSymbol,
	FlagAmount,
}

// MintTokenCmd will create a mint token tx and sign.
func MintTokenCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseMintTokenFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			mintFee := types.NewCetCoins(asset.MintFee)
			if !account.GetCoins().IsAllGTE(mintFee) {
				return fmt.Errorf("address %s doesn't have enough cet to mint token", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token will be minted")
	cmd.Flags().String(FlagAmount, "", "The amount of mint")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range mintTokenFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var burnTokenFlags = []string{
	FlagSymbol,
	FlagAmount,
}

// BurnTokenCmd will create a burn token tx and sign.
func BurnTokenCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseBurnTokenFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			burnFee := types.NewCetCoins(asset.BurnFee)
			if !account.GetCoins().IsAllGTE(burnFee) {
				return fmt.Errorf("address %s doesn't have enough cet to burn token", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token will be burned")
	cmd.Flags().String(FlagAmount, "", "The amount of burn")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range burnTokenFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var symbolFlags = []string{
	FlagSymbol,
}

// ForbidTokenCmd will create a Forbid token tx and sign.
func ForbidTokenCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseForbidTokenFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			forbidFee := types.NewCetCoins(asset.ForbidTokenFee)
			if !account.GetCoins().IsAllGTE(forbidFee) {
				return fmt.Errorf("address %s doesn't have enough cet to forbid token", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token will be forbidden")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range symbolFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

// UnForbidTokenCmd will create a UnForbid token tx and sign.
func UnForbidTokenCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseUnForbidTokenFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			unForbidFee := types.NewCetCoins(asset.UnForbidTokenFee)
			if !account.GetCoins().IsAllGTE(unForbidFee) {
				return fmt.Errorf("address %s doesn't have enough cet to unforbid token", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token will be un forbidden")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range symbolFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var whitelistFlags = []string{
	FlagSymbol,
	FlagWhitelist,
}

// AddTokenWhitelistCmd will create a add token whitelist tx and sign.
func AddTokenWhitelistCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseAddWhitelistFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			fee := types.NewCetCoins(asset.TokenWhitelistAddFee)
			if !account.GetCoins().IsAllGTE(fee) {
				return fmt.Errorf("address %s doesn't have enough cet to add token whitelist", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token whitelist be added")
	cmd.Flags().String(FlagWhitelist, "", "add token whitelist addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range whitelistFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

// RemoveTokenWhitelistCmd will create a remove token whitelist tx and sign.
func RemoveTokenWhitelistCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseRemoveWhitelistFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			fee := types.NewCetCoins(asset.TokenWhitelistRemoveFee)
			if !account.GetCoins().IsAllGTE(fee) {
				return fmt.Errorf("address %s doesn't have enough cet to remove token whitelist", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token whitelist be remove")
	cmd.Flags().String(FlagWhitelist, "", "remove token whitelist addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range whitelistFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var addressesFlags = []string{
	FlagSymbol,
	FlagAddresses,
}

// ForbidAddrCmd will create forbid address tx and sign.
func ForbidAddrCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseForbidAddrFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			fee := types.NewCetCoins(asset.ForbidAddrFee)
			if !account.GetCoins().IsAllGTE(fee) {
				return fmt.Errorf("address %s doesn't have enough cet to forbid addr", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token address be forbidden")
	cmd.Flags().String(FlagAddresses, "", "forbid addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range addressesFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}

// UnForbidAddrCmd will create unforbid address tx and sign.
func UnForbidAddrCmd(cdc *codec.Codec) *cobra.Command {
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
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			owner := cliCtx.GetFromAddress()
			msg, err := parseUnForbidAddrFlags(owner)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			fee := types.NewCetCoins(asset.UnForbidAddrFee)
			if !account.GetCoins().IsAllGTE(fee) {
				return fmt.Errorf("address %s doesn't have enough cet to unforbid addr", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token address be un-forbidden")
	cmd.Flags().String(FlagAddresses, "", "unforbid addresses")

	_ = cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range addressesFlags {
		_ = cmd.MarkFlagRequired(flag)
	}

	return cmd
}
