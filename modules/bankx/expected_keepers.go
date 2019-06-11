package bankx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExpectedAssetStatusKeeper interface {
	IsTokenForbidden(ctx sdk.Context, symbol string) bool
	IsForbiddenByTokenIssuer(ctx sdk.Context, symbol string, addr sdk.AccAddress) bool
}
