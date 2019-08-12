package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/asset"
)

func SimulateMsgIssuerToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx types.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		tokenName, tokenSymbol, tokenURL, tokenDescrip, tokenIdentify := randomTokenSymbol(r)
		issueAmount := randomIssuance(r)
		msg := asset.NewMsgIssueToken(tokenName, tokenSymbol, issueAmount, acc.Address,
			true, true, true, true, tokenURL, tokenDescrip, tokenIdentify)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}

		ctx, write := ctx.CacheContext()
		ok := asset.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}

		opMsg := simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}

func SimulateMsgTransferOwnership() simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		return simulation.OperationMsg{}, nil, nil
	}
}

func RandStringBytes(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func randomTokenSymbol(r *rand.Rand) (name string, symbol string, url string, descrip string, identify string) {
	symbol = RandStringBytes(r, 4)
	name = fmt.Sprintf("simulation-%s", symbol)
	url = fmt.Sprintf("www.%s.com", symbol)
	descrip = fmt.Sprintf("simulation issue token : %s", symbol)
	identify = fmt.Sprintf("%d-simulation", r.Int())
	return
}

func randomIssuance(r *rand.Rand) sdk.Int {
	total := 100000000000000000
	issue := r.Intn(100000000000000000)
	if issue < total/2 {
		issue += total / 4
	}
	return sdk.NewInt(int64(issue))
}
