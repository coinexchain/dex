package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case MsgCreateMarketInfo:
			return handlerMsgCreateMarketinfo(ctx, msg, k)
		case MsgCreateGTEOrder:
			return handlerMsgCreateGTEOrder(ctx, msg, k)
		default:
			errMsg := "Unrecognized market Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerMsgCreateMarketinfo(ctx sdk.Context, msg MsgCreateMarketInfo, keeper Keeper) sdk.Result {
	return sdk.Result{}
}

func handlerMsgCreateGTEOrder(ctx sdk.Context, msg MsgCreateGTEOrder, keeper Keeper) sdk.Result {
	return sdk.Result{}
}
