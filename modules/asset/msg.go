package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgIssueToken
type MsgIssueToken struct {
	Name        string         `json:"name"`         //  Name of the newly issued asset, limited to 32 unicode characters
	Symbol      string         `json:"symbol"`       //  token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply int64          `json:"total_supply"` //  The total supply for this token [0]
	Owner       sdk.AccAddress `json:"owner"`        // The initial issuer of this token [1]

	Mintable         bool `json:"mintable"`          // Whether this token could be minted after the issuing
	Burnable         bool `json:"burnable"`          // Whether this token could be burned
	AddrForbiddable  bool `json:"addr_forbiddable"`  // whether could forbid some addresses to forbid transaction
	TokenForbiddable bool `json:"token_forbiddable"` // whether token could be global forbid
}

var _ sdk.Msg = MsgIssueToken{}

// NewMsgIssueToken
func NewMsgIssueToken(name string, symbol string, amt int64, owner sdk.AccAddress,
	mintable bool, burnable bool, addrForbiddable bool, tokenForbiddable bool) MsgIssueToken {

	return MsgIssueToken{
		name,
		symbol,
		amt,
		owner,
		mintable,
		burnable,
		addrForbiddable,
		tokenForbiddable,
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
	_, err := NewToken(msg.Name, msg.Symbol, msg.TotalSupply, msg.Owner,
		msg.Mintable, msg.Burnable, msg.AddrForbiddable, msg.TokenForbiddable)
	return err
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

func NewMsgTransferOwnership(symbol string, originalOwner sdk.AccAddress, newOwner sdk.AccAddress) MsgTransferOwnership {
	return MsgTransferOwnership{
		symbol,
		originalOwner,
		newOwner,
	}
}

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
	if msg.OriginalOwner.Empty() || msg.NewOwner.Empty() {
		return ErrorInvalidTokenOwner("transfer owner ship need a valid addr")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransferOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgTransferOwnership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OriginalOwner}
}

// MsgForbidAddress
type MsgForbidAddress struct {
	Symbol  string
	address sdk.AccAddress
}

var _ sdk.Msg = MsgForbidAddress{}

// Route Implements Msg.
func (msg MsgForbidAddress) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgForbidAddress) Type() string {
	return "forbidAddress"
}

// ValidateBasic Implements Msg.
func (msg MsgForbidAddress) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgForbidAddress) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgForbidAddress) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgUnforbidAddress
type MsgUnforbidAddress struct {
	Symbol  string
	address sdk.AccAddress
}

var _ sdk.Msg = MsgUnforbidAddress{}

// Route Implements Msg.
func (msg MsgUnforbidAddress) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgUnforbidAddress) Type() string {
	return "unforbidAddress"
}

// ValidateBasic Implements Msg.
func (msg MsgUnforbidAddress) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgUnforbidAddress) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgUnforbidAddress) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgForbidToken
type MsgForbidToken struct {
	Symbol  string
	address sdk.AccAddress // Whitelist
}

var _ sdk.Msg = MsgForbidToken{}

// Route Implements Msg.
func (msg MsgForbidToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgForbidToken) Type() string {
	return "forbidToken"
}

// ValidateBasic Implements Msg.
func (msg MsgForbidToken) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgForbidToken) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgForbidToken) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgUnforbidToken
type MsgUnforbidToken struct {
	Symbol  string
	address sdk.AccAddress // Whitelist
}

var _ sdk.Msg = MsgUnforbidToken{}

// Route Implements Msg.
func (msg MsgUnforbidToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgUnforbidToken) Type() string {
	return "unforbidToken"
}

// ValidateBasic Implements Msg.
func (msg MsgUnforbidToken) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgUnforbidToken) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgUnforbidToken) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

// MsgBurnToken
type MsgBurnToken struct {
	Symbol       string
	Amount       int64
	OwnerAddress sdk.AccAddress //token owner address
}

var _ sdk.Msg = MsgBurnToken{}

func NewMsgBurnToken(symbol string, amt int64, owner sdk.AccAddress) MsgBurnToken {
	return MsgBurnToken{
		symbol,
		amt,
		owner,
	}
}

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
	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("burn token need a valid addr")
	}
	if msg.Amount > MaxTokenAmount {
		return ErrorInvalidTokenBurn("token total supply limited to 90 billion")
	}
	if msg.Amount < 0 {
		return ErrorInvalidTokenBurn("burn amount should be positive")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurnToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgMintToken
type MsgMintToken struct {
	Symbol       string
	Amount       int64
	OwnerAddress sdk.AccAddress
}

var _ sdk.Msg = MsgMintToken{}

func NewMsgMintToken(symbol string, amt int64, owner sdk.AccAddress) MsgMintToken {
	return MsgMintToken{
		symbol,
		amt,
		owner,
	}
}

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
	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("mint token need a valid addr")
	}
	if msg.Amount > MaxTokenAmount {
		return ErrorInvalidTokenMint("token total supply limited to 90 billion")
	}
	if msg.Amount < 0 {
		return ErrorInvalidTokenMint("mint amount should be positive")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}
