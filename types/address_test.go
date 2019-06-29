package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAddressPrefixes(t *testing.T) {
	require.Equal(t, "coinexpub", Bech32PrefixAccPub)
	require.Equal(t, "coinexvaloper", Bech32PrefixValAddr)
	require.Equal(t, "coinexvaloperpub", Bech32PrefixValPub)
	require.Equal(t, "coinexvalcons", Bech32PrefixConsAddr)
	require.Equal(t, "coinexvalconspub", Bech32PrefixConsPub)
}

func TestInitSdkConfig(t *testing.T) {

	InitSdkConfig()
	config := sdk.GetConfig()
	require.Equal(t,config.GetBech32AccountAddrPrefix(),"coinex")
	require.Equal(t,config.GetBech32ValidatorAddrPrefix(),"coinexvaloper")
	require.Equal(t,config.GetBech32ValidatorAddrPrefix(),"coinexvaloper")


}

