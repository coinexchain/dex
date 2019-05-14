package app

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"
)

// export the state of gaia for a genesis file
func (app *CetChainApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (
	appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// TODO
	panic("not implemented yet!")
}
