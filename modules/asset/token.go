package asset

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"unicode/utf8"
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

	GetAddrFreezeable() bool
	SetAddrFreezeable(bool)

	GetTokenFreezeable() bool
	SetTokenFreezeable(bool)

	GetTotalBurn() int64
	SetTotalBurn(int64) error

	GetTotalMint() int64
	SetTotalMint(int64) error

	GetIsFrozen() bool
	SetIsFrozen(bool)

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

	Mintable        bool `json:"mintable"`         // Whether this token could be minted after the issuing
	Burnable        bool `json:"burnable"`         // Whether this token could be burned
	AddrFreezeable  bool `json:"addr_freezeable"`  // whether could freeze some addresses to forbid transaction
	TokenFreezeable bool `json:"token_freezeable"` // whether token could be global freeze

	TotalBurn int64 `json:"total_burn"` // Total amount of burn
	TotalMint int64 `json:"total_mint"` // Total amount of mint
	IsFrozen  bool  `json:"is_frozen"`  // Whether token being frozen currently
}

// NewToken - new base token
func NewToken(name string, symbol string, amt int64, owner sdk.AccAddress,
	mintable bool, burnable bool, addrfreezeable bool, tokenfreezeable bool) (*BaseToken, sdk.Error) {

	t := &BaseToken{}
	if err := t.SetName(name); err != nil {
		return nil, ErrorInvalidTokenName(CodeSpaceAsset, err.Error())
	}
	if err := t.SetSymbol(symbol); err != nil {
		return nil, ErrorInvalidTokenSymbol(CodeSpaceAsset, err.Error())
	}
	if err := t.SetOwner(owner); err != nil {
		return nil, ErrorInvalidTokenOwner(CodeSpaceAsset, err.Error())
	}
	if err := t.SetTotalSupply(amt); err != nil {
		return nil, ErrorInvalidTokenSupply(CodeSpaceAsset, err.Error())
	}

	t.SetMintable(mintable)
	t.SetBurnable(burnable)
	t.SetAddrFreezeable(addrfreezeable)
	t.SetTokenFreezeable(tokenfreezeable)

	t.SetTotalMint(0)
	t.SetTotalBurn(0)
	t.SetIsFrozen(false)

	return t, nil
}

func (t BaseToken) GetName() string {
	return t.Name
}

func (t *BaseToken) SetName(name string) error {
	if utf8.RuneCountInString(name) > 32 {
		return errors.New("issue token name limited to 32 unicode characters")
	}
	t.Name = name

	return nil
}

func (t BaseToken) GetSymbol() string {
	return t.Symbol
}

func (t *BaseToken) SetSymbol(symbol string) error {
	if m, _ := regexp.MatchString("^[a-z][a-z0-9]{1,7}$", symbol); !m {
		return errors.New("issue token symbol limited to [a-z][a-z0-9]{1,7}")
	}
	t.Symbol = symbol

	return nil
}

func (t BaseToken) GetTotalSupply() int64 {
	return t.TotalSupply
}

func (t *BaseToken) SetTotalSupply(amt int64) error {
	if amt > MaxTokenAmount {
		return errors.New("token total supply limited to 90 billion")
	}
	if amt < 0 {
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
		return errors.New("issue token must set a valid token owner")
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

func (t BaseToken) GetAddrFreezeable() bool {
	return t.AddrFreezeable
}

func (t *BaseToken) SetAddrFreezeable(enable bool) {
	t.AddrFreezeable = enable
}

func (t BaseToken) GetTokenFreezeable() bool {
	return t.TokenFreezeable
}

func (t *BaseToken) SetTokenFreezeable(enable bool) {
	t.TokenFreezeable = enable
}

func (t BaseToken) GetTotalBurn() int64 {
	return t.TotalBurn
}

func (t *BaseToken) SetTotalBurn(amt int64) error {
	if amt > MaxTokenAmount {
		return errors.New("token total supply limited to 90 billion")
	}
	t.TotalBurn = amt
	return nil
}

func (t BaseToken) GetTotalMint() int64 {
	return t.TotalMint
}

func (t *BaseToken) SetTotalMint(amt int64) error {
	if amt > MaxTokenAmount {
		return errors.New("token total supply limited to 90 billion")
	}
	t.TotalMint = amt
	return nil
}

func (t BaseToken) GetIsFrozen() bool {
	return t.IsFrozen
}

func (t *BaseToken) SetIsFrozen(enable bool) {
	t.IsFrozen = enable
}

func (t BaseToken) String() string {
	return fmt.Sprintf(`Token Info: [
  Name:            %s
  Symbol:          %s
  TotalSupply:     %d
  Owner:           %s
  Mintable:        %t
  Burnable:        %t 
  AddrFreezeable:  %t
  TokenFreezeable: %t
  TotalBurn:       %d
  TotalMint:       %d
  IsFrozen:        %t ]`,
		t.Name, t.Symbol, t.TotalSupply, t.Owner.String(), t.Mintable, t.Burnable,
		t.AddrFreezeable, t.TokenFreezeable, t.TotalBurn, t.TotalMint, t.IsFrozen,
	)
}

func NewTokenCoin(denom string, amount int64) sdk.Coin {
	return sdk.NewCoin(denom, sdk.NewInt(amount))
}

func NewTokenCoins(denom string, amount int64) sdk.Coins {
	return sdk.NewCoins(NewTokenCoin(denom, amount))
}
