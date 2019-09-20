package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error // to tranfer coins
}

// Asset Keeper will implement the interface
type ExpectedAssetStatusKeeper interface {
	IsTokenExists(ctx sdk.Context, denom string) bool // check whether there is a coin named "denom"
}

type ExpectedDistributionxKeeper interface {
	DonateToCommunityPool(ctx sdk.Context, fromAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}

type ExpectedAccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
}
