package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/spf13/cobra"
)

// IssueTokenCmd will create a issue token tx and sign it with the given key.
func IssueTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-issue [name] [symbol] [totalSupply] [mintable] [burnable] [addrFreezeable] [tokenFreezeable]",
		Short: "Create and sign a issue token tx",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			name := args[0]
			if utf8.RuneCountInString(name) > 32 {
				return fmt.Errorf("issue token name limited to 32 unicode characters")
			}

			symbol := args[1]
			if m, _ := regexp.MatchString("^[a-z][a-z0-9]{1,7}$", symbol); !m {
				return fmt.Errorf("issue token symbol limited to [a-z][a-z0-9]{1,7}")
			}

			totalSupply, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
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
			issueFee := asset.CetCoin(asset.IssueTokenFee)
			if !account.GetCoins().IsAllGTE(issueFee) {
				return fmt.Errorf("address %s doesn't have enough cet to issue token", owner)
			}

			mintable, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}
			burnable, err := strconv.ParseBool(args[4])
			if err != nil {
				return err
			}
			addrFreezeable, err := strconv.ParseBool(args[5])
			if err != nil {
				return err
			}
			tokenFreezeable, err := strconv.ParseBool(args[6])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := asset.NewMsgIssueToken(name, symbol, totalSupply, owner, mintable, burnable, addrFreezeable, tokenFreezeable)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}

	return cmd
}
