package app

import (
	"encoding/json"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/asset/client"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/distributionx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
	market_client "github.com/coinexchain/dex/modules/market/client"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/modules/stakingx"
	stakingx_client "github.com/coinexchain/dex/modules/stakingx/client"
)

const (
	appName = "CetChainApp"
	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"
)

// default home directories for expected binaries
var (
	// default home directories for cetcli
	DefaultCLIHome = os.ExpandEnv("$HOME/.cetcli")

	// default home directories for cetd
	DefaultNodeHome = os.ExpandEnv("$HOME/.cetd")

	// The ModuleBasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics module.BasicManager
)

func init() {
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distrclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		asset.AppModuleBasic{},
		//TODO: bankx
		//TODO: distributionx
		//TODO: authx
		market.AppModuleBasic{},
		incentive.AppModuleBasic{},
		bankx.AppModuleBasic{},
	)
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

// Extended ABCI application
type CetChainApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keyMain      *sdk.KVStoreKey
	keyAccount   *sdk.KVStoreKey
	keyAccountX  *sdk.KVStoreKey
	keySupply    *sdk.KVStoreKey
	keyStaking   *sdk.KVStoreKey
	tkeyStaking  *sdk.TransientStoreKey
	keySlashing  *sdk.KVStoreKey
	keyDistr     *sdk.KVStoreKey
	tkeyDistr    *sdk.TransientStoreKey
	keyGov       *sdk.KVStoreKey
	keyParams    *sdk.KVStoreKey
	tkeyParams   *sdk.TransientStoreKey
	keyAsset     *sdk.KVStoreKey
	keyMarket    *sdk.KVStoreKey
	keyIncentive *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountKeeper   auth.AccountKeeper
	accountXKeeper  authx.AccountXKeeper
	bankKeeper      bank.BaseKeeper
	bankxKeeper     bankx.Keeper // TODO rename to bankXKeeper
	supplyKeeper    supply.Keeper
	stakingKeeper   staking.Keeper
	stakingXKeeper  stakingx.Keeper
	slashingKeeper  slashing.Keeper
	distrKeeper     distr.Keeper
	distrxKeeper    distributionx.Keeper
	govKeeper       gov.Keeper
	crisisKeeper    crisis.Keeper
	incentiveKeeper incentive.Keeper
	assetKeeper     asset.BaseKeeper
	tokenKeeper     asset.TokenKeeper
	paramsKeeper    params.Keeper
	marketKeeper    market.Keeper
	msgQueProducer  msgqueue.Producer

	// the module manager
	mm *module.Manager
}

// NewCetChainApp returns a reference to an initialized CetChainApp.
func NewCetChainApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *CetChainApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	app := newCetChainApp(bApp, cdc, invCheckPeriod)
	app.initKeepers(invCheckPeriod)
	app.InitModules()
	app.mountStores()

	ah := authx.NewAnteHandler(app.accountKeeper, app.supplyKeeper, app.accountXKeeper,
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
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keyMain:        sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:     sdk.NewKVStoreKey(auth.StoreKey),
		keyAccountX:    sdk.NewKVStoreKey(authx.StoreKey),
		keySupply:      sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:     sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:    sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:       sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:      sdk.NewTransientStoreKey(distr.TStoreKey),
		keySlashing:    sdk.NewKVStoreKey(slashing.StoreKey),
		keyGov:         sdk.NewKVStoreKey(gov.StoreKey),
		keyParams:      sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:     sdk.NewTransientStoreKey(params.TStoreKey),
		keyAsset:       sdk.NewKVStoreKey(asset.StoreKey),
		keyMarket:      sdk.NewKVStoreKey(market.StoreKey),
		keyIncentive:   sdk.NewKVStoreKey(incentive.StoreKey),
	}
}

