package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/viper"
)

const FlagDefaultHTTP = "default-http"
const FlagSwaggerHost = "swagger-host"

func registerSwaggerUI(rs *lcd.RestServer) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rs.Mux.PathPrefix("/swagger/swagger.yaml").Handler(serveSwagger(statikFS))
	rs.Mux.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", staticServer))
}

func serveSwagger(fs http.FileSystem) http.HandlerFunc {
	file, _ := fs.Open("/swagger.yaml")
	buf, _ := ioutil.ReadAll(file)
	swagger := string(buf)

	swaggerHost := viper.GetString(FlagSwaggerHost)
	if swaggerHost != "" {
		swagger = strings.Replace(swagger, "host: dex-api.coinex.org", "host: "+swaggerHost, -1)
	}

	if viper.GetBool(FlagDefaultHTTP) {
		swagger = strings.Replace(swagger,
			"schemes:\n  - https\n  - http", "schemes:\n  - http\n  - https", -1)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(swagger))
	})
}

func bindSwaggerFlags(cmd *cobra.Command) error {
	if err := viper.BindPFlag(FlagDefaultHTTP, cmd.PersistentFlags().Lookup(FlagDefaultHTTP)); err != nil {
		return err
	}

	return viper.BindPFlag(FlagSwaggerHost, cmd.PersistentFlags().Lookup(FlagSwaggerHost))
}
