package incentive

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

const (
	RouterKey = "incentive"
	StoreKey  = RouterKey
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

func (k Keeper) GetParam(ctx sdk.Context) (param Params) {
	k.paramSubspace.GetParamSet(ctx, &param)
	return
}
func (k Keeper) SetParam(ctx sdk.Context, params Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

func (k Keeper) GetState(ctx sdk.Context) (state State) {

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

func (k Keeper) SetState(ctx sdk.Context, state State) sdk.Error {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalBinaryBare(state)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	store.Set(StateKey, bz)
	return nil
}

func (k Keeper) AddNewPlan(ctx sdk.Context, plan Plan) {
	param := k.GetParam(ctx)
	param.Plans = append(param.Plans, plan)
	k.SetParam(ctx, param)
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSubspace params.Subspace,
	bk types.BankKeeper, supplyKeeper authtypes.SupplyKeeper, feeCollectorName string) Keeper {

	return Keeper{
		cdc:              cdc,
		key:              key,
		paramSubspace:    paramSubspace.WithKeyTable(ParamKeyTable()),
		bankKeeper:       bk,
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}
