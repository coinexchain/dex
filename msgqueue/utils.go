package msgqueue

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FillMsgs(ctx sdk.Context, key string, msg interface{}) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(EventTypeMsgQueue, sdk.NewAttribute(key, string(bytes))))
}
