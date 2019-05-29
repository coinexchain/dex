package init

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/app"
)

func TestAddGenesisAccount(t *testing.T) {
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	type args struct {
		appState     app.GenesisState
		addr         sdk.AccAddress
		coins        sdk.Coins
		vestingAmt   sdk.Coins
		vestingStart int64
		vestingEnd   int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid account",
			args{
				app.GenesisState{},
				addr1,
				sdk.NewCoins(),
				sdk.NewCoins(),
				0,
				0,
			},
			false,
		},
		{
			"dup account",
			args{
				app.GenesisState{Accounts: []app.GenesisAccount{{Address: addr1}}},
				addr1,
				sdk.NewCoins(),
				sdk.NewCoins(),
				0,
				0,
			},
			true,
		},
		{
			"invalid vesting amount",
			args{
				app.GenesisState{},
				addr1,
				sdk.NewCoins(sdk.NewInt64Coin("stake", 50)),
				sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				0,
				0,
			},
			true,
		},
		{
			"invalid vesting times",
			args{
				app.GenesisState{},
				addr1,
				sdk.NewCoins(sdk.NewInt64Coin("stake", 50)),
				sdk.NewCoins(sdk.NewInt64Coin("stake", 50)),
				1654668078,
				1554668078,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := addGenesisAccount(
				tt.args.appState,
				&accountInfo{tt.args.addr, tt.args.coins,
					tt.args.vestingAmt, tt.args.vestingStart, tt.args.vestingEnd},
			)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestAddGenesisAccountX(t *testing.T) {
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	type args struct {
		appState     app.GenesisState
		addr         sdk.AccAddress
		coins        sdk.Coins
		vestingAmt   sdk.Coins
		vestingStart int64
		vestingEnd   int64
	}

	testcase := []args{
		{
			app.GenesisState{},
			addr1,
			sdk.NewCoins(),
			sdk.NewCoins(),
			0,
			0,
		},
	}

	newstate, err := addGenesisAccount(testcase[0].appState, &accountInfo{testcase[0].addr, testcase[0].coins, testcase[0].vestingAmt, testcase[0].vestingStart, testcase[0].vestingEnd})

	require.Nil(t, err)
	require.Equal(t, 1, len(newstate.Accounts))
	require.Equal(t, false, newstate.Accounts[0].MemoRequired)
	require.Nil(t, newstate.Accounts[0].LockedCoins)
}
