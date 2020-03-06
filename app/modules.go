package app

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"

	dex "github.com/coinexchain/cet-sdk/types"
)

type AuthModuleBasic struct {
	auth.AppModuleBasic
}

func (amb AuthModuleBasic) DefaultGenesis() json.RawMessage {
	return auth.ModuleCdc.MustMarshalJSON(GetDefaultAuthGenesisState())
}

func GetDefaultAuthGenesisState() auth.GenesisState {
	genState := auth.DefaultGenesisState()
	genState.Params.MaxMemoCharacters = DefaultMaxMemoCharacters
	genState.Params.TxSizeCostPerByte = DefaultTxSizeCostPerByte
	genState.Params.SigVerifyCostED25519 = DefaultSigVerifyCostED25519
	genState.Params.SigVerifyCostSecp256k1 = DefaultSigVerifyCostSecp256k1
	return genState
}

type StakingModuleBasic struct {
	staking.AppModuleBasic
}

func (StakingModuleBasic) DefaultGenesis() json.RawMessage {
	genState := staking.DefaultGenesisState()
	genState.Params.UnbondingTime = DefaultUnbondingTime
	genState.Params.MaxValidators = DefaultMaxValidators
	genState.Params.BondDenom = dex.DefaultBondDenom
	return staking.ModuleCdc.MustMarshalJSON(genState)
}

type SlashingModuleBasic struct {
	slashing.AppModuleBasic
}

func (SlashingModuleBasic) DefaultGenesis() json.RawMessage {
	genState := slashing.DefaultGenesisState()
	genState.Params.MaxEvidenceAge = DefaultMaxEvidenceAge
	genState.Params.SignedBlocksWindow = DefaultSignedBlocksWindow
	genState.Params.MinSignedPerWindow = DefaultMinSignedPerWindow
	genState.Params.SlashFractionDoubleSign = DefaultSlashFractionDoubleSign
	genState.Params.SlashFractionDowntime = DefaultSlashFractionDowntime
	return slashing.ModuleCdc.MustMarshalJSON(genState)
}

type GovModuleBasic struct {
	gov.AppModuleBasic
}

func (GovModuleBasic) DefaultGenesis() json.RawMessage {
	genState := gov.DefaultGenesisState()
	genState.DepositParams.MinDeposit[0].Denom = dex.DefaultBondDenom
	genState.DepositParams.MinDeposit[0].Amount = DefaultGovMinDeposit
	genState.DepositParams.MaxDepositPeriod = DefaultPeriod
	genState.VotingParams.VotingPeriod = VotingPeriod
	genState.TallyParams = gov.TallyParams{
		Quorum:    sdk.NewDecWithPrec(4, 1),
		Threshold: sdk.NewDecWithPrec(5, 1),
		Veto:      sdk.NewDecWithPrec(334, 3),
	}
	return gov.ModuleCdc.MustMarshalJSON(genState)
}

type CrisisModuleBasic struct {
	crisis.AppModuleBasic
}

func (CrisisModuleBasic) DefaultGenesis() json.RawMessage {
	genState := crisis.DefaultGenesisState()
	genState.ConstantFee.Denom = dex.DefaultBondDenom
	genState.ConstantFee.Amount = DefaultCrisisConstantFee
	return crisis.ModuleCdc.MustMarshalJSON(genState)
}
