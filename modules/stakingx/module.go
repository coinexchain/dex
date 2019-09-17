package stakingx

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	staking_cli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	staking_types "github.com/cosmos/cosmos-sdk/x/staking/types"

	stakingx_cli "github.com/coinexchain/dex/modules/stakingx/client/cli"
	stakingx_rest "github.com/coinexchain/dex/modules/stakingx/client/rest"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return staking_types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := staking_types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}

	return data.ValidateGenesis()
}

// register rest routes
func (amb AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	stakingx_rest.RegisterRoutes(ctx, rtr)
}

// get the root tx command of this module
func (amb AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return staking_cli.GetTxCmd(staking_types.StoreKey, cdc)
}

// get the root query command of this module
func (amb AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return stakingx_cli.GetQueryCmd(cdc)
}

//___________________________
// app module object
type AppModule struct {
	AppModuleBasic
	stakingXKeeper Keeper //TODO: rename to StakingXKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(stakingxKeeper Keeper) AppModule {
	return AppModule{
		stakingXKeeper: stakingxKeeper,
	}
}

// module name
func (AppModule) Name() string {
	return ModuleName
}

// register invariants
func (am AppModule) RegisterInvariants(invReg sdk.InvariantRegistry) {}

// module message route name
func (AppModule) Route() string { return "" }

// module handler
func (AppModule) NewHandler() sdk.Handler { return nil }

// module querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.stakingXKeeper, staking_types.ModuleCdc)
}

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	staking_types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.stakingXKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.stakingXKeeper)
	return staking_types.ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
