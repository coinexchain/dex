package testapp

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
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	dexapp "github.com/coinexchain/dex/app"
	"github.com/coinexchain/dex/modules/alias"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bancorlite"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/comment"
	"github.com/coinexchain/dex/modules/distributionx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/coinexchain/dex/modules/supplyx"
	"github.com/coinexchain/dex/msgqueue"
)

// Fake App for unit tests
type TestApp struct {
	Cdc          *codec.Codec
	Cms          types.MultiStore
	keyMain      *sdk.KVStoreKey
	keyAccount   *sdk.KVStoreKey
	keyAccountX  *sdk.KVStoreKey
	keySupply    *sdk.KVStoreKey
	keyStaking   *sdk.KVStoreKey
	keyStakingX  *sdk.KVStoreKey
	tkeyStaking  *sdk.TransientStoreKey
	keySlashing  *sdk.KVStoreKey
	keyDistr     *sdk.KVStoreKey
	keyGov       *sdk.KVStoreKey
	keyParams    *sdk.KVStoreKey
	tkeyParams   *sdk.TransientStoreKey
	keyAsset     *sdk.KVStoreKey
	keyMarket    *sdk.KVStoreKey
	keyBancor    *sdk.KVStoreKey
	keyIncentive *sdk.KVStoreKey
	keyAlias     *sdk.KVStoreKey
	keyComment   *sdk.KVStoreKey

	// Manage getting and setting accounts
	AccountKeeper   auth.AccountKeeper
	AccountXKeeper  authx.AccountXKeeper
	BankKeeper      bank.BaseKeeper
	BankxKeeper     bankx.Keeper // TODO rename to BankxKeeper
	SupplyKeeper    supply.Keeper
	StakingKeeper   staking.Keeper
	StakingXKeeper  stakingx.Keeper
	SlashingKeeper  slashing.Keeper
	DistrKeeper     dist.Keeper
	DistrxKeeper    distributionx.Keeper
	GovKeeper       gov.Keeper
	CrisisKeeper    crisis.Keeper
	IncentiveKeeper incentive.Keeper
	AssetKeeper     asset.Keeper
	TokenKeeper     asset.TokenKeeper
	ParamsKeeper    params.Keeper
	MarketKeeper    market.Keeper
	BancorKeeper    bancorlite.Keeper
	MsgQueProducer  msgqueue.MsgSender
	AliasKeeper     alias.Keeper
	CommentKeeper   comment.Keeper
}

func NewTestApp() *TestApp {
	Cdc := dexapp.MakeCodec()
	app := newTestApp(Cdc)
	app.initKeepers(0)
	app.mountStores()
	return app
}

func newTestApp(Cdc *codec.Codec) *TestApp {
	return &TestApp{
		Cdc:          Cdc,
		keyMain:      sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:   sdk.NewKVStoreKey(auth.StoreKey),
		keyAccountX:  sdk.NewKVStoreKey(authx.StoreKey),
		keySupply:    sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:   sdk.NewKVStoreKey(staking.StoreKey),
		keyStakingX:  sdk.NewKVStoreKey(stakingx.StoreKey),
		tkeyStaking:  sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:     sdk.NewKVStoreKey(dist.StoreKey),
		keySlashing:  sdk.NewKVStoreKey(slashing.StoreKey),
		keyGov:       sdk.NewKVStoreKey(gov.StoreKey),
		keyParams:    sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:   sdk.NewTransientStoreKey(params.TStoreKey),
		keyAsset:     sdk.NewKVStoreKey(asset.StoreKey),
		keyMarket:    sdk.NewKVStoreKey(market.StoreKey),
		keyBancor:    sdk.NewKVStoreKey(bancorlite.StoreKey),
		keyIncentive: sdk.NewKVStoreKey(incentive.StoreKey),
		keyAlias:     sdk.NewKVStoreKey(alias.StoreKey),
		keyComment:   sdk.NewKVStoreKey(comment.StoreKey),
	}
}

