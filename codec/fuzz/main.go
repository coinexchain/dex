package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/coinexchain/codon"
	dexcodec "github.com/coinexchain/dex/codec"
	"github.com/coinexchain/randsrc"
)

var Count = 100 * 10000

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s filename\n", os.Args[0])
		return
	}
	r := randsrc.NewRandSrcFromFile(os.Args[1])

	leafTypes := dexcodec.GetLeafTypes()

	for i := 0; i < Count; i++ {
		if i%10000 == 0 {
			fmt.Printf("=== %d ===\n", i)
		}

		ifc := dexcodec.RandAny(r)
		origS, _ := json.Marshal(ifc)
		var buf bytes.Buffer
		err := dexcodec.EncodeAny(&buf, ifc)
		if err != nil {
			fmt.Printf("Now: %d\n", i)
			codon.ShowInfoForVar(leafTypes, ifc)
			panic(err)
		}
		ifcDec, _, err := dexcodec.DecodeAny(buf.Bytes())
		if err != nil {
			fmt.Printf("Now: %d\n", i)
			codon.ShowInfoForVar(leafTypes, ifc)
			panic(err)
		}
		decS, _ := json.Marshal(ifcDec)
		if !bytes.Equal(origS, decS) {
			fmt.Printf("Now: %d\n%s\n%s\n", i, string(origS), string(decS))
			codon.ShowInfoForVar(leafTypes, ifc)
			panic("Mismatch!")
		}
	}
}
