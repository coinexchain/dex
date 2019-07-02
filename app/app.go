package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/coinexchain/dex/modules/crisisx"
	"github.com/coinexchain/dex/modules/distributionx"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/govx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/modules/stakingx"
)

const (
	appName = "CetChainApp"
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cetcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.cetd")
)

// Extended ABCI application
type CetChainApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyAccountX      *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
	keySlashing      *sdk.KVStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyGov           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey
	keyAsset         *sdk.KVStoreKey
	keyMarket        *sdk.KVStoreKey
	keyIncentive     *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountKeeper       auth.AccountKeeper
	accountXKeeper      authx.AccountXKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.BaseKeeper
	bankxKeeper         bankx.Keeper // TODO rename to bankXKeeper
	stakingKeeper       staking.Keeper
	stakingXKeeper      stakingx.Keeper
	slashingKeeper      slashing.Keeper
	distrKeeper         distr.Keeper
	distrxKeeper        distributionx.Keeper
	govKeeper           gov.Keeper
	crisisKeeper        crisis.Keeper
	incentiveKeeper     incentive.Keeper
	assetKeeper         asset.BaseKeeper
	tokenKeeper         asset.TokenKeeper
	paramsKeeper        params.Keeper
	marketKeeper        market.Keeper
	msgQueProducer      msgqueue.Producer
}

// NewCetChainApp returns a reference to an initialized CetChainApp.
func NewCetChainApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *CetChainApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	app := newCetChainApp(bApp, cdc, invCheckPeriod)
	app.initKeepers()
	app.registerCrisisRoutes()
	app.registerMessageRoutes()
	app.mountStores()

	ah := authx.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper, app.accountXKeeper,
		newAnteHelper(app.accountXKeeper, app.stakingXKeeper))

	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ah)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keyMain)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

func newCetChainApp(bApp *bam.BaseApp, cdc *codec.Codec, invCheckPeriod uint) *CetChainApp {
	return &CetChainApp{
		BaseApp:          bApp,
		cdc:              cdc,
		invCheckPeriod:   invCheckPeriod,
		keyMain:          sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keyAccountX:      sdk.NewKVStoreKey(authx.StoreKey),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:         sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:        sdk.NewTransientStoreKey(distr.TStoreKey),
		keySlashing:      sdk.NewKVStoreKey(slashing.StoreKey),
		keyGov:           sdk.NewKVStoreKey(gov.StoreKey),
		keyFeeCollection: sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyParams:        sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:       sdk.NewTransientStoreKey(params.TStoreKey),
		keyAsset:         sdk.NewKVStoreKey(asset.StoreKey),
		keyMarket:        sdk.NewKVStoreKey(market.StoreKey),
		keyIncentive:     sdk.NewKVStoreKey(incentive.StoreKey),
	}
}

