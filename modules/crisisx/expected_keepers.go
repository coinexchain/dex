package crisisx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExpectBankxKeeper interface {
	TotalAmountOfCoin(ctx sdk.Context, denom string) sdk.Int
}
