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

func SimulateMsgIssueToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx types.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		tokenName, tokenSymbol, tokenURL, tokenDescrip, tokenIdentify := randomTokenAttrs(r)
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

func SimulateMsgTransferOwnership(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		symbol := RandStringBytes(r, SymbolLenth)
		originOwner := simulation.RandomAcc(r, accounts)
		newOwner := simulation.RandomAcc(r, accounts)
		for newOwner.Equals(originOwner) {
			newOwner = simulation.RandomAcc(r, accounts)
		}
		msg := asset.NewMsgTransferOwnership(symbol, originOwner.Address, newOwner.Address)
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

func SimulateMsgMintToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		mintAmount := r.Intn(50000000) + 100000
		msg := asset.NewMsgMintToken(symbol, sdk.NewInt(int64(mintAmount)), acc.Address)
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

func SimulateMsgBurnToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		mintAmount := r.Intn(5000000) + 100000
		msg := asset.NewMsgBurnToken(symbol, sdk.NewInt(int64(mintAmount)), acc.Address)
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

func SimulateMsgForbidToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		msg := asset.NewMsgForbidToken(symbol, acc.Address)
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

func SimulateMsgUnForbidToken(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		msg := asset.NewMsgUnForbidToken(symbol, acc.Address)
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

func SimulateMsgAddTokenWhitelist(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		whiteList := make([]types.AccAddress, len(accounts)/2)
		for i := 0; i < len(whiteList); i++ {
			whiteList[i] = simulation.RandomAcc(r, accounts).Address
		}
		msg := asset.NewMsgAddTokenWhitelist(symbol, acc.Address, whiteList)
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

func SimulateMsgRemoveTokenWhitelist(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		whiteList := make([]types.AccAddress, len(accounts)/3)
		for i := 0; i < len(whiteList); i++ {
			whiteList[i] = simulation.RandomAcc(r, accounts).Address
		}
		msg := asset.NewMsgRemoveTokenWhitelist(symbol, acc.Address, whiteList)
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

func SimulateMsgForbidAddr(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		forbidList := make([]types.AccAddress, len(accounts)/4)
		for i := 0; i < len(forbidList); i++ {
			forbidList[i] = simulation.RandomAcc(r, accounts).Address
		}
		msg := asset.NewMsgForbidAddr(symbol, acc.Address, forbidList)
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

func SimulateMsgUnForbidAddr(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		forbidList := make([]types.AccAddress, len(accounts)/4)
		for i := 0; i < len(forbidList); i++ {
			forbidList[i] = simulation.RandomAcc(r, accounts).Address
		}
		msg := asset.NewMsgUnForbidAddr(symbol, acc.Address, forbidList)
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

func SimulateMsgModifyTokenInfo(k asset.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)
		symbol := RandStringBytes(r, SymbolLenth)
		url := fmt.Sprintf("www.%s.org", symbol)
		describe := fmt.Sprintf("simulation modify info %s", symbol)
		msg := asset.NewMsgModifyTokenInfo(symbol, url, describe, acc.Address)
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

func RandStringBytes(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func randomTokenAttrs(r *rand.Rand) (name string, symbol string, url string, descrip string, identify string) {
	symbol = RandStringBytes(r, SymbolLenth)
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
