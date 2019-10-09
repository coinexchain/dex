package types

import (
	"strings"
)

const (
	SymbolSeparator = "/"
)

func GetSymbol(stock, money string) string {
	return stock + SymbolSeparator + money
}

func SplitSymbol(symbol string) (stock, money string) {
	values := strings.Split(symbol, SymbolSeparator)
	stock, money = values[0], values[1]
	return
}
