package bankx

import (
	"github.com/coinexchain/dex/modules/authx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper struct {
	paramSubspace params.Subspace
	axk           authx.AccountXKeeper
	bk            bank.BaseKeeper
	ak            auth.AccountKeeper
	fck           auth.FeeCollectionKeeper
}

func NewKeeper(paramSubspace params.Subspace, axk authx.AccountXKeeper, bk bank.BaseKeeper, ak auth.AccountKeeper, fck auth.FeeCollectionKeeper) Keeper {
	return Keeper{
		paramSubspace: paramSubspace.WithKeyTable(ParamKeyTable()),
		axk:           axk,
		bk:            bk,
		ak:            ak,
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
