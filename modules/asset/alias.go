package asset

import (
	"github.com/coinexchain/dex/modules/asset/keeper"
	"github.com/coinexchain/dex/modules/asset/types"
)

const (
	DefaultParamspace    = types.DefaultParamspace
	ModuleName           = types.ModuleName
	StoreKey             = types.StoreKey
	QuerierRoute         = types.QuerierRoute
	RouterKey            = types.RouterKey
	QueryToken           = types.QueryToken
	QueryTokenList       = types.QueryTokenList
	QueryWhitelist       = types.QueryWhitelist
	QueryForbiddenAddr   = types.QueryForbiddenAddr
	QueryReservedSymbols = types.QueryReservedSymbols
)

var (
	// functions aliases
	NewQuerier                  = keeper.NewQuerier
	NewBaseKeeper               = keeper.NewBaseKeeper
	NewBaseTokenKeeper          = keeper.NewBaseTokenKeeper
	RegisterCodec               = types.RegisterCodec
	DefaultGenesisState         = types.DefaultGenesisState
	NewGenesisState             = types.NewGenesisState
	NewQueryAssetParams         = types.NewQueryAssetParams
	NewMsgIssueToken            = types.NewMsgIssueToken
	NewMsgTransferOwnership     = types.NewMsgTransferOwnership

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
