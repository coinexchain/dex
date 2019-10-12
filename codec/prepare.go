package codec

import (
	"io"

	"github.com/coinexchain/codon"
)

func ShowInfo() {
	leafTypes := GetLeafTypes()

	//ShowInfo("",Account{})

	codon.ShowInfoForVar(leafTypes, DuplicateVoteEvidence{})
	codon.ShowInfoForVar(leafTypes, PrivKeyEd25519{})
	codon.ShowInfoForVar(leafTypes, PrivKeySecp256k1{})
	codon.ShowInfoForVar(leafTypes, PubKeyEd25519{})
	codon.ShowInfoForVar(leafTypes, PubKeySecp256k1{})
	codon.ShowInfoForVar(leafTypes, PubKeyMultisigThreshold{})

	codon.ShowInfoForVar(leafTypes, BaseVestingAccount{})
	codon.ShowInfoForVar(leafTypes, ContinuousVestingAccount{})
	codon.ShowInfoForVar(leafTypes, DelayedVestingAccount{})
	codon.ShowInfoForVar(leafTypes, ModuleAccount{})
	codon.ShowInfoForVar(leafTypes, StdTx{})
	codon.ShowInfoForVar(leafTypes, MsgBeginRedelegate{})
	codon.ShowInfoForVar(leafTypes, MsgCreateValidator{})
	codon.ShowInfoForVar(leafTypes, MsgDelegate{})
	codon.ShowInfoForVar(leafTypes, MsgEditValidator{})
	codon.ShowInfoForVar(leafTypes, MsgSetWithdrawAddress{})
	codon.ShowInfoForVar(leafTypes, MsgUndelegate{})
	codon.ShowInfoForVar(leafTypes, MsgUnjail{})
	codon.ShowInfoForVar(leafTypes, MsgWithdrawDelegatorReward{})
	codon.ShowInfoForVar(leafTypes, MsgWithdrawValidatorCommission{})
	codon.ShowInfoForVar(leafTypes, MsgDeposit{})
	codon.ShowInfoForVar(leafTypes, MsgSubmitProposal{})
	codon.ShowInfoForVar(leafTypes, MsgVote{})
	codon.ShowInfoForVar(leafTypes, ParameterChangeProposal{})
	codon.ShowInfoForVar(leafTypes, SoftwareUpgradeProposal{})
	codon.ShowInfoForVar(leafTypes, TextProposal{})
	codon.ShowInfoForVar(leafTypes, CommunityPoolSpendProposal{})
	codon.ShowInfoForVar(leafTypes, MsgMultiSend{})
	codon.ShowInfoForVar(leafTypes, MsgSend{})
	codon.ShowInfoForVar(leafTypes, MsgVerifyInvariant{})
	codon.ShowInfoForVar(leafTypes, Supply{})

	codon.ShowInfoForVar(leafTypes, AccountX{})
	codon.ShowInfoForVar(leafTypes, MsgMultiSendX{})
	codon.ShowInfoForVar(leafTypes, MsgSendX{})
	codon.ShowInfoForVar(leafTypes, MsgSetMemoRequired{})
	codon.ShowInfoForVar(leafTypes, BaseToken{})
	codon.ShowInfoForVar(leafTypes, MsgAddTokenWhitelist{})
	codon.ShowInfoForVar(leafTypes, MsgBurnToken{})
	codon.ShowInfoForVar(leafTypes, MsgForbidAddr{})
	codon.ShowInfoForVar(leafTypes, MsgForbidToken{})
	codon.ShowInfoForVar(leafTypes, MsgIssueToken{})
	codon.ShowInfoForVar(leafTypes, MsgMintToken{})
	codon.ShowInfoForVar(leafTypes, MsgModifyTokenInfo{})
	codon.ShowInfoForVar(leafTypes, MsgRemoveTokenWhitelist{})
	codon.ShowInfoForVar(leafTypes, MsgTransferOwnership{})
	codon.ShowInfoForVar(leafTypes, MsgUnForbidAddr{})
	codon.ShowInfoForVar(leafTypes, MsgUnForbidToken{})
	codon.ShowInfoForVar(leafTypes, MsgBancorCancel{})
	codon.ShowInfoForVar(leafTypes, MsgBancorInit{})
	codon.ShowInfoForVar(leafTypes, MsgBancorTrade{})
	codon.ShowInfoForVar(leafTypes, MsgCancelOrder{})
	codon.ShowInfoForVar(leafTypes, MsgCancelTradingPair{})
	codon.ShowInfoForVar(leafTypes, MsgCreateOrder{})
	codon.ShowInfoForVar(leafTypes, MsgCreateTradingPair{})
	codon.ShowInfoForVar(leafTypes, MsgModifyPricePrecision{})
	codon.ShowInfoForVar(leafTypes, Order{})
	codon.ShowInfoForVar(leafTypes, MarketInfo{})
	codon.ShowInfoForVar(leafTypes, &MsgDonateToCommunityPool{})
	codon.ShowInfoForVar(leafTypes, &MsgCommentToken{})
	codon.ShowInfoForVar(leafTypes, &State{})
	codon.ShowInfoForVar(leafTypes, &MsgAliasUpdate{})
}

