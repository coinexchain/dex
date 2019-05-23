package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/testutil"
	"github.com/coinexchain/dex/types"
)

const (
	defaultGenesisAccountNum = 7
)

var (
	defaultBondDenom = "cet"

	//TODO:update these address to be under coinex's control
	_, _, testAddress1 = testutil.KeyPubAddr()
	_, _, testAddress2 = testutil.KeyPubAddr()
	_, _, testAddress3 = testutil.KeyPubAddr()
	_, _, testAddress4 = testutil.KeyPubAddr()
	_, _, testAddress5 = testutil.KeyPubAddr()
	_, _, testAddress6 = testutil.KeyPubAddr()
	_, _, testAddress7 = testutil.KeyPubAddr()
)

// State to Unmarshal
type GenesisState struct {
	Accounts     []GenesisAccount          `json:"accounts"`
	AccountsX    []authx.AccountX          `json:"accountsx"` // additional account info bused by CoinEx chain
	AuthData     auth.GenesisState         `json:"auth"`
	BankData     bank.GenesisState         `json:"bank"`
	BankXData    bankx.GenesisState        `json:"bankx"`
	StakingData  staking.GenesisState      `json:"staking"`
	DistrData    distribution.GenesisState `json:"distr"`
	GovData      gov.GenesisState          `json:"gov"`
	CrisisData   crisis.GenesisState       `json:"crisis"`
	SlashingData slashing.GenesisState     `json:"slashing"`
	AssetData    asset.GenesisState        `json:"asset"`
	GenTxs       []json.RawMessage         `json:"gentxs"`
}

// NewDefaultGenesisState generates the default state for coindex.
func NewDefaultGenesisState() GenesisState {
	gs := GenesisState{
		Accounts:     DefaultCETGenesisAccount(),
		AccountsX:    nil,
		AuthData:     auth.DefaultGenesisState(),
		BankData:     bank.DefaultGenesisState(),
		BankXData:    bankx.DefaultGenesisState(),
		StakingData:  staking.DefaultGenesisState(),
		DistrData:    distribution.DefaultGenesisState(),
		GovData:      gov.DefaultGenesisState(),
		CrisisData:   crisis.DefaultGenesisState(),
		SlashingData: slashing.DefaultGenesisState(),
		AssetData:    asset.DefaultGenesisState(),
		GenTxs:       nil,
	}
	// TODO: create staking.GenesisState & gov.GenesisState & crisis.GenesisState from scratch
	gs.StakingData.Params.BondDenom = defaultBondDenom
	gs.GovData.DepositParams.MinDeposit[0].Denom = defaultBondDenom
	gs.CrisisData.ConstantFee.Denom = defaultBondDenom
	return gs
}

func NewGenesisState(
	accounts []GenesisAccount,
	accountsX []authx.AccountX,
	authData auth.GenesisState,
	bankData bank.GenesisState,
	bankxData bankx.GenesisState,
	stakingData staking.GenesisState,
	distrData distribution.GenesisState,
	govData gov.GenesisState,
	crisisData crisis.GenesisState,
	slashingData slashing.GenesisState,
	assetData asset.GenesisState) GenesisState {

	return GenesisState{
		Accounts:     accounts,
		AccountsX:    accountsX,
		AuthData:     authData,
		BankData:     bankData,
		BankXData:    bankxData,
		StakingData:  stakingData,
		DistrData:    distrData,
		GovData:      govData,
		CrisisData:   crisisData,
		SlashingData: slashingData,
		AssetData:    assetData,
	}
}

// Sanitize sorts accounts and coin sets.
func (gs GenesisState) Sanitize() {
	sort.Slice(gs.Accounts, func(i, j int) bool {
		return gs.Accounts[i].AccountNumber < gs.Accounts[j].AccountNumber
	})

	for _, acc := range gs.Accounts {
		acc.Coins = acc.Coins.Sort()
	}
}

// ValidateGenesisState ensures that the genesis state obeys the expected invariants
// TODO: No validators are both bonded and jailed (#2088)
// TODO: Error if there is a duplicate validator (#1708)
// TODO: Ensure all state machine parameters are in genesis (#1704)
func ValidateGenesisState(genesisState GenesisState) error {
	if err := validateGenesisStateAccounts(genesisState.Accounts); err != nil {
		return err
	}

	if err := asset.ValidateGenesis(genesisState.AssetData); err != nil {
		return err
	}

	// skip stakingData validation as genesis is created from txs
	if len(genesisState.GenTxs) > 0 {
		return nil
	}

	if err := auth.ValidateGenesis(genesisState.AuthData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(genesisState.BankData); err != nil {
		return err
	}
	if err := bankx.ValidateGenesis(genesisState.BankXData); err != nil {
		return err
	}
	if err := staking.ValidateGenesis(genesisState.StakingData); err != nil {
		return err
	}
	if err := distribution.ValidateGenesis(genesisState.DistrData); err != nil {
		return err
	}
	if err := gov.ValidateGenesis(genesisState.GovData); err != nil {
		return err
	}
	if err := crisis.ValidateGenesis(genesisState.CrisisData); err != nil {
		return err
	}

	return slashing.ValidateGenesis(genesisState.SlashingData)
}

// validateGenesisStateAccounts performs validation of genesis accounts. It
// ensures that there are no duplicate accounts in the genesis state and any
// provided vesting accounts are valid.
func validateGenesisStateAccounts(accs []GenesisAccount) error {
	addrMap := make(map[string]bool, len(accs))
	for _, acc := range accs {
		addrStr := acc.Address.String()

		// disallow any duplicate accounts
		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}

		// validate any vesting fields
		if !acc.OriginalVesting.IsZero() {
			if acc.EndTime == 0 {
				return fmt.Errorf("missing end time for vesting account; address: %s", addrStr)
			}

			if acc.StartTime >= acc.EndTime {
				return fmt.Errorf(
					"vesting start time must before end time; address: %s, start: %s, end: %s",
					addrStr,
					time.Unix(acc.StartTime, 0).UTC().Format(time.RFC3339),
					time.Unix(acc.EndTime, 0).UTC().Format(time.RFC3339),
				)
			}
		}

		addrMap[addrStr] = true
	}

	return nil
}

