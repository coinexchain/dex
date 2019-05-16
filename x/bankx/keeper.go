package bankx

import (
	"github.com/coinexchain/dex/x/authx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper struct {
	paramSubspace params.Subspace
	axk           authx.AccountXKeeper
	bk            bank.BaseKeeper
	fck           auth.FeeCollectionKeeper
}

func NewKeeper(paramSubspace params.Subspace, axk authx.AccountXKeeper, bk bank.BaseKeeper, fck auth.FeeCollectionKeeper) Keeper {
	return Keeper{
		paramSubspace: paramSubspace.WithKeyTable(ParamKeyTable()),
		axk:           axk,
		bk:            bk,
		fck:           fck,
	}
}

func (k Keeper) GetParam(ctx sdk.Context) (param Param) {
	k.paramSubspace.Get(ctx, ParamStoreKeyActivatedFee, &param)
	return
}
func (k Keeper) SetParam(ctx sdk.Context, param Param) {
	k.paramSubspace.Set(ctx, ParamStoreKeyActivatedFee, param)
}
