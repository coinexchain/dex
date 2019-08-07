package app

import (
	"encoding/json"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	dtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type NotificationForSigners struct {
	Signers      []sdk.AccAddress `json:"signers"`
	Recipients   []string         `json:"recipients"`
	SerialNumber int64            `json:"serial_number"`
	MsgTypes     []string         `json:"msg_types"`
	Tx           *auth.StdTx      `json:"tx"`
}

func getType(myvar interface{}) string {
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}

func (app *CetChainApp) notifySigners(req abci.RequestDeliverTx, events []abci.Event) {
	recipients := make([]string, 0, 10)
	for _, event := range events {
		for _, attr := range event.Attributes {
			if string(attr.Key) == "recipient" {
				recipients = append(recipients, string(attr.Value))
			}
		}
	}

	defer func() {
		app.txCount++
	}()

	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return
	}

	stdTx, ok := tx.(auth.StdTx)
	if !ok {
		return
	}
	msgTypes := make([]string, len(stdTx.Msgs))
	for i, msg := range stdTx.Msgs {
		msgTypes[i] = getType(msg)
	}

	n4s := &NotificationForSigners{
		Signers:      stdTx.GetSigners(),
		Recipients:   recipients,
		SerialNumber: app.txCount,
		Tx:           &stdTx,
		MsgTypes:     msgTypes,
	}

	bytes, errJSON := json.Marshal(n4s)
	if errJSON != nil {
		return
	}

	PubMsgs = append(PubMsgs, PubMsg{Key: []byte("notify_signers"), Value: bytes})
}

type WithdrawRewardInfo struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Amount    string `json:"amount"`
}

func getWithdrawRewardInfo(dualEvent []abci.Event) []byte {
	var res WithdrawRewardInfo
	for _, attr := range dualEvent[0].Attributes {
		if string(attr.Key) == dtypes.AttributeKeyValidator {
			res.Validator = string(attr.Value)
		} else if string(attr.Key) == dtypes.AttributeKeyAmount {
			res.Amount = string(attr.Value)
		}
	}
	for _, attr := range dualEvent[1].Attributes {
		if string(attr.Key) == sdk.AttributeKeySender {
			res.Delegator = string(attr.Value)
		}
	}
	bytes, errJSON := json.Marshal(res)
	if errJSON != nil {
		return []byte{}
	}
	return bytes
}

type NotificationBeginRedelegation struct {
	Delegator      string `json:"delegator"`
	ValidatorSrc   string `json:"src"`
	ValidatorDst   string `json:"dst"`
	Amount         string `json:"amount"`
	CompletionTime string `json:"completion_time"`
}

