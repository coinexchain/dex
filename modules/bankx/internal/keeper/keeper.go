package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

type Keeper struct {
	ParamSubspace params.Subspace
	Axk           types.ExpectedAccountXKeeper
	Bk            bank.Keeper
	Ak            auth.AccountKeeper
	Tk            types.ExpectedAssetStatusKeeper
	Sk            types.SupplyKeeper
	MsgProducer   msgqueue.MsgSender
}

func NewKeeper(paramSubspace params.Subspace, axk authx.AccountXKeeper,
	bk bank.BaseKeeper, ak auth.AccountKeeper,
	tk types.ExpectedAssetStatusKeeper, sk types.SupplyKeeper, msgProducer msgqueue.MsgSender) Keeper {

	return Keeper{
		ParamSubspace: paramSubspace.WithKeyTable(types.ParamKeyTable()),
		Axk:           axk,
		Bk:            bk,
		Ak:            ak,
		Tk:            tk,
		Sk:            sk,
		MsgProducer:   msgProducer,
	}
}

func (k Keeper) GetParams(ctx sdk.Context) (param types.Params) {
	k.ParamSubspace.GetParamSet(ctx, &param)
	return
}
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.ParamSubspace.SetParamSet(ctx, &params)
}

func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return k.Bk.HasCoins(ctx, addr, amt)
}

func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	ret := k.Bk.SendCoins(ctx, from, to, amt)
	return ret
}

func (k Keeper) FreezeCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, err := k.Bk.SubtractCoins(ctx, addr, amt)
	if err != nil {
		return err
	}
	accx := k.Axk.GetOrCreateAccountX(ctx, addr)
	frozenCoins := accx.FrozenCoins.Add(amt)
	accx.FrozenCoins = frozenCoins
	k.Axk.SetAccountX(ctx, accx)

	return nil
}

func (k Keeper) UnFreezeCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	accx, ok := k.Axk.GetAccountX(ctx, addr)
	if !ok {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	frozenCoins, neg := accx.FrozenCoins.SafeSub(amt)
	if neg {
		return sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze")
	}

	accx.FrozenCoins = frozenCoins
	k.Axk.SetAccountX(ctx, accx)

	_, err := k.Bk.AddCoins(ctx, addr, amt)
	return err
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, err := k.Bk.SubtractCoins(ctx, addr, amt)
	return err
}
func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if _, err := k.Bk.AddCoins(ctx, addr, amt); err != nil {
		return err
	}
	return nil
}

func (k Keeper) DeductInt64CetFee(ctx sdk.Context, addr sdk.AccAddress, amt int64) sdk.Error {
	return k.DeductFee(ctx, addr, dex.NewCetCoins(amt))
}

func (k Keeper) DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	err := k.Sk.SendCoinsFromAccountToModule(ctx, addr, auth.FeeCollectorName, amt)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) IsSendForbidden(ctx sdk.Context, amt sdk.Coins, addr sdk.AccAddress) bool {
	for _, coin := range amt {
		if k.Tk.IsForbiddenByTokenIssuer(ctx, coin.Denom, addr) {
			return true
		}
	}
	return false
}

func (k Keeper) GetTotalCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc := k.Ak.GetAccount(ctx, addr)
	accx, found := k.Axk.GetAccountX(ctx, addr)
	var coins = sdk.Coins{}
	if acc != nil {
		coins = acc.GetCoins()
	}
	if found {
		coins = coins.Add(accx.GetAllCoins())
	}
	return coins

}

func (k Keeper) TotalAmountOfCoin(ctx sdk.Context, denom string) sdk.Int {
	var (
		axkTotalAmount = sdk.ZeroInt()
		akTotalAmount  = sdk.ZeroInt()
	)
	axkProcess := func(acc authx.AccountX) bool {
		val := acc.GetAllCoins().AmountOf(denom)
		axkTotalAmount = axkTotalAmount.Add(val)
		return false
	}

	akProcess := func(acc auth.Account) bool {
		val := acc.GetCoins().AmountOf(denom)
		akTotalAmount = akTotalAmount.Add(val)
		return false
	}

	k.Axk.IterateAccounts(ctx, axkProcess)
	k.Ak.IterateAccounts(ctx, akProcess)

	return axkTotalAmount.Add(akTotalAmount)
}

func (k Keeper) BlacklistedAddr(addr sdk.AccAddress) bool {
	return k.Bk.BlacklistedAddr(addr)
}
