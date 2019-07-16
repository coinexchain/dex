package market

import (
	"github.com/coinexchain/dex/modules/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool                          // to check whether have sufficient coins in special address
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error // to tranfer coins
	FreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                   // freeze some coins when creating orders
	UnFreezeCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) sdk.Error                 // unfreeze coins and then orders can be executed
}

// Asset Keeper will implement the interface
type ExpectedAssetStatusKeeper interface {
	IsTokenForbidden(ctx sdk.Context, denom string) bool // the coin's issuer has forbidden "denom", forbiding transmission and exchange.
	IsTokenExists(ctx sdk.Context, denom string) bool    // check whether there is a coin named "denom"
	IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
	IsForbiddenByTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool
	GetToken(ctx sdk.Context, symbol string) types.Token
}

type ExpectFeeKeeper interface {
	AddCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins
}

type ExpectedChargeFeeKeeper interface {
	SubtractFeeAndCollectFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
}
