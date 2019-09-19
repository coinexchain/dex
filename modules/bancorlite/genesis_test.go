package bancorlite_test

import (
	"testing"

	"github.com/coinexchain/dex/modules/bancorlite"
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"

	"github.com/stretchr/testify/require"
)

func TestDefaultGenesisState(t *testing.T) {
	genesisState := bancorlite.DefaultGenesisState()
	require.Equal(t, bancorlite.DefaultParams(), genesisState.Params)
	require.Equal(t, 0, len(genesisState.BancorInfoMap))
}

func TestGenesisState_Validate(t *testing.T) {

	type args struct {
		gs bancorlite.GenesisState
	}
	testCases := []struct {
		name string
		args args
		want bool
	}{
		{"negative bancor fee",
			args{
				bancorlite.GenesisState{
					types.Params{
						1,
						-1,
						0,
					},
					make(map[string]bancorlite.BancorInfo),
				},
			},
			false,
		},
		{"pass case",
			args{
				bancorlite.GenesisState{
					types.Params{
						1,
						10,
						100,
					},
					make(map[string]bancorlite.BancorInfo),
				},
			},
			true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.gs.Validate()
			if (got == nil) != tt.want {
				t.Errorf("genesisState validate() = %v, want %v", got, tt.want)
			}
		})
	}

}
