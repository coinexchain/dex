package app

import (
	"encoding/json"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tm-db"
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

	"github.com/coinexchain/dex/modules/alias"
	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bancorlite"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/comment"
	"github.com/coinexchain/dex/modules/distributionx"
	"github.com/coinexchain/dex/modules/incentive"
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/msgqueue"
	"github.com/coinexchain/dex/modules/stakingx"
	"github.com/coinexchain/dex/modules/supplyx"
)

const (
	appName = "CoinExChainApp"
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
	ModuleBasics OrderedBasicManager

	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		authx.ModuleName:          nil,
		asset.ModuleName:          {supply.Burner, supply.Minter},
	}
)

type impKeeper struct {
	*CetChainApp
}

func (ik impKeeper) GetMarketLastExePrice(ctx sdk.Context, symbol string) (sdk.Dec, error) {
	return ik.marketKeeper.GetMarketLastExePrice(ctx, symbol)
}

func (ik impKeeper) IsMarketExist(ctx sdk.Context, symbol string) bool {
	return ik.marketKeeper.IsMarketExist(ctx, symbol)
}

func init() {
	modules := []module.AppModuleBasic{
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		params.AppModuleBasic{},
		authx.AppModuleBasic{}, //before `bank` to override `/bank/balances/{address}`
		bankx.AppModuleBasic{},
		bank.AppModuleBasic{},
		distr.AppModuleBasic{},
		supply.AppModuleBasic{},
		AuthModuleBasic{},
		stakingx.AppModuleBasic{}, //before `staking` to override `cetcli q staking pool` command
		StakingModuleBasic{},
		SlashingModuleBasic{},
		CrisisModuleBasic{},
		GovModuleBasic{gov.NewAppModuleBasic(paramsclient.ProposalHandler, distrclient.ProposalHandler)},
		distributionx.AppModuleBasic{},
		incentive.AppModuleBasic{},
		asset.AppModuleBasic{},
		market.AppModuleBasic{},
		bancorlite.AppModuleBasic{},
		alias.AppModuleBasic{},
		comment.AppModuleBasic{},
	}

	ModuleBasics = NewOrderedBasicManager(modules)
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
	cdc       *codec.Codec
	txDecoder sdk.TxDecoder // unmarshal []byte into sdk.Tx
	txCount   int64

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
	assetKeeper     asset.Keeper
	tokenKeeper     asset.TokenKeeper
	paramsKeeper    params.Keeper
	marketKeeper    market.Keeper
	bancorKeeper    bancorlite.Keeper
	msgQueProducer  msgqueue.Producer
	aliasKeeper     alias.Keeper
	commentKeeper   comment.Keeper

	// the module manager
	mm *module.Manager
}

