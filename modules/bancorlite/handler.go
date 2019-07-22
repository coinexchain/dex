package bancorlite

import (
	"bytes"
	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgBancorTrade:
			return handleMsgBancorTrade(ctx, k, msg)
		case types.MsgBancorInit:
			return handleMsgBancorInit(ctx, k, msg)
		default:
			errMsg := "Unrecognized bancorlite Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgBancorInit(ctx sdk.Context, k Keeper, msg types.MsgBancorInit) sdk.Result {
	if bi := k.Bik.Load(ctx, msg.Token); bi != nil {
		return types.ErrBancorAlreadyExists().Result()
	}
	if !k.Axk.IsTokenExists(ctx, msg.Token) {
		return types.ErrNoSuchToken().Result()
	}
	if !k.Axk.IsTokenIssuer(ctx, msg.Token, msg.Owner) {
		return types.ErrNonOwnerIsProhibited().Result()
	}
	coins := sdk.Coins{sdk.Coin{Denom: msg.Token, Amount: msg.MaxSupply}}
	if err := k.Bxk.FreezeCoins(ctx, msg.Owner, coins); err != nil {
		return err.Result()
	}
	bi := &keepers.BancorInfo{
		Owner:       msg.Owner,
		Token:       msg.Token,
		MaxSupply:   msg.MaxSupply,
		MaxPrice:    msg.MaxPrice,
		StockInPool: msg.MaxSupply,
		MoneyInPool: sdk.ZeroInt(),
	}
	k.Bik.Save(ctx, bi)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeBancorlite,
			sdk.NewAttribute(AttributeKeyCreateFor, bi.Token),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeOwner, bi.Owner.String()),
			sdk.NewAttribute(AttributeMaxSupply, bi.MaxSupply.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBancorTrade(ctx sdk.Context, k Keeper, msg types.MsgBancorTrade) sdk.Result {
	bi := k.Bik.Load(ctx, msg.Token)
	if bi == nil {
		return types.ErrNoBancorExists().Result()
	}
	if bytes.Equal(bi.Owner, msg.Sender) {
		return types.ErrOwnerIsProhibited().Result()
	}

	stockInPool := bi.StockInPool.AddRaw(msg.Amount)
	if msg.IsBuy {
		stockInPool = bi.StockInPool.SubRaw(msg.Amount)
	}
	biNew := *bi
	if ok := biNew.UpdateStockInPool(stockInPool); !ok {
		return types.ErrStockInPoolOutofBound().Result()
	}

	diff := bi.MoneyInPool.Sub(biNew.MoneyInPool)
	coinsFromPool := sdk.Coins{sdk.Coin{Denom: "cet", Amount: diff}}
	coinsToPool := sdk.Coins{sdk.Coin{Denom: msg.Token, Amount: sdk.NewInt(msg.Amount)}}
	moneyCrossLimit := msg.MoneyLimit > 0 && diff.LT(sdk.NewInt(msg.MoneyLimit))
	moneyErr := "less than"
	if msg.IsBuy {
		diff = biNew.MoneyInPool.Sub(bi.MoneyInPool)
		coinsToPool = sdk.Coins{sdk.Coin{Denom: "cet", Amount: diff}}
		coinsFromPool = sdk.Coins{sdk.Coin{Denom: msg.Token, Amount: sdk.NewInt(msg.Amount)}}
		moneyCrossLimit = msg.MoneyLimit > 0 && diff.GT(sdk.NewInt(msg.MoneyLimit))
		moneyErr = "more than"
	}

	if moneyCrossLimit {
		return types.ErrMoneyCrossLimit(moneyErr).Result()
	}

	if err := k.Bxk.SendCoins(ctx, msg.Sender, bi.Owner, coinsToPool); err != nil {
		return err.Result()
	}
	if err := k.Bxk.FreezeCoins(ctx, bi.Owner, coinsToPool); err != nil {
		return err.Result()
	}
	if err := k.Bxk.UnFreezeCoins(ctx, bi.Owner, coinsFromPool); err != nil {
		return err.Result()
	}
	if err := k.Bxk.SendCoins(ctx, bi.Owner, msg.Sender, coinsFromPool); err != nil {
		return err.Result()
	}
	k.Bik.Save(ctx, &biNew)

	side := "Sell"
	if msg.IsBuy {
		side = "Buy"
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeBancorlite,
			sdk.NewAttribute(AttributeKeyTradeFor, bi.Token),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeNewStockInPool, biNew.StockInPool.String()),
			sdk.NewAttribute(AttributeNewMoneyInPool, biNew.MoneyInPool.String()),
			sdk.NewAttribute(AttributeNewPrice, biNew.Price.String()),
			sdk.NewAttribute(AttributeNewPrice, biNew.Price.String()),
			sdk.NewAttribute(AttributeTradeSide, side),
			sdk.NewAttribute(AttributeCoinsFromPool, coinsFromPool.String()),
			sdk.NewAttribute(AttributeCoinsToPool, coinsToPool.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