func (app *CetChainApp) initKeepers(invCheckPeriod uint) {
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams, params.DefaultCodespace)
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

	// account permissions
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     {supply.Basic},
		distr.ModuleName:          {supply.Basic},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}

	app.supplyKeeper = supply.NewKeeper(app.cdc, app.keySupply, app.accountKeeper,
		app.bankKeeper, supply.DefaultCodespace, maccPerms)

	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking, app.tkeyStaking,
		app.supplyKeeper,
		app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		&stakingKeeper,
		app.supplyKeeper,
		distr.DefaultCodespace,
		auth.FeeCollectorName,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))

	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.paramsKeeper, app.paramsKeeper.Subspace(gov.DefaultParamspace),
		app.supplyKeeper,
		&stakingKeeper,
		gov.DefaultCodespace,
		govRouter,
	)

	app.crisisKeeper = crisis.NewKeeper(
		app.paramsKeeper.Subspace(crisis.DefaultParamspace),
		invCheckPeriod,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)

	// cet keepers
	app.accountXKeeper = authx.NewKeeper(
		app.cdc,
		app.keyAccountX,
		app.paramsKeeper.Subspace(authx.DefaultParamspace),
	)

	app.stakingXKeeper = stakingx.NewKeeper(
		app.paramsKeeper.Subspace(stakingx.DefaultParamspace),
		app.assetKeeper,
		&stakingKeeper,
		app.distrKeeper,
		app.accountKeeper,
		app.bankxKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		app.stakingXKeeper,
		app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.incentiveKeeper = incentive.NewKeeper(
		app.cdc, app.keyIncentive,
		app.paramsKeeper.Subspace(incentive.DefaultParamspace),
		app.bankKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
	)
	app.tokenKeeper = asset.NewBaseTokenKeeper(
		app.cdc, app.keyAsset,
	)
	app.bankxKeeper = bankx.NewKeeper(
		app.paramsKeeper.Subspace(bankx.DefaultParamspace),
		app.accountXKeeper, app.bankKeeper, app.accountKeeper,
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

func (app *CetChainApp) InitModules() {
	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		//TODO: authx
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		//TODO: bankx
		//TODO: distributionx
		crisis.NewAppModule(app.crisisKeeper),
		incentive.NewAppModule(app.incentiveKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),
		stakingx.NewAppModule(app.stakingXKeeper, stakingx_client.NewStakingXModuleClient()),
		asset.NewAppModule(app.assetKeeper, client.NewAssetModuleClient()),
		market.NewAppModule(app.marketKeeper, market_client.NewMarketModuleClient()),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(market.ModuleName, incentive.ModuleName, distr.ModuleName, slashing.ModuleName)

	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName, authx.ModuleName, market.ModuleName, crisis.ModuleName)

	// genutils must occur after staking so that pools are properly
	// initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(genaccounts.ModuleName, distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName, slashing.ModuleName,
		gov.ModuleName, supply.ModuleName, crisis.ModuleName,
		//TODO: authx.ModuleName,
		//TODO: bankx.ModuleName,
		//TODO: distributionx.ModuleName
		stakingx.ModuleName,
		asset.ModuleName,
		market.ModuleName,
		incentive.ModuleName,
		genutil.ModuleName, //call DeliverGenTxs in genutil at last
	)

	//TODO: set export genesis order

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())
}

// initialize BaseApp
func (app *CetChainApp) mountStores() {
	app.MountStores(app.keyMain, app.keyAccount, app.keyStaking, app.keyDistr,
		app.keySlashing, app.keyGov, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistr,
		app.keyAccountX, app.keyAsset, app.keyMarket, app.keyIncentive,
	)
}

// application updates every end block
func (app *CetChainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
// nolint: unparam
func (app *CetChainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// custom logic for coindex initialization
func (app *CetChainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState map[string]json.RawMessage
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	if err := ModuleBasics.ValidateGenesis(genesisState); err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

// initialize store from a genesis state
func (app *CetChainApp) initFromGenesisState(ctx sdk.Context, genesisState GenesisState) {
	// load the accounts
	app.loadGenesisAccounts(ctx, genesisState)
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

// load a particular height
func (app *CetChainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}
