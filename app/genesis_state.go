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
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	DistrData    distribution.GenesisState `json:"distribution"`
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

func (app *CetChainApp) ExportGenesisState(ctx sdk.Context) GenesisState {
	g := app.mm.ExportGenesis(ctx)
	gs := GenesisState{}

	app.cdc.MustUnmarshalJSON(g[genaccounts.ModuleName], &gs.Accounts)
	app.cdc.MustUnmarshalJSON(g[auth.ModuleName], &gs.AuthData)
	app.cdc.MustUnmarshalJSON(g[authx.ModuleName], &gs.AuthXData)
	app.cdc.MustUnmarshalJSON(g[bank.ModuleName], &gs.BankData)
	app.cdc.MustUnmarshalJSON(g[bankx.ModuleName], &gs.BankXData)
	app.cdc.MustUnmarshalJSON(g[staking.ModuleName], &gs.StakingData)
	app.cdc.MustUnmarshalJSON(g[stakingx.ModuleName], &gs.StakingXData)
	app.cdc.MustUnmarshalJSON(g[distribution.ModuleName], &gs.DistrData)
	app.cdc.MustUnmarshalJSON(g[gov.ModuleName], &gs.GovData)
	app.cdc.MustUnmarshalJSON(g[crisis.ModuleName], &gs.CrisisData)
	app.cdc.MustUnmarshalJSON(g[slashing.ModuleName], &gs.SlashingData)
	app.cdc.MustUnmarshalJSON(g[asset.ModuleName], &gs.AssetData)
	app.cdc.MustUnmarshalJSON(g[market.ModuleName], &gs.MarketData)
	app.cdc.MustUnmarshalJSON(g[incentive.ModuleName], &gs.Incentive)
	app.cdc.MustUnmarshalJSON(g[supply.ModuleName], &gs.Supply)
	app.cdc.MustUnmarshalJSON(g[genutil.ModuleName], &gs.GenUtil)

	return gs
}
