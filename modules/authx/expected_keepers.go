package authx

import sdk "github.com/cosmos/cosmos-sdk/types"

// SupplyKeeper defines the expected supply keeper (noalias)
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	//SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}
