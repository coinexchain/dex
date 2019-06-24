package bankx

import (
	"fmt"
	"github.com/coinexchain/dex/modules/msgqueue"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/authx"
)

type Keeper struct {
	paramSubspace params.Subspace
	axk           authx.AccountXKeeper
	bk            bank.BaseKeeper
	ak            auth.AccountKeeper
	fck           auth.FeeCollectionKeeper
	tk            ExpectedAssetStatusKeeper
	msgProducer   msgqueue.Producer
}

func NewKeeper(paramSubspace params.Subspace, axk authx.AccountXKeeper,
	bk bank.BaseKeeper, ak auth.AccountKeeper, fck auth.FeeCollectionKeeper,
	tk ExpectedAssetStatusKeeper, msgProducer msgqueue.Producer) Keeper {

	return Keeper{
		paramSubspace: paramSubspace.WithKeyTable(ParamKeyTable()),
		axk:           axk,
		bk:            bk,
		ak:            ak,
		fck:           fck,
		tk:            tk,
		msgProducer:   msgProducer,
	}
}

func (k Keeper) GetParam(ctx sdk.Context) (param Params) {
	k.paramSubspace.GetParamSet(ctx, &param)
	return
}
func (k Keeper) SetParam(ctx sdk.Context, params Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return k.bk.HasCoins(ctx, addr, amt)
}

func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, ret := k.bk.SendCoins(ctx, from, to, amt)
	return ret
}

func (k Keeper) FreezeCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, _, err := k.bk.SubtractCoins(ctx, addr, amt)
	if err != nil {
		return err
	}

	accx := k.axk.GetOrCreateAccountX(ctx, addr)
	frozenCoins := accx.FrozenCoins.Add(amt)
	accx.FrozenCoins = frozenCoins
	k.axk.SetAccountX(ctx, accx)

	return nil
}

func (k Keeper) UnFreezeCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	accx, ok := k.axk.GetAccountX(ctx, addr)
	if !ok {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	frozenCoins, neg := accx.FrozenCoins.SafeSub(amt)
	if neg {
		return sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze")
	}

	accx.FrozenCoins = frozenCoins
	k.axk.SetAccountX(ctx, accx)

	_, _, err := k.bk.AddCoins(ctx, addr, amt)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, _, err := k.bk.SubtractCoins(ctx, addr, amt)
	return err
}
func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if _, _, err := k.bk.AddCoins(ctx, addr, amt); err != nil {
		return err
	}
	return nil
}

func (k Keeper) DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if _, _, err := k.bk.SubtractCoins(ctx, addr, amt); err != nil {
		return err
	}

	k.fck.AddCollectedFees(ctx, amt)
	return nil
}

func (k Keeper) IsSendForbidden(ctx sdk.Context, amt sdk.Coins, addr sdk.AccAddress) bool {
	for _, coin := range amt {
		if k.tk.IsForbiddenByTokenIssuer(ctx, coin.Denom, addr) {
			return true
		}
	}
	return false
}
