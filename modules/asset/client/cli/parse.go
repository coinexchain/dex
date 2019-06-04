package cli

import (
	"fmt"
	"strings"

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

func parseForbidTokenFlags(owner sdk.AccAddress) (*asset.MsgForbidToken, error) {
	if err := checkFlags(symbolFlags, "$ cetcli tx asset forbid-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgForbidToken(
		viper.GetString(FlagSymbol),
		owner,
	)

	return &msg, nil
}

func parseUnForbidTokenFlags(owner sdk.AccAddress) (*asset.MsgUnForbidToken, error) {
	if err := checkFlags(symbolFlags, "$ cetcli tx asset unforbid-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgUnForbidToken(
		viper.GetString(FlagSymbol),
		owner,
	)

	return &msg, nil
}

func parseAddWhitelistFlags(owner sdk.AccAddress) (*asset.MsgAddTokenWhitelist, error) {
	var addr sdk.AccAddress
	whitelist := make([]sdk.AccAddress, 0)
	var err error

	if err := checkFlags(symbolFlags, "$ cetcli tx asset add-whitelist -h"); err != nil {
		return nil, err
	}

	str := strings.Split(viper.GetString(FlagWhitelist), ",")
	for _, s := range str {
		if addr, err = sdk.AccAddressFromBech32(s); err != nil {
			return nil, err
		}
		whitelist = append(whitelist, addr)
	}

	msg := asset.NewMsgAddTokenWhitelist(
		viper.GetString(FlagSymbol),
		owner,
		whitelist,
	)

	return &msg, nil
}

func parseRemoveWhitelistFlags(owner sdk.AccAddress) (*asset.MsgRemoveTokenWhitelist, error) {
	var addr sdk.AccAddress
	whitelist := make([]sdk.AccAddress, 0)
	var err error

	if err := checkFlags(symbolFlags, "$ cetcli tx asset remove-whitelist -h"); err != nil {
		return nil, err
	}

	str := strings.Split(viper.GetString(FlagWhitelist), ",")
	for _, s := range str {
		if addr, err = sdk.AccAddressFromBech32(s); err != nil {
			return nil, err
		}
		whitelist = append(whitelist, addr)
	}

	msg := asset.NewMsgRemoveTokenWhitelist(
		viper.GetString(FlagSymbol),
		owner,
		whitelist,
	)

	return &msg, nil
}
