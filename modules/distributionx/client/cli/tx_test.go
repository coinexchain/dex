package cli

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/distributionx/types"
	dex "github.com/coinexchain/dex/types"
)

var testAddrBech32 = "coinex12kcupm2x8fw0gglgcz8850kw0k2kx0ff8sr3rn"

func TestDonateTxCmd(t *testing.T) {
	cmdFactory := func() *cobra.Command {
		return DonateTxCmd(nil)
	}

	testAddr, _ := sdk.AccAddressFromBech32(testAddrBech32)
	args := fmt.Sprintf("1000cet --from=%s", testAddr)
	msg := types.NewMsgDonateToCommunityPool(testAddr, dex.NewCetCoins(1000))

	cliutil.TestTxCmd(t, cmdFactory, args, &msg)
}
