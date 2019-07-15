package market

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/types"
)

// app module basics object
type AppModuleBasic struct {
	apc types.ModuleClient
}

func (AppModuleBasic) Name() string {
	return ModuleName
}
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// genesis
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(data json.RawMessage) error {
	var state GenesisState
	if err := ModuleCdc.UnmarshalJSON(data, &state); err != nil {
		return err
	}
	return state.Validate()
}

// client functionality
func (amb AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	amb.apc.RegisterRESTRoutes(ctx, rtr)
}

func (amb AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return amb.apc.GetTxCmd(cdc)
}

func (amb AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return amb.apc.GetQueryCmd(cdc)
}

// ___________________________
// app module object
type AppModule struct {
	AppModuleBasic
	marketKeeper Keeper //TODO: rename to AssetKeeper
	apc          types.ModuleClient
}

// NewAppModule creates a new AppModule object
func NewAppModule(marketKeeper Keeper, apc types.ModuleClient) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{apc: apc},
		marketKeeper:   marketKeeper,
		apc:            apc,
	}
}

// registers
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// routes
func (am AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.marketKeeper)
}

func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.marketKeeper, nil)
}

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.marketKeeper)
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.marketKeeper)
	// TODO. will check the return val
	return nil
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.marketKeeper, genesisState)

	// TODO. will check the return value
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.marketKeeper)
	return ModuleCdc.MustMarshalJSON(gs)
}