// NewCetChainApp returns a reference to an initialized CetChainApp.
func NewCetChainApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *CetChainApp {

	cdc := MakeCodec()

	txDecoder := auth.DefaultTxDecoder(cdc)
	bApp := bam.NewBaseApp(appName, logger, db, txDecoder, baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	app := newCetChainApp(bApp, cdc, invCheckPeriod, txDecoder)
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

func newCetChainApp(bApp *bam.BaseApp, cdc *codec.Codec, invCheckPeriod uint, txDecoder sdk.TxDecoder) *CetChainApp {
	return &CetChainApp{
		BaseApp:        bApp,
		txDecoder:      txDecoder,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keyMain:        sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:     sdk.NewKVStoreKey(auth.StoreKey),
		keyAccountX:    sdk.NewKVStoreKey(authx.StoreKey),
		keySupply:      sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:     sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:    sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:       sdk.NewKVStoreKey(distr.StoreKey),
		keySlashing:    sdk.NewKVStoreKey(slashing.StoreKey),
		keyGov:         sdk.NewKVStoreKey(gov.StoreKey),
		keyParams:      sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:     sdk.NewTransientStoreKey(params.TStoreKey),
		keyAsset:       sdk.NewKVStoreKey(asset.StoreKey),
		keyMarket:      sdk.NewKVStoreKey(market.StoreKey),
		keyBancor:      sdk.NewKVStoreKey(bancorlite.StoreKey),
		keyIncentive:   sdk.NewKVStoreKey(incentive.StoreKey),
		keyAlias:       sdk.NewKVStoreKey(alias.StoreKey),
		keyComment:     sdk.NewKVStoreKey(comment.StoreKey),
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
		bank.DefaultCodespace, app.ModuleAccountAddrs(),
	)

	app.supplyKeeper = supply.NewKeeper(app.cdc, app.keySupply, app.accountKeeper,
		app.bankKeeper, maccPerms)

	var stakingKeeper staking.Keeper

	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		&stakingKeeper,
		app.supplyKeeper,
		distr.DefaultCodespace,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)
	supplyxKeeper := supplyx.NewKeeper(app.supplyKeeper, app.distrKeeper)

	stakingKeeper = staking.NewKeeper(
		app.cdc,
		app.keyStaking, app.tkeyStaking,
		supplyxKeeper,
		//app.supplyKeeper,
		app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
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
		//app.supplyKeeper,
		supplyxKeeper,
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
	eventTypeMsgQueue := ""
	if app.msgQueProducer.IsSubScribe(authx.ModuleName) {
		eventTypeMsgQueue = msgqueue.EventTypeMsgQueue
	}
	app.accountXKeeper = authx.NewKeeper(
		app.cdc,
		app.keyAccountX,
		app.paramsKeeper.Subspace(authx.DefaultParamspace),
		app.supplyKeeper,
		app.accountKeeper,
		eventTypeMsgQueue,
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		//app.stakingXKeeper,
		&stakingKeeper,
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
		app.supplyKeeper,
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
		app.supplyKeeper,
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

	ik := impKeeper{app}
	app.bancorKeeper = bancorlite.NewBaseKeeper(
		bancorlite.NewBancorInfoKeeper(app.keyBancor, app.cdc, app.paramsKeeper.Subspace(bancorlite.StoreKey)),
		app.bankxKeeper,
		app.assetKeeper,
		ik,
		app.msgQueProducer)

	app.marketKeeper = market.NewBaseKeeper(
		app.keyMarket,
		app.tokenKeeper,
		app.bankxKeeper,
		app.cdc,
		app.msgQueProducer,
		app.paramsKeeper.Subspace(market.StoreKey),
		app.bancorKeeper,
	)
	// register the staking hooks
	// NOTE: The stakingKeeper above is passed by reference, so that it can be
	// modified like below:
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()))

	eventTypeMsgQueue = ""
	if app.msgQueProducer.IsSubScribe(comment.ModuleName) {
		eventTypeMsgQueue = msgqueue.EventTypeMsgQueue
	}
	app.commentKeeper = *comment.NewBaseKeeper(
		app.keyComment,
		app.bankxKeeper,
		app.assetKeeper,
		app.distrxKeeper,
		eventTypeMsgQueue,
	)
	app.aliasKeeper = alias.NewBaseKeeper(
		app.keyAlias,
		app.bankxKeeper,
		app.assetKeeper,
		app.paramsKeeper.Subspace(alias.StoreKey),
	)
}

func (app *CetChainApp) InitModules() {

	modules := []module.AppModule{
		genaccounts.NewAppModule(app.accountKeeper),
		auth.NewAppModule(app.accountKeeper),
		authx.NewAppModule(app.accountXKeeper, app.tokenKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		bankx.NewAppModule(app.bankxKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		incentive.NewAppModule(app.incentiveKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		distributionx.NewAppModule(app.distrxKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),
		stakingx.NewAppModule(app.stakingXKeeper),
		asset.NewAppModule(app.assetKeeper),
		market.NewAppModule(app.marketKeeper),
		bancorlite.NewAppModule(app.bancorKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		alias.NewAppModule(app.aliasKeeper),
		comment.NewAppModule(app.commentKeeper),
	}

	app.mm = module.NewManager(modules...)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(market.ModuleName, incentive.ModuleName, distr.ModuleName, slashing.ModuleName)

	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName, authx.ModuleName, market.ModuleName, crisis.ModuleName)

	initGenesisOrder := []string{
		genaccounts.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		supply.ModuleName,
		authx.ModuleName,
		bankx.ModuleName,
		incentive.ModuleName,
		stakingx.ModuleName,
		asset.ModuleName,
		market.ModuleName,
		bancorlite.ModuleName,
		crisis.ModuleName,
		genutil.ModuleName, //call DeliverGenTxs in genutil at last
		alias.ModuleName,
		comment.ModuleName,
	}

	// genutils must occur after staking so that pools are properly
	// initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(initGenesisOrder...)

	exportGenesisOrder := initGenesisOrder
	app.mm.SetOrderExportGenesis(exportGenesisOrder...)

	app.crisisKeeper.RegisterRoute(authx.ModuleName, "pre-total-supply", authx.PreTotalSupplyInvariant(app.accountXKeeper))
	app.mm.RegisterInvariants(&app.crisisKeeper)

	//crisis module should be reset since invariants has been registered to crisis keeper
	app.replaceEmptyCrisisModule(&modules)

	registerRoutesWithOrder(modules, app.Router(), app.QueryRouter())
}

func (app *CetChainApp) replaceEmptyCrisisModule(modules *[]module.AppModule) {
	crisisWithInvariants := crisis.NewAppModule(&app.crisisKeeper)

	app.mm.Modules[crisis.ModuleName] = crisisWithInvariants

	for i, module := range *modules {
		if module.Name() == crisis.ModuleName {
			(*modules)[i] = crisisWithInvariants
		}
	}
}

func registerRoutesWithOrder(modules []module.AppModule, router sdk.Router, queryRouter sdk.QueryRouter) {
	for _, module := range modules {
		if module.Route() != "" && module.Route() != "bank" {
			router.AddRoute(module.Route(), module.NewHandler())
		}
		if module.QuerierRoute() != "" {
			queryRouter.AddRoute(module.QuerierRoute(), module.NewQuerierHandler())
		}
	}
}

// initialize BaseApp
func (app *CetChainApp) mountStores() {
	app.MountStores(app.keyMain, app.keyAccount, app.keySupply, app.keyStaking, app.keyDistr,
		app.keySlashing, app.keyGov, app.keyParams,
		app.tkeyParams, app.tkeyStaking,
		app.keyAccountX, app.keyAsset, app.keyMarket, app.keyIncentive,
		app.keyBancor, app.keyAlias, app.keyComment,
	)
}

// application updates every end block
func (app *CetChainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	PubMsgs = make([]PubMsg, 0, 10000)
	ret := app.mm.BeginBlock(ctx, req)
	if app.msgQueProducer.IsOpenToggle() {
		app.txCount = req.Header.TotalTxs - req.Header.NumTxs
		ret.Events = FilterMsgsOnlyKafka(ret.Events)
	}
	return ret
}

func (app *CetChainApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	ret := app.BaseApp.DeliverTx(req)
	if app.msgQueProducer.IsOpenToggle() {
		if ret.Code == uint32(sdk.CodeOK) {
			ret.Events = FilterMsgsOnlyKafka(ret.Events)
			app.notifySigners(req, ret.Events)
			app.notifyInTx(ret.Events)
		} else {
			ret.Events = RemoveMsgsOnlyKafka(ret.Events)
		}
	}
	return ret
}

// application updates every end block
// nolint: unparam
func (app *CetChainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	ret := app.mm.EndBlock(ctx, req)
	if app.msgQueProducer.IsOpenToggle() {
		ret.Events = FilterMsgsOnlyKafka(ret.Events)
		app.notifyComplete(ret.Events)
	}
	return ret
}

func (app *CetChainApp) Commit() abci.ResponseCommit {
	for _, msg := range PubMsgs {
		app.msgQueProducer.SendMsg(msg.Key, msg.Value)
	}
	return app.BaseApp.Commit()
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

// load a particular height
func (app *CetChainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *CetChainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.supplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
