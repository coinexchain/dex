package types

import (
	"bytes"
	"strconv"

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
	_ sdk.Msg = &MsgModifyTokenInfo{}
)

// MsgIssueToken
type MsgIssueToken struct {
	Name             string         `json:"name" yaml:"name"`                           // Name of the newly issued asset, limited to 32 unicode characters
	Symbol           string         `json:"symbol" yaml:"symbol"`                       // token symbol, [a-z][a-z0-9]{1,7}
	TotalSupply      sdk.Int        `json:"total_supply" yaml:"total_supply"`           // The total supply for this token [0]
	Owner            sdk.AccAddress `json:"owner" yaml:"owner"`                         // The initial issuer of this token [1]
	Mintable         bool           `json:"mintable" yaml:"mintable"`                   // Whether this token could be minted after the issuing
	Burnable         bool           `json:"burnable" yaml:"burnable"`                   // Whether this token could be burned
	AddrForbiddable  bool           `json:"addr_forbiddable" yaml:"addr_forbiddable"`   // whether could forbid some addresses to forbid transaction
	TokenForbiddable bool           `json:"token_forbiddable" yaml:"token_forbiddable"` // whether token could be global forbid
	URL              string         `json:"url" yaml:"url"`                             //URL of token website
	Description      string         `json:"description" yaml:"description"`             //Description of token info
	Identity         string         `json:"identity" yaml:"identity"`                   //Identity of token
}

// NewMsgIssueToken
func NewMsgIssueToken(name string, symbol string, amt sdk.Int, owner sdk.AccAddress,
	mintable bool, burnable bool, addrForbiddable bool, tokenForbiddable bool,
	url string, description string, identity string) MsgIssueToken {

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
		identity,
	}
}

func (msg *MsgIssueToken) SetAccAddress(addr sdk.AccAddress) {
	msg.Owner = addr
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
		msg.Mintable, msg.Burnable, msg.AddrForbiddable, msg.TokenForbiddable, msg.URL, msg.Description, msg.Identity)
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

func (msg *MsgTransferOwnership) SetAccAddress(addr sdk.AccAddress) {
	msg.OriginalOwner = addr
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
	Amount       sdk.Int        `json:"amount" yaml:"amount"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"`
}

func NewMsgMintToken(symbol string, amt sdk.Int, owner sdk.AccAddress) MsgMintToken {
	return MsgMintToken{
		symbol,
		amt,
		owner,
	}
}

func (msg *MsgMintToken) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddress = addr
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
	if !amt.IsPositive() {
		return ErrInvalidTokenMintAmt(amt.String())
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
	Amount       sdk.Int        `json:"amount" yaml:"amount"`
	OwnerAddress sdk.AccAddress `json:"owner_address" yaml:"owner_address"` //token owner address
}

func NewMsgBurnToken(symbol string, amt sdk.Int, owner sdk.AccAddress) MsgBurnToken {
	return MsgBurnToken{
		symbol,
		amt,
		owner,
	}
}

func (msg *MsgBurnToken) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddress = addr
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
	if !amt.IsPositive() {
		return ErrInvalidTokenBurnAmt(amt.String())
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

func (msg *MsgForbidToken) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddress = addr
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

func (msg *MsgUnForbidToken) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddress = addr
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

func (msg *MsgAddTokenWhitelist) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddress = addr
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

	for _, addr := range msg.Whitelist {
		if !addr.Empty() {
			return nil
		}
	}
	return ErrNilTokenWhitelist()
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

func (msg *MsgRemoveTokenWhitelist) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddress = addr
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
	for _, addr := range msg.Whitelist {
		if !addr.Empty() {
			return nil
		}
	}
	return ErrNilTokenWhitelist()
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

func (msg *MsgForbidAddr) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddr = addr
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
	for _, address := range msg.Addresses {
		if bytes.Equal(address, msg.OwnerAddr) {
			return ErrTokenOwnerSelfForbidden()
		}
	}

	for _, addr := range msg.Addresses {
		if !addr.Empty() {
			return nil
		}
	}
	return ErrNilForbiddenAddress()
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

func (msg *MsgUnForbidAddr) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddr = addr
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

	for _, addr := range msg.Addresses {
		if !addr.Empty() {
			return nil
		}
	}
	return ErrNilForbiddenAddress()
}

