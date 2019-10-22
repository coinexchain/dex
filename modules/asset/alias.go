package asset

import (
	"github.com/coinexchain/dex/modules/asset/internal/keepers"
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
	QueryParameters      = types.QueryParameters
	QueryReservedSymbols = types.QueryReservedSymbols
	MaxTokenAmount       = types.MaxTokenAmount
	RareSymbolLength     = types.RareSymbolLength

	DefaultIssueTokenFee     = types.DefaultIssueTokenFee
	DefaultIssueRareTokenFee = types.DefaultIssueRareTokenFee
)

var (
	// functions aliases

	NewQuerier                 = keepers.NewQuerier
	NewBaseKeeper              = keepers.NewBaseKeeper
	NewBaseTokenKeeper         = keepers.NewBaseTokenKeeper
	RegisterCodec              = types.RegisterCodec
	DefaultGenesisState        = types.DefaultGenesisState
	NewGenesisState            = types.NewGenesisState
	NewQueryAssetParams        = types.NewQueryAssetParams
	NewToken                   = types.NewToken
	NewMsgIssueToken           = types.NewMsgIssueToken
	NewMsgTransferOwnership    = types.NewMsgTransferOwnership
	NewMsgMintToken            = types.NewMsgMintToken
	NewMsgBurnToken            = types.NewMsgBurnToken
	NewMsgForbidToken          = types.NewMsgForbidToken
	NewMsgUnForbidToken        = types.NewMsgUnForbidToken
	NewMsgAddTokenWhitelist    = types.NewMsgAddTokenWhitelist
	NewMsgRemoveTokenWhitelist = types.NewMsgRemoveTokenWhitelist
	NewMsgForbidAddr           = types.NewMsgForbidAddr
	NewMsgUnForbidAddr         = types.NewMsgUnForbidAddr
	NewMsgModifyTokenInfo      = types.NewMsgModifyTokenInfo
	TestIdentityString         = types.TestIdentityString
	ValidateTokenSymbol        = types.ValidateTokenSymbol

	DefaultParams = types.DefaultParams

	// variable aliases

	ModuleCdc = types.ModuleCdc
)

type (
	Keeper                  = keepers.BaseKeeper
	BaseTokenKeeper         = keepers.BaseTokenKeeper
	TokenKeeper             = keepers.TokenKeeper
	Params                  = types.Params
	GenesisState            = types.GenesisState
	Token                   = types.Token
	BaseToken               = types.BaseToken
	MsgForbidToken          = types.MsgForbidToken
	MsgForbidAddr           = types.MsgForbidAddr
	MsgIssueToken           = types.MsgIssueToken
	MsgTransferOwnership    = types.MsgTransferOwnership
	MsgMintToken            = types.MsgMintToken
	MsgBurnToken            = types.MsgBurnToken
	MsgUnForbidToken        = types.MsgUnForbidToken
	MsgAddTokenWhitelist    = types.MsgAddTokenWhitelist
	MsgRemoveTokenWhitelist = types.MsgRemoveTokenWhitelist
	MsgUnForbidAddr         = types.MsgUnForbidAddr
	MsgModifyTokenInfo      = types.MsgModifyTokenInfo
)
