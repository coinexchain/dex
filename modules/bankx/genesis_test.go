package bankx_test

import (
	"github.com/coinexchain/dex/modules/bankx"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidate(t *testing.T) {
	genes := bankx.DefaultGenesisState()
	err := genes.ValidateGenesis()
	require.Equal(t, nil, err)
}

func TestInitGenesis(t *testing.T) {
	genes := bankx.DefaultGenesisState()
	bkx, _, ctx := defaultContext()
	bankx.InitGenesis(ctx, *bkx, genes)
	gen := bankx.ExportGenesis(ctx, *bkx)
	require.Equal(t, genes, gen)
}
