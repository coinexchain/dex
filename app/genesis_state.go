package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/stakingx"
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
	return FromMap(app.cdc, g)
}

func FromMap(cdc *codec.Codec, g map[string]json.RawMessage) GenesisState {
	gs := GenesisState{}

	cdc.MustUnmarshalJSON(g[genaccounts.ModuleName], &gs.Accounts)
	cdc.MustUnmarshalJSON(g[auth.ModuleName], &gs.AuthData)
	cdc.MustUnmarshalJSON(g[authx.ModuleName], &gs.AuthXData)
	cdc.MustUnmarshalJSON(g[bank.ModuleName], &gs.BankData)
	cdc.MustUnmarshalJSON(g[bankx.ModuleName], &gs.BankXData)
	cdc.MustUnmarshalJSON(g[staking.ModuleName], &gs.StakingData)
	cdc.MustUnmarshalJSON(g[stakingx.ModuleName], &gs.StakingXData)
	cdc.MustUnmarshalJSON(g[distribution.ModuleName], &gs.DistrData)
	cdc.MustUnmarshalJSON(g[gov.ModuleName], &gs.GovData)
	cdc.MustUnmarshalJSON(g[crisis.ModuleName], &gs.CrisisData)
	cdc.MustUnmarshalJSON(g[slashing.ModuleName], &gs.SlashingData)
	cdc.MustUnmarshalJSON(g[asset.ModuleName], &gs.AssetData)
	cdc.MustUnmarshalJSON(g[market.ModuleName], &gs.MarketData)
	cdc.MustUnmarshalJSON(g[incentive.ModuleName], &gs.Incentive)
	cdc.MustUnmarshalJSON(g[supply.ModuleName], &gs.Supply)
	cdc.MustUnmarshalJSON(g[genutil.ModuleName], &gs.GenUtil)

	return gs
}
