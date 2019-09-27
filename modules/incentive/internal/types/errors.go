package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceIncentive sdk.CodespaceType = "incentive"

	// 701 ～ 799
	CodeInvalidAdjustmentHeight      sdk.CodeType = 701
	CodeInvalidDefaultRewardPerBlock sdk.CodeType = 702
	CodeInvalidPlanHeight            sdk.CodeType = 703
	CodeInvalidRewardPerBlock        sdk.CodeType = 704
	CodeInvalidTotalIncentive        sdk.CodeType = 705
	CodeInvalidPlanToAdd             sdk.CodeType = 706
)
