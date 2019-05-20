package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewCetCoin(amount int64) sdk.Coin {
	return sdk.NewCoin(CET, sdk.NewInt(amount))
}

func NewCetCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(NewCetCoin(amount))
}
