package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
	"strings"
)

const (
	FlagName           = "name"
	FlagSymbol         = "symbol"
	FlagTotalSupply    = "total-supply"
	FlagMintable       = "mintable"
	FlagBurnable       = "burnable"
	FlagAddrFreezable  = "addr-freezable"
	FlagTokenFreezable = "token-freezable"

	FlagNewOwner = "new-owner"
	FlagAmount   = "amount"
)

var issueTokenFlags = []string{
	FlagName,
	FlagSymbol,
	FlagTotalSupply,
	FlagMintable,
	FlagBurnable,
	FlagAddrFreezable,
	FlagTokenFreezable,
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
	--addr-freezable=false \
	--token-freezable=false \
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
	cmd.Flags().Bool(FlagAddrFreezable, false, " Whether the token holder address can be frozen by token owner")
	cmd.Flags().Bool(FlagTokenFreezable, false, "Whether the token can be frozen")

	cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range issueTokenFlags {
		cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var transferOwnershipFlags = []string{
	FlagSymbol,
	FlagNewOwner,
}

// TransferOwnershipCmd will create a transfer token  owner tx and sign.
func TransferOwnershipCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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

			bz, err := cdc.MarshalJSON(asset.NewQueryAssetParams(msg.Symbol))
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryToken)
			if res, _ := cliCtx.QueryWithData(route, bz); res == nil {
				return fmt.Errorf("transfer invalid token`s ownership")
			}

			// ensure account has enough coins
			account, err := cliCtx.GetAccount(originalOwner)
			if err != nil {
				return err
			}

			issueFee := types.NewCetCoins(asset.TransferOwnershipFee)
			if !account.GetCoins().IsAllGTE(issueFee) {
				return fmt.Errorf("address %s doesn't have enough cet to issue token", originalOwner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token`s ownership be transferred")
	cmd.Flags().String(FlagNewOwner, "", "Who do you want to transfer to ?")

	cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range transferOwnershipFlags {
		cmd.MarkFlagRequired(flag)
	}

	return cmd
}

var mintTokenFlags = []string{
	FlagSymbol,
	FlagAmount,
}

// MintTokenCmd will create a mint token tx and sign.
func MintTokenCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
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

			issueFee := types.NewCetCoins(asset.TransferOwnershipFee)
			if !account.GetCoins().IsAllGTE(issueFee) {
				return fmt.Errorf("address %s doesn't have enough cet to issue token", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagSymbol, "", "Which token will be minted")
	cmd.Flags().String(FlagAmount, "", "The amount of mint")

	cmd.MarkFlagRequired(client.FlagFrom)
	for _, flag := range mintTokenFlags {
		cmd.MarkFlagRequired(flag)
	}

	return cmd
}