func GenerateCodecFile(w io.Writer) {
	list := []codon.AliasAndValue{
		{Alias: "PubKey", Value: (*PubKey)(nil)},
		{Alias: "Msg", Value: (*Msg)(nil)},
		{Alias: "Account", Value: (*Account)(nil)},
		{Alias: "Content", Value: (*Content)(nil)},

		{Alias: "DuplicateVoteEvidence", Value: DuplicateVoteEvidence{}},
		{Alias: "PrivKeyEd25519", Value: PrivKeyEd25519{}},
		{Alias: "PrivKeySecp256k1", Value: PrivKeySecp256k1{}},
		{Alias: "PubKeyEd25519", Value: PubKeyEd25519{}},
		{Alias: "PubKeySecp256k1", Value: PubKeySecp256k1{}},
		{Alias: "PubKeyMultisigThreshold", Value: PubKeyMultisigThreshold{}},
		{Alias: "SignedMsgType", Value: SignedMsgType(0)},
		{Alias: "VoteOption", Value: VoteOption(0)},
		{Alias: "Vote", Value: Vote{}},

		{Alias: "Coin", Value: Coin{}},
		{Alias: "LockedCoin", Value: LockedCoin{}},
		{Alias: "StdSignature", Value: StdSignature{}},
		{Alias: "ParamChange", Value: ParamChange{}},
		{Alias: "Input", Value: Input{}},
		{Alias: "Output", Value: Output{}},
		{Alias: "AccAddress", Value: AccAddress{}},
		{Alias: "CommentRef", Value: CommentRef{}},

		{Alias: "BaseAccount", Value: BaseAccount{}},
		{Alias: "BaseVestingAccount", Value: BaseVestingAccount{}},
		{Alias: "ContinuousVestingAccount", Value: ContinuousVestingAccount{}},
		{Alias: "DelayedVestingAccount", Value: DelayedVestingAccount{}},
		{Alias: "ModuleAccount", Value: ModuleAccount{}},
		{Alias: "StdTx", Value: StdTx{}},
		{Alias: "MsgBeginRedelegate", Value: MsgBeginRedelegate{}},
		{Alias: "MsgCreateValidator", Value: MsgCreateValidator{}},
		{Alias: "MsgDelegate", Value: MsgDelegate{}},
		{Alias: "MsgEditValidator", Value: MsgEditValidator{}},
		{Alias: "MsgSetWithdrawAddress", Value: MsgSetWithdrawAddress{}},
		{Alias: "MsgUndelegate", Value: MsgUndelegate{}},
		{Alias: "MsgUnjail", Value: MsgUnjail{}},
		{Alias: "MsgWithdrawDelegatorReward", Value: MsgWithdrawDelegatorReward{}},
		{Alias: "MsgWithdrawValidatorCommission", Value: MsgWithdrawValidatorCommission{}},
		{Alias: "MsgDeposit", Value: MsgDeposit{}},
		{Alias: "MsgSubmitProposal", Value: MsgSubmitProposal{}},
		{Alias: "MsgVote", Value: MsgVote{}},
		{Alias: "ParameterChangeProposal", Value: ParameterChangeProposal{}},
		{Alias: "SoftwareUpgradeProposal", Value: SoftwareUpgradeProposal{}},
		{Alias: "TextProposal", Value: TextProposal{}},
		{Alias: "CommunityPoolSpendProposal", Value: CommunityPoolSpendProposal{}},
		{Alias: "MsgMultiSend", Value: MsgMultiSend{}},
		{Alias: "MsgSend", Value: MsgSend{}},
		{Alias: "MsgVerifyInvariant", Value: MsgVerifyInvariant{}},
		{Alias: "Supply", Value: Supply{}},

		{Alias: "AccountX", Value: AccountX{}},
		{Alias: "MsgMultiSendX", Value: MsgMultiSendX{}},
		{Alias: "MsgSendX", Value: MsgSendX{}},
		{Alias: "MsgSetMemoRequired", Value: MsgSetMemoRequired{}},
		{Alias: "BaseToken", Value: BaseToken{}},
		{Alias: "MsgAddTokenWhitelist", Value: MsgAddTokenWhitelist{}},
		{Alias: "MsgBurnToken", Value: MsgBurnToken{}},
		{Alias: "MsgForbidAddr", Value: MsgForbidAddr{}},
		{Alias: "MsgForbidToken", Value: MsgForbidToken{}},
		{Alias: "MsgIssueToken", Value: MsgIssueToken{}},
		{Alias: "MsgMintToken", Value: MsgMintToken{}},
		{Alias: "MsgModifyTokenInfo", Value: MsgModifyTokenInfo{}},
		{Alias: "MsgRemoveTokenWhitelist", Value: MsgRemoveTokenWhitelist{}},
		{Alias: "MsgTransferOwnership", Value: MsgTransferOwnership{}},
		{Alias: "MsgUnForbidAddr", Value: MsgUnForbidAddr{}},
		{Alias: "MsgUnForbidToken", Value: MsgUnForbidToken{}},
		{Alias: "MsgBancorCancel", Value: MsgBancorCancel{}},
		{Alias: "MsgBancorInit", Value: MsgBancorInit{}},
		{Alias: "MsgBancorTrade", Value: MsgBancorTrade{}},
		{Alias: "MsgCancelOrder", Value: MsgCancelOrder{}},
		{Alias: "MsgCancelTradingPair", Value: MsgCancelTradingPair{}},
		{Alias: "MsgCreateOrder", Value: MsgCreateOrder{}},
		{Alias: "MsgCreateTradingPair", Value: MsgCreateTradingPair{}},
		{Alias: "MsgModifyPricePrecision", Value: MsgModifyPricePrecision{}},
		{Alias: "Order", Value: Order{}},
		{Alias: "MarketInfo", Value: MarketInfo{}},
		{Alias: "MsgDonateToCommunityPool", Value: MsgDonateToCommunityPool{}},
		{Alias: "MsgCommentToken", Value: MsgCommentToken{}},
		{Alias: "State", Value: State{}},
		{Alias: "MsgAliasUpdate", Value: MsgAliasUpdate{}},
	}

	extraImports := []string{`"time"`, `sdk "github.com/cosmos/cosmos-sdk/types"`}
	ignoreImpl := make(map[string]string)
	ignoreImpl["StdSignature"] = "PubKey"
	ignoreImpl["PubKeyMultisigThreshold"] = "PubKey"
	codon.GenerateCodecFile(w, GetLeafTypes(), ignoreImpl, list, extraLogics, extraImports)
}

