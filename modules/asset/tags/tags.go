package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Category   = "category"
	TxCategory = "asset"

	Issue = "issue"
	Token = "token"
	Owner = "owner"

	TransferOwnership = "TransferOwnership"
	OriginalOwner     = "OriginalOwner"
	NewOwner          = "NewOwner"

	MintToken = "MintToken"
	Amt       = "Amount"
)

// Tag keys and values
var (
	Action = sdk.TagAction
)
