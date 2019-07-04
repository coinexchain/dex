package init

import (
	"github.com/coinexchain/dex/app"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestValidateGenesisCmd(t *testing.T) {

	ctx := server.NewDefaultContext()
	cdc := app.MakeCodec()
	cmd := ValidateGenesisCmd(ctx, cdc)
	os.Remove("./genesis.json")
	require.Equal(t,
		"Error loading genesis doc from config/genesis.json: open config/genesis.json: no such file or directory",
		cmd.RunE(nil, []string{}).Error())

	defer os.Remove("./genesis.json")
	_, _, err := initializeGenesisFile(cdc, "./genesis.json")
	require.NoError(t, err)
	cmd = ValidateGenesisCmd(ctx, cdc)
	require.Equal(t, nil, cmd.RunE(nil, []string{"./genesis.json"}))
}
