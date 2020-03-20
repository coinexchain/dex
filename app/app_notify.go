package app

import (
	"encoding/json"
	"reflect"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	sltypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	dex "github.com/coinexchain/cet-sdk/types"
)

type TxExtraInfo struct {
	Code      uint32       `json:"code,omitempty"`
	Data      []byte       `json:"data,omitempty"`
	Log       string       `json:"log,omitempty"`
	Info      string       `json:"info,omitempty"`
	GasWanted int64        `json:"gas_wanted,omitempty"`
	GasUsed   int64        `json:"gas_used,omitempty"`
	Events    []abci.Event `json:"events,omitempty"`
	Codespace string       `json:"codespace,omitempty"`
}

type NewHeightInfo struct {
	ChainID       string       `json:"chain_id"`
	Height        int64        `json:"height"`
	TimeStamp     int64        `json:"timestamp"`
	LastBlockHash cmn.HexBytes `json:"last_block_hash"`
}

func (app *CetChainApp) pushNewHeightInfo(ctx sdk.Context) {
	msg := NewHeightInfo{
		ChainID:       ctx.BlockHeader().ChainID,
		Height:        ctx.BlockHeight(),
		TimeStamp:     ctx.BlockHeader().Time.Unix(),
		LastBlockHash: ctx.BlockHeader().LastBlockId.Hash,
	}
	bytes := dex.SafeJSONMarshal(msg)
	app.appendPubMsgKV("height_info", bytes)
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
	ExtraInfo    string           `json:"extra_info,omitempty"`
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

func (app *CetChainApp) notifyTx(req abci.RequestDeliverTx, stdTx auth.StdTx, ret abci.ResponseDeliverTx) {
	events := ret.Events
	transfers := make([]TransferRecord, 0, 10)
	ok := ret.Code == uint32(sdk.CodeOK)
	unbondingMsgList := make([][]byte, 0, 10)
	redelegationMsgList := make([][]byte, 0, 10)
	for i := 0; ok && i < len(events); i++ {
		if events[i].Type == stypes.EventTypeUnbond {
			if i+1 <= len(events) {
				val := getNotificationBeginUnbonding(events[i : i+2])
				unbondingMsgList = append(unbondingMsgList, val)
				i++
			}
		} else if events[i].Type == stypes.EventTypeRedelegate {
			if i+1 <= len(events) {
				val := getNotificationBeginRedelegation(events[i : i+2])
				redelegationMsgList = append(redelegationMsgList, val)
				i++
			}
		} else if events[i].Type == "transfer" && i+2 <= len(events) {
			val := getTransferRecord(events[i : i+2])
			transfers = append(transfers, val)
			i++
		}
	}

	defer func() {
		app.txCount++
	}()

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
		Hash:         tmtypes.Tx(req.Tx).Hash(),
	}

	if ret.Code != uint32(sdk.CodeOK) {
		txExtraInfo := &TxExtraInfo{
			Code:      ret.Code,
			Data:      ret.Data,
			Log:       ret.Log,
			Info:      ret.Info,
			GasWanted: ret.GasWanted,
			GasUsed:   ret.GasUsed,
			Events:    ret.Events,
			Codespace: ret.Codespace,
		}
		bytes, errJSON = json.Marshal(txExtraInfo)
		if errJSON == nil {
			n4s.ExtraInfo = string(bytes)
		}
	}

	bytes, errJSON = json.Marshal(n4s)
	if errJSON != nil {
		return
	}

	app.appendPubMsgKV("notify_tx", bytes)
	for _, val := range unbondingMsgList {
		app.appendPubMsgKV("begin_unbonding", val)
	}
	for _, val := range redelegationMsgList {
		app.appendPubMsgKV("begin_redelegation", val)
	}
}

