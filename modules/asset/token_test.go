package asset

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNewToken(t *testing.T) {
	type args struct {
		name           string
		symbol         string
		amt            int64
		owner          sdk.AccAddress
		mintable       bool
		burnable       bool
		addrfreezable  bool
		tokenfreezable bool
	}
	tests := []struct {
		name    string
		args    args
		want    *BaseToken
		wantErr sdk.Error
	}{
		{
			"base-case",
			args{
				"ABC Token",
				"abc",
				210000000000,
				tAccAddr,
				false,
				false,
				false,
				false,
			},
			&BaseToken{
				"ABC Token",
				"abc",
				210000000000,
				tAccAddr,
				false,
				false,
				false,
				false,
				0,
				0,
				false,
			},
			nil,
		},
		{
			"caseMissOwner",
			args{
				"ABC Token",
				"abc",
				210000000000,
				sdk.AccAddress{},
				false,
				false,
				false,
				false,
			},
			nil,
			ErrorInvalidTokenOwner("token owner is invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewToken(tt.args.name, tt.args.symbol, tt.args.amt, tt.args.owner, tt.args.mintable, tt.args.burnable, tt.args.addrfreezable, tt.args.tokenfreezable)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewToken() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantErr) {
				t.Errorf("NewToken() got1 = %v, want %v", got1, tt.wantErr)
			}
		})
	}
}

func TestBaseToken_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		token   *BaseToken
		wantErr error
	}{
		{
			"base-case",
			&BaseToken{
				"ABC Token",
				"abc",
				210000000000,
				tAccAddr,
				false,
				false,
				false,
				false,
				0,
				0,
				false,
			},
			nil,
		},
		{
			"case-invalid",
			&BaseToken{
				"ABC Token",
				"abc",
				210000000000,
				tAccAddr,
				false,
				false,
				false,
				false,
				0,
				-100000000,
				false,
			},
			ErrorInvalidTokenMint("Invalid total mint: -100000000"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := &BaseToken{
				Name:           tt.token.Name,
				Symbol:         tt.token.Symbol,
				TotalSupply:    tt.token.TotalSupply,
				Owner:          tt.token.Owner,
				Mintable:       tt.token.Mintable,
				Burnable:       tt.token.Burnable,
				AddrFreezable:  tt.token.AddrFreezable,
				TokenFreezable: tt.token.TokenFreezable,
				TotalBurn:      tt.token.TotalBurn,
				TotalMint:      tt.token.TotalMint,
				IsFrozen:       tt.token.IsFrozen,
			}
			if err := base.IsValid(); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("BaseToken.IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
