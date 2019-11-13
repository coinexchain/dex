package comment

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/comment/internal/keepers"
	"github.com/coinexchain/dex/modules/comment/internal/types"
)

func GetModuleCdc() *codec.Codec {
	return types.ModuleCdc
}

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	NewBaseKeeper = keepers.NewKeeper
)

type (
	Keeper          = keepers.Keeper
	TokenComment    = types.TokenComment
	CommentRef      = types.CommentRef
	MsgCommentToken = types.MsgCommentToken
)
