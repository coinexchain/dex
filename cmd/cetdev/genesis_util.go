package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"

	"github.com/coinexchain/cet-sdk/modules/asset"
	dex "github.com/coinexchain/cet-sdk/types"
)

func newBaseGenesisAccount(address string, amt int64) genaccounts.GenesisAccount {
	return genaccounts.NewGenesisAccount(&auth.BaseAccount{
		Address: accAddressFromBech32(address),
		Coins:   dex.NewCetCoins(amt),
	})
}

func newVestingGenesisAccount(address string, amt int64, endTime int64) genaccounts.GenesisAccount {
	acc, err := genaccounts.NewGenesisAccountI(&auth.DelayedVestingAccount{
		BaseVestingAccount: &auth.BaseVestingAccount{
			BaseAccount: &auth.BaseAccount{
				Address: accAddressFromBech32(address),
				Coins:   dex.NewCetCoins(amt),
			},
			OriginalVesting: dex.NewCetCoins(amt),
			EndTime:         endTime,
		},
	})
	if err != nil {
		panic(err)
	}
	return acc
}

func accAddressFromBech32(address string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return addr
}

func createCetToken(ownerAddr string) asset.Token {
	token := &asset.BaseToken{
		Name:             "CoinEx Chain Native Token",
		Symbol:           dex.CET,
		TotalSupply:      sdk.NewInt(587767527061317189),
		Owner:            accAddressFromBech32(ownerAddr),
		SendLock:         sdk.ZeroInt(),
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  false,
		TokenForbiddable: false,
		TotalBurn:        sdk.NewInt(412232472938682811),
		TotalMint:        sdk.ZeroInt(),
		IsForbidden:      false,
		URL:              "https://www.coinex.org",
		Description:      "Decentralized public chain ecosystem, Born for financial liberalization",
	}
	if err := token.Validate(); err != nil {
		panic(err)
	}

	return token
}
