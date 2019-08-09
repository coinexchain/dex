package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error // to tranfer coins
	FreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                   // freeze some coins when creating orders
	UnFreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                 // unfreeze coins and then orders can be executed
	DeductFee(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error
	DeductInt64CetFee(ctx sdk.Context, addr sdk.AccAddress, amt int64) sdk.Error
}

// Asset Keeper will implement the interface
type ExpectedAssetStatusKeeper interface {
	IsTokenExists(ctx sdk.Context, denom string) bool // check whether there is a coin named "denom"
	IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
}

// market keeper will implement the interface
type ExpectedMarketKeeper interface {
	GetMarketLastExePrice(ctx sdk.Context, symbol string) (sdk.Dec, error)
	IsMarketExist(ctx sdk.Context, symbol string) bool
}
