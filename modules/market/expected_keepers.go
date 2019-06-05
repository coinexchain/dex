package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool                          // to check whether have sufficient coins in special address
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error // to tranfer coins
	FreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                   // freeze some coins when creating orders
	UnFreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                 // unfreeze coins and then orders can be executed
}

// Asset Keeper will implement the interface
type ExpectedAssertStatusKeeper interface {
	IsTokenForbidden(ctx sdk.Context, denom string) bool // the coin's issuer has forbidden "denom", forbiding transmission and exchange.
	IsTokenExists(ctx sdk.Context, denom string) bool    // check whether there is a coin named "denom"
	IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
	IsForbiddenByTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
}
