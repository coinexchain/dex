package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type MergedParams struct {
	MaxMemoCharacters      uint64  `json:"max_memo_characters" yaml:"max_memo_characters"`
	TxSigLimit             uint64  `json:"tx_sig_limit" yaml:"tx_sig_limit"`
	TxSizeCostPerByte      uint64  `json:"tx_size_cost_per_byte" yaml:"tx_size_cost_per_byte"`
	SigVerifyCostED25519   uint64  `json:"sig_verify_cost_ed25519" yaml:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256k1 uint64  `json:"sig_verify_cost_secp256k1" yaml:"sig_verify_cost_secp256k1"`
	MinGasPriceLimit       sdk.Dec `json:"min_gas_price_limit" yaml:"min_gas_price_limit"`
}

func NewMergedParams(params auth.Params, paramsx Params) MergedParams {
	return MergedParams{
		MaxMemoCharacters:      params.MaxMemoCharacters,
		TxSigLimit:             params.TxSigLimit,
		TxSizeCostPerByte:      params.TxSizeCostPerByte,
		SigVerifyCostED25519:   params.SigVerifyCostED25519,
		SigVerifyCostSecp256k1: params.SigVerifyCostSecp256k1,
		MinGasPriceLimit:       paramsx.MinGasPriceLimit,
	}
}

// String implements the stringer interface.
func (p MergedParams) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("MaxMemoCharacters: %d\n", p.MaxMemoCharacters))
	sb.WriteString(fmt.Sprintf("TxSigLimit: %d\n", p.TxSigLimit))
	sb.WriteString(fmt.Sprintf("TxSizeCostPerByte: %d\n", p.TxSizeCostPerByte))
	sb.WriteString(fmt.Sprintf("SigVerifyCostED25519: %d\n", p.SigVerifyCostED25519))
	sb.WriteString(fmt.Sprintf("SigVerifyCostSecp256k1: %d\n", p.SigVerifyCostSecp256k1))
	sb.WriteString(fmt.Sprintf("MinGasPriceLimit: %s\n", p.MinGasPriceLimit))
	return sb.String()
}
