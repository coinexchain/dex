package denoms

import "strings"

const CET = "cet"

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
