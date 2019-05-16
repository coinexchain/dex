package asset

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Token is an interface used to store asset at a given token within state.
// Many complex conditions can be used in the concrete struct which implements Token.
type Token interface {
	GetName() string
	SetName(string)

	GetSymbol() string
	SetSymbol(string)

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
	Name        string         //  Name of the newly issued asset, limited to 32 unicode characters.
	Symbol      string         //  token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply int64          //  The total supply for this token [0]
	Owner       sdk.AccAddress // The initial issuer of this token

	Mintable        bool // Whether this token could be minted after the issuing
	Burnable        bool // Whether this token could be burned
	AddrFreezeable  bool // whether could freeze some addresses to forbid transaction
	TokenFreezeable bool // whether token could be global freeze

	TotalBurn int64 // Total amount of burn
	TotalMint int64 // Total amount of mint
	IsFrozen  bool  // Whether token being frozen currently
}

// NewBaseToken - default Mintable/Burnable/AddrFreezeable/TokenFreezeable/IsFrozen
func NewToken() BaseToken {
	return BaseToken{
		Mintable:        false,
		Burnable:        false,
		AddrFreezeable:  false,
		TokenFreezeable: false,
	}
}

func (t BaseToken) GetName() string {
	return t.Name
}

func (t BaseToken) SetName(name string) {
	t.Name = name
}

func (t BaseToken) GetSymbol() string {
	return t.Symbol
}

func (t BaseToken) SetSymbol(symbol string) {
	t.Symbol = symbol
}

func (t BaseToken) GetTotalSupply() int64 {
	return t.TotalSupply
}

func (t BaseToken) SetTotalSupply(amt int64) error {
	if amt > MaxTokenAmount {
		return errors.New("token total supply limit to 90 billion")
	}
	t.TotalSupply = amt
	return nil
}

func (t BaseToken) GetOwner() sdk.AccAddress {
	return t.Owner
}

func (t BaseToken) SetOwner(addr sdk.AccAddress) error {
	if addr.Empty() {
		return errors.New("must set a valid token owner")
	}
	t.Owner = addr
	return nil
}

func (t BaseToken) GetMintable() bool {
	return t.Mintable
}

func (t BaseToken) SetMintable(enable bool) {
	t.Mintable = enable
}

func (t BaseToken) GetBurnable() bool {
	return t.Burnable
}

func (t BaseToken) SetBurnable(enable bool) {
	t.Burnable = enable
}

func (t BaseToken) GetAddrFreezeable() bool {
	return t.AddrFreezeable
}

func (t BaseToken) SetAddrFreezeable(enable bool) {
	t.AddrFreezeable = enable
}

func (t BaseToken) GetTokenFreezeable() bool {
	return t.TokenFreezeable
}

func (t BaseToken) SetTokenFreezeable(enable bool) {
	t.TokenFreezeable = enable
}

func (t BaseToken) GetTotalBurn() int64 {
	return t.TotalBurn
}

func (t BaseToken) SetTotalBurn(amt int64) error {
	if amt > MaxTokenAmount {
		return errors.New("token total supply limit to 90 billion")
	}
	t.TotalBurn = amt
	return nil
}

func (t BaseToken) GetTotalMint() int64 {
	return t.TotalMint
}

func (t BaseToken) SetTotalMint(amt int64) error {
	if amt > MaxTokenAmount {
		return errors.New("token total supply limit to 90 billion")
	}
	t.TotalMint = amt
	return nil
}

func (t BaseToken) GetIsFrozen() bool {
	return t.IsFrozen
}

func (t BaseToken) SetIsFrozen(enable bool) {
	t.IsFrozen = enable
}

func (t BaseToken) String() string {
	panic("implement me")
}
