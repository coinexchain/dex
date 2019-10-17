package keeper

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

type Keeper struct {
	paramSubspace params.Subspace
	axk           types.ExpectedAccountXKeeper
	bk            bank.Keeper
	ak            auth.AccountKeeper
	tk            types.ExpectedAssetStatusKeeper
	sk            types.SupplyKeeper
	MsgProducer   msgqueue.MsgSender
}

func NewKeeper(paramSubspace params.Subspace, axk authx.AccountXKeeper,
	bk bank.BaseKeeper, ak auth.AccountKeeper,
	tk types.ExpectedAssetStatusKeeper, sk types.SupplyKeeper, msgProducer msgqueue.MsgSender) Keeper {

	return Keeper{
		paramSubspace: paramSubspace.WithKeyTable(types.ParamKeyTable()),
		axk:           axk,
		bk:            bk,
		ak:            ak,
		tk:            tk,
		sk:            sk,
		MsgProducer:   msgProducer,
	}
}

func (k Keeper) GetParams(ctx sdk.Context) (param types.Params) {
	k.paramSubspace.GetParamSet(ctx, &param)
	return
}
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return k.bk.HasCoins(ctx, addr, amt)
}

func (k Keeper) SendCoins(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if k.IsSendForbidden(ctx, amt, from) {
		return types.ErrTokenForbiddenByOwner()
	}
	ret := k.bk.SendCoins(ctx, from, to, amt)
	return ret
}

func (k Keeper) SendLockedCoins(ctx sdk.Context, fromAddr, toAddr, supervisor sdk.AccAddress, amt sdk.Coins,
	unlockTime int64, reward int64) sdk.Error {
	if k.IsSendForbidden(ctx, amt, fromAddr) {
		return types.ErrTokenForbiddenByOwner()
	}
	if k.ak.GetAccount(ctx, toAddr) == nil {
		if err := k.AddCoins(ctx, toAddr, sdk.Coins{}); err != nil {
			return err
		}
	}

	lockDuration := (unlockTime - ctx.BlockHeader().Time.Unix()) * int64(time.Second)
	if lockDuration > k.GetParams(ctx).LockCoinsFreeTime {
		exceededDays := (lockDuration-k.GetParams(ctx).LockCoinsFreeTime-1)/(24*int64(time.Hour)) + 1
		lockCoinsFee := dex.NewCetCoins(k.GetParams(ctx).LockCoinsFeePerDay * exceededDays)
		if err := k.DeductFee(ctx, fromAddr, lockCoinsFee); err != nil {
			return err
		}
	}

	if err := k.SubtractCoins(ctx, fromAddr, amt); err != nil {
		return err
	}

	ax := k.axk.GetOrCreateAccountX(ctx, toAddr)
	for _, coin := range amt {
		ax.LockedCoins = append(ax.LockedCoins, authx.LockedCoin{
			FromAddress: fromAddr,
			Supervisor:  supervisor,
			Coin:        coin,
			UnlockTime:  unlockTime,
			Reward:      reward,
		})
		if err := k.tk.UpdateTokenSendLock(ctx, coin.Denom, coin.Amount, true); err != nil {
			return err
		}
	}
	k.axk.SetAccountX(ctx, ax)

	if !amt.Empty() {
		k.axk.InsertUnlockedCoinsQueue(ctx, unlockTime, toAddr)
	}
	return nil
}

func (k Keeper) ReturnLockedCoins(ctx sdk.Context, fromAddr, toAddr, supervisor sdk.AccAddress, amt *sdk.Coin,
	unlockTime int64, reward int64) sdk.Error {

	lockedCoin := authx.NewLockedCoin(fromAddr, supervisor, amt.Denom, amt.Amount, unlockTime, reward)
	if err := k.earlierUnlockCoin(ctx, toAddr, &lockedCoin); err != nil {
		return err
	}
	returnAmt := sdk.NewCoin(amt.Denom, amt.Amount.Sub(sdk.NewInt(reward)))
	if err := k.AddCoins(ctx, fromAddr, sdk.NewCoins(returnAmt)); err != nil {
		return err
	}
	rewardAmt := sdk.NewCoin(amt.Denom, sdk.NewInt(reward))
	if err := k.AddCoins(ctx, supervisor, sdk.NewCoins(rewardAmt)); err != nil {
		return err
	}

	return nil
}

