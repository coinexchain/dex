package simulation

import (
	"math"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/bancorlite"
)

func SimulateMsgBancorInit(assetKeeper asset.Keeper, blk bancorlite.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		tokenOwner := simulation.RandomAcc(r, accs).Address
		newSymbol := randomSymbol(r, "bl", 3)
		amount := 1e10 + r.Int63n(1e12)

		// issue new token
		issueTokenMsg := createMsgIssueToken(newSymbol, amount, tokenOwner)
		issueTokenOK := simulateHandleMsg(issueTokenMsg, asset.NewHandler(assetKeeper), ctx)
		if issueTokenOK {
			return simulation.NoOpMsg(), nil, nil
		}

		// create bancor
		bancorInitMsg := createMsgBancorInit(r, tokenOwner, newSymbol, amount)
		bancorInitOk := simulateHandleMsg(bancorInitMsg, bancorlite.NewHandler(blk), ctx)
		opMsg = simulation.NewOperationMsg(bancorInitMsg, bancorInitOk, "")
		return opMsg, nil, nil
	}
}

func createMsgIssueToken(newSymbol string, amount int64, tokenOwner sdk.AccAddress) asset.MsgIssueToken {
	return asset.NewMsgIssueToken(newSymbol, newSymbol, sdk.NewInt(amount), tokenOwner,
		false, false, false, false, "", "", "")
}

func createMsgBancorInit(r *rand.Rand,
	owner sdk.AccAddress, stockSymbol string, stockTotalAmount int64) bancorlite.MsgBancorInit {

	maxSupply := stockTotalAmount / 2
	initPrice := r.Int63n(5) // give 0 more chances
	if initPrice > 0 {
		initPrice = r.Int63n(1000)
	}
	maxPrice := initPrice * 1000

	return bancorlite.MsgBancorInit{
		Owner:              owner,
		Stock:              stockSymbol,
		Money:              sdk.DefaultBondDenom, // TODO
		InitPrice:          sdk.NewDec(initPrice),
		MaxPrice:           sdk.NewDec(maxPrice),
		MaxSupply:          sdk.NewInt(maxSupply),
		EarliestCancelTime: 0, // TODO
	}
}

func SimulateMsgBancorTrade(blk bancorlite.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		sender := simulation.RandomAcc(r, accs).Address
		bancorInfo := randomBancorInfo(r, blk, ctx)
		if bancorInfo == nil || !bancorInfo.StockInPool.IsPositive() {
			return simulation.NoOpMsg(), nil, nil
		}

		amount := r.Int63n(bancorInfo.StockInPool.Int64() / 2)
		msg := bancorlite.MsgBancorTrade{
			Sender:     sender,
			Stock:      bancorInfo.Stock,
			Money:      bancorInfo.Money,
			IsBuy:      true, // TODO
			Amount:     amount,
			MoneyLimit: math.MaxInt64, // TODO
		}
		ok := simulateHandleMsg(msg, bancorlite.NewHandler(blk), ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
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
		msg := bancorlite.MsgBancorCancel{
			Owner: bancorInfo.Owner,
			Stock: bancorInfo.Stock,
			Money: bancorInfo.Money,
		}

		ok := simulateHandleMsg(msg, bancorlite.NewHandler(blk), ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
