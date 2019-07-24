package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
}

// Asset Keeper will implement the interface
type ExpectedAssetStatusKeeper interface {
	IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
}
