package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/distributionx"
	dexsim "github.com/coinexchain/dex/simulation"
	dex "github.com/coinexchain/dex/types"
)

func SimulateMsgDonateToCommunityPool(ak auth.AccountKeeper, dxk distributionx.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		dbAcc := ak.GetAccount(ctx, acc.Address)
		cetAmt := dbAcc.GetCoins().AmountOf(dex.CET)
		rdmAmt := simulation.RandomAmount(r, cetAmt)

		if rdmAmt.LT(sdk.OneInt()) {
			return simulation.NoOpMsg(distributionx.ModuleName), nil, nil
		}

		msg := distributionx.MsgDonateToCommunityPool{
			FromAddr: acc.Address,
			Amount:   sdk.NewCoins(sdk.NewCoin(dex.CET, rdmAmt)),
		}

		ok := dexsim.SimulateHandleMsg(msg, distributionx.NewHandler(dxk), ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
