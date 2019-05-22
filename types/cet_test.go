package types

import (
	"testing"
)

func TestNewCetCoin(t *testing.T) {

	coin := NewCetCoin(1)

	if coin.Amount.Int64() != 1 {
		t.Error("coin is not 1")
	}
	coin = NewCetCoin(0)

	if coin.Amount.Int64() != 0 {
		t.Error("coin is not 0")
	}
}

func TestNewCetCoins(t *testing.T) {
	coins := NewCetCoins(1)

	if coins[0].Amount.Int64() != 1 {
		t.Error("coin is not 1")
	}
}
