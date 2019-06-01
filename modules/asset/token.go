package asset

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Token is an interface used to store asset at a given token within state.
// Many complex conditions can be used in the concrete struct which implements Token.
type Token interface {
	GetName() string
	SetName(string) error

	GetSymbol() string
	SetSymbol(string) error

	GetTotalSupply() int64
	SetTotalSupply(int64) error

	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress) error

	GetMintable() bool
	SetMintable(bool)

	GetBurnable() bool
	SetBurnable(bool)

	GetAddrForbiddable() bool
	SetAddrForbiddable(bool)

	GetTokenForbiddable() bool
	SetTokenForbiddable(bool)

	GetTotalBurn() int64
	SetTotalBurn(int64) error

	GetTotalMint() int64
	SetTotalMint(int64) error

	GetIsForbidden() bool
	SetIsForbidden(bool)

	IsValid() error
	// Ensure that account implements stringer
	String() string
}

//-----------------------------------------------------------------------------
// BaseAccount

var _ Token = (*BaseToken)(nil)

// BaseToken - a base Token structure.
type BaseToken struct {
	Name        string         `json:"name"`         //  Name of the newly issued asset, limited to 32 unicode characters.
	Symbol      string         `json:"symbol"`       //  token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply int64          `json:"total_supply"` //  The total supply for this token [0]
	Owner       sdk.AccAddress `json:"owner"`        // The initial issuer of this token

	Mintable         bool `json:"mintable"`          // Whether this token could be minted after the issuing
	Burnable         bool `json:"burnable"`          // Whether this token could be burned
	AddrForbiddable  bool `json:"addr_forbiddable"`  // whether could forbid some addresses to forbid transaction
	TokenForbiddable bool `json:"token_forbiddable"` // whether token could be global forbid

	TotalBurn   int64 `json:"total_burn"`   // Total amount of burn
	TotalMint   int64 `json:"total_mint"`   // Total amount of mint
	IsForbidden bool  `json:"is_forbidden"` // Whether token being forbidden currently
}

var (
	// TokenSymbolRegex : Token symbol can be 2 ~ 8 characters long.
	TokenSymbolRegex = regexp.MustCompile("^[a-z][a-z0-9]{1,7}$")
)

// NewToken - new base token
func NewToken(name string, symbol string, totalSupply int64, owner sdk.AccAddress,
	mintable bool, burnable bool, addrForbiddable bool, tokenForbiddable bool) (*BaseToken, sdk.Error) {

	t := &BaseToken{}
	if err := t.SetName(name); err != nil {
		return nil, ErrorInvalidTokenName(err.Error())
	}
	if err := t.SetSymbol(symbol); err != nil {
		return nil, ErrorInvalidTokenSymbol(err.Error())
	}
	if err := t.SetOwner(owner); err != nil {
		return nil, ErrorInvalidTokenOwner(err.Error())
	}
	if err := t.SetTotalSupply(totalSupply); err != nil {
		return nil, ErrorInvalidTokenSupply(err.Error())
	}

	t.SetMintable(mintable)
	t.SetBurnable(burnable)
	t.SetAddrForbiddable(addrForbiddable)
	t.SetTokenForbiddable(tokenForbiddable)

	if err := t.SetTotalMint(0); err != nil {
		return nil, ErrorInvalidTokenMint(err.Error())
	}
	if err := t.SetTotalBurn(0); err != nil {
		return nil, ErrorInvalidTokenBurn(err.Error())
	}
	t.SetIsForbidden(false)

	return t, nil
}

func (t *BaseToken) IsValid() error {
	_, err := NewToken(t.Name, t.Symbol, t.TotalSupply, t.Owner,
		t.Mintable, t.Burnable, t.AddrForbiddable, t.TokenForbiddable)

	if err != nil {
		return err
	}

	if !t.TokenForbiddable && t.IsForbidden {
		return ErrorInvalidForbiddenState("Invalid Forbidden state")
	}

	if t.TotalMint < 0 {
		return ErrorInvalidTokenMint(fmt.Sprintf("Invalid total mint: %d", t.TotalMint))
	}

	if t.TotalBurn < 0 {
		return ErrorInvalidTokenMint(fmt.Sprintf("Invalid total burn: %d", t.TotalBurn))
	}

	return nil
}

func (t BaseToken) GetName() string {
	return t.Name
}

func (t *BaseToken) SetName(name string) error {
	if utf8.RuneCountInString(name) > 32 {
		return errors.New("token name limited to 32 unicode characters")
	}
	t.Name = name

	return nil
}

func (t BaseToken) GetSymbol() string {
	return t.Symbol
}

func (t *BaseToken) SetSymbol(symbol string) error {
	if !TokenSymbolRegex.MatchString(symbol) {
		return errors.New("token symbol limited to [a-z][a-z0-9]{1,7}")
	}
	t.Symbol = symbol

	return nil
}

func (t BaseToken) GetTotalSupply() int64 {
	return t.TotalSupply
}

func (t *BaseToken) SetTotalSupply(amt int64) error {
	if amt > MaxTokenAmount {
		return errors.New("token total supply before 1e8 boosting should be less than 90 billion")
	}
	if amt <= 0 {
		return errors.New("token total supply must a positive")
	}
	t.TotalSupply = amt
	return nil
}

func (t BaseToken) GetOwner() sdk.AccAddress {
	return t.Owner
}

func (t *BaseToken) SetOwner(addr sdk.AccAddress) error {
	if addr.Empty() {
		return errors.New("token owner is invalid")
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

func (t BaseToken) GetTotalBurn() int64 {
	return t.TotalBurn
}

func (t *BaseToken) SetTotalBurn(amt int64) error {
	if amt > MaxTokenAmount || amt < 0 {
		return errors.New("token total burn amt is invalid")
	}
	t.TotalBurn = amt
	return nil
}

func (t BaseToken) GetTotalMint() int64 {
	return t.TotalMint
}

func (t *BaseToken) SetTotalMint(amt int64) error {
	if amt > MaxTokenAmount || amt < 0 {
		return errors.New("token total mint amt is invalid")
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
	return fmt.Sprintf(`Token Info: [
  Name:           %s
  Symbol:         %s
  TotalSupply:    %d
  Owner:          %s
  Mintable:       %t
  Burnable:       %t 
  AddrForbiddable:  %t
  TokenForbiddable: %t
  TotalBurn:      %d
  TotalMint:      %d
  IsForbidden:       %t ]`,
		t.Name, t.Symbol, t.TotalSupply, t.Owner.String(), t.Mintable, t.Burnable,
		t.AddrForbiddable, t.TokenForbiddable, t.TotalBurn, t.TotalMint, t.IsForbidden,
	)
}

func NewTokenCoin(denom string, amount int64) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(amount))
}

func NewTokenCoins(denom string, amount int64) sdk.Coins {
	return sdk.NewCoins(NewTokenCoin(denom, amount))
}