// CetAppGenState but with JSON
func CetAppGenStateJSON(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {
	// create the final app state
	genesisState, err := CetAppGenState(cdc, genDoc, appGenTxs)
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(cdc, genesisState)
}

// Create the core parameters for genesis initialization for gaia
// note that the pubkey input is this machines pubkey
func CetAppGenState(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	genesisState GenesisState, err error) {

	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return genesisState, errors.New("there must be at least one genesis tx")
	}

	stakingData := genesisState.StakingData
	for i, genTx := range appGenTxs {
		var tx auth.StdTx
		if err := cdc.UnmarshalJSON(genTx, &tx); err != nil {
			return genesisState, err
		}

		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return genesisState, errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}

		if _, ok := msgs[0].(staking.MsgCreateValidator); !ok {
			return genesisState, fmt.Errorf(
				"Genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}

	for _, acc := range genesisState.Accounts {
		for _, coin := range acc.Coins {
			if coin.Denom == genesisState.StakingData.Params.BondDenom {
				stakingData.Pool.NotBondedTokens = stakingData.Pool.NotBondedTokens.
					Add(coin.Amount) // increase the supply
			}
		}
	}

	genesisState.StakingData = stakingData
	genesisState.GenTxs = appGenTxs

	return genesisState, nil
}

// GenesisAccount defines an account initialized at genesis.
type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting"`  // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free"`    // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64     `json:"start_time"`        // vesting start time (UNIX Epoch time)
	EndTime          int64     `json:"end_time"`          // vesting end time (UNIX Epoch time)
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() auth.Account {
	bacc := &auth.BaseAccount{
		Address:       ga.Address,
		Coins:         ga.Coins.Sort(),
		AccountNumber: ga.AccountNumber,
		Sequence:      ga.Sequence,
	}

	if !ga.OriginalVesting.IsZero() {
		baseVestingAcc := &auth.BaseVestingAccount{
			BaseAccount:      bacc,
			OriginalVesting:  ga.OriginalVesting,
			DelegatedFree:    ga.DelegatedFree,
			DelegatedVesting: ga.DelegatedVesting,
			EndTime:          ga.EndTime,
		}

		if ga.StartTime != 0 && ga.EndTime != 0 {
			return &auth.ContinuousVestingAccount{
				BaseVestingAccount: baseVestingAcc,
				StartTime:          ga.StartTime,
			}
		} else if ga.EndTime != 0 {
			return &auth.DelayedVestingAccount{
				BaseVestingAccount: baseVestingAcc,
			}
		} else {
			panic(fmt.Sprintf("invalid genesis vesting account: %+v", ga))
		}
	}

	return bacc
}
func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
}

func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	vacc, ok := acc.(auth.VestingAccount)
	if ok {
		gacc.OriginalVesting = vacc.GetOriginalVesting()
		gacc.DelegatedFree = vacc.GetDelegatedFree()
		gacc.DelegatedVesting = vacc.GetDelegatedVesting()
		gacc.StartTime = vacc.GetStartTime()
		gacc.EndTime = vacc.GetEndTime()
	}

	return gacc
}

func defaultCETAccount() []auth.Account {

	accs := make([]auth.Account, defaultGenesisAccountNum)

	ba1 := auth.NewBaseAccountWithAddress(testAddress1)
	ba1.SetCoins(types.NewCetCoins(2888000000))
	accs[0] = &ba1

	ba2 := auth.NewBaseAccountWithAddress(testAddress2)
	ba2.SetCoins(types.NewCetCoins(1200000000))
	accs[1] = &ba2

	ba3 := auth.NewBaseAccountWithAddress(testAddress3)
	ba3.SetCoins(types.NewCetCoins(360000000))
	accs[2] = auth.NewDelayedVestingAccount(&ba3, 1577836800)

	ba4 := auth.NewBaseAccountWithAddress(testAddress4)
	ba4.SetCoins(types.NewCetCoins(360000000))
	accs[3] = auth.NewDelayedVestingAccount(&ba4, 1609459200)

	ba5 := auth.NewBaseAccountWithAddress(testAddress5)
	ba5.SetCoins(types.NewCetCoins(360000000))
	accs[4] = auth.NewDelayedVestingAccount(&ba5, 1640995200)

	ba6 := auth.NewBaseAccountWithAddress(testAddress6)
	ba6.SetCoins(types.NewCetCoins(360000000))
	accs[5] = auth.NewDelayedVestingAccount(&ba6, 1672531200)

	ba7 := auth.NewBaseAccountWithAddress(testAddress7)
	ba7.SetCoins(types.NewCetCoins(360000000))
	accs[6] = auth.NewDelayedVestingAccount(&ba7, 1704067200)

	return accs
}

func DefaultCETGenesisAccount() []GenesisAccount {
	accs := defaultCETAccount()
	genesisAccs := make([]GenesisAccount, len(accs))

	for i, acc := range accs {
		bacc, ok := acc.(*auth.BaseAccount)
		if ok {
			genesisAccs[i] = NewGenesisAccount(bacc)
			continue
		}

		genesisAccs[i] = NewGenesisAccountI(acc)
	}

	return genesisAccs
}