func (app *CetChainApp) initKeepers() {
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)
	app.msgQueProducer = msgqueue.NewProducer()
	// define the accountKeeper
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	// add handlers
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		app.cdc,
		app.keyFeeCollection,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking, app.tkeyStaking,
		app.bankKeeper, app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakingKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	govBankKeeper := govx.NewKeeper(
		app.bankKeeper,
		app.accountKeeper,
		app.distrKeeper,
	)
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.paramsKeeper, app.paramsKeeper.Subspace(gov.DefaultParamspace), &govBankKeeper, &stakingKeeper,
		gov.DefaultCodespace,
	)
	app.crisisKeeper = crisis.NewKeeper(
		app.paramsKeeper.Subspace(crisis.DefaultParamspace),
		app.distrKeeper, app.bankKeeper,
		app.feeCollectionKeeper,
	)

	// cet keepers
	app.accountXKeeper = authx.NewKeeper(
		app.cdc,
		app.keyAccountX,
		app.paramsKeeper.Subspace(authx.DefaultParamspace),
	)
	app.stakingXKeeper = stakingx.NewKeeper(
		app.paramsKeeper.Subspace(stakingx.DefaultParamspace), &stakingKeeper, app.distrKeeper, app.accountKeeper)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		app.stakingXKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.incentiveKeeper = incentive.NewKeeper(
		app.cdc, app.keyIncentive, app.paramsKeeper.Subspace(incentive.DefaultParamspace), app.feeCollectionKeeper, app.bankKeeper,
	)
	app.tokenKeeper = asset.NewBaseTokenKeeper(
		app.cdc, app.keyAsset,
	)
	app.bankxKeeper = bankx.NewKeeper(
		app.paramsKeeper.Subspace(bankx.DefaultParamspace),
		app.accountXKeeper, app.bankKeeper, app.accountKeeper,
		app.feeCollectionKeeper,
		app.tokenKeeper,
		app.msgQueProducer,
	)
	app.distrxKeeper = distributionx.NewKeeper(
		app.bankxKeeper,
		app.distrKeeper,
	)
	app.assetKeeper = asset.NewBaseKeeper(
		app.cdc,
		app.keyAsset,
		app.paramsKeeper.Subspace(asset.DefaultParamspace),
		app.bankxKeeper,
		&app.stakingKeeper,
	)
	app.marketKeeper = market.NewKeeper(
		app.keyMarket,
		app.tokenKeeper,
		app.bankxKeeper,
		app.feeCollectionKeeper,
		app.cdc,
		app.msgQueProducer,
		app.paramsKeeper.Subspace(market.StoreKey),
	)

	// register the staking hooks
	// NOTE: The stakingKeeper above is passed by reference, so that it can be
	// modified like below:
	app.stakingKeeper = *stakingKeeper.SetHooks(
		NewStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

}

func (app *CetChainApp) registerCrisisRoutes() {
	crisisx.RegisterInvariants(&app.crisisKeeper, app.assetKeeper, app.bankxKeeper, app.feeCollectionKeeper, app.distrKeeper, app.stakingKeeper)
	bank.RegisterInvariants(&app.crisisKeeper, app.accountKeeper)
	distr.RegisterInvariants(&app.crisisKeeper, app.distrKeeper, app.stakingKeeper)

	//Invariants checks of staking module will adjust and included by stakingx.RegisterInvariants
	stakingx.RegisterInvariants(&app.crisisKeeper, app.stakingXKeeper, app.assetKeeper, app.stakingKeeper)
}

func (app *CetChainApp) registerMessageRoutes() {
	app.Router().
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
		AddRoute(distr.RouterKey, distr.NewHandler(app.distrKeeper)).
		AddRoute(slashing.RouterKey, slashing.NewHandler(app.slashingKeeper)).
		AddRoute(gov.RouterKey, gov.NewHandler(app.govKeeper)).
		AddRoute(crisis.RouterKey, crisis.NewHandler(app.crisisKeeper)).
		AddRoute(bankx.RouterKey, bankx.NewHandler(app.bankxKeeper)).
		AddRoute(asset.RouterKey, asset.NewHandler(app.assetKeeper)).
		AddRoute(market.RouterKey, market.NewHandler(app.marketKeeper)).
		AddRoute(distributionx.RouterKey, distributionx.NewHandler(app.distrxKeeper))

	app.QueryRouter().
		AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
		AddRoute(authx.QuerierRoute, authx.NewQuerier(app.accountXKeeper)).
		AddRoute(distr.QuerierRoute, distr.NewQuerier(app.distrKeeper)).
		AddRoute(gov.QuerierRoute, gov.NewQuerier(app.govKeeper)).
		AddRoute(slashing.QuerierRoute, slashing.NewQuerier(app.slashingKeeper, app.cdc)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper, app.cdc)).
		AddRoute(stakingx.QuerierRoute, stakingx.NewQuerier(app.stakingXKeeper, app.cdc)).
		AddRoute(asset.QuerierRoute, asset.NewQuerier(app.tokenKeeper, app.cdc)).
		AddRoute(market.StoreKey, market.NewQuerier(app.marketKeeper, app.cdc))
}

// initialize BaseApp
func (app *CetChainApp) mountStores() {
	app.MountStores(app.keyMain, app.keyAccount, app.keyStaking, app.keyDistr,
		app.keySlashing, app.keyGov, app.keyFeeCollection, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistr,
		app.keyAccountX, app.keyAsset, app.keyMarket, app.keyIncentive,
	)
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bankx.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	distributionx.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	authx.RegisterCodec(cdc)
	crisis.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	asset.RegisterCodec(cdc)
	market.RegisterCodec(cdc)
	incentive.RegisterCodec(cdc)
	return cdc
}

