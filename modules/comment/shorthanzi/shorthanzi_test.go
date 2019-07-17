package shorthanzi

import (
	"fmt"
	"testing"
)

func runTestFail(t *testing.T, text string) {
	s := Transform(text)
	data, ok := EncodeHanzi(text)
	if ok {
		t.Errorf("Encoding should fail, but it does not fail.")
	}
	data2, _ := compressText(text)
	fmt.Printf("text:%d transform:%d compressed:%d direct_compress:%d character_count:%d\n", len(text), len(s), len(data), len(data2), len([]rune(text)))
}

func runTest(t *testing.T, text string) {
	s := Transform(text)
	data, ok := EncodeHanzi(text)
	if !ok {
		t.Errorf("Fail in encoding")
	}
	data2, _ := compressText(text)
	fmt.Printf("text:%d transform:%d compressed:%d direct_compress:%d character_count:%d\n", len(text), len(s), len(data), len(data2), len([]rune(text)))
	textOut, _ := DecodeHanzi(data)
	if textOut != text {
		t.Errorf("Not equal!")
	}
}

func Test1(t *testing.T) {
	runTestFail(t, Text0)
	runTest(t, Text1)
	runTest(t, Text2)
	runTest(t, Text3)
}
