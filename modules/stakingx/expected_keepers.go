package stakingx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

type DistributionKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool types.FeePool)
	SetFeePool(ctx sdk.Context, feePool types.FeePool)
	GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins
}

type ExpectBankxKeeper interface {
	TotalAmountOfCoin(ctx sdk.Context, denom string) sdk.Int
}

type ExpectSupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
}
