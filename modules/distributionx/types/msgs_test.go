package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/testutil"
	"github.com/coinexchain/dex/types"
)

var validCoins = types.NewCetCoins(10e8)

func TestMain(m *testing.M) {
	types.InitSdkConfig()
	os.Exit(m.Run())
}

func TestDonateToCommunityPoolRoute(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgDonateToCommunityPool(addr, types.NewCetCoins(1e8))
	require.Equal(t, Route(), "distrx")
	require.Equal(t, Type(), "donate_to_community_pool")
}

func TestDonateToCommunityPoolValidation(t *testing.T) {
	validAddr := sdk.AccAddress([]byte("addr"))
	var emptyAddr sdk.AccAddress

	var invalidDenomCoins = sdk.NewCoins(sdk.NewCoin("abc", sdk.NewInt(1e8)))
	var invalidLenCoins = append(validCoins, types.NewCetCoin(1e8))
	var invalidAmount = sdk.NewCoins(sdk.NewCoin("cet", sdk.NewInt(10)))
	invalidAmount[0].Amount = sdk.ZeroInt()

	testutil.ValidateBasic(t, []testutil.TestCase{
		{Valid: true, Msg: NewMsgDonateToCommunityPool(validAddr, validCoins)},
		{Valid: false, Msg: NewMsgDonateToCommunityPool(validAddr, invalidDenomCoins)},
		{Valid: false, Msg: NewMsgDonateToCommunityPool(validAddr, invalidLenCoins)},
		{Valid: false, Msg: NewMsgDonateToCommunityPool(validAddr, invalidAmount)},
		{Valid: false, Msg: NewMsgDonateToCommunityPool(emptyAddr, validCoins)},
	})
}

func TestDonateToCommunityPoolGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("addr")))
	msg := NewMsgDonateToCommunityPool(addr, validCoins)
	sign := GetSignBytes()

	expected := `{"type":"distrx/MsgDonateToCommunityPool","value":{"amount":[{"amount":"1000000000","denom":"cet"}],"from_addr":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"}}`
	require.Equal(t, expected, string(sign))
}

func TestDonateToCommunityPoolGetSigners(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgDonateToCommunityPool(addr, validCoins)
	signers := GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, addr, signers[0])
}
