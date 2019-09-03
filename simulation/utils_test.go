package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
)

func TestReplaceDenom(t *testing.T) {
	gaccs := []genaccounts.GenesisAccount{
		{
			Coins: sdk.Coins{
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)),
				sdk.NewCoin("btc", sdk.NewInt(100)),
			},
		},
	}

	for _, gacc := range gaccs {
		ReplaceDenom(gacc, sdk.DefaultBondDenom, "cet")
	}

	require.Equal(t, "cet", gaccs[0].Coins[0].Denom)
	require.Equal(t, "btc", gaccs[0].Coins[1].Denom)
}
