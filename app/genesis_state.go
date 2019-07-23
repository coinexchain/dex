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
	"github.com/coinexchain/dex/modules/bancorlite"
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
	BancorData   bancorlite.GenesisState   `json:"bancorlite"`
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
		BancorData:   bancorlite.DefaultGenesisState(),
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

	unmarshalField(cdc, g[genaccounts.ModuleName], &gs.Accounts)
	unmarshalField(cdc, g[auth.ModuleName], &gs.AuthData)
	unmarshalField(cdc, g[authx.ModuleName], &gs.AuthXData)
	unmarshalField(cdc, g[bank.ModuleName], &gs.BankData)
	unmarshalField(cdc, g[bankx.ModuleName], &gs.BankXData)
	unmarshalField(cdc, g[staking.ModuleName], &gs.StakingData)
	unmarshalField(cdc, g[stakingx.ModuleName], &gs.StakingXData)
	unmarshalField(cdc, g[distribution.ModuleName], &gs.DistrData)
	unmarshalField(cdc, g[gov.ModuleName], &gs.GovData)
	unmarshalField(cdc, g[crisis.ModuleName], &gs.CrisisData)
	unmarshalField(cdc, g[slashing.ModuleName], &gs.SlashingData)
	unmarshalField(cdc, g[asset.ModuleName], &gs.AssetData)
	unmarshalField(cdc, g[market.ModuleName], &gs.MarketData)
	unmarshalField(cdc, g[bancorlite.ModuleName], &gs.BancorData)
	unmarshalField(cdc, g[incentive.ModuleName], &gs.Incentive)
	unmarshalField(cdc, g[supply.ModuleName], &gs.Supply)
	unmarshalField(cdc, g[genutil.ModuleName], &gs.GenUtil)

	return gs
}

func unmarshalField(cdc *codec.Codec, bz []byte, ptr interface{}) {
	if bz != nil {
		cdc.MustUnmarshalJSON(bz, ptr)
	}
}
