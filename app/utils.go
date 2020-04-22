package app

import (
	"os"
	"strings"

	"github.com/coinexchain/cet-sdk/msgqueue"

	cfg "github.com/tendermint/tendermint/config"

	toml "github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

const (
	TSDirCfg = "dir"
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
	config, err := toml.LoadFile(filePath)
	if err != nil {
		return config, err
	}
	config.Set(TSDirCfg, getPreFixBks(msgqueue.CfgPrefixPrune))
	return config, err
}

func isOpenTs() bool {
	bkCfg := getPreFixBks(msgqueue.CfgPrefixPrune)
	return len(bkCfg) > 0
}

func getPreFixBks(prefix string) string {
	brokers := viper.GetStringSlice(msgqueue.FlagBrokers)
	for _, b := range brokers {
		if strings.HasPrefix(b, prefix) {
			return b
		}
	}
	return ""
}
