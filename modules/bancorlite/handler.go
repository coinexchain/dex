package bancorlite

import (
	"bytes"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/msgqueue"
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
			errMsg := "Unrecognized bancorlite Msg type: " + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgBancorInit(ctx sdk.Context, k Keeper, msg types.MsgBancorInit) sdk.Result {
	if bi := k.Bik.Load(ctx, msg.Stock+keepers.SymbolSeparator+msg.Money); bi != nil {
		return types.ErrBancorAlreadyExists().Result()
	}
	if !k.Axk.IsTokenExists(ctx, msg.Stock) || !k.Axk.IsTokenExists(ctx, msg.Money) {
		return types.ErrNoSuchToken().Result()
	}
	if !k.Axk.IsTokenIssuer(ctx, msg.Stock, msg.Owner) {
		return types.ErrNonOwnerIsProhibited().Result()
	}
	if msg.Money != "cet" &&
		!k.Mk.IsMarketExist(ctx, msg.Stock+keepers.SymbolSeparator+"cet") {
		return types.ErrNonMarketExist().Result()
	}
	suppliedCoins := sdk.Coins{sdk.NewCoin(msg.Stock, msg.MaxSupply)}
	if err := k.Bxk.FreezeCoins(ctx, msg.Owner, suppliedCoins); err != nil {
		return err.Result()
	}
	fee := k.Bik.GetParams(ctx).CreateBancorFee
	if err := k.Bxk.DeductInt64CetFee(ctx, msg.Owner, fee); err != nil {
		return err.Result()
	}
	bi := &keepers.BancorInfo{
		Owner:              msg.Owner,
		Stock:              msg.Stock,
		Money:              msg.Money,
		InitPrice:          msg.InitPrice,
		MaxSupply:          msg.MaxSupply,
		MaxPrice:           msg.MaxPrice,
		Price:              msg.InitPrice,
		StockInPool:        msg.MaxSupply,
		MoneyInPool:        sdk.ZeroInt(),
		EarliestCancelTime: msg.EarliestCancelTime,
	}
	k.Bik.Save(ctx, bi)

	fillMsgQueue(ctx, k, KafkaBancorInfo, *bi)

	return sdk.Result{}
}

func handleMsgBancorCancel(ctx sdk.Context, k Keeper, msg types.MsgBancorCancel) sdk.Result {
	bi := k.Bik.Load(ctx, msg.Stock+keepers.SymbolSeparator+msg.Money)
	if bi == nil {
		return types.ErrNoBancorExists().Result()
	}
	if !bytes.Equal(bi.Owner, msg.Owner) {
		return types.ErrNotBancorOwner().Result()
	}
	if ctx.BlockHeader().Time.Unix() < bi.EarliestCancelTime {
		return types.ErrEarliestCancelTimeNotArrive().Result()
	}
	if !k.Mk.IsMarketExist(ctx, msg.Stock+keepers.SymbolSeparator+"cet") {
		return types.ErrNonMarketExist().Result()
	}
	fee := k.Bik.GetParams(ctx).CancelBancorFee
	if err := k.Bxk.DeductInt64CetFee(ctx, msg.Owner, fee); err != nil {
		return err.Result()
	}
	k.Bik.Remove(ctx, bi)
	if err := k.Bxk.UnFreezeCoins(ctx, bi.Owner, sdk.NewCoins(sdk.NewCoin(bi.Stock, bi.StockInPool))); err != nil {
		return err.Result()
	}
	if err := k.Bxk.UnFreezeCoins(ctx, bi.Owner, sdk.NewCoins(sdk.NewCoin(bi.Money, bi.MoneyInPool))); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgBancorTrade(ctx sdk.Context, k Keeper, msg types.MsgBancorTrade) sdk.Result {
	bi := k.Bik.Load(ctx, msg.Stock+keepers.SymbolSeparator+msg.Money)
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
	if msg.IsBuy {
		diff = biNew.MoneyInPool.Sub(bi.MoneyInPool)
	}
	coinsFromPool := sdk.Coins{sdk.NewCoin(msg.Money, diff)}
	coinsToPool := sdk.Coins{sdk.NewCoin(msg.Stock, sdk.NewInt(msg.Amount))}
	moneyCrossLimit := msg.MoneyLimit > 0 && diff.LT(sdk.NewInt(msg.MoneyLimit))
	moneyErr := "less than"
	if msg.IsBuy {
		coinsToPool = sdk.Coins{sdk.NewCoin(msg.Money, diff)}
		coinsFromPool = sdk.Coins{sdk.NewCoin(msg.Stock, sdk.NewInt(msg.Amount))}
		moneyCrossLimit = msg.MoneyLimit > 0 && diff.GT(sdk.NewInt(msg.MoneyLimit))
		moneyErr = "more than"
	}

	if moneyCrossLimit {
		return types.ErrMoneyCrossLimit(moneyErr).Result()
	}

	commission, err := getTradeFee(ctx, k, msg, diff)
	if err != nil {
		return err.Result()
	}

	if err := k.Bxk.DeductFee(ctx, msg.Sender, sdk.NewCoins(sdk.NewCoin("cet", commission))); err != nil {
		return err.Result()
	}

	if err = swapStockAndMoney(ctx, k, msg.Sender, bi.Owner, coinsFromPool, coinsToPool); err != nil {
		return err.Result()
	}

	k.Bik.Save(ctx, &biNew)

	sideStr := "sell"
	side := market.SELL
	if msg.IsBuy {
		sideStr = "buy"
		side = market.BUY
	}

	m := types.MsgBancorTradeInfoForKafka{
		Sender:     msg.Sender,
		Stock:      msg.Stock,
		Money:      msg.Money,
		Amount:     msg.Amount,
		Side:       byte(side),
		MoneyLimit: msg.MoneyLimit,
		TxPrice: biNew.Price.Add(bi.Price).QuoInt64(2),
		BlockHeight: ctx.BlockHeight(),
	}
	fillMsgQueue(ctx, k, KafkaBancorTrade, m)
	fillMsgQueue(ctx, k, KafkaBancorInfo, biNew)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeBancorlite,
			sdk.NewAttribute(AttributeKeyTradeFor, bi.Stock+keepers.SymbolSeparator+bi.Money),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(AttributeNewStockInPool, biNew.StockInPool.String()),
			sdk.NewAttribute(AttributeNewMoneyInPool, biNew.MoneyInPool.String()),
			sdk.NewAttribute(AttributeNewPrice, biNew.Price.String()),
			sdk.NewAttribute(AttributeTradeSide, sideStr),
			sdk.NewAttribute(AttributeCoinsFromPool, coinsFromPool.String()),
			sdk.NewAttribute(AttributeCoinsToPool, coinsToPool.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func fillMsgQueue(ctx sdk.Context, keeper Keeper, key string, msg interface{}) {
	if keeper.MsgProducer.IsSubscribed(types.Topic) {
		b, err := json.Marshal(msg)
		if err != nil {
			return
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(msgqueue.EventTypeMsgQueue,
			sdk.NewAttribute(key, string(b))))
	}
}

func getTradeFee(ctx sdk.Context, k keepers.Keeper, msg types.MsgBancorTrade,
	amountOfMoney sdk.Int) (sdk.Int, sdk.Error) {

	var commission sdk.Int
	if msg.Money == "cet" {
		commission = amountOfMoney.
			Mul(sdk.NewInt(k.Bik.GetParams(ctx).TradeFeeRate)).
			Quo(sdk.NewInt(10000))
	} else {
		price, err := k.Mk.GetMarketLastExePrice(ctx, msg.Stock+keepers.SymbolSeparator+"cet")
		if err != nil {
			return commission, types.ErrGetMarketPrice(err.Error())
		}
		commission = price.
			MulInt(sdk.NewInt(msg.Amount)).
			MulInt(sdk.NewInt(k.Bik.GetParams(ctx).TradeFeeRate)).
			QuoInt(sdk.NewInt(10000)).RoundInt()
	}

	if commission.Int64() < k.Mk.GetMarketFeeMin(ctx) {
		return commission, types.ErrTradeQuantityToSmall(commission.Int64())
	}
	return commission, nil
}

func swapStockAndMoney(ctx sdk.Context, k keepers.Keeper, trader sdk.AccAddress, owner sdk.AccAddress,
	coinsFromPool sdk.Coins, coinsToPool sdk.Coins) sdk.Error {
	if err := k.Bxk.SendCoins(ctx, trader, owner, coinsToPool); err != nil {
		return err
	}
	if err := k.Bxk.FreezeCoins(ctx, owner, coinsToPool); err != nil {
		return err
	}
	if err := k.Bxk.UnFreezeCoins(ctx, owner, coinsFromPool); err != nil {
		return err
	}
	if err := k.Bxk.SendCoins(ctx, owner, trader, coinsFromPool); err != nil {
		return err
	}
	return nil
}
