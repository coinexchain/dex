package bankx

import sdk "github.com/cosmos/cosmos-sdk/types"

const (

	CodeSpaceBankx ="bankx"

	CodeFirstTransferNotCET=19

)

func ErrorFirstTransferNotCET(codespace sdk.CodespaceType) sdk.Error{
	return sdk.NewError(codespace,CodeFirstTransferNotCET,"first transfer must be CET")
}