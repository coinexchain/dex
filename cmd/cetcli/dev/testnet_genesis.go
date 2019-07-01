package dev

import (
	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/coinexchain/dex/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

func createTestnetGenesisState(cdc *codec.Codec) app.GenesisState {
	genState := app.NewDefaultGenesisState()
	genState.Accounts = createGenesisAccounts()
	genState.AssetData = createTestnetGenesisAssetData()

	addNonBondableAddresses(&genState.StakingXData.Params)

	checkGenState(&genState)

	return genState
}

func addNonBondableAddresses(stakingxParam *stakingx.Params) {
	addNonBondableAddress(stakingxParam, incentive.PoolAddr.String())
	addNonBondableAddress(stakingxParam, viper.GetString(flagAddrCoinExFoundation))
	addNonBondableAddress(stakingxParam, viper.GetString(flagAddrVesting2020))
	addNonBondableAddress(stakingxParam, viper.GetString(flagAddrVesting2021))
	addNonBondableAddress(stakingxParam, viper.GetString(flagAddrVesting2022))
	addNonBondableAddress(stakingxParam, viper.GetString(flagAddrVesting2023))
	addNonBondableAddress(stakingxParam, viper.GetString(flagAddrVesting2024))
}

func addNonBondableAddress(params *stakingx.Params, address string) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err.Error())
	}

	params.NonBondableAddresses = append(params.NonBondableAddresses, addr)
}

func createTestnetGenesisAssetData() asset.GenesisState {
	state := asset.DefaultGenesisState()
	cetOwnerAddr := viper.GetString(flagAddrCoinExFoundation)
	state.Tokens = append(state.Tokens, createCetToken(cetOwnerAddr))
	return state
}

func createGenesisAccounts() (accs []app.GenesisAccount) {
	accs = append(accs,
		newBaseGenesisAccount(incentive.PoolAddr.String(), 31536000000000000),
		newBaseGenesisAccount(viper.GetString(flagAddrCirculation), 288788547005740000),
		newBaseGenesisAccount(viper.GetString(flagAddrCoinExFoundation), 88464000000000000),
		newVestingGenesisAccount(viper.GetString(flagAddrVesting2020), 36000000000000000, 1577836800),
		newVestingGenesisAccount(viper.GetString(flagAddrVesting2021), 36000000000000000, 1609459200),
		newVestingGenesisAccount(viper.GetString(flagAddrVesting2022), 36000000000000000, 1640995200),
		newVestingGenesisAccount(viper.GetString(flagAddrVesting2023), 36000000000000000, 1672531200),
		newVestingGenesisAccount(viper.GetString(flagAddrVesting2024), 36000000000000000, 1704067200),
	)

	return
}

func checkGenState(genState *app.GenesisState) {
	tokens := genState.AssetData.Tokens
	if len(tokens) != 1 || tokens[0].GetSymbol() != types.CET {
		panic("only CET token should exists during network initial genesis")
	}

	if tokens[0].GetOwner().String() != viper.GetString(flagAddrCoinExFoundation) {
		panic("owner of CET should be addr of CoinEx Foundation")
	}
}