// application updates every end block
func (app *CetChainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	//block incentive for the previous block
	_ = incentive.BeginBlocker(ctx, app.incentiveKeeper)

	// distribute rewards for the previous block
	distr.BeginBlocker(ctx, req, app.distrKeeper)

	// slash anyone who double signed.
	// NOTE: This should happen after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool,
	// so as to keep the CanWithdrawInvariant invariant.
	// TODO: This should really happen at EndBlocker.
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
// nolint: unparam
func (app *CetChainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, app.govKeeper)
	validatorUpdates, endBlockerTags := staking.EndBlocker(ctx, app.stakingKeeper)
	tags = append(tags, endBlockerTags...)
	authx.EndBlocker(ctx, app.accountXKeeper, app.accountKeeper)
	market.EndBlocker(ctx, app.marketKeeper)

	if app.invCheckPeriod != 0 && ctx.BlockHeight()%int64(app.invCheckPeriod) == 0 {
		app.assertRuntimeInvariants()
	}

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// custom logic for coindex initialization
func (app *CetChainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	// TODO is this now the whole genesis file?

	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	// for i := range genesisState.Accounts {
	// 	fmt.Printf("genesis accout address : %s, coins : %s\n", genesisState.Accounts[i].Address.String(), genesisState.Accounts[i].Coins.String())
	// }

	validators := app.initFromGenesisState(ctx, genesisState)

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(validators)))
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	// assert runtime invariants
	app.assertRuntimeInvariants()

	return abci.ResponseInitChain{
		Validators: validators,
	}
}

// initialize store from a genesis state
func (app *CetChainApp) initFromGenesisState(ctx sdk.Context, genesisState GenesisState) []abci.ValidatorUpdate {
	genesisState.Sanitize()

	// load the accounts
	app.loadGenesisAccounts(ctx, genesisState)

	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)

	// load the initial staking information
	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	// initialize module-specific stores
	app.initModuleStores(ctx, genesisState)

	// validate genesis state
	if err := genesisState.Validate(); err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	if len(genesisState.GenTxs) > 0 {
		app.deliverGenTxs(genesisState.GenTxs)

		validators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}
	return validators
}

func (app *CetChainApp) loadGenesisAccounts(ctx sdk.Context, genesisState GenesisState) {
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc = app.accountKeeper.NewAccount(ctx, acc) // set account number
		app.accountKeeper.SetAccount(ctx, acc)

		accx := authx.AccountX{Address: gacc.Address, MemoRequired: gacc.MemoRequired, LockedCoins: gacc.LockedCoins}
		app.accountXKeeper.SetAccountX(ctx, accx)
	}
}

func (app *CetChainApp) initModuleStores(ctx sdk.Context, genesisState GenesisState) {
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	authx.InitGenesis(ctx, app.accountXKeeper, genesisState.AuthXData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	bankx.InitGenesis(ctx, app.bankxKeeper, genesisState.BankXData)
	stakingx.InitGenesis(ctx, app.stakingXKeeper, genesisState.StakingXData)
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakingData.Validators.ToSDKValidators())
	gov.InitGenesis(ctx, app.govKeeper, genesisState.GovData)
	crisis.InitGenesis(ctx, app.crisisKeeper, genesisState.CrisisData)
	asset.InitGenesis(ctx, app.assetKeeper, genesisState.AssetData)
	market.InitGenesis(ctx, app.marketKeeper, genesisState.MarketData)
	incentive.InitGenesis(ctx, app.incentiveKeeper, genesisState.Incentive)
}

func (app *CetChainApp) deliverGenTxs(genTxs []json.RawMessage) {
	for _, genTx := range genTxs {
		var tx auth.StdTx
		err := app.cdc.UnmarshalJSON(genTx, &tx)
		if err != nil {
			panic(err)
		}
		bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
		res := app.BaseApp.DeliverTx(bz)
		if !res.IsOK() {
			panic(res.Log)
		}
	}
}

// load a particular height
func (app *CetChainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}
