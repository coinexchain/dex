package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bankx Keeper will implement the interface
type ExpectedBankxKeeper interface {
	SendCoins(from, to sdk.AccAddress, amt sdk.Coins) error     // to tranfer coins
	FreezeCoins(acc sdk.AccAddress, amt sdk.Coins) error        // freeze some coins when creating orders
	UnfreezeCoins(acc sdk.AccAddress, amt sdk.Coins) error      // unfreeze coins and then orders can be executed
	IsFrozenByCoinOwner(acc sdk.AccAddress, denom string) error // The owner of the coin named "denom" has frozen acc to trade or transfer his coin.
}

// Asset Keeper will implement the interface
type ExpectedAssertStatusKeeper interface {
	IsFrozen(denom string) error // the coin's owner has frozen "denom", forbiding transmission and exchange.
	Exists(denom string) error   // check whether there is a coin named "denom"
	IsTokenIssuer(denom string, addr sdk.AccAddress) error
}
