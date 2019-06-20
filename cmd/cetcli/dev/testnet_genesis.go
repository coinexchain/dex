package dev

import (
	"github.com/coinexchain/dex/app"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createTestnetGenesisState(cdc *codec.Codec) app.GenesisState {
	genState := app.NewDefaultGenesisState()
	genState.Accounts = createExampleGenesisAccounts()
	genState.StakingData.Pool.NotBondedTokens = sdk.NewInt(588788547005740000)
	genState.AssetData = createExampleGenesisAssetData()
	genState.MarketData = createExampleGenesisMarketData()
	//genState.GenTxs = append(genState.GenTxs, createExampleGenTx(cdc))
	return genState
}
