package main

import (
	"os"
	"fmt"

	"github.com/coinexchain/dex/codec"
)

func main() {
	usage := "usage: %s [info|codec|proto|ser]\n"
	if len(os.Args) != 2 {
		fmt.Printf(usage, os.Args[0])
		return
	}
	switch os.Args[1] {
	case "info":
		codec.ShowInfo()
	case "codec":
		codec.GenerateCodecFile(os.Stdout)
	case "proto":
		codec.GenerateProtoFile()
	case "ser":
		codec.GenerateSerializableImpl(os.Stdout)
	default:
		fmt.Printf(usage, os.Args[0])
	}
}

