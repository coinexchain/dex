package main

import (
	"testing"
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/coinexchain/dex/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//https://etherscan.io/token/0x081f67afa0ccf8c7b17540767bbe95df2ba8d97f
//date: 2019/08/06 total:5,877,675,270.61317189
const expectedTotalSupply = int64(587767527061317189)

func TestCetSupply(t *testing.T) {
	testAddr := sdk.AccAddress(crypto.AddressHash([]byte("test_addr"))).String()
	cetToken := createCetToken(testAddr)

	require.Equal(t, expectedTotalSupply, cetToken.GetTotalSupply().Int64())

	historicalTotal := cetToken.GetTotalSupply().Add(cetToken.GetTotalBurn())
	require.Equal(t, int64(1000000000000000000), historicalTotal.Int64())
}

func TestTotalCetInGenesisAccounts(t *testing.T) {
	total := sdk.NewCoins()
	for _, account := range createGenesisAccounts() {
		total = total.Add(account.Coins)
	}

	require.Equal(t, expectedTotalSupply, total.AmountOf(types.CET).Int64())
}

func TestUnlockTimeOfVestingAccounts(t *testing.T) {
	accounts := createGenesisAccounts()

	var endTimes []string
	for _, account := range accounts {
		if account.EndTime != 0 {
			require.Equal(t, int64(36000000000000000), account.Coins.AmountOf(types.CET).Int64())

			endTime := time.Unix(account.EndTime, 0).UTC().String()
			endTimes = append(endTimes, endTime)
		}
	}

	expectedEndTimes := []string{
		"2020-01-01 00:00:00 +0000 UTC",
		"2021-01-01 00:00:00 +0000 UTC",
		"2022-01-01 00:00:00 +0000 UTC",
		"2023-01-01 00:00:00 +0000 UTC",
		"2024-01-01 00:00:00 +0000 UTC",
	}
	require.Equal(t, expectedEndTimes, endTimes)
}
