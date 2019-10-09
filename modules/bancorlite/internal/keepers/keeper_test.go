package keepers_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

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
				Price:              tt.fields.Price,
				StockInPool:        tt.fields.StockInPool,
				MoneyInPool:        tt.fields.MoneyInPool,
				EarliestCancelTime: tt.fields.EarliestCancelTime,
			}
			if got := bi.UpdateStockInPool(tt.args.stockInPool); got != tt.want {
				t.Errorf("BancorInfo.UpdateStockInPool() = %v, want %v", got, tt.want)
			}
		})
	}
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
			Price:              sdk.NewDec(1),
			StockInPool:        sdk.NewInt(90),
			MoneyInPool:        sdk.NewInt(5),
			EarliestCancelTime: 0,
		},
	}
	for _, p := range bi {
		keeper.Bik.Save(ctx, &p)
	}

	for i := range bi {
		loadBI := keeper.Bik.Load(ctx, bi[i].GetSymbol())
		require.True(t, reflect.DeepEqual(*loadBI, bi[i]))
	}

	keeper.Bik.Remove(ctx, &bi[0])
	require.Nil(t, keeper.Bik.Load(ctx, bi[0].GetSymbol()))
	require.False(t, keeper.IsBancorExist(ctx, bch))
	require.True(t, keeper.IsBancorExist(ctx, abc))
}
