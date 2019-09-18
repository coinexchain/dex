package types

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ResponseFrom(err sdk.Error) abci.ResponseCheckTx {
	result := err.Result()
	return abci.ResponseCheckTx{
		Codespace: string(result.Codespace),
		Code:      uint32(result.Code),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted),
		GasUsed:   int64(result.GasUsed),
		Events:    result.Events.ToABCIEvents(),
	}
}

func SafeJSONMarshal(msg interface{}) []byte {
	bytes, errJSON := json.Marshal(msg)
	if errJSON != nil {
		bytes = []byte{}
	}
	return bytes
}
