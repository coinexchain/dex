package incentive

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = "incentive"
)

var _ params.ParamSet = &Params{}

var (
	KeyIncentiveBlockInterval         = []byte("incentiveBlockInterval")
	KeyIncentiveDefaultRewardPerBlock = []byte("incentiveDefaultRewardPerBlock")
	KeyIncentivePlans                 = []byte("incentivePlans")
)

type Params struct {
	IncentiveBlockInterval uint16 `json:"incentive_block_interval"`
	DefaultRewardPerBlock  uint16 `json:"default_reward_per_block"`
	Plans                  []Plan `json:"plans"`
}

type Plan struct {
	StartHeight    int64 `json:"start_height"`
	EndHeight      int64 `json:"end_height"`
	RewardPerBlock int64 `json:"reward_per_block"`
	TotalIncentive int64 `json:"total_incentive"`
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyIncentiveBlockInterval, Value: &p.IncentiveBlockInterval},
		{Key: KeyIncentiveDefaultRewardPerBlock, Value: &p.DefaultRewardPerBlock},
		{Key: KeyIncentivePlans, Value: &p.Plans},
	}
}

func DefaultParams() Params {
	return Params{
		IncentiveBlockInterval: 3,
		DefaultRewardPerBlock:  0,
		Plans: []Plan{
			{0, 10512000, 10, 105120000},
			{10512000, 21024000, 8, 84096000},
			{21024000, 31536000, 6, 63072000},
			{31536000, 42048000, 4, 42048000},
			{42048000, 52560000, 2, 21024000},
		},
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}
