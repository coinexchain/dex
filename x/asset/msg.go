package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"unicode/utf8"
)

// RouterKey is they name of the asset module
const RouterKey = "asset"

// MsgIssueToken
type MsgIssueToken struct {
	Name        string         //  Name of the newly issued asset, limited to 32 unicode characters
	Symbol      string         //  token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply int64          //  The total supply for this token [0]
	Owner       sdk.AccAddress // The initial issuer of this token [1]

	Mintable        bool // Whether this token could be minted after the issuing
	Burnable        bool // Whether this token could be burned
	AddrFreezeable  bool // whether could freeze some addresses to forbid transaction
	TokenFreezeable bool // whether token could be global freeze
}

var _ sdk.Msg = MsgIssueToken{}

// NewMsgIssueToken
func NewMsgIssueToken(name string, symbol string, amt int64, owner sdk.AccAddress,
	mintable bool, burnable bool, addrFreezeable bool, tokenFreezeable bool) MsgIssueToken {

	return MsgIssueToken{
		name,
		symbol,
		amt,
		owner,
		mintable,
		burnable,
		addrFreezeable,
		tokenFreezeable,
	}
}

// Route Implements Msg.
func (msg MsgIssueToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgIssueToken) Type() string {
	return "issueToken"
}

// ValidateBasic Implements Msg.
func (msg MsgIssueToken) ValidateBasic() sdk.Error {

	if utf8.RuneCountInString(msg.Name) > 32 {
		return sdk.ErrTxDecode("issue token name limited to 32 unicode characters")
	}
	if m, _ := regexp.MatchString("[a-z][a-z0-9]{1,7}", msg.Symbol); !m {
		return sdk.ErrTxDecode("issue token symbol limited to [a-z][a-z0-9]{1,7}")
	}
	if msg.TotalSupply > MaxTokenAmount {
		return sdk.ErrTxDecode("issue token supply amt limited to 90 billion")
	}
	if msg.TotalSupply < 0 {
		return sdk.ErrTxDecode("issue token supply amt should be positive")
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("missing owner address")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTransferOwnership
type MsgTransferOwnership struct {
	Symbol        string
	OriginalOwner sdk.AccAddress
	NewOwner      sdk.AccAddress
}

var _ sdk.Msg = MsgTransferOwnership{}

// Route Implements Msg.
func (msg MsgTransferOwnership) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgTransferOwnership) Type() string {
	return "transferOwnerShip"
}

// ValidateBasic Implements Msg.
func (msg MsgTransferOwnership) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgTransferOwnership) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgTransferOwnership) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgFreezeAddress
type MsgFreezeAddress struct {
	Symbol  string
	address sdk.AccAddress
}

var _ sdk.Msg = MsgFreezeAddress{}

// Route Implements Msg.
func (msg MsgFreezeAddress) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgFreezeAddress) Type() string {
	return "freezeAddress"
}

// ValidateBasic Implements Msg.
func (msg MsgFreezeAddress) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgFreezeAddress) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgFreezeAddress) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgUnfreezeAddress
type MsgUnfreezeAddress struct {
	Symbol  string
	address sdk.AccAddress
}

var _ sdk.Msg = MsgUnfreezeAddress{}

// Route Implements Msg.
func (msg MsgUnfreezeAddress) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgUnfreezeAddress) Type() string {
	return "unfreezeAddress"
}

// ValidateBasic Implements Msg.
func (msg MsgUnfreezeAddress) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgUnfreezeAddress) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgUnfreezeAddress) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgFreezeToken
type MsgFreezeToken struct {
	Symbol  string
	address sdk.AccAddress // Whitelist
}

var _ sdk.Msg = MsgFreezeToken{}

// Route Implements Msg.
func (msg MsgFreezeToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgFreezeToken) Type() string {
	return "freezeToken"
}

// ValidateBasic Implements Msg.
func (msg MsgFreezeToken) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgFreezeToken) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgFreezeToken) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgUnfreezeToken
type MsgUnfreezeToken struct {
	Symbol  string
	address sdk.AccAddress // Whitelist
}

var _ sdk.Msg = MsgUnfreezeToken{}

// Route Implements Msg.
func (msg MsgUnfreezeToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgUnfreezeToken) Type() string {
	return "unfreezeToken"
}

// ValidateBasic Implements Msg.
func (msg MsgUnfreezeToken) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgUnfreezeToken) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgUnfreezeToken) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgBurnToken
type MsgBurnToken struct {
	Symbol       string
	Amount       uint64         //[0]
	ownerAddress sdk.AccAddress //token owner address
}

var _ sdk.Msg = MsgBurnToken{}

// Route Implements Msg.
func (msg MsgBurnToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgBurnToken) Type() string {
	return "burnToken"
}

// ValidateBasic Implements Msg.
func (msg MsgBurnToken) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgBurnToken) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgMintToken
type MsgMintToken struct {
	Symbol       string
	Amount       uint64 //[0]
	ownerAddress sdk.AccAddress
}

var _ sdk.Msg = MsgMintToken{}

// Route Implements Msg.
func (msg MsgMintToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgMintToken) Type() string {
	return "mintToken"
}

// ValidateBasic Implements Msg.
func (msg MsgMintToken) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgMintToken) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	panic("implement me")
}
