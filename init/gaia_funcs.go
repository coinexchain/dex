package init

import (
	gaia_init "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
)

var (
	ExportGenesisFile            = gaia_init.ExportGenesisFile
	ExportGenesisFileWithTime    = gaia_init.ExportGenesisFileWithTime
	InitializeNodeValidatorFiles = gaia_init.InitializeNodeValidatorFiles
	LoadGenesisDoc               = gaia_init.LoadGenesisDoc
)
