package cli

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/distributionx/types"
)

var testAddr = "coinex12kcupm2x8fw0gglgcz8850kw0k2kx0ff8sr3rn"

func TestDonateTxCmd(t *testing.T) {

	cmd := DonateTxCmd(nil)
	amount := "1000cet"
	fromFlag := fmt.Sprintf("--from=%s", testAddr)
	cmd.SetArgs([]string{amount, fromFlag})
	//cliutil.SetViperWithArgs([]string{amount,fromFlag})

	executed := false
	oldCliRun := cliutil.CliRunCommand
	defer func() {
		cliutil.CliRunCommand = oldCliRun
	}()

	cliutil.CliRunCommand = func(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
		executed = true
		require.Equal(t, amount, msg.(*types.MsgDonateToCommunityPool).Amount.String())
		return nil
	}

	err := cmd.Execute()
	require.Nil(t, err)
	require.True(t, executed)
}
