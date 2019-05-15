package asset

import (
sdk "github.com/cosmos/cosmos-sdk/types"
)

// Token is an interface used to store asset at a given token within state.
// Many complex conditions can be used in the concrete struct which implements Token.
type Token interface {
	GetName() string
	SetName() string

	GetSymbol() string
	SetSymbol() string

	GetTotalSupply() int64
	SetTotalSupply() int64

	GetOwner() sdk.AccAddress
	SetOwner(sdk.AccAddress)

	GetMintable() bool
	SetMintable() bool

	GetBurnable() bool
	SetBurnable() bool

	GetAddrFreezeable() bool
	SetAddrFreezeable() bool

	GetTokenFreezeable() bool
	SetTokenFreezeable() bool

	GetTotalBurn() int64
	SetTotalBurn() int64

	GetTotalMint() int64
	SetTotalMint() int64

	GetisFrozen() bool
	SetisFrozen() bool

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
	isFrozen  bool  // Whether token being frozen currently
}

func (token BaseToken) GetName() string {
	panic("implement me")
}

func (token BaseToken) SetName() string {
	panic("implement me")
}

func (token BaseToken) GetSymbol() string {
	panic("implement me")
}

func (token BaseToken) SetSymbol() string {
	panic("implement me")
}

func (token BaseToken) GetTotalSupply() int64 {
	panic("implement me")
}

func (token BaseToken) SetTotalSupply() int64 {
	panic("implement me")
}

func (token BaseToken) GetOwner() sdk.AccAddress {
	panic("implement me")
}

func (token BaseToken) SetOwner(sdk.AccAddress) {
	panic("implement me")
}

func (token BaseToken) GetMintable() bool {
	panic("implement me")
}

func (token BaseToken) SetMintable() bool {
	panic("implement me")
}

func (token BaseToken) GetBurnable() bool {
	panic("implement me")
}

func (token BaseToken) SetBurnable() bool {
	panic("implement me")
}

func (token BaseToken) GetAddrFreezeable() bool {
	panic("implement me")
}

func (token BaseToken) SetAddrFreezeable() bool {
	panic("implement me")
}

func (token BaseToken) GetTokenFreezeable() bool {
	panic("implement me")
}

func (token BaseToken) SetTokenFreezeable() bool {
	panic("implement me")
}

func (token BaseToken) GetTotalBurn() int64 {
	panic("implement me")
}

func (token BaseToken) SetTotalBurn() int64 {
	panic("implement me")
}

func (token BaseToken) GetTotalMint() int64 {
	panic("implement me")
}

func (token BaseToken) SetTotalMint() int64 {
	panic("implement me")
}

func (token BaseToken) GetisFrozen() bool {
	panic("implement me")
}

func (token BaseToken) SetisFrozen() bool {
	panic("implement me")
}

func (token BaseToken) String() string {
	panic("implement me")
}

