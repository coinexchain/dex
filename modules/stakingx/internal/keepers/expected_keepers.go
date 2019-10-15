package keepers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	supply "github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/coinexchain/dex/modules/asset"
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
	GetModuleAccount(ctx sdk.Context, name string) supply.ModuleAccountI
	GetSupply(ctx sdk.Context) (supply supply.SupplyI)
}

type AssetViewKeeper interface {
	GetToken(ctx sdk.Context, symbol string) asset.Token
}
