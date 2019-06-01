package cli

import (
	"fmt"

	"github.com/coinexchain/dex/modules/asset"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

func checkFlags(flags []string, help string) error {
	for _, flag := range flags {
		if viper.GetString(flag) == "" {
			return fmt.Errorf("--%s flag is a noop, please see help : "+help, flag)
		}
	}

	return nil
}

func parseIssueFlags(owner sdk.AccAddress) (*asset.MsgIssueToken, error) {
	if err := checkFlags(issueTokenFlags, "$ cetcli tx asset issue-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgIssueToken(
		viper.GetString(FlagName),
		viper.GetString(FlagSymbol),
		viper.GetInt64(FlagTotalSupply),
		owner,
		viper.GetBool(FlagMintable),
		viper.GetBool(FlagBurnable),
		viper.GetBool(FlagAddrForbiddable),
		viper.GetBool(FlagTokenForbiddable))

	return &msg, nil
}

func parseTransferOwnershipFlags(orginalOwner sdk.AccAddress) (*asset.MsgTransferOwnership, error) {
	if err := checkFlags(transferOwnershipFlags, "$ cetcli tx asset transfer-ownership -h"); err != nil {
		return nil, err
	}

	newOwner, _ := sdk.AccAddressFromBech32(viper.GetString(FlagNewOwner))
	msg := asset.NewMsgTransferOwnership(
		viper.GetString(FlagSymbol),
		orginalOwner,
		newOwner,
	)

	return &msg, nil
}

func parseMintTokenFlags(owner sdk.AccAddress) (*asset.MsgMintToken, error) {
	if err := checkFlags(mintTokenFlags, "$ cetcli tx asset mint-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgMintToken(
		viper.GetString(FlagSymbol),
		viper.GetInt64(FlagAmount),
		owner,
	)

	return &msg, nil
}

func parseBurnTokenFlags(owner sdk.AccAddress) (*asset.MsgBurnToken, error) {
	if err := checkFlags(burnTokenFlags, "$ cetcli tx asset burn-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgBurnToken(
		viper.GetString(FlagSymbol),
		viper.GetInt64(FlagAmount),
		owner,
	)

	return &msg, nil
}
