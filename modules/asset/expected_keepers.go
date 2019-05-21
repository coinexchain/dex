package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper - expected bank keeper,follow sdk principle of least authority.
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	SetSendEnabled(ctx sdk.Context, enabled bool)
}
