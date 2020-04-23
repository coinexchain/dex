package app

import (
	"os"
	"strings"

	"github.com/gorilla/mux"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	"github.com/coinexchain/cet-sdk/msgqueue"
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
	filePath := conf.RootDir + "/config/trade-server.toml"
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return nil, err
	}
	config, err := toml.LoadFile(filePath)
	if err != nil {
		return config, err
	}
	path := strings.Split(getPreFixBks(msgqueue.CfgPrefixPrune), msgqueue.CfgPrefixPrune)[1]
	config.Set(TSDirCfg, path)
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

func CreateContextAndRegisterRoutes(router *mux.Router) {
	var ctx = newCLIContextForEmbeddedLDC()
	client.RegisterRoutes(ctx, router)
	authrest.RegisterTxRoutes(ctx, router)
	ModuleBasics.RegisterRESTRoutes(ctx, router)
}

// see cosmos-sdk/client/context/context.go#NewCLIContextWithFrom()
func newCLIContextForEmbeddedLDC() context.CLIContext {
	var cdc = MakeCodec()
	var nodeURI = viper.GetString(flags.FlagNode)
	if nodeURI == "" {
		nodeURI = "tcp://localhost:26657"
	}
	var rpc = rpcclient.NewHTTP(nodeURI, "/websocket")

	// fill members of ctx
	return context.CLIContext{
		Codec:     cdc,
		Client:    rpc,
		NodeURI:   nodeURI,
		TrustNode: true,

		// default values is enough?
		//Output:        os.Stdout,
		//From:          viper.GetString(flags.FlagFrom),
		//OutputFormat:  viper.GetString(cli.OutputFlag),
		//Height:        viper.GetInt64(flags.FlagHeight),
		//UseLedger:     viper.GetBool(flags.FlagUseLedger),
		//BroadcastMode: viper.GetString(flags.FlagBroadcastMode),
		//Verifier:      verifier,
		//Simulate:      viper.GetBool(flags.FlagDryRun),
		//GenerateOnly:  genOnly,
		//FromAddress:   fromAddress,
		//FromName:      fromName,
		//Indent:        viper.GetBool(flags.FlagIndentResponse),
		//SkipConfirm:   viper.GetBool(flags.FlagSkipConfirmation),
	}
}
