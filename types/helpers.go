package types

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/types"
)

func ResponseFrom(err types.Error) abci.ResponseCheckTx {
	result := err.Result()
	return abci.ResponseCheckTx{
		Code:      uint32(result.Code),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted),
		GasUsed:   int64(result.GasUsed),
		Events:    result.Events.ToABCIEvents(),
	}
}
