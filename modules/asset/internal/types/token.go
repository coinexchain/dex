package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"unicode/utf8"
)

const (
	MaxTokenIdentityLength    = 3000
	MaxTokenURLLength         = 100
	MaxTokenDescriptionLength = 1024

	// constant used in flags to indicate that token info field should not be updated
	DoNotModifyTokenInfo = "[no-change]"
)

// Token is an interface used to store asset at a given token within state.
// Many complex conditions can be used in the concrete struct which implements Token.
type Token interface {
	SetName(string) sdk.Error

	GetSymbol() string
	SetSymbol(string) sdk.Error

	GetTotalSupply() int64
	SetTotalSupply(int64) sdk.Error

	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress) sdk.Error

	GetMintable() bool
	SetMintable(bool)

	GetBurnable() bool
	SetBurnable(bool)

	GetAddrForbiddable() bool
	SetAddrForbiddable(bool)

	GetTokenForbiddable() bool
	SetTokenForbiddable(bool)

	GetTotalBurn() int64
	SetTotalBurn(int64) sdk.Error

	GetTotalMint() int64
	SetTotalMint(int64) sdk.Error

	GetIsForbidden() bool
	SetIsForbidden(bool)

	GetURL() string
	SetURL(string) sdk.Error

	GetDescription() string
	SetDescription(string) sdk.Error

	GetIdentity() string
	SetIdentity(string) sdk.Error

	Validate() sdk.Error
	// Ensure that token implements stringer
	String() string
}

//-----------------------------------------------------------------------------
var _ Token = (*BaseToken)(nil)

