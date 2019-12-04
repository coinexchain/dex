package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrMsg(t *testing.T) {
	err := ErrInvalidSymbol()
	require.Equal(t, CodeInvalidSymbol, err.Code())
	err = ErrBancorAlreadyExists()
	require.Equal(t, CodeBancorAlreadyExists, err.Code())
	err = ErrNonOwnerIsProhibited()
	require.Equal(t, CodeNonOwnerIsProhibited, err.Code())
	err = ErrNoSuchToken()
	require.Equal(t, CodeNoSuchToken, err.Code())
	err = ErrNotBancorOwner()
	require.Equal(t, CodeNotBancorOwner, err.Code())
	err = ErrNegativePrice()
	require.Equal(t, CodeNegativeInitPrice, err.Code())
	err = ErrEarliestCancelTimeIsNegative()
	require.Equal(t, CodeCancelEnableTimeNegative, err.Code())
	err = ErrEarliestCancelTimeNotArrive()
	require.Equal(t, CodeCancelTimeNotArrived, err.Code())

}
