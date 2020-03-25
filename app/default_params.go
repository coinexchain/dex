package app

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const Day = 24 * time.Hour

// auth
const (
	DefaultMaxMemoCharacters      uint64 = 512
	DefaultTxSizeCostPerByte      uint64 = 20
	DefaultSigVerifyCostED25519   uint64 = 11800
	DefaultSigVerifyCostSecp256k1 uint64 = 20000
)

// staking
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	DefaultUnbondingTime = 21 * Day

	// Default maximum number of bonded validators
	DefaultMaxValidators uint16 = 42
)

// slashing
const (
	DefaultMaxEvidenceAge           = 21 * Day
	DefaultSignedBlocksWindow int64 = 10000
)

// consensus
const (
	DefaultEvidenceMaxAge int64 = 1000000
)

var (
	DefaultMinSignedPerWindow      = sdk.NewDecWithPrec(5, 2)             // 0.05
	DefaultSlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))    // 0.05
	DefaultSlashFractionDowntime   = sdk.NewDec(1).Quo(sdk.NewDec(10000)) // 0.0001

	DefaultGovMinDeposit = sdk.NewInt(10000e8)

	DefaultCrisisConstantFee = sdk.NewInt(100000e8)
)

// gov
const (
	// Default period for deposits & voting
	DefaultPeriod = 14 * Day // TODO
	VotingPeriod  = 7 * Day
)

// staking
const (
	MinSelfDelegation = 1000000e8
)
