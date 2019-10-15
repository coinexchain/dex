package plugin

import (
	"os/exec"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func TestTogglePlugin(t *testing.T) {
	cmd := exec.Command("/bin/bash", "./test_plugin/build_test_plugin.sh")
	err := cmd.Run()
	require.Nil(t, err)

	defer func() {
		cmd = exec.Command("rm", "./test_plugin/data/plugin.so")
		_ = cmd.Run()
	}()

	logger := log.NewNopLogger()
	holder := Holder{}
	holder.WaitPluginToggleSignal(logger)

	// invalid path
	viper.Set(flags.FlagHome, "./invalid/")
	holder.togglePlugin()
	require.Equal(t, int32(0), holder.isEnabled)
	require.Nil(t, holder.GetPlugin())

	// valid path
	viper.Set(flags.FlagHome, "./test_plugin/")
	holder.togglePlugin()
	require.Equal(t, int32(1), holder.isEnabled)
	require.NotNil(t, holder.GetPlugin())

	holder.togglePlugin()
	require.Equal(t, int32(0), holder.isEnabled)
	require.Nil(t, holder.GetPlugin())

	holder.togglePlugin()
	require.Equal(t, int32(1), holder.isEnabled)
	require.NotNil(t, holder.GetPlugin())
}
