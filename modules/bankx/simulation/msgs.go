package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
	dexsim "github.com/coinexchain/dex/simulation"
	dex "github.com/coinexchain/dex/types"
)

// SendTx tests and runs a single msg send where both
// accounts already exist.
func SimulateMsgSend(mapper auth.AccountKeeper, bk bankx.Keeper) simulation.Operation {
	handler := bankx.NewHandler(bk)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		lockCoinsFee := dex.NewCetCoins(bk.GetParams(ctx).LockCoinsFee)

		fromAcc, comment, msg, ok := createMsgSend(r, ctx, accs, mapper, lockCoinsFee)
		opMsg = simulation.NewOperationMsg(msg, ok, comment)
		if !ok {
			return opMsg, nil, nil
		}

		fOps, err = sendAndVerifyMsgSend(app, mapper, msg, ctx, []crypto.PrivKey{fromAcc.PrivKey}, handler, lockCoinsFee)
		if err != nil {
			return opMsg, fOps, err
		}
		return opMsg, fOps, nil
	}
}

func createMsgSend(r *rand.Rand, ctx sdk.Context, accs []simulation.Account, mapper auth.AccountKeeper, lockCoinsFee sdk.Coins) (
	fromAcc simulation.Account, comment string, msg types.MsgSend, ok bool) {

	fromAcc = simulation.RandomAcc(r, accs)
	toAcc := simulation.RandomAcc(r, accs)
	// Disallow sending money to yourself
	for {
		if !fromAcc.PubKey.Equals(toAcc.PubKey) {
			break
		}
		toAcc = simulation.RandomAcc(r, accs)
	}
	initFromCoins := mapper.GetAccount(ctx, fromAcc.Address).SpendableCoins(ctx.BlockHeader().Time)

	if len(initFromCoins) == 0 {
		return fromAcc, "skipping, no coins at all", msg, false
	}

	denomIndex := r.Intn(len(initFromCoins))
	amt, goErr := simulation.RandPositiveInt(r, initFromCoins[denomIndex].Amount)
	if goErr != nil {
		return fromAcc, "skipping bank send due to account having no coins of denomination " + initFromCoins[denomIndex].Denom, msg, false
	}

	coins := sdk.Coins{sdk.NewCoin(initFromCoins[denomIndex].Denom, amt)}

	var unlockTime int64
	if initFromCoins.Sub(coins).IsAllGT(lockCoinsFee) {
		if r.Intn(10) > 0 { // 10% of sends are locked within 1 minute
			unlockTime = ctx.BlockHeader().Time.Unix() + r.Int63n(60) + 1
		}
	}

	msg = bankx.NewMsgSend(fromAcc.Address, toAcc.Address, coins, unlockTime)
	return fromAcc, "", msg, true
}

// Sends and verifies the transition of a msg send.
func sendAndVerifyMsgSend(app *baseapp.BaseApp, mapper auth.AccountKeeper, msg types.MsgSend, ctx sdk.Context, privkeys []crypto.PrivKey, handler sdk.Handler, lockCoinsFee sdk.Coins) ([]simulation.FutureOperation, error) {
	fromAcc := mapper.GetAccount(ctx, msg.FromAddress)
	AccountNumbers := []uint64{fromAcc.GetAccountNumber()}
	SequenceNumbers := []uint64{fromAcc.GetSequence()}
	initialFromAddrCoins := fromAcc.GetCoins()

	toAcc := mapper.GetAccount(ctx, msg.ToAddress)
	initialToAddrCoins := toAcc.GetCoins()

	if handler != nil {
		res := handler(ctx, msg)
		if !res.IsOK() {
			if res.Code == bank.CodeSendDisabled || res.Code == types.CodeTokenForbiddenByOwner {
				return nil, nil
			}
			// TODO: Do this in a more 'canonical' way
			return nil, fmt.Errorf("handling msg failed %v", res)
		}
	} else {
		tx := mock.GenTx([]sdk.Msg{msg},
			AccountNumbers,
			SequenceNumbers,
			privkeys...)
		res := app.Deliver(tx)
		if !res.IsOK() {
			// TODO: Do this in a more 'canonical' way
			return nil, fmt.Errorf("Deliver failed %v", res)
		}
	}

	fromAcc = mapper.GetAccount(ctx, msg.FromAddress)
	toAcc = mapper.GetAccount(ctx, msg.ToAddress)

	if msg.UnlockTime == 0 {
		if !initialFromAddrCoins.Sub(msg.Amount).IsEqual(fromAcc.GetCoins()) {
			return nil, fmt.Errorf("fromAddress %s had an incorrect amount of coins", fromAcc.GetAddress())
		}

		if !initialToAddrCoins.Add(msg.Amount).IsEqual(toAcc.GetCoins()) {
			return nil, fmt.Errorf("toAddress %s had an incorrect amount of coins", toAcc.GetAddress())
		}
	} else {
		fOps := []simulation.FutureOperation{
			{
				BlockTime: time.Unix(msg.UnlockTime, 0),
				Op:        checkLockSend(mapper, initialToAddrCoins, msg),
			},
		}
		return fOps, nil
	}

	return nil, nil
}
func checkLockSend(ak auth.AccountKeeper, oldAmt sdk.Coins, msg bankx.MsgSend) simulation.Operation {

	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (OperationMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {
		updatedAmt := ak.GetAccount(ctx, msg.ToAddress).GetCoins()
		if updatedAmt.Sub(msg.Amount).IsAllGTE(oldAmt) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}
		return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("lock send has failed")
	}
}

