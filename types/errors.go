package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ErrUnknownRequest(module string, msg sdk.Msg) sdk.Result {
	//errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
	errMsg := fmt.Sprintf("Unrecognized %s Msg type: %s", module, msg.Type())
	return sdk.ErrUnknownRequest(errMsg).Result()
}
