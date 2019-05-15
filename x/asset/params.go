package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params defines the parameters for the asset module.
type Params struct {
	IssueTokenFee                 sdk.Coins `json:"issue_token_fee"`
	FreezeAddrFee                 sdk.Coins `json:"freeze_address_fee"`
	UnFreezeAddrFee               sdk.Coins `json:"unfreeze_address_fee"`
	FreezeTokenFee                sdk.Coins `json:"freeze_token_fee"`
	UnFreezeTokenFee              sdk.Coins `json:"unfreeze_token_fee"`
	TokenFreezeWhitelistAddFee    sdk.Coins `json:"token_freeze_whitelist_add_fee"`
	TokenFreezeWhitelistRemoveFee sdk.Coins `json:"token_freeze_whitelist_remove_fee"`
	BurnFee                       sdk.Coins `json:"burn_fee"`
	MintFee                       sdk.Coins `json:"mint_fee"`
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		IssueTokenFee:                 cet(10000),
		FreezeAddrFee:                 cet(10),
		UnFreezeAddrFee:               cet(10),
		FreezeTokenFee:                cet(1000),
		UnFreezeTokenFee:              cet(1000),
		TokenFreezeWhitelistAddFee:    cet(100),
		TokenFreezeWhitelistRemoveFee: cet(100),
		BurnFee:                       cet(1000),
		MintFee:                       cet(1000),
	}
}

func cet(amt int64) sdk.Coins {
	return sdk.Coins{
		sdk.NewCoin("cet", sdk.NewInt(amt)),
	}
}
