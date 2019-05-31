package asset

import (
	"bytes"
	"fmt"

	"github.com/coinexchain/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default asset module parameter subspace
const (
	DefaultParamspace = ModuleName
	MaxTokenAmount    = 9E18 // 90 billion * 10 ^ 8

	IssueTokenFee                 = 1E12 // 10000 * 10 ^8
	TransferOwnershipFee          = 1E9  // 10 * 10 ^ 8
	ForbidAddrFee                 = 1E9  // 10 * 10 ^ 8
	UnForbidAddrFee               = 1E9  // 10 * 10 ^ 8
	ForbidTokenFee                = 1E9  // 10 * 10 ^ 8
	UnForbidTokenFee              = 1E9  // 10 * 10 ^ 8
	TokenForbidWhitelistAddFee    = 2E9  // 20 * 10 ^ 8
	TokenForbidWhitelistRemoveFee = 2E9  // 20 * 10 ^ 8
	BurnFee                       = 1E9  // 10 * 10 ^ 8
	MintFee                       = 1E9  // 10 * 10 ^ 8
)

// Parameter keys
var (
	KeyIssueTokenFee                 = []byte("IssueTokenFee")
	KeyTransferOwnershipFee          = []byte("TransferOwnershipFee")
	KeyForbidAddrFee                 = []byte("ForbidAddrFee")
	KeyUnForbidAddrFee               = []byte("UnForbidAddrFee")
	KeyForbidTokenFee                = []byte("ForbidTokenFee")
	KeyUnForbidTokenFee              = []byte("UnForbidTokenFee")
	KeyTokenForbidWhitelistAddFee    = []byte("TokenForbidWhitelistAddFee")
	KeyTokenForbidWhitelistRemoveFee = []byte("TokenForbidWhitelistRemoveFee")
	KeyBurnFee                       = []byte("BurnFee")
	KeyMintFee                       = []byte("MintFee")
)

var _ params.ParamSet = &Params{}

// Params defines the parameters for the asset module.
type Params struct {
	// FeeParams define the rules according to which fee are charged.
	IssueTokenFee                 sdk.Coins `json:"issue_token_fee"`
	TransferOwnershipFee          sdk.Coins `json:"transfer_ownership_fee"`
	ForbidAddrFee                 sdk.Coins `json:"forbid_address_fee"`
	UnForbidAddrFee               sdk.Coins `json:"unforbid_address_fee"`
	ForbidTokenFee                sdk.Coins `json:"forbid_token_fee"`
	UnForbidTokenFee              sdk.Coins `json:"unforbid_token_fee"`
	TokenForbidWhitelistAddFee    sdk.Coins `json:"token_forbid_whitelist_add_fee"`
	TokenForbidWhitelistRemoveFee sdk.Coins `json:"token_forbid_whitelist_remove_fee"`
	BurnFee                       sdk.Coins `json:"burn_fee"`
	MintFee                       sdk.Coins `json:"mint_fee"`
}

// ParamKeyTable for asset module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyIssueTokenFee, Value: &p.IssueTokenFee},
		{Key: KeyTransferOwnershipFee, Value: &p.TransferOwnershipFee},
		{Key: KeyForbidAddrFee, Value: &p.ForbidAddrFee},
		{Key: KeyUnForbidAddrFee, Value: &p.UnForbidAddrFee},
		{Key: KeyForbidTokenFee, Value: &p.ForbidTokenFee},
		{Key: KeyUnForbidTokenFee, Value: &p.UnForbidTokenFee},
		{Key: KeyTokenForbidWhitelistAddFee, Value: &p.TokenForbidWhitelistAddFee},
		{Key: KeyTokenForbidWhitelistRemoveFee, Value: &p.TokenForbidWhitelistRemoveFee},
		{Key: KeyBurnFee, Value: &p.BurnFee},
		{Key: KeyMintFee, Value: &p.MintFee},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := msgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := msgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func (p *Params) ValidateGenesis() error {
	for _, pair := range p.ParamSetPairs() {
		fee := pair.Value.(*sdk.Coins)
		if fee.Empty() || fee.IsAnyNegative() {
			return fmt.Errorf("%s must be a valid sdk.Coins, is %s", pair.Key, fee.String())
		}
	}
	return nil
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {

	return Params{
		types.NewCetCoins(IssueTokenFee),
		types.NewCetCoins(TransferOwnershipFee),
		types.NewCetCoins(ForbidAddrFee),
		types.NewCetCoins(UnForbidAddrFee),
		types.NewCetCoins(ForbidTokenFee),
		types.NewCetCoins(UnForbidTokenFee),
		types.NewCetCoins(TokenForbidWhitelistAddFee),
		types.NewCetCoins(TokenForbidWhitelistRemoveFee),
		types.NewCetCoins(BurnFee),
		types.NewCetCoins(MintFee),
	}
}
