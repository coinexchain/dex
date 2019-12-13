package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"os"
)

func main() {
	buildTable()
}

/*
# y = x^(A)，x = [0, 1], A = [0, 5]
# [51][1001]string: A[0, 5] map to [0, 50], x[0, 1] map to [0, 1000]; y -> string
*/
func buildTable() {

	decimal.DivisionPrecision = 16
	m := [51][1001]decimal.Decimal{}

	f, err := os.Create("./modules/bancorlite/internal/types/CWTable/table.go")
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	f.WriteString("var BancorTable = [51][1001]string")
	f.WriteString("{\n")

	for i := 0; i <= 50; i++ {
		f.WriteString("{\n")
		for j := 0; j <= 1000; j++ {
			a := 0.1 * float64(i)
			x := 0.001 * float64(j)
			m[i][j] = decimal.NewFromFloat(math.Pow(x, a))
			s := fmt.Sprintf("\"%s\",", m[i][j].String())
			f.WriteString(s)
		}
		f.WriteString("},\n")
	}
	f.WriteString("}")
}
