package main

import (
	"github.com/coinexchain/dex/app/plugin"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/cosmos/cosmos-sdk/x/auth"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	CodeSpacePlugin      sdk.CodespaceType = "plugin"
	CodeNotAcceptableMsg sdk.CodeType      = 2000
)

var errNotAcceptableMsg = sdk.NewError(CodeSpacePlugin, CodeNotAcceptableMsg, "")

type MsgFilter struct {
}

func (f MsgFilter) PreCheckTx(req abci.RequestCheckTx, txDecoder sdk.TxDecoder, logger log.Logger) sdk.Error {
	tx, err := txDecoder(req.Tx)
	if err != nil {
		return err
	}

	for _, msg := range tx.GetMsgs() {
		switch msg.(type) {
		case bankx.MsgSend:
			stdTx := tx.(auth.StdTx)
			if stdTx.Fee.Amount.AmountOf("cet").Int64() >= 1000000000000 {
				return errNotAcceptableMsg
			}
		}
	}

	return nil
}

func (f MsgFilter) Name() string {
	return "SimplePlugin"
}

var _ plugin.AppPlugin = (*MsgFilter)(nil)

// Instance is the exported symbol
var Instance MsgFilter
