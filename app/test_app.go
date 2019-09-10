package app

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/alias"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bancorlite"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/stakingx"
)

type TestApp struct {
	*CetChainApp
	Cms types.MultiStore
}

func NewTestApp() *TestApp {
	app := newTestApp()
	app.initKeepers(0)
	app.mountStores()
	return app
}

func newTestApp() *TestApp {
	cdc := MakeCodec()
	txDecoder := auth.DefaultTxDecoder(cdc)

	db := dbm.NewMemDB()
	bApp := bam.NewBaseApp(appName, log.NewNopLogger(), db, txDecoder)
	//bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion("ut")

	return &TestApp{
		CetChainApp: newCetChainApp(bApp, cdc, 0, txDecoder),
		Cms:         store.NewCommitMultiStore(db),
	}
}

func (app *TestApp) mountStores() {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(app.keyMain, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.tkeyParams, sdk.StoreTypeTransient, db)
	cms.MountStoreWithDB(app.keyAccount, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keySupply, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyStaking, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyDistr, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keySlashing, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyParams, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyGov, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyAccountX, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyAsset, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyMarket, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyIncentive, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyBancor, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyAlias, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyComment, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.keyStakingX, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(app.tkeyStaking, sdk.StoreTypeTransient, db)
	_ = cms.LoadLatestVersion()
	app.Cms = cms
}

func (app *TestApp) Cdc() *codec.Codec                    { return app.cdc }
func (app *TestApp) AliasKeeper() alias.Keeper            { return app.aliasKeeper }
func (app *TestApp) AssetKeeper() asset.Keeper            { return app.assetKeeper }
func (app *TestApp) BancorKeeper() bancorlite.Keeper      { return app.bancorKeeper }
func (app *TestApp) MarketKeeper() market.Keeper          { return app.marketKeeper }
func (app *TestApp) IncentiveKeeper() incentive.Keeper    { return app.incentiveKeeper }
func (app *TestApp) AccountXKeeper() authx.AccountXKeeper { return app.accountXKeeper }
func (app *TestApp) BankXKeeper() bankx.Keeper            { return app.bankxKeeper }
func (app *TestApp) StakingXKeeper() stakingx.Keeper      { return app.stakingXKeeper }
func (app *TestApp) SupplyKeeper() supply.Keeper          { return app.supplyKeeper }
func (app *TestApp) StakingKeeper() staking.Keeper        { return app.stakingKeeper }
func (app *TestApp) AccountKeeper() auth.AccountKeeper    { return app.accountKeeper }

func (app *TestApp) NewCtx() sdk.Context {
	return sdk.NewContext(app.Cms,
		abci.Header{ChainID: "test-chain-id", Time: time.Now()},
		false, log.NewNopLogger())
}
