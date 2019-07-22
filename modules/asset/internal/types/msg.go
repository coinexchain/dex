package types

import (
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgIssueToken{}
	_ sdk.Msg = &MsgTransferOwnership{}
	_ sdk.Msg = &MsgMintToken{}
	_ sdk.Msg = &MsgBurnToken{}
	_ sdk.Msg = &MsgForbidToken{}
	_ sdk.Msg = &MsgUnForbidToken{}
	_ sdk.Msg = &MsgAddTokenWhitelist{}
	_ sdk.Msg = &MsgRemoveTokenWhitelist{}
	_ sdk.Msg = &MsgForbidAddr{}
	_ sdk.Msg = &MsgUnForbidAddr{}
	_ sdk.Msg = &MsgModifyTokenDescription{}
	_ sdk.Msg = &MsgModifyTokenURL{}
)

// MsgIssueToken
type MsgIssueToken struct {
	Name             string         `json:"name" yaml:"name"`              // Name of the newly issued asset, limited to 32 unicode characters
	Symbol           string         `json:"symbol" yaml:"symbol"`            // token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply      int64          `json:"total_supply" yaml:"total_supply"`      // The total supply for this token [0]
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`             // The initial issuer of this token [1]
	Mintable         bool           `json:"mintable" yaml:"mintable"`          // Whether this token could be minted after the issuing
	Burnable         bool           `json:"burnable" yaml:"burnable"`          // Whether this token could be burned
	AddrForbiddable  bool           `json:"addr_forbiddable" yaml:"addr_forbiddable"`  // whether could forbid some addresses to forbid transaction
	TokenForbiddable bool           `json:"token_forbiddable" yaml:"token_forbiddable"` // whether token could be global forbid
	URL              string         `json:"url" yaml:"url"`               //URL of token website
	Description      string         `json:"description" yaml:"description"`       //Description of token info
}

// NewMsgIssueToken
func NewMsgIssueToken(name string, symbol string, amt int64, owner sdk.AccAddress,
	mintable bool, burnable bool, addrForbiddable bool, tokenForbiddable bool, url string, description string) MsgIssueToken {

	return MsgIssueToken{
		name,
		symbol,
		amt,
		owner,
		mintable,
		burnable,
		addrForbiddable,
		tokenForbiddable,
		url,
		description,
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
		msg.Mintable, msg.Burnable, msg.AddrForbiddable, msg.TokenForbiddable, msg.URL, msg.Description)
	return err
}

// GetSignBytes Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgTransferOwnership
type MsgTransferOwnership struct {
	Symbol        string         `json:"symbol" yaml:"symbol"`
	OriginalOwner sdk.AccAddress `json:"original_owner" yaml:"original_owner"`
	NewOwner      sdk.AccAddress `json:"new_owner" yaml:"new_owner"`
}

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
	return "transfer_ownership"
}

// ValidateBasic Implements Msg.
func (msg MsgTransferOwnership) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return err
	}

	if msg.OriginalOwner.Empty() || msg.NewOwner.Empty() {
		return ErrNilTokenOwner()
	}

	if msg.OriginalOwner.Equals(msg.NewOwner) {
		return ErrTransferSelfTokenOwner()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransferOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgTransferOwnership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OriginalOwner}
}

// MsgMintToken
type MsgMintToken struct {
	Symbol       string         `json:"symbol" yaml:"symbol"`
	Amount       int64          `json:"amount" yaml:"amount"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"`
}

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
		return err
	}

	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}

	amt := msg.Amount
	if amt > MaxTokenAmount || amt <= 0 {
		return ErrInvalidTokenMintAmt(amt)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgBurnToken
type MsgBurnToken struct {
	Symbol       string         `json:"symbol" yaml:"symbol"`
	Amount       int64          `json:"amount" yaml:"amount"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"` //token owner address
}

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
		return err
	}

	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}

	amt := msg.Amount
	if amt > MaxTokenAmount || amt <= 0 {
		return ErrInvalidTokenBurnAmt(amt)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurnToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgForbidToken
type MsgForbidToken struct {
	Symbol       string         `json:"symbol" yaml:"symbol"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"`
}

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
		return err
	}
	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgForbidToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgForbidToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgUnForbidToken
type MsgUnForbidToken struct {
	Symbol       string         `json:"symbol" yaml:"symbol"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"`
}

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
		return err
	}
	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUnForbidToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgUnForbidToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgAddWhitelist
