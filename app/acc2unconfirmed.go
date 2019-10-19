package app

import (
	"bytes"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	CodeSpaceUnconfirmedLimit sdk.CodespaceType = "unconfirmed_limit"
	CodeTooManyUnconfirmedTx  sdk.CodeType      = 2100
)

var errTooManyUnconfirmedTx = sdk.NewError(CodeSpaceUnconfirmedLimit, CodeTooManyUnconfirmedTx, "Too Many Unconfirmed Transactions")

const (
	SameTxExist      = 1
	OtherTxExist     = 2
	NoTxExist        = 3
	SweepPeriod      = 15 * 60 // 15 minutes
	DefaultLimitTime = 60      // a minute
)

type UnconfirmedTx struct {
	HashID    []byte
	Timestamp int64
	ABCIResp  abci.ResponseCheckTx
}

type Account2UnconfirmedTx struct {
	auMap         map[string]UnconfirmedTx
	limitTime     int64
	removeList    []sdk.AccAddress
	lastSweepTime int64
	mutex         sync.RWMutex
}

func NewAccount2UnconfirmedTx(limitTime int64) *Account2UnconfirmedTx {
	return &Account2UnconfirmedTx{
		auMap:         make(map[string]UnconfirmedTx),
		limitTime:     limitTime,
		removeList:    make([]sdk.AccAddress, 0, 5000),
		lastSweepTime: 0,
	}
}

func (acc2unc *Account2UnconfirmedTx) Lookup(addr sdk.AccAddress, hashid []byte, timestamp int64) (int, *abci.ResponseCheckTx) {
	acc2unc.mutex.RLock()
	unconfirmedTx, ok := acc2unc.auMap[string(addr)]
	acc2unc.mutex.RUnlock()
	expired := timestamp-unconfirmedTx.Timestamp > acc2unc.limitTime
	if !ok || expired {
		return NoTxExist, nil
	}
	if bytes.Equal(unconfirmedTx.HashID, hashid) {
		return SameTxExist, &unconfirmedTx.ABCIResp
	}
	return OtherTxExist, nil
}

func (acc2unc *Account2UnconfirmedTx) Add(addr sdk.AccAddress, hashid []byte, timestamp int64, resp abci.ResponseCheckTx) {
	acc2unc.mutex.Lock()
	acc2unc.auMap[string(addr)] = UnconfirmedTx{HashID: hashid, Timestamp: timestamp, ABCIResp: resp}
	acc2unc.mutex.Unlock()
}

func (acc2unc *Account2UnconfirmedTx) AddToRemoveList(addrs []sdk.AccAddress) {
	//acc2unc.mutex.Lock()
	acc2unc.removeList = append(acc2unc.removeList, addrs...)
	//acc2unc.mutex.Unlock()
}

func (acc2unc *Account2UnconfirmedTx) CommitRemove(timestamp int64) {
	//acc2unc.mutex.Lock()
	for _, addr := range acc2unc.removeList {
		s := string(addr)
		delete(acc2unc.auMap, s) // will do nothing if key not existing
	}
	if timestamp-acc2unc.lastSweepTime > SweepPeriod {
		acc2unc.mutex.Lock()
		for acc, unconfirmedTx := range acc2unc.auMap {
			expired := timestamp-unconfirmedTx.Timestamp > acc2unc.limitTime
			if expired {
				delete(acc2unc.auMap, acc)
			}
		}
		acc2unc.mutex.Unlock()
		acc2unc.lastSweepTime = timestamp
	}
}

func (acc2unc *Account2UnconfirmedTx) ClearRemoveList() {
	//acc2unc.mutex.Lock()
	acc2unc.removeList = acc2unc.removeList[:0]
	//acc2unc.mutex.Unlock()
}
