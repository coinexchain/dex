package server

import (
	"fmt"
	"github.com/coinexchain/dex/modules/authx/types"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"

	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"

	dex "github.com/coinexchain/dex/types"
)

func PersistentPreRunEFn(context *sdkserver.Context) func(*cobra.Command, []string) error {
	fn := sdkserver.PersistentPreRunEFn(context)
	return func(cmd *cobra.Command, args []string) error {
		if err := fn(cmd, args); err != nil {
			return err
		}

		tmpConf := cfg.DefaultConfig()
		err := viper.Unmarshal(tmpConf)
		if err != nil {
			// TODO: Handle with #870
			panic(err)
		}
		appConfigFilePath := filepath.Join(tmpConf.RootDir, "config/app.toml")
		if state, err := os.Stat(appConfigFilePath); os.IsNotExist(err) {
			_ = adjustAppConfig()
		} else {
			if state.IsDir() {
				return fmt.Errorf("he specified path is a directory, not a file")
			}
		}
		return err
	}
}

func adjustAppConfig() error {
	tmpConf := cfg.DefaultConfig()
	err := viper.Unmarshal(tmpConf)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}
	rootDir := tmpConf.RootDir
	appConfigFilePath := filepath.Join(rootDir, "config/app.toml")

	appConf, _ := config.ParseConfig()
	// use network_min_gas_price as default value for node_mini_gas_price
	appConf.MinGasPrices = fmt.Sprintf("%s%s", types.DefaultMinGasPriceLimit, dex.DefaultBondDenom)
	config.WriteConfigFile(appConfigFilePath, appConf)
	return nil
}