func (app *TestApp) initKeepers(invCheckPeriod uint) {
	app.ParamsKeeper = params.NewKeeper(app.Cdc, app.keyParams, app.tkeyParams, params.DefaultCodespace)
	app.MsgQueProducer = msgqueue.NewProducer()
	// define the AccountKeeper
	app.AccountKeeper = auth.NewAccountKeeper(
		app.Cdc,
		app.keyAccount,
		app.ParamsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	// add handlers
	app.BankKeeper = bank.NewBaseKeeper(
		app.AccountKeeper,
		app.ParamsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace, app.ModuleAccountAddrs(),
	)

	app.SupplyKeeper = supply.NewKeeper(app.Cdc, app.keySupply, app.AccountKeeper,
		app.BankKeeper, dexapp.MaccPerms)

	var StakingKeeper staking.Keeper

	app.DistrKeeper = dist.NewKeeper(
		app.Cdc,
		app.keyDistr,
		app.ParamsKeeper.Subspace(dist.DefaultParamspace),
		&StakingKeeper,
		app.SupplyKeeper,
		dist.DefaultCodespace,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	supplyxKeeper := supplyx.NewKeeper(app.SupplyKeeper, app.DistrKeeper)

	StakingKeeper = staking.NewKeeper(
		app.Cdc,
		app.keyStaking, app.tkeyStaking,
		supplyxKeeper,
		//app.SupplyKeeper,
		app.ParamsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(dist.RouterKey, dist.NewCommunityPoolSpendProposalHandler(app.DistrKeeper))

	app.GovKeeper = gov.NewKeeper(
		app.Cdc,
		app.keyGov,
		app.ParamsKeeper, app.ParamsKeeper.Subspace(gov.DefaultParamspace),
		//app.SupplyKeeper,
		supplyxKeeper,
		&StakingKeeper,
		gov.DefaultCodespace,
		govRouter,
	)

	app.CrisisKeeper = crisis.NewKeeper(
		app.ParamsKeeper.Subspace(crisis.DefaultParamspace),
		invCheckPeriod,
		app.SupplyKeeper,
		auth.FeeCollectorName,
	)

	// cet keepers
	eventTypeMsgQueue := ""
	if app.MsgQueProducer.IsSubscribed(authx.ModuleName) {
		eventTypeMsgQueue = msgqueue.EventTypeMsgQueue
	}
	app.AccountXKeeper = authx.NewKeeper(
		app.Cdc,
		app.keyAccountX,
		app.ParamsKeeper.Subspace(authx.DefaultParamspace),
		app.SupplyKeeper,
		app.AccountKeeper,
		eventTypeMsgQueue,
	)

	app.SlashingKeeper = slashing.NewKeeper(
		app.Cdc,
		app.keySlashing,
		//app.StakingXKeeper,
		&StakingKeeper,
		app.ParamsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.IncentiveKeeper = incentive.NewKeeper(
		app.Cdc, app.keyIncentive,
		app.ParamsKeeper.Subspace(incentive.DefaultParamspace),
		app.BankKeeper,
		app.SupplyKeeper,
		auth.FeeCollectorName,
	)
	app.TokenKeeper = asset.NewBaseTokenKeeper(
		app.Cdc, app.keyAsset,
	)
	app.BankxKeeper = bankx.NewKeeper(
		app.ParamsKeeper.Subspace(bankx.DefaultParamspace),
		app.AccountXKeeper, app.BankKeeper, app.AccountKeeper,
		app.TokenKeeper,
		app.SupplyKeeper,
		app.MsgQueProducer,
	)
	app.DistrxKeeper = distributionx.NewKeeper(
		app.BankxKeeper,
		app.DistrKeeper,
	)
	app.AssetKeeper = asset.NewBaseKeeper(
		app.Cdc,
		app.keyAsset,
		app.ParamsKeeper.Subspace(asset.DefaultParamspace),
		app.BankxKeeper,
		app.SupplyKeeper,
	)
	app.StakingXKeeper = stakingx.NewKeeper(
		app.keyStakingX,
		app.Cdc,
		app.ParamsKeeper.Subspace(stakingx.DefaultParamspace),
		app.AssetKeeper,
		&StakingKeeper,
		app.DistrKeeper,
		app.AccountKeeper,
		app.BankxKeeper,
		app.SupplyKeeper,
		auth.FeeCollectorName,
	)

	app.BancorKeeper = bancorlite.NewBaseKeeper(
		bancorlite.NewBancorInfoKeeper(app.keyBancor, app.Cdc, app.ParamsKeeper.Subspace(bancorlite.StoreKey)),
		app.BankxKeeper,
		app.AssetKeeper,
		&app.MarketKeeper,
		app.MsgQueProducer)

	app.MarketKeeper = market.NewBaseKeeper(
		app.keyMarket,
		app.TokenKeeper,
		app.BankxKeeper,
		app.Cdc,
		app.MsgQueProducer,
		app.ParamsKeeper.Subspace(market.StoreKey),
		app.BancorKeeper,
		app.AccountKeeper,
	)
	// register the staking hooks
	// NOTE: The StakingKeeper above is passed by reference, so that it can be
	// modified like below:
	app.StakingKeeper = *StakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()))

	eventTypeMsgQueue = ""
	if app.MsgQueProducer.IsSubscribed(comment.ModuleName) {
		eventTypeMsgQueue = msgqueue.EventTypeMsgQueue
	}
	app.CommentKeeper = *comment.NewBaseKeeper(
		app.keyComment,
		app.BankxKeeper,
		app.AssetKeeper,
		app.AccountKeeper,
		app.DistrxKeeper,
		eventTypeMsgQueue,
	)
	app.AliasKeeper = alias.NewBaseKeeper(
		app.keyAlias,
		app.BankxKeeper,
		app.AssetKeeper,
		app.ParamsKeeper.Subspace(alias.StoreKey),
	)
}

func (app *TestApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range dexapp.MaccPerms {
		modAccAddrs[app.SupplyKeeper.GetModuleAddress(acc).String()] = true
	}
	return modAccAddrs
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

func (app *TestApp) NewCtx() sdk.Context {
	return sdk.NewContext(app.Cms,
		abci.Header{ChainID: "test-chain-id", Time: time.Now()},
		false, log.NewNopLogger())
}
