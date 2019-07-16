package bankx

import (
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

const (
	DefaultCodespace = types.CodeSpaceBankx

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

var (
	// functions aliases
	RegisterCodec                        = types.RegisterCodec
	ParamKeyTable                        = types.ParamKeyTable
	DefaultParams                        = types.DefaultParams
	NewKeeper                            = keeper.NewKeeper
	NewMsgSend                           = types.NewMsgSend
	NewMsgSetTransferMemoRequired        = types.NewMsgSetTransferMemoRequired
	ErrMemoMissing                       = types.ErrMemoMissing
	ErrorInsufficientCETForActivatingFee = types.ErrorInsufficientCETForActivatingFee
	// variable aliases
	ModuleCdc       = types.ModuleCdc
	CodeMemoMissing = types.CodeMemoMissing
)

type (
	Keeper             = keeper.Keeper
	MsgSend            = types.MsgSend
	MsgSetMemoRequired = types.MsgSetMemoRequired
)
