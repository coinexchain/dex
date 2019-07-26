package keepers

import (
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"time"
)

var (
	BancorInfoKey    = []byte{0x10}
	BancorInfoKeyEnd = []byte{0x11}
)

const SymbolSeparator = "/"

type BancorInfo struct {
	Owner            sdk.AccAddress `json:"sender"`
	Stock            string         `json:"stock"`
	Money            string         `json:"money"`
	InitPrice        sdk.Dec        `json:"init_price"`
	MaxSupply        sdk.Int        `json:"max_supply"`
	MaxPrice         sdk.Dec        `json:"max_price"`
	Price            sdk.Dec        `json:"price"`
	StockInPool      sdk.Int        `json:"stock_in_pool"`
	MoneyInPool      sdk.Int        `json:"money_in_pool"`
	EnableCancelTime int64          `json:"enable_cancel_time"`
}

func (bi *BancorInfo) UpdateStockInPool(stockInPool sdk.Int) bool {
	if stockInPool.IsNegative() || stockInPool.GT(bi.MaxSupply) {
		return false
	}

	bi.StockInPool = stockInPool
	suppliedStock := bi.MaxSupply.Sub(bi.StockInPool)
	bi.Price = bi.MaxPrice.Sub(bi.InitPrice).MulInt(suppliedStock).QuoInt(bi.MaxSupply).Add(bi.InitPrice)
	bi.MoneyInPool = bi.Price.Add(bi.InitPrice).MulInt(suppliedStock).QuoInt64(2).RoundInt()
	return true
}

func (bi *BancorInfo) IsConsistent() bool {
	if bi.StockInPool.IsNegative() || bi.StockInPool.GT(bi.MaxSupply) {
		return false
	}
	suppliedStock := bi.MaxSupply.Sub(bi.StockInPool)
	price := bi.MaxPrice.Sub(bi.InitPrice).MulInt(suppliedStock).QuoInt(bi.MaxSupply).Add(bi.InitPrice)
	moneyInPool := price.Add(bi.InitPrice).MulInt(suppliedStock).QuoInt64(2).RoundInt()
	return price.Equal(bi.Price) && moneyInPool.Equal(bi.MoneyInPool)
}

type BancorInfoDisplay struct {
	Stock            string `json:"stock"`
	Money            string `json:"money"`
	InitPrice        string `json:"init_price"`
	MaxSupply        string `json:"max_supply"`
	MaxPrice         string `json:"max_price"`
	StockInPool      string `json:"stock_in_pool"`
	MoneyInPool      string `json:"money_in_pool"`
	EnableCancelTime string `json:"enable_cancel_time"`
}

func NewBancorInfoDisplay(bi *BancorInfo) BancorInfoDisplay {
	return BancorInfoDisplay{
		Stock:            bi.Stock,
		Money:            bi.Money,
		InitPrice:        bi.InitPrice.String(),
		MaxSupply:        bi.MaxSupply.String(),
		MaxPrice:         bi.MaxPrice.String(),
		StockInPool:      bi.StockInPool.String(),
		MoneyInPool:      bi.MoneyInPool.String(),
		EnableCancelTime: time.Unix(bi.EnableCancelTime, 0).Format(time.RFC3339),
	}
}

type BancorInfoKeeper struct {
	biKey         sdk.StoreKey
	codec         *codec.Codec
	paramSubspace params.Subspace
}

func NewBancorInfoKeeper(key sdk.StoreKey, cdc *codec.Codec, paramSubspace params.Subspace) *BancorInfoKeeper {
	return &BancorInfoKeeper{
		biKey:         key,
		codec:         cdc,
		paramSubspace: paramSubspace.WithKeyTable(types.ParamKeyTable()),
	}
}

func (keeper *BancorInfoKeeper) SetParam(ctx sdk.Context, params types.Params) {
	keeper.paramSubspace.SetParamSet(ctx, &params)
}

func (keeper *BancorInfoKeeper) GetParam(ctx sdk.Context) (param types.Params) {
	keeper.paramSubspace.GetParamSet(ctx, &param)
	return
}

func (keeper *BancorInfoKeeper) Save(ctx sdk.Context, bi *BancorInfo) {
	store := ctx.KVStore(keeper.biKey)
	value := keeper.codec.MustMarshalBinaryBare(bi)
	key := append(BancorInfoKey, []byte(bi.Stock+SymbolSeparator+bi.Money)...)
	store.Set(key, value)
}

