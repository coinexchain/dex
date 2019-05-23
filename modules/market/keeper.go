package market

import (
	"github.com/coinexchain/dex/modules/incentive"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	marketIdetifierPrefix     = []byte{0x01}
	orderBookIdetifierPrefix  = []byte{0x02}
	orderQueueIdetifierPrefix = []byte{0x03}
	askListIdetifierPrefix    = []byte{0x04}
	bidListIdetifierPrefix    = []byte{0x05}
)

type Keeper struct {
	markeyKey sdk.StoreKey
	axk       ExpectedAssertStatusKeeper
	bnk       ExpectedBankxKeeper
	fek       incentive.FeeCollectionKeeper
}

func NewKeeper(key sdk.StoreKey, axkVal ExpectedAssertStatusKeeper, bnkVal ExpectedBankxKeeper, fekVal incentive.FeeCollectionKeeper) Keeper {
	return Keeper{markeyKey: key, axk: axkVal, bnk: bnkVal, fek: fekVal}
}
