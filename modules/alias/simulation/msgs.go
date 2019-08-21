package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/coinexchain/dex/modules/alias"
	"github.com/coinexchain/dex/modules/alias/internal/keepers"
	"github.com/coinexchain/dex/modules/alias/internal/types"
	simulationx "github.com/coinexchain/dex/simulation"
)

// TODO

func SimulateMsgAliasUpdate(k keepers.Keeper) simulation.Operation {

	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		isAdd := simulationx.RandomBool(r)
		fromAcc, msg, err := createMsgAliasUpdate(r, ctx, k, accs, isAdd)
		if err != nil {
			return simulation.NoOpMsg(alias.ModuleName), nil, nil
		}

		ok := handleAlias(ctx, app, k, msg, fromAcc)
		opMsg = simulation.NewOperationMsg(msg, ok, "")

		if !ok {
			return opMsg, nil, nil
		}

		if err = verify(ctx, k, msg); err != nil {
			return opMsg, nil, err
		}

		return opMsg, nil, nil
	}
}

func createMsgAliasUpdate(r *rand.Rand, ctx sdk.Context, k keepers.Keeper, accs []simulation.Account, isAdd bool) (
	fromAcc simulation.Account, msg types.MsgAliasUpdate, err error) {

	fromAcc = simulation.RandomAcc(r, accs)
	if isAdd {
		alias := randomAlias(r)
		asDefault := simulationx.RandomBool(r)

		msg = types.MsgAliasUpdate{
			Owner:     fromAcc.Address,
			Alias:     alias,
			IsAdd:     isAdd,
			AsDefault: asDefault,
		}
	} else {
		aliasList := getAccountAlias(ctx, k, fromAcc.Address)
		if len(aliasList) == 0 {
			return fromAcc, msg, fmt.Errorf("no alias to remove")
		}
		alias := randomAliasFromList(r, aliasList)
		msg = types.MsgAliasUpdate{
			Owner: fromAcc.Address,
			Alias: alias,
			IsAdd: isAdd,
		}
	}

	if msg.ValidateBasic() != nil {
		return fromAcc, msg, fmt.Errorf("expect msg to pass validation check")
	}

	return fromAcc, msg, nil
}

func handleAlias(ctx sdk.Context, app *baseapp.BaseApp, k keepers.Keeper, msg types.MsgAliasUpdate, fromAcc simulation.Account) bool {

	handler := alias.NewHandler(k)
	cachectx, write := ctx.CacheContext()
	ok := handler(cachectx, msg).IsOK()

	if ok {
		write()
	}
	return ok
}

func verify(ctx sdk.Context, k keepers.Keeper, msg types.MsgAliasUpdate) error {

	addr, asDefault := k.AliasKeeper.GetAddressFromAlias(ctx, msg.Alias)

	if msg.IsAdd {
		if !msg.Owner.Equals(sdk.AccAddress(addr)) || asDefault != msg.AsDefault {
			return fmt.Errorf("alias added operation failed")
		}
	} else {
		if addr != nil {
			return fmt.Errorf("alias remove operation failed")
		}
	}

	return nil
}
func getAccountAlias(ctx sdk.Context, k keepers.Keeper, addr sdk.AccAddress) []string {
	return k.AliasKeeper.GetAliasListOfAccount(ctx, addr)
}
func randomAlias(r *rand.Rand) string {
	aliasLength := r.Intn(aliasMaxLength-1) + aliasMinLength
	return randomSymbol(r, aliasLength)
}
func randomAliasFromList(r *rand.Rand, aliasList []string) string {
	return aliasList[r.Intn(len(aliasList))]
}
