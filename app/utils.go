package app

import (
	"os"

	cfg "github.com/tendermint/tendermint/config"

	toml "github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

func initConf() (*toml.Tree, error) {
	conf := cfg.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	filePath := conf.RootDir + "config/trade-server.toml"
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return nil, err
	}
	return toml.LoadFile(filePath)
}
