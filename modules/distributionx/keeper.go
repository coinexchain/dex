package distributionx

import (
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/cosmos/cosmos-sdk/x/distribution"
)

type Keeper struct {
	bxk bankx.Keeper
	dk  distribution.Keeper
}

func NewKeeper(bxk bankx.Keeper, dk distribution.Keeper) Keeper {
	return Keeper{
		bxk,
		dk,
	}
}
