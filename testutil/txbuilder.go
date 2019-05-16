package testutil

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type stdTxBuilder struct {
	chainId  string
	msgs     []sdk.Msg
	privKeys []crypto.PrivKey
	accNums  []uint64
	seqs     []uint64
	fee      auth.StdFee
}

func NewStdTxBuilder(chainId string) *stdTxBuilder {
	return &stdTxBuilder{chainId: chainId}
}

func (builder *stdTxBuilder) Msgs(msgs ...sdk.Msg) *stdTxBuilder {
	builder.msgs = msgs
	return builder
}
func (builder *stdTxBuilder) AccNumSeqKey(num, seq uint64, key crypto.PrivKey) *stdTxBuilder {
	builder.accNums = append(builder.accNums, num)
	builder.seqs = append(builder.seqs, seq)
	builder.privKeys = append(builder.privKeys, key)
	return builder
}
func (builder *stdTxBuilder) Fee(fee auth.StdFee) *stdTxBuilder {
	builder.fee = fee
	return builder
}

func (builder *stdTxBuilder) Build() auth.StdTx {
	sigs := make([]auth.StdSignature, len(builder.privKeys))
	for i, privKey := range builder.privKeys {
		signBytes := auth.StdSignBytes(builder.chainId,
			builder.accNums[i], builder.seqs[i], builder.fee, builder.msgs, "")

		sig, err := privKey.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{PubKey: privKey.PubKey(), Signature: sig}
	}

	tx := auth.NewStdTx(builder.msgs, builder.fee, sigs, "")
	return tx
}
