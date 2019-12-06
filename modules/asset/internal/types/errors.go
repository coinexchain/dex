package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceAsset sdk.CodespaceType = ModuleName

	// 501 ~ 599
	CodeInvalidTokenName             sdk.CodeType = 501
	CodeInvalidTokenSymbol           sdk.CodeType = 502
	CodeInvalidTokenSupply           sdk.CodeType = 503
	CodeInvalidTokenOwner            sdk.CodeType = 504
	CodeInvalidTokenMintAmt          sdk.CodeType = 505
	CodeInvalidTokenBurnAmt          sdk.CodeType = 506
	CodeInvalidTokenForbidden        sdk.CodeType = 507
	CodeInvalidTokenUnForbidden      sdk.CodeType = 508
	CodeInvalidTokenWhitelist        sdk.CodeType = 509
	CodeInvalidForbiddenAddress      sdk.CodeType = 510
	CodeInvalidTokenURL              sdk.CodeType = 511
	CodeInvalidTokenDescription      sdk.CodeType = 512
	CodeTokenNotFound                sdk.CodeType = 513
	CodeDuplicateTokenSymbol         sdk.CodeType = 514
	CodeTransferSelfTokenOwner       sdk.CodeType = 515
	CodeNilTokenOwner                sdk.CodeType = 516
	CodeNeedTokenOwner               sdk.CodeType = 517
	CodeInvalidIssueOwner            sdk.CodeType = 518
	CodeTokenMintNotSupported        sdk.CodeType = 519
	CodeTokenBurnNotSupported        sdk.CodeType = 520
	CodeTokenForbiddenNotSupported   sdk.CodeType = 521
	CodeAddressForbiddenNotSupported sdk.CodeType = 522
	CodeNilTokenWhitelist            sdk.CodeType = 523
	CodeNilForbiddenAddress          sdk.CodeType = 524
	CodeInvalidTokenIdentity         sdk.CodeType = 526
	CodeInvalidSendLockAmt           sdk.CodeType = 527
	CodeAddrInBlackList              sdk.CodeType = 528
	CodeNilTokenIdentity             sdk.CodeType = 529
	CodeTokenOwnerSelfForbidden      sdk.CodeType = 530
	CodeInvalidTokenInfo             sdk.CodeType = 531
	CodeTokenInfoSealed              sdk.CodeType = 532
)

func ErrInvalidTokenName(name string) sdk.Error {
	msg := fmt.Sprintf("invalid name %s ： token name is limited to 32 unicode characters", name)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenName, msg)
}
func ErrInvalidTokenSymbol(symbol string) sdk.Error {
	msg := fmt.Sprintf("invalid symbol %s : token symbol not match with [a-z][a-z0-9]{1,15}", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenSymbol, msg)
}
func ErrInvalidTokenSupply(amt string) sdk.Error {
	msg := fmt.Sprintf("invalid supply %s : token total supply is limited int256 and must be positive", amt)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenSupply, msg)
}
func ErrInvalidSendLockAmt(amt string) sdk.Error {
	msg := fmt.Sprintf("invalid send lock %s : token send lock amount must be positive", amt)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidSendLockAmt, msg)
}
func ErrInvalidTokenOwner(addr sdk.Address) sdk.Error {
	msg := fmt.Sprintf("invalid owner %s : token owner is invalid", addr.String())
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenOwner, msg)
}
func ErrInvalidTokenMintAmt(amt string) sdk.Error {
	msg := fmt.Sprintf("invalid mint amount %s : token total supply is limited int256 and must be positive", amt)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenMintAmt, msg)
}
func ErrInvalidTokenBurnAmt(amt string) sdk.Error {
	msg := fmt.Sprintf("invalid burn amount %s : token total supply is limited int256 and must be positive", amt)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenBurnAmt, msg)
}
func ErrInvalidTokenForbidden(symbol string) sdk.Error {
	msg := fmt.Sprintf("invalid forbid %s : token has been forbidden", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenForbidden, msg)
}
func ErrInvalidTokenUnForbidden(symbol string) sdk.Error {
	msg := fmt.Sprintf("invalid unforbid %s : token has not been forbidden", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenUnForbidden, msg)
}
func ErrInvalidTokenWhitelist() sdk.Error {
	msg := fmt.Sprintf("whitelist : token whitelist is invalid")
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenWhitelist, msg)
}
func ErrInvalidForbiddenAddress() sdk.Error {
	msg := fmt.Sprintf("forbidden address : address is invalid")
	return sdk.NewError(CodeSpaceAsset, CodeInvalidForbiddenAddress, msg)
}
func ErrInvalidTokenURL(url string) sdk.Error {
	msg := fmt.Sprintf("invalid url %s : token url is limited to %d unicode characters", url, MaxTokenURLLength)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenURL, msg)
}
func ErrInvalidTokenDescription(description string) sdk.Error {
	msg := fmt.Sprintf("invalid description %s : token description is limited to %d bytes size", description, MaxTokenDescriptionLength)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenDescription, msg)
}
func ErrInvalidTokenIdentity(identity string) sdk.Error {
	msg := fmt.Sprintf("invalid identity %s : token identity is limited to %d", identity, MaxTokenIdentityLength)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenIdentity, msg)
}
func ErrNilTokenIdentity() sdk.Error {
	msg := fmt.Sprintf("token identity is nil")
	return sdk.NewError(CodeSpaceAsset, CodeNilTokenIdentity, msg)
}

