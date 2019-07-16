package asset

import (
	"github.com/coinexchain/dex/modules/asset/keeper"
	"github.com/coinexchain/dex/modules/asset/types"
)

const (
	DefaultParamspace = types.DefaultParamspace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
)

var (
	// functions aliases
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState
	NewGenesisState     = types.NewGenesisState

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper = keeper.BaseKeeper
	TokenKeeper = keeper.TokenKeeper
	Params = types.Params
	GenesisState = types.GenesisState
	Token = types.Token
)
