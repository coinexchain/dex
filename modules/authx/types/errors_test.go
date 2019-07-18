package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestErrMsg(t *testing.T) {
	err1 := ErrInvalidMinGasPriceLimit(sdk.NewDec(100))
	require.True(t, strings.Contains(err1.Error(), "invalid minimum gas price limit: 100.0"))

	err2 := ErrGasPriceTooLow(sdk.NewDec(100), sdk.NewDec(60))
	require.True(t, strings.Contains(err2.Error(), "gas price too low: 60.000000000000000000 < 100.000000000000000000"))
}
