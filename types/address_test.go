package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAddressPrefixes(t *testing.T) {
	InitSdkConfig()
	config := sdk.GetConfig()

	require.Equal(t, "coinex", config.GetBech32AccountAddrPrefix())
	require.Equal(t, "coinexpub", config.GetBech32AccountPubPrefix())
	require.Equal(t, "coinexvaloper", config.GetBech32ValidatorAddrPrefix())
	require.Equal(t, "coinexvaloperpub", config.GetBech32ValidatorPubPrefix())
	require.Equal(t, "coinexvalcons", config.GetBech32ConsensusAddrPrefix())
	require.Equal(t, "coinexvalconspub", config.GetBech32ConsensusPubPrefix())
}
