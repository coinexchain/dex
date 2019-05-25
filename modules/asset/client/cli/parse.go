package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

func parseIssueFlags(owner sdk.AccAddress) (*asset.MsgIssueToken, error) {
	for _, flag := range issueTokenFlags {
		if viper.GetString(flag) == "" {
			return nil, fmt.Errorf("--%s flag is a noop, pls see help : "+
				"$ cetcli tx asset issue-token -h", flag)
		}
	}

	msg := asset.NewMsgIssueToken(
		viper.GetString(FlagName),
		viper.GetString(FlagSymbol),
		viper.GetInt64(FlagTotalSupply),
		owner,
		viper.GetBool(FlagMintable),
		viper.GetBool(FlagBurnable),
		viper.GetBool(FlagAddrFreezable),
		viper.GetBool(FlagTokenFreezable))

	return &msg, nil
}