func GetLeafTypes() map[string]string {
	leafTypes := make(map[string]string, 20)
	leafTypes["github.com/cosmos/cosmos-sdk/types.Int"] = "sdk.Int"
	leafTypes["github.com/cosmos/cosmos-sdk/types.Dec"] = "sdk.Dec"
	leafTypes["time.Time"] = "time.Time"
	return leafTypes
}

type RandSrc interface {
	GetBool() bool
	GetInt() int
	GetInt8() int8
	GetInt16() int16
	GetInt32() int32
	GetInt64() int64
	GetUint() uint
	GetUint8() uint8
	GetUint16() uint16
	GetUint32() uint32
	GetUint64() uint64
	GetFloat32() float32
	GetFloat64() float64
	GetString(n int) string
	GetBytes(n int) []byte
}

const MaxSliceLength = 10
const MaxStringLength = 100

var extraLogics = `
func EncodeTime(w io.Writer, t time.Time) error {
	t = t.UTC()
	sec := t.Unix()
	var buf [10]byte
	n := binary.PutVarint(buf[:], sec)
	_, err := w.Write(buf[0:n])
	if err != nil {
		return err
	}

	nanosec := t.Nanosecond()
	n = binary.PutVarint(buf[:], int64(nanosec))
	_, err = w.Write(buf[0:n])
	if err != nil {
		return err
	}
	return nil
}

func DecodeTime(bz []byte) (time.Time, int, error) {
	sec, n := binary.Varint(bz)
	var err error
	if n == 0 {
		// buf too small
		err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		err = errors.New("EOF decoding varint")
	}
	if err!=nil {
		return time.Unix(sec,0), n, err
	}

	nanosec, m := binary.Varint(bz[n:])
	if m == 0 {
		// buf too small
		err = errors.New("buffer too small")
	} else if m < 0 {
		// value larger than 64 bits (overflow)
		// and -m is the number of bytes read
		m = -m
		err = errors.New("EOF decoding varint")
	}
	if err!=nil {
		return time.Unix(sec,nanosec), n+m, err
	}

	return time.Unix(sec, nanosec).UTC(), n+m, nil
}

func RandTime(r RandSrc) time.Time {
	return time.Unix(r.GetInt64(), r.GetInt64()).UTC()
}

func EncodeInt(w io.Writer, v sdk.Int) error {
	s, err := v.MarshalAmino()
	if err!=nil {
		return err
	}
	return codonEncodeString(w, s)
}

func DecodeInt(bz []byte) (sdk.Int, int, error) {
	v := sdk.ZeroInt()
	var n int
	var err error
	s := codonDecodeString(bz, &n, &err)
	if err!=nil {
		return v, n, err
	}

	err = (&v).UnmarshalAmino(s)
	if err!=nil {
		return v, n, err
	}

	return v, n, nil
}

func RandInt(r RandSrc) sdk.Int {
	res := sdk.NewInt(r.GetInt64())
	count := int(r.GetInt64()%3)
	for i:=0; i<count; i++ {
		res = res.MulRaw(r.GetInt64())
	}
	return res
}

func EncodeDec(w io.Writer, v sdk.Dec) error {
	s, err := v.MarshalAmino()
	if err!=nil {
		return err
	}
	return codonEncodeString(w, s)
}

func DecodeDec(bz []byte) (sdk.Dec, int, error) {
	v := sdk.ZeroDec()
	var n int
	var err error
	s := codonDecodeString(bz, &n, &err)
	if err!=nil {
		return v, n, err
	}

	err = (&v).UnmarshalAmino(s)
	if err!=nil {
		return v, n, err
	}

	return v, n, nil
}

func RandDec(r RandSrc) sdk.Dec {
	res := sdk.NewDec(r.GetInt64())
	count := int(r.GetInt64()%3)
	for i:=0; i<count; i++ {
		res = res.MulInt64(r.GetInt64())
	}
	res = res.QuoInt64(r.GetInt64()&0xFFFFFFFF)
	return res
}
`

