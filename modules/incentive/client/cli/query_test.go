package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/incentive/internal/keepers"
	"github.com/coinexchain/dex/modules/incentive/internal/types"
)

func TestQueryParamsCmd(t *testing.T) {
	expectedPath := fmt.Sprintf("custom/%s/%s", types.ModuleName, keepers.QueryParameters)
	oldCliQuery := cliutil.CliQuery
	executed := false
	cliutil.CliQuery = func(cdc *codec.Codec, path string, param interface{}) error {
		require.Equal(t, expectedPath, path)
		require.Equal(t, nil, param)
		executed = true
		return nil
	}
	defer func() {
		cliutil.CliQuery = oldCliQuery
	}()
	cmd := GetQueryCmd(nil)
	cmd.SetArgs([]string{"params"})
	err := cmd.Execute()
	require.Nil(t, nil, err)
	require.True(t, executed)
}
