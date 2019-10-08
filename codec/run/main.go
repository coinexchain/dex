package main

import (
	"os"

	"github.com/coinexchain/dex/codec"
)

func main() {
	//codec.ShowInfo()
	codec.GenerateCodecFile(os.Stdout)
}
