package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/types"
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
	FlagName            = "name"
	FlagSymbol          = "symbol"
	FlagTotalSupply     = "total-supply"
	FlagMintable        = "mintable"
	FlagBurnable        = "burnable"
	FlagAddrFreezeable  = "addr-freezeable"
	FlagTokenFreezeable = "token-freezeable"
)

type issue struct {
	Name            string
	Symbol          string
	TotalSupply     int64
	Mintable        bool
	Burnable        bool
	AddrFreezeable  bool
	TokenFreezeable bool
}

var issueFlags = []string{
	FlagName,
	FlagSymbol,
	FlagTotalSupply,
	FlagMintable,
	FlagBurnable,
	FlagAddrFreezeable,
	FlagTokenFreezeable,
}

// IssueTokenCmd will create a issue token tx and sign.
func IssueTokenCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-token",
		Short: "Create and sign a issue-token tx",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create and sign a issue-token tx, broadcast to nodes.

Example:
$ cetcli tx asset issue-token --name="my first token" \
	--symbol="token1" \
	--total-supply=2100000000000000 \
	--mintable=false \
	--burnable=false \
	--addr-freezeable=false \
	--token-freezeable=false --from mykey
`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			issue, err := parseIssueFlags()
			if err != nil {
				return err
			}

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			name := issue.Name
			if utf8.RuneCountInString(name) > 32 {
				return fmt.Errorf("issue token name limited to 32 unicode characters")
			}

			symbol := issue.Symbol
			if m, _ := regexp.MatchString("^[a-z][a-z0-9]{1,7}$", symbol); !m {
				return fmt.Errorf("issue token symbol limited to [a-z][a-z0-9]{1,7}")
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, asset.QueryTokenList)
			if res, _ := cliCtx.QueryWithData(route, nil); res != nil {
				var tokens []asset.Token
				cdc.MustUnmarshalJSON(res, &tokens)

				for _, t := range tokens {
					if symbol == t.GetSymbol() {
						return fmt.Errorf("token symbol already exists，pls query tokens and issue another symbol")
					}
				}
			}

			totalSupply := issue.TotalSupply
			if totalSupply > asset.MaxTokenAmount {
				return fmt.Errorf("issue token totalSupply limited to 9E18")
			}
			if totalSupply < 0 {
				return fmt.Errorf("issue token totalSupply should be a positive")
			}

			owner := cliCtx.GetFromAddress()
			account, err := cliCtx.GetAccount(owner)
			if err != nil {
				return err
			}

			// ensure account has enough coins
			issueFee := types.NewCetCoins(asset.IssueTokenFee)
			if !account.GetCoins().IsAllGTE(issueFee) {
				return fmt.Errorf("address %s doesn't have enough cet to issue token", owner)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := asset.NewMsgIssueToken(name, symbol, totalSupply, owner, issue.Mintable,
				issue.Burnable, issue.AddrFreezeable, issue.TokenFreezeable)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	cmd.Flags().String(FlagName, "", "issue token name limited to 32 unicode characters")
	cmd.Flags().String(FlagSymbol, "", "issue token symbol limited to [a-z][a-z0-9]{1,7}")
	cmd.Flags().String(FlagTotalSupply, "", "issue token totalSupply limited to 9E18")
	cmd.Flags().String(FlagMintable, "", "whether this token could be minted after the issuing")
	cmd.Flags().String(FlagBurnable, "", "whether this token could be burned")
	cmd.Flags().String(FlagAddrFreezeable, "", " whether could freeze some addresses to forbid transaction")
	cmd.Flags().String(FlagTokenFreezeable, "", "whether token could be global freeze")

	return cmd
}
