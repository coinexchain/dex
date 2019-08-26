package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"testing"
)

func TestBaseToken_Validate(t *testing.T) {
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
				sdk.NewInt(2100),
				sdk.ZeroInt(),
				testAddr,
				false,
				false,
				false,
				false,
				sdk.ZeroInt(),
				sdk.ZeroInt(),
				false,
				"",
				"",
				TestIdentityString,
			},
			nil,
		},
		{
			"case-invalid-total-mint",
			&BaseToken{
				"ABC Token",
				"abc",
				sdk.NewInt(2100),
				sdk.ZeroInt(),
				testAddr,
				false,
				false,
				false,
				false,
				sdk.ZeroInt(),
				sdk.NewInt(2100),
				false,
				"",
				"",
				TestIdentityString,
			},
			ErrTokenMintNotSupported("abc"),
		},
		{
			"case-invalid-total-burn",
			&BaseToken{
				"ABC Token",
				"abc",
				sdk.NewInt(2100),
				sdk.ZeroInt(),
				testAddr,
				false,
				false,
				false,
				false,
				sdk.NewInt(2100),
				sdk.ZeroInt(),
				false,
				"",
				"",
				TestIdentityString,
			},
			ErrTokenBurnNotSupported("abc"),
		},
		{
			"case-invalid-forbidden-state",
			&BaseToken{
				"ABC Token",
				"abc",
				sdk.NewInt(2100),
				sdk.ZeroInt(),
				testAddr,
				false,
				false,
				false,
				false,
				sdk.ZeroInt(),
				sdk.ZeroInt(),
				true,
				"",
				"",
				TestIdentityString,
			},
			ErrTokenForbiddenNotSupported("abc"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.token.Validate(); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("BaseToken.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
