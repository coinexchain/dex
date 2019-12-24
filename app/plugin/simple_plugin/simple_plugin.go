package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/coinexchain/cet-sdk/modules/bankx"
	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app/plugin"
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
			if stdTx.Fee.Amount.AmountOf(dex.CET).Int64() >= 1000000000000 {
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
