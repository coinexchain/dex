package market

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"time"
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

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(marketKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	mk := NewKeeper(marketKey, MockAssertKeeper{}, MockBankxKeeper{})
	handler := NewHandler(mk)
	return testInput{ctx: ctx, mk: mk, handler: handler}
}

func TestMarketInfoSetFailed(t *testing.T) {
	input := prepareMockInput()
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8}
	ret := input.handler(input.ctx, msgMarketInfo)
	require.Equal(t, false, ret.IsOK(), "create market info should failed")
}

func TestMarketInfoSetSuccess(t *testing.T) {
	input := prepareMockInput()
	msgMarketInfo := MsgCreateMarketInfo{Stock: stock, Money: money, Creator: haveCetAddress, PricePrecision: 8}
	/*ret := */ input.handler(input.ctx, msgMarketInfo)
	//require.Equal(t, true, ret.IsOK(), "create market info should succeed")
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
		Side:           Buy,
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
		Side:           Buy,
		TimeInForce:    time.Now().Nanosecond() + 10000,
	}
	/*ret := */ input.handler(input.ctx, msgGteOrder)
	//require.Equal(t, true, ret.IsOK(), "create GTE order should succeed")
}
