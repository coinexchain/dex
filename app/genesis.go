package app

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	gaia_app "github.com/cosmos/cosmos-sdk/cmd/gaia/app"

	"github.com/coinexchain/dex/x/asset"
	"github.com/coinexchain/dex/x/bankx"
)

var (
	// bonded tokens given to genesis validators/accounts
	freeTokensPerAcc = sdk.TokensFromTendermintPower(150)
	defaultBondDenom = "cet"
)

// State to Unmarshal
type GenesisState struct {
	Accounts     []gaia_app.GenesisAccount `json:"accounts"`
	AuthData     auth.GenesisState         `json:"auth"`
	BankData     bank.GenesisState         `json:"bank"`
	BankXData    bankx.GenesisState        `json:"bankx"`
	StakingData  staking.GenesisState      `json:"staking"`
	DistrData    distribution.GenesisState `json:"distr"`
	GovData      gov.GenesisState          `json:"gov"`
	CrisisData   crisis.GenesisState       `json:"crisis"`
	SlashingData slashing.GenesisState     `json:"slashing"`
	AssetData    asset.GenesisState        `json:"asset"`
	GenTxs       []json.RawMessage         `json:"gentxs"`
}

// NewDefaultGenesisState generates the default state for gaia.
func NewDefaultGenesisState() GenesisState {
	gs := GenesisState{
		Accounts:     nil,
		AuthData:     auth.DefaultGenesisState(),
		BankData:     bank.DefaultGenesisState(),
		BankXData:    bankx.DefaultGenesisState(),
		StakingData:  staking.DefaultGenesisState(),
		DistrData:    distribution.DefaultGenesisState(),
		GovData:      gov.DefaultGenesisState(),
		CrisisData:   crisis.DefaultGenesisState(),
		SlashingData: slashing.DefaultGenesisState(),
		AssetData:    asset.DefaultGenesisState(),
		GenTxs:       nil,
	}
	// TODO: create staking.GenesisState from scratch
	gs.StakingData.Params.BondDenom = defaultBondDenom
	return gs
}

// Sanitize sorts accounts and coin sets.
func (gs GenesisState) Sanitize() {
	sort.Slice(gs.Accounts, func(i, j int) bool {
		return gs.Accounts[i].AccountNumber < gs.Accounts[j].AccountNumber
	})

	for _, acc := range gs.Accounts {
		acc.Coins = acc.Coins.Sort()
	}
}

// ValidateGenesisState ensures that the genesis state obeys the expected invariants
// TODO: No validators are both bonded and jailed (#2088)
// TODO: Error if there is a duplicate validator (#1708)
// TODO: Ensure all state machine parameters are in genesis (#1704)
func ValidateGenesisState(genesisState GenesisState) error {
	if err := validateGenesisStateAccounts(genesisState.Accounts); err != nil {
		return err
	}

	// skip stakingData validation as genesis is created from txs
	if len(genesisState.GenTxs) > 0 {
		return nil
	}

	if err := auth.ValidateGenesis(genesisState.AuthData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(genesisState.BankData); err != nil {
		return err
	}
	if err := bankx.ValidateGenesis(genesisState.BankXData); err != nil {
		return err
	}
	if err := staking.ValidateGenesis(genesisState.StakingData); err != nil {
		return err
	}
	//if err := mint.ValidateGenesis(genesisState.MintData); err != nil {
	//	return err
	//}
	if err := distribution.ValidateGenesis(genesisState.DistrData); err != nil {
		return err
	}
	if err := gov.ValidateGenesis(genesisState.GovData); err != nil {
		return err
	}
	if err := crisis.ValidateGenesis(genesisState.CrisisData); err != nil {
		return err
	}
	if err := asset.ValidateGenesis(genesisState.AssetData); err != nil {
		return err
	}

	return slashing.ValidateGenesis(genesisState.SlashingData)
}

// validateGenesisStateAccounts performs validation of genesis accounts. It
// ensures that there are no duplicate accounts in the genesis state and any
// provided vesting accounts are valid.
func validateGenesisStateAccounts(accs []gaia_app.GenesisAccount) error {
	addrMap := make(map[string]bool, len(accs))
	for _, acc := range accs {
		addrStr := acc.Address.String()

		// disallow any duplicate accounts
		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}

		// validate any vesting fields
		if !acc.OriginalVesting.IsZero() {
			if acc.EndTime == 0 {
				return fmt.Errorf("missing end time for vesting account; address: %s", addrStr)
			}

			if acc.StartTime >= acc.EndTime {
				return fmt.Errorf(
					"vesting start time must before end time; address: %s, start: %s, end: %s",
					addrStr,
					time.Unix(acc.StartTime, 0).UTC().Format(time.RFC3339),
					time.Unix(acc.EndTime, 0).UTC().Format(time.RFC3339),
				)
			}
		}

		addrMap[addrStr] = true
	}

	return nil
}