// -----------------------------------------------------------------------------
func ErrTokenNotFound(symbol string) sdk.Error {
	msg := fmt.Sprintf("token %s is not in store", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeTokenNotFound, msg)
}

func ErrDuplicateTokenSymbol(symbol string) sdk.Error {
	msg := fmt.Sprintf("token symbol %s already exists in store", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeDuplicateTokenSymbol, msg)
}

func ErrTransferSelfTokenOwner() sdk.Error {
	msg := fmt.Sprintf("can not and no need to transfer ownership to self")
	return sdk.NewError(CodeSpaceAsset, CodeTransferSelfTokenOwner, msg)
}
func ErrNilTokenOwner() sdk.Error {
	msg := fmt.Sprintf("token owner is nil")
	return sdk.NewError(CodeSpaceAsset, CodeNilTokenOwner, msg)
}
func ErrNeedTokenOwner(addr sdk.Address) sdk.Error {
	msg := fmt.Sprintf("only token owner %s can operate this", addr.String())
	return sdk.NewError(CodeSpaceAsset, CodeNeedTokenOwner, msg)
}
func ErrInvalidIssueOwner() sdk.Error {
	msg := fmt.Sprintf("only coinex dex foundation can issue reserved symbol token, you can run \n" +
		"$ cetcli query asset reserved-symbol \n" +
		"to query reserved token symbol")
	return sdk.NewError(CodeSpaceAsset, CodeInvalidIssueOwner, msg)
}

func ErrTokenMintNotSupported(symbol string) sdk.Error {
	msg := fmt.Sprintf("token %s do not support mint", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeTokenMintNotSupported, msg)
}

func ErrTokenBurnNotSupported(symbol string) sdk.Error {
	msg := fmt.Sprintf("token %s do not support burn", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeTokenBurnNotSupported, msg)
}

func ErrTokenForbiddenNotSupported(symbol string) sdk.Error {
	msg := fmt.Sprintf("token %s do not support token forbid", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeTokenForbiddenNotSupported, msg)
}

func ErrAddressForbiddenNotSupported(symbol string) sdk.Error {
	msg := fmt.Sprintf("token %s do not support address forbid", symbol)
	return sdk.NewError(CodeSpaceAsset, CodeAddressForbiddenNotSupported, msg)
}

func ErrNilTokenWhitelist() sdk.Error {
	msg := fmt.Sprintf("whitelist is nil")
	return sdk.NewError(CodeSpaceAsset, CodeNilTokenWhitelist, msg)
}

func ErrNilForbiddenAddress() sdk.Error {
	msg := fmt.Sprintf("forbidden address is nil")
	return sdk.NewError(CodeSpaceAsset, CodeNilForbiddenAddress, msg)
}

func ErrAccInBlackList(addr sdk.AccAddress) sdk.Error {
	msg := fmt.Sprintf("%s is in blacklist", addr.String())
	return sdk.NewError(CodeSpaceAsset, CodeAddrInBlackList, msg)
}

func ErrTokenOwnerSelfForbidden() sdk.Error {
	msg := fmt.Sprintf("Token owner can not be forbidden")
	return sdk.NewError(CodeSpaceAsset, CodeTokenOwnerSelfForbidden, msg)
}

func ErrInvalidTokenInfo(field, val string) sdk.Error {
	msg := fmt.Sprintf("invalid token %s: %s", field, val)
	return sdk.NewError(CodeSpaceAsset, CodeInvalidTokenInfo, msg)
}
func ErrCodeTokenInfoSealed(field string) sdk.Error {
	msg := fmt.Sprintf("token %s sealed", field)
	return sdk.NewError(CodeSpaceAsset, CodeTokenInfoSealed, msg)
}
