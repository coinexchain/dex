package types

import (
	"strings"

	dex "github.com/coinexchain/dex/types"
)

func IsOnlyForCoinEx(alias string) bool {
	if strings.HasPrefix(alias, "coinex") ||
		strings.HasSuffix(alias, "coinex") ||
		strings.HasSuffix(alias, "coinex.org") ||
		strings.HasSuffix(alias, "coinex.com") {
		return true
	}

	return alias == dex.CET || alias == "viabtc" || alias == "cetdac"
}

func IsValidAlias(alias string) bool {
	if len(alias) < 2 || len(alias) > 100 {
		return false
	}
	for _, c := range alias {
		if !isValidChar(c) {
			return false
		}
	}
	return true
}

func isValidChar(c rune) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	if 'a' <= c && c <= 'z' {
		return true
	}
	if c == '-' || c == '_' || c == '.' || c == '@' {
		return true
	}
	return false
}
