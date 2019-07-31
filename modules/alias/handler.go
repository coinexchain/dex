package alias

import (
	"bytes"
	"fmt"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	dexsdk "github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgAliasUpdate:
			return handleMsgAliasUpdate(ctx, k, msg)
		default:
			errMsg := "Unrecognized alias Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgAliasUpdate(ctx sdk.Context, k Keeper, msg types.MsgAliasUpdate) sdk.Result {
	if msg.IsAdd {
		if types.IsOnlyForCoinEx(msg.Alias) && !k.AssetKeeper.IsTokenIssuer(ctx, "cet", msg.Owner) {
			return types.ErrCanOnlyBeUsedByCetOwner(msg.Alias).Result()
		}
		addr, asDefault := k.AliasKeeper.GetAddressFromAlias(ctx, msg.Alias)
		if len(addr) != 0 && asDefault == msg.AsDefault {
			return types.ErrAliasAlreadyExists().Result()
		}
		aliasParams := k.GetParams(ctx)
		ok, addNewAlias := k.AliasKeeper.AddAlias(ctx, msg.Alias, msg.Owner, msg.AsDefault, aliasParams.MaxAliasCount)
		if !ok {
			return types.ErrMaxAliasCountReached().Result()
		} else if addNewAlias {
			var coins sdk.Coins
			if len(msg.Alias) == 2 {
				coins = dexsdk.NewCetCoins(aliasParams.FeeForAliasLength2)
			} else if len(msg.Alias) == 3 {
				coins = dexsdk.NewCetCoins(aliasParams.FeeForAliasLength3)
			} else if len(msg.Alias) == 4 {
				coins = dexsdk.NewCetCoins(aliasParams.FeeForAliasLength4)
			} else if len(msg.Alias) == 5 {
				coins = dexsdk.NewCetCoins(aliasParams.FeeForAliasLength5)
			} else if len(msg.Alias) == 6 {
				coins = dexsdk.NewCetCoins(aliasParams.FeeForAliasLength6)
			} else {
				coins = dexsdk.NewCetCoins(aliasParams.FeeForAliasLength7OrHigher)
			}
			err := k.BankKeeper.DeductFee(ctx, msg.Owner, coins)
			if err != nil {
				return err.Result()
			}

		}
	} else {
		addr, _ := k.AliasKeeper.GetAddressFromAlias(ctx, msg.Alias)
		if !bytes.Equal(addr, msg.Owner) {
			//fmt.Printf("%x vs %x\n", addr, []byte(msg.Owner))
			return types.ErrNoSuchAlias().Result()
		}
		k.AliasKeeper.RemoveAlias(ctx, msg.Alias, msg.Owner)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	)
	if msg.IsAdd {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAddAlias,
				sdk.NewAttribute(types.AttributeKeyAlias, msg.Alias),
				sdk.NewAttribute(types.AttributeKeyAsDefault, fmt.Sprintf("%t", msg.AsDefault)),
			),
		)
	} else {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRemoveAlias,
				sdk.NewAttribute(types.AttributeKeyAlias, msg.Alias),
			),
		)
	}

	return sdk.Result{
		Codespace: types.CodeSpaceAlias,
		Events:    ctx.EventManager().Events(),
	}
}
