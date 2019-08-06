package app

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"reflect"
)

//	sktypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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

//type NotificationCompleteRedelegation struct {
//	Delegator    string `json:"delegator"`
//	ValidatorSrc string `json:"src"`
//	ValidatorDst string `json:"dst"`
//}
//
//type NotificationCompleteUnbonding struct {
//	Delegator    string `json:"delegator"`
//	Validator    string `json:"validator"`
//}
//
//func (app *CetChainApp) notifyDelegators(req abci.RequestDeliverTx, events []abci.Event) {
//sktypes.EventTypeCompleteRedelegation,
//sktypes.EventTypeCompleteUnbonding,
//EventTypeUnbond,
//EventTypeRedelegate,
//
//// UnbondingDelegation stores all of a single delegator's unbonding bonds
//// for a single validator in an time-ordered list
//type UnbondingDelegation struct {
//	DelegatorAddress sdk.AccAddress             `json:"delegator_address" yaml:"delegator_address"` // delegator
//	ValidatorAddress sdk.ValAddress             `json:"validator_address" yaml:"validator_address"` // validator unbonding from operator addr
//	Entries          []UnbondingDelegationEntry `json:"entries" yaml:"entries"`                     // unbonding delegation entries
//}
//
//// UnbondingDelegationEntry - entry to an UnbondingDelegation
//type UnbondingDelegationEntry struct {
//	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height which the unbonding took place
//	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the unbonding delegation will complete
//	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"` // atoms initially scheduled to receive at completion
//	Balance        sdk.Int   `json:"balance" yaml:"balance"`                 // atoms to receive at completion
//}
//
//// Redelegation contains the list of a particular delegator's
//// redelegating bonds from a particular source validator to a
//// particular destination validator
//type Redelegation struct {
//	DelegatorAddress    sdk.AccAddress      `json:"delegator_address" yaml:"delegator_address"`         // delegator
//	ValidatorSrcAddress sdk.ValAddress      `json:"validator_src_address" yaml:"validator_src_address"` // validator redelegation source operator addr
//	ValidatorDstAddress sdk.ValAddress      `json:"validator_dst_address" yaml:"validator_dst_address"` // validator redelegation destination operator addr
//	Entries             []RedelegationEntry `json:"entries" yaml:"entries"`                             // redelegation entries
//}
//
//// RedelegationEntry - entry to a Redelegation
//type RedelegationEntry struct {
//	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height at which the redelegation took place
//	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the redelegation will complete
//	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"` // initial balance when redelegation started
//	SharesDst      sdk.Dec   `json:"shares_dst" yaml:"shares_dst"`           // amount of destination-validator shares created by redelegation
//}
//
//
//// return a given amount of all the delegator unbonding-delegations
//func (k Keeper) GetUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress,
//	maxRetrieve uint16) (unbondingDelegations []types.UnbondingDelegation) {
//// return a given amount of all the delegator redelegations
//func (k Keeper) GetRedelegations(ctx sdk.Context, delegator sdk.AccAddress,
//	maxRetrieve uint16) (redelegations []types.Redelegation) {
//
//}
//
