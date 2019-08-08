package stakingx

import (
	"github.com/coinexchain/dex/modules/asset"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

type DistributionKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool types.FeePool)
	SetFeePool(ctx sdk.Context, feePool types.FeePool)
	GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins

	IterateValidatorOutstandingRewards(ctx sdk.Context, handler func(val sdk.ValAddress, rewards types.ValidatorOutstandingRewards) (stop bool))
}

type ExpectBankxKeeper interface {
	TotalAmountOfCoin(ctx sdk.Context, denom string) sdk.Int
}

type ExpectSupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
}

type AssetViewKeeper interface {
	GetToken(ctx sdk.Context, symbol string) asset.Token
}
