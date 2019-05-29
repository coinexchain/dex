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

func (mb MockBankxKeeper) HaveSufficientCoins(addr sdk.AccAddress, amt sdk.Coins) bool {
	return true
}

//-----------------------------------------------------------

type MockAssertKeeper struct {
}

func (ma MockAssertKeeper) IsTokenFrozen(ctx sdk.Context, denom string) bool {
	return true
}

func (ma MockAssertKeeper) IsTokenExists(ctx sdk.Context, denom string) bool {
	return true
}

func (ma MockAssertKeeper) IsTokenIssuer(ctx sdk.Context, denom string, addr sdk.AccAddress) bool {
	return true
}
