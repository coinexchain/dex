package bankx

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/modules/bankx/internal/types"
)

func GetModuleCdc() *codec.Codec {
	return types.ModuleCdc
}

const (
	DefaultCodespace = types.CodeSpaceBankx

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

const (
	Create                    = types.Create
	Return                    = types.Return
	EarlierUnlockBySender     = types.EarlierUnlockBySender
	EarlierUnlockBySupervisor = types.EarlierUnlockBySupervisor
)

var (
	// functions aliases

	RegisterCodec                      = types.RegisterCodec
	ParamKeyTable                      = types.ParamKeyTable
	DefaultParams                      = types.DefaultParams
	NewParams                          = types.NewParams
	NewKeeper                          = keeper.NewKeeper
	NewMsgSend                         = types.NewMsgSend
	NewMsgSetTransferMemoRequired      = types.NewMsgSetTransferMemoRequired
	NewMsgMultiSend                    = types.NewMsgMultiSend
	ErrMemoMissing                     = types.ErrMemoMissing
	ErrInsufficientCETForActivatingFee = types.ErrInsufficientCETForActivatingFee

	// variable aliases

	CodeMemoMissing = types.CodeMemoMissing
)

type (
	Keeper             = keeper.Keeper
	MsgSend            = types.MsgSend
	MsgSetMemoRequired = types.MsgSetMemoRequired
	MsgMultiSend       = types.MsgMultiSend
	MsgSupervisedSend  = types.MsgSupervisedSend
)