func (k Keeper) EarlierUnlockBySender(ctx sdk.Context, fromAddr, toAddr, supervisor sdk.AccAddress, amt *sdk.Coin,
	unlockTime int64, reward int64) sdk.Error {

	lockedCoin := authx.NewLockedCoin(fromAddr, supervisor, amt.Denom, amt.Amount, unlockTime, reward)
	if err := k.earlierUnlockCoin(ctx, toAddr, &lockedCoin); err != nil {
		return err
	}
	if supervisor.Empty() {
		if err := k.AddCoins(ctx, toAddr, sdk.NewCoins(*amt)); err != nil {
			return err
		}
	} else {
		receivedAmt := sdk.NewCoin(amt.Denom, amt.Amount.Sub(sdk.NewInt(reward)))
		if err := k.AddCoins(ctx, toAddr, sdk.NewCoins(receivedAmt)); err != nil {
			return err
		}
		rewardAmt := sdk.NewCoin(amt.Denom, sdk.NewInt(reward))
		if err := k.AddCoins(ctx, supervisor, sdk.NewCoins(rewardAmt)); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) EarlierUnlockBySupervisor(ctx sdk.Context, fromAddr, toAddr, supervisor sdk.AccAddress, amt *sdk.Coin,
	unlockTime int64, reward int64) sdk.Error {

	lockedCoin := authx.NewLockedCoin(fromAddr, supervisor, amt.Denom, amt.Amount, unlockTime, reward)
	if err := k.earlierUnlockCoin(ctx, toAddr, &lockedCoin); err != nil {
		return err
	}
	receivedAmt := sdk.NewCoin(amt.Denom, amt.Amount.Sub(sdk.NewInt(reward)))
	if err := k.AddCoins(ctx, toAddr, sdk.NewCoins(receivedAmt)); err != nil {
		return err
	}
	rewardAmt := sdk.NewCoin(amt.Denom, sdk.NewInt(reward))
	if err := k.AddCoins(ctx, supervisor, sdk.NewCoins(rewardAmt)); err != nil {
		return err
	}

	return nil
}

func (k Keeper) earlierUnlockCoin(ctx sdk.Context, addr sdk.AccAddress, amt *authx.LockedCoin) sdk.Error {
	ax, ok := k.axk.GetAccountX(ctx, addr)
	if !ok {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	coinIndex := -1
	hasOther := false
	for i, lockedCoin := range ax.LockedCoins {
		if bytes.Equal(amt.FromAddress, lockedCoin.FromAddress) &&
			bytes.Equal(amt.Supervisor, lockedCoin.Supervisor) &&
			amt.Coin.IsEqual(lockedCoin.Coin) &&
			amt.UnlockTime == lockedCoin.UnlockTime &&
			amt.Reward == lockedCoin.Reward {
			coinIndex = i
		} else if amt.UnlockTime == lockedCoin.UnlockTime {
			hasOther = true
		}
	}

	if coinIndex < 0 {
		return types.ErrorLockedCoinNotFound()
	}

	if err := k.tk.UpdateTokenSendLock(ctx, amt.Coin.Denom, amt.Coin.Amount, false); err != nil {
		return err
	}

	ax.LockedCoins = append(ax.LockedCoins[:coinIndex], ax.LockedCoins[coinIndex+1:]...)
	k.axk.SetAccountX(ctx, ax)

	if !hasOther {
		k.axk.RemoveFromUnlockedCoinsQueue(ctx, amt.UnlockTime, addr)
	}

	return nil
}

func (k Keeper) FreezeCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if k.IsSendForbidden(ctx, amt, addr) {
		return types.ErrTokenForbiddenByOwner()
	}
	_, err := k.bk.SubtractCoins(ctx, addr, amt)
	if err != nil {
		return err
	}
	accx := k.axk.GetOrCreateAccountX(ctx, addr)
	accx.FrozenCoins = accx.FrozenCoins.Add(amt)
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

	_, err := k.bk.AddCoins(ctx, addr, amt)
	return err
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	_, err := k.bk.SubtractCoins(ctx, addr, amt)
	return err
}

func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	if _, err := k.bk.AddCoins(ctx, addr, amt); err != nil {
		return err
	}
	return nil
}

func (k Keeper) MockAddLockedCoins(ctx sdk.Context, addr sdk.AccAddress, lockedCoins authx.LockedCoins) {
	ax := k.axk.GetOrCreateAccountX(ctx, addr)
	ax.LockedCoins = append(ax.LockedCoins, lockedCoins...)
	k.axk.SetAccountX(ctx, ax)
}

