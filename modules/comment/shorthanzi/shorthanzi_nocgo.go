// +build !lz4usecgo

package shorthanzi

import (
	"github.com/bkaradzic/go-lz4"
)

func CompressText(data string) ([]byte, bool) {
	buf := make([]byte, len(data)*2)
	buf, err := lz4.Encode(buf, []byte(data))
	if err != nil {
		return nil, false
	}
	if len(buf) >= len(data) || len(buf) == 0 {
		return nil, false
	}
	return buf, true // compressed data
}

func DecompressText(buf []byte) (string, bool) {
	outSize := 10 * len(buf)
	out := make([]byte, outSize)
	out, err := lz4.Decode(out, buf)
	if err == nil {
		return string(out), true // uncompressed data
	}
	return "", false
}
