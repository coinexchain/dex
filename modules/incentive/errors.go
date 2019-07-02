package incentive

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	CodeSpaceIncentive sdk.CodespaceType = "incentive"

	// 701 ～ 799
	CodeInvalidDefaultRewardPerBlock sdk.CodeType = 701
	CodeInvalidPlanHeight            sdk.CodeType = 702
	CodeInvalidRewardPerBlock        sdk.CodeType = 703
	CodeInvalidTotalIncentive        sdk.CodeType = 704
)
