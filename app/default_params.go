package app

import "time"

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
