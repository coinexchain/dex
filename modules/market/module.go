package market

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/coinexchain/dex/modules/market/client/cli"
	"github.com/coinexchain/dex/modules/market/client/rest"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	types2 "github.com/coinexchain/dex/modules/market/internal/types"
	"github.com/coinexchain/dex/types"
)

// app module basics object
type AppModuleBasic struct {
}

func (AppModuleBasic) Name() string {
	return types2.ModuleName
}
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types2.RegisterCodec(cdc)
}

// genesis
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types2.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(data json.RawMessage) error {
	var state GenesisState
	if err := types2.ModuleCdc.UnmarshalJSON(data, &state); err != nil {
		return err
	}
	return state.Validate()
}

// client functionality
func (amb AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, types2.ModuleCdc)
}

func (amb AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

func (amb AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// ___________________________
// app module object
type AppModule struct {
	AppModuleBasic
	marketKeeper keepers.Keeper //TODO: rename to AssetKeeper
	apc          types.ModuleClient
}

// NewAppModule creates a new AppModule object
func NewAppModule(marketKeeper keepers.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		marketKeeper:   marketKeeper,
	}
}

// registers
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// routes
func (am AppModule) Route() string {
	return types2.RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.marketKeeper)
}

func (am AppModule) QuerierRoute() string {
	return types2.QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keepers.NewQuerier(am.marketKeeper, nil)
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
	types2.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.marketKeeper, genesisState)

	// TODO. will check the return value
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.marketKeeper)
	return types2.ModuleCdc.MustMarshalJSON(gs)
}
