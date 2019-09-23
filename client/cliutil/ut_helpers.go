package cliutil

import (
	"strings"

	"github.com/spf13/viper"
)

func SetViperWithArgs(args []string) {
	viper.Reset()
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		idx := strings.Index(arg, "=")
		if idx < 0 {
			continue
		}
		viper.Set(arg[2:idx], arg[idx+1:])
	}
}
