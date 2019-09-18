package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func TestNewCetCoinE8(t *testing.T) {
	type args struct {
		amount int64
	}
	tests := []struct {
		name string
		args args
		want sdk.Coin
	}{
		{name: "cet", args: args{1}, want: sdk.NewInt64Coin("cet", E8)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCetCoinE8(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCetCoinE8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCetCoinsE8(t *testing.T) {
	type args struct {
		amount int64
	}
	tests := []struct {
		name string
		args args
		want sdk.Coins
	}{
		{name: "cet", args: args{1}, want: []sdk.Coin{sdk.NewInt64Coin("cet", E8)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCetCoinsE8(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCetCoinsE8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCET(t *testing.T) {
	type args struct {
		coin sdk.Coin
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "cet", args: args{sdk.NewInt64Coin("cet", 1)}, want: true},
		{name: "btc", args: args{sdk.NewInt64Coin("btc", 1)}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCET(tt.args.coin); got != tt.want {
				t.Errorf("IsCET() = %v, want %v", got, tt.want)
			}
		})
	}
}
