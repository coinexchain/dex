package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "asset"

	// StoreKey is string representation of the store key for asset
	StoreKey = ModuleName

	// RouterKey is the message route for asset
	RouterKey = ModuleName

	// QuerierRoute is the querier route for asset
	QuerierRoute = ModuleName

	DefaultParamspace = ModuleName
)

var (
	SeparateKey      = []byte{0x3A}
	TokenKey         = []byte{0x01}
	WhitelistKey     = []byte{0x02}
	ForbiddenAddrKey = []byte{0x03}
)

// GetTokenStoreKey - TokenKey | symbol
func GetTokenStoreKey(symbol string) []byte {
	return append(TokenKey, symbol...)
}

// GetWhitelistStoreKey - WhitelistKey | Symbol | : | AccAddress
func GetWhitelistStoreKey(symbol string, addr sdk.AccAddress) []byte {
	return append(append(append(WhitelistKey, symbol...), SeparateKey...), addr...)
}

// GetWhitelistKeyPrefix - Prefix WhitelistKey | Symbol | :
func GetWhitelistKeyPrefix(symbol string) []byte {
	return append(append(WhitelistKey, symbol...), SeparateKey...)
}

// GetWhitelistKeyPrefixLength -  WhitelistKey length
func GetWhitelistKeyPrefixLength(symbol string) int {
	return len(GetWhitelistKeyPrefix(symbol))
}

// GetForbiddenAddrStoreKey - ForbiddenAddrKey | Symbol | : | AccAddress
func GetForbiddenAddrStoreKey(symbol string, addr sdk.AccAddress) []byte {
	return append(append(append(ForbiddenAddrKey, symbol...), SeparateKey...), addr...)
}

// GetForbiddenAddrKeyPrefix - ForbiddenAddrKey | Symbol | :
func GetForbiddenAddrKeyPrefix(symbol string) []byte {
	return append(append(ForbiddenAddrKey, symbol...), SeparateKey...)
}

// GetForbiddenAddrKeyPrefixLength - ForbiddenAddrKey length
func GetForbiddenAddrKeyPrefixLength(symbol string) int {
	return len(GetForbiddenAddrKeyPrefix(symbol))
}
