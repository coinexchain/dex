package dev

import (
	"encoding/json"
	"fmt"

	"github.com/coinexchain/dex/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking"

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
	genState := createExampleGenesisState(cdc)
	gneStateBytes, err := codec.MarshalJSONIndent(cdc, genState)
	if err != nil {
		return err
	}

	genDoc := tm.GenesisDoc{
		ChainID:    "coinexdex",
		Validators: nil,
		AppState:   gneStateBytes,
	}
	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	genDocBytes, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(genDocBytes))
	return nil
}

func createExampleGenesisState(cdc *codec.Codec) app.GenesisState {
	genState := app.NewDefaultGenesisState()
	genState.Accounts = createGenesisAccounts()
	genState.StakingData.Pool.NotBondedTokens = sdk.NewInt(588788547005740000)
	genState.AssetData = createGenesisAssetData()
	genState.MarketData = createGenesisMarketData()
	genState.GenTxs = append(genState.GenTxs, createExampleGenTx(cdc))
	return genState
}

func createGenesisAccounts() (accs []app.GenesisAccount) {
	accs = append(accs,
		newBaseGenesisAccount(incentive.IncentivePoolAddr.String(), 30000000000000000),
		newBaseGenesisAccount("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke", 288788547005740000),
		newBaseGenesisAccount("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h", 90000000000000000),
		newVestingGenesisAccount("coinex1zvf0hx6rpz0n7dkuzu34s39dnsyr8eygqs8h3q", 36000000000000000, 1577836800),
		newVestingGenesisAccount("coinex1rfeae36tmm9t3gzacfq59hnv9j7fnaed3m4hhg", 36000000000000000, 1609459200),
		newVestingGenesisAccount("coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47", 36000000000000000, 1640995200),
		newVestingGenesisAccount("coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g", 36000000000000000, 1672531200),
		newVestingGenesisAccount("coinex1qyy6tvx7ymw44t4444sfmexpvczchr0tcp2p6p", 36000000000000000, 1704067200),
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
		Owner:            accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
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
		Owner:            accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
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
		Sender:      accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
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
		Sender:      accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
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
		Creator:           accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
		PricePrecision:    8,
		LastExecutedPrice: sdk.NewDec(8568),
	}

	state := market.DefaultGenesisState()
	state.Orders = append(state.Orders, order0, order1)
	state.MarketInfos = append(state.MarketInfos, market0)

	return state
}

func createExampleGenTx(cdc *codec.Codec) json.RawMessage {
	key, pk, addr := testutil.KeyPubAddr()

	amount := dex.NewCetCoin(10000000000000000)
	description := staking.NewDescription("node0", "node0", "http://node0.coinexchain.org", "")

	rate, _ := sdk.NewDecFromStr("0.1")
	maxRate, _ := sdk.NewDecFromStr("0.2")
	maxChangeRate, _ := sdk.NewDecFromStr("0.01")
	commissionMsg := staking.NewCommissionMsg(rate, maxRate, maxChangeRate)

	minSelfDelegation := sdk.NewInt(10000000000000000)

	msg := staking.NewMsgCreateValidator(
		sdk.ValAddress(addr), pk, amount, description, commissionMsg, minSelfDelegation,
	)

	stdTx := testutil.NewStdTxBuilder("coinexdex").
		Msgs(msg).
		AccNumSeqKey(0, 0, key).
		GasAndFee(200000, 10).
		Build()

	txBytes, err := codec.MarshalJSONIndent(cdc, stdTx)
	if err != nil {
		panic(err)
	}
	return txBytes
}