package keeper

import (
	dex "github.com/coinexchain/dex/types"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}
