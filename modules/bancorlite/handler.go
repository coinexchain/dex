package bancorlite

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/msgqueue"
	dex "github.com/coinexchain/dex/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgBancorInit:
			return handleMsgBancorInit(ctx, k, msg)
		case types.MsgBancorTrade:
			return handleMsgBancorTrade(ctx, k, msg)
		case types.MsgBancorCancel:
			return handleMsgBancorCancel(ctx, k, msg)
		default:
			return dex.ErrUnknownRequest(ModuleName, msg)
		}
	}
}

func handleMsgBancorInit(ctx sdk.Context, k Keeper, msg types.MsgBancorInit) sdk.Result {
	if bi := k.Load(ctx, msg.GetSymbol()); bi != nil {
		return types.ErrBancorAlreadyExists().Result()
	}
	if !k.IsTokenExists(ctx, msg.Stock) || !k.IsTokenExists(ctx, msg.Money) {
		return types.ErrNoSuchToken().Result()
	}
	if !k.IsTokenIssuer(ctx, msg.Stock, msg.Owner) {
		return types.ErrNonOwnerIsProhibited().Result()
	}
	suppliedCoins := sdk.Coins{sdk.NewCoin(msg.Stock, msg.MaxSupply)}
	if err := k.FreezeCoins(ctx, msg.Owner, suppliedCoins); err != nil {
		return err.Result()
	}
	fee := k.GetParams(ctx).CreateBancorFee
	if err := k.DeductInt64CetFee(ctx, msg.Owner, fee); err != nil {
		return err.Result()
	}
	var precision byte
	if msg.StockPrecision <= 8 {
		precision = msg.StockPrecision
	}
	initPrice, err := sdk.NewDecFromStr(msg.InitPrice)
	if err != nil {
		return types.ErrPriceFmt().Result()
	}
	maxPrice, err := sdk.NewDecFromStr(msg.MaxPrice)
	if err != nil {
		return types.ErrPriceFmt().Result()
	}

	ar := types.CalculateAR(msg, initPrice, maxPrice)

	bi := &keepers.BancorInfo{
		Owner:              msg.Owner,
		Stock:              msg.Stock,
		Money:              msg.Money,
		InitPrice:          initPrice,
		MaxSupply:          msg.MaxSupply,
		StockPrecision:     precision,
		MaxPrice:           maxPrice,
		MaxMoney:           msg.MaxMoney,
		AR:                 ar,
		Price:              initPrice,
		StockInPool:        msg.MaxSupply,
		MoneyInPool:        sdk.ZeroInt(),
		EarliestCancelTime: msg.EarliestCancelTime,
	}
	k.Save(ctx, bi)

	fillMsgQueue(ctx, k, KafkaBancorInfo, *bi)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeKeyBancorInit,
			sdk.NewAttribute(AttributeSymbol, bi.GetSymbol()),
			sdk.NewAttribute(AttributeOwner, bi.Owner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBancorCancel(ctx sdk.Context, k Keeper, msg types.MsgBancorCancel) sdk.Result {
	bi := k.Load(ctx, msg.GetSymbol())
	if bi == nil {
		return types.ErrNoBancorExists().Result()
	}
	if !bytes.Equal(bi.Owner, msg.Owner) {
		return types.ErrNotBancorOwner().Result()
	}
	if ctx.BlockHeader().Time.Unix() < bi.EarliestCancelTime {
		return types.ErrEarliestCancelTimeNotArrive().Result()
	}
	fee := k.GetParams(ctx).CancelBancorFee
	if err := k.DeductInt64CetFee(ctx, msg.Owner, fee); err != nil {
		return err.Result()
	}
	k.Remove(ctx, bi)
	if err := k.UnFreezeCoins(ctx, bi.Owner, sdk.NewCoins(sdk.NewCoin(bi.Stock, bi.StockInPool))); err != nil {
		return err.Result()
	}
	if err := k.UnFreezeCoins(ctx, bi.Owner, sdk.NewCoins(sdk.NewCoin(bi.Money, bi.MoneyInPool))); err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeKeyBancorCancel,
			sdk.NewAttribute(AttributeSymbol, bi.GetSymbol()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgBancorTrade(ctx sdk.Context, k Keeper, msg types.MsgBancorTrade) sdk.Result {
	bi := k.Load(ctx, msg.GetSymbol())
	if bi == nil {
		return types.ErrNoBancorExists().Result()
	}
	if bytes.Equal(bi.Owner, msg.Sender) {
		return types.ErrOwnerIsProhibited().Result()
	}
	if k.IsForbiddenByTokenIssuer(ctx, bi.Stock, msg.Sender) ||
		k.IsForbiddenByTokenIssuer(ctx, bi.Money, msg.Sender) ||
		k.IsForbiddenByTokenIssuer(ctx, bi.Stock, bi.Owner) ||
		k.IsForbiddenByTokenIssuer(ctx, bi.Money, bi.Owner) {
		return types.ErrTokenForbiddenByOwner().Result()
	}
	if !types.CheckStockPrecision(sdk.NewInt(msg.Amount), bi.StockPrecision) {
		return types.ErrStockAmountPrecisionNotMatch().Result()
	}
	stockInPool := bi.StockInPool.AddRaw(msg.Amount)
	if msg.IsBuy {
		stockInPool = bi.StockInPool.SubRaw(msg.Amount)
	}
	biNew := *bi
	if ok := biNew.UpdateStockInPool(stockInPool); !ok {
		return types.ErrStockInPoolOutofBound().Result()
	}

	var (
		diff            sdk.Int
		coinsFromPool   sdk.Coins
		coinsToPool     sdk.Coins
		moneyCrossLimit bool
		moneyErr        string
	)

	if msg.IsBuy {
		diff = biNew.MoneyInPool.Sub(bi.MoneyInPool)
		if !diff.IsPositive() {
			return types.ErrTradeMoneyNotPositive().Result()
		}
		coinsToPool = sdk.Coins{sdk.NewCoin(msg.Money, diff)}
		coinsFromPool = sdk.Coins{sdk.NewCoin(msg.Stock, sdk.NewInt(msg.Amount))}
		moneyCrossLimit = msg.MoneyLimit > 0 && diff.GT(sdk.NewInt(msg.MoneyLimit))
		moneyErr = "more than"
	} else {
		diff = bi.MoneyInPool.Sub(biNew.MoneyInPool)
		if !diff.IsPositive() {
			return types.ErrTradeMoneyNotPositive().Result()
		}
		coinsFromPool = sdk.Coins{sdk.NewCoin(msg.Money, diff)}
		coinsToPool = sdk.Coins{sdk.NewCoin(msg.Stock, sdk.NewInt(msg.Amount))}
		moneyCrossLimit = msg.MoneyLimit > 0 && diff.LT(sdk.NewInt(msg.MoneyLimit))
		moneyErr = "less than"
	}

	if moneyCrossLimit {
		return types.ErrMoneyCrossLimit(moneyErr).Result()
	}

	commission, err := getTradeFee(ctx, k, msg, diff)
	if err != nil {
		return err.Result()
	}

	if err := k.DeductFee(ctx, msg.Sender, sdk.NewCoins(sdk.NewCoin(dex.CET, commission))); err != nil {
		return err.Result()
	}

	if err = swapStockAndMoney(ctx, k, msg.Sender, bi.Owner, coinsFromPool, coinsToPool); err != nil {
		return err.Result()
	}

	k.Save(ctx, &biNew)

	sideStr := "sell"
	side := market.SELL
	if msg.IsBuy {
		sideStr = "buy"
		side = market.BUY
	}

	m := types.MsgBancorTradeInfoForKafka{
		Sender:      msg.Sender,
		Stock:       msg.Stock,
		Money:       msg.Money,
		Amount:      msg.Amount,
		Side:        byte(side),
		MoneyLimit:  msg.MoneyLimit,
		TxPrice:     sdk.NewDecFromInt(diff).QuoInt64(msg.Amount),
		BlockHeight: ctx.BlockHeight(),
	}
	fillMsgQueue(ctx, k, KafkaBancorTrade, m)
	fillMsgQueue(ctx, k, KafkaBancorInfo, biNew)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeKeyBancorTrade,
			sdk.NewAttribute(AttributeSymbol, bi.GetSymbol()),
			sdk.NewAttribute(AttributeNewStockInPool, biNew.StockInPool.String()),
			sdk.NewAttribute(AttributeNewMoneyInPool, biNew.MoneyInPool.String()),
			sdk.NewAttribute(AttributeNewPrice, biNew.Price.String()),
			sdk.NewAttribute(AttributeTradeSide, sideStr),
			sdk.NewAttribute(AttributeCoinsFromPool, coinsFromPool.String()),
			sdk.NewAttribute(AttributeCoinsToPool, coinsToPool.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func fillMsgQueue(ctx sdk.Context, keeper Keeper, key string, msg interface{}) {
	if keeper.IsSubscribed(types.Topic) {
		msgqueue.FillMsgs(ctx, key, msg)
	}
}

func getTradeFee(ctx sdk.Context, k keepers.Keeper, msg types.MsgBancorTrade,
	amountOfMoney sdk.Int) (sdk.Int, sdk.Error) {

	var commission sdk.Int
	if msg.Money == "cet" {
		commission = amountOfMoney.
			Mul(sdk.NewInt(k.GetParams(ctx).TradeFeeRate)).
			Quo(sdk.NewInt(10000))
	} else {
		price, err := k.GetMarketLastExePrice(ctx, dex.GetSymbol(msg.Stock, dex.CET))
		if err != nil {
			return commission, types.ErrGetMarketPrice(err.Error())
		}
		commission = price.
			MulInt(sdk.NewInt(msg.Amount)).
			MulInt(sdk.NewInt(k.GetParams(ctx).TradeFeeRate)).
			QuoInt(sdk.NewInt(10000)).RoundInt()
	}
	if commission.Int64() < k.GetMarketFeeMin(ctx) {
		return commission, types.ErrTradeQuantityTooSmall(commission.Int64())
	}
	return commission, nil
}

func swapStockAndMoney(ctx sdk.Context, k keepers.Keeper, trader sdk.AccAddress, owner sdk.AccAddress,
	coinsFromPool sdk.Coins, coinsToPool sdk.Coins) sdk.Error {
	if err := k.SendCoins(ctx, trader, owner, coinsToPool); err != nil {
		return err
	}
	if err := k.FreezeCoins(ctx, owner, coinsToPool); err != nil {
		return err
	}
	if err := k.UnFreezeCoins(ctx, owner, coinsFromPool); err != nil {
		return err
	}
	if err := k.SendCoins(ctx, owner, trader, coinsFromPool); err != nil {
		return err
	}
	return nil
}
