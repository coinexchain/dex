package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

func TestQueryCmds(t *testing.T) {
	testQueryCmd(t, "params", "custom/asset/parameters", nil)
	testQueryCmd(t, "token abc", "custom/asset/token-info", types.NewQueryAssetParams("abc"))
	testQueryCmd(t, "tokens", "custom/asset/token-list", nil)
	testQueryCmd(t, "whitelist abc", "custom/asset/token-whitelist", types.NewQueryWhitelistParams("abc"))
	testQueryCmd(t, "forbidden-addresses abc", "custom/asset/addr-forbidden", types.NewQueryForbiddenAddrParams("abc"))
	testQueryCmd(t, "reserved-symbols", "custom/asset/reserved-symbols", nil)
}

func testQueryCmd(t *testing.T, args string, expectedPath string, expectedParam interface{}) {
	cmdFactory := func() *cobra.Command {
		return GetQueryCmd(types.ModuleCdc)
	}
	cliutil.TestQueryCmd(t, cmdFactory, args, expectedPath, expectedParam)
}
