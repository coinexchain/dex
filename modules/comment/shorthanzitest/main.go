package main

import (
	"bufio"
	"fmt"
	"github.com/coinexchain/dex/modules/comment/shorthanzi"
	"log"
	"math/rand"
	"unicode/utf8"
	"os"
	"strings"
)

const BufSize = 64*1024*1024

func testShortHanzi(seed int64, n int32) {
	r := rand.New(rand.NewSource(seed))
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
		os.Exit(2)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte,BufSize), BufSize)
	for {
		randN := r.Int31n(n) + 1
		ok, text := getText(scanner, randN)
		if !ok {
			break
		}
		if !utf8.ValidString(text) || len(text)==0 {
			continue
		}
		testText(text, lineCount)

		lineCount += int(randN)
		for i:=0; i<int(randN); i++ {
			lineCount++
			if lineCount%100000 == 0 {
				fmt.Printf("Line:%d\n", lineCount)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getText(scanner *bufio.Scanner, randN int32) (bool, string) {
	res := make([]string, 0, randN)
	ok := true
	for i := 0; i < int(randN); i++ {
		scanOk := scanner.Scan()
		if !scanOk {
			ok = false
			break
		}
		res = append(res, scanner.Text())
	}
	return ok, strings.Join(res, "")
}

func testText(line string, lineCount int) {
	tline := shorthanzi.Transform(line)
	ttline := shorthanzi.Transform(tline)
	if ttline != line {
		fmt.Printf("TT %d: %s\n", lineCount, line)
	}

	bz, ok := shorthanzi.EncodeHanzi(line)
	if !ok {
		return //When line is too short, it will fail. It is just normal.
	}
	ttline, ok = shorthanzi.DecodeHanzi(bz)
	if !ok {
		fmt.Printf("DE0 %d: %s\n", lineCount, line)
		fmt.Printf("DE1 %d: %s\n", lineCount, ttline)
	}
	if ttline != line {
		fmt.Printf("== %d\n", lineCount)
		fmt.Printf("ref--|%s\n", line)
		fmt.Printf("imp--|%s\n", ttline)
	}
}

func main() {
	testShortHanzi(0, 10)
}