// GetSignBytes Implements Msg.
func (msg MsgUnForbidAddr) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgUnForbidAddr) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddr}
}

// MsgModifyTokenInfo
type MsgModifyTokenInfo struct {
	Symbol           string         `json:"symbol" yaml:"symbol"`
	OwnerAddress     sdk.AccAddress `json:"owner_address" yaml:"owner_address"`
	URL              string         `json:"url" yaml:"url"`
	Description      string         `json:"description" yaml:"description"`
	Identity         string         `json:"identity" yaml:"identity"`
	Name             string         `json:"name" yaml:"name"`
	TotalSupply      string         `json:"total_supply" yaml:"total_supply"`
	Mintable         string         `json:"mintable" yaml:"mintable"`
	Burnable         string         `json:"burnable" yaml:"burnable"`
	AddrForbiddable  string         `json:"addr_forbiddable" yaml:"addr_forbiddable"`
	TokenForbiddable string         `json:"token_forbiddable" yaml:"token_forbiddable"`
}

func NewMsgModifyTokenInfo(symbol, url, description, identity string, owner sdk.AccAddress,
	name, totalSupply, mintable, burnable, addrForbiddable, tokenForbiddable string) MsgModifyTokenInfo {
	return MsgModifyTokenInfo{
		Symbol:           symbol,
		URL:              url,
		Description:      description,
		Identity:         identity,
		OwnerAddress:     owner,
		Name:             name,
		TotalSupply:      totalSupply,
		Mintable:         mintable,
		Burnable:         burnable,
		AddrForbiddable:  addrForbiddable,
		TokenForbiddable: tokenForbiddable,
	}
}

func (msg *MsgModifyTokenInfo) SetAccAddress(addr sdk.AccAddress) {
	msg.OwnerAddress = addr
}

// Route Implements Msg.
func (msg MsgModifyTokenInfo) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgModifyTokenInfo) Type() string {
	return "modify_token_info"
}

// ValidateBasic Implements Msg.
func (msg MsgModifyTokenInfo) ValidateBasic() sdk.Error {
	tmpToken := BaseToken{}
	if err := tmpToken.SetSymbol(msg.Symbol); err != nil {
		return err
	}
	if err := tmpToken.SetOwner(msg.OwnerAddress); err != nil {
		return err
	}
	if msg.URL != DoNotModifyTokenInfo {
		if err := tmpToken.SetURL(msg.URL); err != nil {
			return err
		}
	}
	if msg.Description != DoNotModifyTokenInfo {
		if err := tmpToken.SetDescription(msg.Description); err != nil {
			return err
		}
	}
	if msg.Identity != DoNotModifyTokenInfo {
		if err := tmpToken.SetIdentity(msg.Identity); err != nil {
			return err
		}
	}
	if msg.Name != DoNotModifyTokenInfo {
		if err := tmpToken.SetName(msg.Name); err != nil {
			return err
		}
	}
	if msg.TotalSupply != DoNotModifyTokenInfo {
		if supply, ok := sdk.NewIntFromString(msg.TotalSupply); !ok {
			return ErrInvalidTokenInfo("TotalSupply", msg.TotalSupply)
		} else if err := tmpToken.SetTotalSupply(supply); err != nil {
			return err
		}
	}
	if err := validateBoolField("Mintable", msg.Mintable); err != nil {
		return err
	}
	if err := validateBoolField("Burnable", msg.Burnable); err != nil {
		return err
	}
	if err := validateBoolField("AddrForbiddable", msg.AddrForbiddable); err != nil {
		return err
	}
	if err := validateBoolField("TokenForbiddable", msg.TokenForbiddable); err != nil {
		return err
	}

	return nil
}

func validateBoolField(fieldName, valStr string) sdk.Error {
	if valStr != DoNotModifyTokenInfo {
		if _, err := strconv.ParseBool(valStr); err != nil {
			return ErrInvalidTokenInfo(fieldName, valStr)
		}
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgModifyTokenInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgModifyTokenInfo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}
