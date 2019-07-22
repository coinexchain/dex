package types

import (
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
				210000000000,
				testAddr,
				false,
				false,
				false,
				false,
				0,
				0,
				false,
				"",
				"",
				"",
			},
			nil,
		},
		{
			"case-invalid-total-mint",
			&BaseToken{
				"ABC Token",
				"abc",
				210000000000,
				testAddr,
				false,
				false,
				false,
				false,
				0,
				2100,
				false,
				"",
				"",
				"",
			},
			ErrTokenMintNotSupported("abc"),
		},
		{
			"case-invalid-total-burn",
			&BaseToken{
				"ABC Token",
				"abc",
				210000000000,
				testAddr,
				false,
				false,
				false,
				false,
				2100,
				0,
				false,
				"",
				"",
				"",
			},
			ErrTokenBurnNotSupported("abc"),
		},
		{
			"case-invalid-forbidden-state",
			&BaseToken{
				"ABC Token",
				"abc",
				210000000000,
				testAddr,
				false,
				false,
				false,
				false,
				0,
				0,
				true,
				"",
				"",
				"",
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
