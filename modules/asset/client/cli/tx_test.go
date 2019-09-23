package cli

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

const testAddrBech32 = "coinex1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97"

func TestTxCmds(t *testing.T) {
	testAddr, _ := sdk.AccAddressFromBech32(testAddrBech32)

	testTxCmd(t, "transfer-ownership --symbol=abc --new-owner={testAddrBech32}",
		types.NewMsgTransferOwnership("abc", nil, testAddr))

	testTxCmd(t, "mint-token --symbol=abc --amount=10000000000000000",
		types.NewMsgMintToken("abc", sdk.NewInt(10000000000000000), nil))

	testTxCmd(t, "burn-token --symbol=abc --amount=10000000000000000",
		types.NewMsgBurnToken("abc", sdk.NewInt(10000000000000000), nil))

	testTxCmd(t, "forbid-token --symbol=abc",
		types.NewMsgForbidToken("abc", nil))

	testTxCmd(t, "unforbid-token --symbol=abc",
		types.NewMsgUnForbidToken("abc", nil))

	testTxCmd(t, "add-whitelist --symbol=abc --whitelist={testAddrBech32}",
		types.NewMsgAddTokenWhitelist("abc", nil, []sdk.AccAddress{testAddr}))

	testTxCmd(t, "remove-whitelist --symbol=abc --whitelist={testAddrBech32}",
		types.NewMsgRemoveTokenWhitelist("abc", nil, []sdk.AccAddress{testAddr}))

	testTxCmd(t, "forbid-addr --symbol=abc --addresses={testAddrBech32}",
		types.NewMsgForbidAddr("abc", nil, []sdk.AccAddress{testAddr}))

	testTxCmd(t, "unforbid-addr --symbol=abc --addresses={testAddrBech32}",
		types.NewMsgUnForbidAddr("abc", nil, []sdk.AccAddress{testAddr}))

	testTxCmd(t, "modify-token-info --symbol=abc --url=coinex.org --description=cool --identity=CET",
		types.NewMsgModifyTokenInfo("abc", "coinex.org", "cool", "CET", nil))
}

func testTxCmd(t *testing.T, args string, expectedMsg interface{}) {
	executed := false
	oldCliRun := cliutil.CliRunCommand
	cliutil.CliRunCommand = func(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
		executed = true
		require.Equal(t, val2ptr(expectedMsg), msg)
		return nil
	}
	defer func() {
		cliutil.CliRunCommand = oldCliRun
	}()

	args1 := strings.Replace(args, "{testAddrBech32}", testAddrBech32, -1)
	argArr := strings.Split(args1, " ")
	cliutil.SetViperWithArgs(argArr)

	cmd := GetTxCmd(types.ModuleCdc)
	cmd.SetArgs(argArr)
	err := cmd.Execute()
	require.NoError(t, err)
	require.True(t, executed)
}

func val2ptr(msg interface{}) cliutil.MsgWithAccAddress {
	vp := reflect.New(reflect.TypeOf(msg))
	vp.Elem().Set(reflect.ValueOf(msg))
	return vp.Interface().(cliutil.MsgWithAccAddress)
}
