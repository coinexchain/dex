package codec

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ptypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/coinexchain/cet-sdk/modules/alias"
	"github.com/coinexchain/cet-sdk/modules/asset"
	"github.com/coinexchain/cet-sdk/modules/authx"
	"github.com/coinexchain/cet-sdk/modules/bancorlite"
	"github.com/coinexchain/cet-sdk/modules/bankx"
	"github.com/coinexchain/cet-sdk/modules/comment"
	distrx "github.com/coinexchain/cet-sdk/modules/distributionx"
	"github.com/coinexchain/cet-sdk/modules/incentive"
	"github.com/coinexchain/cet-sdk/modules/market"
)

type (
	PubKey  = crypto.PubKey
	Msg     = sdk.Msg
	Account = auth.Account
	Content = govtypes.Content

	DuplicateVoteEvidence   = tmtypes.DuplicateVoteEvidence
	PrivKeyEd25519          = ed25519.PrivKeyEd25519
	PrivKeySecp256k1        = secp256k1.PrivKeySecp256k1
	PubKeyEd25519           = ed25519.PubKeyEd25519
	PubKeySecp256k1         = secp256k1.PubKeySecp256k1
	PubKeyMultisigThreshold = multisig.PubKeyMultisigThreshold
	SignedMsgType           = tmtypes.SignedMsgType
	VoteOption              = govtypes.VoteOption
	Vote                    = tmtypes.Vote

	Int = sdk.Int
	Dec = sdk.Dec

	Coin         = sdk.Coin
	StdSignature = auth.StdSignature
	ParamChange  = ptypes.ParamChange
	Input        = bank.Input
	Output       = bank.Output
	LockedCoin   = authx.LockedCoin
	AccAddress   = sdk.AccAddress
	CommentRef   = comment.CommentRef

	BaseAccount                    = auth.BaseAccount
	BaseVestingAccount             = auth.BaseVestingAccount
	ContinuousVestingAccount       = auth.ContinuousVestingAccount
	DelayedVestingAccount          = auth.DelayedVestingAccount
	StdTx                          = auth.StdTx
	MsgBeginRedelegate             = staking.MsgBeginRedelegate
	MsgCreateValidator             = staking.MsgCreateValidator
	MsgDelegate                    = staking.MsgDelegate
	MsgEditValidator               = staking.MsgEditValidator
	MsgUndelegate                  = staking.MsgUndelegate
	MsgUnjail                      = slashing.MsgUnjail
	MsgSetWithdrawAddress          = distr.MsgSetWithdrawAddress
	MsgWithdrawDelegatorReward     = distr.MsgWithdrawDelegatorReward
	MsgWithdrawValidatorCommission = distr.MsgWithdrawValidatorCommission
	MsgDeposit                     = gov.MsgDeposit
	MsgSubmitProposal              = gov.MsgSubmitProposal
	MsgVote                        = gov.MsgVote
	SoftwareUpgradeProposal        = gov.SoftwareUpgradeProposal
	TextProposal                   = gov.TextProposal
	ParameterChangeProposal        = ptypes.ParameterChangeProposal
	CommunityPoolSpendProposal     = distr.CommunityPoolSpendProposal
	MsgMultiSend                   = bank.MsgMultiSend
	MsgSend                        = bank.MsgSend
	MsgVerifyInvariant             = crisis.MsgVerifyInvariant
	Supply                         = supply.Supply
	ModuleAccount                  = supply.ModuleAccount

	AccountX                 = authx.AccountX
	MsgMultiSendX            = bankx.MsgMultiSend
	MsgSendX                 = bankx.MsgSend
	MsgSetMemoRequired       = bankx.MsgSetMemoRequired
	BaseToken                = asset.BaseToken
	MsgAddTokenWhitelist     = asset.MsgAddTokenWhitelist
	MsgBurnToken             = asset.MsgBurnToken
	MsgForbidAddr            = asset.MsgForbidAddr
	MsgForbidToken           = asset.MsgForbidToken
	MsgIssueToken            = asset.MsgIssueToken
	MsgMintToken             = asset.MsgMintToken
	MsgModifyTokenInfo       = asset.MsgModifyTokenInfo
	MsgRemoveTokenWhitelist  = asset.MsgRemoveTokenWhitelist
	MsgTransferOwnership     = asset.MsgTransferOwnership
	MsgUnForbidAddr          = asset.MsgUnForbidAddr
	MsgUnForbidToken         = asset.MsgUnForbidToken
	MsgBancorCancel          = bancorlite.MsgBancorCancel
	MsgBancorInit            = bancorlite.MsgBancorInit
	MsgBancorTrade           = bancorlite.MsgBancorTrade
	MsgCancelOrder           = market.MsgCancelOrder
	MsgCancelTradingPair     = market.MsgCancelTradingPair
	MsgCreateOrder           = market.MsgCreateOrder
	MsgCreateTradingPair     = market.MsgCreateTradingPair
	MsgModifyPricePrecision  = market.MsgModifyPricePrecision
	Order                    = market.Order
	MarketInfo               = market.MarketInfo
	MsgDonateToCommunityPool = distrx.MsgDonateToCommunityPool
	MsgCommentToken          = comment.MsgCommentToken
	State                    = incentive.State
	MsgAliasUpdate           = alias.MsgAliasUpdate
)
