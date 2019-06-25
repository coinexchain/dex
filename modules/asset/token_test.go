package asset

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
				tAccAddr,
				false,
				false,
				false,
				false,
				0,
				0,
				false,
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
				tAccAddr,
				false,
				false,
				false,
				false,
				0,
				-1,
				false,
				"",
				"",
			},
			ErrorInvalidTokenMint("Invalid total mint: -1"),
		},
		{
			"case-invalid-total-burn",
			&BaseToken{
				"ABC Token",
				"abc",
				210000000000,
				tAccAddr,
				false,
				false,
				false,
				false,
				9E18 + 1,
				0,
				false,
				"",
				"",
			},
			ErrorInvalidTokenBurn("Invalid total burn: 9000000000000000001"),
		},
		{
			"case-invalid-forbidden-state",
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
				true,
				"",
				"",
			},
			ErrorInvalidTokenForbidden("Invalid Forbidden state"),
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
