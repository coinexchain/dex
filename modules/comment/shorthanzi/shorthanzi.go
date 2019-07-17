package shorthanzi

import "github.com/pierrec/lz4"

var charMap map[rune]rune

func init() {
	charMap = make(map[rune]rune)
	for i := 0; i < len(oneBytePairs); i += 2 {
		a := oneBytePairs[i]
		b := oneBytePairs[i+1]
		if _, ok := charMap[a]; ok {
			panic("Character is already mapped!")
		}
		if _, ok := charMap[b]; ok {
			panic("Character is already mapped!")
		}
		charMap[a] = b
		charMap[b] = a
	}
	for i := 0; i < len(twoBytePairs); i += 2 {
		a := twoBytePairs[i]
		b := twoBytePairs[i+1]
		if _, ok := charMap[a]; ok {
			panic("Character is already mapped!")
		}
		if _, ok := charMap[b]; ok {
			panic("Character is already mapped!")
		}
		charMap[a] = b
		charMap[b] = a
	}
}

func Transform(in string) string {
	runes := make([]rune, 0, len(in)/2)
	for _, c := range in {
		rep, ok := charMap[c]
		if ok {
			runes = append(runes, rep)
		} else {
			runes = append(runes, c)
		}
	}
	return string(runes)
}

func EncodeHanzi(in string) ([]byte, bool) {
	data := Transform(in)
	return compressText(data)
}

func DecodeHanzi(in []byte) (string, bool) {
	s, ok := decompressText(in)
	if !ok {
		return "", false
	}
	return Transform(s), true
}

func compressText(data string) ([]byte, bool) {
	buf := make([]byte, len(data))

	n, err := lz4.CompressBlockHC([]byte(data), buf, 0)
	if err != nil {
		return nil, false
	}
	if n >= len(data) || n == 0 {
		return nil, false
	}
	return buf[:n], true // compressed data
}

func decompressText(buf []byte) (string, bool) {
	out := make([]byte, 10*len(buf))
	n, err := lz4.UncompressBlock(buf, out)
	if err != nil {
		return "", false
	}
	return string(out[:n]), true // uncompressed data
}
