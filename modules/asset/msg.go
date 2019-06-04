package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgIssueToken
type MsgIssueToken struct {
	Name             string         `json:"name"`              // Name of the newly issued asset, limited to 32 unicode characters
	Symbol           string         `json:"symbol"`            // token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply      int64          `json:"total_supply"`      // The total supply for this token [0]
	Owner            sdk.AccAddress `json:"owner"`             // The initial issuer of this token [1]
	Mintable         bool           `json:"mintable"`          // Whether this token could be minted after the issuing
	Burnable         bool           `json:"burnable"`          // Whether this token could be burned
	AddrForbiddable  bool           `json:"addr_forbiddable"`  // whether could forbid some addresses to forbid transaction
	TokenForbiddable bool           `json:"token_forbiddable"` // whether token could be global forbid
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
	return "issue_token"
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
	Symbol        string         `json:"symbol"`
	OriginalOwner sdk.AccAddress `json:"original_owner"`
	NewOwner      sdk.AccAddress `json:"new_owner"`
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
	return "transfer_ownerShip"
}

// ValidateBasic Implements Msg.
func (msg MsgTransferOwnership) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrorInvalidTokenSymbol(err.Error())
	}

	if msg.OriginalOwner.Empty() || msg.NewOwner.Empty() {
		return ErrorInvalidTokenOwner("transfer owner ship need a valid addr")
	}

	if msg.OriginalOwner.Equals(msg.NewOwner) {
		return ErrorInvalidTokenOwner("Can not and no need to transfer ownership to self")
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

// MsgMintToken
type MsgMintToken struct {
	Symbol       string         `json:"symbol"`
	Amount       int64          `json:"amount"`
	OwnerAddress sdk.AccAddress `json:"owner_address"`
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
	return "mint_token"
}

// ValidateBasic Implements Msg.
func (msg MsgMintToken) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrorInvalidTokenSymbol(err.Error())
	}

	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("mint token need a valid owner addr")
	}

	if msg.Amount > MaxTokenAmount {
		return ErrorInvalidTokenMint("token total supply before 1e8 boosting should be less than 90 billion")
	}

	if msg.Amount <= 0 {
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

// MsgBurnToken
type MsgBurnToken struct {
	Symbol       string         `json:"symbol"`
	Amount       int64          `json:"amount"`
	OwnerAddress sdk.AccAddress `json:"owner_address"` //token owner address
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
	return "burn_token"
}

// ValidateBasic Implements Msg.
func (msg MsgBurnToken) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrorInvalidTokenSymbol(err.Error())
	}

	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("burn token need a valid owner addr")
	}

	if msg.Amount > MaxTokenAmount {
		return ErrorInvalidTokenBurn("token total supply before 1e8 boosting should be less than 90 billion")
	}

	if msg.Amount <= 0 {
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

// MsgForbidToken
type MsgForbidToken struct {
	Symbol       string         `json:"symbol"`
	OwnerAddress sdk.AccAddress `json:"owner_address"`
}

var _ sdk.Msg = MsgForbidToken{}

func NewMsgForbidToken(symbol string, owner sdk.AccAddress) MsgForbidToken {
	return MsgForbidToken{
		symbol,
		owner,
	}
}

// Route Implements Msg.
func (msg MsgForbidToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgForbidToken) Type() string {
	return "forbid_token"
}

// ValidateBasic Implements Msg.
func (msg MsgForbidToken) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrorInvalidTokenSymbol(err.Error())
	}
	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("forbid token need a valid owner addr")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgForbidToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgForbidToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgUnForbidToken
type MsgUnForbidToken struct {
	Symbol       string         `json:"symbol"`
	OwnerAddress sdk.AccAddress `json:"owner_address"`
}

var _ sdk.Msg = MsgUnForbidToken{}

func NewMsgUnForbidToken(symbol string, owner sdk.AccAddress) MsgUnForbidToken {
	return MsgUnForbidToken{
		symbol,
		owner,
	}
}

// Route Implements Msg.
func (msg MsgUnForbidToken) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgUnForbidToken) Type() string {
	return "unforbid_token"
}

// ValidateBasic Implements Msg.
func (msg MsgUnForbidToken) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrorInvalidTokenSymbol(err.Error())
	}
	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("forbid token need a valid owner addr")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUnForbidToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgUnForbidToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgAddWhitelist
type MsgAddTokenWhitelist struct {
	Symbol       string           `json:"symbol"`
	OwnerAddress sdk.AccAddress   `json:"owner_address"`
	Whitelist    []sdk.AccAddress `json:"whitelist"`
}

var _ sdk.Msg = MsgAddTokenWhitelist{}

func NewMsgAddTokenWhitelist(symbol string, owner sdk.AccAddress, whitelist []sdk.AccAddress) MsgAddTokenWhitelist {
	return MsgAddTokenWhitelist{
		symbol,
		owner,
		whitelist,
	}
}

// Route Implements Msg.
func (msg MsgAddTokenWhitelist) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgAddTokenWhitelist) Type() string {
	return "add_token_whitelist"
}

// ValidateBasic Implements Msg.
func (msg MsgAddTokenWhitelist) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrorInvalidTokenSymbol(err.Error())
	}
	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("add token whitelist need a valid owner addr")
	}
	if len(msg.Whitelist) == 0 {
		return ErrorInvalidTokenWhitelist("add nil token whitelist")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddTokenWhitelist) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAddTokenWhitelist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgRemoveWhitelist
type MsgRemoveTokenWhitelist struct {
	Symbol       string           `json:"symbol"`
	OwnerAddress sdk.AccAddress   `json:"owner_address"`
	Whitelist    []sdk.AccAddress `json:"whitelist"`
}

var _ sdk.Msg = MsgRemoveTokenWhitelist{}

func NewMsgRemoveTokenWhitelist(symbol string, owner sdk.AccAddress, whitelist []sdk.AccAddress) MsgRemoveTokenWhitelist {
	return MsgRemoveTokenWhitelist{
		symbol,
		owner,
		whitelist,
	}
}

// Route Implements Msg.
func (msg MsgRemoveTokenWhitelist) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgRemoveTokenWhitelist) Type() string {
	return "remove_token_whitelist"
}

// ValidateBasic Implements Msg.
func (msg MsgRemoveTokenWhitelist) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrorInvalidTokenSymbol(err.Error())
	}
	if msg.OwnerAddress.Empty() {
		return ErrorInvalidTokenOwner("remove token whitelist need a valid owner addr")
	}
	if len(msg.Whitelist) == 0 {
		return ErrorInvalidTokenWhitelist("remove nil token whitelist")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRemoveTokenWhitelist) GetSignBytes() []byte {
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRemoveTokenWhitelist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
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
	return "forbid_address"
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

// MsgUnForbidAddress
type MsgUnForbidAddress struct {
	Symbol  string
	address sdk.AccAddress
}

var _ sdk.Msg = MsgUnForbidAddress{}

// Route Implements Msg.
func (msg MsgUnForbidAddress) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgUnForbidAddress) Type() string {
	return "unforbid_address"
}

// ValidateBasic Implements Msg.
func (msg MsgUnForbidAddress) ValidateBasic() sdk.Error {
	panic("implement me")
}

// GetSignBytes Implements Msg.
func (msg MsgUnForbidAddress) GetSignBytes() []byte {
	panic("implement me")
}

// GetSigners Implements Msg.
func (msg MsgUnForbidAddress) GetSigners() []sdk.AccAddress {
	panic("implement me")
}
