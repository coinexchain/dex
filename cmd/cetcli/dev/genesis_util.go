package dev

import (
	"github.com/coinexchain/dex/modules/asset"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/app"
	dex "github.com/coinexchain/dex/types"
)

func newBaseGenesisAccount(address string, amt int64) app.GenesisAccount {
	return app.NewGenesisAccount(&auth.BaseAccount{
		Address: accAddressFromBech32(address),
		Coins:   dex.NewCetCoins(amt),
	})
}

func newVestingGenesisAccount(address string, amt int64, endTime int64) app.GenesisAccount {
	return app.NewGenesisAccountI(&auth.DelayedVestingAccount{
		BaseVestingAccount: &auth.BaseVestingAccount{
			BaseAccount: &auth.BaseAccount{
				Address: accAddressFromBech32(address),
				Coins:   dex.NewCetCoins(amt),
			},
			OriginalVesting: dex.NewCetCoins(amt),
			EndTime:         endTime,
		},
	})
}

func accAddressFromBech32(address string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return addr
}

func createCetToken() asset.Token {
	token := &asset.BaseToken{
		Name:             "CoinEx Chain Native Token",
		Symbol:           "cet",
		TotalSupply:      588788547005740000,
		Owner:            accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  false,
		TokenForbiddable: false,
		TotalBurn:        411211452994260000,
		TotalMint:        0,
		IsForbidden:      false,
	}
	if err := token.Validate(); err != nil {
		panic(err)
	}
	return token
}
