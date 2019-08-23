package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/bankx"
	dexsim "github.com/coinexchain/dex/simulation"
)

func SimulateMsgSetMemoRequired(k bankx.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx types.Context, accounts []simulation.Account) (
		opMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)

		msg := bankx.NewMsgSetTransferMemoRequired(acc.Address, r.Intn(2) > 0)
		ok := dexsim.SimulateHandleMsg(msg, bankx.NewHandler(k), ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
