package authx

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type testInput struct {
	ctx sdk.Context
	axk AccountXKeeper
	ak  auth.AccountKeeper
	cdc *codec.Codec
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()
	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authXKey := sdk.NewKVStoreKey("authXKey")
	authKey := sdk.NewKVStoreKey("authKey")
	skey := sdk.NewKVStoreKey("params")
	tkey := sdk.NewTransientStoreKey("transient_params")
	paramsKeeper := params.NewKeeper(cdc, skey, tkey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authXKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	axk := NewKeeper(cdc, authXKey, paramsKeeper.Subspace(bank.DefaultParamspace))
	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return testInput{ctx: ctx, axk: axk, ak: ak, cdc: cdc}
}

func TestGetSetParams(t *testing.T) {
	input := setupTestInput()
	params := DefaultParams()
	input.axk.SetParams(input.ctx, params)
	params2 := input.axk.GetParams(input.ctx)
	require.True(t, params.Equal(params2))
}

func TestAccountXGetSet(t *testing.T) {
	input := setupTestInput()
	addr := sdk.AccAddress([]byte("some-address"))

	_, ok := input.axk.GetAccountX(input.ctx, addr)
	require.False(t, ok)

	//create account
	acc := NewAccountXWithAddress(addr)
	require.Equal(t, addr, acc.Address)

	input.axk.SetAccountX(input.ctx, acc)

	acc, ok = input.axk.GetAccountX(input.ctx, addr)
	require.True(t, ok)

	acc.MemoRequired = false
	input.axk.SetAccountX(input.ctx, acc)
	acc, _ = input.axk.GetAccountX(input.ctx, addr)
	require.Equal(t, false, acc.MemoRequired)

	lockedcoin := acc.LockedCoins
	require.Nil(t, lockedcoin)
}

func TestAddressStoreKey(t *testing.T) {
	addr := sdk.AccAddress([]byte("some-address1"))
	addrStoreKey := AddressStoreKey(addr)
	expectedOutput := []byte{0x1, 0x73, 0x6f, 0x6d, 0x65, 0x2d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x31}
	require.Equal(t, expectedOutput, addrStoreKey)
}
