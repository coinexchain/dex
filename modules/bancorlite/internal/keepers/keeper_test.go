package keepers_test

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/coinexchain/dex/modules/bancorlite/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/coinexchain/dex/modules/bancorlite/internal/keepers"
	"github.com/coinexchain/dex/testapp"
)

var (
	owner = sdk.AccAddress("user")
	bch   = "bch"
	cet   = "cet"
	abc   = "abc"
)

func defaultContext() (keepers.Keeper, sdk.Context) {
	app := testapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{}, false, log.NewNopLogger())
	return app.BancorKeeper, ctx
}
func TestBancorInfo_UpdateStockInPool(t *testing.T) {
	type fields struct {
		Owner              sdk.AccAddress
		Stock              string
		Money              string
		InitPrice          sdk.Dec
		MaxSupply          sdk.Int
		MaxPrice           sdk.Dec
		Price              sdk.Dec
		StockInPool        sdk.Int
		MoneyInPool        sdk.Int
		EarliestCancelTime int64
	}
	type args struct {
		stockInPool sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "positive",
			fields: fields{
				Owner:              owner,
				Stock:              bch,
				Money:              cet,
				InitPrice:          sdk.NewDec(0),
				MaxSupply:          sdk.NewInt(100),
				MaxPrice:           sdk.NewDec(10),
				StockInPool:        sdk.NewInt(10),
				EarliestCancelTime: 100,
			},
			args: args{
				stockInPool: sdk.NewInt(20),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &keepers.BancorInfo{
				Owner:              tt.fields.Owner,
				Stock:              tt.fields.Stock,
				Money:              tt.fields.Money,
				InitPrice:          tt.fields.InitPrice,
				MaxSupply:          tt.fields.MaxSupply,
				MaxPrice:           tt.fields.MaxPrice,
				MaxMoney:           sdk.ZeroInt(),
				Price:              tt.fields.Price,
				StockInPool:        tt.fields.StockInPool,
				MoneyInPool:        tt.fields.MoneyInPool,
				AR:                 0,
				EarliestCancelTime: tt.fields.EarliestCancelTime,
			}
			if got := bi.UpdateStockInPool(tt.args.stockInPool); got != tt.want {
				t.Errorf("BancorInfo.UpdateStockInPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBancorInfo_UpdateStockInPool2(t *testing.T) {

	rand.Seed(time.Now().Unix())
	random := func(min, max int64) int64 {
		return rand.Int63n(max-min) + min
	}
	maxSupply := int64(10000000000000)
	times := 1000000
	count := 0
	count2 := 0
	count3 := 0
	for i := 0; i < times; i++ {
		initPrice := random(0, 10)
		maxPrice := random(10, 20)
		randomAr := random(0, 5000)
		var tmpAr float64
		tmpAr = float64(randomAr) / 1000
		if randomAr == 0 {
			tmpAr = float64(randomAr) + 0.1
		}
		maxMoney := int64((float64(maxSupply*maxPrice) + tmpAr*float64(initPrice)*float64(maxSupply)) / (tmpAr + 1))
		supply := random(0, maxSupply)

		bi := keepers.BancorInfo{
			Owner:              nil,
			Stock:              "",
			Money:              "",
			InitPrice:          sdk.NewDec(initPrice),
			MaxSupply:          sdk.NewInt(maxSupply),
			StockPrecision:     0,
			MaxPrice:           sdk.NewDec(maxPrice),
			MaxMoney:           sdk.NewInt(maxMoney),
			AR:                 0,
			Price:              sdk.NewDec(initPrice),
			StockInPool:        sdk.NewInt(maxSupply),
			MoneyInPool:        sdk.NewInt(0),
			EarliestCancelTime: 0,
		}

		ar, money, _ := CalculateMoney(float64(supply), float64(maxSupply), float64(maxMoney), float64(initPrice), float64(maxPrice))
		bi.AR = int64(ar * types.ARSamples)
		bi.UpdateStockInPool(sdk.NewInt(maxSupply - supply))
		diffMoney := math.Abs(money - float64(bi.MoneyInPool.Int64()))
		if diffMoney/money > 0.000001 {
			//fmt.Printf("money is rough: ar:%f, AR in pool:%d, money diff ratio:%f, maxMoney:%d, maxSupply:%d, initPrice:%d, maxPrice:%d," +
			//	" supply:%d, money:%f, moneyInPool:%d, price:%f, priceInPool:%s" +
			//	"\n",
			//	ar, bi.AR, diffMoney/money, maxMoney, maxSupply, initPrice, maxPrice, supply, money, bi.MoneyInPool.Int64(),
			//	price, bi.Price.String(), )
			count++
		}
		if diffMoney/money > 0.00001 {
			count2++
		}
		if diffMoney/money > 0.0001 {
			count3++
		}
		//s := fmt.Sprintf("%f", price)
		//priceDec, _ := sdk.NewDecFromStr(s)
		//fmt.Printf("priceDec:%s\n", priceDec.String())
		//if priceDec.Sub(bi.Price).GT(sdk.NewDec(1)) || bi.Price.Sub(priceDec).GT(sdk.NewDec(1)) {
		//	fmt.Printf("price is rough: ar:%d, maxMoney:%d, maxSupply:%d, initPrice:%d, maxPrice:%d, supply:%d, price:%f, priceInPool:%s\n",
		//		bi.AR, maxMoney, maxSupply, initPrice, maxPrice, supply, price, bi.Price.String())
		//}
	}
	//fmt.Printf("percent of pass 0.000_001: %f\n", 1 - float64(count)/float64(times))
	//fmt.Printf("percent of pass 0.000_01: %f\n", 1 - float64(count2)/float64(times))
	//fmt.Printf("percent of pass 0.000_1: %f\n", 1 - float64(count3)/float64(times))
	require.True(t, float64(count3)/float64(times) < 0.05)
}

func CalculateMoney(supply, maxSupply, maxMoney, initPrice, maxPrice float64) (ar, money, price float64) {
	ar = (maxSupply*maxPrice - maxMoney) / (maxMoney - initPrice*maxSupply)
	money = math.Pow(supply/maxSupply, ar+1)*(maxMoney-initPrice*maxSupply) + initPrice*supply
	price = math.Pow(supply/maxSupply, ar)*(maxPrice-initPrice) + initPrice
	return
}

func TestBancorInfo_IsConsistent(t *testing.T) {
	type fields struct {
		Owner              sdk.AccAddress
		Stock              string
		Money              string
		InitPrice          sdk.Dec
		MaxSupply          sdk.Int
		MaxPrice           sdk.Dec
		Price              sdk.Dec
		StockInPool        sdk.Int
		MoneyInPool        sdk.Int
		EarliestCancelTime int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "positive",
			fields: fields{
				Owner:              owner,
				Stock:              bch,
				Money:              cet,
				InitPrice:          sdk.NewDec(0),
				MaxSupply:          sdk.NewInt(100),
				MaxPrice:           sdk.NewDec(10),
				Price:              sdk.NewDec(1),
				StockInPool:        sdk.NewInt(90),
				MoneyInPool:        sdk.NewInt(5),
				EarliestCancelTime: 100,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &keepers.BancorInfo{
				Owner:              tt.fields.Owner,
				Stock:              tt.fields.Stock,
				Money:              tt.fields.Money,
				InitPrice:          tt.fields.InitPrice,
				MaxSupply:          tt.fields.MaxSupply,
				MaxPrice:           tt.fields.MaxPrice,
				MaxMoney:           sdk.ZeroInt(),
				Price:              tt.fields.Price,
				StockInPool:        tt.fields.StockInPool,
				MoneyInPool:        tt.fields.MoneyInPool,
				EarliestCancelTime: tt.fields.EarliestCancelTime,
			}
			if got := bi.IsConsistent(); got != tt.want {
				t.Errorf("BancorInfo.IsConsistent() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestBancorInfoKeeper(t *testing.T) {
	keeper, ctx := defaultContext()
	bi := []keepers.BancorInfo{
		{
			Owner:              owner,
			Stock:              bch,
			Money:              cet,
			InitPrice:          sdk.NewDec(0),
			MaxSupply:          sdk.NewInt(100),
			MaxPrice:           sdk.NewDec(10),
			MaxMoney:           sdk.ZeroInt(),
			AR:                 1,
			Price:              sdk.NewDec(1),
			StockInPool:        sdk.NewInt(90),
			MoneyInPool:        sdk.NewInt(5),
			EarliestCancelTime: 100,
		},
		{
			Owner:              owner,
			Stock:              abc,
			Money:              cet,
			InitPrice:          sdk.NewDec(0),
			MaxSupply:          sdk.NewInt(100),
			MaxPrice:           sdk.NewDec(10),
			MaxMoney:           sdk.ZeroInt(),
			AR:                 1,
			Price:              sdk.NewDec(1),
			StockInPool:        sdk.NewInt(90),
			MoneyInPool:        sdk.NewInt(5),
			EarliestCancelTime: 0,
		},
	}
	for _, p := range bi {
		keeper.Save(ctx, &p)
	}

	for i := range bi {
		loadBI := keeper.Load(ctx, bi[i].GetSymbol())
		require.True(t, reflect.DeepEqual(*loadBI, bi[i]))
	}

	keeper.Remove(ctx, &bi[0])
	require.Nil(t, keeper.Load(ctx, bi[0].GetSymbol()))
	require.False(t, keeper.IsBancorExist(ctx, bch))
	require.True(t, keeper.IsBancorExist(ctx, abc))
}
