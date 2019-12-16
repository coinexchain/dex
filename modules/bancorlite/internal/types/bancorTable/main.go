package main

import (
	"fmt"
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
	f, err := os.Create("./modules/bancorlite/internal/types/table.go")
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	f.WriteString(`
package types
import (
	"math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)
const MaxAR = 50
const SupplyRatioSamples = 1000
func TableLookup(x,y int64) sdk.Dec {
	return sdk.NewDec(int64(BancorTable[int(x)][int(y)])).Quo(sdk.NewDec(int64(math.MaxInt32)))
}
`)
	f.WriteString("var BancorTable = [61][1001]int32")
	f.WriteString("{\n")

	for i := 0; i <= 60; i++ {
		f.WriteString("{\n")
		for j := 0; j <= 1000; j++ {
			a := 0.1 * float64(i)
			x := 0.001 * float64(j)
			v := int32(math.Pow(x, a) * float64(math.MaxInt32))
			s := fmt.Sprintf("%d,", v)
			if j == 1000 {
				s = fmt.Sprintf("%d},\n", v)
			} else if j%10 == 9 {
				s = fmt.Sprintf("%d,\n", v)
			}
			f.WriteString(s)
		}
	}
	f.WriteString("}")
}
