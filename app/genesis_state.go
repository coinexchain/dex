package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
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

	"github.com/coinexchain/cet-sdk/modules/alias"
	"github.com/coinexchain/cet-sdk/modules/asset"
	"github.com/coinexchain/cet-sdk/modules/authx"
	"github.com/coinexchain/cet-sdk/modules/bancorlite"
	"github.com/coinexchain/cet-sdk/modules/bankx"
	"github.com/coinexchain/cet-sdk/modules/comment"
	"github.com/coinexchain/cet-sdk/modules/incentive"
	"github.com/coinexchain/cet-sdk/modules/market"
	"github.com/coinexchain/cet-sdk/modules/stakingx"
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
	CommentData  comment.GenesisState      `json:"comment"`
	AliasData    alias.GenesisState        `json:"alias"`
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
		CommentData:  comment.DefaultGenesisState(),
		AliasData:    alias.DefaultGenesisState(),
		Incentive:    incentive.DefaultGenesisState(),
		Supply:       supply.DefaultGenesisState(),
		GenUtil:      genutil.GenesisState{},
	}
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
	unmarshalField(cdc, g[comment.ModuleName], &gs.CommentData)
	unmarshalField(cdc, g[alias.ModuleName], &gs.AliasData)
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

func (gs GenesisState) toMap(cdc *codec.Codec) map[string]json.RawMessage {
	m := make(map[string]json.RawMessage)
	m[genaccounts.ModuleName] = cdc.MustMarshalJSON(gs.Accounts)
	m[auth.ModuleName] = cdc.MustMarshalJSON(gs.AuthData)
	m[authx.ModuleName] = cdc.MustMarshalJSON(gs.AuthXData)
	m[bank.ModuleName] = cdc.MustMarshalJSON(gs.BankData)
	m[bankx.ModuleName] = cdc.MustMarshalJSON(gs.BankXData)
	m[staking.ModuleName] = cdc.MustMarshalJSON(gs.StakingData)
	m[stakingx.ModuleName] = cdc.MustMarshalJSON(gs.StakingXData)
	m[distribution.ModuleName] = cdc.MustMarshalJSON(gs.DistrData)
	m[gov.ModuleName] = cdc.MustMarshalJSON(gs.GovData)
	m[crisis.ModuleName] = cdc.MustMarshalJSON(gs.CrisisData)
	m[slashing.ModuleName] = cdc.MustMarshalJSON(gs.SlashingData)
	m[asset.ModuleName] = cdc.MustMarshalJSON(gs.AssetData)
	m[market.ModuleName] = cdc.MustMarshalJSON(gs.MarketData)
	m[bancorlite.ModuleName] = cdc.MustMarshalJSON(gs.BancorData)
	m[comment.ModuleName] = cdc.MustMarshalJSON(gs.CommentData)
	m[alias.ModuleName] = cdc.MustMarshalJSON(gs.AliasData)
	m[incentive.ModuleName] = cdc.MustMarshalJSON(gs.Incentive)
	m[supply.ModuleName] = cdc.MustMarshalJSON(gs.Supply)
	m[genutil.ModuleName] = cdc.MustMarshalJSON(gs.GenUtil)
	return m
}
