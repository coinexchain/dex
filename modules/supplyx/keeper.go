package supplyx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

//wrap supply keeper to modify BurnCoins func
type Keeper struct {
	sk supply.Keeper
	dk distribution.Keeper
}

func NewKeeper(sk supply.Keeper, dk distribution.Keeper) Keeper {
	return Keeper{
		sk: sk,
		dk: dk,
	}
}
func (k Keeper) GetSupply(ctx sdk.Context) supply.Supply {
	return k.sk.GetSupply(ctx)
}

func (k Keeper) GetModuleAddress(name string) sdk.AccAddress {
	return k.sk.GetModuleAddress(name)
}
func (k Keeper) GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI {
	return k.sk.GetModuleAccount(ctx, moduleName)
}

func (k Keeper) SetModuleAccount(ctx sdk.Context, macc supplyexported.ModuleAccountI) {
	k.sk.SetModuleAccount(ctx, macc)
}

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) sdk.Error {
	return k.sk.SendCoinsFromModuleToModule(ctx, senderPool, recipientPool, amt)
}
func (k Keeper) UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.sk.UndelegateCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}
func (k Keeper) DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error {
	return k.sk.DelegateCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}

func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string,
	recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return k.sk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}
func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
	recipientModule string, amt sdk.Coins) sdk.Error {
	return k.sk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}

func (k Keeper) BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error {
	err := k.sk.SendCoinsFromModuleToModule(ctx, name, distribution.ModuleName, amt)
	if err != nil {
		return err
	}

	feePool := k.dk.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoins(amt))
	k.dk.SetFeePool(ctx, feePool)
	return nil
}
