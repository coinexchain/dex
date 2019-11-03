package bankx

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)
		case types.MsgSetMemoRequired:
			return handleMsgSetMemoRequired(ctx, k, msg)
		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)
		default:
			return dex.ErrUnknownRequest(ModuleName, msg)
		}
	}
}

func handleMsgMultiSend(ctx sdk.Context, k Keeper, msg types.MsgMultiSend) sdk.Result {
	if err := k.GetSendEnabled(ctx); err != nil {
		return err.Result()
	}

	for _, out := range msg.Outputs {
		if k.BlacklistedAddr(out.Address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", out.Address)).Result()
		}
	}

	for _, input := range msg.Inputs {
		if k.IsSendForbidden(ctx, input.Coins, input.Address) {
			return types.ErrTokenForbiddenByOwner().Result()
		}
	}

	addrs := k.PreCheckFreshAccounts(ctx, msg.Outputs)

	if err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs); err != nil {
		return err.Result()
	}

	if err := k.DeductActivationFeeForFreshAccounts(ctx, addrs); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	//fillMsgQueue(ctx, k, "multi_send", msg)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}

}

var SendEnabled bool = true

func handleMsgSend(ctx sdk.Context, k Keeper, msg types.MsgSend) sdk.Result {
	//if err := k.GetSendEnabled(ctx); err != nil {
	//	return err.Result()
	//}
	if !SendEnabled {
		return bank.ErrSendDisabled(types.CodeSpaceBankx).Result()
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.ToAddress)).Result()
	}
	amt := msg.Amount
	if !amt.IsValid() {
		//tmp for pprof test
		return sdk.Result{Code: 2}
	}
	if k.IsSendForbidden(ctx, amt, msg.FromAddress) {
		return types.ErrTokenForbiddenByOwner().Result()
	}

	ak, _ := k.GetAccount()
	acc := ak.GetAccount(ctx, msg.FromAddress)
	if acc == nil {
		//tmp for pprof test
		return sdk.Result{Code: 3}
	}
	srcCoins := acc.GetCoins()
	srcSpendableCoins := acc.SpendableCoins(ctx.BlockTime())
	var activationFee sdk.Coins
	var dstCoins sdk.Coins
	accDst := ak.GetAccount(ctx, msg.ToAddress)
	if accDst == nil {
		activationFee = dex.NewCetCoins(k.GetParams(ctx).ActivationFee)
	} else {
		dstCoins = accDst.GetCoins()
	}
	if srcSpendableCoins.IsAllLT(amt) {
		return sdk.ErrInsufficientCoins("sender has insufficient coins for the transfer").Result()
	}

	//check whether toAccount exist
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	time := msg.UnlockTime
	if time != 0 {
		//tmp invalid for pprof test
		return lockedSend(ctx, k, msg.FromAddress, msg.ToAddress, amt, time)
	}
	return normalSend(ctx, ak, acc, accDst, srcCoins, dstCoins, amt, activationFee)
}

func lockedSend(ctx sdk.Context, k Keeper, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins, unlockTime int64) sdk.Result {
	if unlockTime < ctx.BlockHeader().Time.Unix() {
		return types.ErrUnlockTime("Invalid Unlock Time:" +
			fmt.Sprintf("%d < %d", unlockTime, ctx.BlockHeader().Time.Unix())).Result()
	}

	if err := k.SendLockedCoins(ctx, fromAddr, toAddr, amt, unlockTime); err != nil {
		return err.Result()
	}

	fillMsgQueue(ctx, k, "send_lock_coins", types.NewMsgSend(fromAddr, toAddr, amt, unlockTime))

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func normalSend(ctx sdk.Context, k auth.AccountKeeper, fromAcc, toAcc exported.Account, srcCoins, dstCoins, amt, activationFee sdk.Coins) sdk.Result {

	srcCoins = srcCoins.Sub(amt)
	dstCoins = dstCoins.Add(amt).Sub(activationFee)
	_ = fromAcc.SetCoins(srcCoins)
	_ = toAcc.SetCoins(dstCoins)
	k.SetAccount(ctx, fromAcc)
	k.SetAccount(ctx, toAcc)
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgSetMemoRequired(ctx sdk.Context, k Keeper, msg types.MsgSetMemoRequired) sdk.Result {
	addr := msg.Address
	required := msg.Required

	if k.BlacklistedAddr(msg.Address) {
		return sdk.ErrUnauthorized(fmt.Sprintf("%s is not allowed to receive transactions", msg.Address)).Result()
	}

	if err := k.SetMemoRequired(ctx, addr, required); err != nil {
		return err.Result()
	}

	//fillMsgQueue(ctx, k, "send_memo_required", msg)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAddr, addr.String()),
			sdk.NewAttribute(types.AttributeKeyMemoRequired, fmt.Sprintf("%v", required)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func fillMsgQueue(ctx sdk.Context, keeper Keeper, key string, msg interface{}) {
	if keeper.MsgProducer.IsSubscribed(types.Topic) {
		bytes, err := json.Marshal(msg)
		if err != nil {
			return
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(msgqueue.EventTypeMsgQueue,
			sdk.NewAttribute(key, string(bytes))))
	}
}
