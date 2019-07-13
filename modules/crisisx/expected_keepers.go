package crisisx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

type ExpectBankxKeeper interface {
	TotalAmountOfCoin(ctx sdk.Context, denom string) sdk.Int
}

type ExpectSupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
}