/*
func EncodeDuplicateVoteEvidence(w io.Writer, v DuplicateVoteEvidence) error {
// codon version: 1
var err error
err = EncodePubKey(w, v.PubKey)
}

func EncodePubKeyMultisigThreshold(w io.Writer, v PubKeyMultisigThreshold) error {
// codon version: 1
var err error
err = codonEncodeUvarint(w, uint64(v.K))
if err != nil {return err}
err = codonEncodeVarint(w, int64(len(v.PubKeys)))
if err != nil {return err}
for _0:=0; _0<len(v.PubKeys); _0++ {
err = EncodePubKey(w, v.PubKeys[_0])
if err != nil {return err} // interface_encode

func EncodeStdSignature(w io.Writer, v StdSignature) error {
// codon version: 1
var err error
err = EncodePubKey(w, v.PubKey)
if err != nil {return err} // interface_encode

func EncodeBaseVestingAccount(w io.Writer, v BaseVestingAccount) error {
err = EncodePubKey(w, v.BaseAccount.PubKey)
if err != nil {return err} // interface_encode

func EncodeContinuousVestingAccount(w io.Writer, v ContinuousVestingAccount) error {
err = EncodePubKey(w, v.BaseVestingAccount.BaseAccount.PubKey)
if err != nil {return err} // interface_encode

func EncodeDelayedVestingAccount(w io.Writer, v DelayedVestingAccount) error {
err = EncodePubKey(w, v.BaseVestingAccount.BaseAccount.PubKey)
if err != nil {return err} // interface_encode

func EncodeModuleAccount(w io.Writer, v ModuleAccount) error {
err = EncodePubKey(w, v.BaseVestingAccount.BaseAccount.PubKey)
if err != nil {return err} // interface_encode

func EncodeStdTx(w io.Writer, v StdTx) error {
for _0:=0; _0<len(v.Msgs); _0++ {
err = EncodeMsg(w, v.Msgs[_0])
if err != nil {return err} // interface_encode
}
for _0:=0; _0<len(v.Signatures); _0++ {
err = EncodePubKey(w, v.Signatures[_0].PubKey)
if err != nil {return err} // interface_encode

func EncodeMsgCreateValidator(w io.Writer, v MsgCreateValidator) error {
err = EncodePubKey(w, v.PubKey)
if err != nil {return err} // interface_encode

*/
