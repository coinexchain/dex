package cliutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const FlagGenerateUnsignedTx = "generate-unsigned-tx"

func PrintUnsignedTx(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, msgs []sdk.Msg, from sdk.AccAddress) error {
	num, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(from)
	if err != nil {
		return err
	}
	if txBldr.AccountNumber() == 0 {
		txBldr = txBldr.WithAccountNumber(num)
	}
	if txBldr.Sequence() == 0 {
		txBldr = txBldr.WithSequence(seq)
	}
	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(cliCtx.Output, "%s\n", stdSignMsg.Bytes())
	return nil
}
