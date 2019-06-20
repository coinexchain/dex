package app

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const Day = 24 * time.Hour

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
	DefaultSignedBlocksWindow int64 = 1000
)

var (
	DefaultMinSignedPerWindow      = sdk.NewDecWithPrec(5, 2)             // 0.05
	DefaultSlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))    // 0.05
	DefaultSlashFractionDowntime   = sdk.NewDec(1).Quo(sdk.NewDec(10000)) // 0.0001
)

// gov
const (
	// Default period for deposits & voting
	DefaultPeriod = 14 * Day
)
