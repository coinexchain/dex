package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/bancorlite"
	dexsim "github.com/coinexchain/dex/simulation"
	dex "github.com/coinexchain/dex/types"
)

var symbolPrefix = "bl" // bancor_lite

func SimulateMsgBancorInit(assetKeeper asset.Keeper, blk bancorlite.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		addr := simulation.RandomAcc(r, accs).Address
		newSymbol := dexsim.RandomSymbol(r, symbolPrefix, 3)
		identity := dexsim.RandomSymbol(r, "", 30)
		amount := 1e10 + r.Int63n(1e12)

		// issue new token
		issueTokenMsg := createMsgIssueToken(newSymbol, amount, addr, identity)
		issueTokenOK := dexsim.SimulateHandleMsg(issueTokenMsg, asset.NewHandler(assetKeeper), ctx)
		if !issueTokenOK {
			return simulation.NoOpMsg(asset.ModuleName), nil, nil
		}

		// init trading pair
		bancorInitMsg := createMsgBancorInit(r, addr, newSymbol, amount)
		bancorInitOk := dexsim.SimulateHandleMsg(bancorInitMsg, bancorlite.NewHandler(blk), ctx)
		opMsg = simulation.NewOperationMsg(bancorInitMsg, bancorInitOk, "")
		if !bancorInitOk {
			return opMsg, nil, nil
		}

		//verify bancor init
		ok := verifyBancorInit(ctx, blk, bancorInitMsg)
		if !ok {
			return simulation.NewOperationMsg(bancorInitMsg, ok, ""), nil, fmt.Errorf("bancor initialization failed")
		}
		return simulation.NewOperationMsg(bancorInitMsg, ok, ""), nil, nil
	}
}

func createMsgIssueToken(newSymbol string, amount int64, tokenOwner sdk.AccAddress, identity string) asset.MsgIssueToken {
	return asset.NewMsgIssueToken(newSymbol, newSymbol, sdk.NewInt(amount), tokenOwner,
		false, false, false, false, "", "", identity)
}

func createMsgBancorInit(r *rand.Rand,
	owner sdk.AccAddress, stockSymbol string, stockTotalAmount int64) bancorlite.MsgBancorInit {

	maxSupply := stockTotalAmount / 2
	initPrice := r.Int63n(5) // give 0 more chances
	if initPrice > 0 {
		initPrice = r.Int63n(1000)
	}
	maxPrice := (initPrice + 1) * 1000 //  make maxPrice not 0 when initPrice is 0

	return bancorlite.MsgBancorInit{
		Owner:              owner,
		Stock:              stockSymbol,
		Money:              dex.DefaultBondDenom,
		InitPrice:          sdk.NewDec(initPrice),
		MaxPrice:           sdk.NewDec(maxPrice),
		MaxSupply:          sdk.NewInt(maxSupply),
		EarliestCancelTime: 0, // TODO
	}
}
func verifyBancorInit(ctx sdk.Context, keeper bancorlite.Keeper, msg bancorlite.MsgBancorInit) bool {
	bancorInfo := keeper.Bik.Load(ctx, msg.GetSymbol())
	return bancorInfo.Stock == msg.Stock &&
		bancorInfo.Money == msg.Money &&
		bancorInfo.Owner.Equals(msg.Owner) &&
		bancorInfo.EarliestCancelTime == msg.EarliestCancelTime &&
		bancorInfo.InitPrice.Equal(msg.InitPrice) &&
		bancorInfo.MaxPrice.Equal(msg.MaxPrice) &&
		bancorInfo.MaxSupply.Equal(msg.MaxSupply) &&
		bancorInfo.StockInPool.Equal(msg.MaxSupply) &&
		bancorInfo.MoneyInPool.IsZero() &&
		bancorInfo.Price.Equal(msg.InitPrice)
}

func SimulateMsgBancorTrade(ak auth.AccountKeeper, blk bancorlite.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		addr := simulation.RandomAcc(r, accs).Address
		dbAcc := ak.GetAccount(ctx, addr)
		saleableCoins := getSaleableCoins(dbAcc)
		if len(saleableCoins) > 0 && r.Intn(2) > 0 {
			return simulateMsgBancorSell(blk, r, ctx, addr, saleableCoins)
		}
		return simulateMsgBancorBuy(blk, r, ctx, addr)
	}
}

