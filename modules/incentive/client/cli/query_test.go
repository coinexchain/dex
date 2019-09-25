package cli

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/incentive/internal/keepers"
	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

func TestQueryParamsCmd(t *testing.T) {
	cmdFactory := func() *cobra.Command {
		return GetQueryCmd(nil)
	}

	cliutil.TestQueryCmd(t, cmdFactory, "params",
		fmt.Sprintf("custom/%s/%s", types.ModuleName, keepers.QueryParameters), nil)
}
