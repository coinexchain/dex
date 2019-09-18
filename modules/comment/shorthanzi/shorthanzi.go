package shorthanzi

const maxOutputSize = 160 * 1024

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
	return CompressText(data)
}

func DecodeHanzi(in []byte) (string, bool) {
	s, ok := DecompressText(in)
	if !ok {
		return "", false
	}
	return Transform(s), true
}
