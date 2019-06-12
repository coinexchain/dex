package testutil

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	dex "github.com/coinexchain/dex/types"
)

type StdTxBuilder struct {
	chainID  string
	msgs     []sdk.Msg
	privKeys []crypto.PrivKey
	accNums  []uint64
	seqs     []uint64
	fee      auth.StdFee
}

func NewStdTxBuilder(chainID string) *StdTxBuilder {
	return &StdTxBuilder{chainID: chainID}
}

func (builder *StdTxBuilder) Msgs(msgs ...sdk.Msg) *StdTxBuilder {
	builder.msgs = msgs
	return builder
}
func (builder *StdTxBuilder) AccNumSeqKey(num, seq uint64, key crypto.PrivKey) *StdTxBuilder {
	builder.accNums = append(builder.accNums, num)
	builder.seqs = append(builder.seqs, seq)
	builder.privKeys = append(builder.privKeys, key)
	return builder
}
func (builder *StdTxBuilder) GasAndFee(gas uint64, cet int64) *StdTxBuilder {
	builder.fee = auth.NewStdFee(gas, dex.NewCetCoins(cet))
	return builder
}

func (builder *StdTxBuilder) Build() auth.StdTx {
	sigs := make([]auth.StdSignature, len(builder.privKeys))
	for i, privKey := range builder.privKeys {
		signBytes := auth.StdSignBytes(builder.chainID,
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
