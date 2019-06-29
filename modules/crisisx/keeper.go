package crisisx

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/crisis"
)

const (
	ModuleName = "crisisx"
)

type Keeper struct {
	ck   crisis.Keeper
	bk   ExpectBankxKeeper
	feek auth.FeeCollectionKeeper
}

func NewKeeper(ckVal crisis.Keeper) Keeper {
	return Keeper{
		ck: ckVal,
	}
}
