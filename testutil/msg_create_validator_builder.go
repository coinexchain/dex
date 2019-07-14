package testutil

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	dex "github.com/coinexchain/dex/types"
)

type MsgCreateValidatorBuilder struct {
	description       staking.Description
	commission        staking.CommissionRates
	minSelfDelegation sdk.Int
	valAddr           sdk.ValAddress
	pubKey            crypto.PubKey
	selfDelegation    sdk.Coin
}

func NewMsgCreateValidatorBuilder(valAddr sdk.ValAddress, pubKey crypto.PubKey) *MsgCreateValidatorBuilder {
	return &MsgCreateValidatorBuilder{
		description: staking.NewDescription("node", "node", "www.node.org", "node"),
		commission:  staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		valAddr:     valAddr,
		pubKey:      pubKey,
	}
}

func (builder *MsgCreateValidatorBuilder) Description(moniker, identity, website, details string) *MsgCreateValidatorBuilder {
	builder.description = staking.NewDescription(moniker, identity, website, details)
	return builder
}
func (builder *MsgCreateValidatorBuilder) Commission(rate, maxRate, maxChangeRate string) *MsgCreateValidatorBuilder {
	builder.commission = staking.NewCommissionRates(
		sdk.MustNewDecFromStr(rate),
		sdk.MustNewDecFromStr(maxRate),
		sdk.MustNewDecFromStr(maxChangeRate),
	)
	return builder
}
func (builder *MsgCreateValidatorBuilder) MinSelfDelegation(minSelfDelegation int64) *MsgCreateValidatorBuilder {
	builder.minSelfDelegation = sdk.NewInt(minSelfDelegation)
	return builder
}
func (builder *MsgCreateValidatorBuilder) SelfDelegation(selfDelegation int64) *MsgCreateValidatorBuilder {
	builder.selfDelegation = dex.NewCetCoin(selfDelegation)
	return builder
}

func (builder *MsgCreateValidatorBuilder) Build() staking.MsgCreateValidator {
	return staking.NewMsgCreateValidator(builder.valAddr, builder.pubKey,
		builder.selfDelegation, builder.description, builder.commission, builder.minSelfDelegation)
}
