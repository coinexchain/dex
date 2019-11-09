package market

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/modules/market/client/cli"
	"github.com/coinexchain/dex/modules/market/client/rest"
	"github.com/coinexchain/dex/modules/market/internal/keepers"
	"github.com/coinexchain/dex/modules/market/internal/types"
)

// app module basics object
type AppModuleBasic struct {
}

func (AppModuleBasic) Name() string {
	return ModuleName
}
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// genesis
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(data json.RawMessage) error {
	var state GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(data, &state); err != nil {
		return err
	}
	return state.Validate()
}

// client functionality
func (amb AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, types.ModuleCdc)
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
	marketKeeper keepers.Keeper
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
	return types.RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.marketKeeper)
}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keepers.NewQuerier(am.marketKeeper)
}

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.marketKeeper)
	return nil
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.marketKeeper, genesisState)

	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.marketKeeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}
