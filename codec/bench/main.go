package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	dexcodec "github.com/coinexchain/dex/codec"
	"github.com/coinexchain/randsrc"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s filename\n", os.Args[0])
		return
	}
	r := randsrc.NewRandSrcFromFile(os.Args[1])
	accounts := make([]dexcodec.AccountX, 1000)
	accountsJ := make([][]byte, 1000)
	for i := 0; i < len(accounts); i++ {
		accounts[i] = dexcodec.RandAccountX(r)
		s, _ := json.Marshal(accounts[i])
		accountsJ[i] = s
		//fmt.Printf("Here %s\n", s)
	}

	// Check correctness of codon
	var err error
	bzList := make([][]byte, 1000)
	for i := 0; i < len(accounts); i++ {
		var buf bytes.Buffer
		err = dexcodec.BareEncodeAny(&buf, accounts[i])
		if err != nil {
			panic(err)
		}
		bzList[i] = buf.Bytes()
	}
	for i := 0; i < len(accounts); i++ {
		var v dexcodec.AccountX
		_, err = dexcodec.BareDecodeAny(bzList[i], &v)
		if err != nil {
			panic(err)
		}
		s, _ := json.Marshal(v)
		if !bytes.Equal(s, accountsJ[i]) {
			fmt.Printf("%s\n%s\n%d mismatch!\n", string(s), string(accountsJ[i]), i)
		}

	}

	cdc := codec.New()
	totalBytes := 0
	nanoSecCount := time.Now().UnixNano()
	for j := 0; j < 300; j++ {
		for i := 0; i < len(accounts); i++ {
			bzList[i], err = cdc.MarshalBinaryBare(accounts[i])
			totalBytes += len(bzList[i])
			if err != nil {
				panic(err)
			}
		}
		for i := 0; i < len(accounts); i++ {
			err = cdc.UnmarshalBinaryBare(bzList[i], &accounts[i])
			if err != nil {
				panic(err)
			}
		}
	}
	span := time.Now().UnixNano() - nanoSecCount
	fmt.Printf("Amino: time = %d, bytes = %d, bytes/ns = %f\n", span, totalBytes, float64(totalBytes)/float64(span))

	totalBytes = 0
	nanoSecCount = time.Now().UnixNano()
	for j := 0; j < 300; j++ {
		for i := 0; i < len(accounts); i++ {
			var buf bytes.Buffer
			err = dexcodec.BareEncodeAny(&buf, accounts[i])
			if err != nil {
				panic(err)
			}
			bzList[i] = buf.Bytes()
			totalBytes += len(bzList[i])
		}
		for i := 0; i < len(accounts); i++ {
			_, err = dexcodec.BareDecodeAny(bzList[i], &accounts[i])
			if err != nil {
				panic(err)
			}
		}
	}
	span = time.Now().UnixNano() - nanoSecCount
	fmt.Printf("Codon: time = %d, bytes = %d, bytes/ns = %f\n", span, totalBytes, float64(totalBytes)/float64(span))
}
