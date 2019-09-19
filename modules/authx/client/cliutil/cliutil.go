package cliutil

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	authxutils "github.com/coinexchain/dex/modules/authx/client/utils"
)

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

var CliRunCommand = func(cdc *codec.Codec, senderPtr *sdk.AccAddress, msg sdk.Msg) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	*senderPtr = cliCtx.GetFromAddress()
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	generateUnsignedTx := viper.GetBool(authxutils.FlagGenerateUnsignedTx)
	if generateUnsignedTx {
		return authxutils.PrintUnsignedTx(cliCtx, txBldr, []sdk.Msg{msg}, *senderPtr)
	}
	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
