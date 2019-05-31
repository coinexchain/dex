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
	)
	return
}

func newBaseGenesisAccount(address string, amt int64) app.GenesisAccount {
	return app.NewGenesisAccount(&auth.BaseAccount{
		Address: accAddressFromBech32(address),
		Coins:   dex.NewCetCoins(amt),
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