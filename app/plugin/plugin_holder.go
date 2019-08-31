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

type Holder struct {
	isEnabled      int32
	pluginInstance AppPlugin
}

func (loader *Holder) isPluginLoaded() bool {
	return loader.pluginInstance != nil
}

func (loader *Holder) GetPlugin() AppPlugin {
	isEnabled := atomic.LoadInt32(&loader.isEnabled) == 1
	if isEnabled {
		return loader.pluginInstance
	}

	return nil
}

func (loader *Holder) togglePlugin(logger log.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("toggle plugin failed: %s", string(debug.Stack())))
		}
	}()

	if !loader.isPluginLoaded() {
		loader.loadAndEnablePlugin(logger)
		return
	}

	isEnabled := atomic.LoadInt32(&loader.isEnabled) == 1
	if isEnabled {
		loader.disablePlugin(logger)
	} else {
		loader.enablePlugin(logger)
	}
}

func (loader *Holder) enablePlugin(logger log.Logger) {
	atomic.StoreInt32(&loader.isEnabled, 1)

	if p := loader.GetPlugin(); p != nil {
		logger.Info(fmt.Sprintf("plugin %s isEnabled", p.Name()))
	}
}

func (loader *Holder) disablePlugin(logger log.Logger) {
	atomic.StoreInt32(&loader.isEnabled, 0)

	if loader.isPluginLoaded() {
		logger.Info(fmt.Sprintf("plugin %s disabled", loader.pluginInstance.Name()))
	}
}

func (loader *Holder) loadAndEnablePlugin(logger log.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("load plugin failed: %s", string(debug.Stack())))
		}
	}()

	rootDir := viper.GetString(flags.FlagHome)
	pluginPath := path.Join(rootDir, "data/plugin.so")

	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		logger.Info(fmt.Sprintf("plugin %s not exists", pluginPath))
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