// SingleInputSendMsg tests and runs a single msg multisend, with one input and one output, where both
// accounts already exist.
func SimulateSingleInputMsgMultiSend(mapper auth.AccountKeeper, bk bankx.Keeper) simulation.Operation {
	handler := bankx.NewHandler(bk)
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		fromAcc, comment, msg, ok := createSingleInputMsgMultiSend(r, ctx, accs, mapper)
		opMsg = simulation.NewOperationMsg(msg, ok, comment)
		if !ok {
			return opMsg, nil, nil
		}
		err = sendAndVerifyMsgMultiSend(app, mapper, msg, ctx, []crypto.PrivKey{fromAcc.PrivKey}, handler)
		if err != nil {
			return opMsg, nil, err
		}
		return opMsg, nil, nil
	}
}

func createSingleInputMsgMultiSend(r *rand.Rand, ctx sdk.Context, accs []simulation.Account, mapper auth.AccountKeeper) (
	fromAcc simulation.Account, comment string, msg types.MsgMultiSend, ok bool) {

	fromAcc = simulation.RandomAcc(r, accs)
	toAcc := simulation.RandomAcc(r, accs)
	// Disallow sending money to yourself
	for {
		if !fromAcc.PubKey.Equals(toAcc.PubKey) {
			break
		}
		toAcc = simulation.RandomAcc(r, accs)
	}
	toAddr := toAcc.Address
	initFromCoins := mapper.GetAccount(ctx, fromAcc.Address).SpendableCoins(ctx.BlockHeader().Time)

	if len(initFromCoins) == 0 {
		return fromAcc, "skipping, no coins at all", msg, false
	}

	denomIndex := r.Intn(len(initFromCoins))
	amt, goErr := simulation.RandPositiveInt(r, initFromCoins[denomIndex].Amount)
	if goErr != nil {
		return fromAcc, "skipping bank send due to account having no coins of denomination " + initFromCoins[denomIndex].Denom, msg, false
	}

	coins := sdk.Coins{sdk.NewCoin(initFromCoins[denomIndex].Denom, amt)}
	msg = types.MsgMultiSend{
		Inputs:  []bank.Input{bank.NewInput(fromAcc.Address, coins)},
		Outputs: []bank.Output{bank.NewOutput(toAddr, coins)},
	}
	return fromAcc, "", msg, true
}

// Sends and verifies the transition of a msg multisend. This fails if there are repeated inputs or outputs
// pass in handler as nil to handle txs, otherwise handle msgs
func sendAndVerifyMsgMultiSend(app *baseapp.BaseApp, mapper auth.AccountKeeper, msg types.MsgMultiSend,
	ctx sdk.Context, privkeys []crypto.PrivKey, handler sdk.Handler) error {

	initialInputAddrCoins := make([]sdk.Coins, len(msg.Inputs))
	initialOutputAddrCoins := make([]sdk.Coins, len(msg.Outputs))
	AccountNumbers := make([]uint64, len(msg.Inputs))
	SequenceNumbers := make([]uint64, len(msg.Inputs))

	for i := 0; i < len(msg.Inputs); i++ {
		acc := mapper.GetAccount(ctx, msg.Inputs[i].Address)
		AccountNumbers[i] = acc.GetAccountNumber()
		SequenceNumbers[i] = acc.GetSequence()
		initialInputAddrCoins[i] = acc.GetCoins()
	}
	for i := 0; i < len(msg.Outputs); i++ {
		acc := mapper.GetAccount(ctx, msg.Outputs[i].Address)
		initialOutputAddrCoins[i] = acc.GetCoins()
	}
	if handler != nil {
		res := handler(ctx, msg)
		if !res.IsOK() {
			if res.Code == bank.CodeSendDisabled || res.Code == types.CodeTokenForbiddenByOwner {
				return nil
			}
			// TODO: Do this in a more 'canonical' way
			return fmt.Errorf("handling msg failed %v", res)
		}
	} else {
		tx := mock.GenTx([]sdk.Msg{msg},
			AccountNumbers,
			SequenceNumbers,
			privkeys...)
		res := app.Deliver(tx)
		if !res.IsOK() {
			// TODO: Do this in a more 'canonical' way
			return fmt.Errorf("Deliver failed %v", res)
		}
	}

	for i := 0; i < len(msg.Inputs); i++ {
		terminalInputCoins := mapper.GetAccount(ctx, msg.Inputs[i].Address).GetCoins()
		if !initialInputAddrCoins[i].Sub(msg.Inputs[i].Coins).IsEqual(terminalInputCoins) {
			return fmt.Errorf("input #%d had an incorrect amount of coins", i)
		}
	}
	for i := 0; i < len(msg.Outputs); i++ {
		terminalOutputCoins := mapper.GetAccount(ctx, msg.Outputs[i].Address).GetCoins()
		if !terminalOutputCoins.IsEqual(initialOutputAddrCoins[i].Add(msg.Outputs[i].Coins)) {
			return fmt.Errorf("output #%d had an incorrect amount of coins", i)
		}
	}
	return nil
}

func SimulateMsgSetMemoRequired(k bankx.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simulation.Account) (
		opMsg simulation.OperationMsg, futureOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accounts)

		msg := bankx.NewMsgSetTransferMemoRequired(acc.Address, r.Intn(2) > 0)
		ok := dexsim.SimulateHandleMsg(msg, bankx.NewHandler(k), ctx)
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
