package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// SupplyKeeper defines the expected supply keeper (noalias)
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	//SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}

type ExpectedAccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
	SetAccount(ctx sdk.Context, acc exported.Account)
}