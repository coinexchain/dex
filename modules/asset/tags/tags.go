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
	AddWhitelist    = "add-whitelist"
	RemoveWhitelist = "remove-whitelist"
	Addresses       = "addresses"
	URL             = "url"
)

// Tag keys and values
var (
	Action = sdk.TagAction
)
