package types

import "strings"

const CET = "cet"

// default bond denomination
const DefaultBondDenom = CET

// reserved token names
var reserved = []string{
	CET,
	"btc", "eth", "eos",
}

func IsReserved(denom string) bool {
	for _, d := range reserved {
		if strings.ToLower(denom) == d {
			return true
		}
	}
	return false
}
