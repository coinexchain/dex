package cli

import (
	"fmt"
	"github.com/coinexchain/dex/modules/asset/internal/types"
	"strings"

	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func checkFlags(flags []string, help string) error {
	for _, flag := range flags {
		if viper.GetString(flag) == "" {
			return fmt.Errorf("--%s flag is a noop, please see help : "+help, flag)
		}
	}

	return nil
}

func parseIssueFlags(owner sdk.AccAddress) (*types.MsgIssueToken, error) {
	if err := checkFlags(issueTokenFlags, "$ cetcli tx asset issue-token -h"); err != nil {
		return nil, err
	}

	msg := types.NewMsgIssueToken(
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

func parseTransferOwnershipFlags(orginalOwner sdk.AccAddress) (*types.MsgTransferOwnership, error) {
	if err := checkFlags(transferOwnershipFlags, "$ cetcli tx asset transfer-ownership -h"); err != nil {
		return nil, err
	}

	newOwner, _ := sdk.AccAddressFromBech32(viper.GetString(flagNewOwner))
	msg := types.NewMsgTransferOwnership(
		viper.GetString(flagSymbol),
		orginalOwner,
		newOwner,
	)

	return &msg, nil
}

func parseMintTokenFlags(owner sdk.AccAddress) (*types.MsgMintToken, error) {
	if err := checkFlags(mintTokenFlags, "$ cetcli tx asset mint-token -h"); err != nil {
		return nil, err
	}

	msg := types.NewMsgMintToken(
		viper.GetString(flagSymbol),
		viper.GetInt64(flagAmount),
		owner,
	)

	return &msg, nil
}

func parseBurnTokenFlags(owner sdk.AccAddress) (*types.MsgBurnToken, error) {
	if err := checkFlags(burnTokenFlags, "$ cetcli tx asset burn-token -h"); err != nil {
		return nil, err
	}

	msg := types.NewMsgBurnToken(
		viper.GetString(flagSymbol),
		viper.GetInt64(flagAmount),
		owner,
	)

	return &msg, nil
}

func parseForbidTokenFlags(owner sdk.AccAddress) (*types.MsgForbidToken, error) {
	if err := checkFlags(symbolFlags, "$ cetcli tx asset forbid-token -h"); err != nil {
		return nil, err
	}

	msg := types.NewMsgForbidToken(
		viper.GetString(flagSymbol),
		owner,
	)

	return &msg, nil
}

func parseUnForbidTokenFlags(owner sdk.AccAddress) (*types.MsgUnForbidToken, error) {
	if err := checkFlags(symbolFlags, "$ cetcli tx asset unforbid-token -h"); err != nil {
		return nil, err
	}

	msg := types.NewMsgUnForbidToken(
		viper.GetString(flagSymbol),
		owner,
	)

	return &msg, nil
}

func parseAddWhitelistFlags(owner sdk.AccAddress) (*types.MsgAddTokenWhitelist, error) {
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

	msg := types.NewMsgAddTokenWhitelist(
		viper.GetString(flagSymbol),
		owner,
		whitelist,
	)

	return &msg, nil
}

func parseRemoveWhitelistFlags(owner sdk.AccAddress) (*types.MsgRemoveTokenWhitelist, error) {
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

	msg := types.NewMsgRemoveTokenWhitelist(
		viper.GetString(flagSymbol),
		owner,
		whitelist,
	)

	return &msg, nil
}

func parseForbidAddrFlags(owner sdk.AccAddress) (*types.MsgForbidAddr, error) {
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

	msg := types.NewMsgForbidAddr(
		viper.GetString(flagSymbol),
		owner,
		addresses,
	)

	return &msg, nil
}

func parseUnForbidAddrFlags(owner sdk.AccAddress) (*types.MsgUnForbidAddr, error) {
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

	msg := types.NewMsgUnForbidAddr(
		viper.GetString(flagSymbol),
		owner,
		addresses,
	)

	return &msg, nil
}

func parseModifyTokenURLFlags(owner sdk.AccAddress) (*types.MsgModifyTokenURL, error) {
	if err := checkFlags(modifyTokenURLFlags, "$ cetcli tx asset modify-token-url -h"); err != nil {
		return nil, err
	}

	msg := types.NewMsgModifyTokenURL(
		viper.GetString(flagSymbol),
		viper.GetString(flagTokenURL),
		owner,
	)

	return &msg, nil
}

func parseModifyTokenDescriptionFlags(owner sdk.AccAddress) (*types.MsgModifyTokenDescription, error) {
	if err := checkFlags(modifyTokenDescriptionFlags, "$ cetcli tx asset modify-token-description -h"); err != nil {
		return nil, err
	}

	msg := types.NewMsgModifyTokenDescription(
		viper.GetString(flagSymbol),
		viper.GetString(flagTokenDescription),
		owner,
	)

	return &msg, nil
}
