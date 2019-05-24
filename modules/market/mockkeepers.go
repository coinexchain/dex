package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MockBankxKeeper struct {
}

func (mb MockBankxKeeper) SendCoins(from, to sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (mb MockBankxKeeper) FreezeCoins(acc sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (mb MockBankxKeeper) UnfreezeCoins(acc sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (mb MockBankxKeeper) DeductFeeFromAddressAndCollectFeetoIncentive(acc sdk.AccAddress, coins sdk.Coins) error {
	return nil
}

func (mb MockBankxKeeper) HaveSufficientCoins(addr sdk.AccAddress, amt sdk.Coins) bool {
	return true
}

//-----------------------------------------------------------

type MockAssertKeeper struct {
}

func (ma MockAssertKeeper) IsTokenFrozen(addr sdk.AccAddress, denom string) bool {
	return true
}

func (ma MockAssertKeeper) IsTokenExists(denom string) bool {
	return true
}

func (ma MockAssertKeeper) IsTokenIssuer(denom string, addr sdk.AccAddress) bool {
	return true
}
