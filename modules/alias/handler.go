package alias

import (
	"bytes"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx.WithEventManager(sdk.NewEventManager())
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
		if addr := k.GetAddressFromAlias(ctx, msg.Alias); len(addr) != 0 {
			return types.ErrAliasAlreadyExists().Result()
		}
		k.AddAlias(ctx, msg.Alias, msg.Owner)
	} else {
		if addr := k.GetAddressFromAlias(ctx, msg.Alias); !bytes.Equal(addr, msg.Owner) {
			//fmt.Printf("%x vs %x\n", addr, []byte(msg.Owner))
			return types.ErrNoSuchAlias().Result()
		}
		k.RemoveAlias(ctx, msg.Alias, msg.Owner)
	}
	return sdk.Result{} //TODO
}
