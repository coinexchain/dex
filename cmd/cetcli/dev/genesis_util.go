package dev

import (
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
