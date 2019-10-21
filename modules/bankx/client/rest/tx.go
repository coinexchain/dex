package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/restutil"
)

// sendRequestHandlerFn - http request handler to send coins to a address.
func sendTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	checker := func(cdc *codec.Codec, cliCtx context.CLIContext, req restutil.RestReq) error {
		currentTime := time.Now().Unix()
		unlockTime := req.(*sendReq).UnlockTime

		if unlockTime < 0 {
			return fmt.Errorf("invalid unlock time: %d", unlockTime)
		}
		if unlockTime > 0 && unlockTime < currentTime {
			return fmt.Errorf("unlock time should be later than the current time")
		}

		return nil
	}
	return restutil.NewRestHandlerBuilder(cdc, cliCtx, new(sendReq)).Build(checker)
}

func sendRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return restutil.NewRestHandler(cdc, cliCtx, new(memoReq))
}

func sendSupervisedTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	checker := func(cdc *codec.Codec, cliCtx context.CLIContext, req restutil.RestReq) error {
		currentTime := time.Now().Unix()
		unlockTime := req.(*sendSupervisedReq).UnlockTime

		if unlockTime < currentTime {
			return fmt.Errorf("unlock time should be later than the current time")
		}

		return nil
	}
	return restutil.NewRestHandlerBuilder(cdc, cliCtx, new(sendSupervisedReq)).Build(checker)
}
