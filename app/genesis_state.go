package app

import (
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/cosmos-sdk/x/genaccounts"
)

// State to Unmarshal
type GenesisState struct {
	Accounts     genaccounts.GenesisState  `json:"accounts"`
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
	Supply       supply.GenesisState       `json:"supply"`
	GenUtil      genutil.GenesisState      `json:"genutil"`
}

func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts:     genaccounts.GenesisState{},
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
		Supply:       supply.DefaultGenesisState(),
		GenUtil:      genutil.GenesisState{},
	}
}
