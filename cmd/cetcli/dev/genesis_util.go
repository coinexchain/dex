package dev

import (
	"github.com/coinexchain/dex/modules/asset"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"

	dex "github.com/coinexchain/dex/types"
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
		Symbol:           "cet",
		TotalSupply:      588788547005740000,
		Owner:            accAddressFromBech32(ownerAddr),
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  false,
		TokenForbiddable: false,
		TotalBurn:        411211452994260000,
		TotalMint:        0,
		IsForbidden:      false,
		URL:              "https://www.coinex.org",
		Description:      "Decentralized public chain ecosystem, Born for financial liberalization",
	}
	if err := token.Validate(); err != nil {
		panic(err)
	}

	return token
}
