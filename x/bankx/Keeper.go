package bankx

import (
	"github.com/coinexchain/dex/x/authx"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct{
	axk authx.AccountXKeeper
	bk bank.BaseKeeper
	fck auth.FeeCollectionKeeper

}

func NewKeeper(axk authx.AccountXKeeper, bk bank.BaseKeeper, fck auth.FeeCollectionKeeper) Keeper{
	return Keeper{
		axk:axk,
		bk:bk,
		fck:fck,
	}
}