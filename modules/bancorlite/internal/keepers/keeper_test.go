package keepers

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var owner = sdk.AccAddress("user")

func TestBancorInfo_UpdateStockInPool(t *testing.T) {
	type fields struct {
		Owner            sdk.AccAddress
		Stock            string
		Money            string
		InitPrice        sdk.Dec
		MaxSupply        sdk.Int
		MaxPrice         sdk.Dec
		Price            sdk.Dec
		StockInPool      sdk.Int
		MoneyInPool      sdk.Int
		EnableCancelTime int64
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
				Owner:            owner,
				Stock:            "bch",
				Money:            "cet",
				InitPrice:        sdk.NewDec(0),
				MaxSupply:        sdk.NewInt(100),
				MaxPrice:         sdk.NewDec(10),
				StockInPool:      sdk.NewInt(10),
				EnableCancelTime: 100,
			},
			args: args{
				stockInPool: sdk.NewInt(20),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BancorInfo{
				Owner:            tt.fields.Owner,
				Stock:            tt.fields.Stock,
				Money:            tt.fields.Money,
				InitPrice:        tt.fields.InitPrice,
				MaxSupply:        tt.fields.MaxSupply,
				MaxPrice:         tt.fields.MaxPrice,
				Price:            tt.fields.Price,
				StockInPool:      tt.fields.StockInPool,
				MoneyInPool:      tt.fields.MoneyInPool,
				EnableCancelTime: tt.fields.EnableCancelTime,
			}
			if got := bi.UpdateStockInPool(tt.args.stockInPool); got != tt.want {
				t.Errorf("BancorInfo.UpdateStockInPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBancorInfo_IsConsistent(t *testing.T) {
	type fields struct {
		Owner            sdk.AccAddress
		Stock            string
		Money            string
		InitPrice        sdk.Dec
		MaxSupply        sdk.Int
		MaxPrice         sdk.Dec
		Price            sdk.Dec
		StockInPool      sdk.Int
		MoneyInPool      sdk.Int
		EnableCancelTime int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "positive",
			fields: fields{
				Owner:            owner,
				Stock:            "bch",
				Money:            "cet",
				InitPrice:        sdk.NewDec(0),
				MaxSupply:        sdk.NewInt(100),
				MaxPrice:         sdk.NewDec(10),
				Price:            sdk.NewDec(1),
				StockInPool:      sdk.NewInt(90),
				MoneyInPool:      sdk.NewInt(5),
				EnableCancelTime: 100,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BancorInfo{
				Owner:            tt.fields.Owner,
				Stock:            tt.fields.Stock,
				Money:            tt.fields.Money,
				InitPrice:        tt.fields.InitPrice,
				MaxSupply:        tt.fields.MaxSupply,
				MaxPrice:         tt.fields.MaxPrice,
				Price:            tt.fields.Price,
				StockInPool:      tt.fields.StockInPool,
				MoneyInPool:      tt.fields.MoneyInPool,
				EnableCancelTime: tt.fields.EnableCancelTime,
			}
			if got := bi.IsConsistent(); got != tt.want {
				t.Errorf("BancorInfo.IsConsistent() = %v, want %v", got, tt.want)
			}
		})
	}
}
