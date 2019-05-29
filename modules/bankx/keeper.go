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
func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return k.bk.HasCoins(ctx, addr, amt)
}
func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, ret := k.bk.SendCoins(ctx, from, to, amt)
	return ret
}
func (k Keeper) FreezeCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {

	acc := k.ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.ErrInvalidAddress("account doesn't exist yet")
	}

	newCoins, neg := acc.GetCoins().SafeSub(amt)
	if neg {
		return sdk.ErrInsufficientCoins("account has insufficient coins to freeze")
	}
	acc.SetCoins(newCoins)
	k.ak.SetAccount(ctx, acc)

	accx, _ := k.axk.GetAccountX(ctx, addr)
	frozenCoins := accx.FrozenCoins.Add(amt)
	accx.FrozenCoins = frozenCoins
	k.axk.SetAccountX(ctx, accx)

	return nil
}

func (k Keeper) UnFreezeCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {

	accx, ok := k.axk.GetAccountX(ctx, addr)
	if !ok {
		return sdk.ErrInvalidAddress("account doesn't exist yet")
	}
	frozenCoins, neg := accx.FrozenCoins.SafeSub(amt)
	if neg {
		return sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze")
	}
	accx.FrozenCoins = frozenCoins
	k.axk.SetAccountX(ctx, accx)

	acc := k.ak.GetAccount(ctx, addr)
	newcoins := acc.GetCoins().Add(amt)
	acc.SetCoins(newcoins)
	k.ak.SetAccount(ctx, acc)

	return nil
}
