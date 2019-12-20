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
# [51][1001]string: A[0, 5] map to [0, 500], x[0, 1] map to [0, 1000]; y -> string
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
const MaxAR = 5000
const ARSamples = 1000
const SupplyRatioSamples = 1000
func TableLookup(x,y int64) sdk.Dec {
	return sdk.NewDec(int64(BancorTable[int(x)][int(y)])).Quo(sdk.NewDec(int64(math.MaxInt32)))
}
`)
	f.WriteString("var BancorTable = [6001][1001]int32")
	f.WriteString("{\n")
	maxAr := 5000
	arSamples := 1000
	supplyRatioSamples := 1000
	for i := 0; i <= maxAr+maxAr/5; i++ {
		f.WriteString("{\n")
		for j := 0; j <= supplyRatioSamples; j++ {
			a := 1.0 / float64(arSamples) * float64(i)
			//fmt.Println(a)
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
