package incentive

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	dex "github.com/coinexchain/dex/types"
)

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}

func TestIncentiveCoinsAddress(t *testing.T) {
	require.Equal(t, "coinex1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97", IncentivePoolAddr.String())
}
