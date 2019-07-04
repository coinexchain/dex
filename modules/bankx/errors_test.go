package bankx

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestErrMsg(t *testing.T) {
	err := ErrMemoMissing()
	require.Equal(t, CodeMemoMissing, err.Code())
	err = ErrorInsufficientCETForActivatingFee()
	require.Equal(t, CodeInsufficientCETForActivationFee, err.Code())
	err = ErrUnlockTime("")
	require.Equal(t, CodeInvalidUnlockTime, err.Code())
	err = ErrTokenForbiddenByOwner("")
	require.Equal(t, CodeTokenForbiddenByOwner, err.Code())
}
