package keepers

import (
	"github.com/coinexchain/dex/modules/bancorlite/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	BancorInfoKey    = []byte{0x10}
	BancorInfoKeyEnd = []byte{0x11}
)

type BancorInfo struct {
	Owner       sdk.AccAddress `json:"sender"`
	Token       string         `json:"token"`
	MaxSupply   sdk.Int        `json:"max_supply"`
	MaxPrice    sdk.Dec        `json:"max_price"`
	Price       sdk.Dec        `json:"price"`
	StockInPool sdk.Int        `json:"stock_in_pool"`
	MoneyInPool sdk.Int        `json:"money_in_pool"`
}

func (bi *BancorInfo) UpdateStockInPool(stockInPool sdk.Int) bool {
	if stockInPool.IsNegative() || stockInPool.GT(bi.MaxSupply) {
		return false
	}

	bi.StockInPool = stockInPool
	suppliedStock := bi.MaxSupply.Sub(bi.StockInPool)
	bi.Price = bi.MaxPrice.MulInt(suppliedStock).QuoInt(bi.MaxSupply)
	bi.MoneyInPool = bi.Price.MulInt(suppliedStock).QuoInt64(2).RoundInt()
	return true
}

func (bi *BancorInfo) IsConsistent() bool {
	if bi.StockInPool.IsNegative() || bi.StockInPool.GT(bi.MaxSupply) {
		return false
	}
	suppliedStock := bi.MaxSupply.Sub(bi.StockInPool)
	price := bi.MaxPrice.MulInt(suppliedStock).QuoInt(bi.MaxSupply)
	moneyInPool := price.MulInt(suppliedStock).QuoInt64(2).RoundInt()
	return price.Equal(bi.Price) && moneyInPool.Equal(bi.MoneyInPool)
}

type BancorInfoDisplay struct {
	Token       string `json:"token"`
	MaxSupply   string `json:"max_supply"`
	MaxPrice    string `json:"max_price"`
	StockInPool string `json:"stock_in_pool"`
	MoneyInPool string `json:"money_in_pool"`
}

func NewBancorInfoDisplay(bi *BancorInfo) BancorInfoDisplay {
	return BancorInfoDisplay{
		Token:       bi.Token,
		MaxSupply:   bi.MaxSupply.String(),
		MaxPrice:    bi.MaxPrice.String(),
		StockInPool: bi.StockInPool.String(),
		MoneyInPool: bi.MoneyInPool.String(),
	}
}

type BancorInfoKeeper struct {
	biKey sdk.StoreKey
	codec *codec.Codec
}

func NewBancorInfoKeeper(key sdk.StoreKey, cdc *codec.Codec) *BancorInfoKeeper {
	return &BancorInfoKeeper{
		biKey: key,
		codec: cdc,
	}
}

func (keeper *BancorInfoKeeper) Save(ctx sdk.Context, bi *BancorInfo) {
	store := ctx.KVStore(keeper.biKey)
	value := keeper.codec.MustMarshalBinaryBare(bi)
	key := append(BancorInfoKey, []byte(bi.Token)...)
	store.Set(key, value)
}

func (keeper *BancorInfoKeeper) Load(ctx sdk.Context, token string) *BancorInfo {
	store := ctx.KVStore(keeper.biKey)
	bi := &BancorInfo{}
	key := append(BancorInfoKey, []byte(token)...)
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
	Bik *BancorInfoKeeper
	Bxk types.ExpectedBankxKeeper
	Axk types.ExpectedAssetStatusKeeper
}

func NewKeeper(bik *BancorInfoKeeper,
	bxk types.ExpectedBankxKeeper,
	axk types.ExpectedAssetStatusKeeper) *Keeper {
	return &Keeper{
		Bik: bik,
		Bxk: bxk,
		Axk: axk,
	}
}
