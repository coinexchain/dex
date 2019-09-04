package app

import (
	"encoding/json"
	"reflect"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sltypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type NewHeightInfo struct {
	Height        int64        `json:"height"`
	TimeStamp     time.Time    `json:"timestamp"`
	LastBlockHash cmn.HexBytes `json:"last_block_hash"`
}

func (app *CetChainApp) pushNewHeightInfo(ctx sdk.Context) {
	msg := NewHeightInfo{
		Height:        ctx.BlockHeight(),
		TimeStamp:     ctx.BlockHeader().Time,
		LastBlockHash: ctx.BlockHeader().LastBlockId.Hash,
	}
	bytes, errJSON := json.Marshal(msg)
	if errJSON != nil {
		bytes = []byte{}
	}
	PubMsgs = append(PubMsgs, PubMsg{Key: []byte("height_info"), Value: bytes})
}

type TransferRecord struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
}

type NotificationTx struct {
	Signers      []sdk.AccAddress `json:"signers"`
	Transfers    []TransferRecord `json:"transfers"`
	SerialNumber int64            `json:"serial_number"`
	MsgTypes     []string         `json:"msg_types"`
	TxJSON       string           `json:"tx_json"`
	Height       int64            `json:"height"`
	Hash         []byte           `json:"hash"`
}

func getTransferRecord(dualEvent []abci.Event) TransferRecord {
	var res TransferRecord
	for _, attr := range dualEvent[0].Attributes {
		if string(attr.Key) == "recipient" {
			res.Recipient = string(attr.Value)
		} else if string(attr.Key) == "amount" {
			res.Amount = string(attr.Value)
		}
	}
	for _, attr := range dualEvent[1].Attributes {
		if string(attr.Key) == sdk.AttributeKeySender {
			res.Sender = string(attr.Value)
		}
	}
	return res
}

func getType(myvar interface{}) string {
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}

func (app *CetChainApp) notifyTx(req abci.RequestDeliverTx, events []abci.Event) {
	transfers := make([]TransferRecord, 0, 10)
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
			//} else if events[i].Type == dtypes.EventTypeWithdrawRewards {
			//	if i+1 <= len(events) {
			//		val := getWithdrawRewardInfo(events[i : i+2])
			//		PubMsgs = append(PubMsgs, PubMsg{Key: []byte("withdraw_reward"), Value: val})
			//		i++
			//	}
		} else if events[i].Type == "transfer" && i+1 <= len(events) {
			val := getTransferRecord(events[i : i+2])
			transfers = append(transfers, val)
			i++
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

	bytes, errJSON := json.Marshal(&stdTx)
	if errJSON != nil {
		return
	}

	n4s := &NotificationTx{
		Signers:      stdTx.GetSigners(),
		Transfers:    transfers,
		SerialNumber: app.txCount,
		TxJSON:       string(bytes),
		MsgTypes:     msgTypes,
		Height:       app.height,
		Hash:         req.Tx.Hash(),
	}

	bytes, errJSON = json.Marshal(n4s)
	if errJSON != nil {
		return
	}

	PubMsgs = append(PubMsgs, PubMsg{Key: []byte("notify_tx"), Value: bytes})
}

//type WithdrawRewardInfo struct {
//	Delegator string `json:"delegator"`
//	Validator string `json:"validator"`
//	Amount    string `json:"amount"`
//}
//
//func getWithdrawRewardInfo(dualEvent []abci.Event) []byte {
//	var res WithdrawRewardInfo
//	for _, attr := range dualEvent[0].Attributes {
//		if string(attr.Key) == dtypes.AttributeKeyValidator {
//			res.Validator = string(attr.Value)
//		} else if string(attr.Key) == dtypes.AttributeKeyAmount {
//			res.Amount = string(attr.Value)
//		}
//	}
//	for _, attr := range dualEvent[1].Attributes {
//		if string(attr.Key) == sdk.AttributeKeySender {
//			res.Delegator = string(attr.Value)
//		}
//	}
//	bytes, errJSON := json.Marshal(res)
//	if errJSON != nil {
//		return []byte{}
//	}
//	return bytes
//}

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
		} else if string(attr.Key) == sdk.AttributeKeyAmount {
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
		} else if string(attr.Key) == sdk.AttributeKeyAmount {
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

type NotificationSlash struct {
	Validator string `json:"validator"`
	Power     string `json:"power"`
	Reason    string `json:"reason"`
	Jailed    bool   `json:"jailed"`
}

func getNotificationSlash(event abci.Event) []byte {
	var res NotificationSlash
	for _, attr := range event.Attributes {
		if string(attr.Key) == sltypes.AttributeKeyAddress {
			res.Validator = string(attr.Value)
		} else if string(attr.Key) == sltypes.AttributeKeyPower {
			res.Power = string(attr.Value)
		} else if string(attr.Key) == sltypes.AttributeKeyReason {
			res.Reason = string(attr.Value)
		} else if string(attr.Key) == sltypes.AttributeKeyJailed {
			res.Jailed = true
		}
	}
	bytes, errJSON := json.Marshal(res)
	if errJSON != nil {
		return []byte{}
	}
	return bytes
}

func (app *CetChainApp) notifyBeginBlock(events []abci.Event) {
	//fmt.Printf("========== BeginBlock events ============\n")
	for _, event := range events {
		//fmt.Printf("= Event: %s\n", event.Type)
		//for _, attr := range event.Attributes {
		//	fmt.Printf("= K: %s; V: %s\n", attr.Key, attr.Value)
		//}
		if event.Type == sltypes.EventTypeSlash {
			val := getNotificationSlash(event)
			PubMsgs = append(PubMsgs, PubMsg{Key: []byte("slash"), Value: val})
		}
	}
}

func (app *CetChainApp) notifyEndBlock(events []abci.Event) {
	//fmt.Printf("========== EndBlock events ============\n")
	for _, event := range events {
		//fmt.Printf("= Event: %s\n", event.Type)
		//for _, attr := range event.Attributes {
		//	fmt.Printf("= K: %s; V: %s\n", attr.Key, attr.Value)
		//}
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
//		sdk.NewEvent(
//			types.EventTypeTransfer,
//			sdk.NewAttribute(types.AttributeKeyRecipient, toAddr.String()),
//			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
//		),
//		sdk.NewEvent(
//			sdk.EventTypeMessage,
//			sdk.NewAttribute(types.AttributeKeySender, fromAddr.String()),
//		),

//					types.EventTypeSlash,
//					sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
//					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
//					sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
//					sdk.NewAttribute(types.AttributeKeyJailed, consAddr.String()),
