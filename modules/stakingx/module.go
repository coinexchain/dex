package stakingx

import (
	"encoding/json"

	"github.com/coinexchain/dex/types"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	staking_types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct {
	apc types.ModuleClient
}

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
	amb.apc.RegisterRESTRoutes(ctx, rtr)
}

// get the root tx command of this module
func (amb AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return amb.apc.GetTxCmd(cdc)
}

// get the root query command of this module
func (amb AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return amb.apc.GetQueryCmd(cdc)
}

//___________________________
// app module object
type AppModule struct {
	AppModuleBasic
	stakingXKeeper Keeper //TODO: rename to StakingXKeeper
	apc            types.ModuleClient
}

// NewAppModule creates a new AppModule object
func NewAppModule(stakingxKeeper Keeper, apc types.ModuleClient) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{apc: apc},
		stakingXKeeper: stakingxKeeper,
		apc:            apc,
	}
}

// module name
func (AppModule) Name() string {
	return ModuleName
}

// register invariants
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

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