func getNotificationBeginRedelegation(dualEvent []abci.Event) []byte {
	var res NotificationBeginRedelegation
	for _, attr := range dualEvent[0].Attributes {
		if string(attr.Key) == stypes.AttributeKeySrcValidator {
			res.ValidatorSrc = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyDstValidator {
			res.ValidatorDst = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyAmount {
			res.Amount = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyCompletionTime {
			res.CompletionTime = string(attr.Value)
		}
	}
	for _, attr := range dualEvent[1].Attributes {
		if string(attr.Key) == sdk.AttributeKeySender {
			res.Delegator = string(attr.Value)
		}
	}
	bytes, errJSON := json.Marshal(res)
	if errJSON != nil {
		return []byte{}
	}
	return bytes
}

type NotificationBeginUnbonding struct {
	Delegator      string `json:"delegator"`
	Validator      string `json:"validator"`
	Amount         string `json:"amount"`
	CompletionTime string `json:"completion_time"`
}

func getNotificationBeginUnbonding(dualEvent []abci.Event) []byte {
	var res NotificationBeginUnbonding
	for _, attr := range dualEvent[0].Attributes {
		if string(attr.Key) == stypes.AttributeKeyValidator {
			res.Validator = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyAmount {
			res.Amount = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyCompletionTime {
			res.CompletionTime = string(attr.Value)
		}
	}
	for _, attr := range dualEvent[1].Attributes {
		if string(attr.Key) == sdk.AttributeKeySender {
			res.Delegator = string(attr.Value)
		}
	}
	bytes, errJSON := json.Marshal(res)
	if errJSON != nil {
		return []byte{}
	}
	return bytes
}

type NotificationCompleteRedelegation struct {
	Delegator    string `json:"delegator"`
	ValidatorSrc string `json:"src"`
	ValidatorDst string `json:"dst"`
}

func getNotificationCompleteRedelegation(event abci.Event) []byte {
	var res NotificationCompleteRedelegation
	for _, attr := range event.Attributes {
		if string(attr.Key) == stypes.AttributeKeyDstValidator {
			res.ValidatorDst = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeySrcValidator {
			res.ValidatorSrc = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyDelegator {
			res.Delegator = string(attr.Value)
		}
	}
	bytes, errJSON := json.Marshal(res)
	if errJSON != nil {
		return []byte{}
	}
	return bytes
}

type NotificationCompleteUnbonding struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
}

func getNotificationCompleteUnbonding(event abci.Event) []byte {
	var res NotificationCompleteUnbonding
	for _, attr := range event.Attributes {
		if string(attr.Key) == stypes.AttributeKeyValidator {
			res.Validator = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyDelegator {
			res.Delegator = string(attr.Value)
		}
	}
	bytes, errJSON := json.Marshal(res)
	if errJSON != nil {
		return []byte{}
	}
	return bytes
}

func (app *CetChainApp) notifyInTx(events []abci.Event) {
	for i := 0; i < len(events); i++ {
		if events[i].Type == stypes.EventTypeUnbond {
			if i+1 <= len(events) {
				val := getNotificationBeginUnbonding(events[i : i+2])
				PubMsgs = append(PubMsgs, PubMsg{Key: []byte("begin_unbonding"), Value: val})
				i++
			}
		} else if events[i].Type == stypes.EventTypeRedelegate {
			if i+1 <= len(events) {
				val := getNotificationBeginRedelegation(events[i : i+2])
				PubMsgs = append(PubMsgs, PubMsg{Key: []byte("begin_redelegation"), Value: val})
				i++
			}
		} else if events[i].Type == dtypes.EventTypeWithdrawRewards {
			if i+1 <= len(events) {
				val := getWithdrawRewardInfo(events[i : i+2])
				PubMsgs = append(PubMsgs, PubMsg{Key: []byte("withdraw_reward"), Value: val})
				i++
			}
		}
	}
}

func (app *CetChainApp) notifyComplete(events []abci.Event) {
	for _, event := range events {
		if event.Type == stypes.EventTypeCompleteUnbonding {
			val := getNotificationCompleteUnbonding(event)
			PubMsgs = append(PubMsgs, PubMsg{Key: []byte("complete_unbonding"), Value: val})
		} else if event.Type == stypes.EventTypeCompleteRedelegation {
			val := getNotificationCompleteRedelegation(event)
			PubMsgs = append(PubMsgs, PubMsg{Key: []byte("complete_redelegation"), Value: val})
		}
	}
}

//		sdk.NewEvent(
//			types.EventTypeWithdrawRewards,
//			sdk.NewAttribute(types.AttributeKeyAmount, rewards.String()),
//			sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
//		),
//		sdk.NewEvent(
//			sdk.EventTypeMessage,
//			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
//			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
//		),
//		sdk.NewEvent(
//			stypes.EventTypeUnbond,
//			sdk.NewAttribute(stypes.AttributeKeyValidator, msg.ValidatorAddress.String()),
//			sdk.NewAttribute(stypes.AttributeKeyAmount, msg.Amount.Amount.String()),
//			sdk.NewAttribute(stypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
//		),
//		sdk.NewEvent(
//			sdk.EventTypeMessage,
//			sdk.NewAttribute(sdk.AttributeKeyModule, stypes.AttributeValueCategory),
//			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
//		),
//		sdk.NewEvent(
//			stypes.EventTypeRedelegate,
//			sdk.NewAttribute(stypes.AttributeKeySrcValidator, msg.ValidatorSrcAddress.String()),
//			sdk.NewAttribute(stypes.AttributeKeyDstValidator, msg.ValidatorDstAddress.String()),
//			sdk.NewAttribute(stypes.AttributeKeyAmount, msg.Amount.Amount.String()),
//			sdk.NewAttribute(stypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
//		),
//		sdk.NewEvent(
//			sdk.EventTypeMessage,
//			sdk.NewAttribute(sdk.AttributeKeyModule, stypes.AttributeValueCategory),
//			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
//		),
//			sdk.NewEvent(
//				stypes.EventTypeCompleteUnbonding,
//				sdk.NewAttribute(stypes.AttributeKeyValidator, dvPair.ValidatorAddress.String()),
//				sdk.NewAttribute(stypes.AttributeKeyDelegator, dvPair.DelegatorAddress.String()),
//			),
//			sdk.NewEvent(
//				stypes.EventTypeCompleteRedelegation,
//				sdk.NewAttribute(stypes.AttributeKeyDelegator, dvvTriplet.DelegatorAddress.String()),
//				sdk.NewAttribute(stypes.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress.String()),
//				sdk.NewAttribute(stypes.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress.String()),
//			),
