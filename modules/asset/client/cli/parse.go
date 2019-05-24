package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

func parseIssueFlags(owner sdk.AccAddress) (*asset.MsgIssueToken, error) {
	for _, flag := range issueFlags {
		if viper.GetString(flag) == "" {
			return nil, fmt.Errorf("--%s flag is a noop, pls see help : "+
				"$ cetcli tx asset issue-token -h", flag)
		}
	}

	name := viper.GetString(FlagName)
	symbol := viper.GetString(FlagSymbol)
	totalSupply := viper.GetInt64(FlagTotalSupply)
	mintable := viper.GetBool(FlagMintable)
	burnable := viper.GetBool(FlagBurnable)
	addrFreezeable := viper.GetBool(FlagAddrFreezeable)
	tokenFreezeable := viper.GetBool(FlagTokenFreezeable)

	msg := asset.NewMsgIssueToken(name, symbol, totalSupply, owner,
		mintable, burnable, addrFreezeable, tokenFreezeable)

	return &msg, nil
}
