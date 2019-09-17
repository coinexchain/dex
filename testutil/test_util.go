package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	dex "github.com/coinexchain/dex/types"
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

func KeyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func ToAccAddress(addr string) sdk.AccAddress {
	return sdk.AccAddress([]byte(addr))
}

func NewBaseAccount(cet int64, accNum, seq uint64) (crypto.PrivKey, auth.BaseAccount) {
	key, pubKey, addr := KeyPubAddr()
	acc := auth.BaseAccount{
		Address:       addr,
		Coins:         dex.NewCetCoins(cet),
		PubKey:        pubKey,
		AccountNumber: accNum,
		Sequence:      seq,
	}
	return key, acc
}