type MsgAddTokenWhitelist struct {
	Symbol       string           `json:"symbol" yaml:"symbol"`
	OwnerAddress sdk.AccAddress   `json:"owner_address" yaml:"owner_address"`
	Whitelist    []sdk.AccAddress `json:"whitelist" yaml:"whitelist"`
}

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
		return err
	}
	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}
	if len(msg.Whitelist) == 0 {
		return ErrNilTokenWhitelist()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddTokenWhitelist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAddTokenWhitelist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgRemoveWhitelist
type MsgRemoveTokenWhitelist struct {
	Symbol       string           `json:"symbol" yaml:"symbol"`
	OwnerAddress sdk.AccAddress   `json:"owner_address" yaml:"owner_address"`
	Whitelist    []sdk.AccAddress `json:"whitelist" yaml:"whitelist"`
}

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
		return err
	}
	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}
	if len(msg.Whitelist) == 0 {
		return ErrNilTokenWhitelist()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRemoveTokenWhitelist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRemoveTokenWhitelist) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgForbidAddr
type MsgForbidAddr struct {
	Symbol    string           `json:"symbol" yaml:"symbol"`
	OwnerAddr sdk.AccAddress   `json:"owner_address" yaml:"owner_address"`
	Addresses []sdk.AccAddress `json:"addresses" yaml:"addresses"`
}

func NewMsgForbidAddr(symbol string, owner sdk.AccAddress, addresses []sdk.AccAddress) MsgForbidAddr {
	return MsgForbidAddr{
		symbol,
		owner,
		addresses,
	}
}

// Route Implements Msg.
func (msg MsgForbidAddr) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgForbidAddr) Type() string {
	return "forbid_addr"
}

// ValidateBasic Implements Msg.
func (msg MsgForbidAddr) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return err
	}
	if msg.OwnerAddr.Empty() {
		return ErrNilTokenOwner()
	}
	if len(msg.Addresses) == 0 {
		return ErrNilForbiddenAddress()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgForbidAddr) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgForbidAddr) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddr}
}

// MsgUnForbidAddr
type MsgUnForbidAddr struct {
	Symbol    string           `json:"symbol" yaml:"symbol"`
	OwnerAddr sdk.AccAddress   `json:"owner_address" yaml:"owner_address"`
	Addresses []sdk.AccAddress `json:"addresses" yaml:"addresses"`
}

func NewMsgUnForbidAddr(symbol string, owner sdk.AccAddress, addresses []sdk.AccAddress) MsgUnForbidAddr {
	return MsgUnForbidAddr{
		symbol,
		owner,
		addresses,
	}
}

// Route Implements Msg.
func (msg MsgUnForbidAddr) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgUnForbidAddr) Type() string {
	return "unforbid_addr"
}

// ValidateBasic Implements Msg.
func (msg MsgUnForbidAddr) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return err
	}
	if msg.OwnerAddr.Empty() {
		return ErrNilTokenOwner()
	}
	if len(msg.Addresses) == 0 {
		return ErrNilForbiddenAddress()
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUnForbidAddr) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgUnForbidAddr) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddr}
}

// MsgModifyURL
type MsgModifyTokenURL struct {
	Symbol       string         `json:"symbol" yaml:"symbol"`
	URL          string         `json:"url" yaml:"url"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"` //token owner address
}

func NewMsgModifyTokenURL(symbol string, url string, owner sdk.AccAddress) MsgModifyTokenURL {
	return MsgModifyTokenURL{
		symbol,
		url,
		owner,
	}
}

// Route Implements Msg.
func (msg MsgModifyTokenURL) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgModifyTokenURL) Type() string {
	return "modify_token_url"
}

// ValidateBasic Implements Msg.
func (msg MsgModifyTokenURL) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return err
	}

	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}

	if utf8.RuneCountInString(msg.URL) > 100 {
		return ErrInvalidTokenURL(msg.URL)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgModifyTokenURL) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgModifyTokenURL) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

// MsgModifyTokenDescription
type MsgModifyTokenDescription struct {
	Symbol       string         `json:"symbol" yaml:"symbol"`
	Description  string         `json:"description" yaml:"description"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"` //token owner address
}

func NewMsgModifyTokenDescription(symbol string, description string, owner sdk.AccAddress) MsgModifyTokenDescription {
	return MsgModifyTokenDescription{
		symbol,
		description,
		owner,
	}
}

// Route Implements Msg.
func (msg MsgModifyTokenDescription) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgModifyTokenDescription) Type() string {
	return "modify_token_description"
}

// ValidateBasic Implements Msg.
func (msg MsgModifyTokenDescription) ValidateBasic() sdk.Error {
	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return err
	}

	if msg.OwnerAddress.Empty() {
		return ErrNilTokenOwner()
	}

	if len(msg.Description) > 1024 {
		return ErrInvalidTokenDescription(msg.Description)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgModifyTokenDescription) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgModifyTokenDescription) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}