func simulateMsgBancorSell(blk bancorlite.Keeper,
	r *rand.Rand, ctx sdk.Context, addr sdk.AccAddress, saleableCoins sdk.Coins,
) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

	saleableCoin := saleableCoins[r.Intn(len(saleableCoins))]
	amount := r.Int63n(saleableCoin.Amount.Int64() / 2)
	msg := bancorlite.MsgBancorTrade{
		Sender:     addr,
		Stock:      saleableCoin.Denom,
		Money:      dex.DefaultBondDenom,
		IsBuy:      false,
		Amount:     amount,
		MoneyLimit: 0, // TODO
	}
	oldBancorInfo := blk.Bik.Load(ctx, msg.GetSymbol())
	ok := dexsim.SimulateHandleMsg(msg, bancorlite.NewHandler(blk), ctx)
	opMsg = simulation.NewOperationMsg(msg, ok, "")
	if !ok {
		return opMsg, nil, nil
	}
	ok = verifyBancorSellBuy(ctx, blk, oldBancorInfo, msg, false)
	if !ok {
		return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("bancor sell failed")
	}
	return simulation.NewOperationMsg(msg, ok, ""), nil, nil

}
func verifyBancorSellBuy(ctx sdk.Context, blk bancorlite.Keeper, oldBancorInfo *bancorlite.BancorInfo, msg bancorlite.MsgBancorTrade, isbuy bool) bool {
	updatedBancorInfo := blk.Bik.Load(ctx, msg.GetSymbol())
	stockInPool := oldBancorInfo.StockInPool.Add(sdk.NewInt(msg.Amount))
	if isbuy {
		stockInPool = oldBancorInfo.StockInPool.Sub(sdk.NewInt(msg.Amount))
	}
	ok := oldBancorInfo.UpdateStockInPool(stockInPool)
	return ok &&
		oldBancorInfo.StockInPool.Equal(updatedBancorInfo.StockInPool) &&
		oldBancorInfo.MoneyInPool.Equal(updatedBancorInfo.MoneyInPool) &&
		oldBancorInfo.Price.Equal(updatedBancorInfo.Price)
}
func simulateMsgBancorBuy(blk bancorlite.Keeper,
	r *rand.Rand, ctx sdk.Context, addr sdk.AccAddress,
) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

	bancorInfo := randomBancorInfo(r, blk, ctx)
	if bancorInfo == nil || !bancorInfo.StockInPool.IsPositive() {
		return simulation.NoOpMsg(bancorlite.ModuleName), nil, nil
	}

	amount := r.Int63n(bancorInfo.StockInPool.Int64() / 2)
	msg := bancorlite.MsgBancorTrade{
		Sender:     addr,
		Stock:      bancorInfo.Stock,
		Money:      bancorInfo.Money,
		IsBuy:      true,
		Amount:     amount,
		MoneyLimit: math.MaxInt64, // TODO
	}
	oldBancorInfo := blk.Bik.Load(ctx, msg.GetSymbol())
	ok := dexsim.SimulateHandleMsg(msg, bancorlite.NewHandler(blk), ctx)
	opMsg = simulation.NewOperationMsg(msg, ok, "")
	if !ok {
		return opMsg, nil, nil
	}
	ok = verifyBancorSellBuy(ctx, blk, oldBancorInfo, msg, true)
	if !ok {
		return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("bancor buy failed")
	}
	return simulation.NewOperationMsg(msg, ok, ""), nil, nil

}

func getSaleableCoins(acc auth.Account) sdk.Coins {
	saleableCoins := sdk.Coins{}
	for _, coin := range acc.GetCoins() {
		if strings.HasPrefix(coin.Denom, symbolPrefix) && coin.Amount.Int64() > 1 {
			saleableCoins = append(saleableCoins, coin)
		}
	}
	return saleableCoins
}

func randomBancorInfo(r *rand.Rand, blk bancorlite.Keeper, ctx sdk.Context) *bancorlite.BancorInfo {
	bis := getAllBancorInfos(blk, ctx)
	switch n := len(bis); n {
	case 0:
		return nil
	case 1:
		return bis[0]
	default:
		return bis[r.Intn(n)]
	}
}

func getAllBancorInfos(blk bancorlite.Keeper, ctx sdk.Context) []*bancorlite.BancorInfo {
	bis := make([]*bancorlite.BancorInfo, 0, 100)
	blk.Bik.Iterate(ctx, func(bi *bancorlite.BancorInfo) {
		bis = append(bis, bi)
	})
	return bis
}

func SimulateMsgBancorCancel(blk bancorlite.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		bancorInfo := randomBancorInfo(r, blk, ctx)
		if bancorInfo == nil {
			return simulation.NoOpMsg(bancorlite.ModuleName), nil, nil
		}

		msg := bancorlite.MsgBancorCancel{
			Owner: bancorInfo.Owner,
			Stock: bancorInfo.Stock,
			Money: bancorInfo.Money,
		}

		ok := dexsim.SimulateHandleMsg(msg, bancorlite.NewHandler(blk), ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		if !ok {
			return opMsg, nil, nil
		}
		ok = verifyBancorCancel(ctx, blk, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("bancor cancel failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}
func verifyBancorCancel(ctx sdk.Context, blk bancorlite.Keeper, msg bancorlite.MsgBancorCancel) bool {
	bancorInfo := blk.Bik.Load(ctx, msg.GetSymbol())
	return bancorInfo == nil
}
