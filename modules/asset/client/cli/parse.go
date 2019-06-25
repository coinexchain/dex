package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/asset"
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
		viper.GetString(flagName),
		viper.GetString(flagSymbol),
		viper.GetInt64(flagTotalSupply),
		owner,
		viper.GetBool(flagMintable),
		viper.GetBool(flagBurnable),
		viper.GetBool(flagAddrForbiddable),
		viper.GetBool(flagTokenForbiddable),
		viper.GetString(flagTokenURL),
		viper.GetString(flagTokenDescription),
	)

	return &msg, nil
}

func parseTransferOwnershipFlags(orginalOwner sdk.AccAddress) (*asset.MsgTransferOwnership, error) {
	if err := checkFlags(transferOwnershipFlags, "$ cetcli tx asset transfer-ownership -h"); err != nil {
		return nil, err
	}

	newOwner, _ := sdk.AccAddressFromBech32(viper.GetString(flagNewOwner))
	msg := asset.NewMsgTransferOwnership(
		viper.GetString(flagSymbol),
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
		viper.GetString(flagSymbol),
		viper.GetInt64(flagAmount),
		owner,
	)

	return &msg, nil
}

func parseBurnTokenFlags(owner sdk.AccAddress) (*asset.MsgBurnToken, error) {
	if err := checkFlags(burnTokenFlags, "$ cetcli tx asset burn-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgBurnToken(
		viper.GetString(flagSymbol),
		viper.GetInt64(flagAmount),
		owner,
	)

	return &msg, nil
}

func parseForbidTokenFlags(owner sdk.AccAddress) (*asset.MsgForbidToken, error) {
	if err := checkFlags(symbolFlags, "$ cetcli tx asset forbid-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgForbidToken(
		viper.GetString(flagSymbol),
		owner,
	)

	return &msg, nil
}

func parseUnForbidTokenFlags(owner sdk.AccAddress) (*asset.MsgUnForbidToken, error) {
	if err := checkFlags(symbolFlags, "$ cetcli tx asset unforbid-token -h"); err != nil {
		return nil, err
	}

	msg := asset.NewMsgUnForbidToken(
		viper.GetString(flagSymbol),
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

	str := strings.Split(viper.GetString(flagWhitelist), ",")
	for _, s := range str {
		if addr, err = sdk.AccAddressFromBech32(s); err != nil {
			return nil, err
		}
		whitelist = append(whitelist, addr)
	}

	msg := asset.NewMsgAddTokenWhitelist(
		viper.GetString(flagSymbol),
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

	str := strings.Split(viper.GetString(flagWhitelist), ",")
	for _, s := range str {
		if addr, err = sdk.AccAddressFromBech32(s); err != nil {
			return nil, err
		}
		whitelist = append(whitelist, addr)
	}

	msg := asset.NewMsgRemoveTokenWhitelist(
		viper.GetString(flagSymbol),
		owner,
		whitelist,
	)

	return &msg, nil
}

func parseForbidAddrFlags(owner sdk.AccAddress) (*asset.MsgForbidAddr, error) {
	var addr sdk.AccAddress
	addresses := make([]sdk.AccAddress, 0)
	var err error

	if err := checkFlags(symbolFlags, "$ cetcli tx asset forbid-addr -h"); err != nil {
		return nil, err
	}

	str := strings.Split(viper.GetString(flagAddresses), ",")
	for _, s := range str {
		if addr, err = sdk.AccAddressFromBech32(s); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	msg := asset.NewMsgForbidAddr(
		viper.GetString(flagSymbol),
		owner,
		addresses,
	)

	return &msg, nil
}

func parseUnForbidAddrFlags(owner sdk.AccAddress) (*asset.MsgUnForbidAddr, error) {
	var addr sdk.AccAddress
	addresses := make([]sdk.AccAddress, 0)
	var err error

	if err := checkFlags(symbolFlags, "$ cetcli tx asset unforbid-addr -h"); err != nil {
		return nil, err
	}

	str := strings.Split(viper.GetString(flagAddresses), ",")
	for _, s := range str {
		if addr, err = sdk.AccAddressFromBech32(s); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	msg := asset.NewMsgUnForbidAddr(
		viper.GetString(flagSymbol),
		owner,
		addresses,
	)

	return &msg, nil
}
