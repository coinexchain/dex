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

func (mb MockBankxKeeper) IsFrozenByCoinOwner(acc sdk.AccAddress, denom string) error {
	return nil
}

func (mb MockBankxKeeper) HaveSufficientCoins(addr sdk.AccAddress, amt sdk.Coins) bool {
	return true
}

//-----------------------------------------------------------

type MockAssertKeeper struct {
}

func (ma MockAssertKeeper) IsFrozen(denom string) error {
	return nil
}

func (ma MockAssertKeeper) Exists(denom string) error {
	return nil
}

func (ma MockAssertKeeper) IsTokenIssuer(denom string, addr sdk.AccAddress) error {
	return nil
}