type NotificationBeginRedelegation struct {
	Delegator      string `json:"delegator"`
	ValidatorSrc   string `json:"src"`
	ValidatorDst   string `json:"dst"`
	Amount         string `json:"amount"`
	CompletionTime int64  `json:"completion_time"`
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
			if tmp, err := time.Parse(time.RFC3339, string(attr.Value)); err != nil {
				res.CompletionTime = tmp.Unix()
			}
		}
	}
	for _, attr := range dualEvent[1].Attributes {
		if string(attr.Key) == sdk.AttributeKeySender {
			res.Delegator = string(attr.Value)
		}
	}
	return dex.SafeJSONMarshal(res)
}

type NotificationBeginUnbonding struct {
	Delegator      string `json:"delegator"`
	Validator      string `json:"validator"`
	Amount         string `json:"amount"`
	CompletionTime int64  `json:"completion_time"`
}

func getNotificationBeginUnbonding(dualEvent []abci.Event) []byte {
	var res NotificationBeginUnbonding
	for _, attr := range dualEvent[0].Attributes {
		if string(attr.Key) == stypes.AttributeKeyValidator {
			res.Validator = string(attr.Value)
		} else if string(attr.Key) == sdk.AttributeKeyAmount {
			res.Amount = string(attr.Value)
		} else if string(attr.Key) == stypes.AttributeKeyCompletionTime {
			if tmp, err := time.Parse(time.RFC3339, string(attr.Value)); err != nil {
				res.CompletionTime = tmp.Unix()
			}
		}
	}
	for _, attr := range dualEvent[1].Attributes {
		if string(attr.Key) == sdk.AttributeKeySender {
			res.Delegator = string(attr.Value)
		}
	}
	return dex.SafeJSONMarshal(res)
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
	return dex.SafeJSONMarshal(res)
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
	return dex.SafeJSONMarshal(res)
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
	return dex.SafeJSONMarshal(res)
}

func (app *CetChainApp) notifyBeginBlock(events []abci.Event) {
	//fmt.Printf("========== BeginBlock events ============\n")
	subscribedDistr := app.msgQueProducer.IsSubscribed(distr.ModuleName)
	for _, event := range events {
		//fmt.Printf("= Event: %s\n", event.Type)
		//for _, attr := range event.Attributes {
		//	fmt.Printf("= K: %s; V: %s\n", attr.Key, attr.Value)
		//}
		if event.Type == sltypes.EventTypeSlash {
			val := getNotificationSlash(event)
			app.appendPubMsgKV("slash", val)
		} else if subscribedDistr && event.Type == distrtypes.EventTypeCommission {
			val := getValidatorCommissionMsg(event)
			app.appendPubMsgKV("validator_commission", val)
		} else if subscribedDistr && event.Type == distrtypes.EventTypeRewards {
			val := getDelegatorRewardsMsg(event)
			app.appendPubMsgKV("delegator_rewards", val)
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
			app.appendPubMsgKV("complete_unbonding", val)
		} else if event.Type == stypes.EventTypeCompleteRedelegation {
			val := getNotificationCompleteRedelegation(event)
			app.appendPubMsgKV("complete_redelegation", val)
		}
	}
}

type NotificationValidatorCommission struct {
	Validator  string `json:"validator"`
	Commission string `json:"commission"`
}

func getValidatorCommissionMsg(event abci.Event) []byte {
	var res NotificationValidatorCommission
	for _, attr := range event.Attributes {
		if string(attr.Key) == distrtypes.AttributeKeyValidator {
			res.Validator = string(attr.Value)
		} else if string(attr.Key) == sdk.AttributeKeyAmount {
			res.Commission = string(attr.Value)
		}
	}
	return dex.SafeJSONMarshal(res)
}

type NotificationDelegatorRewards struct {
	Validator string `json:"validator"`
	Rewards   string `json:"rewards"`
}

func getDelegatorRewardsMsg(event abci.Event) []byte {
	var res NotificationDelegatorRewards
	for _, attr := range event.Attributes {
		if string(attr.Key) == distrtypes.AttributeKeyValidator {
			res.Validator = string(attr.Value)
		} else if string(attr.Key) == sdk.AttributeKeyAmount {
			res.Rewards = string(attr.Value)
		}
	}
	return dex.SafeJSONMarshal(res)
}
