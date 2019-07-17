package asset

import (
	"github.com/coinexchain/dex/modules/asset/internal/keeper"
	"github.com/coinexchain/dex/modules/asset/internal/types"
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
	MaxTokenAmount       = types.MaxTokenAmount
	RareSymbolLength     = types.RareSymbolLength

	IssueTokenFee     = types.IssueTokenFee
	IssueRareTokenFee = types.IssueRareTokenFee
)

var (
	// functions aliases

	NewQuerier                   = keeper.NewQuerier
	NewBaseKeeper                = keeper.NewBaseKeeper
	NewBaseTokenKeeper           = keeper.NewBaseTokenKeeper
	RegisterCodec                = types.RegisterCodec
	DefaultGenesisState          = types.DefaultGenesisState
	NewGenesisState              = types.NewGenesisState
	NewQueryAssetParams          = types.NewQueryAssetParams
	NewToken                     = types.NewToken
	NewMsgIssueToken             = types.NewMsgIssueToken
	NewMsgTransferOwnership      = types.NewMsgTransferOwnership
	NewMsgMintToken              = types.NewMsgMintToken
	NewMsgBurnToken              = types.NewMsgBurnToken
	NewMsgForbidToken            = types.NewMsgForbidToken
	NewMsgUnForbidToken          = types.NewMsgUnForbidToken
	NewMsgAddTokenWhitelist      = types.NewMsgAddTokenWhitelist
	NewMsgRemoveTokenWhitelist   = types.NewMsgRemoveTokenWhitelist
	NewMsgForbidAddr             = types.NewMsgForbidAddr
	NewMsgUnForbidAddr           = types.NewMsgUnForbidAddr
	NewMsgModifyTokenURL         = types.NewMsgModifyTokenURL
	NewMsgModifyTokenDescription = types.NewMsgModifyTokenDescription

	DefaultParams = types.DefaultParams

	// variable aliases

	ModuleCdc = types.ModuleCdc
)

type (
	Keeper         = keeper.BaseKeeper
	TokenKeeper    = keeper.TokenKeeper
	Params         = types.Params
	GenesisState   = types.GenesisState
	Token          = types.Token
	BaseToken      = types.BaseToken
	MsgForbidToken = types.MsgForbidToken
	MsgForbidAddr  = types.MsgForbidAddr
)
