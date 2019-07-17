package types

import (
	dex "github.com/coinexchain/dex/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParams_Equal(t *testing.T) {
	p1 := DefaultParams()
	p2 := DefaultParams()
	require.Equal(t, p1, p2)
	require.True(t, p1.Equal(p2))

	// mount should equal
	cet := dex.NewCetCoins(10)
	p1.IssueTokenFee = cet
	require.NotEqual(t, p1, p2)

	// denom should equal
	abc := NewTokenCoins("abc", 1E12)
	p1.IssueTokenFee = abc
	require.NotEqual(t, p1, p2)
}
