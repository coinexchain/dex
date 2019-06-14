package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/modules/stakingx"
	dex "github.com/coinexchain/dex/types"
)

// State to Unmarshal
type GenesisState struct {
	Accounts     []GenesisAccount          `json:"accounts"`
	AuthData     auth.GenesisState         `json:"auth"`
	AuthXData    authx.GenesisState        `json:"authx"`
	BankData     bank.GenesisState         `json:"bank"`
	BankXData    bankx.GenesisState        `json:"bankx"`
	StakingData  staking.GenesisState      `json:"staking"`
	StakingXData stakingx.GenesisState     `json:"stakingx"`
	DistrData    distribution.GenesisState `json:"distr"`
	GovData      gov.GenesisState          `json:"gov"`
	CrisisData   crisis.GenesisState       `json:"crisis"`
	SlashingData slashing.GenesisState     `json:"slashing"`
	AssetData    asset.GenesisState        `json:"asset"`
	MarketData   market.GenesisState       `json:"market"`
	GenTxs       []json.RawMessage         `json:"gentxs"`
	MsgQueData   msgqueue.GenesisState     `json:"msgqueue"`
}

// NewDefaultGenesisState generates the default state for coindex.
func NewDefaultGenesisState() GenesisState {
	gs := GenesisState{
		Accounts:     nil,
		AuthData:     auth.DefaultGenesisState(),
		AuthXData:    authx.DefaultGenesisState(),
		BankData:     bank.DefaultGenesisState(),
		BankXData:    bankx.DefaultGenesisState(),
		StakingData:  staking.DefaultGenesisState(),
		StakingXData: stakingx.DefaultGenesisState(),
		DistrData:    distribution.DefaultGenesisState(),
		GovData:      gov.DefaultGenesisState(),
		CrisisData:   crisis.DefaultGenesisState(),
		SlashingData: slashing.DefaultGenesisState(),
		AssetData:    asset.DefaultGenesisState(),
		MarketData:   market.DefaultGenesisState(),
		GenTxs:       nil,
	}
	// TODO: create staking.GenesisState & gov.GenesisState & crisis.GenesisState from scratch
	gs.StakingData.Params.UnbondingTime = stakingx.DefaultUnbondingTime
	gs.StakingData.Params.MaxValidators = stakingx.DefaultMaxValidators
	gs.StakingData.Params.BondDenom = dex.DefaultBondDenom
	gs.GovData.DepositParams.MinDeposit[0].Denom = dex.DefaultBondDenom
	gs.CrisisData.ConstantFee.Denom = dex.DefaultBondDenom
	return gs
}

func NewGenesisState(
	accounts []GenesisAccount,
	authData auth.GenesisState,
	authxData authx.GenesisState,
	bankData bank.GenesisState,
	bankxData bankx.GenesisState,
	stakingData staking.GenesisState,
	stakingxData stakingx.GenesisState,
	distrData distribution.GenesisState,
	govData gov.GenesisState,
	crisisData crisis.GenesisState,
	slashingData slashing.GenesisState,
	assetData asset.GenesisState,
	marketData market.GenesisState) GenesisState {

	return GenesisState{
		Accounts:     accounts,
		AuthData:     authData,
		AuthXData:    authxData,
		BankData:     bankData,
		BankXData:    bankxData,
		StakingData:  stakingData,
		StakingXData: stakingxData,
		DistrData:    distrData,
		GovData:      govData,
		CrisisData:   crisisData,
		SlashingData: slashingData,
		AssetData:    assetData,
		MarketData:   marketData,
	}
}

// Sanitize sorts accounts and coin sets.
func (gs GenesisState) Sanitize() {
	sort.Slice(gs.Accounts, func(i, j int) bool {
		return gs.Accounts[i].AccountNumber < gs.Accounts[j].AccountNumber
	})

	for _, acc := range gs.Accounts {
		acc.Coins = acc.Coins.Sort()
		acc.FrozenCoins = acc.FrozenCoins.Sort()
	}
}

// ValidateGenesisState ensures that the genesis state obeys the expected invariants
// TODO: No validators are both bonded and jailed (#2088)
// TODO: Error if there is a duplicate validator (#1708)
// TODO: Ensure all state machine parameters are in genesis (#1704)
func (gs GenesisState) Validate() error {
	if err := validateGenesisStateAccounts(gs.Accounts); err != nil {
		return err
	}

	if err := gs.AuthXData.Validate(); err != nil {
		return err
	}
	if err := gs.BankXData.Validate(); err != nil {
		return err
	}
	if err := gs.StakingXData.Validate(); err != nil {
		return err
	}
	if err := gs.AssetData.Validate(); err != nil {
		return err
	}
	if err := gs.MarketData.Validate(); err != nil {
		return err
	}

	// skip stakingData validation as genesis is created from txs
	if len(gs.GenTxs) > 0 {
		return nil
	}

	if err := auth.ValidateGenesis(gs.AuthData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(gs.BankData); err != nil {
		return err
	}
	if err := staking.ValidateGenesis(gs.StakingData); err != nil {
		return err
	}
	if err := distribution.ValidateGenesis(gs.DistrData); err != nil {
		return err
	}
	if err := gov.ValidateGenesis(gs.GovData); err != nil {
		return err
	}
	if err := crisis.ValidateGenesis(gs.CrisisData); err != nil {
		return err
	}
	if err := slashing.ValidateGenesis(gs.SlashingData); err != nil {
		return err
	}

	return nil
}

// validateGenesisStateAccounts performs validation of genesis accounts. It
// ensures that there are no duplicate accounts in the genesis state and any
// provided vesting accounts are valid.
func validateGenesisStateAccounts(accs []GenesisAccount) error {
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

// CetAppGenState but with JSON
func CetAppGenStateJSON(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {
	// create the final app state
	genesisState, err := CetAppGenState(cdc, genDoc, appGenTxs)
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(cdc, genesisState)
}

// Create the core parameters for genesis initialization for gaia
// note that the pubkey input is this machines pubkey
func CetAppGenState(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	genesisState GenesisState, err error) {

	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return genesisState, errors.New("there must be at least one genesis tx")
	}

	stakingData := genesisState.StakingData
	for i, genTx := range appGenTxs {
		var tx auth.StdTx
		if err := cdc.UnmarshalJSON(genTx, &tx); err != nil {
			return genesisState, err
		}

		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return genesisState, errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}

		if _, ok := msgs[0].(staking.MsgCreateValidator); !ok {
			return genesisState, fmt.Errorf(
				"Genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}

	for _, acc := range genesisState.Accounts {
		for _, coin := range acc.Coins {
			if coin.Denom == genesisState.StakingData.Params.BondDenom {
				stakingData.Pool.NotBondedTokens = stakingData.Pool.NotBondedTokens.
					Add(coin.Amount) // increase the supply
			}
		}
	}

	genesisState.StakingData = stakingData
	genesisState.GenTxs = appGenTxs

	return genesisState, nil
}
