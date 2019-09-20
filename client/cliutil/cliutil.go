package cliutil

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func SetViperWithArgs(args []string) {
	viper.Reset()
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		idx := strings.Index(arg, "=")
		if idx < 0 {
			continue
		}
		viper.Set(arg[2:idx], arg[idx+1:])
	}
}

type MsgWithAccAddress interface {
	sdk.Msg
	SetAccAddress(address sdk.AccAddress)
}

var CliQuery = func(cdc *codec.Codec, query string, param interface{}) error {
	var bz []byte
	var err error
	bz = nil
	if param != nil {
		bz, err = cdc.MarshalJSON(param)
		if err != nil {
			return err
		}
	}

	cliCtx := context.NewCLIContext().WithCodec(cdc)
	res, _, err := cliCtx.QueryWithData(query, bz)
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

var CliRunCommand = func(cdc *codec.Codec, msg MsgWithAccAddress) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	senderAddr := cliCtx.GetFromAddress()
	msg.SetAccAddress(senderAddr)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	generateUnsignedTx := viper.GetBool(FlagGenerateUnsignedTx)
	if generateUnsignedTx {
		return PrintUnsignedTx(cliCtx, txBldr, []sdk.Msg{msg}, senderAddr)
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
