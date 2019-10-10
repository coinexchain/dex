package types

import (
	"encoding/json"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Merge several byte slices into one
func ConcatKeys(keys ...[]byte) []byte {
	var totalLen int
	for _, key := range keys {
		totalLen += len(key)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, key := range keys {
		i += copy(tmp[i:], key)
	}
	return tmp
}

func ErrUnknownRequest(module string, msg sdk.Msg) sdk.Result {
	//errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
	errMsg := fmt.Sprintf("Unrecognized %s Msg type: %s", module, msg.Type())
	return sdk.ErrUnknownRequest(errMsg).Result()
}

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
