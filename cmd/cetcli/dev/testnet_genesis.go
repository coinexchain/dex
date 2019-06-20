package dev

import (
	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createTestnetGenesisState(cdc *codec.Codec) app.GenesisState {
	genState := app.NewDefaultGenesisState()
	genState.Accounts = createExampleGenesisAccounts()
	genState.StakingData.Pool.NotBondedTokens = sdk.NewInt(588788547005740000)
	genState.AssetData = createTestnetGenesisAssetData()
	return genState
}

func createTestnetGenesisAssetData() asset.GenesisState {
	state := asset.DefaultGenesisState()
	state.Tokens = append(state.Tokens, createCetToken())
	return state
}