func (k Keeper) MockAddFrozenCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) {
	ax := k.axk.GetOrCreateAccountX(ctx, addr)
	ax.FrozenCoins = ax.FrozenCoins.Add(amt)
	k.axk.SetAccountX(ctx, ax)
}

func (k Keeper) DeductInt64CetFee(ctx sdk.Context, addr sdk.AccAddress, amt int64) sdk.Error {
	return k.DeductFee(ctx, addr, dex.NewCetCoins(amt))
}

func (k Keeper) DeductActivationFee(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, transfer sdk.Coins) (sdk.Coins, sdk.Error) {
	//toAccount doesn't exist yet
	if k.ak.GetAccount(ctx, to) == nil {
		activationFee := dex.NewCetCoins(k.GetParams(ctx).ActivationFee)
		amt, neg := transfer.SafeSub(activationFee)
		if neg {
			return transfer, types.ErrorInsufficientCETForActivatingFee()
		}
		return amt, k.DeductFee(ctx, from, activationFee)
	}
	return transfer, nil
}
func (k Keeper) PreCheckFreshAccounts(ctx sdk.Context, outputs []bank.Output) (addrs []sdk.AccAddress) {
	addrsMap := make(map[string]bool)
	for _, output := range outputs {
		//toAccount doesn't exist yet
		if k.ak.GetAccount(ctx, output.Address) == nil && !addrsMap[output.Address.String()] {
			addrs = append(addrs, output.Address)
			addrsMap[output.Address.String()] = true
		}
	}
	return addrs
}
func (k Keeper) DeductActivationFeeForFreshAccounts(ctx sdk.Context, addrs []sdk.AccAddress) sdk.Error {
	fee := dex.NewCetCoins(k.GetParams(ctx).ActivationFee)
	for _, addr := range addrs {
		if err := k.DeductFee(ctx, addr, fee); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) DeductFee(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.sk.SendCoinsFromAccountToModule(ctx, addr, auth.FeeCollectorName, amt)
}

func (k Keeper) DonateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.sk.SendCoinsFromAccountToModule(ctx, addr, distribution.ModuleName, amt)
}

func (k Keeper) IsSendForbidden(ctx sdk.Context, amt sdk.Coins, addr sdk.AccAddress) bool {
	for _, coin := range amt {
		if k.tk.IsForbiddenByTokenIssuer(ctx, coin.Denom, addr) {
			return true
		}
	}
	return false
}

func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.bk.GetCoins(ctx, addr)
}

func (k Keeper) GetFrozenCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	accX, _ := k.axk.GetAccountX(ctx, addr)
	return accX.FrozenCoins
}

func (k Keeper) GetLockedCoins(ctx sdk.Context, addr sdk.AccAddress) authx.LockedCoins {
	accX, _ := k.axk.GetAccountX(ctx, addr)
	return accX.LockedCoins
}

func (k Keeper) GetTotalCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc := k.ak.GetAccount(ctx, addr)
	accx, found := k.axk.GetAccountX(ctx, addr)
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

	k.axk.IterateAccounts(ctx, axkProcess)
	k.ak.IterateAccounts(ctx, akProcess)

	return axkTotalAmount.Add(akTotalAmount)
}

func (k Keeper) BlacklistedAddr(addr sdk.AccAddress) bool {
	return k.bk.BlacklistedAddr(addr)
}

func (k Keeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	k.bk.SetSendEnabled(ctx, enabled)
}

func (k Keeper) GetSendEnabled(ctx sdk.Context) bool {
	return k.bk.GetSendEnabled(ctx)
}

func (k Keeper) InputOutputCoins(ctx sdk.Context, inputs []bank.Input, outputs []bank.Output) sdk.Error {
	return k.bk.InputOutputCoins(ctx, inputs, outputs)
}

func (k Keeper) SetMemoRequired(ctx sdk.Context, addr sdk.AccAddress, required bool) sdk.Error {
	account := k.ak.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	accountX := k.axk.GetOrCreateAccountX(ctx, addr)
	accountX.MemoRequired = required
	k.axk.SetAccountX(ctx, accountX)

	return nil
}

func (k Keeper) GetMemoRequired(ctx sdk.Context, addr sdk.AccAddress) bool {
	if accX, ok := k.axk.GetAccountX(ctx, addr); ok {
		return accX.MemoRequired
	}
	return false
}
