package main

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/cet-sdk/modules/asset"
	"github.com/coinexchain/cet-sdk/modules/incentive"
	"github.com/coinexchain/cet-sdk/modules/market"
	"github.com/coinexchain/cet-sdk/testutil"
	dex "github.com/coinexchain/cet-sdk/types"
	"github.com/coinexchain/dex/app"
)

func createExampleGenesisState(cdc *codec.Codec) app.GenesisState {
	gsMap := app.ModuleBasics.DefaultGenesis()
	genState := app.FromMap(cdc, gsMap)
	genState.Accounts = createExampleGenesisAccounts()
	//genState.StakingData.Pool.NotBondedTokens = sdk.NewInt(588788547005740000)
	genState.AssetData = createExampleGenesisAssetData()
	genState.MarketData = createExampleGenesisMarketData()
	genState.GenUtil.GenTxs = append(genState.GenUtil.GenTxs, createExampleGenTx(cdc))
	return genState
}

func createExampleGenesisAccounts() (accs []genaccounts.GenesisAccount) {
	accs = append(accs,
		newBaseGenesisAccount(incentive.PoolAddr.String(), 31500000000000000),
		newBaseGenesisAccount("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke", 288800000000000000),
		newBaseGenesisAccount("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h", 88500000000000000),
		newVestingGenesisAccount("coinex1zvf0hx6rpz0n7dkuzu34s39dnsyr8eygqs8h3q", 36000000000000000, 1577836800),
		newVestingGenesisAccount("coinex1rfeae36tmm9t3gzacfq59hnv9j7fnaed3m4hhg", 36000000000000000, 1609459200),
		newVestingGenesisAccount("coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47", 36000000000000000, 1640995200),
		newVestingGenesisAccount("coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g", 36000000000000000, 1672531200),
		newVestingGenesisAccount("coinex1qyy6tvx7ymw44t4444sfmexpvczchr0tcp2p6p", 36000000000000000, 1704067200),
	)
	return
}

func createExampleGenesisAssetData() asset.GenesisState {
	cet := createCetToken("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	abc := createAbcToken()

	state := asset.DefaultGenesisState()
	state.Tokens = append(state.Tokens, cet, abc)
	return state
}

func createAbcToken() asset.Token {
	token := &asset.BaseToken{
		Name:             "ABC Chain Native Token",
		Symbol:           "abc",
		TotalSupply:      sdk.NewInt(588788547005740000),
		SendLock:         sdk.ZeroInt(),
		Owner:            accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
		Mintable:         false,
		Burnable:         true,
		AddrForbiddable:  false,
		TokenForbiddable: false,
		TotalBurn:        sdk.NewInt(411211452994260000),
		TotalMint:        sdk.ZeroInt(),
		IsForbidden:      false,
	}
	if err := token.Validate(); err != nil {
		panic(err)
	}
	return token
}

func createExampleGenesisMarketData() market.GenesisState {
	order0 := &market.Order{
		Sender:      accAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd"),
		Sequence:    100,
		TradingPair: "abc/cet",
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
		TradingPair: "btc/cet",
		OrderType:   2,
		Price:       sdk.NewDec(121920),
		Quantity:    100000,
		Side:        1,
		TimeInForce: 1002682839,
		Height:      100,
	}

	market0 := market.MarketInfo{
		Stock:             "abc",
		Money:             dex.CET,
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
	commissionMsg := staking.NewCommissionRates(rate, maxRate, maxChangeRate)

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
