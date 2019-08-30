package plugin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

type AppPlugin interface {
	PreCheckTx(abci.RequestCheckTx, sdk.TxDecoder, log.Logger) sdk.Error
	Name() string
}
