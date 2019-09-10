package app

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	genesisFile        string
	paramsFile         string
	exportParamsPath   string
	exportParamsHeight int
	exportStatePath    string
	exportStatsPath    string
	seed               int64
	initialBlockHeight int
	numBlocks          int
	blockSize          int
	enabled            bool
	verbose            bool
	lean               bool
	commit             bool
	period             int
	onOperation        bool // TODO Remove in favor of binary search for invariant violation
	allInvariants      bool
	genesisTime        int64
)

// AppStateFromGenesisFileFn util function to generate the genesis AppState
// from a genesis.json file
func AppStateFromGenesisFileFn(
	r *rand.Rand, _ []simulation.Account, _ time.Time,
) (json.RawMessage, []simulation.Account, string) {

	var genesis tmtypes.GenesisDoc
	cdc := MakeCodec()

	bytes, err := ioutil.ReadFile(genesisFile)
	if err != nil {
		panic(err)
	}

	cdc.MustUnmarshalJSON(bytes, &genesis)

	var appState GenesisState
	cdc.MustUnmarshalJSON(genesis.AppState, &appState)

	accounts := genaccounts.GetGenesisStateFromAppState(cdc, appState.toMap(cdc))

	var newAccs []simulation.Account
	for _, acc := range accounts {
		// Pick a random private key, since we don't know the actual key
		// This should be fine as it's only used for mock Tendermint validators
		// and these keys are never actually used to sign by mock Tendermint.
		privkeySeed := make([]byte, 15)
		r.Read(privkeySeed)

		privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
		newAccs = append(newAccs, simulation.Account{PrivKey: privKey, PubKey: privKey.PubKey(), Address: acc.Address})
	}

	return genesis.AppState, newAccs, genesis.ChainID
}
