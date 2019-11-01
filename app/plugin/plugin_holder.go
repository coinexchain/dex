package plugin

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"plugin"
	"runtime/debug"
	"sync/atomic"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
)

var reloadPluginSignal os.Signal

func SetReloadPluginSignal(signal os.Signal) {
	reloadPluginSignal = signal
}

func (loader *Holder) WaitPluginToggleSignal(logger log.Logger) {
	loader.logger = logger
	togglePlugin := func(c chan os.Signal) {
		for {
			<-c
			loader.togglePlugin()
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, reloadPluginSignal)
	go togglePlugin(c)
}

type Holder struct {
	isEnabled      int32
	pluginInstance AppPlugin
	logger         log.Logger
}

func (loader *Holder) isPluginLoaded() bool {
	return loader.pluginInstance != nil
}

func (loader *Holder) GetPlugin() AppPlugin {
	if loader.isPluginEnabled() {
		return loader.pluginInstance
	}

	return nil
}

func (loader *Holder) togglePlugin() {
	defer func() {
		if r := recover(); r != nil {
			loader.logger.Error(fmt.Sprintf("toggle plugin failed: %s", string(debug.Stack())))
		}
	}()

	if !loader.isPluginLoaded() {
		loader.loadAndEnablePlugin()
		return
	}

	if loader.isPluginEnabled() {
		loader.disablePlugin()
	} else {
		loader.enablePlugin()
	}
}

func (loader *Holder) isPluginEnabled() bool {
	return atomic.LoadInt32(&loader.isEnabled) == 1
}

func (loader *Holder) enablePlugin() {
	atomic.StoreInt32(&loader.isEnabled, 1)

	if loader.pluginInstance != nil {
		loader.logger.Info(fmt.Sprintf("plugin %s is enabled", loader.pluginInstance.Name()))
	}
}

func (loader *Holder) disablePlugin() {
	atomic.StoreInt32(&loader.isEnabled, 0)

	if loader.pluginInstance != nil {
		loader.logger.Info(fmt.Sprintf("plugin %s is disabled", loader.pluginInstance.Name()))
	}
}

func (loader *Holder) loadAndEnablePlugin() {
	defer func() {
		if r := recover(); r != nil {
			loader.logger.Error(fmt.Sprintf("load plugin failed: %s", string(debug.Stack())))
		}
	}()

	rootDir := viper.GetString(flags.FlagHome)
	pluginPath := path.Join(rootDir, "data/plugin.so")

	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		loader.logger.Error(fmt.Sprintf("plugin %s not exists", pluginPath))
		return
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		loader.logger.Error(fmt.Sprintf("plugin %s open failed, %s", pluginPath, err.Error()))
		return
	}

	symbol, err := p.Lookup("Instance")
	if err != nil {
		loader.logger.Error(fmt.Sprintf("Lookup Instance in plugin %s failed", pluginPath))
		return
	}

	instance, ok := symbol.(AppPlugin)
	if !ok {
		loader.logger.Error(fmt.Sprintf("Instance in plugin %s is invalid", pluginPath))
		return
	}

	loader.pluginInstance = instance
	loader.enablePlugin()
}
