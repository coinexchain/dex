package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	amino "github.com/tendermint/go-amino"

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
	runRandTest(r)
}

func runRandTest(r dexcodec.RandSrc) {
	codec.RunInitFuncList()
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

func registerAll(cdcImp *dexcodec.CodecImp, cdcAmino *amino.Codec) {
	for _, entry := range dexcodec.TypeEntryList {
		v := entry.Value
		name := entry.Alias
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Interface {
			cdcImp.RegisterInterface(v, nil)
			cdcAmino.RegisterInterface(v, nil)
		} else {
			cdcImp.RegisterConcrete(v, name, nil)
			cdcAmino.RegisterConcrete(v, name, nil)
		}
	}
}

func findMismatch(a, b []byte) int {
	length := len(a)
	if len(b) < len(a) {
		length = len(b)
	}
	for i := 0; i < length; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	if len(b) != len(a) {
		return length
	}
	return -1
}
