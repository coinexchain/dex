package dev

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	tm "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/incentive"
	dex "github.com/coinexchain/dex/types"
)

func ExampleGenesisCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example-genesis",
		Short: "Print example genesis JSON",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printExampleGenesis(cdc)
		},
	}
	return cmd
}

func printExampleGenesis(cdc *codec.Codec) error {
	genState := createExampleGenesisState()
	gneStateBytes, err := codec.MarshalJSONIndent(cdc, genState)
	if err != nil {
		return err
	}

	genDoc := tm.GenesisDoc{
		ChainID:    "coinexdex",
		Validators: nil,
		AppState:   gneStateBytes,
	}
	genDoc.ValidateAndComplete()

	genDocBytes, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(genDocBytes))
	return nil
}

func createExampleGenesisState() app.GenesisState {
	genState := app.NewDefaultGenesisState()
	genState.Accounts = createGenesisAccounts()
	return genState
}

func createGenesisAccounts() (accs []app.GenesisAccount) {
	accs = append(accs,
		newBaseGenesisAccount(incentive.IncentiveCoinsAccAddr.String(), 30000000000000000),
		newBaseGenesisAccount("cosmos1c79cqwzah604v0pqg0h88g99p5zg08hgf0cspy", 258788547005740000),
		newBaseGenesisAccount("cosmos1n3n5w8mqjf339xse0rwvl0u7nqgp8e5d0nwt20", 120000000000000000),
		newVestingGenesisAccount("cosmos1xtpex9x7yq8n9d7f8dpgu5mfajrv2thvr6u34q", 36000000000000000, 1577836800),
		newVestingGenesisAccount("cosmos1966f22al7r23h3melq8yt8tnglhweunrxkcezl", 36000000000000000, 1609459200),
		newVestingGenesisAccount("cosmos12kt3yq0kdvu3zm0pq65dkd83hy3j9wgd2m9hfv", 36000000000000000, 1640995200),
		newVestingGenesisAccount("cosmos1r0z8lf82euwlxx0fuvny3jfl0jj2tmdxwuutxj", 36000000000000000, 1672531200),
		newVestingGenesisAccount("cosmos1wezn7xuu5ha39t089mwfeypx0rxvxsutnr0h9p", 36000000000000000, 1704067200),
	)
	return
}

func newBaseGenesisAccount(address string, amt int64) app.GenesisAccount {
	return app.NewGenesisAccount(&auth.BaseAccount{
		Address: accAddressFromBech32(address),
		Coins:   dex.NewCetCoins(amt),
	})
}

func newVestingGenesisAccount(address string, amt int64, endTime int64) app.GenesisAccount {
	return app.NewGenesisAccountI(&auth.DelayedVestingAccount{
		BaseVestingAccount: &auth.BaseVestingAccount{
			BaseAccount: &auth.BaseAccount{
				Address: accAddressFromBech32(address),
				Coins:   dex.NewCetCoins(amt),
			},
			OriginalVesting: dex.NewCetCoins(amt),
			EndTime:         endTime,
		},
	})
}

func accAddressFromBech32(address string) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return addr
}

func randomAccAddress() sdk.AccAddress {
	return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
}
