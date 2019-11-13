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
	runBench(r)
}

func runBench(r dexcodec.RandSrc) {
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
		buf := make([]byte, 0, 1024)
		dexcodec.EncodeAny(&buf, accounts[i])
		if err != nil {
			panic(err)
		}
		bzList[i] = buf
	}
	for i := 0; i < len(accounts); i++ {
		obj, _, err := dexcodec.DecodeAny(bzList[i])
		v := obj.(dexcodec.AccountX)
		if err != nil {
			panic(err)
		}
		s, _ := json.Marshal(v)
		if !bytes.Equal(s, accountsJ[i]) {
			fmt.Printf("%s\n%s\n%d mismatch!\n", string(s), string(accountsJ[i]), i)
		}

	}
	println("========== Check 0 Finished ==========")
	extraCheck1(accounts, accountsJ)
	println("========== Check 1 Finished ==========")
	extraCheck2(accounts, accountsJ)
	println("========== Check 2 Finished ==========")

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
			bzList[i] = bzList[i][:0]
			dexcodec.EncodeAny(&bzList[i], accounts[i])
			totalBytes += len(bzList[i])
		}
		for i := 0; i < len(accounts); i++ {
			_, _, err = dexcodec.DecodeAny(bzList[i])
			if err != nil {
				panic(err)
			}
		}
	}
	span = time.Now().UnixNano() - nanoSecCount
	fmt.Printf("Codon: time = %d, bytes = %d, bytes/ns = %f\n", span, totalBytes, float64(totalBytes)/float64(span))
}

func extraCheck1(accounts []dexcodec.AccountX, accountsJ [][]byte) {
	var err error
	stub := dexcodec.CodonStub{}
	bzList := make([][]byte, 1000)
	for i := 0; i < len(accounts); i++ {
		bzList[i], err = stub.MarshalBinaryBare(accounts[i])
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < len(accounts); i++ {
		var v dexcodec.AccountX
		err := stub.UnmarshalBinaryBare(bzList[i], &v)
		if err != nil {
			panic(err)
		}
		s, _ := json.Marshal(v)
		if !bytes.Equal(s, accountsJ[i]) {
			fmt.Printf("%s\n%s\n%d mismatch!\n", string(s), string(accountsJ[i]), i)
		}

	}
}

func extraCheck2(accounts []dexcodec.AccountX, accountsJ [][]byte) {
	var err error
	stub := dexcodec.CodonStub{}
	bzList := make([][]byte, 1000)
	for i := 0; i < len(accounts); i++ {
		bzList[i], err = stub.MarshalBinaryLengthPrefixed(accounts[i])
		if err != nil {
			panic(err)
		}
	}
	for i := 0; i < len(accounts); i++ {
		var v dexcodec.AccountX
		err := stub.UnmarshalBinaryLengthPrefixed(bzList[i], &v)
		if err != nil {
			panic(err)
		}
		s, _ := json.Marshal(v)
		if !bytes.Equal(s, accountsJ[i]) {
			fmt.Printf("%s\n%s\n%d mismatch!\n", string(s), string(accountsJ[i]), i)
		}

	}
}
