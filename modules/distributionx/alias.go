package distributionx

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/modules/distributionx/types"
)

func GetModuleCdc() *codec.Codec {
	return types.ModuleCdc
}

const (
	ModuleName = types.ModuleName
)

type (
	MsgDonateToCommunityPool = types.MsgDonateToCommunityPool
)