// BaseToken - a base Token structure.
type BaseToken struct {
	Name             string         `json:"name" yaml:"name"`                           //  Name of the newly issued asset, limited to 32 unicode characters.
	Symbol           string         `json:"symbol" yaml:"symbol"`                       //  token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply      int64          `json:"total_supply" yaml:"total_supply"`           //  The total supply for this token [0]
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`                         // The initial issuer of this token
	Mintable         bool           `json:"mintable" yaml:"mintable"`                   // Whether this token could be minted after the issuing
	Burnable         bool           `json:"burnable" yaml:"burnable"`                   // Whether this token could be burned
	AddrForbiddable  bool           `json:"addr_forbiddable" yaml:"addr_forbiddable"`   // whether could forbid some addresses to forbid transaction
	TokenForbiddable bool           `json:"token_forbiddable" yaml:"token_forbiddable"` // whether token could be global forbid
	TotalBurn        int64          `json:"total_burn" yaml:"total_burn"`               // Total amount of burn
	TotalMint        int64          `json:"total_mint" yaml:"total_mint"`               // Total amount of mint
	IsForbidden      bool           `json:"is_forbidden" yaml:"is_forbidden"`           // Whether token being forbidden currently
	URL              string         `json:"url" yaml:"url"`                             //URL of token website
	Description      string         `json:"description" yaml:"description"`             //Description of token info
	Identity         string         `json:"identity" yaml:"identity"`                   //Identity of token
}

//nolint
var (

	// suffixSymbolRegex : Only CET owner can issue .suffix token
	suffixSymbolRegex = regexp.MustCompile("^[a-z][a-z0-9]{1,7}\\.[a-z]$")

	// tokenSymbolRegex : Token symbol can be 2 ~ 8 characters long.
	tokenSymbolRegex = regexp.MustCompile("^[a-z][a-z0-9]{1,7}(\\.[a-z])?$")
)

// NewToken - new base token
func NewToken(name string, symbol string, totalSupply int64, owner sdk.AccAddress,
	mintable bool, burnable bool, addrForbiddable bool, tokenForbiddable bool,
	url string, description string, identity string) (*BaseToken, sdk.Error) {

	t := &BaseToken{}
	var err sdk.Error
	if err = t.SetName(name); err != nil {
		return nil, err
	}
	if err = t.SetOwner(owner); err != nil {
		return nil, err
	}
	if err = t.SetSymbol(symbol); err != nil {
		return nil, err
	}
	if err = t.SetTotalSupply(totalSupply); err != nil {
		return nil, err
	}
	if err = t.SetURL(url); err != nil {
		return nil, err
	}
	if err = t.SetDescription(description); err != nil {
		return nil, err
	}
	if err = t.SetIdentity(identity); err != nil {
		return nil, err
	}

	t.SetMintable(mintable)
	t.SetBurnable(burnable)
	t.SetAddrForbiddable(addrForbiddable)
	t.SetTokenForbiddable(tokenForbiddable)

	if err = t.SetTotalMint(0); err != nil {
		return nil, err
	}
	if err = t.SetTotalBurn(0); err != nil {
		return nil, err
	}
	t.SetIsForbidden(false)

	return t, nil
}

func (t *BaseToken) Validate() sdk.Error {
	_, err := NewToken(t.Name, t.Symbol, t.TotalSupply, t.Owner,
		t.Mintable, t.Burnable, t.AddrForbiddable, t.TokenForbiddable, t.URL, t.Description, t.Identity)

	if err != nil {
		return err
	}

	if !t.TokenForbiddable && t.IsForbidden {
		return ErrTokenForbiddenNotSupported(t.GetSymbol())
	}

	if !t.Mintable && t.TotalMint > 0 {
		return ErrTokenMintNotSupported(t.GetSymbol())
	}

	if !t.Burnable && t.TotalBurn > 0 {
		return ErrTokenBurnNotSupported(t.GetSymbol())
	}

	return nil
}

func (t *BaseToken) SetName(name string) sdk.Error {
	if utf8.RuneCountInString(name) > 32 {
		return ErrInvalidTokenName(name)
	}

	t.Name = name
	return nil
}

func (t BaseToken) GetSymbol() string {
	return t.Symbol
}

func ValidateTokenSymbol(symbol string) sdk.Error {
	if !tokenSymbolRegex.MatchString(symbol) {
		return ErrInvalidTokenSymbol(symbol)
	}
	return nil
}

func IsSuffixSymbol(symbol string) bool {
	return suffixSymbolRegex.MatchString(symbol)
}

func (t *BaseToken) SetSymbol(symbol string) sdk.Error {
	if err := ValidateTokenSymbol(symbol); err != nil {
		return err
	}

	t.Symbol = symbol
	return nil
}

func (t BaseToken) GetTotalSupply() int64 {
	return t.TotalSupply
}

func (t *BaseToken) SetTotalSupply(amt int64) sdk.Error {
	if amt > MaxTokenAmount || amt <= 0 {
		return ErrInvalidTokenSupply(amt)
	}
	t.TotalSupply = amt
	return nil
}

func (t BaseToken) GetOwner() sdk.AccAddress {
	return t.Owner
}

func (t *BaseToken) SetOwner(addr sdk.AccAddress) sdk.Error {
	if addr.Empty() {
		return ErrNilTokenOwner()
	}

	t.Owner = addr
	return nil
}

func (t BaseToken) GetMintable() bool {
	return t.Mintable
}

func (t *BaseToken) SetMintable(enable bool) {
	t.Mintable = enable
}

func (t BaseToken) GetBurnable() bool {
	return t.Burnable
}

func (t *BaseToken) SetBurnable(enable bool) {
	t.Burnable = enable
}

func (t BaseToken) GetAddrForbiddable() bool {
	return t.AddrForbiddable
}

func (t *BaseToken) SetAddrForbiddable(enable bool) {
	t.AddrForbiddable = enable
}

func (t BaseToken) GetTokenForbiddable() bool {
	return t.TokenForbiddable
}

func (t *BaseToken) SetTokenForbiddable(enable bool) {
	t.TokenForbiddable = enable
}

func (t BaseToken) GetURL() string {
	return t.URL
}

func (t *BaseToken) SetURL(url string) sdk.Error {
	if utf8.RuneCountInString(url) > MaxTokenURLLength {
		return ErrInvalidTokenURL(url)
	}
	t.URL = url
	return nil
}

func (t BaseToken) GetDescription() string {
	return t.Description
}

func (t *BaseToken) SetDescription(description string) sdk.Error {
	if len(description) > MaxTokenDescriptionLength {
		return ErrInvalidTokenDescription(description)
	}
	t.Description = description
	return nil
}

func (t BaseToken) GetIdentity() string {
	return t.Identity
}

func (t *BaseToken) SetIdentity(identity string) sdk.Error {
	if len(identity) > MaxTokenIdentityLength {
		return ErrInvalidTokenIdentity(identity)
	}
	t.Identity = identity
	return nil
}

func (t BaseToken) GetTotalBurn() int64 {
	return t.TotalBurn
}

func (t *BaseToken) SetTotalBurn(amt int64) sdk.Error {
	if amt > MaxTokenAmount || amt < 0 {
		return ErrInvalidTokenBurnAmt(amt)
	}
	t.TotalBurn = amt
	return nil
}

func (t BaseToken) GetTotalMint() int64 {
	return t.TotalMint
}

func (t *BaseToken) SetTotalMint(amt int64) sdk.Error {
	if amt > MaxTokenAmount || amt < 0 {
		return ErrInvalidTokenMintAmt(amt)
	}
	t.TotalMint = amt
	return nil
}

func (t BaseToken) GetIsForbidden() bool {
	return t.IsForbidden
}

func (t *BaseToken) SetIsForbidden(enable bool) {
	t.IsForbidden = enable
}

func (t BaseToken) String() string {
	return fmt.Sprintf(`Token Info: 
[
  Name:             %s
  Symbol:           %s
  TotalSupply:      %d
  Owner:            %s
  Mintable:         %t
  Burnable:         %t
  AddrForbiddable:  %t
  TokenForbiddable: %t
  TotalBurn:        %d
  TotalMint:        %d
  IsForbidden:      %t
  URL:              %s
  Description:      %s
  Identity:			%s
]`,
		t.Name, t.Symbol, t.TotalSupply, t.Owner.String(), t.Mintable, t.Burnable,
		t.AddrForbiddable, t.TokenForbiddable, t.TotalBurn, t.TotalMint, t.IsForbidden,
		t.URL, t.Description, t.Identity,
	)
}

func MustMarshalToken(cdc *codec.Codec, token Token) []byte {
	return cdc.MustMarshalBinaryBare(token)
}

func MustUnmarshalToken(cdc *codec.Codec, value []byte) Token {
	validator, err := UnmarshalToken(cdc, value)
	if err != nil {
		panic(err)
	}
	return validator
}

func UnmarshalToken(cdc *codec.Codec, value []byte) (token Token, err error) {
	err = cdc.UnmarshalBinaryBare(value, &token)
	return token, err
}

func NewTokenCoins(denom string, amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewInt64Coin(denom, amount))
}