func (keeper *BancorInfoKeeper) Remove(ctx sdk.Context, bi *BancorInfo) {
	store := ctx.KVStore(keeper.biKey)
	key := append(BancorInfoKey, []byte(bi.Stock+SymbolSeparator+bi.Money)...)
	value := store.Get(key)
	if value != nil {
		store.Delete(key)
	}
}

func (keeper *BancorInfoKeeper) Load(ctx sdk.Context, symbol string) *BancorInfo {
	store := ctx.KVStore(keeper.biKey)
	bi := &BancorInfo{}
	key := append(BancorInfoKey, []byte(symbol)...)
	biBytes := store.Get(key)
	if biBytes == nil {
		return nil
	}
	keeper.codec.MustUnmarshalBinaryBare(biBytes, bi)
	return bi
}

func (keeper *BancorInfoKeeper) Iterate(ctx sdk.Context, biProc func(bi *BancorInfo)) {
	store := ctx.KVStore(keeper.biKey)
	iter := store.Iterator(BancorInfoKey, BancorInfoKeyEnd)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		bi := &BancorInfo{}
		keeper.codec.MustUnmarshalBinaryBare(iter.Value(), bi)
		biProc(bi)
	}
}

type Keeper struct {
	Bik         *BancorInfoKeeper
	Bxk         types.ExpectedBankxKeeper
	Axk         types.ExpectedAssetStatusKeeper
	Mk          types.ExpectedMarketKeeper
	MsgProducer msgqueue.Producer
}

func NewKeeper(bik *BancorInfoKeeper,
	bxk types.ExpectedBankxKeeper,
	axk types.ExpectedAssetStatusKeeper,
	mk types.ExpectedMarketKeeper,
	mq msgqueue.Producer) Keeper {
	return Keeper{
		Bik:         bik,
		Bxk:         bxk,
		Axk:         axk,
		Mk:          mk,
		MsgProducer: mq,
	}
}

func (k Keeper) IsBancorExist(ctx sdk.Context, stock string) bool {
	store := ctx.KVStore(k.Bik.biKey)
	key := append(BancorInfoKey, []byte(stock+SymbolSeparator)...)
	iter := store.Iterator(key, append(key, 0xff))
	defer iter.Close()
	iter.Domain()
	for iter.Valid() {
		return true
	}
	return false
}

//kafka msg
type MsgBancorCreateForKafka struct {
	Owner            sdk.AccAddress `json:"owner"`
	Stock            string         `json:"stock"`
	Money            string         `json:"money"`
	InitPrice        sdk.Dec        `json:"init_price"`
	MaxSupply        sdk.Int        `json:"max_supply"`
	MaxPrice         sdk.Dec        `json:"max_price"`
	EnableCancelTime int64          `json:"enable_cancel_time"`
	BlockHeight      int64          `json:"block_height"`
}

type MsgBancorInfoForKafka struct {
	Owner            sdk.AccAddress `json:"sender"`
	Stock            string         `json:"stock"`
	Money            string         `json:"money"`
	InitPrice        sdk.Dec        `json:"init_price"`
	MaxSupply        sdk.Int        `json:"max_supply"`
	MaxPrice         sdk.Dec        `json:"max_price"`
	Price            sdk.Dec        `json:"price"`
	StockInPool      sdk.Int        `json:"stock_in_pool"`
	MoneyInPool      sdk.Int        `json:"money_in_pool"`
	EnableCancelTime int64          `json:"enable_cancel_time"`
	BlockHeight      int64          `json:"block_height"`
}

type MsgBancorTradeInfoForKafka struct {
	Sender      sdk.AccAddress `json:"sender"`
	Stock       string         `json:"stock"`
	Money       string         `json:"money"`
	Amount      int64          `json:"amount"`
	Side        string         `json:"side"`
	MoneyLimit  int64          `json:"money_limit"`
	TxPrice     sdk.Dec        `json:"transaction_price"`
	BlockHeight int64          `json:"block_height"`
}

type MsgBancorCancelForKafka struct {
	Owner       sdk.AccAddress `json:"owner"`
	Stock       string         `json:"stock"`
	Money       string         `json:"money"`
	BlockHeight int64          `json:"block_height"`
}
