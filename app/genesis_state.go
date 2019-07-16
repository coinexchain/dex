package app

import (
	"encoding/json"
)

// State to Unmarshal
type OrderedGenesisState struct {
	Accounts     json.RawMessage `json:"accounts"`
	AuthData     json.RawMessage `json:"auth"`
	AuthXData    json.RawMessage `json:"authx"`
	BankData     json.RawMessage `json:"bank"`
	BankXData    json.RawMessage `json:"bankx"`
	StakingData  json.RawMessage `json:"staking"`
	StakingXData json.RawMessage `json:"stakingx"`
	DistrData    json.RawMessage `json:"distr"`
	GovData      json.RawMessage `json:"gov"`
	CrisisData   json.RawMessage `json:"crisis"`
	SlashingData json.RawMessage `json:"slashing"`
	AssetData    json.RawMessage `json:"asset"`
	MarketData   json.RawMessage `json:"market"`
	Incentive    json.RawMessage `json:"incentive"`
	GenUtil      json.RawMessage `json:"genutil"`
}

func NewOrderedGenesisState(unordered map[string]json.RawMessage) OrderedGenesisState {
	return OrderedGenesisState{
		Accounts     : getAndDelete(unordered, "accounts"),
		AuthData     : getAndDelete(unordered, "auth"),
		AuthXData    : getAndDelete(unordered, "authx"),
		BankData     : getAndDelete(unordered, "bank"),
		BankXData    : getAndDelete(unordered, "bankx"),
		StakingData  : getAndDelete(unordered, "staking"),
		StakingXData : getAndDelete(unordered, "stakingx"),
		DistrData    : getAndDelete(unordered, "distr"),
		GovData      : getAndDelete(unordered, "gov"),
		CrisisData   : getAndDelete(unordered, "crisis"),
		SlashingData : getAndDelete(unordered, "slashing"),
		AssetData    : getAndDelete(unordered, "asset"),
		MarketData   : getAndDelete(unordered, "market"),
		Incentive    : getAndDelete(unordered, "incentive"),
		GenUtil      : getAndDelete(unordered, "genutil"),
	}
}

func getAndDelete(m map[string]json.RawMessage, key string) json.RawMessage {
	if val, ok := m[key]; ok {
		delete(m, key)
		return val
	}
	panic("key not exist: " + key)
}

//func adjustDefaultParams(gs *GenesisState) {
//	gs.AuthData.Params.MaxMemoCharacters = DefaultMaxMemoCharacters

//	gs.StakingData.Params.UnbondingTime = DefaultUnbondingTime
//	gs.StakingData.Params.MaxValidators = DefaultMaxValidators
//	gs.StakingData.Params.BondDenom = dex.DefaultBondDenom

//	gs.SlashingData.Params.MaxEvidenceAge = DefaultMaxEvidenceAge
//	gs.SlashingData.Params.SignedBlocksWindow = DefaultSignedBlocksWindow
//	gs.SlashingData.Params.MinSignedPerWindow = DefaultMinSignedPerWindow
//	gs.SlashingData.Params.SlashFractionDoubleSign = DefaultSlashFractionDoubleSign
//	gs.SlashingData.Params.SlashFractionDowntime = DefaultSlashFractionDowntime

//	gs.GovData.DepositParams.MinDeposit[0].Denom = dex.DefaultBondDenom
//	gs.GovData.DepositParams.MinDeposit[0].Amount = DefaultGovMinDeposit
//	gs.GovData.DepositParams.MaxDepositPeriod = DefaultPeriod
//	gs.GovData.VotingParams.VotingPeriod = DefaultPeriod
//	gs.GovData.TallyParams = gov.TallyParams{
//		Quorum:    sdk.NewDecWithPrec(4, 1),
//		Threshold: sdk.NewDecWithPrec(5, 1),
//		Veto:      sdk.NewDecWithPrec(334, 3),
//	}

//	gs.CrisisData.ConstantFee.Denom = dex.DefaultBondDenom
//	gs.CrisisData.ConstantFee.Amount = DefaultCrisisConstantFee
//}
