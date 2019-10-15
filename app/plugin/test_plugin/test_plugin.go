package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

type TestMsgFilter struct {
}

func (f TestMsgFilter) PreCheckTx(req abci.RequestCheckTx, txDecoder sdk.TxDecoder, logger log.Logger) sdk.Error {
	return nil
}

func (f TestMsgFilter) Name() string {
	return "TestPlugin"
}

// Instance is the exported symbol
var Instance TestMsgFilter
