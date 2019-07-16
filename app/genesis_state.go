package app

import (
	"encoding/json"

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
	//Accounts     []GenesisAccount          `json:"accounts"`
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
		//Accounts:     nil,
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
	//accounts []GenesisAccount,
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
		//Accounts:     accounts,
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
