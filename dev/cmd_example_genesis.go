package dev

import (
	"fmt"

	"github.com/spf13/cobra"
	tm "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
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
	genState.AssetData = createGenesisAssetData()
	genState.MarketData = createGenesisMarketData()
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

func createGenesisAssetData() asset.GenesisState {
	t0 := &asset.BaseToken{
		Name:             "CoinEx Chain Native Token",
		Symbol:           "cet",
		TotalSupply:      588788547005740000,
		Owner:            accAddressFromBech32("cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v"),
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  false,
		TokenForbiddable: false,
		TotalBurn:        411211452994260000,
		TotalMint:        0,
		IsForbidden:      false,
	}
	t1 := &asset.BaseToken{
		Name:             "ABC Chain Native Token",
		Symbol:           "abc",
		TotalSupply:      588788547005740000,
		Owner:            accAddressFromBech32("cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v"),
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  false,
		TokenForbiddable: false,
		TotalBurn:        411211452994260000,
		TotalMint:        0,
		IsForbidden:      false,
	}

	state := asset.DefaultGenesisState()
	state.Tokens = append(state.Tokens, t0, t1)
	return state
}

func createGenesisMarketData() market.GenesisState {
	order0 := &market.Order{
		Sender:      accAddressFromBech32("cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v"),
		Sequence:    100,
		Symbol:      "abc/cet",
		OrderType:   2,
		Price:       sdk.NewDec(100),
		Quantity:    100000,
		Side:        1,
		TimeInForce: 10092839,
		Height:      100,
	}
	order1 := &market.Order{
		Sender:      accAddressFromBech32("cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v"),
		Sequence:    170,
		Symbol:      "btc/cet",
		OrderType:   2,
		Price:       sdk.NewDec(121920),
		Quantity:    100000,
		Side:        1,
		TimeInForce: 1002682839,
		Height:      100,
	}

	market0 := market.MarketInfo{
		Stock:             "abc",
		Money:             "cet",
		Creator:           accAddressFromBech32("cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v"),
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(8568),
	}

	state := market.DefaultGenesisState()
	state.Orders = append(state.Orders, order0, order1)
	state.MarketInfos = append(state.MarketInfos, market0)

	return state
}
