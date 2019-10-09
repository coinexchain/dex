package keepers

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

var (
	StateKey = []byte{0x01}
)

type Keeper struct {
	cdc              *codec.Codec
	key              sdk.StoreKey
	paramSubspace    params.Subspace
	bankKeeper       types.BankKeeper
	supplyKeeper     authtypes.SupplyKeeper
	feeCollectorName string
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSubspace params.Subspace,
	bk types.BankKeeper, supplyKeeper authtypes.SupplyKeeper, feeCollectorName string) Keeper {

	return Keeper{
		cdc:              cdc,
		key:              key,
		paramSubspace:    paramSubspace.WithKeyTable(types.ParamKeyTable()),
		bankKeeper:       bk,
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}

func (k Keeper) GetParams(ctx sdk.Context) (param types.Params) {
	k.paramSubspace.GetParamSet(ctx, &param)
	return
}
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

func (k Keeper) GetState(ctx sdk.Context) (state types.State) {
	store := ctx.KVStore(k.key)
	bz := store.Get(StateKey)
	if bz == nil {
		panic("cannot load the adjustment height for incentive pool")
	}
	if err := k.cdc.UnmarshalBinaryBare(bz, &state); err != nil {
		panic(err)
	}
	return
}

func (k Keeper) SetState(ctx sdk.Context, state types.State) sdk.Error {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalBinaryBare(state)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(StateKey, bz)
	return nil
}

func (k Keeper) AddNewPlan(ctx sdk.Context, plan types.Plan) sdk.Error {
	if err := types.CheckPlans([]types.Plan{plan}); err != nil {
		return sdk.NewError(types.CodeSpaceIncentive, types.CodeInvalidPlanToAdd, "new plan is invalid")
	}
	param := k.GetParams(ctx)
	param.Plans = append(param.Plans, plan)
	k.SetParams(ctx, param)
	return nil
}

func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}
func (k Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return k.bankKeeper.HasCoins(ctx, addr, amt)
}
