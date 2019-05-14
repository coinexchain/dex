package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestCase struct {
	Valid bool
	Msg   sdk.Msg
}

func ValidateBasic(t *testing.T, cases []TestCase) {
	for _, tc := range cases {
		err := tc.Msg.ValidateBasic()
		if tc.Valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}
