package init

import (
	//gaia_app "github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	//gaia_init "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
)

var (
	CollectStdTxs                = gaia_app.CollectStdTxs
	ExportGenesisFile            = gaia_init.ExportGenesisFile
	ExportGenesisFileWithTime    = gaia_init.ExportGenesisFileWithTime
	InitializeNodeValidatorFiles = gaia_init.InitializeNodeValidatorFiles
	LoadGenesisDoc               = gaia_init.LoadGenesisDoc
)
