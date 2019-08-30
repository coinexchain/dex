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
	"unsafe"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
)

var loadedPlugin AppPlugin

type Holder struct {
	plugin unsafe.Pointer
}

func (loader *Holder) GetPlugin() AppPlugin {
	p := (*AppPlugin)(atomic.LoadPointer(&loader.plugin))
	if p == nil {
		return nil
	}

	return *p
}

func (loader *Holder) togglePlugin(logger log.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("toggle plugin failed: %s", string(debug.Stack())))
		}
	}()

	p := loader.GetPlugin()
	if p != nil {
		name := p.Name()
		atomic.StorePointer(&loader.plugin, unsafe.Pointer(nil))
		loadedPlugin = nil
		logger.Info(fmt.Sprintf("plugin %s unloaded", name))
	} else {
		loader.LoadPlugin(logger)
	}
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
		logger.Info(fmt.Sprintf("plugin %s does not exist", pluginPath))
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

	loadedPlugin = symbol.(AppPlugin)
	atomic.StorePointer(&loader.plugin, unsafe.Pointer(&loadedPlugin))

	logger.Info(fmt.Sprintf("plugin %v loaded", loader.GetPlugin().Name()))
}
