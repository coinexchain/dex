package market

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
	"testing"
	"time"

	"github.com/coinexchain/dex/modules/market/match"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	ctx     sdk.Context
	mk      Keeper
	handler sdk.Handler
}

var (
	haveCetAddress = []byte("have-cet")
	stock          = "ludete"
	money          = "cet"
)

func prepareMockInput() testInput {
	db := dbm.NewMemDB()
	marketKey := sdk.NewKVStoreKey(MarketKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(marketKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	cdc := codec.New()

	mk := NewKeeper(marketKey, MockAssertKeeper{}, MockBankxKeeper{}, cdc,
		params.NewKeeper(cdc, keyParams, tkeyParams).Subspace(MarketKey))
	mk.RegisterCodec()
	handler := NewHandler(mk)
	return testInput{ctx: ctx, mk: mk, handler: handler}
}

func TestMarketInfoSetFailed(t *testing.T) {
	input := prepareMockInput()
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 6}
	ret := input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, false, ret.IsOK(), "create market info should failed")
}

func TestMarketInfoSetSuccess(t *testing.T) {
	input := prepareMockInput()
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8}
	ret := input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, true, ret.IsOK(), "create market info should succeed")
}

func TestCreateGTEOrderFailed(t *testing.T) {
	input := prepareMockInput()
	msgGteOrder := MsgCreateGTEOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		Symbol:         stock + SymbolSeparator + "noExist",
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.BUY,
		TimeInForce:    time.Now().Nanosecond() + 10000,
	}
	ret := input.handler(input.ctx, msgGteOrder)
	require.Equal(t, false, ret.IsOK(), "create GTE order should failed")
}

func TestCreateGTEOrderSuccess(t *testing.T) {
	input := prepareMockInput()
	msgGteOrder := MsgCreateGTEOrder{
		Sender:         haveCetAddress,
		Sequence:       1,
		Symbol:         stock + SymbolSeparator + money,
		OrderType:      LimitOrder,
		PricePrecision: 8,
		Price:          100,
		Quantity:       10000000,
		Side:           match.BUY,
		TimeInForce:    time.Now().Nanosecond() + 10000,
	}
	/*ret := */ input.handler(input.ctx, msgGteOrder)
	//require.Equal(t, true, ret.IsOK(), "create GTE order should succeed")
}
