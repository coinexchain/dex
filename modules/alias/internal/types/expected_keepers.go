package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	DeductInt64CetFee(ctx sdk.Context, addr sdk.AccAddress, amt int64) sdk.Error
}

// Asset Keeper will implement the interface
type ExpectedAssetStatusKeeper interface {
	IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
}
