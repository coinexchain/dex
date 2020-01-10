package main

import (
	"os"
	"fmt"

	"github.com/coinexchain/dex/codec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s [info|codec|proto]\n", os.Args[0])
		return
	}
	switch os.Args[1] {
	case "info":
		codec.ShowInfo()
	case "codec":
		genCodec()
	case "proto":
		codec.GenerateProtoFile()
	default:
		fmt.Printf("usage: %s [info|codec|proto]\n", os.Args[0])
	}
}

func genCodec() {
	codec.GenerateCodecFile(os.Stdout)
}
