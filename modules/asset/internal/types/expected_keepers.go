package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error

	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	GetTotalCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// Supply Keeper will implement the interface
type ExpectedSupplyKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
}
