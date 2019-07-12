package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"

	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"

	"github.com/coinexchain/dex/modules/authx"
	dex "github.com/coinexchain/dex/types"
)

func PersistentPreRunEFn(context *sdkserver.Context) func(*cobra.Command, []string) error {
	fn := sdkserver.PersistentPreRunEFn(context)
	return func(cmd *cobra.Command, args []string) error {
		createAppConfigFile()
		return fn(cmd, args)
	}
}

func createAppConfigFile() {
	tmpConf := cfg.DefaultConfig()
	err := viper.Unmarshal(tmpConf)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}
	rootDir := tmpConf.RootDir

	appConfigFilePath := filepath.Join(rootDir, "config/app.toml")
	if _, err := os.Stat(appConfigFilePath); os.IsNotExist(err) {
		appConf, _ := config.ParseConfig()
		// use network_min_gas_price as default value for node_mini_gas_price
		appConf.MinGasPrices = fmt.Sprintf("%s%s", authx.DefaultMinGasPriceLimit, dex.DefaultBondDenom)
		config.WriteConfigFile(appConfigFilePath, appConf)
	}
}
