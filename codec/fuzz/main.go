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

	buf := make([]byte, 0, 4096)
	for i := 0; i < Count; i++ {
		if i%10000 == 0 {
			fmt.Printf("=== %d ===\n", i)
		}

		ifc := dexcodec.RandAny(r)
		origS, _ := json.Marshal(ifc)
		buf = buf[:0]
		dexcodec.EncodeAny(&buf, ifc)
		ifcDec, _, err := dexcodec.DecodeAny(buf)
		if err != nil {
			fmt.Printf("Now: %d\n", i)
			codon.ShowInfoForVar(leafTypes, ifc)
			panic(err)
		}
		cpDec := dexcodec.DeepCopyAny(ifcDec)
		decS, _ := json.Marshal(cpDec)
		if !bytes.Equal(origS, decS) {
			fmt.Printf("Now: %d\n%s\n%s\n", i, string(origS), string(decS))
			codon.ShowInfoForVar(leafTypes, ifc)
			panic("Mismatch!")
		}
	}
}
