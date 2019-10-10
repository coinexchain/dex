package randsrc

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"math"
	"os"
)

type RandBytesSrcFromFile struct {
	fname   string
	file    *os.File
	scanner *bufio.Scanner
	h       hash.Hash
	sum     []byte
	idx     int
}

func NewRandBytesSrcFromFile(fname string) RandBytesSrcFromFile {
	rs := RandBytesSrcFromFile{}
	rs.fname = fname
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	rs.file = file
	rs.scanner = bufio.NewScanner(rs.file)
	rs.scanner.Buffer(make([]byte, 32), 32)
	rs.h = sha256.New()
	rs.step()
	return rs
}

func (rs *RandBytesSrcFromFile) Close() {
	rs.file.Close()
}

func (rs *RandBytesSrcFromFile) step() {
	if !rs.scanner.Scan() {
		file, err := os.Open(rs.fname)
		if err != nil {
			panic(err)
		}
		rs.file.Close()
		rs.file = file
		rs.scanner = bufio.NewScanner(rs.file)
		rs.scanner.Buffer(make([]byte, 32), 32)
		rs.scanner.Scan()
	}
	rs.h.Write(rs.scanner.Bytes())
	rs.sum = rs.h.Sum(nil)
	rs.idx = 0
}

func (rs *RandBytesSrcFromFile) GetBytes(n int) []byte {
	res := make([]byte, 0, n)
	for len(res) < n {
		res = append(res, rs.sum[rs.idx])
		rs.idx++
		if rs.idx == len(rs.sum) {
			rs.step()
		}
	}
	return res
}

func (rs *RandBytesSrcFromFile) GetString(n int) string {
	return string(rs.GetBytes(n))
}

type RandSrcFromFile struct {
	RandBytesSrcFromFile
}

func NewRandSrcFromFile(fname string) *RandSrcFromFile {
	var res RandSrcFromFile
	res.RandBytesSrcFromFile = NewRandBytesSrcFromFile(fname)
	return &res
}

func (rs *RandSrcFromFile) GetBool() bool {
	bz := rs.GetBytes(1)
	return bz[0] != 0
}

func (rs *RandSrcFromFile) GetUint8() uint8 {
	bz := rs.GetBytes(1)
	return bz[0]
}

func (rs *RandSrcFromFile) GetUint16() uint16 {
	return binary.LittleEndian.Uint16(rs.GetBytes(2))
}

func (rs *RandSrcFromFile) GetUint32() uint32 {
	return binary.LittleEndian.Uint32(rs.GetBytes(4))
}

func (rs *RandSrcFromFile) GetUint64() uint64 {
	return binary.LittleEndian.Uint64(rs.GetBytes(8))
}

func (rs *RandSrcFromFile) GetInt64() int64 {
	return int64(rs.GetUint64())
}
func (rs *RandSrcFromFile) GetInt32() int32 {
	return int32(rs.GetUint32())
}
func (rs *RandSrcFromFile) GetInt16() int16 {
	return int16(rs.GetUint16())
}
func (rs *RandSrcFromFile) GetInt8() int8 {
	return int8(rs.GetUint8())
}

func (rs *RandSrcFromFile) GetInt() int {
	return int(rs.GetUint64())
}
func (rs *RandSrcFromFile) GetUint() uint {
	return uint(rs.GetUint64())
}

func (rs *RandSrcFromFile) GetFloat64() float64 {
	return math.Float64frombits(rs.GetUint64())
}
func (rs *RandSrcFromFile) GetFloat32() float32 {
	return math.Float32frombits(rs.GetUint32())
}

type RandSrc interface {
	GetBool() bool
	GetInt8() int8
	GetInt16() int16
	GetInt32() int32
	GetInt64() int64
	GetUint8() uint8
	GetUint16() uint16
	GetUint32() uint32
	GetUint64() uint64
	GetFloat32() float32
	GetFloat64() float64
	GetString(n int) string
	GetBytes(n int) []byte
}

var _ RandSrc = &RandSrcFromFile{}
