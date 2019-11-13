package main

import (
	"os"

	"github.com/coinexchain/dex/codec"
)

func main() {
	//codec.ShowInfo()
	genCode()
}

func genCode() {
	codec.GenerateCodecFile(os.Stdout)
}
