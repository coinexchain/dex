package cli

import (
	"github.com/coinexchain/dex/modules/asset/internal/types"
	dex "github.com/coinexchain/dex/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}
func TestAddGenesisToken(t *testing.T) {
	token := &types.BaseToken{}
	_ = token.SetName("aaa")
	_ = token.SetSymbol("aaa")

	genesis := types.GenesisState{
		Tokens: []types.Token{token},
	}
	err := addGenesisToken(&genesis, token)
	assert.Error(t, err)

	token = &types.BaseToken{}
	_ = token.SetName("bbb")
	_ = token.SetSymbol("bbb")
	_ = addGenesisToken(&genesis, token)
	require.Equal(t, token.GetSymbol(), genesis.Tokens[1].GetSymbol())
}

func TestParseTokenInfo(t *testing.T) {
	defer os.RemoveAll("./keys")
	_, err := parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagOwner, "owner")
	_, err = parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagOwner, "coinex1paehyhx9sxdfwc3rjf85vwn6kjnmzjemtedpnl")
	viper.Set(flagName, "1")
	_, err = parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagName, "aaa")
	viper.Set(flagSymbol, "1")
	_, err = parseTokenInfo()
	assert.Error(t, err)

	viper.Set(flagSymbol, "aaa")
	viper.Set(flagTotalSupply, int64(100))
	viper.Set(flagTotalBurn, int64(100))
	viper.Set(flagTotalMint, int64(100))
	token, _ := parseTokenInfo()
	require.Equal(t, "aaa", token.GetSymbol())
}
