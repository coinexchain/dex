package cli

import (
	"fmt"
	"github.com/spf13/viper"
)

func parseIssueFlags() (*issue, error) {
	issue := &issue{}

	for _, flag := range issueFlags {
		if viper.GetString(flag) == "" {
			return nil, fmt.Errorf("--%s flag is a noop, pls see help : " +
				"$ cetcli tx asset issue-token -h", flag)
		}
	}

	issue.Name = viper.GetString(FlagName)
	issue.Symbol = viper.GetString(FlagSymbol)
	issue.TotalSupply = viper.GetInt64(FlagTotalSupply)
	issue.Mintable = viper.GetBool(FlagMintable)
	issue.Burnable = viper.GetBool(FlagBurnable)
	issue.AddrFreezeable = viper.GetBool(FlagAddrFreezeable)
	issue.TokenFreezeable = viper.GetBool(FlagTokenFreezeable)

	return issue, nil

}
