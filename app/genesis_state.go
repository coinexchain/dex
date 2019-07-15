package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
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
	Incentive    incentive.GenesisState    `json:"incentive"`
	GenTxs       []json.RawMessage         `json:"gentxs"`
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
		Incentive:    incentive.DefaultGenesisState(),
		GenTxs:       nil,
	}
	// TODO: create staking.GenesisState & gov.GenesisState & crisis.GenesisState from scratch
	adjustDefaultParams(&gs)
	return gs
}

func adjustDefaultParams(gs *GenesisState) {
	gs.AuthData.Params.MaxMemoCharacters = DefaultMaxMemoCharacters
	gs.StakingData.Params.UnbondingTime = DefaultUnbondingTime
	gs.StakingData.Params.MaxValidators = DefaultMaxValidators
	gs.StakingData.Params.BondDenom = dex.DefaultBondDenom
	gs.SlashingData.Params.MaxEvidenceAge = DefaultMaxEvidenceAge
	gs.SlashingData.Params.SignedBlocksWindow = DefaultSignedBlocksWindow
	gs.SlashingData.Params.MinSignedPerWindow = DefaultMinSignedPerWindow
	gs.SlashingData.Params.SlashFractionDoubleSign = DefaultSlashFractionDoubleSign
	gs.SlashingData.Params.SlashFractionDowntime = DefaultSlashFractionDowntime
	gs.GovData.DepositParams.MinDeposit[0].Denom = dex.DefaultBondDenom
	gs.GovData.DepositParams.MinDeposit[0].Amount = DefaultGovMinDeposit
	gs.GovData.DepositParams.MaxDepositPeriod = DefaultPeriod
	gs.GovData.VotingParams.VotingPeriod = DefaultPeriod
	gs.GovData.TallyParams = gov.TallyParams{
		Quorum:    sdk.NewDecWithPrec(4, 1),
		Threshold: sdk.NewDecWithPrec(5, 1),
		Veto:      sdk.NewDecWithPrec(334, 3),
	}
	gs.CrisisData.ConstantFee.Denom = dex.DefaultBondDenom
	gs.CrisisData.ConstantFee.Amount = DefaultCrisisConstantFee
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
	marketData market.GenesisState,
	incentive incentive.GenesisState,
) GenesisState {

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
		Incentive:    incentive,
	}
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

// Create the core parameters for genesis initialization for CoinEx chain
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
				//stakingData.Pool.NotBondedTokens = stakingData.Pool.NotBondedTokens.
				//	Add(coin.Amount) // increase the supply
			}
		}
	}

	genesisState.StakingData = stakingData
	genesisState.GenTxs = appGenTxs

	return genesisState, nil
}
