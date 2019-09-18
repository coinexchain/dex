package types

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type OrderedBasicManager struct {
	module.BasicManager
	modules []module.AppModuleBasic
}

func NewOrderedBasicManager(modules []module.AppModuleBasic) OrderedBasicManager {
	return OrderedBasicManager{
		BasicManager: module.NewBasicManager(modules...),
		modules:      modules,
	}
}

func (bm OrderedBasicManager) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	for _, m := range bm.modules {
		m.RegisterRESTRoutes(ctx, rtr)
	}
}

func (bm OrderedBasicManager) AddTxCommands(rootTxCmd *cobra.Command, cdc *codec.Codec) {
	for _, m := range bm.modules {
		if cmd := m.GetTxCmd(cdc); cmd != nil {
			if !isDuplicatedTxCmd(m.Name()) {
				rootTxCmd.AddCommand(cmd)
			}
		}
	}
}

func isDuplicatedTxCmd(module string) bool {
	return module == "distrx" || //mounted (do not use distrx.ModuleName to prevent circle import)
		module == auth.ModuleName || //mounted
		module == bank.ModuleName || //overwritten
		module == staking.ModuleName //overwritten
}

func (bm OrderedBasicManager) AddQueryCommands(rootQueryCmd *cobra.Command, cdc *codec.Codec) {
	for _, m := range bm.modules {
		if cmd := m.GetQueryCmd(cdc); cmd != nil {
			if !isDuplicatedQueryCmd(m.Name()) {
				rootQueryCmd.AddCommand(cmd)
			}
		}
	}
}

func isDuplicatedQueryCmd(module string) bool {
	return module == auth.ModuleName || //overwritten
		module == staking.ModuleName //overwritten
}

func (bm OrderedBasicManager) ValidateGenesis(genesis map[string]json.RawMessage) error {
	for _, m := range bm.modules {
		if isEmptyDataForGenutil(genesis, m) {
			continue
		}

		if err := m.ValidateGenesis(genesis[m.Name()]); err != nil {
			return err
		}
	}
	return nil
}

func isEmptyDataForGenutil(genesis map[string]json.RawMessage, m module.AppModuleBasic) bool {
	return m.Name() == genutil.ModuleName && len(genesis[m.Name()]) == 0
}
