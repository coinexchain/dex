package bankx

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidate(t *testing.T) {
	genes := DefaultGenesisState()
	err := genes.Validate()
	require.Equal(t, nil, err)
}

func TestInitGenesis(t *testing.T) {
	genes := DefaultGenesisState()
	input := setupTestInput()
	InitGenesis(input.ctx, input.bxk, genes)
	gen := ExportGenesis(input.ctx, input.bxk)
	require.Equal(t, genes, gen)
}
