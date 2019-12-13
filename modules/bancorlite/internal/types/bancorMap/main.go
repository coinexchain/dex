package main

import (
"fmt"
"github.com/shopspring/decimal"
"math"
"os"
"sort"
)

func main() {
	buildMap()
}

/*
# y = x^(A)，x = [0, 1], A = (0, 5]
# map[int64]map[int64]string: A(0, 5] map to (0, 50], x[0, 1] map to (0, 1000); y -> string
*/
func buildMap() {
	decimal.DivisionPrecision = 16
	m := make(map[int]map[int]decimal.Decimal, 10)

	for i := 1; i <= 50; i++ {
		sub := make(map[int]decimal.Decimal)
		for j := 0; j <= 1000; j++ {

			a := 0.1 * float64(i)
			x := 0.001 * float64(j)
			sub[j] = decimal.NewFromFloat(math.Pow(x, a))
		}
		m[i] = sub
	}
	f, err := os.Create("./modules/bancorlite/internal/types/bancorMap/table.go")
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	f.WriteString("var BancorMap = map[int64]map[int64]string")
	_, _ = f.WriteString("{\n")
	var outerKeys []int
	for key := range m {
		outerKeys = append(outerKeys, key)
	}
	sort.Ints(outerKeys)
	for _, k := range outerKeys {
		s := fmt.Sprintf("%d:{\n", k)
		f.WriteString(s)
		var keys []int
		for key := range m[k] {
			keys = append(keys, key)
		}
		sort.Ints(keys)
		for _, innerK := range keys {
			s := fmt.Sprintf("%d:\"%s\",", innerK, m[k][innerK].String())
			f.WriteString(s)
		}
		f.WriteString("},\n")
	}
	f.WriteString("}")
}
