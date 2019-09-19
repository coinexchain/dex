package asset_test

import (
	"os"
	"testing"

	"github.com/coinexchain/dex/testapp"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/asset/internal/types"
	dex "github.com/coinexchain/dex/types"
)

func TestMain(m *testing.M) {
	dex.InitSdkConfig()
	os.Exit(m.Run())
}

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	tk  asset.Keeper
}

func createTestInput() testInput {

	app := testapp.NewTestApp()
	ctx := app.NewCtx()
	app.AssetKeeper.SetParams(ctx, types.DefaultParams())

	initSupply := dex.NewCetCoinsE8(10000)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(initSupply))
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	_ = notBondedPool.SetCoins(initSupply)
	app.SupplyKeeper.SetModuleAccount(ctx, notBondedPool)

	return testInput{app.Cdc, ctx, app.AssetKeeper}
}

var _, _, testAddr = keyPubAddr()

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func mockAddrList() (list []sdk.AccAddress) {
	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addr3, _ = sdk.AccAddressFromBech32("coinex1zvf0hx6rpz0n7dkuzu34s39dnsyr8eygqs8h3q")

	list = append(list, addr1)
	list = append(list, addr2)
	list = append(list, addr3)
	return
}
