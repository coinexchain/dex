package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressPrefixes(t *testing.T) {
	require.Equal(t, "coinexpub", Bech32PrefixAccPub)
	require.Equal(t, "coinexvaloper", Bech32PrefixValAddr)
	require.Equal(t, "coinexvaloperpub", Bech32PrefixValPub)
	require.Equal(t, "coinexvalcons", Bech32PrefixConsAddr)
	require.Equal(t, "coinexvalconspub", Bech32PrefixConsPub)
}
