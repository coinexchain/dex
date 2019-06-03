package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Category   = "category"
	TxCategory = "asset"

	Issue           = "issue"
	Token           = "token"
	Owner           = "owner"
	OriginalOwner   = "original-owner"
	NewOwner        = "new-owner"
	Amt             = "amount"
	AddWhitelist    = "add-white-list"
	RemoveWhitelist = "remove-white-list"
)

// Tag keys and values
var (
	Action = sdk.TagAction
)
