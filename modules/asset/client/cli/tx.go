package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/client"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
)

const (
	FlagName           = "name"
	FlagSymbol         = "symbol"
	FlagTotalSupply    = "total-supply"
	FlagMintable       = "mintable"
	FlagBurnable       = "burnable"
	FlagAddrFreezable  = "addr-freezable"
	FlagTokenFreezable = "token-freezable"
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

			if utf8.RuneCountInString(msg.Name) > 32 {
				return fmt.Errorf("issue token name limited to 32 unicode characters")
			}

			if m, _ := regexp.MatchString("^[a-z][a-z0-9]{1,7}$", msg.Symbol); !m {
				return fmt.Errorf("issue token symbol limited to [a-z][a-z0-9]{1,7}")
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryTokenList)
			if res, _ := cliCtx.QueryWithData(route, nil); res != nil {
				var tokens []asset.Token
				cdc.MustUnmarshalJSON(res, &tokens)

				for _, t := range tokens {
					if msg.Symbol == t.GetSymbol() {
						return fmt.Errorf("token symbol already exists，pls query tokens and issue another symbol")
					}
				}
			}

			if msg.TotalSupply > asset.MaxTokenAmount {
				return fmt.Errorf("issue token totalSupply limited to 9E18")
			}
			if msg.TotalSupply <= 0 {
				return fmt.Errorf("issue token totalSupply should be positive")
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
