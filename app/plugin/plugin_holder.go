package plugin

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"plugin"
	"runtime/debug"
	"sync/atomic"
	"syscall"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
)

type Holder struct {
	isLoadedFlag int32

	pluginInstance AppPlugin
}

func (loader *Holder) GetPlugin() AppPlugin {
	isLoaded := atomic.LoadInt32(&loader.isLoadedFlag)
	if isLoaded == 0 {
		return nil
	}

	return loader.pluginInstance
}

func (loader *Holder) togglePlugin(logger log.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("toggle plugin failed: %s", string(debug.Stack())))
		}
	}()

	p := loader.GetPlugin()
	if p != nil {
		loader.disablePlugin(p, logger)
	} else {
		loader.enablePlugin(logger)
	}
}

func (loader *Holder) enablePlugin(logger log.Logger) {
	atomic.StoreInt32(&loader.isLoadedFlag, 1)
	p := loader.GetPlugin()
	if p != nil {
		logger.Info(fmt.Sprintf("plugin %s enabled", p.Name()))
	}
}

func (loader *Holder) disablePlugin(p AppPlugin, logger log.Logger) {
	atomic.StoreInt32(&loader.isLoadedFlag, 0)
	logger.Info(fmt.Sprintf("plugin %s disabled", p.Name()))
}

func (loader *Holder) WaitPluginToggleSignal(logger log.Logger) {
	togglePlugin := func(c chan os.Signal) {
		for {
			<-c
			loader.togglePlugin(logger)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go togglePlugin(c)
}

func (loader *Holder) LoadPlugin(logger log.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("load plugin failed: %s", string(debug.Stack())))
		}
	}()

	rootDir := viper.GetString(flags.FlagHome)
	pluginPath := path.Join(rootDir, "data/plugin.so")

	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		logger.Error(fmt.Sprintf("plugin %s open failed, %s", pluginPath, err.Error()))
		return
	}

	symbol, err := p.Lookup("Instance")
	if err != nil {
		logger.Error(fmt.Sprintf("Lookup Instance in plugin %s failed", pluginPath))
		return
	}

	instance, ok := symbol.(AppPlugin)
	if !ok {
		logger.Error(fmt.Sprintf("Instance in plugin %s is invalid", pluginPath))
		return
	}

	loader.pluginInstance = instance
	loader.enablePlugin(logger)
}
