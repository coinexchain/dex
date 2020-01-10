//nolint
package codec

import (
	"encoding/binary"
	"errors"
	"fmt"
	amino "github.com/coinexchain/codon/wrap-amino"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"math/big"
	"reflect"
	"time"
)

type RandSrc interface {
	GetBool() bool
	GetInt() int
	GetInt8() int8
	GetInt16() int16
	GetInt32() int32
	GetInt64() int64
	GetUint() uint
	GetUint8() uint8
	GetUint16() uint16
	GetUint32() uint32
	GetUint64() uint64
	GetFloat32() float32
	GetFloat64() float64
	GetString(n int) string
	GetBytes(n int) []byte
}

func codonWriteVarint(w *[]byte, v int64) {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], v)
	*w = append(*w, buf[0:n]...)
}
func codonWriteUvarint(w *[]byte, v uint64) {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], v)
	*w = append(*w, buf[0:n]...)
}

func codonEncodeBool(n int, w *[]byte, v bool) {
	codonWriteUvarint(w, uint64(n)<<3)
	if v {
		codonWriteUvarint(w, uint64(1))
	} else {
		codonWriteUvarint(w, uint64(0))
	}
}
func codonEncodeVarint(n int, w *[]byte, v int64) {
	codonWriteUvarint(w, uint64(n)<<3)
	codonWriteVarint(w, int64(v))
}
func codonEncodeInt8(n int, w *[]byte, v int8) {
	codonWriteUvarint(w, uint64(n)<<3)
	codonWriteVarint(w, int64(v))
}
func codonEncodeInt16(n int, w *[]byte, v int16) {
	codonWriteUvarint(w, uint64(n)<<3)
	codonWriteVarint(w, int64(v))
}
func codonEncodeUvarint(n int, w *[]byte, v uint64) {
	codonWriteUvarint(w, uint64(n)<<3)
	codonWriteUvarint(w, v)
}
func codonEncodeUint8(n int, w *[]byte, v uint8) {
	codonWriteUvarint(w, uint64(n)<<3)
	codonWriteUvarint(w, uint64(v))
}
func codonEncodeUint16(n int, w *[]byte, v uint16) {
	codonWriteUvarint(w, uint64(n)<<3)
	codonWriteUvarint(w, uint64(v))
}

func codonEncodeByteSlice(n int, w *[]byte, v []byte) {
	codonWriteUvarint(w, (uint64(n)<<3)|2)
	codonWriteUvarint(w, uint64(len(v)))
	*w = append(*w, v...)
}
func codonEncodeString(n int, w *[]byte, v string) {
	codonEncodeByteSlice(n, w, []byte(v))
}
func codonDecodeBool(bz []byte, n *int, err *error) bool {
	return codonDecodeInt64(bz, n, err) != 0
}
func codonDecodeInt(bz []byte, n *int, err *error) int {
	return int(codonDecodeInt64(bz, n, err))
}
func codonDecodeInt8(bz []byte, n *int, err *error) int8 {
	return int8(codonDecodeInt64(bz, n, err))
}
func codonDecodeInt16(bz []byte, n *int, err *error) int16 {
	return int16(codonDecodeInt64(bz, n, err))
}
func codonDecodeInt32(bz []byte, n *int, err *error) int32 {
	return int32(codonDecodeInt64(bz, n, err))
}
func codonDecodeInt64(bz []byte, m *int, err *error) int64 {
	i, n := binary.Varint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	*err = nil
	return int64(i)
}
func codonDecodeUint(bz []byte, n *int, err *error) uint {
	return uint(codonDecodeUint64(bz, n, err))
}
func codonDecodeUint8(bz []byte, n *int, err *error) uint8 {
	return uint8(codonDecodeUint64(bz, n, err))
}
func codonDecodeUint16(bz []byte, n *int, err *error) uint16 {
	return uint16(codonDecodeUint64(bz, n, err))
}
func codonDecodeUint32(bz []byte, n *int, err *error) uint32 {
	return uint32(codonDecodeUint64(bz, n, err))
}
func codonDecodeUint64(bz []byte, m *int, err *error) uint64 {
	i, n := binary.Uvarint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	*err = nil
	return uint64(i)
}
func codonGetByteSlice(res *[]byte, bz []byte) (int, error) {
	length, n := binary.Uvarint(bz)
	if n == 0 {
		// buf too small
		return n, errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		return n, errors.New("EOF decoding varint")
	}
	if length == 0 {
		*res = nil
		return 0, nil
	}
	bz = bz[n:]
	if len(bz) < int(length) {
		*res = nil
		return 0, errors.New("Not enough bytes to read")
	}
	if *res == nil {
		*res = append(*res, bz[:length]...)
	} else {
		*res = append((*res)[:0], bz[:length]...)
	}
	return n + int(length), nil
}
func codonDecodeString(bz []byte, n *int, err *error) string {
	var res []byte
	*n, *err = codonGetByteSlice(&res, bz)
	return string(res)
}

func init() {
	codec.SetFirstInitFunc(func() {
		amino.Stub = &CodonStub{}
	})
}
func EncodeTime(t time.Time) []byte {
	t = t.UTC()
	sec := t.Unix()
	var buf [20]byte
	n := binary.PutVarint(buf[:], sec)

	nanosec := t.Nanosecond()
	m := binary.PutVarint(buf[n:], int64(nanosec))
	return buf[:n+m]
}

func DecodeTime(bz []byte) (time.Time, int, error) {
	sec, n := binary.Varint(bz)
	var err error
	if n == 0 {
		// buf too small
		err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		err = errors.New("EOF decoding varint")
	}
	if err != nil {
		return time.Unix(sec, 0), n, err
	}

	nanosec, m := binary.Varint(bz[n:])
	if m == 0 {
		// buf too small
		err = errors.New("buffer too small")
	} else if m < 0 {
		// value larger than 64 bits (overflow)
		// and -m is the number of bytes read
		m = -m
		err = errors.New("EOF decoding varint")
	}
	if err != nil {
		return time.Unix(sec, nanosec), n + m, err
	}

	return time.Unix(sec, nanosec).UTC(), n + m, nil
}

func StringToTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

var maxSec = StringToTime("9999-09-29T08:02:06.647266Z").Unix()

func RandTime(r RandSrc) time.Time {
	sec := r.GetInt64()
	nanosec := r.GetInt64()
	if sec < 0 {
		sec = -sec
	}
	if nanosec < 0 {
		nanosec = -nanosec
	}
	nanosec = nanosec % (1000 * 1000 * 1000)
	sec = sec % maxSec
	return time.Unix(sec, nanosec).UTC()
}

func DeepCopyTime(t time.Time) time.Time {
	return t.Add(time.Duration(0))
}

func ByteSliceWithLengthPrefix(bz []byte) []byte {
	buf := make([]byte, binary.MaxVarintLen64+len(bz))
	n := binary.PutUvarint(buf[:], uint64(len(bz)))
	return append(buf[0:n], bz...)
}

func EncodeInt(v sdk.Int) []byte {
	b := byte(0)
	if v.BigInt().Sign() < 0 {
		b = byte(1)
	}
	bz := v.BigInt().Bytes()
	return append(bz, b)
}

func DecodeInt(bz []byte) (v sdk.Int, n int, err error) {
	isNeg := bz[len(bz)-1] != 0
	n = len(bz)
	x := big.NewInt(0)
	z := big.NewInt(0)
	x.SetBytes(bz[:len(bz)-1])
	if isNeg {
		z.Neg(x)
		v = sdk.NewIntFromBigInt(z)
	} else {
		v = sdk.NewIntFromBigInt(x)
	}
	return
}

func RandInt(r RandSrc) sdk.Int {
	res := sdk.NewInt(r.GetInt64())
	count := int(r.GetInt64() % 3)
	for i := 0; i < count; i++ {
		res = res.MulRaw(r.GetInt64())
	}
	return res
}

func DeepCopyInt(i sdk.Int) sdk.Int {
	return i.AddRaw(0)
}

func EncodeDec(v sdk.Dec) []byte {
	b := byte(0)
	if v.Int.Sign() < 0 {
		b = byte(1)
	}
	bz := v.Int.Bytes()
	return append(bz, b)
}

func DecodeDec(bz []byte) (v sdk.Dec, n int, err error) {
	isNeg := bz[len(bz)-1] != 0
	n = len(bz)
	v = sdk.ZeroDec()
	v.Int.SetBytes(bz[:len(bz)-1])
	if isNeg {
		v.Int.Neg(v.Int)
	}
	return
}

func RandDec(r RandSrc) sdk.Dec {
	res := sdk.NewDec(r.GetInt64())
	count := int(r.GetInt64() % 3)
	for i := 0; i < count; i++ {
		res = res.MulInt64(r.GetInt64())
	}
	res = res.QuoInt64(r.GetInt64() & 0xFFFFFFFF)
	return res
}

func DeepCopyDec(d sdk.Dec) sdk.Dec {
	return d.MulInt64(1)
}

// ========= BridgeBegin ============
type CodecImp struct {
	sealed bool
}

var _ amino.Sealer = &CodecImp{}
var _ amino.CodecIfc = &CodecImp{}

func (cdc *CodecImp) MarshalBinaryBare(o interface{}) ([]byte, error) {
	s := CodonStub{}
	return s.MarshalBinaryBare(o)
}
func (cdc *CodecImp) MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error) {
	s := CodonStub{}
	return s.MarshalBinaryLengthPrefixed(o)
}
func (cdc *CodecImp) MarshalBinaryLengthPrefixedWriter(w io.Writer, o interface{}) (n int64, err error) {
	bz, err := cdc.MarshalBinaryLengthPrefixed(o)
	m, err := w.Write(bz)
	return int64(m), err
}
func (cdc *CodecImp) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	s := CodonStub{}
	return s.UnmarshalBinaryBare(bz, ptr)
}
func (cdc *CodecImp) UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error {
	s := CodonStub{}
	return s.UnmarshalBinaryLengthPrefixed(bz, ptr)
}
func (cdc *CodecImp) UnmarshalBinaryLengthPrefixedReader(r io.Reader, ptr interface{}, maxSize int64) (n int64, err error) {
	if maxSize < 0 {
		panic("maxSize cannot be negative.")
	}

	// Read byte-length prefix.
	var l int64
	var buf [binary.MaxVarintLen64]byte
	for i := 0; i < len(buf); i++ {
		_, err = r.Read(buf[i : i+1])
		if err != nil {
			return
		}
		n += 1
		if buf[i]&0x80 == 0 {
			break
		}
		if n >= maxSize {
			err = fmt.Errorf("Read overflow, maxSize is %v but uvarint(length-prefix) is itself greater than maxSize.", maxSize)
		}
	}
	u64, _ := binary.Uvarint(buf[:])
	if err != nil {
		return
	}
	if maxSize > 0 {
		if uint64(maxSize) < u64 {
			err = fmt.Errorf("Read overflow, maxSize is %v but this amino binary object is %v bytes.", maxSize, u64)
			return
		}
		if (maxSize - n) < int64(u64) {
			err = fmt.Errorf("Read overflow, maxSize is %v but this length-prefixed amino binary object is %v+%v bytes.", maxSize, n, u64)
			return
		}
	}
	l = int64(u64)
	if l < 0 {
		err = fmt.Errorf("Read overflow, this implementation can't read this because, why would anyone have this much data?")
	}

	// Read that many bytes.
	var bz = make([]byte, l, l)
	_, err = io.ReadFull(r, bz)
	if err != nil {
		return
	}
	n += l

	// Decode.
	err = cdc.UnmarshalBinaryBare(bz, ptr)
	return
}

//------

func (cdc *CodecImp) MustMarshalBinaryBare(o interface{}) []byte {
	bz, err := cdc.MarshalBinaryBare(o)
	if err != nil {
		panic(err)
	}
	return bz
}
func (cdc *CodecImp) MustMarshalBinaryLengthPrefixed(o interface{}) []byte {
	bz, err := cdc.MarshalBinaryLengthPrefixed(o)
	if err != nil {
		panic(err)
	}
	return bz
}
func (cdc *CodecImp) MustUnmarshalBinaryBare(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinaryBare(bz, ptr)
	if err != nil {
		panic(err)
	}
}
func (cdc *CodecImp) MustUnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
	if err != nil {
		panic(err)
	}
}

// ====================
func derefPtr(v interface{}) reflect.Type {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func (cdc *CodecImp) PrintTypes(out io.Writer) error {
	for _, entry := range GetSupportList() {
		_, err := out.Write([]byte(entry))
		if err != nil {
			return err
		}
		_, err = out.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
func (cdc *CodecImp) RegisterConcrete(o interface{}, name string, copts *amino.ConcreteOptions) {
	if cdc.sealed {
		panic("Codec is already sealed")
	}
	t := derefPtr(o)
	path := t.PkgPath() + "." + t.Name()
	found := false
	for _, entry := range GetSupportList() {
		if path == entry {
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Sprintf("%s is not supported", path))
	}
}
func (cdc *CodecImp) RegisterInterface(o interface{}, _ *amino.InterfaceOptions) {
	if cdc.sealed {
		panic("Codec is already sealed")
	}
	t := derefPtr(o)
	path := t.PkgPath() + "." + t.Name()
	found := false
	for _, entry := range GetSupportList() {
		if path == entry {
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Sprintf("%s is not supported", path))
	}
}
func (cdc *CodecImp) SealImp() {
	if cdc.sealed {
		panic("Codec is already sealed")
	}
	cdc.sealed = true
}

// ========================================

type CodonStub struct {
}

func (_ *CodonStub) NewCodecImp() amino.CodecIfc {
	return &CodecImp{}
}
func (_ *CodonStub) DeepCopy(o interface{}) (r interface{}) {
	r = DeepCopyAny(o)
	return
}

func (_ *CodonStub) MarshalBinaryBare(o interface{}) ([]byte, error) {
	if _, ok := getMagicNumOfVar(o); !ok {
		return nil, errors.New("Not Supported Type")
	}
	buf := make([]byte, 0, 1024)
	EncodeAny(&buf, o)
	return buf, nil
}
func (s *CodonStub) MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error) {
	if _, ok := getMagicNumOfVar(o); !ok {
		return nil, errors.New("Not Supported Type")
	}
	bz, err := s.MarshalBinaryBare(o)
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], uint64(len(bz)))
	return append(buf[:n], bz...), err
}
func (_ *CodonStub) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		panic("Unmarshal expects a pointer")
	}

	if len(bz) <= 4 {
		return fmt.Errorf("Byte slice is too short: %d", len(bz))
	}
	o, _, err := DecodeAny(bz)
	if rv.Elem().Kind() == reflect.Interface {
		AssignIfcPtrFromStruct(ptr, o)
	} else {
		rv.Elem().Set(reflect.ValueOf(o))
	}
	return err
}
func (s *CodonStub) UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error {
	if len(bz) == 0 {
		return errors.New("UnmarshalBinaryLengthPrefixed cannot decode empty bytes")
	}
	// Read byte-length prefix.
	u64, n := binary.Uvarint(bz)
	if n < 0 {
		return fmt.Errorf("Error reading msg byte-length prefix: got code %v", n)
	}
	if u64 > uint64(len(bz)-n) {
		return fmt.Errorf("Not enough bytes to read in UnmarshalBinaryLengthPrefixed, want %v more bytes but only have %v",
			u64, len(bz)-n)
	} else if u64 < uint64(len(bz)-n) {
		return fmt.Errorf("Bytes left over in UnmarshalBinaryLengthPrefixed, should read %v more bytes but have %v",
			u64, len(bz)-n)
	}
	bz = bz[n:]
	return s.UnmarshalBinaryBare(bz, ptr)
}
func (s *CodonStub) MustMarshalBinaryLengthPrefixed(o interface{}) []byte {
	bz, err := s.MarshalBinaryLengthPrefixed(o)
	if err != nil {
		panic(err)
	}
	return bz
}

// ========================================
func (_ *CodonStub) UvarintSize(u uint64) int {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], u)
	return n
}
func (_ *CodonStub) EncodeByteSlice(w io.Writer, bz []byte) error {
	_, err := w.Write(ByteSliceWithLengthPrefix(bz))
	return err
}
func (s *CodonStub) ByteSliceSize(bz []byte) int {
	return s.UvarintSize(uint64(len(bz))) + len(bz)
}
func (_ *CodonStub) EncodeVarint(w io.Writer, i int64) error {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], i)
	_, err := w.Write(buf[:n])
	return err
}
func (s *CodonStub) EncodeInt8(w io.Writer, i int8) error {
	return s.EncodeVarint(w, int64(i))
}
func (s *CodonStub) EncodeInt16(w io.Writer, i int16) error {
	return s.EncodeVarint(w, int64(i))
}
func (s *CodonStub) EncodeInt32(w io.Writer, i int32) error {
	return s.EncodeVarint(w, int64(i))
}
func (s *CodonStub) EncodeInt64(w io.Writer, i int64) error {
	return s.EncodeVarint(w, int64(i))
}
func (_ *CodonStub) EncodeUvarint(w io.Writer, u uint64) error {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], u)
	_, err := w.Write(buf[:n])
	return err
}
func (s *CodonStub) EncodeByte(w io.Writer, b byte) error {
	return s.EncodeUvarint(w, uint64(b))
}
func (s *CodonStub) EncodeUint8(w io.Writer, u uint8) error {
	return s.EncodeUvarint(w, uint64(u))
}
func (s *CodonStub) EncodeUint16(w io.Writer, u uint16) error {
	return s.EncodeUvarint(w, uint64(u))
}
func (s *CodonStub) EncodeUint32(w io.Writer, u uint32) error {
	return s.EncodeUvarint(w, uint64(u))
}
func (s *CodonStub) EncodeUint64(w io.Writer, u uint64) error {
	return s.EncodeUvarint(w, uint64(u))
}
func (_ *CodonStub) EncodeBool(w io.Writer, b bool) error {
	u := byte(0)
	if b {
		u = byte(1)
	}
	_, err := w.Write([]byte{u})
	return err
}
func (s *CodonStub) EncodeString(w io.Writer, str string) error {
	return s.EncodeByteSlice(w, []byte(str))
}
func (_ *CodonStub) DecodeInt8(bz []byte) (i int8, n int, err error) {
	i = codonDecodeInt8(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeInt16(bz []byte) (i int16, n int, err error) {
	i = codonDecodeInt16(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeInt32(bz []byte) (i int32, n int, err error) {
	i = codonDecodeInt32(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeInt64(bz []byte) (i int64, n int, err error) {
	i = codonDecodeInt64(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeVarint(bz []byte) (i int64, n int, err error) {
	i = codonDecodeInt64(bz, &n, &err)
	return
}
func (s *CodonStub) DecodeByte(bz []byte) (b byte, n int, err error) {
	b = codonDecodeUint8(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint8(bz []byte) (u uint8, n int, err error) {
	u = codonDecodeUint8(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint16(bz []byte) (u uint16, n int, err error) {
	u = codonDecodeUint16(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint32(bz []byte) (u uint32, n int, err error) {
	u = codonDecodeUint32(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint64(bz []byte) (u uint64, n int, err error) {
	u = codonDecodeUint64(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUvarint(bz []byte) (u uint64, n int, err error) {
	u = codonDecodeUint64(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeBool(bz []byte) (b bool, n int, err error) {
	b = codonDecodeBool(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeByteSlice(bz []byte) (bz2 []byte, n int, err error) {
	m, err := codonGetByteSlice(&bz2, bz)
	n += m
	return
}
func (_ *CodonStub) DecodeString(bz []byte) (s string, n int, err error) {
	s = codonDecodeString(bz, &n, &err)
	return
}
func (_ *CodonStub) VarintSize(i int64) int {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], i)
	return n
}

// ========= BridgeEnd ============
// Non-Interface
func EncodePrivKeyEd25519(w *[]byte, v PrivKeyEd25519) {
	codonEncodeByteSlice(0, w, v[:])
} //End of EncodePrivKeyEd25519

func DecodePrivKeyEd25519(bz []byte) (v PrivKeyEd25519, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			o := v[:]
			n, err = codonGetByteSlice(&o, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodePrivKeyEd25519

func RandPrivKeyEd25519(r RandSrc) PrivKeyEd25519 {
	var length int
	var v PrivKeyEd25519
	length = 64
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPrivKeyEd25519

func DeepCopyPrivKeyEd25519(in PrivKeyEd25519) (out PrivKeyEd25519) {
	var length int
	length = len(in)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		out[_0] = in[_0]
	}
	return
} //End of DeepCopyPrivKeyEd25519

// Non-Interface
func EncodePrivKeySecp256k1(w *[]byte, v PrivKeySecp256k1) {
	codonEncodeByteSlice(0, w, v[:])
} //End of EncodePrivKeySecp256k1

func DecodePrivKeySecp256k1(bz []byte) (v PrivKeySecp256k1, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			o := v[:]
			n, err = codonGetByteSlice(&o, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodePrivKeySecp256k1

func RandPrivKeySecp256k1(r RandSrc) PrivKeySecp256k1 {
	var length int
	var v PrivKeySecp256k1
	length = 32
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPrivKeySecp256k1

func DeepCopyPrivKeySecp256k1(in PrivKeySecp256k1) (out PrivKeySecp256k1) {
	var length int
	length = len(in)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		out[_0] = in[_0]
	}
	return
} //End of DeepCopyPrivKeySecp256k1

// Non-Interface
func EncodePubKeyEd25519(w *[]byte, v PubKeyEd25519) {
	codonEncodeByteSlice(0, w, v[:])
} //End of EncodePubKeyEd25519

func DecodePubKeyEd25519(bz []byte) (v PubKeyEd25519, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			o := v[:]
			n, err = codonGetByteSlice(&o, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodePubKeyEd25519

func RandPubKeyEd25519(r RandSrc) PubKeyEd25519 {
	var length int
	var v PubKeyEd25519
	length = 32
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPubKeyEd25519

func DeepCopyPubKeyEd25519(in PubKeyEd25519) (out PubKeyEd25519) {
	var length int
	length = len(in)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		out[_0] = in[_0]
	}
	return
} //End of DeepCopyPubKeyEd25519

// Non-Interface
func EncodePubKeySecp256k1(w *[]byte, v PubKeySecp256k1) {
	codonEncodeByteSlice(0, w, v[:])
} //End of EncodePubKeySecp256k1

func DecodePubKeySecp256k1(bz []byte) (v PubKeySecp256k1, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			o := v[:]
			n, err = codonGetByteSlice(&o, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodePubKeySecp256k1

func RandPubKeySecp256k1(r RandSrc) PubKeySecp256k1 {
	var length int
	var v PubKeySecp256k1
	length = 33
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPubKeySecp256k1

func DeepCopyPubKeySecp256k1(in PubKeySecp256k1) (out PubKeySecp256k1) {
	var length int
	length = len(in)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		out[_0] = in[_0]
	}
	return
} //End of DeepCopyPubKeySecp256k1

// Non-Interface
func EncodePubKeyMultisigThreshold(w *[]byte, v PubKeyMultisigThreshold) {
	codonEncodeUvarint(1, w, uint64(v.K))
	for _0 := 0; _0 < len(v.PubKeys); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			w := make([]byte, 0, 64)
			EncodePubKey(&w, v.PubKeys[_0]) // interface_encode
			return w
		}()) // end of v.PubKeys[_0]
	}
} //End of EncodePubKeyMultisigThreshold

func DecodePubKeyMultisigThreshold(bz []byte) (v PubKeyMultisigThreshold, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.K
			v.K = uint(codonDecodeUint(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.PubKeys
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp PubKey
			tmp, n, err = DecodePubKey(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.PubKeys = append(v.PubKeys, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodePubKeyMultisigThreshold

func RandPubKeyMultisigThreshold(r RandSrc) PubKeyMultisigThreshold {
	var length int
	var v PubKeyMultisigThreshold
	v.K = r.GetUint()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.PubKeys = nil
	} else {
		v.PubKeys = make([]PubKey, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		v.PubKeys[_0] = RandPubKey(r)
	}
	return v
} //End of RandPubKeyMultisigThreshold

func DeepCopyPubKeyMultisigThreshold(in PubKeyMultisigThreshold) (out PubKeyMultisigThreshold) {
	var length int
	out.K = in.K
	length = len(in.PubKeys)
	if length == 0 {
		out.PubKeys = nil
	} else {
		out.PubKeys = make([]PubKey, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		out.PubKeys[_0] = DeepCopyPubKey(in.PubKeys[_0])
	}
	return
} //End of DeepCopyPubKeyMultisigThreshold

// Non-Interface
func EncodeSignedMsgType(w *[]byte, v SignedMsgType) {
	codonEncodeUint8(0, w, uint8(v))
} //End of EncodeSignedMsgType

func DecodeSignedMsgType(bz []byte) (v SignedMsgType, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			v = SignedMsgType(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeSignedMsgType

func RandSignedMsgType(r RandSrc) SignedMsgType {
	var v SignedMsgType
	v = SignedMsgType(r.GetUint8())
	return v
} //End of RandSignedMsgType

func DeepCopySignedMsgType(in SignedMsgType) (out SignedMsgType) {
	out = in
	return
} //End of DeepCopySignedMsgType

// Non-Interface
func EncodeVoteOption(w *[]byte, v VoteOption) {
	codonEncodeUint8(0, w, uint8(v))
} //End of EncodeVoteOption

func DecodeVoteOption(bz []byte) (v VoteOption, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			v = VoteOption(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeVoteOption

func RandVoteOption(r RandSrc) VoteOption {
	var v VoteOption
	v = VoteOption(r.GetUint8())
	return v
} //End of RandVoteOption

func DeepCopyVoteOption(in VoteOption) (out VoteOption) {
	out = in
	return
} //End of DeepCopyVoteOption

// Non-Interface
func EncodeVote(w *[]byte, v Vote) {
	codonEncodeUint8(1, w, uint8(v.Type))
	codonEncodeVarint(2, w, int64(v.Height))
	codonEncodeVarint(3, w, int64(v.Round))
	codonEncodeByteSlice(4, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeByteSlice(1, w, v.BlockID.Hash[:])
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeVarint(1, w, int64(v.BlockID.PartsHeader.Total))
			codonEncodeByteSlice(2, w, v.BlockID.PartsHeader.Hash[:])
			return wBuf
		}()) // end of v.BlockID.PartsHeader
		return wBuf
	}()) // end of v.BlockID
	codonEncodeByteSlice(5, w, EncodeTime(v.Timestamp))
	codonEncodeByteSlice(6, w, v.ValidatorAddress[:])
	codonEncodeVarint(7, w, int64(v.ValidatorIndex))
	codonEncodeByteSlice(8, w, v.Signature[:])
} //End of EncodeVote

func DecodeVote(bz []byte) (v Vote, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Type
			v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Height
			v.Height = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Round
			v.Round = int(codonDecodeInt(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.BlockID
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.BlockID.Hash
						var tmpBz []byte
						n, err = codonGetByteSlice(&tmpBz, bz)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						v.BlockID.Hash = tmpBz
					case 1: // v.BlockID.PartsHeader
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						func(bz []byte) {
							for len(bz) != 0 {
								tag := codonDecodeUint64(bz, &n, &err)
								if err != nil {
									return
								}
								bz = bz[n:]
								total += n
								tag = tag >> 3
								switch tag {
								case 0: // v.BlockID.PartsHeader.Total
									v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
								case 1: // v.BlockID.PartsHeader.Hash
									var tmpBz []byte
									n, err = codonGetByteSlice(&tmpBz, bz)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									v.BlockID.PartsHeader.Hash = tmpBz
								default:
									err = errors.New("Unknown Field")
									return
								}
							} // end for
						}(bz[:l]) // end func
						if err != nil {
							return
						}
						bz = bz[l:]
						n += int(l)
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 4: // v.Timestamp
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.Timestamp, n, err = DecodeTime(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 5: // v.ValidatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddress = tmpBz
		case 6: // v.ValidatorIndex
			v.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 7: // v.Signature
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Signature = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeVote

func RandVote(r RandSrc) Vote {
	var length int
	var v Vote
	v.Type = SignedMsgType(r.GetUint8())
	v.Height = r.GetInt64()
	v.Round = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BlockID.Hash = r.GetBytes(length)
	v.BlockID.PartsHeader.Total = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BlockID.PartsHeader.Hash = r.GetBytes(length)
	// end of v.BlockID.PartsHeader
	// end of v.BlockID
	v.Timestamp = RandTime(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	v.ValidatorIndex = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Signature = r.GetBytes(length)
	return v
} //End of RandVote

func DeepCopyVote(in Vote) (out Vote) {
	var length int
	out.Type = in.Type
	out.Height = in.Height
	out.Round = in.Round
	length = len(in.BlockID.Hash)
	if length == 0 {
		out.BlockID.Hash = nil
	} else {
		out.BlockID.Hash = make([]uint8, length)
	}
	copy(out.BlockID.Hash[:], in.BlockID.Hash[:])
	out.BlockID.PartsHeader.Total = in.BlockID.PartsHeader.Total
	length = len(in.BlockID.PartsHeader.Hash)
	if length == 0 {
		out.BlockID.PartsHeader.Hash = nil
	} else {
		out.BlockID.PartsHeader.Hash = make([]uint8, length)
	}
	copy(out.BlockID.PartsHeader.Hash[:], in.BlockID.PartsHeader.Hash[:])
	// end of .BlockID.PartsHeader
	// end of .BlockID
	out.Timestamp = DeepCopyTime(in.Timestamp)
	length = len(in.ValidatorAddress)
	if length == 0 {
		out.ValidatorAddress = nil
	} else {
		out.ValidatorAddress = make([]uint8, length)
	}
	copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
	out.ValidatorIndex = in.ValidatorIndex
	length = len(in.Signature)
	if length == 0 {
		out.Signature = nil
	} else {
		out.Signature = make([]uint8, length)
	}
	copy(out.Signature[:], in.Signature[:])
	return
} //End of DeepCopyVote

// Non-Interface
func EncodeSdkInt(w *[]byte, v SdkInt) {
	codonEncodeByteSlice(0, w, EncodeInt(v))
} //End of EncodeSdkInt

func DecodeSdkInt(bz []byte) (v SdkInt, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeSdkInt

func RandSdkInt(r RandSrc) SdkInt {
	var v SdkInt
	v = RandInt(r)
	return v
} //End of RandSdkInt

func DeepCopySdkInt(in SdkInt) (out SdkInt) {
	out = DeepCopyInt(in)
	return
} //End of DeepCopySdkInt

// Non-Interface
func EncodeSdkDec(w *[]byte, v SdkDec) {
	codonEncodeByteSlice(0, w, EncodeDec(v))
} //End of EncodeSdkDec

func DecodeSdkDec(bz []byte) (v SdkDec, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v, n, err = DecodeDec(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeSdkDec

func RandSdkDec(r RandSrc) SdkDec {
	var v SdkDec
	v = RandDec(r)
	return v
} //End of RandSdkDec

func DeepCopySdkDec(in SdkDec) (out SdkDec) {
	out = DeepCopyDec(in)
	return
} //End of DeepCopySdkDec

// Non-Interface
func Encodeuint64(w *[]byte, v uint64) {
	codonEncodeUvarint(0, w, uint64(v))
} //End of Encodeuint64

func Decodeuint64(bz []byte) (v uint64, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			v = uint64(codonDecodeUint64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of Decodeuint64

func Randuint64(r RandSrc) uint64 {
	var v uint64
	v = r.GetUint64()
	return v
} //End of Randuint64

func DeepCopyuint64(in uint64) (out uint64) {
	out = in
	return
} //End of DeepCopyuint64

// Non-Interface
func Encodeint64(w *[]byte, v int64) {
	codonEncodeVarint(0, w, int64(v))
} //End of Encodeint64

func Decodeint64(bz []byte) (v int64, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			v = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of Decodeint64

func Randint64(r RandSrc) int64 {
	var v int64
	v = r.GetInt64()
	return v
} //End of Randint64

func DeepCopyint64(in int64) (out int64) {
	out = in
	return
} //End of DeepCopyint64

// Non-Interface
func EncodeConsAddress(w *[]byte, v ConsAddress) {
	codonEncodeByteSlice(0, w, v[:])
} //End of EncodeConsAddress

func DecodeConsAddress(bz []byte) (v ConsAddress, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeConsAddress

func RandConsAddress(r RandSrc) ConsAddress {
	var length int
	var v ConsAddress
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v = r.GetBytes(length)
	return v
} //End of RandConsAddress

func DeepCopyConsAddress(in ConsAddress) (out ConsAddress) {
	var length int
	length = len(in)
	if length == 0 {
		out = nil
	} else {
		out = make([]uint8, length)
	}
	copy(out[:], in[:])
	return
} //End of DeepCopyConsAddress

// Non-Interface
func EncodeCoin(w *[]byte, v Coin) {
	codonEncodeString(1, w, v.Denom)
	codonEncodeByteSlice(2, w, EncodeInt(v.Amount))
} //End of EncodeCoin

func DecodeCoin(bz []byte) (v Coin, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Denom
			v.Denom = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.Amount, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeCoin

func RandCoin(r RandSrc) Coin {
	var v Coin
	v.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = RandInt(r)
	return v
} //End of RandCoin

func DeepCopyCoin(in Coin) (out Coin) {
	out.Denom = in.Denom
	out.Amount = DeepCopyInt(in.Amount)
	return
} //End of DeepCopyCoin

// Non-Interface
func EncodeDecCoin(w *[]byte, v DecCoin) {
	codonEncodeString(1, w, v.Denom)
	codonEncodeByteSlice(2, w, EncodeDec(v.Amount))
} //End of EncodeDecCoin

func DecodeDecCoin(bz []byte) (v DecCoin, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Denom
			v.Denom = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.Amount, n, err = DecodeDec(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeDecCoin

func RandDecCoin(r RandSrc) DecCoin {
	var v DecCoin
	v.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = RandDec(r)
	return v
} //End of RandDecCoin

func DeepCopyDecCoin(in DecCoin) (out DecCoin) {
	out.Denom = in.Denom
	out.Amount = DeepCopyDec(in.Amount)
	return
} //End of DeepCopyDecCoin

// Non-Interface
func EncodeLockedCoin(w *[]byte, v LockedCoin) {
	codonEncodeByteSlice(1, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Coin.Denom)
		codonEncodeByteSlice(2, w, EncodeInt(v.Coin.Amount))
		return wBuf
	}()) // end of v.Coin
	codonEncodeVarint(2, w, int64(v.UnlockTime))
	codonEncodeByteSlice(3, w, v.FromAddress[:])
	codonEncodeByteSlice(4, w, v.Supervisor[:])
	codonEncodeVarint(5, w, int64(v.Reward))
} //End of EncodeLockedCoin

func DecodeLockedCoin(bz []byte) (v LockedCoin, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Coin
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Coin.Denom
						v.Coin.Denom = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Coin.Amount
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Coin.Amount, n, err = DecodeInt(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 1: // v.UnlockTime
			v.UnlockTime = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.FromAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.FromAddress = tmpBz
		case 3: // v.Supervisor
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Supervisor = tmpBz
		case 4: // v.Reward
			v.Reward = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeLockedCoin

func RandLockedCoin(r RandSrc) LockedCoin {
	var length int
	var v LockedCoin
	v.Coin.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Coin.Amount = RandInt(r)
	// end of v.Coin
	v.UnlockTime = r.GetInt64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Supervisor = r.GetBytes(length)
	v.Reward = r.GetInt64()
	return v
} //End of RandLockedCoin

func DeepCopyLockedCoin(in LockedCoin) (out LockedCoin) {
	var length int
	out.Coin.Denom = in.Coin.Denom
	out.Coin.Amount = DeepCopyInt(in.Coin.Amount)
	// end of .Coin
	out.UnlockTime = in.UnlockTime
	length = len(in.FromAddress)
	if length == 0 {
		out.FromAddress = nil
	} else {
		out.FromAddress = make([]uint8, length)
	}
	copy(out.FromAddress[:], in.FromAddress[:])
	length = len(in.Supervisor)
	if length == 0 {
		out.Supervisor = nil
	} else {
		out.Supervisor = make([]uint8, length)
	}
	copy(out.Supervisor[:], in.Supervisor[:])
	out.Reward = in.Reward
	return
} //End of DeepCopyLockedCoin

// Non-Interface
func EncodeStdSignature(w *[]byte, v StdSignature) {
	codonEncodeByteSlice(1, w, func() []byte {
		w := make([]byte, 0, 64)
		EncodePubKey(&w, v.PubKey) // interface_encode
		return w
	}()) // end of v.PubKey
	codonEncodeByteSlice(2, w, v.Signature[:])
} //End of EncodeStdSignature

func DecodeStdSignature(bz []byte) (v StdSignature, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.PubKey
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.PubKey, n, err = DecodePubKey(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n // interface_decode
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 1: // v.Signature
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Signature = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeStdSignature

func RandStdSignature(r RandSrc) StdSignature {
	var length int
	var v StdSignature
	v.PubKey = RandPubKey(r) // interface_decode
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Signature = r.GetBytes(length)
	return v
} //End of RandStdSignature

func DeepCopyStdSignature(in StdSignature) (out StdSignature) {
	var length int
	out.PubKey = DeepCopyPubKey(in.PubKey)
	length = len(in.Signature)
	if length == 0 {
		out.Signature = nil
	} else {
		out.Signature = make([]uint8, length)
	}
	copy(out.Signature[:], in.Signature[:])
	return
} //End of DeepCopyStdSignature

// Non-Interface
func EncodeParamChange(w *[]byte, v ParamChange) {
	codonEncodeString(1, w, v.Subspace)
	codonEncodeString(2, w, v.Key)
	codonEncodeString(3, w, v.Subkey)
	codonEncodeString(4, w, v.Value)
} //End of EncodeParamChange

func DecodeParamChange(bz []byte) (v ParamChange, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Subspace
			v.Subspace = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Key
			v.Key = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Subkey
			v.Subkey = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.Value
			v.Value = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeParamChange

func RandParamChange(r RandSrc) ParamChange {
	var v ParamChange
	v.Subspace = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Key = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Subkey = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Value = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandParamChange

func DeepCopyParamChange(in ParamChange) (out ParamChange) {
	out.Subspace = in.Subspace
	out.Key = in.Key
	out.Subkey = in.Subkey
	out.Value = in.Value
	return
} //End of DeepCopyParamChange

// Non-Interface
func EncodeInput(w *[]byte, v Input) {
	codonEncodeByteSlice(1, w, v.Address[:])
	for _0 := 0; _0 < len(v.Coins); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Coins[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Coins[_0].Amount))
			return wBuf
		}()) // end of v.Coins[_0]
	}
} //End of EncodeInput

func DecodeInput(bz []byte) (v Input, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Address
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Address = tmpBz
		case 1: // v.Coins
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Coins = append(v.Coins, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeInput

func RandInput(r RandSrc) Input {
	var length int
	var v Input
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Coins = nil
	} else {
		v.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0] = RandCoin(r)
	}
	return v
} //End of RandInput

func DeepCopyInput(in Input) (out Input) {
	var length int
	length = len(in.Address)
	if length == 0 {
		out.Address = nil
	} else {
		out.Address = make([]uint8, length)
	}
	copy(out.Address[:], in.Address[:])
	length = len(in.Coins)
	if length == 0 {
		out.Coins = nil
	} else {
		out.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Coins[_0] = DeepCopyCoin(in.Coins[_0])
	}
	return
} //End of DeepCopyInput

// Non-Interface
func EncodeOutput(w *[]byte, v Output) {
	codonEncodeByteSlice(1, w, v.Address[:])
	for _0 := 0; _0 < len(v.Coins); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Coins[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Coins[_0].Amount))
			return wBuf
		}()) // end of v.Coins[_0]
	}
} //End of EncodeOutput

func DecodeOutput(bz []byte) (v Output, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Address
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Address = tmpBz
		case 1: // v.Coins
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Coins = append(v.Coins, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeOutput

func RandOutput(r RandSrc) Output {
	var length int
	var v Output
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Coins = nil
	} else {
		v.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0] = RandCoin(r)
	}
	return v
} //End of RandOutput

func DeepCopyOutput(in Output) (out Output) {
	var length int
	length = len(in.Address)
	if length == 0 {
		out.Address = nil
	} else {
		out.Address = make([]uint8, length)
	}
	copy(out.Address[:], in.Address[:])
	length = len(in.Coins)
	if length == 0 {
		out.Coins = nil
	} else {
		out.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Coins[_0] = DeepCopyCoin(in.Coins[_0])
	}
	return
} //End of DeepCopyOutput

// Non-Interface
func EncodeAccAddress(w *[]byte, v AccAddress) {
	codonEncodeByteSlice(0, w, v[:])
} //End of EncodeAccAddress

func DecodeAccAddress(bz []byte) (v AccAddress, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeAccAddress

func RandAccAddress(r RandSrc) AccAddress {
	var length int
	var v AccAddress
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v = r.GetBytes(length)
	return v
} //End of RandAccAddress

func DeepCopyAccAddress(in AccAddress) (out AccAddress) {
	var length int
	length = len(in)
	if length == 0 {
		out = nil
	} else {
		out = make([]uint8, length)
	}
	copy(out[:], in[:])
	return
} //End of DeepCopyAccAddress

// Non-Interface
func EncodeCommentRef(w *[]byte, v CommentRef) {
	codonEncodeUvarint(1, w, uint64(v.ID))
	codonEncodeByteSlice(2, w, v.RewardTarget[:])
	codonEncodeString(3, w, v.RewardToken)
	codonEncodeVarint(4, w, int64(v.RewardAmount))
	for _0 := 0; _0 < len(v.Attitudes); _0++ {
		codonEncodeVarint(5, w, int64(v.Attitudes[_0]))
	}
} //End of EncodeCommentRef

func DecodeCommentRef(bz []byte) (v CommentRef, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.ID
			v.ID = uint64(codonDecodeUint64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.RewardTarget
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.RewardTarget = tmpBz
		case 2: // v.RewardToken
			v.RewardToken = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.RewardAmount
			v.RewardAmount = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.Attitudes
			var tmp int32
			tmp = int32(codonDecodeInt32(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Attitudes = append(v.Attitudes, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeCommentRef

func RandCommentRef(r RandSrc) CommentRef {
	var length int
	var v CommentRef
	v.ID = r.GetUint64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.RewardTarget = r.GetBytes(length)
	v.RewardToken = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.RewardAmount = r.GetInt64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Attitudes = nil
	} else {
		v.Attitudes = make([]int32, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of int32
		v.Attitudes[_0] = r.GetInt32()
	}
	return v
} //End of RandCommentRef

func DeepCopyCommentRef(in CommentRef) (out CommentRef) {
	var length int
	out.ID = in.ID
	length = len(in.RewardTarget)
	if length == 0 {
		out.RewardTarget = nil
	} else {
		out.RewardTarget = make([]uint8, length)
	}
	copy(out.RewardTarget[:], in.RewardTarget[:])
	out.RewardToken = in.RewardToken
	out.RewardAmount = in.RewardAmount
	length = len(in.Attitudes)
	if length == 0 {
		out.Attitudes = nil
	} else {
		out.Attitudes = make([]int32, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of int32
		out.Attitudes[_0] = in.Attitudes[_0]
	}
	return
} //End of DeepCopyCommentRef

// Non-Interface
func EncodeBaseAccount(w *[]byte, v BaseAccount) {
	codonEncodeByteSlice(1, w, v.Address[:])
	for _0 := 0; _0 < len(v.Coins); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Coins[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Coins[_0].Amount))
			return wBuf
		}()) // end of v.Coins[_0]
	}
	codonEncodeByteSlice(3, w, func() []byte {
		w := make([]byte, 0, 64)
		EncodePubKey(&w, v.PubKey) // interface_encode
		return w
	}()) // end of v.PubKey
	codonEncodeUvarint(4, w, uint64(v.AccountNumber))
	codonEncodeUvarint(5, w, uint64(v.Sequence))
} //End of EncodeBaseAccount

func DecodeBaseAccount(bz []byte) (v BaseAccount, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Address
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Address = tmpBz
		case 1: // v.Coins
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Coins = append(v.Coins, tmp)
		case 2: // v.PubKey
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.PubKey, n, err = DecodePubKey(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n // interface_decode
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 3: // v.AccountNumber
			v.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.Sequence
			v.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeBaseAccount

func RandBaseAccount(r RandSrc) BaseAccount {
	var length int
	var v BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Coins = nil
	} else {
		v.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0] = RandCoin(r)
	}
	v.PubKey = RandPubKey(r) // interface_decode
	v.AccountNumber = r.GetUint64()
	v.Sequence = r.GetUint64()
	return v
} //End of RandBaseAccount

func DeepCopyBaseAccount(in BaseAccount) (out BaseAccount) {
	var length int
	length = len(in.Address)
	if length == 0 {
		out.Address = nil
	} else {
		out.Address = make([]uint8, length)
	}
	copy(out.Address[:], in.Address[:])
	length = len(in.Coins)
	if length == 0 {
		out.Coins = nil
	} else {
		out.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Coins[_0] = DeepCopyCoin(in.Coins[_0])
	}
	out.PubKey = DeepCopyPubKey(in.PubKey)
	out.AccountNumber = in.AccountNumber
	out.Sequence = in.Sequence
	return
} //End of DeepCopyBaseAccount

// Non-Interface
func EncodeBaseVestingAccount(w *[]byte, v BaseVestingAccount) {
	codonEncodeByteSlice(1, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeByteSlice(1, w, v.BaseAccount.Address[:])
		for _0 := 0; _0 < len(v.BaseAccount.Coins); _0++ {
			codonEncodeByteSlice(2, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseAccount.Coins[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseAccount.Coins[_0].Amount))
				return wBuf
			}()) // end of v.BaseAccount.Coins[_0]
		}
		codonEncodeByteSlice(3, w, func() []byte {
			w := make([]byte, 0, 64)
			EncodePubKey(&w, v.BaseAccount.PubKey) // interface_encode
			return w
		}()) // end of v.BaseAccount.PubKey
		codonEncodeUvarint(4, w, uint64(v.BaseAccount.AccountNumber))
		codonEncodeUvarint(5, w, uint64(v.BaseAccount.Sequence))
		return wBuf
	}()) // end of v.BaseAccount
	for _0 := 0; _0 < len(v.OriginalVesting); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.OriginalVesting[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.OriginalVesting[_0].Amount))
			return wBuf
		}()) // end of v.OriginalVesting[_0]
	}
	for _0 := 0; _0 < len(v.DelegatedFree); _0++ {
		codonEncodeByteSlice(3, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.DelegatedFree[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.DelegatedFree[_0].Amount))
			return wBuf
		}()) // end of v.DelegatedFree[_0]
	}
	for _0 := 0; _0 < len(v.DelegatedVesting); _0++ {
		codonEncodeByteSlice(4, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.DelegatedVesting[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.DelegatedVesting[_0].Amount))
			return wBuf
		}()) // end of v.DelegatedVesting[_0]
	}
	codonEncodeVarint(5, w, int64(v.EndTime))
} //End of EncodeBaseVestingAccount

func DecodeBaseVestingAccount(bz []byte) (v BaseVestingAccount, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.BaseAccount
			v.BaseAccount = &BaseAccount{}
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.BaseAccount.Address
						var tmpBz []byte
						n, err = codonGetByteSlice(&tmpBz, bz)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						v.BaseAccount.Address = tmpBz
					case 1: // v.BaseAccount.Coins
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseAccount.Coins = append(v.BaseAccount.Coins, tmp)
					case 2: // v.BaseAccount.PubKey
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.BaseAccount.PubKey, n, err = DecodePubKey(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n // interface_decode
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					case 3: // v.BaseAccount.AccountNumber
						v.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 4: // v.BaseAccount.Sequence
						v.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 1: // v.OriginalVesting
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.OriginalVesting = append(v.OriginalVesting, tmp)
		case 2: // v.DelegatedFree
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.DelegatedFree = append(v.DelegatedFree, tmp)
		case 3: // v.DelegatedVesting
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.DelegatedVesting = append(v.DelegatedVesting, tmp)
		case 4: // v.EndTime
			v.EndTime = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeBaseVestingAccount

func RandBaseVestingAccount(r RandSrc) BaseVestingAccount {
	var length int
	var v BaseVestingAccount
	v.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseAccount.Coins = nil
	} else {
		v.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.OriginalVesting = nil
	} else {
		v.OriginalVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.OriginalVesting[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.DelegatedFree = nil
	} else {
		v.DelegatedFree = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.DelegatedFree[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.DelegatedVesting = nil
	} else {
		v.DelegatedVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.DelegatedVesting[_0] = RandCoin(r)
	}
	v.EndTime = r.GetInt64()
	return v
} //End of RandBaseVestingAccount

func DeepCopyBaseVestingAccount(in BaseVestingAccount) (out BaseVestingAccount) {
	var length int
	out.BaseAccount = &BaseAccount{}
	length = len(in.BaseAccount.Address)
	if length == 0 {
		out.BaseAccount.Address = nil
	} else {
		out.BaseAccount.Address = make([]uint8, length)
	}
	copy(out.BaseAccount.Address[:], in.BaseAccount.Address[:])
	length = len(in.BaseAccount.Coins)
	if length == 0 {
		out.BaseAccount.Coins = nil
	} else {
		out.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseAccount.Coins[_0] = DeepCopyCoin(in.BaseAccount.Coins[_0])
	}
	out.BaseAccount.PubKey = DeepCopyPubKey(in.BaseAccount.PubKey)
	out.BaseAccount.AccountNumber = in.BaseAccount.AccountNumber
	out.BaseAccount.Sequence = in.BaseAccount.Sequence
	// end of .BaseAccount
	length = len(in.OriginalVesting)
	if length == 0 {
		out.OriginalVesting = nil
	} else {
		out.OriginalVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.OriginalVesting[_0] = DeepCopyCoin(in.OriginalVesting[_0])
	}
	length = len(in.DelegatedFree)
	if length == 0 {
		out.DelegatedFree = nil
	} else {
		out.DelegatedFree = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.DelegatedFree[_0] = DeepCopyCoin(in.DelegatedFree[_0])
	}
	length = len(in.DelegatedVesting)
	if length == 0 {
		out.DelegatedVesting = nil
	} else {
		out.DelegatedVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.DelegatedVesting[_0] = DeepCopyCoin(in.DelegatedVesting[_0])
	}
	out.EndTime = in.EndTime
	return
} //End of DeepCopyBaseVestingAccount

// Non-Interface
func EncodeContinuousVestingAccount(w *[]byte, v ContinuousVestingAccount) {
	codonEncodeByteSlice(1, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeByteSlice(1, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, v.BaseVestingAccount.BaseAccount.Address[:])
			for _0 := 0; _0 < len(v.BaseVestingAccount.BaseAccount.Coins); _0++ {
				codonEncodeByteSlice(2, w, func() []byte {
					wBuf := make([]byte, 0, 64)
					w := &wBuf
					codonEncodeString(1, w, v.BaseVestingAccount.BaseAccount.Coins[_0].Denom)
					codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.BaseAccount.Coins[_0].Amount))
					return wBuf
				}()) // end of v.BaseVestingAccount.BaseAccount.Coins[_0]
			}
			codonEncodeByteSlice(3, w, func() []byte {
				w := make([]byte, 0, 64)
				EncodePubKey(&w, v.BaseVestingAccount.BaseAccount.PubKey) // interface_encode
				return w
			}()) // end of v.BaseVestingAccount.BaseAccount.PubKey
			codonEncodeUvarint(4, w, uint64(v.BaseVestingAccount.BaseAccount.AccountNumber))
			codonEncodeUvarint(5, w, uint64(v.BaseVestingAccount.BaseAccount.Sequence))
			return wBuf
		}()) // end of v.BaseVestingAccount.BaseAccount
		for _0 := 0; _0 < len(v.BaseVestingAccount.OriginalVesting); _0++ {
			codonEncodeByteSlice(2, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseVestingAccount.OriginalVesting[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.OriginalVesting[_0].Amount))
				return wBuf
			}()) // end of v.BaseVestingAccount.OriginalVesting[_0]
		}
		for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedFree); _0++ {
			codonEncodeByteSlice(3, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseVestingAccount.DelegatedFree[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.DelegatedFree[_0].Amount))
				return wBuf
			}()) // end of v.BaseVestingAccount.DelegatedFree[_0]
		}
		for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedVesting); _0++ {
			codonEncodeByteSlice(4, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseVestingAccount.DelegatedVesting[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.DelegatedVesting[_0].Amount))
				return wBuf
			}()) // end of v.BaseVestingAccount.DelegatedVesting[_0]
		}
		codonEncodeVarint(5, w, int64(v.BaseVestingAccount.EndTime))
		return wBuf
	}()) // end of v.BaseVestingAccount
	codonEncodeVarint(2, w, int64(v.StartTime))
} //End of EncodeContinuousVestingAccount

func DecodeContinuousVestingAccount(bz []byte) (v ContinuousVestingAccount, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.BaseVestingAccount
			v.BaseVestingAccount = &BaseVestingAccount{}
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.BaseVestingAccount.BaseAccount
						v.BaseVestingAccount.BaseAccount = &BaseAccount{}
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						func(bz []byte) {
							for len(bz) != 0 {
								tag := codonDecodeUint64(bz, &n, &err)
								if err != nil {
									return
								}
								bz = bz[n:]
								total += n
								tag = tag >> 3
								switch tag {
								case 0: // v.BaseVestingAccount.BaseAccount.Address
									var tmpBz []byte
									n, err = codonGetByteSlice(&tmpBz, bz)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									v.BaseVestingAccount.BaseAccount.Address = tmpBz
								case 1: // v.BaseVestingAccount.BaseAccount.Coins
									l := codonDecodeUint64(bz, &n, &err)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									if int(l) > len(bz) {
										err = errors.New("Length Too Large")
										return
									}
									var tmp Coin
									tmp, n, err = DecodeCoin(bz[:l])
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									if int(l) != n {
										err = errors.New("Length Mismatch")
										return
									}
									v.BaseVestingAccount.BaseAccount.Coins = append(v.BaseVestingAccount.BaseAccount.Coins, tmp)
								case 2: // v.BaseVestingAccount.BaseAccount.PubKey
									l := codonDecodeUint64(bz, &n, &err)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									if int(l) > len(bz) {
										err = errors.New("Length Too Large")
										return
									}
									v.BaseVestingAccount.BaseAccount.PubKey, n, err = DecodePubKey(bz[:l])
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n // interface_decode
									if int(l) != n {
										err = errors.New("Length Mismatch")
										return
									}
								case 3: // v.BaseVestingAccount.BaseAccount.AccountNumber
									v.BaseVestingAccount.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
								case 4: // v.BaseVestingAccount.BaseAccount.Sequence
									v.BaseVestingAccount.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
								default:
									err = errors.New("Unknown Field")
									return
								}
							} // end for
						}(bz[:l]) // end func
						if err != nil {
							return
						}
						bz = bz[l:]
						n += int(l)
					case 1: // v.BaseVestingAccount.OriginalVesting
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseVestingAccount.OriginalVesting = append(v.BaseVestingAccount.OriginalVesting, tmp)
					case 2: // v.BaseVestingAccount.DelegatedFree
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseVestingAccount.DelegatedFree = append(v.BaseVestingAccount.DelegatedFree, tmp)
					case 3: // v.BaseVestingAccount.DelegatedVesting
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseVestingAccount.DelegatedVesting = append(v.BaseVestingAccount.DelegatedVesting, tmp)
					case 4: // v.BaseVestingAccount.EndTime
						v.BaseVestingAccount.EndTime = int64(codonDecodeInt64(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 1: // v.StartTime
			v.StartTime = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeContinuousVestingAccount

func RandContinuousVestingAccount(r RandSrc) ContinuousVestingAccount {
	var length int
	var v ContinuousVestingAccount
	v.BaseVestingAccount = &BaseVestingAccount{}
	v.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.BaseAccount.Coins = nil
	} else {
		v.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseVestingAccount.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseVestingAccount.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseVestingAccount.BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.OriginalVesting = nil
	} else {
		v.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.OriginalVesting[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.DelegatedFree = nil
	} else {
		v.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedFree[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.DelegatedVesting = nil
	} else {
		v.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedVesting[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.EndTime = r.GetInt64()
	// end of v.BaseVestingAccount
	v.StartTime = r.GetInt64()
	return v
} //End of RandContinuousVestingAccount

func DeepCopyContinuousVestingAccount(in ContinuousVestingAccount) (out ContinuousVestingAccount) {
	var length int
	out.BaseVestingAccount = &BaseVestingAccount{}
	out.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = len(in.BaseVestingAccount.BaseAccount.Address)
	if length == 0 {
		out.BaseVestingAccount.BaseAccount.Address = nil
	} else {
		out.BaseVestingAccount.BaseAccount.Address = make([]uint8, length)
	}
	copy(out.BaseVestingAccount.BaseAccount.Address[:], in.BaseVestingAccount.BaseAccount.Address[:])
	length = len(in.BaseVestingAccount.BaseAccount.Coins)
	if length == 0 {
		out.BaseVestingAccount.BaseAccount.Coins = nil
	} else {
		out.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.BaseAccount.Coins[_0] = DeepCopyCoin(in.BaseVestingAccount.BaseAccount.Coins[_0])
	}
	out.BaseVestingAccount.BaseAccount.PubKey = DeepCopyPubKey(in.BaseVestingAccount.BaseAccount.PubKey)
	out.BaseVestingAccount.BaseAccount.AccountNumber = in.BaseVestingAccount.BaseAccount.AccountNumber
	out.BaseVestingAccount.BaseAccount.Sequence = in.BaseVestingAccount.BaseAccount.Sequence
	// end of .BaseVestingAccount.BaseAccount
	length = len(in.BaseVestingAccount.OriginalVesting)
	if length == 0 {
		out.BaseVestingAccount.OriginalVesting = nil
	} else {
		out.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.OriginalVesting[_0] = DeepCopyCoin(in.BaseVestingAccount.OriginalVesting[_0])
	}
	length = len(in.BaseVestingAccount.DelegatedFree)
	if length == 0 {
		out.BaseVestingAccount.DelegatedFree = nil
	} else {
		out.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.DelegatedFree[_0] = DeepCopyCoin(in.BaseVestingAccount.DelegatedFree[_0])
	}
	length = len(in.BaseVestingAccount.DelegatedVesting)
	if length == 0 {
		out.BaseVestingAccount.DelegatedVesting = nil
	} else {
		out.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.DelegatedVesting[_0] = DeepCopyCoin(in.BaseVestingAccount.DelegatedVesting[_0])
	}
	out.BaseVestingAccount.EndTime = in.BaseVestingAccount.EndTime
	// end of .BaseVestingAccount
	out.StartTime = in.StartTime
	return
} //End of DeepCopyContinuousVestingAccount

// Non-Interface
func EncodeDelayedVestingAccount(w *[]byte, v DelayedVestingAccount) {
	codonEncodeByteSlice(1, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeByteSlice(1, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, v.BaseVestingAccount.BaseAccount.Address[:])
			for _0 := 0; _0 < len(v.BaseVestingAccount.BaseAccount.Coins); _0++ {
				codonEncodeByteSlice(2, w, func() []byte {
					wBuf := make([]byte, 0, 64)
					w := &wBuf
					codonEncodeString(1, w, v.BaseVestingAccount.BaseAccount.Coins[_0].Denom)
					codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.BaseAccount.Coins[_0].Amount))
					return wBuf
				}()) // end of v.BaseVestingAccount.BaseAccount.Coins[_0]
			}
			codonEncodeByteSlice(3, w, func() []byte {
				w := make([]byte, 0, 64)
				EncodePubKey(&w, v.BaseVestingAccount.BaseAccount.PubKey) // interface_encode
				return w
			}()) // end of v.BaseVestingAccount.BaseAccount.PubKey
			codonEncodeUvarint(4, w, uint64(v.BaseVestingAccount.BaseAccount.AccountNumber))
			codonEncodeUvarint(5, w, uint64(v.BaseVestingAccount.BaseAccount.Sequence))
			return wBuf
		}()) // end of v.BaseVestingAccount.BaseAccount
		for _0 := 0; _0 < len(v.BaseVestingAccount.OriginalVesting); _0++ {
			codonEncodeByteSlice(2, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseVestingAccount.OriginalVesting[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.OriginalVesting[_0].Amount))
				return wBuf
			}()) // end of v.BaseVestingAccount.OriginalVesting[_0]
		}
		for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedFree); _0++ {
			codonEncodeByteSlice(3, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseVestingAccount.DelegatedFree[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.DelegatedFree[_0].Amount))
				return wBuf
			}()) // end of v.BaseVestingAccount.DelegatedFree[_0]
		}
		for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedVesting); _0++ {
			codonEncodeByteSlice(4, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseVestingAccount.DelegatedVesting[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseVestingAccount.DelegatedVesting[_0].Amount))
				return wBuf
			}()) // end of v.BaseVestingAccount.DelegatedVesting[_0]
		}
		codonEncodeVarint(5, w, int64(v.BaseVestingAccount.EndTime))
		return wBuf
	}()) // end of v.BaseVestingAccount
} //End of EncodeDelayedVestingAccount

func DecodeDelayedVestingAccount(bz []byte) (v DelayedVestingAccount, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.BaseVestingAccount
			v.BaseVestingAccount = &BaseVestingAccount{}
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.BaseVestingAccount.BaseAccount
						v.BaseVestingAccount.BaseAccount = &BaseAccount{}
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						func(bz []byte) {
							for len(bz) != 0 {
								tag := codonDecodeUint64(bz, &n, &err)
								if err != nil {
									return
								}
								bz = bz[n:]
								total += n
								tag = tag >> 3
								switch tag {
								case 0: // v.BaseVestingAccount.BaseAccount.Address
									var tmpBz []byte
									n, err = codonGetByteSlice(&tmpBz, bz)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									v.BaseVestingAccount.BaseAccount.Address = tmpBz
								case 1: // v.BaseVestingAccount.BaseAccount.Coins
									l := codonDecodeUint64(bz, &n, &err)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									if int(l) > len(bz) {
										err = errors.New("Length Too Large")
										return
									}
									var tmp Coin
									tmp, n, err = DecodeCoin(bz[:l])
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									if int(l) != n {
										err = errors.New("Length Mismatch")
										return
									}
									v.BaseVestingAccount.BaseAccount.Coins = append(v.BaseVestingAccount.BaseAccount.Coins, tmp)
								case 2: // v.BaseVestingAccount.BaseAccount.PubKey
									l := codonDecodeUint64(bz, &n, &err)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									if int(l) > len(bz) {
										err = errors.New("Length Too Large")
										return
									}
									v.BaseVestingAccount.BaseAccount.PubKey, n, err = DecodePubKey(bz[:l])
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n // interface_decode
									if int(l) != n {
										err = errors.New("Length Mismatch")
										return
									}
								case 3: // v.BaseVestingAccount.BaseAccount.AccountNumber
									v.BaseVestingAccount.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
								case 4: // v.BaseVestingAccount.BaseAccount.Sequence
									v.BaseVestingAccount.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
								default:
									err = errors.New("Unknown Field")
									return
								}
							} // end for
						}(bz[:l]) // end func
						if err != nil {
							return
						}
						bz = bz[l:]
						n += int(l)
					case 1: // v.BaseVestingAccount.OriginalVesting
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseVestingAccount.OriginalVesting = append(v.BaseVestingAccount.OriginalVesting, tmp)
					case 2: // v.BaseVestingAccount.DelegatedFree
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseVestingAccount.DelegatedFree = append(v.BaseVestingAccount.DelegatedFree, tmp)
					case 3: // v.BaseVestingAccount.DelegatedVesting
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseVestingAccount.DelegatedVesting = append(v.BaseVestingAccount.DelegatedVesting, tmp)
					case 4: // v.BaseVestingAccount.EndTime
						v.BaseVestingAccount.EndTime = int64(codonDecodeInt64(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeDelayedVestingAccount

func RandDelayedVestingAccount(r RandSrc) DelayedVestingAccount {
	var length int
	var v DelayedVestingAccount
	v.BaseVestingAccount = &BaseVestingAccount{}
	v.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.BaseAccount.Coins = nil
	} else {
		v.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseVestingAccount.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseVestingAccount.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseVestingAccount.BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.OriginalVesting = nil
	} else {
		v.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.OriginalVesting[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.DelegatedFree = nil
	} else {
		v.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedFree[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseVestingAccount.DelegatedVesting = nil
	} else {
		v.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedVesting[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.EndTime = r.GetInt64()
	// end of v.BaseVestingAccount
	return v
} //End of RandDelayedVestingAccount

func DeepCopyDelayedVestingAccount(in DelayedVestingAccount) (out DelayedVestingAccount) {
	var length int
	out.BaseVestingAccount = &BaseVestingAccount{}
	out.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = len(in.BaseVestingAccount.BaseAccount.Address)
	if length == 0 {
		out.BaseVestingAccount.BaseAccount.Address = nil
	} else {
		out.BaseVestingAccount.BaseAccount.Address = make([]uint8, length)
	}
	copy(out.BaseVestingAccount.BaseAccount.Address[:], in.BaseVestingAccount.BaseAccount.Address[:])
	length = len(in.BaseVestingAccount.BaseAccount.Coins)
	if length == 0 {
		out.BaseVestingAccount.BaseAccount.Coins = nil
	} else {
		out.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.BaseAccount.Coins[_0] = DeepCopyCoin(in.BaseVestingAccount.BaseAccount.Coins[_0])
	}
	out.BaseVestingAccount.BaseAccount.PubKey = DeepCopyPubKey(in.BaseVestingAccount.BaseAccount.PubKey)
	out.BaseVestingAccount.BaseAccount.AccountNumber = in.BaseVestingAccount.BaseAccount.AccountNumber
	out.BaseVestingAccount.BaseAccount.Sequence = in.BaseVestingAccount.BaseAccount.Sequence
	// end of .BaseVestingAccount.BaseAccount
	length = len(in.BaseVestingAccount.OriginalVesting)
	if length == 0 {
		out.BaseVestingAccount.OriginalVesting = nil
	} else {
		out.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.OriginalVesting[_0] = DeepCopyCoin(in.BaseVestingAccount.OriginalVesting[_0])
	}
	length = len(in.BaseVestingAccount.DelegatedFree)
	if length == 0 {
		out.BaseVestingAccount.DelegatedFree = nil
	} else {
		out.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.DelegatedFree[_0] = DeepCopyCoin(in.BaseVestingAccount.DelegatedFree[_0])
	}
	length = len(in.BaseVestingAccount.DelegatedVesting)
	if length == 0 {
		out.BaseVestingAccount.DelegatedVesting = nil
	} else {
		out.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseVestingAccount.DelegatedVesting[_0] = DeepCopyCoin(in.BaseVestingAccount.DelegatedVesting[_0])
	}
	out.BaseVestingAccount.EndTime = in.BaseVestingAccount.EndTime
	// end of .BaseVestingAccount
	return
} //End of DeepCopyDelayedVestingAccount

// Non-Interface
func EncodeModuleAccount(w *[]byte, v ModuleAccount) {
	codonEncodeByteSlice(1, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeByteSlice(1, w, v.BaseAccount.Address[:])
		for _0 := 0; _0 < len(v.BaseAccount.Coins); _0++ {
			codonEncodeByteSlice(2, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.BaseAccount.Coins[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.BaseAccount.Coins[_0].Amount))
				return wBuf
			}()) // end of v.BaseAccount.Coins[_0]
		}
		codonEncodeByteSlice(3, w, func() []byte {
			w := make([]byte, 0, 64)
			EncodePubKey(&w, v.BaseAccount.PubKey) // interface_encode
			return w
		}()) // end of v.BaseAccount.PubKey
		codonEncodeUvarint(4, w, uint64(v.BaseAccount.AccountNumber))
		codonEncodeUvarint(5, w, uint64(v.BaseAccount.Sequence))
		return wBuf
	}()) // end of v.BaseAccount
	codonEncodeString(2, w, v.Name)
	for _0 := 0; _0 < len(v.Permissions); _0++ {
		codonEncodeString(3, w, v.Permissions[_0])
	}
} //End of EncodeModuleAccount

func DecodeModuleAccount(bz []byte) (v ModuleAccount, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.BaseAccount
			v.BaseAccount = &BaseAccount{}
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.BaseAccount.Address
						var tmpBz []byte
						n, err = codonGetByteSlice(&tmpBz, bz)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						v.BaseAccount.Address = tmpBz
					case 1: // v.BaseAccount.Coins
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.BaseAccount.Coins = append(v.BaseAccount.Coins, tmp)
					case 2: // v.BaseAccount.PubKey
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.BaseAccount.PubKey, n, err = DecodePubKey(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n // interface_decode
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					case 3: // v.BaseAccount.AccountNumber
						v.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 4: // v.BaseAccount.Sequence
						v.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 1: // v.Name
			v.Name = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Permissions
			var tmp string
			tmp = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Permissions = append(v.Permissions, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeModuleAccount

func RandModuleAccount(r RandSrc) ModuleAccount {
	var length int
	var v ModuleAccount
	v.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.BaseAccount.Coins = nil
	} else {
		v.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseAccount
	v.Name = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Permissions = nil
	} else {
		v.Permissions = make([]string, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of string
		v.Permissions[_0] = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	}
	return v
} //End of RandModuleAccount

func DeepCopyModuleAccount(in ModuleAccount) (out ModuleAccount) {
	var length int
	out.BaseAccount = &BaseAccount{}
	length = len(in.BaseAccount.Address)
	if length == 0 {
		out.BaseAccount.Address = nil
	} else {
		out.BaseAccount.Address = make([]uint8, length)
	}
	copy(out.BaseAccount.Address[:], in.BaseAccount.Address[:])
	length = len(in.BaseAccount.Coins)
	if length == 0 {
		out.BaseAccount.Coins = nil
	} else {
		out.BaseAccount.Coins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.BaseAccount.Coins[_0] = DeepCopyCoin(in.BaseAccount.Coins[_0])
	}
	out.BaseAccount.PubKey = DeepCopyPubKey(in.BaseAccount.PubKey)
	out.BaseAccount.AccountNumber = in.BaseAccount.AccountNumber
	out.BaseAccount.Sequence = in.BaseAccount.Sequence
	// end of .BaseAccount
	out.Name = in.Name
	length = len(in.Permissions)
	if length == 0 {
		out.Permissions = nil
	} else {
		out.Permissions = make([]string, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of string
		out.Permissions[_0] = in.Permissions[_0]
	}
	return
} //End of DeepCopyModuleAccount

// Non-Interface
func EncodeStdTx(w *[]byte, v StdTx) {
	for _0 := 0; _0 < len(v.Msgs); _0++ {
		codonEncodeByteSlice(1, w, func() []byte {
			w := make([]byte, 0, 64)
			EncodeMsg(&w, v.Msgs[_0]) // interface_encode
			return w
		}()) // end of v.Msgs[_0]
	}
	codonEncodeByteSlice(2, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		for _0 := 0; _0 < len(v.Fee.Amount); _0++ {
			codonEncodeByteSlice(1, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.Fee.Amount[_0].Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.Fee.Amount[_0].Amount))
				return wBuf
			}()) // end of v.Fee.Amount[_0]
		}
		codonEncodeUvarint(2, w, uint64(v.Fee.Gas))
		return wBuf
	}()) // end of v.Fee
	for _0 := 0; _0 < len(v.Signatures); _0++ {
		codonEncodeByteSlice(3, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, func() []byte {
				w := make([]byte, 0, 64)
				EncodePubKey(&w, v.Signatures[_0].PubKey) // interface_encode
				return w
			}()) // end of v.Signatures[_0].PubKey
			codonEncodeByteSlice(2, w, v.Signatures[_0].Signature[:])
			return wBuf
		}()) // end of v.Signatures[_0]
	}
	codonEncodeString(4, w, v.Memo)
} //End of EncodeStdTx

func DecodeStdTx(bz []byte) (v StdTx, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Msgs
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Msg
			tmp, n, err = DecodeMsg(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Msgs = append(v.Msgs, tmp)
		case 1: // v.Fee
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Fee.Amount
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						var tmp Coin
						tmp, n, err = DecodeCoin(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
						v.Fee.Amount = append(v.Fee.Amount, tmp)
					case 1: // v.Fee.Gas
						v.Fee.Gas = uint64(codonDecodeUint64(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 2: // v.Signatures
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp StdSignature
			tmp, n, err = DecodeStdSignature(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Signatures = append(v.Signatures, tmp)
		case 3: // v.Memo
			v.Memo = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeStdTx

func RandStdTx(r RandSrc) StdTx {
	var length int
	var v StdTx
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Msgs = nil
	} else {
		v.Msgs = make([]Msg, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		v.Msgs[_0] = RandMsg(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Fee.Amount = nil
	} else {
		v.Fee.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Fee.Amount[_0] = RandCoin(r)
	}
	v.Fee.Gas = r.GetUint64()
	// end of v.Fee
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Signatures = nil
	} else {
		v.Signatures = make([]StdSignature, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Signatures[_0] = RandStdSignature(r)
	}
	v.Memo = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandStdTx

func DeepCopyStdTx(in StdTx) (out StdTx) {
	var length int
	length = len(in.Msgs)
	if length == 0 {
		out.Msgs = nil
	} else {
		out.Msgs = make([]Msg, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		out.Msgs[_0] = DeepCopyMsg(in.Msgs[_0])
	}
	length = len(in.Fee.Amount)
	if length == 0 {
		out.Fee.Amount = nil
	} else {
		out.Fee.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Fee.Amount[_0] = DeepCopyCoin(in.Fee.Amount[_0])
	}
	out.Fee.Gas = in.Fee.Gas
	// end of .Fee
	length = len(in.Signatures)
	if length == 0 {
		out.Signatures = nil
	} else {
		out.Signatures = make([]StdSignature, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Signatures[_0] = DeepCopyStdSignature(in.Signatures[_0])
	}
	out.Memo = in.Memo
	return
} //End of DeepCopyStdTx

// Non-Interface
func EncodeMsgBeginRedelegate(w *[]byte, v MsgBeginRedelegate) {
	codonEncodeByteSlice(1, w, v.DelegatorAddress[:])
	codonEncodeByteSlice(2, w, v.ValidatorSrcAddress[:])
	codonEncodeByteSlice(3, w, v.ValidatorDstAddress[:])
	codonEncodeByteSlice(4, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Amount.Denom)
		codonEncodeByteSlice(2, w, EncodeInt(v.Amount.Amount))
		return wBuf
	}()) // end of v.Amount
} //End of EncodeMsgBeginRedelegate

func DecodeMsgBeginRedelegate(bz []byte) (v MsgBeginRedelegate, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.DelegatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.DelegatorAddress = tmpBz
		case 1: // v.ValidatorSrcAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorSrcAddress = tmpBz
		case 2: // v.ValidatorDstAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorDstAddress = tmpBz
		case 3: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Amount.Denom
						v.Amount.Denom = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Amount.Amount
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Amount.Amount, n, err = DecodeInt(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgBeginRedelegate

func RandMsgBeginRedelegate(r RandSrc) MsgBeginRedelegate {
	var length int
	var v MsgBeginRedelegate
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorSrcAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorDstAddress = r.GetBytes(length)
	v.Amount.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount.Amount = RandInt(r)
	// end of v.Amount
	return v
} //End of RandMsgBeginRedelegate

func DeepCopyMsgBeginRedelegate(in MsgBeginRedelegate) (out MsgBeginRedelegate) {
	var length int
	length = len(in.DelegatorAddress)
	if length == 0 {
		out.DelegatorAddress = nil
	} else {
		out.DelegatorAddress = make([]uint8, length)
	}
	copy(out.DelegatorAddress[:], in.DelegatorAddress[:])
	length = len(in.ValidatorSrcAddress)
	if length == 0 {
		out.ValidatorSrcAddress = nil
	} else {
		out.ValidatorSrcAddress = make([]uint8, length)
	}
	copy(out.ValidatorSrcAddress[:], in.ValidatorSrcAddress[:])
	length = len(in.ValidatorDstAddress)
	if length == 0 {
		out.ValidatorDstAddress = nil
	} else {
		out.ValidatorDstAddress = make([]uint8, length)
	}
	copy(out.ValidatorDstAddress[:], in.ValidatorDstAddress[:])
	out.Amount.Denom = in.Amount.Denom
	out.Amount.Amount = DeepCopyInt(in.Amount.Amount)
	// end of .Amount
	return
} //End of DeepCopyMsgBeginRedelegate

// Non-Interface
func EncodeMsgCreateValidator(w *[]byte, v MsgCreateValidator) {
	codonEncodeByteSlice(1, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Description.Moniker)
		codonEncodeString(2, w, v.Description.Identity)
		codonEncodeString(3, w, v.Description.Website)
		codonEncodeString(4, w, v.Description.Details)
		return wBuf
	}()) // end of v.Description
	codonEncodeByteSlice(2, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeByteSlice(1, w, EncodeDec(v.Commission.Rate))
		codonEncodeByteSlice(2, w, EncodeDec(v.Commission.MaxRate))
		codonEncodeByteSlice(3, w, EncodeDec(v.Commission.MaxChangeRate))
		return wBuf
	}()) // end of v.Commission
	codonEncodeByteSlice(3, w, EncodeInt(v.MinSelfDelegation))
	codonEncodeByteSlice(4, w, v.DelegatorAddress[:])
	codonEncodeByteSlice(5, w, v.ValidatorAddress[:])
	codonEncodeByteSlice(6, w, func() []byte {
		w := make([]byte, 0, 64)
		EncodePubKey(&w, v.PubKey) // interface_encode
		return w
	}()) // end of v.PubKey
	codonEncodeByteSlice(7, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Value.Denom)
		codonEncodeByteSlice(2, w, EncodeInt(v.Value.Amount))
		return wBuf
	}()) // end of v.Value
} //End of EncodeMsgCreateValidator

func DecodeMsgCreateValidator(bz []byte) (v MsgCreateValidator, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Description
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Description.Moniker
						v.Description.Moniker = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Description.Identity
						v.Description.Identity = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 2: // v.Description.Website
						v.Description.Website = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 3: // v.Description.Details
						v.Description.Details = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 1: // v.Commission
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Commission.Rate
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Commission.Rate, n, err = DecodeDec(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					case 1: // v.Commission.MaxRate
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Commission.MaxRate, n, err = DecodeDec(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					case 2: // v.Commission.MaxChangeRate
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Commission.MaxChangeRate, n, err = DecodeDec(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 2: // v.MinSelfDelegation
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.MinSelfDelegation, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 3: // v.DelegatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.DelegatorAddress = tmpBz
		case 4: // v.ValidatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddress = tmpBz
		case 5: // v.PubKey
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.PubKey, n, err = DecodePubKey(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n // interface_decode
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 6: // v.Value
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Value.Denom
						v.Value.Denom = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Value.Amount
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Value.Amount, n, err = DecodeInt(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgCreateValidator

func RandMsgCreateValidator(r RandSrc) MsgCreateValidator {
	var length int
	var v MsgCreateValidator
	v.Description.Moniker = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Identity = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Website = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Details = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	// end of v.Description
	v.Commission.Rate = RandDec(r)
	v.Commission.MaxRate = RandDec(r)
	v.Commission.MaxChangeRate = RandDec(r)
	// end of v.Commission
	v.MinSelfDelegation = RandInt(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	v.PubKey = RandPubKey(r) // interface_decode
	v.Value.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Value.Amount = RandInt(r)
	// end of v.Value
	return v
} //End of RandMsgCreateValidator

func DeepCopyMsgCreateValidator(in MsgCreateValidator) (out MsgCreateValidator) {
	var length int
	out.Description.Moniker = in.Description.Moniker
	out.Description.Identity = in.Description.Identity
	out.Description.Website = in.Description.Website
	out.Description.Details = in.Description.Details
	// end of .Description
	out.Commission.Rate = DeepCopyDec(in.Commission.Rate)
	out.Commission.MaxRate = DeepCopyDec(in.Commission.MaxRate)
	out.Commission.MaxChangeRate = DeepCopyDec(in.Commission.MaxChangeRate)
	// end of .Commission
	out.MinSelfDelegation = DeepCopyInt(in.MinSelfDelegation)
	length = len(in.DelegatorAddress)
	if length == 0 {
		out.DelegatorAddress = nil
	} else {
		out.DelegatorAddress = make([]uint8, length)
	}
	copy(out.DelegatorAddress[:], in.DelegatorAddress[:])
	length = len(in.ValidatorAddress)
	if length == 0 {
		out.ValidatorAddress = nil
	} else {
		out.ValidatorAddress = make([]uint8, length)
	}
	copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
	out.PubKey = DeepCopyPubKey(in.PubKey)
	out.Value.Denom = in.Value.Denom
	out.Value.Amount = DeepCopyInt(in.Value.Amount)
	// end of .Value
	return
} //End of DeepCopyMsgCreateValidator

// Non-Interface
func EncodeMsgDelegate(w *[]byte, v MsgDelegate) {
	codonEncodeByteSlice(1, w, v.DelegatorAddress[:])
	codonEncodeByteSlice(2, w, v.ValidatorAddress[:])
	codonEncodeByteSlice(3, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Amount.Denom)
		codonEncodeByteSlice(2, w, EncodeInt(v.Amount.Amount))
		return wBuf
	}()) // end of v.Amount
} //End of EncodeMsgDelegate

func DecodeMsgDelegate(bz []byte) (v MsgDelegate, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.DelegatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.DelegatorAddress = tmpBz
		case 1: // v.ValidatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddress = tmpBz
		case 2: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Amount.Denom
						v.Amount.Denom = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Amount.Amount
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Amount.Amount, n, err = DecodeInt(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgDelegate

func RandMsgDelegate(r RandSrc) MsgDelegate {
	var length int
	var v MsgDelegate
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	v.Amount.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount.Amount = RandInt(r)
	// end of v.Amount
	return v
} //End of RandMsgDelegate

func DeepCopyMsgDelegate(in MsgDelegate) (out MsgDelegate) {
	var length int
	length = len(in.DelegatorAddress)
	if length == 0 {
		out.DelegatorAddress = nil
	} else {
		out.DelegatorAddress = make([]uint8, length)
	}
	copy(out.DelegatorAddress[:], in.DelegatorAddress[:])
	length = len(in.ValidatorAddress)
	if length == 0 {
		out.ValidatorAddress = nil
	} else {
		out.ValidatorAddress = make([]uint8, length)
	}
	copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
	out.Amount.Denom = in.Amount.Denom
	out.Amount.Amount = DeepCopyInt(in.Amount.Amount)
	// end of .Amount
	return
} //End of DeepCopyMsgDelegate

// Non-Interface
func EncodeMsgEditValidator(w *[]byte, v MsgEditValidator) {
	codonEncodeByteSlice(1, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Description.Moniker)
		codonEncodeString(2, w, v.Description.Identity)
		codonEncodeString(3, w, v.Description.Website)
		codonEncodeString(4, w, v.Description.Details)
		return wBuf
	}()) // end of v.Description
	codonEncodeByteSlice(2, w, v.ValidatorAddress[:])
	codonEncodeByteSlice(3, w, EncodeDec(*(v.CommissionRate)))
	codonEncodeByteSlice(4, w, EncodeInt(*(v.MinSelfDelegation)))
} //End of EncodeMsgEditValidator

func DecodeMsgEditValidator(bz []byte) (v MsgEditValidator, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Description
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Description.Moniker
						v.Description.Moniker = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Description.Identity
						v.Description.Identity = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 2: // v.Description.Website
						v.Description.Website = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 3: // v.Description.Details
						v.Description.Details = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 1: // v.ValidatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddress = tmpBz
		case 2: // v.CommissionRate
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.CommissionRate = &SdkDec{}
			*(v.CommissionRate), n, err = DecodeDec(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 3: // v.MinSelfDelegation
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.MinSelfDelegation = &SdkInt{}
			*(v.MinSelfDelegation), n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgEditValidator

func RandMsgEditValidator(r RandSrc) MsgEditValidator {
	var length int
	var v MsgEditValidator
	v.Description.Moniker = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Identity = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Website = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Details = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	// end of v.Description
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	v.CommissionRate = &SdkDec{}
	*(v.CommissionRate) = RandDec(r)
	v.MinSelfDelegation = &SdkInt{}
	*(v.MinSelfDelegation) = RandInt(r)
	return v
} //End of RandMsgEditValidator

func DeepCopyMsgEditValidator(in MsgEditValidator) (out MsgEditValidator) {
	var length int
	out.Description.Moniker = in.Description.Moniker
	out.Description.Identity = in.Description.Identity
	out.Description.Website = in.Description.Website
	out.Description.Details = in.Description.Details
	// end of .Description
	length = len(in.ValidatorAddress)
	if length == 0 {
		out.ValidatorAddress = nil
	} else {
		out.ValidatorAddress = make([]uint8, length)
	}
	copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
	out.CommissionRate = &SdkDec{}
	*(out.CommissionRate) = DeepCopyDec(*(in.CommissionRate))
	out.MinSelfDelegation = &SdkInt{}
	*(out.MinSelfDelegation) = DeepCopyInt(*(in.MinSelfDelegation))
	return
} //End of DeepCopyMsgEditValidator

// Non-Interface
func EncodeMsgSetWithdrawAddress(w *[]byte, v MsgSetWithdrawAddress) {
	codonEncodeByteSlice(1, w, v.DelegatorAddress[:])
	codonEncodeByteSlice(2, w, v.WithdrawAddress[:])
} //End of EncodeMsgSetWithdrawAddress

func DecodeMsgSetWithdrawAddress(bz []byte) (v MsgSetWithdrawAddress, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.DelegatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.DelegatorAddress = tmpBz
		case 1: // v.WithdrawAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.WithdrawAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgSetWithdrawAddress

func RandMsgSetWithdrawAddress(r RandSrc) MsgSetWithdrawAddress {
	var length int
	var v MsgSetWithdrawAddress
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.WithdrawAddress = r.GetBytes(length)
	return v
} //End of RandMsgSetWithdrawAddress

func DeepCopyMsgSetWithdrawAddress(in MsgSetWithdrawAddress) (out MsgSetWithdrawAddress) {
	var length int
	length = len(in.DelegatorAddress)
	if length == 0 {
		out.DelegatorAddress = nil
	} else {
		out.DelegatorAddress = make([]uint8, length)
	}
	copy(out.DelegatorAddress[:], in.DelegatorAddress[:])
	length = len(in.WithdrawAddress)
	if length == 0 {
		out.WithdrawAddress = nil
	} else {
		out.WithdrawAddress = make([]uint8, length)
	}
	copy(out.WithdrawAddress[:], in.WithdrawAddress[:])
	return
} //End of DeepCopyMsgSetWithdrawAddress

// Non-Interface
func EncodeMsgUndelegate(w *[]byte, v MsgUndelegate) {
	codonEncodeByteSlice(1, w, v.DelegatorAddress[:])
	codonEncodeByteSlice(2, w, v.ValidatorAddress[:])
	codonEncodeByteSlice(3, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Amount.Denom)
		codonEncodeByteSlice(2, w, EncodeInt(v.Amount.Amount))
		return wBuf
	}()) // end of v.Amount
} //End of EncodeMsgUndelegate

func DecodeMsgUndelegate(bz []byte) (v MsgUndelegate, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.DelegatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.DelegatorAddress = tmpBz
		case 1: // v.ValidatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddress = tmpBz
		case 2: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Amount.Denom
						v.Amount.Denom = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Amount.Amount
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Amount.Amount, n, err = DecodeInt(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgUndelegate

func RandMsgUndelegate(r RandSrc) MsgUndelegate {
	var length int
	var v MsgUndelegate
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	v.Amount.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount.Amount = RandInt(r)
	// end of v.Amount
	return v
} //End of RandMsgUndelegate

func DeepCopyMsgUndelegate(in MsgUndelegate) (out MsgUndelegate) {
	var length int
	length = len(in.DelegatorAddress)
	if length == 0 {
		out.DelegatorAddress = nil
	} else {
		out.DelegatorAddress = make([]uint8, length)
	}
	copy(out.DelegatorAddress[:], in.DelegatorAddress[:])
	length = len(in.ValidatorAddress)
	if length == 0 {
		out.ValidatorAddress = nil
	} else {
		out.ValidatorAddress = make([]uint8, length)
	}
	copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
	out.Amount.Denom = in.Amount.Denom
	out.Amount.Amount = DeepCopyInt(in.Amount.Amount)
	// end of .Amount
	return
} //End of DeepCopyMsgUndelegate

// Non-Interface
func EncodeMsgUnjail(w *[]byte, v MsgUnjail) {
	codonEncodeByteSlice(1, w, v.ValidatorAddr[:])
} //End of EncodeMsgUnjail

func DecodeMsgUnjail(bz []byte) (v MsgUnjail, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.ValidatorAddr
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddr = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgUnjail

func RandMsgUnjail(r RandSrc) MsgUnjail {
	var length int
	var v MsgUnjail
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddr = r.GetBytes(length)
	return v
} //End of RandMsgUnjail

func DeepCopyMsgUnjail(in MsgUnjail) (out MsgUnjail) {
	var length int
	length = len(in.ValidatorAddr)
	if length == 0 {
		out.ValidatorAddr = nil
	} else {
		out.ValidatorAddr = make([]uint8, length)
	}
	copy(out.ValidatorAddr[:], in.ValidatorAddr[:])
	return
} //End of DeepCopyMsgUnjail

// Non-Interface
func EncodeMsgWithdrawDelegatorReward(w *[]byte, v MsgWithdrawDelegatorReward) {
	codonEncodeByteSlice(1, w, v.DelegatorAddress[:])
	codonEncodeByteSlice(2, w, v.ValidatorAddress[:])
} //End of EncodeMsgWithdrawDelegatorReward

func DecodeMsgWithdrawDelegatorReward(bz []byte) (v MsgWithdrawDelegatorReward, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.DelegatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.DelegatorAddress = tmpBz
		case 1: // v.ValidatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgWithdrawDelegatorReward

func RandMsgWithdrawDelegatorReward(r RandSrc) MsgWithdrawDelegatorReward {
	var length int
	var v MsgWithdrawDelegatorReward
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	return v
} //End of RandMsgWithdrawDelegatorReward

func DeepCopyMsgWithdrawDelegatorReward(in MsgWithdrawDelegatorReward) (out MsgWithdrawDelegatorReward) {
	var length int
	length = len(in.DelegatorAddress)
	if length == 0 {
		out.DelegatorAddress = nil
	} else {
		out.DelegatorAddress = make([]uint8, length)
	}
	copy(out.DelegatorAddress[:], in.DelegatorAddress[:])
	length = len(in.ValidatorAddress)
	if length == 0 {
		out.ValidatorAddress = nil
	} else {
		out.ValidatorAddress = make([]uint8, length)
	}
	copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
	return
} //End of DeepCopyMsgWithdrawDelegatorReward

// Non-Interface
func EncodeMsgWithdrawValidatorCommission(w *[]byte, v MsgWithdrawValidatorCommission) {
	codonEncodeByteSlice(1, w, v.ValidatorAddress[:])
} //End of EncodeMsgWithdrawValidatorCommission

func DecodeMsgWithdrawValidatorCommission(bz []byte) (v MsgWithdrawValidatorCommission, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.ValidatorAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ValidatorAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgWithdrawValidatorCommission

func RandMsgWithdrawValidatorCommission(r RandSrc) MsgWithdrawValidatorCommission {
	var length int
	var v MsgWithdrawValidatorCommission
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	return v
} //End of RandMsgWithdrawValidatorCommission

func DeepCopyMsgWithdrawValidatorCommission(in MsgWithdrawValidatorCommission) (out MsgWithdrawValidatorCommission) {
	var length int
	length = len(in.ValidatorAddress)
	if length == 0 {
		out.ValidatorAddress = nil
	} else {
		out.ValidatorAddress = make([]uint8, length)
	}
	copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
	return
} //End of DeepCopyMsgWithdrawValidatorCommission

// Non-Interface
func EncodeMsgDeposit(w *[]byte, v MsgDeposit) {
	codonEncodeUvarint(1, w, uint64(v.ProposalID))
	codonEncodeByteSlice(2, w, v.Depositor[:])
	for _0 := 0; _0 < len(v.Amount); _0++ {
		codonEncodeByteSlice(3, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Amount[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Amount[_0].Amount))
			return wBuf
		}()) // end of v.Amount[_0]
	}
} //End of EncodeMsgDeposit

func DecodeMsgDeposit(bz []byte) (v MsgDeposit, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.ProposalID
			v.ProposalID = uint64(codonDecodeUint64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Depositor
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Depositor = tmpBz
		case 2: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Amount = append(v.Amount, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgDeposit

func RandMsgDeposit(r RandSrc) MsgDeposit {
	var length int
	var v MsgDeposit
	v.ProposalID = r.GetUint64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Depositor = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Amount = nil
	} else {
		v.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandMsgDeposit

func DeepCopyMsgDeposit(in MsgDeposit) (out MsgDeposit) {
	var length int
	out.ProposalID = in.ProposalID
	length = len(in.Depositor)
	if length == 0 {
		out.Depositor = nil
	} else {
		out.Depositor = make([]uint8, length)
	}
	copy(out.Depositor[:], in.Depositor[:])
	length = len(in.Amount)
	if length == 0 {
		out.Amount = nil
	} else {
		out.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Amount[_0] = DeepCopyCoin(in.Amount[_0])
	}
	return
} //End of DeepCopyMsgDeposit

// Non-Interface
func EncodeMsgSubmitProposal(w *[]byte, v MsgSubmitProposal) {
	codonEncodeByteSlice(1, w, func() []byte {
		w := make([]byte, 0, 64)
		EncodeContent(&w, v.Content) // interface_encode
		return w
	}()) // end of v.Content
	for _0 := 0; _0 < len(v.InitialDeposit); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.InitialDeposit[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.InitialDeposit[_0].Amount))
			return wBuf
		}()) // end of v.InitialDeposit[_0]
	}
	codonEncodeByteSlice(3, w, v.Proposer[:])
} //End of EncodeMsgSubmitProposal

func DecodeMsgSubmitProposal(bz []byte) (v MsgSubmitProposal, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Content
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.Content, n, err = DecodeContent(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n // interface_decode
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 1: // v.InitialDeposit
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.InitialDeposit = append(v.InitialDeposit, tmp)
		case 2: // v.Proposer
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Proposer = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgSubmitProposal

func RandMsgSubmitProposal(r RandSrc) MsgSubmitProposal {
	var length int
	var v MsgSubmitProposal
	v.Content = RandContent(r) // interface_decode
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.InitialDeposit = nil
	} else {
		v.InitialDeposit = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.InitialDeposit[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Proposer = r.GetBytes(length)
	return v
} //End of RandMsgSubmitProposal

func DeepCopyMsgSubmitProposal(in MsgSubmitProposal) (out MsgSubmitProposal) {
	var length int
	out.Content = DeepCopyContent(in.Content)
	length = len(in.InitialDeposit)
	if length == 0 {
		out.InitialDeposit = nil
	} else {
		out.InitialDeposit = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.InitialDeposit[_0] = DeepCopyCoin(in.InitialDeposit[_0])
	}
	length = len(in.Proposer)
	if length == 0 {
		out.Proposer = nil
	} else {
		out.Proposer = make([]uint8, length)
	}
	copy(out.Proposer[:], in.Proposer[:])
	return
} //End of DeepCopyMsgSubmitProposal

// Non-Interface
func EncodeMsgVote(w *[]byte, v MsgVote) {
	codonEncodeUvarint(1, w, uint64(v.ProposalID))
	codonEncodeByteSlice(2, w, v.Voter[:])
	codonEncodeUint8(3, w, uint8(v.Option))
} //End of EncodeMsgVote

func DecodeMsgVote(bz []byte) (v MsgVote, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.ProposalID
			v.ProposalID = uint64(codonDecodeUint64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Voter
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Voter = tmpBz
		case 2: // v.Option
			v.Option = VoteOption(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgVote

func RandMsgVote(r RandSrc) MsgVote {
	var length int
	var v MsgVote
	v.ProposalID = r.GetUint64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Voter = r.GetBytes(length)
	v.Option = VoteOption(r.GetUint8())
	return v
} //End of RandMsgVote

func DeepCopyMsgVote(in MsgVote) (out MsgVote) {
	var length int
	out.ProposalID = in.ProposalID
	length = len(in.Voter)
	if length == 0 {
		out.Voter = nil
	} else {
		out.Voter = make([]uint8, length)
	}
	copy(out.Voter[:], in.Voter[:])
	out.Option = in.Option
	return
} //End of DeepCopyMsgVote

// Non-Interface
func EncodeParameterChangeProposal(w *[]byte, v ParameterChangeProposal) {
	codonEncodeString(1, w, v.Title)
	codonEncodeString(2, w, v.Description)
	for _0 := 0; _0 < len(v.Changes); _0++ {
		codonEncodeByteSlice(3, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Changes[_0].Subspace)
			codonEncodeString(2, w, v.Changes[_0].Key)
			codonEncodeString(3, w, v.Changes[_0].Subkey)
			codonEncodeString(4, w, v.Changes[_0].Value)
			return wBuf
		}()) // end of v.Changes[_0]
	}
} //End of EncodeParameterChangeProposal

func DecodeParameterChangeProposal(bz []byte) (v ParameterChangeProposal, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Title
			v.Title = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Description
			v.Description = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Changes
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp ParamChange
			tmp, n, err = DecodeParamChange(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Changes = append(v.Changes, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeParameterChangeProposal

func RandParameterChangeProposal(r RandSrc) ParameterChangeProposal {
	var length int
	var v ParameterChangeProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Changes = nil
	} else {
		v.Changes = make([]ParamChange, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Changes[_0] = RandParamChange(r)
	}
	return v
} //End of RandParameterChangeProposal

func DeepCopyParameterChangeProposal(in ParameterChangeProposal) (out ParameterChangeProposal) {
	var length int
	out.Title = in.Title
	out.Description = in.Description
	length = len(in.Changes)
	if length == 0 {
		out.Changes = nil
	} else {
		out.Changes = make([]ParamChange, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Changes[_0] = DeepCopyParamChange(in.Changes[_0])
	}
	return
} //End of DeepCopyParameterChangeProposal

// Non-Interface
func EncodeSoftwareUpgradeProposal(w *[]byte, v SoftwareUpgradeProposal) {
	codonEncodeString(1, w, v.Title)
	codonEncodeString(2, w, v.Description)
} //End of EncodeSoftwareUpgradeProposal

func DecodeSoftwareUpgradeProposal(bz []byte) (v SoftwareUpgradeProposal, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Title
			v.Title = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Description
			v.Description = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeSoftwareUpgradeProposal

func RandSoftwareUpgradeProposal(r RandSrc) SoftwareUpgradeProposal {
	var v SoftwareUpgradeProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandSoftwareUpgradeProposal

func DeepCopySoftwareUpgradeProposal(in SoftwareUpgradeProposal) (out SoftwareUpgradeProposal) {
	out.Title = in.Title
	out.Description = in.Description
	return
} //End of DeepCopySoftwareUpgradeProposal

// Non-Interface
func EncodeTextProposal(w *[]byte, v TextProposal) {
	codonEncodeString(1, w, v.Title)
	codonEncodeString(2, w, v.Description)
} //End of EncodeTextProposal

func DecodeTextProposal(bz []byte) (v TextProposal, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Title
			v.Title = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Description
			v.Description = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeTextProposal

func RandTextProposal(r RandSrc) TextProposal {
	var v TextProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandTextProposal

func DeepCopyTextProposal(in TextProposal) (out TextProposal) {
	out.Title = in.Title
	out.Description = in.Description
	return
} //End of DeepCopyTextProposal

// Non-Interface
func EncodeCommunityPoolSpendProposal(w *[]byte, v CommunityPoolSpendProposal) {
	codonEncodeString(1, w, v.Title)
	codonEncodeString(2, w, v.Description)
	codonEncodeByteSlice(3, w, v.Recipient[:])
	for _0 := 0; _0 < len(v.Amount); _0++ {
		codonEncodeByteSlice(4, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Amount[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Amount[_0].Amount))
			return wBuf
		}()) // end of v.Amount[_0]
	}
} //End of EncodeCommunityPoolSpendProposal

func DecodeCommunityPoolSpendProposal(bz []byte) (v CommunityPoolSpendProposal, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Title
			v.Title = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Description
			v.Description = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Recipient
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Recipient = tmpBz
		case 3: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Amount = append(v.Amount, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeCommunityPoolSpendProposal

func RandCommunityPoolSpendProposal(r RandSrc) CommunityPoolSpendProposal {
	var length int
	var v CommunityPoolSpendProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Recipient = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Amount = nil
	} else {
		v.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandCommunityPoolSpendProposal

func DeepCopyCommunityPoolSpendProposal(in CommunityPoolSpendProposal) (out CommunityPoolSpendProposal) {
	var length int
	out.Title = in.Title
	out.Description = in.Description
	length = len(in.Recipient)
	if length == 0 {
		out.Recipient = nil
	} else {
		out.Recipient = make([]uint8, length)
	}
	copy(out.Recipient[:], in.Recipient[:])
	length = len(in.Amount)
	if length == 0 {
		out.Amount = nil
	} else {
		out.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Amount[_0] = DeepCopyCoin(in.Amount[_0])
	}
	return
} //End of DeepCopyCommunityPoolSpendProposal

// Non-Interface
func EncodeMsgMultiSend(w *[]byte, v MsgMultiSend) {
	for _0 := 0; _0 < len(v.Inputs); _0++ {
		codonEncodeByteSlice(1, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, v.Inputs[_0].Address[:])
			for _1 := 0; _1 < len(v.Inputs[_0].Coins); _1++ {
				codonEncodeByteSlice(2, w, func() []byte {
					wBuf := make([]byte, 0, 64)
					w := &wBuf
					codonEncodeString(1, w, v.Inputs[_0].Coins[_1].Denom)
					codonEncodeByteSlice(2, w, EncodeInt(v.Inputs[_0].Coins[_1].Amount))
					return wBuf
				}()) // end of v.Inputs[_0].Coins[_1]
			}
			return wBuf
		}()) // end of v.Inputs[_0]
	}
	for _0 := 0; _0 < len(v.Outputs); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, v.Outputs[_0].Address[:])
			for _1 := 0; _1 < len(v.Outputs[_0].Coins); _1++ {
				codonEncodeByteSlice(2, w, func() []byte {
					wBuf := make([]byte, 0, 64)
					w := &wBuf
					codonEncodeString(1, w, v.Outputs[_0].Coins[_1].Denom)
					codonEncodeByteSlice(2, w, EncodeInt(v.Outputs[_0].Coins[_1].Amount))
					return wBuf
				}()) // end of v.Outputs[_0].Coins[_1]
			}
			return wBuf
		}()) // end of v.Outputs[_0]
	}
} //End of EncodeMsgMultiSend

func DecodeMsgMultiSend(bz []byte) (v MsgMultiSend, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Inputs
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Input
			tmp, n, err = DecodeInput(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Inputs = append(v.Inputs, tmp)
		case 1: // v.Outputs
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Output
			tmp, n, err = DecodeOutput(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Outputs = append(v.Outputs, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgMultiSend

func RandMsgMultiSend(r RandSrc) MsgMultiSend {
	var length int
	var v MsgMultiSend
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Inputs = nil
	} else {
		v.Inputs = make([]Input, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Inputs[_0] = RandInput(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Outputs = nil
	} else {
		v.Outputs = make([]Output, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Outputs[_0] = RandOutput(r)
	}
	return v
} //End of RandMsgMultiSend

func DeepCopyMsgMultiSend(in MsgMultiSend) (out MsgMultiSend) {
	var length int
	length = len(in.Inputs)
	if length == 0 {
		out.Inputs = nil
	} else {
		out.Inputs = make([]Input, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Inputs[_0] = DeepCopyInput(in.Inputs[_0])
	}
	length = len(in.Outputs)
	if length == 0 {
		out.Outputs = nil
	} else {
		out.Outputs = make([]Output, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Outputs[_0] = DeepCopyOutput(in.Outputs[_0])
	}
	return
} //End of DeepCopyMsgMultiSend

// Non-Interface
func EncodeFeePool(w *[]byte, v FeePool) {
	for _0 := 0; _0 < len(v.CommunityPool); _0++ {
		codonEncodeByteSlice(1, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.CommunityPool[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeDec(v.CommunityPool[_0].Amount))
			return wBuf
		}()) // end of v.CommunityPool[_0]
	}
} //End of EncodeFeePool

func DecodeFeePool(bz []byte) (v FeePool, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.CommunityPool
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp DecCoin
			tmp, n, err = DecodeDecCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.CommunityPool = append(v.CommunityPool, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeFeePool

func RandFeePool(r RandSrc) FeePool {
	var length int
	var v FeePool
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.CommunityPool = nil
	} else {
		v.CommunityPool = make([]DecCoin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.CommunityPool[_0] = RandDecCoin(r)
	}
	return v
} //End of RandFeePool

func DeepCopyFeePool(in FeePool) (out FeePool) {
	var length int
	length = len(in.CommunityPool)
	if length == 0 {
		out.CommunityPool = nil
	} else {
		out.CommunityPool = make([]DecCoin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.CommunityPool[_0] = DeepCopyDecCoin(in.CommunityPool[_0])
	}
	return
} //End of DeepCopyFeePool

// Non-Interface
func EncodeMsgSend(w *[]byte, v MsgSend) {
	codonEncodeByteSlice(1, w, v.FromAddress[:])
	codonEncodeByteSlice(2, w, v.ToAddress[:])
	for _0 := 0; _0 < len(v.Amount); _0++ {
		codonEncodeByteSlice(3, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Amount[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Amount[_0].Amount))
			return wBuf
		}()) // end of v.Amount[_0]
	}
} //End of EncodeMsgSend

func DecodeMsgSend(bz []byte) (v MsgSend, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.FromAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.FromAddress = tmpBz
		case 1: // v.ToAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ToAddress = tmpBz
		case 2: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Amount = append(v.Amount, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgSend

func RandMsgSend(r RandSrc) MsgSend {
	var length int
	var v MsgSend
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ToAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Amount = nil
	} else {
		v.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandMsgSend

func DeepCopyMsgSend(in MsgSend) (out MsgSend) {
	var length int
	length = len(in.FromAddress)
	if length == 0 {
		out.FromAddress = nil
	} else {
		out.FromAddress = make([]uint8, length)
	}
	copy(out.FromAddress[:], in.FromAddress[:])
	length = len(in.ToAddress)
	if length == 0 {
		out.ToAddress = nil
	} else {
		out.ToAddress = make([]uint8, length)
	}
	copy(out.ToAddress[:], in.ToAddress[:])
	length = len(in.Amount)
	if length == 0 {
		out.Amount = nil
	} else {
		out.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Amount[_0] = DeepCopyCoin(in.Amount[_0])
	}
	return
} //End of DeepCopyMsgSend

// Non-Interface
func EncodeMsgSupervisedSend(w *[]byte, v MsgSupervisedSend) {
	codonEncodeByteSlice(1, w, v.FromAddress[:])
	codonEncodeByteSlice(2, w, v.Supervisor[:])
	codonEncodeByteSlice(3, w, v.ToAddress[:])
	codonEncodeByteSlice(4, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeString(1, w, v.Amount.Denom)
		codonEncodeByteSlice(2, w, EncodeInt(v.Amount.Amount))
		return wBuf
	}()) // end of v.Amount
	codonEncodeVarint(5, w, int64(v.UnlockTime))
	codonEncodeVarint(6, w, int64(v.Reward))
	codonEncodeUint8(7, w, v.Operation)
} //End of EncodeMsgSupervisedSend

func DecodeMsgSupervisedSend(bz []byte) (v MsgSupervisedSend, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.FromAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.FromAddress = tmpBz
		case 1: // v.Supervisor
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Supervisor = tmpBz
		case 2: // v.ToAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ToAddress = tmpBz
		case 3: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Amount.Denom
						v.Amount.Denom = string(codonDecodeString(bz, &n, &err))
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
					case 1: // v.Amount.Amount
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) > len(bz) {
							err = errors.New("Length Too Large")
							return
						}
						v.Amount.Amount, n, err = DecodeInt(bz[:l])
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						if int(l) != n {
							err = errors.New("Length Mismatch")
							return
						}
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		case 4: // v.UnlockTime
			v.UnlockTime = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 5: // v.Reward
			v.Reward = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 6: // v.Operation
			v.Operation = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgSupervisedSend

func RandMsgSupervisedSend(r RandSrc) MsgSupervisedSend {
	var length int
	var v MsgSupervisedSend
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Supervisor = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ToAddress = r.GetBytes(length)
	v.Amount.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount.Amount = RandInt(r)
	// end of v.Amount
	v.UnlockTime = r.GetInt64()
	v.Reward = r.GetInt64()
	v.Operation = r.GetUint8()
	return v
} //End of RandMsgSupervisedSend

func DeepCopyMsgSupervisedSend(in MsgSupervisedSend) (out MsgSupervisedSend) {
	var length int
	length = len(in.FromAddress)
	if length == 0 {
		out.FromAddress = nil
	} else {
		out.FromAddress = make([]uint8, length)
	}
	copy(out.FromAddress[:], in.FromAddress[:])
	length = len(in.Supervisor)
	if length == 0 {
		out.Supervisor = nil
	} else {
		out.Supervisor = make([]uint8, length)
	}
	copy(out.Supervisor[:], in.Supervisor[:])
	length = len(in.ToAddress)
	if length == 0 {
		out.ToAddress = nil
	} else {
		out.ToAddress = make([]uint8, length)
	}
	copy(out.ToAddress[:], in.ToAddress[:])
	out.Amount.Denom = in.Amount.Denom
	out.Amount.Amount = DeepCopyInt(in.Amount.Amount)
	// end of .Amount
	out.UnlockTime = in.UnlockTime
	out.Reward = in.Reward
	out.Operation = in.Operation
	return
} //End of DeepCopyMsgSupervisedSend

// Non-Interface
func EncodeMsgVerifyInvariant(w *[]byte, v MsgVerifyInvariant) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeString(2, w, v.InvariantModuleName)
	codonEncodeString(3, w, v.InvariantRoute)
} //End of EncodeMsgVerifyInvariant

func DecodeMsgVerifyInvariant(bz []byte) (v MsgVerifyInvariant, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.InvariantModuleName
			v.InvariantModuleName = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.InvariantRoute
			v.InvariantRoute = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgVerifyInvariant

func RandMsgVerifyInvariant(r RandSrc) MsgVerifyInvariant {
	var length int
	var v MsgVerifyInvariant
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.InvariantModuleName = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.InvariantRoute = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandMsgVerifyInvariant

func DeepCopyMsgVerifyInvariant(in MsgVerifyInvariant) (out MsgVerifyInvariant) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.InvariantModuleName = in.InvariantModuleName
	out.InvariantRoute = in.InvariantRoute
	return
} //End of DeepCopyMsgVerifyInvariant

// Non-Interface
func EncodeSupply(w *[]byte, v Supply) {
	for _0 := 0; _0 < len(v.Total); _0++ {
		codonEncodeByteSlice(1, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Total[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Total[_0].Amount))
			return wBuf
		}()) // end of v.Total[_0]
	}
} //End of EncodeSupply

func DecodeSupply(bz []byte) (v Supply, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Total
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Total = append(v.Total, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeSupply

func RandSupply(r RandSrc) Supply {
	var length int
	var v Supply
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Total = nil
	} else {
		v.Total = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Total[_0] = RandCoin(r)
	}
	return v
} //End of RandSupply

func DeepCopySupply(in Supply) (out Supply) {
	var length int
	length = len(in.Total)
	if length == 0 {
		out.Total = nil
	} else {
		out.Total = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Total[_0] = DeepCopyCoin(in.Total[_0])
	}
	return
} //End of DeepCopySupply

// Non-Interface
func EncodeAccountX(w *[]byte, v AccountX) {
	codonEncodeByteSlice(1, w, v.Address[:])
	codonEncodeBool(2, w, v.MemoRequired)
	for _0 := 0; _0 < len(v.LockedCoins); _0++ {
		codonEncodeByteSlice(3, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeString(1, w, v.LockedCoins[_0].Coin.Denom)
				codonEncodeByteSlice(2, w, EncodeInt(v.LockedCoins[_0].Coin.Amount))
				return wBuf
			}()) // end of v.LockedCoins[_0].Coin
			codonEncodeVarint(2, w, int64(v.LockedCoins[_0].UnlockTime))
			codonEncodeByteSlice(3, w, v.LockedCoins[_0].FromAddress[:])
			codonEncodeByteSlice(4, w, v.LockedCoins[_0].Supervisor[:])
			codonEncodeVarint(5, w, int64(v.LockedCoins[_0].Reward))
			return wBuf
		}()) // end of v.LockedCoins[_0]
	}
	for _0 := 0; _0 < len(v.FrozenCoins); _0++ {
		codonEncodeByteSlice(4, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.FrozenCoins[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.FrozenCoins[_0].Amount))
			return wBuf
		}()) // end of v.FrozenCoins[_0]
	}
} //End of EncodeAccountX

func DecodeAccountX(bz []byte) (v AccountX, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Address
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Address = tmpBz
		case 1: // v.MemoRequired
			v.MemoRequired = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.LockedCoins
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp LockedCoin
			tmp, n, err = DecodeLockedCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.LockedCoins = append(v.LockedCoins, tmp)
		case 3: // v.FrozenCoins
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.FrozenCoins = append(v.FrozenCoins, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeAccountX

func RandAccountX(r RandSrc) AccountX {
	var length int
	var v AccountX
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	v.MemoRequired = r.GetBool()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.LockedCoins = nil
	} else {
		v.LockedCoins = make([]LockedCoin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.LockedCoins[_0] = RandLockedCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.FrozenCoins = nil
	} else {
		v.FrozenCoins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.FrozenCoins[_0] = RandCoin(r)
	}
	return v
} //End of RandAccountX

func DeepCopyAccountX(in AccountX) (out AccountX) {
	var length int
	length = len(in.Address)
	if length == 0 {
		out.Address = nil
	} else {
		out.Address = make([]uint8, length)
	}
	copy(out.Address[:], in.Address[:])
	out.MemoRequired = in.MemoRequired
	length = len(in.LockedCoins)
	if length == 0 {
		out.LockedCoins = nil
	} else {
		out.LockedCoins = make([]LockedCoin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.LockedCoins[_0] = DeepCopyLockedCoin(in.LockedCoins[_0])
	}
	length = len(in.FrozenCoins)
	if length == 0 {
		out.FrozenCoins = nil
	} else {
		out.FrozenCoins = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.FrozenCoins[_0] = DeepCopyCoin(in.FrozenCoins[_0])
	}
	return
} //End of DeepCopyAccountX

// Non-Interface
func EncodeMsgMultiSendX(w *[]byte, v MsgMultiSendX) {
	for _0 := 0; _0 < len(v.Inputs); _0++ {
		codonEncodeByteSlice(1, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, v.Inputs[_0].Address[:])
			for _1 := 0; _1 < len(v.Inputs[_0].Coins); _1++ {
				codonEncodeByteSlice(2, w, func() []byte {
					wBuf := make([]byte, 0, 64)
					w := &wBuf
					codonEncodeString(1, w, v.Inputs[_0].Coins[_1].Denom)
					codonEncodeByteSlice(2, w, EncodeInt(v.Inputs[_0].Coins[_1].Amount))
					return wBuf
				}()) // end of v.Inputs[_0].Coins[_1]
			}
			return wBuf
		}()) // end of v.Inputs[_0]
	}
	for _0 := 0; _0 < len(v.Outputs); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeByteSlice(1, w, v.Outputs[_0].Address[:])
			for _1 := 0; _1 < len(v.Outputs[_0].Coins); _1++ {
				codonEncodeByteSlice(2, w, func() []byte {
					wBuf := make([]byte, 0, 64)
					w := &wBuf
					codonEncodeString(1, w, v.Outputs[_0].Coins[_1].Denom)
					codonEncodeByteSlice(2, w, EncodeInt(v.Outputs[_0].Coins[_1].Amount))
					return wBuf
				}()) // end of v.Outputs[_0].Coins[_1]
			}
			return wBuf
		}()) // end of v.Outputs[_0]
	}
} //End of EncodeMsgMultiSendX

func DecodeMsgMultiSendX(bz []byte) (v MsgMultiSendX, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Inputs
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Input
			tmp, n, err = DecodeInput(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Inputs = append(v.Inputs, tmp)
		case 1: // v.Outputs
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Output
			tmp, n, err = DecodeOutput(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Outputs = append(v.Outputs, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgMultiSendX

func RandMsgMultiSendX(r RandSrc) MsgMultiSendX {
	var length int
	var v MsgMultiSendX
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Inputs = nil
	} else {
		v.Inputs = make([]Input, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Inputs[_0] = RandInput(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Outputs = nil
	} else {
		v.Outputs = make([]Output, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Outputs[_0] = RandOutput(r)
	}
	return v
} //End of RandMsgMultiSendX

func DeepCopyMsgMultiSendX(in MsgMultiSendX) (out MsgMultiSendX) {
	var length int
	length = len(in.Inputs)
	if length == 0 {
		out.Inputs = nil
	} else {
		out.Inputs = make([]Input, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Inputs[_0] = DeepCopyInput(in.Inputs[_0])
	}
	length = len(in.Outputs)
	if length == 0 {
		out.Outputs = nil
	} else {
		out.Outputs = make([]Output, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Outputs[_0] = DeepCopyOutput(in.Outputs[_0])
	}
	return
} //End of DeepCopyMsgMultiSendX

// Non-Interface
func EncodeMsgSendX(w *[]byte, v MsgSendX) {
	codonEncodeByteSlice(1, w, v.FromAddress[:])
	codonEncodeByteSlice(2, w, v.ToAddress[:])
	for _0 := 0; _0 < len(v.Amount); _0++ {
		codonEncodeByteSlice(3, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Amount[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Amount[_0].Amount))
			return wBuf
		}()) // end of v.Amount[_0]
	}
	codonEncodeVarint(4, w, int64(v.UnlockTime))
} //End of EncodeMsgSendX

func DecodeMsgSendX(bz []byte) (v MsgSendX, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.FromAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.FromAddress = tmpBz
		case 1: // v.ToAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.ToAddress = tmpBz
		case 2: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Amount = append(v.Amount, tmp)
		case 3: // v.UnlockTime
			v.UnlockTime = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgSendX

func RandMsgSendX(r RandSrc) MsgSendX {
	var length int
	var v MsgSendX
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ToAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Amount = nil
	} else {
		v.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	v.UnlockTime = r.GetInt64()
	return v
} //End of RandMsgSendX

func DeepCopyMsgSendX(in MsgSendX) (out MsgSendX) {
	var length int
	length = len(in.FromAddress)
	if length == 0 {
		out.FromAddress = nil
	} else {
		out.FromAddress = make([]uint8, length)
	}
	copy(out.FromAddress[:], in.FromAddress[:])
	length = len(in.ToAddress)
	if length == 0 {
		out.ToAddress = nil
	} else {
		out.ToAddress = make([]uint8, length)
	}
	copy(out.ToAddress[:], in.ToAddress[:])
	length = len(in.Amount)
	if length == 0 {
		out.Amount = nil
	} else {
		out.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Amount[_0] = DeepCopyCoin(in.Amount[_0])
	}
	out.UnlockTime = in.UnlockTime
	return
} //End of DeepCopyMsgSendX

// Non-Interface
func EncodeMsgSetMemoRequired(w *[]byte, v MsgSetMemoRequired) {
	codonEncodeByteSlice(1, w, v.Address[:])
	codonEncodeBool(2, w, v.Required)
} //End of EncodeMsgSetMemoRequired

func DecodeMsgSetMemoRequired(bz []byte) (v MsgSetMemoRequired, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Address
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Address = tmpBz
		case 1: // v.Required
			v.Required = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgSetMemoRequired

func RandMsgSetMemoRequired(r RandSrc) MsgSetMemoRequired {
	var length int
	var v MsgSetMemoRequired
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	v.Required = r.GetBool()
	return v
} //End of RandMsgSetMemoRequired

func DeepCopyMsgSetMemoRequired(in MsgSetMemoRequired) (out MsgSetMemoRequired) {
	var length int
	length = len(in.Address)
	if length == 0 {
		out.Address = nil
	} else {
		out.Address = make([]uint8, length)
	}
	copy(out.Address[:], in.Address[:])
	out.Required = in.Required
	return
} //End of DeepCopyMsgSetMemoRequired

// Non-Interface
func EncodeBaseToken(w *[]byte, v BaseToken) {
	codonEncodeString(1, w, v.Name)
	codonEncodeString(2, w, v.Symbol)
	codonEncodeByteSlice(3, w, EncodeInt(v.TotalSupply))
	codonEncodeByteSlice(4, w, EncodeInt(v.SendLock))
	codonEncodeByteSlice(5, w, v.Owner[:])
	codonEncodeBool(6, w, v.Mintable)
	codonEncodeBool(7, w, v.Burnable)
	codonEncodeBool(8, w, v.AddrForbiddable)
	codonEncodeBool(9, w, v.TokenForbiddable)
	codonEncodeByteSlice(10, w, EncodeInt(v.TotalBurn))
	codonEncodeByteSlice(11, w, EncodeInt(v.TotalMint))
	codonEncodeBool(12, w, v.IsForbidden)
	codonEncodeString(13, w, v.URL)
	codonEncodeString(14, w, v.Description)
	codonEncodeString(15, w, v.Identity)
} //End of EncodeBaseToken

func DecodeBaseToken(bz []byte) (v BaseToken, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Name
			v.Name = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.TotalSupply
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.TotalSupply, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 3: // v.SendLock
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.SendLock, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 4: // v.Owner
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Owner = tmpBz
		case 5: // v.Mintable
			v.Mintable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 6: // v.Burnable
			v.Burnable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 7: // v.AddrForbiddable
			v.AddrForbiddable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 8: // v.TokenForbiddable
			v.TokenForbiddable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 9: // v.TotalBurn
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.TotalBurn, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 10: // v.TotalMint
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.TotalMint, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 11: // v.IsForbidden
			v.IsForbidden = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 12: // v.URL
			v.URL = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 13: // v.Description
			v.Description = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 14: // v.Identity
			v.Identity = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeBaseToken

func RandBaseToken(r RandSrc) BaseToken {
	var length int
	var v BaseToken
	v.Name = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.TotalSupply = RandInt(r)
	v.SendLock = RandInt(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Owner = r.GetBytes(length)
	v.Mintable = r.GetBool()
	v.Burnable = r.GetBool()
	v.AddrForbiddable = r.GetBool()
	v.TokenForbiddable = r.GetBool()
	v.TotalBurn = RandInt(r)
	v.TotalMint = RandInt(r)
	v.IsForbidden = r.GetBool()
	v.URL = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Identity = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandBaseToken

func DeepCopyBaseToken(in BaseToken) (out BaseToken) {
	var length int
	out.Name = in.Name
	out.Symbol = in.Symbol
	out.TotalSupply = DeepCopyInt(in.TotalSupply)
	out.SendLock = DeepCopyInt(in.SendLock)
	length = len(in.Owner)
	if length == 0 {
		out.Owner = nil
	} else {
		out.Owner = make([]uint8, length)
	}
	copy(out.Owner[:], in.Owner[:])
	out.Mintable = in.Mintable
	out.Burnable = in.Burnable
	out.AddrForbiddable = in.AddrForbiddable
	out.TokenForbiddable = in.TokenForbiddable
	out.TotalBurn = DeepCopyInt(in.TotalBurn)
	out.TotalMint = DeepCopyInt(in.TotalMint)
	out.IsForbidden = in.IsForbidden
	out.URL = in.URL
	out.Description = in.Description
	out.Identity = in.Identity
	return
} //End of DeepCopyBaseToken

// Non-Interface
func EncodeMsgAddTokenWhitelist(w *[]byte, v MsgAddTokenWhitelist) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, v.OwnerAddress[:])
	for _0 := 0; _0 < len(v.Whitelist); _0++ {
		codonEncodeByteSlice(3, w, v.Whitelist[_0][:])
	}
} //End of EncodeMsgAddTokenWhitelist

func DecodeMsgAddTokenWhitelist(bz []byte) (v MsgAddTokenWhitelist, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.OwnerAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddress = tmpBz
		case 2: // v.Whitelist
			var tmp AccAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			tmp = tmpBz
			v.Whitelist = append(v.Whitelist, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgAddTokenWhitelist

func RandMsgAddTokenWhitelist(r RandSrc) MsgAddTokenWhitelist {
	var length int
	var v MsgAddTokenWhitelist
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Whitelist = nil
	} else {
		v.Whitelist = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Whitelist[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgAddTokenWhitelist

func DeepCopyMsgAddTokenWhitelist(in MsgAddTokenWhitelist) (out MsgAddTokenWhitelist) {
	var length int
	out.Symbol = in.Symbol
	length = len(in.OwnerAddress)
	if length == 0 {
		out.OwnerAddress = nil
	} else {
		out.OwnerAddress = make([]uint8, length)
	}
	copy(out.OwnerAddress[:], in.OwnerAddress[:])
	length = len(in.Whitelist)
	if length == 0 {
		out.Whitelist = nil
	} else {
		out.Whitelist = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = len(in.Whitelist[_0])
		if length == 0 {
			out.Whitelist[_0] = nil
		} else {
			out.Whitelist[_0] = make([]uint8, length)
		}
		copy(out.Whitelist[_0][:], in.Whitelist[_0][:])
	}
	return
} //End of DeepCopyMsgAddTokenWhitelist

// Non-Interface
func EncodeMsgBurnToken(w *[]byte, v MsgBurnToken) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, EncodeInt(v.Amount))
	codonEncodeByteSlice(3, w, v.OwnerAddress[:])
} //End of EncodeMsgBurnToken

func DecodeMsgBurnToken(bz []byte) (v MsgBurnToken, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.Amount, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 2: // v.OwnerAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgBurnToken

func RandMsgBurnToken(r RandSrc) MsgBurnToken {
	var length int
	var v MsgBurnToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = RandInt(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgBurnToken

func DeepCopyMsgBurnToken(in MsgBurnToken) (out MsgBurnToken) {
	var length int
	out.Symbol = in.Symbol
	out.Amount = DeepCopyInt(in.Amount)
	length = len(in.OwnerAddress)
	if length == 0 {
		out.OwnerAddress = nil
	} else {
		out.OwnerAddress = make([]uint8, length)
	}
	copy(out.OwnerAddress[:], in.OwnerAddress[:])
	return
} //End of DeepCopyMsgBurnToken

// Non-Interface
func EncodeMsgForbidAddr(w *[]byte, v MsgForbidAddr) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, v.OwnerAddr[:])
	for _0 := 0; _0 < len(v.Addresses); _0++ {
		codonEncodeByteSlice(3, w, v.Addresses[_0][:])
	}
} //End of EncodeMsgForbidAddr

func DecodeMsgForbidAddr(bz []byte) (v MsgForbidAddr, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.OwnerAddr
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddr = tmpBz
		case 2: // v.Addresses
			var tmp AccAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			tmp = tmpBz
			v.Addresses = append(v.Addresses, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgForbidAddr

func RandMsgForbidAddr(r RandSrc) MsgForbidAddr {
	var length int
	var v MsgForbidAddr
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddr = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Addresses = nil
	} else {
		v.Addresses = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Addresses[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgForbidAddr

func DeepCopyMsgForbidAddr(in MsgForbidAddr) (out MsgForbidAddr) {
	var length int
	out.Symbol = in.Symbol
	length = len(in.OwnerAddr)
	if length == 0 {
		out.OwnerAddr = nil
	} else {
		out.OwnerAddr = make([]uint8, length)
	}
	copy(out.OwnerAddr[:], in.OwnerAddr[:])
	length = len(in.Addresses)
	if length == 0 {
		out.Addresses = nil
	} else {
		out.Addresses = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = len(in.Addresses[_0])
		if length == 0 {
			out.Addresses[_0] = nil
		} else {
			out.Addresses[_0] = make([]uint8, length)
		}
		copy(out.Addresses[_0][:], in.Addresses[_0][:])
	}
	return
} //End of DeepCopyMsgForbidAddr

// Non-Interface
func EncodeMsgForbidToken(w *[]byte, v MsgForbidToken) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, v.OwnerAddress[:])
} //End of EncodeMsgForbidToken

func DecodeMsgForbidToken(bz []byte) (v MsgForbidToken, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.OwnerAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgForbidToken

func RandMsgForbidToken(r RandSrc) MsgForbidToken {
	var length int
	var v MsgForbidToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgForbidToken

func DeepCopyMsgForbidToken(in MsgForbidToken) (out MsgForbidToken) {
	var length int
	out.Symbol = in.Symbol
	length = len(in.OwnerAddress)
	if length == 0 {
		out.OwnerAddress = nil
	} else {
		out.OwnerAddress = make([]uint8, length)
	}
	copy(out.OwnerAddress[:], in.OwnerAddress[:])
	return
} //End of DeepCopyMsgForbidToken

// Non-Interface
func EncodeMsgIssueToken(w *[]byte, v MsgIssueToken) {
	codonEncodeString(1, w, v.Name)
	codonEncodeString(2, w, v.Symbol)
	codonEncodeByteSlice(3, w, EncodeInt(v.TotalSupply))
	codonEncodeByteSlice(4, w, v.Owner[:])
	codonEncodeBool(5, w, v.Mintable)
	codonEncodeBool(6, w, v.Burnable)
	codonEncodeBool(7, w, v.AddrForbiddable)
	codonEncodeBool(8, w, v.TokenForbiddable)
	codonEncodeString(9, w, v.URL)
	codonEncodeString(10, w, v.Description)
	codonEncodeString(11, w, v.Identity)
} //End of EncodeMsgIssueToken

func DecodeMsgIssueToken(bz []byte) (v MsgIssueToken, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Name
			v.Name = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.TotalSupply
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.TotalSupply, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 3: // v.Owner
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Owner = tmpBz
		case 4: // v.Mintable
			v.Mintable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 5: // v.Burnable
			v.Burnable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 6: // v.AddrForbiddable
			v.AddrForbiddable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 7: // v.TokenForbiddable
			v.TokenForbiddable = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 8: // v.URL
			v.URL = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 9: // v.Description
			v.Description = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 10: // v.Identity
			v.Identity = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgIssueToken

func RandMsgIssueToken(r RandSrc) MsgIssueToken {
	var length int
	var v MsgIssueToken
	v.Name = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.TotalSupply = RandInt(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Owner = r.GetBytes(length)
	v.Mintable = r.GetBool()
	v.Burnable = r.GetBool()
	v.AddrForbiddable = r.GetBool()
	v.TokenForbiddable = r.GetBool()
	v.URL = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Identity = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandMsgIssueToken

func DeepCopyMsgIssueToken(in MsgIssueToken) (out MsgIssueToken) {
	var length int
	out.Name = in.Name
	out.Symbol = in.Symbol
	out.TotalSupply = DeepCopyInt(in.TotalSupply)
	length = len(in.Owner)
	if length == 0 {
		out.Owner = nil
	} else {
		out.Owner = make([]uint8, length)
	}
	copy(out.Owner[:], in.Owner[:])
	out.Mintable = in.Mintable
	out.Burnable = in.Burnable
	out.AddrForbiddable = in.AddrForbiddable
	out.TokenForbiddable = in.TokenForbiddable
	out.URL = in.URL
	out.Description = in.Description
	out.Identity = in.Identity
	return
} //End of DeepCopyMsgIssueToken

// Non-Interface
func EncodeMsgMintToken(w *[]byte, v MsgMintToken) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, EncodeInt(v.Amount))
	codonEncodeByteSlice(3, w, v.OwnerAddress[:])
} //End of EncodeMsgMintToken

func DecodeMsgMintToken(bz []byte) (v MsgMintToken, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.Amount, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 2: // v.OwnerAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgMintToken

func RandMsgMintToken(r RandSrc) MsgMintToken {
	var length int
	var v MsgMintToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = RandInt(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgMintToken

func DeepCopyMsgMintToken(in MsgMintToken) (out MsgMintToken) {
	var length int
	out.Symbol = in.Symbol
	out.Amount = DeepCopyInt(in.Amount)
	length = len(in.OwnerAddress)
	if length == 0 {
		out.OwnerAddress = nil
	} else {
		out.OwnerAddress = make([]uint8, length)
	}
	copy(out.OwnerAddress[:], in.OwnerAddress[:])
	return
} //End of DeepCopyMsgMintToken

// Non-Interface
func EncodeMsgModifyTokenInfo(w *[]byte, v MsgModifyTokenInfo) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeString(2, w, v.URL)
	codonEncodeString(3, w, v.Description)
	codonEncodeString(4, w, v.Identity)
	codonEncodeByteSlice(5, w, v.OwnerAddress[:])
} //End of EncodeMsgModifyTokenInfo

func DecodeMsgModifyTokenInfo(bz []byte) (v MsgModifyTokenInfo, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.URL
			v.URL = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Description
			v.Description = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.Identity
			v.Identity = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.OwnerAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgModifyTokenInfo

func RandMsgModifyTokenInfo(r RandSrc) MsgModifyTokenInfo {
	var length int
	var v MsgModifyTokenInfo
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.URL = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Identity = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgModifyTokenInfo

func DeepCopyMsgModifyTokenInfo(in MsgModifyTokenInfo) (out MsgModifyTokenInfo) {
	var length int
	out.Symbol = in.Symbol
	out.URL = in.URL
	out.Description = in.Description
	out.Identity = in.Identity
	length = len(in.OwnerAddress)
	if length == 0 {
		out.OwnerAddress = nil
	} else {
		out.OwnerAddress = make([]uint8, length)
	}
	copy(out.OwnerAddress[:], in.OwnerAddress[:])
	return
} //End of DeepCopyMsgModifyTokenInfo

// Non-Interface
func EncodeMsgRemoveTokenWhitelist(w *[]byte, v MsgRemoveTokenWhitelist) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, v.OwnerAddress[:])
	for _0 := 0; _0 < len(v.Whitelist); _0++ {
		codonEncodeByteSlice(3, w, v.Whitelist[_0][:])
	}
} //End of EncodeMsgRemoveTokenWhitelist

func DecodeMsgRemoveTokenWhitelist(bz []byte) (v MsgRemoveTokenWhitelist, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.OwnerAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddress = tmpBz
		case 2: // v.Whitelist
			var tmp AccAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			tmp = tmpBz
			v.Whitelist = append(v.Whitelist, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgRemoveTokenWhitelist

func RandMsgRemoveTokenWhitelist(r RandSrc) MsgRemoveTokenWhitelist {
	var length int
	var v MsgRemoveTokenWhitelist
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Whitelist = nil
	} else {
		v.Whitelist = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Whitelist[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgRemoveTokenWhitelist

func DeepCopyMsgRemoveTokenWhitelist(in MsgRemoveTokenWhitelist) (out MsgRemoveTokenWhitelist) {
	var length int
	out.Symbol = in.Symbol
	length = len(in.OwnerAddress)
	if length == 0 {
		out.OwnerAddress = nil
	} else {
		out.OwnerAddress = make([]uint8, length)
	}
	copy(out.OwnerAddress[:], in.OwnerAddress[:])
	length = len(in.Whitelist)
	if length == 0 {
		out.Whitelist = nil
	} else {
		out.Whitelist = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = len(in.Whitelist[_0])
		if length == 0 {
			out.Whitelist[_0] = nil
		} else {
			out.Whitelist[_0] = make([]uint8, length)
		}
		copy(out.Whitelist[_0][:], in.Whitelist[_0][:])
	}
	return
} //End of DeepCopyMsgRemoveTokenWhitelist

// Non-Interface
func EncodeMsgTransferOwnership(w *[]byte, v MsgTransferOwnership) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, v.OriginalOwner[:])
	codonEncodeByteSlice(3, w, v.NewOwner[:])
} //End of EncodeMsgTransferOwnership

func DecodeMsgTransferOwnership(bz []byte) (v MsgTransferOwnership, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.OriginalOwner
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OriginalOwner = tmpBz
		case 2: // v.NewOwner
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.NewOwner = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgTransferOwnership

func RandMsgTransferOwnership(r RandSrc) MsgTransferOwnership {
	var length int
	var v MsgTransferOwnership
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OriginalOwner = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.NewOwner = r.GetBytes(length)
	return v
} //End of RandMsgTransferOwnership

func DeepCopyMsgTransferOwnership(in MsgTransferOwnership) (out MsgTransferOwnership) {
	var length int
	out.Symbol = in.Symbol
	length = len(in.OriginalOwner)
	if length == 0 {
		out.OriginalOwner = nil
	} else {
		out.OriginalOwner = make([]uint8, length)
	}
	copy(out.OriginalOwner[:], in.OriginalOwner[:])
	length = len(in.NewOwner)
	if length == 0 {
		out.NewOwner = nil
	} else {
		out.NewOwner = make([]uint8, length)
	}
	copy(out.NewOwner[:], in.NewOwner[:])
	return
} //End of DeepCopyMsgTransferOwnership

// Non-Interface
func EncodeMsgUnForbidAddr(w *[]byte, v MsgUnForbidAddr) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, v.OwnerAddr[:])
	for _0 := 0; _0 < len(v.Addresses); _0++ {
		codonEncodeByteSlice(3, w, v.Addresses[_0][:])
	}
} //End of EncodeMsgUnForbidAddr

func DecodeMsgUnForbidAddr(bz []byte) (v MsgUnForbidAddr, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.OwnerAddr
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddr = tmpBz
		case 2: // v.Addresses
			var tmp AccAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			tmp = tmpBz
			v.Addresses = append(v.Addresses, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgUnForbidAddr

func RandMsgUnForbidAddr(r RandSrc) MsgUnForbidAddr {
	var length int
	var v MsgUnForbidAddr
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddr = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Addresses = nil
	} else {
		v.Addresses = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Addresses[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgUnForbidAddr

func DeepCopyMsgUnForbidAddr(in MsgUnForbidAddr) (out MsgUnForbidAddr) {
	var length int
	out.Symbol = in.Symbol
	length = len(in.OwnerAddr)
	if length == 0 {
		out.OwnerAddr = nil
	} else {
		out.OwnerAddr = make([]uint8, length)
	}
	copy(out.OwnerAddr[:], in.OwnerAddr[:])
	length = len(in.Addresses)
	if length == 0 {
		out.Addresses = nil
	} else {
		out.Addresses = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = len(in.Addresses[_0])
		if length == 0 {
			out.Addresses[_0] = nil
		} else {
			out.Addresses[_0] = make([]uint8, length)
		}
		copy(out.Addresses[_0][:], in.Addresses[_0][:])
	}
	return
} //End of DeepCopyMsgUnForbidAddr

// Non-Interface
func EncodeMsgUnForbidToken(w *[]byte, v MsgUnForbidToken) {
	codonEncodeString(1, w, v.Symbol)
	codonEncodeByteSlice(2, w, v.OwnerAddress[:])
} //End of EncodeMsgUnForbidToken

func DecodeMsgUnForbidToken(bz []byte) (v MsgUnForbidToken, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Symbol
			v.Symbol = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.OwnerAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.OwnerAddress = tmpBz
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgUnForbidToken

func RandMsgUnForbidToken(r RandSrc) MsgUnForbidToken {
	var length int
	var v MsgUnForbidToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgUnForbidToken

func DeepCopyMsgUnForbidToken(in MsgUnForbidToken) (out MsgUnForbidToken) {
	var length int
	out.Symbol = in.Symbol
	length = len(in.OwnerAddress)
	if length == 0 {
		out.OwnerAddress = nil
	} else {
		out.OwnerAddress = make([]uint8, length)
	}
	copy(out.OwnerAddress[:], in.OwnerAddress[:])
	return
} //End of DeepCopyMsgUnForbidToken

// Non-Interface
func EncodeMsgBancorCancel(w *[]byte, v MsgBancorCancel) {
	codonEncodeByteSlice(1, w, v.Owner[:])
	codonEncodeString(2, w, v.Stock)
	codonEncodeString(3, w, v.Money)
} //End of EncodeMsgBancorCancel

func DecodeMsgBancorCancel(bz []byte) (v MsgBancorCancel, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Owner
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Owner = tmpBz
		case 1: // v.Stock
			v.Stock = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Money
			v.Money = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgBancorCancel

func RandMsgBancorCancel(r RandSrc) MsgBancorCancel {
	var length int
	var v MsgBancorCancel
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Owner = r.GetBytes(length)
	v.Stock = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Money = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandMsgBancorCancel

func DeepCopyMsgBancorCancel(in MsgBancorCancel) (out MsgBancorCancel) {
	var length int
	length = len(in.Owner)
	if length == 0 {
		out.Owner = nil
	} else {
		out.Owner = make([]uint8, length)
	}
	copy(out.Owner[:], in.Owner[:])
	out.Stock = in.Stock
	out.Money = in.Money
	return
} //End of DeepCopyMsgBancorCancel

// Non-Interface
func EncodeMsgBancorInit(w *[]byte, v MsgBancorInit) {
	codonEncodeByteSlice(1, w, v.Owner[:])
	codonEncodeString(2, w, v.Stock)
	codonEncodeString(3, w, v.Money)
	codonEncodeString(4, w, v.InitPrice)
	codonEncodeByteSlice(5, w, EncodeInt(v.MaxSupply))
	codonEncodeString(6, w, v.MaxPrice)
	codonEncodeUint8(7, w, v.StockPrecision)
	codonEncodeVarint(8, w, int64(v.EarliestCancelTime))
} //End of EncodeMsgBancorInit

func DecodeMsgBancorInit(bz []byte) (v MsgBancorInit, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Owner
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Owner = tmpBz
		case 1: // v.Stock
			v.Stock = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Money
			v.Money = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.InitPrice
			v.InitPrice = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.MaxSupply
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.MaxSupply, n, err = DecodeInt(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 5: // v.MaxPrice
			v.MaxPrice = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 6: // v.StockPrecision
			v.StockPrecision = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 7: // v.EarliestCancelTime
			v.EarliestCancelTime = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgBancorInit

func RandMsgBancorInit(r RandSrc) MsgBancorInit {
	var length int
	var v MsgBancorInit
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Owner = r.GetBytes(length)
	v.Stock = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Money = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.InitPrice = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.MaxSupply = RandInt(r)
	v.MaxPrice = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.StockPrecision = r.GetUint8()
	v.EarliestCancelTime = r.GetInt64()
	return v
} //End of RandMsgBancorInit

func DeepCopyMsgBancorInit(in MsgBancorInit) (out MsgBancorInit) {
	var length int
	length = len(in.Owner)
	if length == 0 {
		out.Owner = nil
	} else {
		out.Owner = make([]uint8, length)
	}
	copy(out.Owner[:], in.Owner[:])
	out.Stock = in.Stock
	out.Money = in.Money
	out.InitPrice = in.InitPrice
	out.MaxSupply = DeepCopyInt(in.MaxSupply)
	out.MaxPrice = in.MaxPrice
	out.StockPrecision = in.StockPrecision
	out.EarliestCancelTime = in.EarliestCancelTime
	return
} //End of DeepCopyMsgBancorInit

// Non-Interface
func EncodeMsgBancorTrade(w *[]byte, v MsgBancorTrade) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeString(2, w, v.Stock)
	codonEncodeString(3, w, v.Money)
	codonEncodeVarint(4, w, int64(v.Amount))
	codonEncodeBool(5, w, v.IsBuy)
	codonEncodeVarint(6, w, int64(v.MoneyLimit))
} //End of EncodeMsgBancorTrade

func DecodeMsgBancorTrade(bz []byte) (v MsgBancorTrade, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.Stock
			v.Stock = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Money
			v.Money = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.Amount
			v.Amount = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.IsBuy
			v.IsBuy = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 5: // v.MoneyLimit
			v.MoneyLimit = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgBancorTrade

func RandMsgBancorTrade(r RandSrc) MsgBancorTrade {
	var length int
	var v MsgBancorTrade
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.Stock = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Money = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = r.GetInt64()
	v.IsBuy = r.GetBool()
	v.MoneyLimit = r.GetInt64()
	return v
} //End of RandMsgBancorTrade

func DeepCopyMsgBancorTrade(in MsgBancorTrade) (out MsgBancorTrade) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.Stock = in.Stock
	out.Money = in.Money
	out.Amount = in.Amount
	out.IsBuy = in.IsBuy
	out.MoneyLimit = in.MoneyLimit
	return
} //End of DeepCopyMsgBancorTrade

// Non-Interface
func EncodeMsgCancelOrder(w *[]byte, v MsgCancelOrder) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeString(2, w, v.OrderID)
} //End of EncodeMsgCancelOrder

func DecodeMsgCancelOrder(bz []byte) (v MsgCancelOrder, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.OrderID
			v.OrderID = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgCancelOrder

func RandMsgCancelOrder(r RandSrc) MsgCancelOrder {
	var length int
	var v MsgCancelOrder
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.OrderID = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandMsgCancelOrder

func DeepCopyMsgCancelOrder(in MsgCancelOrder) (out MsgCancelOrder) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.OrderID = in.OrderID
	return
} //End of DeepCopyMsgCancelOrder

// Non-Interface
func EncodeMsgCancelTradingPair(w *[]byte, v MsgCancelTradingPair) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeString(2, w, v.TradingPair)
	codonEncodeVarint(3, w, int64(v.EffectiveTime))
} //End of EncodeMsgCancelTradingPair

func DecodeMsgCancelTradingPair(bz []byte) (v MsgCancelTradingPair, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.TradingPair
			v.TradingPair = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.EffectiveTime
			v.EffectiveTime = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgCancelTradingPair

func RandMsgCancelTradingPair(r RandSrc) MsgCancelTradingPair {
	var length int
	var v MsgCancelTradingPair
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.TradingPair = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.EffectiveTime = r.GetInt64()
	return v
} //End of RandMsgCancelTradingPair

func DeepCopyMsgCancelTradingPair(in MsgCancelTradingPair) (out MsgCancelTradingPair) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.TradingPair = in.TradingPair
	out.EffectiveTime = in.EffectiveTime
	return
} //End of DeepCopyMsgCancelTradingPair

// Non-Interface
func EncodeMsgCreateOrder(w *[]byte, v MsgCreateOrder) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeUint8(2, w, v.Identify)
	codonEncodeString(3, w, v.TradingPair)
	codonEncodeUint8(4, w, v.OrderType)
	codonEncodeUint8(5, w, v.PricePrecision)
	codonEncodeVarint(6, w, int64(v.Price))
	codonEncodeVarint(7, w, int64(v.Quantity))
	codonEncodeUint8(8, w, v.Side)
	codonEncodeVarint(9, w, int64(v.TimeInForce))
	codonEncodeVarint(10, w, int64(v.ExistBlocks))
} //End of EncodeMsgCreateOrder

func DecodeMsgCreateOrder(bz []byte) (v MsgCreateOrder, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.Identify
			v.Identify = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.TradingPair
			v.TradingPair = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.OrderType
			v.OrderType = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.PricePrecision
			v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 5: // v.Price
			v.Price = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 6: // v.Quantity
			v.Quantity = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 7: // v.Side
			v.Side = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 8: // v.TimeInForce
			v.TimeInForce = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 9: // v.ExistBlocks
			v.ExistBlocks = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgCreateOrder

func RandMsgCreateOrder(r RandSrc) MsgCreateOrder {
	var length int
	var v MsgCreateOrder
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.Identify = r.GetUint8()
	v.TradingPair = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.OrderType = r.GetUint8()
	v.PricePrecision = r.GetUint8()
	v.Price = r.GetInt64()
	v.Quantity = r.GetInt64()
	v.Side = r.GetUint8()
	v.TimeInForce = r.GetInt64()
	v.ExistBlocks = r.GetInt64()
	return v
} //End of RandMsgCreateOrder

func DeepCopyMsgCreateOrder(in MsgCreateOrder) (out MsgCreateOrder) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.Identify = in.Identify
	out.TradingPair = in.TradingPair
	out.OrderType = in.OrderType
	out.PricePrecision = in.PricePrecision
	out.Price = in.Price
	out.Quantity = in.Quantity
	out.Side = in.Side
	out.TimeInForce = in.TimeInForce
	out.ExistBlocks = in.ExistBlocks
	return
} //End of DeepCopyMsgCreateOrder

// Non-Interface
func EncodeMsgCreateTradingPair(w *[]byte, v MsgCreateTradingPair) {
	codonEncodeString(1, w, v.Stock)
	codonEncodeString(2, w, v.Money)
	codonEncodeByteSlice(3, w, v.Creator[:])
	codonEncodeUint8(4, w, v.PricePrecision)
	codonEncodeUint8(5, w, v.OrderPrecision)
} //End of EncodeMsgCreateTradingPair

func DecodeMsgCreateTradingPair(bz []byte) (v MsgCreateTradingPair, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Stock
			v.Stock = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Money
			v.Money = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Creator
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Creator = tmpBz
		case 3: // v.PricePrecision
			v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.OrderPrecision
			v.OrderPrecision = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgCreateTradingPair

func RandMsgCreateTradingPair(r RandSrc) MsgCreateTradingPair {
	var length int
	var v MsgCreateTradingPair
	v.Stock = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Money = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Creator = r.GetBytes(length)
	v.PricePrecision = r.GetUint8()
	v.OrderPrecision = r.GetUint8()
	return v
} //End of RandMsgCreateTradingPair

func DeepCopyMsgCreateTradingPair(in MsgCreateTradingPair) (out MsgCreateTradingPair) {
	var length int
	out.Stock = in.Stock
	out.Money = in.Money
	length = len(in.Creator)
	if length == 0 {
		out.Creator = nil
	} else {
		out.Creator = make([]uint8, length)
	}
	copy(out.Creator[:], in.Creator[:])
	out.PricePrecision = in.PricePrecision
	out.OrderPrecision = in.OrderPrecision
	return
} //End of DeepCopyMsgCreateTradingPair

// Non-Interface
func EncodeMsgModifyPricePrecision(w *[]byte, v MsgModifyPricePrecision) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeString(2, w, v.TradingPair)
	codonEncodeUint8(3, w, v.PricePrecision)
} //End of EncodeMsgModifyPricePrecision

func DecodeMsgModifyPricePrecision(bz []byte) (v MsgModifyPricePrecision, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.TradingPair
			v.TradingPair = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.PricePrecision
			v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgModifyPricePrecision

func RandMsgModifyPricePrecision(r RandSrc) MsgModifyPricePrecision {
	var length int
	var v MsgModifyPricePrecision
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.TradingPair = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.PricePrecision = r.GetUint8()
	return v
} //End of RandMsgModifyPricePrecision

func DeepCopyMsgModifyPricePrecision(in MsgModifyPricePrecision) (out MsgModifyPricePrecision) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.TradingPair = in.TradingPair
	out.PricePrecision = in.PricePrecision
	return
} //End of DeepCopyMsgModifyPricePrecision

// Non-Interface
func EncodeOrder(w *[]byte, v Order) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeUvarint(2, w, uint64(v.Sequence))
	codonEncodeUint8(3, w, v.Identify)
	codonEncodeString(4, w, v.TradingPair)
	codonEncodeUint8(5, w, v.OrderType)
	codonEncodeByteSlice(6, w, EncodeDec(v.Price))
	codonEncodeVarint(7, w, int64(v.Quantity))
	codonEncodeUint8(8, w, v.Side)
	codonEncodeVarint(9, w, int64(v.TimeInForce))
	codonEncodeVarint(10, w, int64(v.Height))
	codonEncodeVarint(11, w, int64(v.FrozenFee))
	codonEncodeVarint(12, w, int64(v.ExistBlocks))
	codonEncodeVarint(13, w, int64(v.LeftStock))
	codonEncodeVarint(14, w, int64(v.Freeze))
	codonEncodeVarint(15, w, int64(v.DealStock))
	codonEncodeVarint(16, w, int64(v.DealMoney))
} //End of EncodeOrder

func DecodeOrder(bz []byte) (v Order, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.Sequence
			v.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Identify
			v.Identify = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.TradingPair
			v.TradingPair = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.OrderType
			v.OrderType = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 5: // v.Price
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.Price, n, err = DecodeDec(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 6: // v.Quantity
			v.Quantity = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 7: // v.Side
			v.Side = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 8: // v.TimeInForce
			v.TimeInForce = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 9: // v.Height
			v.Height = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 10: // v.FrozenFee
			v.FrozenFee = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 11: // v.ExistBlocks
			v.ExistBlocks = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 12: // v.LeftStock
			v.LeftStock = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 13: // v.Freeze
			v.Freeze = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 14: // v.DealStock
			v.DealStock = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 15: // v.DealMoney
			v.DealMoney = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeOrder

func RandOrder(r RandSrc) Order {
	var length int
	var v Order
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.Sequence = r.GetUint64()
	v.Identify = r.GetUint8()
	v.TradingPair = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.OrderType = r.GetUint8()
	v.Price = RandDec(r)
	v.Quantity = r.GetInt64()
	v.Side = r.GetUint8()
	v.TimeInForce = r.GetInt64()
	v.Height = r.GetInt64()
	v.FrozenFee = r.GetInt64()
	v.ExistBlocks = r.GetInt64()
	v.LeftStock = r.GetInt64()
	v.Freeze = r.GetInt64()
	v.DealStock = r.GetInt64()
	v.DealMoney = r.GetInt64()
	return v
} //End of RandOrder

func DeepCopyOrder(in Order) (out Order) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.Sequence = in.Sequence
	out.Identify = in.Identify
	out.TradingPair = in.TradingPair
	out.OrderType = in.OrderType
	out.Price = DeepCopyDec(in.Price)
	out.Quantity = in.Quantity
	out.Side = in.Side
	out.TimeInForce = in.TimeInForce
	out.Height = in.Height
	out.FrozenFee = in.FrozenFee
	out.ExistBlocks = in.ExistBlocks
	out.LeftStock = in.LeftStock
	out.Freeze = in.Freeze
	out.DealStock = in.DealStock
	out.DealMoney = in.DealMoney
	return
} //End of DeepCopyOrder

// Non-Interface
func EncodeMarketInfo(w *[]byte, v MarketInfo) {
	codonEncodeString(1, w, v.Stock)
	codonEncodeString(2, w, v.Money)
	codonEncodeUint8(3, w, v.PricePrecision)
	codonEncodeByteSlice(4, w, EncodeDec(v.LastExecutedPrice))
	codonEncodeUint8(5, w, v.OrderPrecision)
} //End of EncodeMarketInfo

func DecodeMarketInfo(bz []byte) (v MarketInfo, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Stock
			v.Stock = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Money
			v.Money = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.PricePrecision
			v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.LastExecutedPrice
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			v.LastExecutedPrice, n, err = DecodeDec(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
		case 4: // v.OrderPrecision
			v.OrderPrecision = uint8(codonDecodeUint8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMarketInfo

func RandMarketInfo(r RandSrc) MarketInfo {
	var v MarketInfo
	v.Stock = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Money = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.PricePrecision = r.GetUint8()
	v.LastExecutedPrice = RandDec(r)
	v.OrderPrecision = r.GetUint8()
	return v
} //End of RandMarketInfo

func DeepCopyMarketInfo(in MarketInfo) (out MarketInfo) {
	out.Stock = in.Stock
	out.Money = in.Money
	out.PricePrecision = in.PricePrecision
	out.LastExecutedPrice = DeepCopyDec(in.LastExecutedPrice)
	out.OrderPrecision = in.OrderPrecision
	return
} //End of DeepCopyMarketInfo

// Non-Interface
func EncodeMsgDonateToCommunityPool(w *[]byte, v MsgDonateToCommunityPool) {
	codonEncodeByteSlice(1, w, v.FromAddr[:])
	for _0 := 0; _0 < len(v.Amount); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.Amount[_0].Denom)
			codonEncodeByteSlice(2, w, EncodeInt(v.Amount[_0].Amount))
			return wBuf
		}()) // end of v.Amount[_0]
	}
} //End of EncodeMsgDonateToCommunityPool

func DecodeMsgDonateToCommunityPool(bz []byte) (v MsgDonateToCommunityPool, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.FromAddr
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.FromAddr = tmpBz
		case 1: // v.Amount
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp Coin
			tmp, n, err = DecodeCoin(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.Amount = append(v.Amount, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgDonateToCommunityPool

func RandMsgDonateToCommunityPool(r RandSrc) MsgDonateToCommunityPool {
	var length int
	var v MsgDonateToCommunityPool
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddr = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.Amount = nil
	} else {
		v.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandMsgDonateToCommunityPool

func DeepCopyMsgDonateToCommunityPool(in MsgDonateToCommunityPool) (out MsgDonateToCommunityPool) {
	var length int
	length = len(in.FromAddr)
	if length == 0 {
		out.FromAddr = nil
	} else {
		out.FromAddr = make([]uint8, length)
	}
	copy(out.FromAddr[:], in.FromAddr[:])
	length = len(in.Amount)
	if length == 0 {
		out.Amount = nil
	} else {
		out.Amount = make([]Coin, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.Amount[_0] = DeepCopyCoin(in.Amount[_0])
	}
	return
} //End of DeepCopyMsgDonateToCommunityPool

// Non-Interface
func EncodeMsgCommentToken(w *[]byte, v MsgCommentToken) {
	codonEncodeByteSlice(1, w, v.Sender[:])
	codonEncodeString(2, w, v.Token)
	codonEncodeVarint(3, w, int64(v.Donation))
	codonEncodeString(4, w, v.Title)
	codonEncodeByteSlice(5, w, v.Content[:])
	codonEncodeInt8(6, w, v.ContentType)
	for _0 := 0; _0 < len(v.References); _0++ {
		codonEncodeByteSlice(7, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeUvarint(1, w, uint64(v.References[_0].ID))
			codonEncodeByteSlice(2, w, v.References[_0].RewardTarget[:])
			codonEncodeString(3, w, v.References[_0].RewardToken)
			codonEncodeVarint(4, w, int64(v.References[_0].RewardAmount))
			for _1 := 0; _1 < len(v.References[_0].Attitudes); _1++ {
				codonEncodeVarint(5, w, int64(v.References[_0].Attitudes[_1]))
			}
			return wBuf
		}()) // end of v.References[_0]
	}
} //End of EncodeMsgCommentToken

func DecodeMsgCommentToken(bz []byte) (v MsgCommentToken, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Sender
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Sender = tmpBz
		case 1: // v.Token
			v.Token = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.Donation
			v.Donation = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.Title
			v.Title = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 4: // v.Content
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Content = tmpBz
		case 5: // v.ContentType
			v.ContentType = int8(codonDecodeInt8(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 6: // v.References
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp CommentRef
			tmp, n, err = DecodeCommentRef(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.References = append(v.References, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgCommentToken

func RandMsgCommentToken(r RandSrc) MsgCommentToken {
	var length int
	var v MsgCommentToken
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.Token = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Donation = r.GetInt64()
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Content = r.GetBytes(length)
	v.ContentType = r.GetInt8()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.References = nil
	} else {
		v.References = make([]CommentRef, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.References[_0] = RandCommentRef(r)
	}
	return v
} //End of RandMsgCommentToken

func DeepCopyMsgCommentToken(in MsgCommentToken) (out MsgCommentToken) {
	var length int
	length = len(in.Sender)
	if length == 0 {
		out.Sender = nil
	} else {
		out.Sender = make([]uint8, length)
	}
	copy(out.Sender[:], in.Sender[:])
	out.Token = in.Token
	out.Donation = in.Donation
	out.Title = in.Title
	length = len(in.Content)
	if length == 0 {
		out.Content = nil
	} else {
		out.Content = make([]uint8, length)
	}
	copy(out.Content[:], in.Content[:])
	out.ContentType = in.ContentType
	length = len(in.References)
	if length == 0 {
		out.References = nil
	} else {
		out.References = make([]CommentRef, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.References[_0] = DeepCopyCommentRef(in.References[_0])
	}
	return
} //End of DeepCopyMsgCommentToken

// Non-Interface
func EncodeState(w *[]byte, v State) {
	codonEncodeVarint(1, w, int64(v.HeightAdjustment))
} //End of EncodeState

func DecodeState(bz []byte) (v State, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.HeightAdjustment
			v.HeightAdjustment = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeState

func RandState(r RandSrc) State {
	var v State
	v.HeightAdjustment = r.GetInt64()
	return v
} //End of RandState

func DeepCopyState(in State) (out State) {
	out.HeightAdjustment = in.HeightAdjustment
	return
} //End of DeepCopyState

// Non-Interface
func EncodeMsgAliasUpdate(w *[]byte, v MsgAliasUpdate) {
	codonEncodeByteSlice(1, w, v.Owner[:])
	codonEncodeString(2, w, v.Alias)
	codonEncodeBool(3, w, v.IsAdd)
	codonEncodeBool(4, w, v.AsDefault)
} //End of EncodeMsgAliasUpdate

func DecodeMsgAliasUpdate(bz []byte) (v MsgAliasUpdate, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Owner
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			v.Owner = tmpBz
		case 1: // v.Alias
			v.Alias = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 2: // v.IsAdd
			v.IsAdd = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 3: // v.AsDefault
			v.AsDefault = bool(codonDecodeBool(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeMsgAliasUpdate

func RandMsgAliasUpdate(r RandSrc) MsgAliasUpdate {
	var length int
	var v MsgAliasUpdate
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Owner = r.GetBytes(length)
	v.Alias = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.IsAdd = r.GetBool()
	v.AsDefault = r.GetBool()
	return v
} //End of RandMsgAliasUpdate

func DeepCopyMsgAliasUpdate(in MsgAliasUpdate) (out MsgAliasUpdate) {
	var length int
	length = len(in.Owner)
	if length == 0 {
		out.Owner = nil
	} else {
		out.Owner = make([]uint8, length)
	}
	copy(out.Owner[:], in.Owner[:])
	out.Alias = in.Alias
	out.IsAdd = in.IsAdd
	out.AsDefault = in.AsDefault
	return
} //End of DeepCopyMsgAliasUpdate

// Non-Interface
func EncodeAccAddressList(w *[]byte, v AccAddressList) {
	for _0 := 0; _0 < len(v); _0++ {
		codonEncodeByteSlice(0, w, v[_0][:])
	}
} //End of EncodeAccAddressList

func DecodeAccAddressList(bz []byte) (v AccAddressList, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0:
			var tmp AccAddress
			var tmpBz []byte
			n, err = codonGetByteSlice(&tmpBz, bz)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			tmp = tmpBz
			v = append(v, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeAccAddressList

func RandAccAddressList(r RandSrc) AccAddressList {
	var length int
	var v AccAddressList
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v = nil
	} else {
		v = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v[_0] = r.GetBytes(length)
	}
	return v
} //End of RandAccAddressList

func DeepCopyAccAddressList(in AccAddressList) (out AccAddressList) {
	var length int
	length = len(in)
	if length == 0 {
		out = nil
	} else {
		out = make([]AccAddress, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = len(in[_0])
		if length == 0 {
			out[_0] = nil
		} else {
			out[_0] = make([]uint8, length)
		}
		copy(out[_0][:], in[_0][:])
	}
	return
} //End of DeepCopyAccAddressList

// Non-Interface
func EncodeCommitInfo(w *[]byte, v CommitInfo) {
	codonEncodeVarint(1, w, int64(v.Version))
	for _0 := 0; _0 < len(v.StoreInfos); _0++ {
		codonEncodeByteSlice(2, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeString(1, w, v.StoreInfos[_0].Name)
			codonEncodeByteSlice(2, w, func() []byte {
				wBuf := make([]byte, 0, 64)
				w := &wBuf
				codonEncodeByteSlice(1, w, func() []byte {
					wBuf := make([]byte, 0, 64)
					w := &wBuf
					codonEncodeVarint(1, w, int64(v.StoreInfos[_0].Core.CommitID.Version))
					codonEncodeByteSlice(2, w, v.StoreInfos[_0].Core.CommitID.Hash[:])
					return wBuf
				}()) // end of v.StoreInfos[_0].Core.CommitID
				return wBuf
			}()) // end of v.StoreInfos[_0].Core
			return wBuf
		}()) // end of v.StoreInfos[_0]
	}
} //End of EncodeCommitInfo

func DecodeCommitInfo(bz []byte) (v CommitInfo, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Version
			v.Version = int64(codonDecodeInt64(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.StoreInfos
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) > len(bz) {
				err = errors.New("Length Too Large")
				return
			}
			var tmp StoreInfo
			tmp, n, err = DecodeStoreInfo(bz[:l])
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			if int(l) != n {
				err = errors.New("Length Mismatch")
				return
			}
			v.StoreInfos = append(v.StoreInfos, tmp)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeCommitInfo

func RandCommitInfo(r RandSrc) CommitInfo {
	var length int
	var v CommitInfo
	v.Version = r.GetInt64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	if length == 0 {
		v.StoreInfos = nil
	} else {
		v.StoreInfos = make([]StoreInfo, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.StoreInfos[_0] = RandStoreInfo(r)
	}
	return v
} //End of RandCommitInfo

func DeepCopyCommitInfo(in CommitInfo) (out CommitInfo) {
	var length int
	out.Version = in.Version
	length = len(in.StoreInfos)
	if length == 0 {
		out.StoreInfos = nil
	} else {
		out.StoreInfos = make([]StoreInfo, length)
	}
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		out.StoreInfos[_0] = DeepCopyStoreInfo(in.StoreInfos[_0])
	}
	return
} //End of DeepCopyCommitInfo

// Non-Interface
func EncodeStoreInfo(w *[]byte, v StoreInfo) {
	codonEncodeString(1, w, v.Name)
	codonEncodeByteSlice(2, w, func() []byte {
		wBuf := make([]byte, 0, 64)
		w := &wBuf
		codonEncodeByteSlice(1, w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			codonEncodeVarint(1, w, int64(v.Core.CommitID.Version))
			codonEncodeByteSlice(2, w, v.Core.CommitID.Hash[:])
			return wBuf
		}()) // end of v.Core.CommitID
		return wBuf
	}()) // end of v.Core
} //End of EncodeStoreInfo

func DecodeStoreInfo(bz []byte) (v StoreInfo, total int, err error) {
	var n int
	for len(bz) != 0 {
		tag := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		tag = tag >> 3
		switch tag {
		case 0: // v.Name
			v.Name = string(codonDecodeString(bz, &n, &err))
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
		case 1: // v.Core
			l := codonDecodeUint64(bz, &n, &err)
			if err != nil {
				return
			}
			bz = bz[n:]
			total += n
			func(bz []byte) {
				for len(bz) != 0 {
					tag := codonDecodeUint64(bz, &n, &err)
					if err != nil {
						return
					}
					bz = bz[n:]
					total += n
					tag = tag >> 3
					switch tag {
					case 0: // v.Core.CommitID
						l := codonDecodeUint64(bz, &n, &err)
						if err != nil {
							return
						}
						bz = bz[n:]
						total += n
						func(bz []byte) {
							for len(bz) != 0 {
								tag := codonDecodeUint64(bz, &n, &err)
								if err != nil {
									return
								}
								bz = bz[n:]
								total += n
								tag = tag >> 3
								switch tag {
								case 0: // v.Core.CommitID.Version
									v.Core.CommitID.Version = int64(codonDecodeInt64(bz, &n, &err))
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
								case 1: // v.Core.CommitID.Hash
									var tmpBz []byte
									n, err = codonGetByteSlice(&tmpBz, bz)
									if err != nil {
										return
									}
									bz = bz[n:]
									total += n
									v.Core.CommitID.Hash = tmpBz
								default:
									err = errors.New("Unknown Field")
									return
								}
							} // end for
						}(bz[:l]) // end func
						if err != nil {
							return
						}
						bz = bz[l:]
						n += int(l)
					default:
						err = errors.New("Unknown Field")
						return
					}
				} // end for
			}(bz[:l]) // end func
			if err != nil {
				return
			}
			bz = bz[l:]
			n += int(l)
		default:
			err = errors.New("Unknown Field")
			return
		}
	} // end for
	return v, total, nil
} //End of DecodeStoreInfo

func RandStoreInfo(r RandSrc) StoreInfo {
	var length int
	var v StoreInfo
	v.Name = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Core.CommitID.Version = r.GetInt64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Core.CommitID.Hash = r.GetBytes(length)
	// end of v.Core.CommitID
	// end of v.Core
	return v
} //End of RandStoreInfo

func DeepCopyStoreInfo(in StoreInfo) (out StoreInfo) {
	var length int
	out.Name = in.Name
	out.Core.CommitID.Version = in.Core.CommitID.Version
	length = len(in.Core.CommitID.Hash)
	if length == 0 {
		out.Core.CommitID.Hash = nil
	} else {
		out.Core.CommitID.Hash = make([]uint8, length)
	}
	copy(out.Core.CommitID.Hash[:], in.Core.CommitID.Hash[:])
	// end of .Core.CommitID
	// end of .Core
	return
} //End of DeepCopyStoreInfo

// Interface
func DecodePubKey(bz []byte) (v PubKey, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 32929284:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PubKeyEd25519
		tmp, n, err = DecodePubKeyEd25519(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 471601077:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PubKeyMultisigThreshold
		tmp, n, err = DecodePubKeyMultisigThreshold(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 45700893:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PubKeySecp256k1
		tmp, n, err = DecodePubKeySecp256k1(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 496135554:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp StdSignature
		tmp, n, err = DecodeStdSignature(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodePubKey
func EncodePubKey(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case PubKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PubKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyEd25519(w, v)
			return wBuf
		}())
	case *PubKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PubKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyEd25519(w, *v)
			return wBuf
		}())
	case PubKeyMultisigThreshold:
		codonEncodeByteSlice(int(getMagicNum("PubKeyMultisigThreshold")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyMultisigThreshold(w, v)
			return wBuf
		}())
	case *PubKeyMultisigThreshold:
		codonEncodeByteSlice(int(getMagicNum("PubKeyMultisigThreshold")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyMultisigThreshold(w, *v)
			return wBuf
		}())
	case PubKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PubKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeySecp256k1(w, v)
			return wBuf
		}())
	case *PubKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PubKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeySecp256k1(w, *v)
			return wBuf
		}())
	case StdSignature:
		codonEncodeByteSlice(int(getMagicNum("StdSignature")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdSignature(w, v)
			return wBuf
		}())
	case *StdSignature:
		codonEncodeByteSlice(int(getMagicNum("StdSignature")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdSignature(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandPubKey(r RandSrc) PubKey {
	switch r.GetUint() % 2 {
	case 0:
		return RandPubKeyEd25519(r)
	case 1:
		return RandPubKeySecp256k1(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyPubKey(x PubKey) PubKey {
	switch v := x.(type) {
	case PubKeyEd25519:
		res := DeepCopyPubKeyEd25519(v)
		return res
	case *PubKeyEd25519:
		res := DeepCopyPubKeyEd25519(*v)
		return &res
	case PubKeySecp256k1:
		res := DeepCopyPubKeySecp256k1(v)
		return res
	case *PubKeySecp256k1:
		res := DeepCopyPubKeySecp256k1(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodePrivKey(bz []byte) (v PrivKey, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 454630091:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PrivKeyEd25519
		tmp, n, err = DecodePrivKeyEd25519(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 252074570:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PrivKeySecp256k1
		tmp, n, err = DecodePrivKeySecp256k1(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodePrivKey
func EncodePrivKey(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case PrivKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PrivKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeyEd25519(w, v)
			return wBuf
		}())
	case *PrivKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PrivKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeyEd25519(w, *v)
			return wBuf
		}())
	case PrivKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PrivKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeySecp256k1(w, v)
			return wBuf
		}())
	case *PrivKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PrivKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeySecp256k1(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandPrivKey(r RandSrc) PrivKey {
	switch r.GetUint() % 2 {
	case 0:
		return RandPrivKeyEd25519(r)
	case 1:
		return RandPrivKeySecp256k1(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyPrivKey(x PrivKey) PrivKey {
	switch v := x.(type) {
	case PrivKeyEd25519:
		res := DeepCopyPrivKeyEd25519(v)
		return res
	case *PrivKeyEd25519:
		res := DeepCopyPrivKeyEd25519(*v)
		return &res
	case PrivKeySecp256k1:
		res := DeepCopyPrivKeySecp256k1(v)
		return res
	case *PrivKeySecp256k1:
		res := DeepCopyPrivKeySecp256k1(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeMsg(bz []byte) (v Msg, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 153655625:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgAddTokenWhitelist
		tmp, n, err = DecodeMsgAddTokenWhitelist(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 285022960:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgAliasUpdate
		tmp, n, err = DecodeMsgAliasUpdate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 495620655:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBancorCancel
		tmp, n, err = DecodeMsgBancorCancel(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 115719891:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBancorInit
		tmp, n, err = DecodeMsgBancorInit(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 158991245:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBancorTrade
		tmp, n, err = DecodeMsgBancorTrade(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 241505485:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBeginRedelegate
		tmp, n, err = DecodeMsgBeginRedelegate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 449833530:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBurnToken
		tmp, n, err = DecodeMsgBurnToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 304885884:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCancelOrder
		tmp, n, err = DecodeMsgCancelOrder(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 332906025:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCancelTradingPair
		tmp, n, err = DecodeMsgCancelTradingPair(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 56755914:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCommentToken
		tmp, n, err = DecodeMsgCommentToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 377235122:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCreateOrder
		tmp, n, err = DecodeMsgCreateOrder(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 484769215:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCreateTradingPair
		tmp, n, err = DecodeMsgCreateTradingPair(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 374726888:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCreateValidator
		tmp, n, err = DecodeMsgCreateValidator(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 106121128:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgDelegate
		tmp, n, err = DecodeMsgDelegate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 402428987:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgDeposit
		tmp, n, err = DecodeMsgDeposit(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 138691582:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgDonateToCommunityPool
		tmp, n, err = DecodeMsgDonateToCommunityPool(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 502447547:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgEditValidator
		tmp, n, err = DecodeMsgEditValidator(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 488219826:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgForbidAddr
		tmp, n, err = DecodeMsgForbidAddr(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 286625767:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgForbidToken
		tmp, n, err = DecodeMsgForbidToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 336682718:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgIssueToken
		tmp, n, err = DecodeMsgIssueToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 135231766:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgMintToken
		tmp, n, err = DecodeMsgMintToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 380178367:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgModifyPricePrecision
		tmp, n, err = DecodeMsgModifyPricePrecision(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 405536984:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgModifyTokenInfo
		tmp, n, err = DecodeMsgModifyTokenInfo(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 167957709:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgMultiSend
		tmp, n, err = DecodeMsgMultiSend(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 155708190:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgMultiSendX
		tmp, n, err = DecodeMsgMultiSendX(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 94212241:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgRemoveTokenWhitelist
		tmp, n, err = DecodeMsgRemoveTokenWhitelist(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 31944689:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSend
		tmp, n, err = DecodeMsgSend(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 145878995:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSendX
		tmp, n, err = DecodeMsgSendX(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 376054118:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSetMemoRequired
		tmp, n, err = DecodeMsgSetMemoRequired(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 73614168:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSetWithdrawAddress
		tmp, n, err = DecodeMsgSetWithdrawAddress(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 331897574:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSubmitProposal
		tmp, n, err = DecodeMsgSubmitProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 87007006:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSupervisedSend
		tmp, n, err = DecodeMsgSupervisedSend(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 377250605:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgTransferOwnership
		tmp, n, err = DecodeMsgTransferOwnership(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 92120341:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUnForbidAddr
		tmp, n, err = DecodeMsgUnForbidAddr(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 355758274:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUnForbidToken
		tmp, n, err = DecodeMsgUnForbidToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 479724980:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUndelegate
		tmp, n, err = DecodeMsgUndelegate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 261008494:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUnjail
		tmp, n, err = DecodeMsgUnjail(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 121564406:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgVerifyInvariant
		tmp, n, err = DecodeMsgVerifyInvariant(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 314113032:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgVote
		tmp, n, err = DecodeMsgVote(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 23904706:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgWithdrawDelegatorReward
		tmp, n, err = DecodeMsgWithdrawDelegatorReward(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 514264249:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgWithdrawValidatorCommission
		tmp, n, err = DecodeMsgWithdrawValidatorCommission(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeMsg
func EncodeMsg(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case MsgAddTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgAddTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAddTokenWhitelist(w, v)
			return wBuf
		}())
	case *MsgAddTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgAddTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAddTokenWhitelist(w, *v)
			return wBuf
		}())
	case MsgAliasUpdate:
		codonEncodeByteSlice(int(getMagicNum("MsgAliasUpdate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAliasUpdate(w, v)
			return wBuf
		}())
	case *MsgAliasUpdate:
		codonEncodeByteSlice(int(getMagicNum("MsgAliasUpdate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAliasUpdate(w, *v)
			return wBuf
		}())
	case MsgBancorCancel:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorCancel")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorCancel(w, v)
			return wBuf
		}())
	case *MsgBancorCancel:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorCancel")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorCancel(w, *v)
			return wBuf
		}())
	case MsgBancorInit:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorInit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorInit(w, v)
			return wBuf
		}())
	case *MsgBancorInit:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorInit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorInit(w, *v)
			return wBuf
		}())
	case MsgBancorTrade:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorTrade")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorTrade(w, v)
			return wBuf
		}())
	case *MsgBancorTrade:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorTrade")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorTrade(w, *v)
			return wBuf
		}())
	case MsgBeginRedelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgBeginRedelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBeginRedelegate(w, v)
			return wBuf
		}())
	case *MsgBeginRedelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgBeginRedelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBeginRedelegate(w, *v)
			return wBuf
		}())
	case MsgBurnToken:
		codonEncodeByteSlice(int(getMagicNum("MsgBurnToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBurnToken(w, v)
			return wBuf
		}())
	case *MsgBurnToken:
		codonEncodeByteSlice(int(getMagicNum("MsgBurnToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBurnToken(w, *v)
			return wBuf
		}())
	case MsgCancelOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelOrder(w, v)
			return wBuf
		}())
	case *MsgCancelOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelOrder(w, *v)
			return wBuf
		}())
	case MsgCancelTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelTradingPair(w, v)
			return wBuf
		}())
	case *MsgCancelTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelTradingPair(w, *v)
			return wBuf
		}())
	case MsgCommentToken:
		codonEncodeByteSlice(int(getMagicNum("MsgCommentToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCommentToken(w, v)
			return wBuf
		}())
	case *MsgCommentToken:
		codonEncodeByteSlice(int(getMagicNum("MsgCommentToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCommentToken(w, *v)
			return wBuf
		}())
	case MsgCreateOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateOrder(w, v)
			return wBuf
		}())
	case *MsgCreateOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateOrder(w, *v)
			return wBuf
		}())
	case MsgCreateTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateTradingPair(w, v)
			return wBuf
		}())
	case *MsgCreateTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateTradingPair(w, *v)
			return wBuf
		}())
	case MsgCreateValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateValidator(w, v)
			return wBuf
		}())
	case *MsgCreateValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateValidator(w, *v)
			return wBuf
		}())
	case MsgDelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgDelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDelegate(w, v)
			return wBuf
		}())
	case *MsgDelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgDelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDelegate(w, *v)
			return wBuf
		}())
	case MsgDeposit:
		codonEncodeByteSlice(int(getMagicNum("MsgDeposit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDeposit(w, v)
			return wBuf
		}())
	case *MsgDeposit:
		codonEncodeByteSlice(int(getMagicNum("MsgDeposit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDeposit(w, *v)
			return wBuf
		}())
	case MsgDonateToCommunityPool:
		codonEncodeByteSlice(int(getMagicNum("MsgDonateToCommunityPool")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDonateToCommunityPool(w, v)
			return wBuf
		}())
	case *MsgDonateToCommunityPool:
		codonEncodeByteSlice(int(getMagicNum("MsgDonateToCommunityPool")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDonateToCommunityPool(w, *v)
			return wBuf
		}())
	case MsgEditValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgEditValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgEditValidator(w, v)
			return wBuf
		}())
	case *MsgEditValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgEditValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgEditValidator(w, *v)
			return wBuf
		}())
	case MsgForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidAddr(w, v)
			return wBuf
		}())
	case *MsgForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidAddr(w, *v)
			return wBuf
		}())
	case MsgForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidToken(w, v)
			return wBuf
		}())
	case *MsgForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidToken(w, *v)
			return wBuf
		}())
	case MsgIssueToken:
		codonEncodeByteSlice(int(getMagicNum("MsgIssueToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgIssueToken(w, v)
			return wBuf
		}())
	case *MsgIssueToken:
		codonEncodeByteSlice(int(getMagicNum("MsgIssueToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgIssueToken(w, *v)
			return wBuf
		}())
	case MsgMintToken:
		codonEncodeByteSlice(int(getMagicNum("MsgMintToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMintToken(w, v)
			return wBuf
		}())
	case *MsgMintToken:
		codonEncodeByteSlice(int(getMagicNum("MsgMintToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMintToken(w, *v)
			return wBuf
		}())
	case MsgModifyPricePrecision:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyPricePrecision")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyPricePrecision(w, v)
			return wBuf
		}())
	case *MsgModifyPricePrecision:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyPricePrecision")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyPricePrecision(w, *v)
			return wBuf
		}())
	case MsgModifyTokenInfo:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyTokenInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyTokenInfo(w, v)
			return wBuf
		}())
	case *MsgModifyTokenInfo:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyTokenInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyTokenInfo(w, *v)
			return wBuf
		}())
	case MsgMultiSend:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSend(w, v)
			return wBuf
		}())
	case *MsgMultiSend:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSend(w, *v)
			return wBuf
		}())
	case MsgMultiSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSendX(w, v)
			return wBuf
		}())
	case *MsgMultiSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSendX(w, *v)
			return wBuf
		}())
	case MsgRemoveTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgRemoveTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgRemoveTokenWhitelist(w, v)
			return wBuf
		}())
	case *MsgRemoveTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgRemoveTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgRemoveTokenWhitelist(w, *v)
			return wBuf
		}())
	case MsgSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSend(w, v)
			return wBuf
		}())
	case *MsgSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSend(w, *v)
			return wBuf
		}())
	case MsgSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSendX(w, v)
			return wBuf
		}())
	case *MsgSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSendX(w, *v)
			return wBuf
		}())
	case MsgSetMemoRequired:
		codonEncodeByteSlice(int(getMagicNum("MsgSetMemoRequired")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetMemoRequired(w, v)
			return wBuf
		}())
	case *MsgSetMemoRequired:
		codonEncodeByteSlice(int(getMagicNum("MsgSetMemoRequired")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetMemoRequired(w, *v)
			return wBuf
		}())
	case MsgSetWithdrawAddress:
		codonEncodeByteSlice(int(getMagicNum("MsgSetWithdrawAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetWithdrawAddress(w, v)
			return wBuf
		}())
	case *MsgSetWithdrawAddress:
		codonEncodeByteSlice(int(getMagicNum("MsgSetWithdrawAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetWithdrawAddress(w, *v)
			return wBuf
		}())
	case MsgSubmitProposal:
		codonEncodeByteSlice(int(getMagicNum("MsgSubmitProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSubmitProposal(w, v)
			return wBuf
		}())
	case *MsgSubmitProposal:
		codonEncodeByteSlice(int(getMagicNum("MsgSubmitProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSubmitProposal(w, *v)
			return wBuf
		}())
	case MsgSupervisedSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSupervisedSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSupervisedSend(w, v)
			return wBuf
		}())
	case *MsgSupervisedSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSupervisedSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSupervisedSend(w, *v)
			return wBuf
		}())
	case MsgTransferOwnership:
		codonEncodeByteSlice(int(getMagicNum("MsgTransferOwnership")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgTransferOwnership(w, v)
			return wBuf
		}())
	case *MsgTransferOwnership:
		codonEncodeByteSlice(int(getMagicNum("MsgTransferOwnership")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgTransferOwnership(w, *v)
			return wBuf
		}())
	case MsgUnForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidAddr(w, v)
			return wBuf
		}())
	case *MsgUnForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidAddr(w, *v)
			return wBuf
		}())
	case MsgUnForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidToken(w, v)
			return wBuf
		}())
	case *MsgUnForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidToken(w, *v)
			return wBuf
		}())
	case MsgUndelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgUndelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUndelegate(w, v)
			return wBuf
		}())
	case *MsgUndelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgUndelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUndelegate(w, *v)
			return wBuf
		}())
	case MsgUnjail:
		codonEncodeByteSlice(int(getMagicNum("MsgUnjail")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnjail(w, v)
			return wBuf
		}())
	case *MsgUnjail:
		codonEncodeByteSlice(int(getMagicNum("MsgUnjail")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnjail(w, *v)
			return wBuf
		}())
	case MsgVerifyInvariant:
		codonEncodeByteSlice(int(getMagicNum("MsgVerifyInvariant")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVerifyInvariant(w, v)
			return wBuf
		}())
	case *MsgVerifyInvariant:
		codonEncodeByteSlice(int(getMagicNum("MsgVerifyInvariant")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVerifyInvariant(w, *v)
			return wBuf
		}())
	case MsgVote:
		codonEncodeByteSlice(int(getMagicNum("MsgVote")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVote(w, v)
			return wBuf
		}())
	case *MsgVote:
		codonEncodeByteSlice(int(getMagicNum("MsgVote")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVote(w, *v)
			return wBuf
		}())
	case MsgWithdrawDelegatorReward:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawDelegatorReward")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawDelegatorReward(w, v)
			return wBuf
		}())
	case *MsgWithdrawDelegatorReward:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawDelegatorReward")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawDelegatorReward(w, *v)
			return wBuf
		}())
	case MsgWithdrawValidatorCommission:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawValidatorCommission")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawValidatorCommission(w, v)
			return wBuf
		}())
	case *MsgWithdrawValidatorCommission:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawValidatorCommission")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawValidatorCommission(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandMsg(r RandSrc) Msg {
	switch r.GetUint() % 41 {
	case 0:
		return RandMsgAddTokenWhitelist(r)
	case 1:
		return RandMsgAliasUpdate(r)
	case 2:
		return RandMsgBancorCancel(r)
	case 3:
		return RandMsgBancorInit(r)
	case 4:
		return RandMsgBancorTrade(r)
	case 5:
		return RandMsgBeginRedelegate(r)
	case 6:
		return RandMsgBurnToken(r)
	case 7:
		return RandMsgCancelOrder(r)
	case 8:
		return RandMsgCancelTradingPair(r)
	case 9:
		return RandMsgCommentToken(r)
	case 10:
		return RandMsgCreateOrder(r)
	case 11:
		return RandMsgCreateTradingPair(r)
	case 12:
		return RandMsgCreateValidator(r)
	case 13:
		return RandMsgDelegate(r)
	case 14:
		return RandMsgDeposit(r)
	case 15:
		return RandMsgDonateToCommunityPool(r)
	case 16:
		return RandMsgEditValidator(r)
	case 17:
		return RandMsgForbidAddr(r)
	case 18:
		return RandMsgForbidToken(r)
	case 19:
		return RandMsgIssueToken(r)
	case 20:
		return RandMsgMintToken(r)
	case 21:
		return RandMsgModifyPricePrecision(r)
	case 22:
		return RandMsgModifyTokenInfo(r)
	case 23:
		return RandMsgMultiSend(r)
	case 24:
		return RandMsgMultiSendX(r)
	case 25:
		return RandMsgRemoveTokenWhitelist(r)
	case 26:
		return RandMsgSend(r)
	case 27:
		return RandMsgSendX(r)
	case 28:
		return RandMsgSetMemoRequired(r)
	case 29:
		return RandMsgSetWithdrawAddress(r)
	case 30:
		return RandMsgSubmitProposal(r)
	case 31:
		return RandMsgSupervisedSend(r)
	case 32:
		return RandMsgTransferOwnership(r)
	case 33:
		return RandMsgUnForbidAddr(r)
	case 34:
		return RandMsgUnForbidToken(r)
	case 35:
		return RandMsgUndelegate(r)
	case 36:
		return RandMsgUnjail(r)
	case 37:
		return RandMsgVerifyInvariant(r)
	case 38:
		return RandMsgVote(r)
	case 39:
		return RandMsgWithdrawDelegatorReward(r)
	case 40:
		return RandMsgWithdrawValidatorCommission(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyMsg(x Msg) Msg {
	switch v := x.(type) {
	case MsgAddTokenWhitelist:
		res := DeepCopyMsgAddTokenWhitelist(v)
		return res
	case *MsgAddTokenWhitelist:
		res := DeepCopyMsgAddTokenWhitelist(*v)
		return &res
	case MsgAliasUpdate:
		res := DeepCopyMsgAliasUpdate(v)
		return res
	case *MsgAliasUpdate:
		res := DeepCopyMsgAliasUpdate(*v)
		return &res
	case MsgBancorCancel:
		res := DeepCopyMsgBancorCancel(v)
		return res
	case *MsgBancorCancel:
		res := DeepCopyMsgBancorCancel(*v)
		return &res
	case MsgBancorInit:
		res := DeepCopyMsgBancorInit(v)
		return res
	case *MsgBancorInit:
		res := DeepCopyMsgBancorInit(*v)
		return &res
	case MsgBancorTrade:
		res := DeepCopyMsgBancorTrade(v)
		return res
	case *MsgBancorTrade:
		res := DeepCopyMsgBancorTrade(*v)
		return &res
	case MsgBeginRedelegate:
		res := DeepCopyMsgBeginRedelegate(v)
		return res
	case *MsgBeginRedelegate:
		res := DeepCopyMsgBeginRedelegate(*v)
		return &res
	case MsgBurnToken:
		res := DeepCopyMsgBurnToken(v)
		return res
	case *MsgBurnToken:
		res := DeepCopyMsgBurnToken(*v)
		return &res
	case MsgCancelOrder:
		res := DeepCopyMsgCancelOrder(v)
		return res
	case *MsgCancelOrder:
		res := DeepCopyMsgCancelOrder(*v)
		return &res
	case MsgCancelTradingPair:
		res := DeepCopyMsgCancelTradingPair(v)
		return res
	case *MsgCancelTradingPair:
		res := DeepCopyMsgCancelTradingPair(*v)
		return &res
	case MsgCommentToken:
		res := DeepCopyMsgCommentToken(v)
		return res
	case *MsgCommentToken:
		res := DeepCopyMsgCommentToken(*v)
		return &res
	case MsgCreateOrder:
		res := DeepCopyMsgCreateOrder(v)
		return res
	case *MsgCreateOrder:
		res := DeepCopyMsgCreateOrder(*v)
		return &res
	case MsgCreateTradingPair:
		res := DeepCopyMsgCreateTradingPair(v)
		return res
	case *MsgCreateTradingPair:
		res := DeepCopyMsgCreateTradingPair(*v)
		return &res
	case MsgCreateValidator:
		res := DeepCopyMsgCreateValidator(v)
		return res
	case *MsgCreateValidator:
		res := DeepCopyMsgCreateValidator(*v)
		return &res
	case MsgDelegate:
		res := DeepCopyMsgDelegate(v)
		return res
	case *MsgDelegate:
		res := DeepCopyMsgDelegate(*v)
		return &res
	case MsgDeposit:
		res := DeepCopyMsgDeposit(v)
		return res
	case *MsgDeposit:
		res := DeepCopyMsgDeposit(*v)
		return &res
	case MsgDonateToCommunityPool:
		res := DeepCopyMsgDonateToCommunityPool(v)
		return res
	case *MsgDonateToCommunityPool:
		res := DeepCopyMsgDonateToCommunityPool(*v)
		return &res
	case MsgEditValidator:
		res := DeepCopyMsgEditValidator(v)
		return res
	case *MsgEditValidator:
		res := DeepCopyMsgEditValidator(*v)
		return &res
	case MsgForbidAddr:
		res := DeepCopyMsgForbidAddr(v)
		return res
	case *MsgForbidAddr:
		res := DeepCopyMsgForbidAddr(*v)
		return &res
	case MsgForbidToken:
		res := DeepCopyMsgForbidToken(v)
		return res
	case *MsgForbidToken:
		res := DeepCopyMsgForbidToken(*v)
		return &res
	case MsgIssueToken:
		res := DeepCopyMsgIssueToken(v)
		return res
	case *MsgIssueToken:
		res := DeepCopyMsgIssueToken(*v)
		return &res
	case MsgMintToken:
		res := DeepCopyMsgMintToken(v)
		return res
	case *MsgMintToken:
		res := DeepCopyMsgMintToken(*v)
		return &res
	case MsgModifyPricePrecision:
		res := DeepCopyMsgModifyPricePrecision(v)
		return res
	case *MsgModifyPricePrecision:
		res := DeepCopyMsgModifyPricePrecision(*v)
		return &res
	case MsgModifyTokenInfo:
		res := DeepCopyMsgModifyTokenInfo(v)
		return res
	case *MsgModifyTokenInfo:
		res := DeepCopyMsgModifyTokenInfo(*v)
		return &res
	case MsgMultiSend:
		res := DeepCopyMsgMultiSend(v)
		return res
	case *MsgMultiSend:
		res := DeepCopyMsgMultiSend(*v)
		return &res
	case MsgMultiSendX:
		res := DeepCopyMsgMultiSendX(v)
		return res
	case *MsgMultiSendX:
		res := DeepCopyMsgMultiSendX(*v)
		return &res
	case MsgRemoveTokenWhitelist:
		res := DeepCopyMsgRemoveTokenWhitelist(v)
		return res
	case *MsgRemoveTokenWhitelist:
		res := DeepCopyMsgRemoveTokenWhitelist(*v)
		return &res
	case MsgSend:
		res := DeepCopyMsgSend(v)
		return res
	case *MsgSend:
		res := DeepCopyMsgSend(*v)
		return &res
	case MsgSendX:
		res := DeepCopyMsgSendX(v)
		return res
	case *MsgSendX:
		res := DeepCopyMsgSendX(*v)
		return &res
	case MsgSetMemoRequired:
		res := DeepCopyMsgSetMemoRequired(v)
		return res
	case *MsgSetMemoRequired:
		res := DeepCopyMsgSetMemoRequired(*v)
		return &res
	case MsgSetWithdrawAddress:
		res := DeepCopyMsgSetWithdrawAddress(v)
		return res
	case *MsgSetWithdrawAddress:
		res := DeepCopyMsgSetWithdrawAddress(*v)
		return &res
	case MsgSubmitProposal:
		res := DeepCopyMsgSubmitProposal(v)
		return res
	case *MsgSubmitProposal:
		res := DeepCopyMsgSubmitProposal(*v)
		return &res
	case MsgSupervisedSend:
		res := DeepCopyMsgSupervisedSend(v)
		return res
	case *MsgSupervisedSend:
		res := DeepCopyMsgSupervisedSend(*v)
		return &res
	case MsgTransferOwnership:
		res := DeepCopyMsgTransferOwnership(v)
		return res
	case *MsgTransferOwnership:
		res := DeepCopyMsgTransferOwnership(*v)
		return &res
	case MsgUnForbidAddr:
		res := DeepCopyMsgUnForbidAddr(v)
		return res
	case *MsgUnForbidAddr:
		res := DeepCopyMsgUnForbidAddr(*v)
		return &res
	case MsgUnForbidToken:
		res := DeepCopyMsgUnForbidToken(v)
		return res
	case *MsgUnForbidToken:
		res := DeepCopyMsgUnForbidToken(*v)
		return &res
	case MsgUndelegate:
		res := DeepCopyMsgUndelegate(v)
		return res
	case *MsgUndelegate:
		res := DeepCopyMsgUndelegate(*v)
		return &res
	case MsgUnjail:
		res := DeepCopyMsgUnjail(v)
		return res
	case *MsgUnjail:
		res := DeepCopyMsgUnjail(*v)
		return &res
	case MsgVerifyInvariant:
		res := DeepCopyMsgVerifyInvariant(v)
		return res
	case *MsgVerifyInvariant:
		res := DeepCopyMsgVerifyInvariant(*v)
		return &res
	case MsgVote:
		res := DeepCopyMsgVote(v)
		return res
	case *MsgVote:
		res := DeepCopyMsgVote(*v)
		return &res
	case MsgWithdrawDelegatorReward:
		res := DeepCopyMsgWithdrawDelegatorReward(v)
		return res
	case *MsgWithdrawDelegatorReward:
		res := DeepCopyMsgWithdrawDelegatorReward(*v)
		return &res
	case MsgWithdrawValidatorCommission:
		res := DeepCopyMsgWithdrawValidatorCommission(v)
		return res
	case *MsgWithdrawValidatorCommission:
		res := DeepCopyMsgWithdrawValidatorCommission(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeAccount(bz []byte) (v Account, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 26591049:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp BaseAccount
		tmp, n, err = DecodeBaseAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = &tmp
		return
	case 351017031:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp BaseVestingAccount
		tmp, n, err = DecodeBaseVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 82892118:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ContinuousVestingAccount
		tmp, n, err = DecodeContinuousVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 12374068:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp DelayedVestingAccount
		tmp, n, err = DecodeDelayedVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 321193714:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ModuleAccount
		tmp, n, err = DecodeModuleAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeAccount
func EncodeAccount(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case BaseAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseAccount(w, v)
			return wBuf
		}())
	case *BaseAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseAccount(w, *v)
			return wBuf
		}())
	case BaseVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseVestingAccount(w, v)
			return wBuf
		}())
	case *BaseVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseVestingAccount(w, *v)
			return wBuf
		}())
	case ContinuousVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("ContinuousVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeContinuousVestingAccount(w, v)
			return wBuf
		}())
	case *ContinuousVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("ContinuousVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeContinuousVestingAccount(w, *v)
			return wBuf
		}())
	case DelayedVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("DelayedVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDelayedVestingAccount(w, v)
			return wBuf
		}())
	case *DelayedVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("DelayedVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDelayedVestingAccount(w, *v)
			return wBuf
		}())
	case ModuleAccount:
		codonEncodeByteSlice(int(getMagicNum("ModuleAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeModuleAccount(w, v)
			return wBuf
		}())
	case *ModuleAccount:
		codonEncodeByteSlice(int(getMagicNum("ModuleAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeModuleAccount(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandAccount(r RandSrc) Account {
	switch r.GetUint() % 5 {
	case 0:
		tmp := RandBaseAccount(r)
		return &tmp
	case 1:
		return RandBaseVestingAccount(r)
	case 2:
		return RandContinuousVestingAccount(r)
	case 3:
		return RandDelayedVestingAccount(r)
	case 4:
		return RandModuleAccount(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyAccount(x Account) Account {
	switch v := x.(type) {
	case *BaseAccount:
		res := DeepCopyBaseAccount(*v)
		return &res
	case BaseVestingAccount:
		res := DeepCopyBaseVestingAccount(v)
		return res
	case *BaseVestingAccount:
		res := DeepCopyBaseVestingAccount(*v)
		return &res
	case ContinuousVestingAccount:
		res := DeepCopyContinuousVestingAccount(v)
		return res
	case *ContinuousVestingAccount:
		res := DeepCopyContinuousVestingAccount(*v)
		return &res
	case DelayedVestingAccount:
		res := DeepCopyDelayedVestingAccount(v)
		return res
	case *DelayedVestingAccount:
		res := DeepCopyDelayedVestingAccount(*v)
		return &res
	case ModuleAccount:
		res := DeepCopyModuleAccount(v)
		return res
	case *ModuleAccount:
		res := DeepCopyModuleAccount(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeVestingAccount(bz []byte) (v VestingAccount, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 82892118:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ContinuousVestingAccount
		tmp, n, err = DecodeContinuousVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = &tmp
		return
	case 12374068:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp DelayedVestingAccount
		tmp, n, err = DecodeDelayedVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = &tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeVestingAccount
func EncodeVestingAccount(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case ContinuousVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("ContinuousVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeContinuousVestingAccount(w, v)
			return wBuf
		}())
	case *ContinuousVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("ContinuousVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeContinuousVestingAccount(w, *v)
			return wBuf
		}())
	case DelayedVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("DelayedVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDelayedVestingAccount(w, v)
			return wBuf
		}())
	case *DelayedVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("DelayedVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDelayedVestingAccount(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandVestingAccount(r RandSrc) VestingAccount {
	switch r.GetUint() % 2 {
	case 0:
		tmp := RandContinuousVestingAccount(r)
		return &tmp
	case 1:
		tmp := RandDelayedVestingAccount(r)
		return &tmp
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyVestingAccount(x VestingAccount) VestingAccount {
	switch v := x.(type) {
	case *ContinuousVestingAccount:
		res := DeepCopyContinuousVestingAccount(*v)
		return &res
	case *DelayedVestingAccount:
		res := DeepCopyDelayedVestingAccount(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeContent(bz []byte) (v Content, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 23847857:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp CommunityPoolSpendProposal
		tmp, n, err = DecodeCommunityPoolSpendProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 84061114:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ParameterChangeProposal
		tmp, n, err = DecodeParameterChangeProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 261665006:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp SoftwareUpgradeProposal
		tmp, n, err = DecodeSoftwareUpgradeProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 458846397:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp TextProposal
		tmp, n, err = DecodeTextProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeContent
func EncodeContent(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case CommunityPoolSpendProposal:
		codonEncodeByteSlice(int(getMagicNum("CommunityPoolSpendProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommunityPoolSpendProposal(w, v)
			return wBuf
		}())
	case *CommunityPoolSpendProposal:
		codonEncodeByteSlice(int(getMagicNum("CommunityPoolSpendProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommunityPoolSpendProposal(w, *v)
			return wBuf
		}())
	case ParameterChangeProposal:
		codonEncodeByteSlice(int(getMagicNum("ParameterChangeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeParameterChangeProposal(w, v)
			return wBuf
		}())
	case *ParameterChangeProposal:
		codonEncodeByteSlice(int(getMagicNum("ParameterChangeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeParameterChangeProposal(w, *v)
			return wBuf
		}())
	case SoftwareUpgradeProposal:
		codonEncodeByteSlice(int(getMagicNum("SoftwareUpgradeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSoftwareUpgradeProposal(w, v)
			return wBuf
		}())
	case *SoftwareUpgradeProposal:
		codonEncodeByteSlice(int(getMagicNum("SoftwareUpgradeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSoftwareUpgradeProposal(w, *v)
			return wBuf
		}())
	case TextProposal:
		codonEncodeByteSlice(int(getMagicNum("TextProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeTextProposal(w, v)
			return wBuf
		}())
	case *TextProposal:
		codonEncodeByteSlice(int(getMagicNum("TextProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeTextProposal(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandContent(r RandSrc) Content {
	switch r.GetUint() % 4 {
	case 0:
		return RandCommunityPoolSpendProposal(r)
	case 1:
		return RandParameterChangeProposal(r)
	case 2:
		return RandSoftwareUpgradeProposal(r)
	case 3:
		return RandTextProposal(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyContent(x Content) Content {
	switch v := x.(type) {
	case CommunityPoolSpendProposal:
		res := DeepCopyCommunityPoolSpendProposal(v)
		return res
	case *CommunityPoolSpendProposal:
		res := DeepCopyCommunityPoolSpendProposal(*v)
		return &res
	case ParameterChangeProposal:
		res := DeepCopyParameterChangeProposal(v)
		return res
	case *ParameterChangeProposal:
		res := DeepCopyParameterChangeProposal(*v)
		return &res
	case SoftwareUpgradeProposal:
		res := DeepCopySoftwareUpgradeProposal(v)
		return res
	case *SoftwareUpgradeProposal:
		res := DeepCopySoftwareUpgradeProposal(*v)
		return &res
	case TextProposal:
		res := DeepCopyTextProposal(v)
		return res
	case *TextProposal:
		res := DeepCopyTextProposal(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeTx(bz []byte) (v Tx, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 420461272:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp StdTx
		tmp, n, err = DecodeStdTx(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeTx
func EncodeTx(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case StdTx:
		codonEncodeByteSlice(int(getMagicNum("StdTx")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdTx(w, v)
			return wBuf
		}())
	case *StdTx:
		codonEncodeByteSlice(int(getMagicNum("StdTx")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdTx(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandTx(r RandSrc) Tx {
	switch r.GetUint() % 1 {
	case 0:
		return RandStdTx(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyTx(x Tx) Tx {
	switch v := x.(type) {
	case StdTx:
		res := DeepCopyStdTx(v)
		return res
	case *StdTx:
		res := DeepCopyStdTx(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeModuleAccountI(bz []byte) (v ModuleAccountI, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 321193714:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ModuleAccount
		tmp, n, err = DecodeModuleAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeModuleAccountI
func EncodeModuleAccountI(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case ModuleAccount:
		codonEncodeByteSlice(int(getMagicNum("ModuleAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeModuleAccount(w, v)
			return wBuf
		}())
	case *ModuleAccount:
		codonEncodeByteSlice(int(getMagicNum("ModuleAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeModuleAccount(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandModuleAccountI(r RandSrc) ModuleAccountI {
	switch r.GetUint() % 1 {
	case 0:
		return RandModuleAccount(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyModuleAccountI(x ModuleAccountI) ModuleAccountI {
	switch v := x.(type) {
	case ModuleAccount:
		res := DeepCopyModuleAccount(v)
		return res
	case *ModuleAccount:
		res := DeepCopyModuleAccount(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeSupplyI(bz []byte) (v SupplyI, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 288082478:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp Supply
		tmp, n, err = DecodeSupply(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeSupplyI
func EncodeSupplyI(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case Supply:
		codonEncodeByteSlice(int(getMagicNum("Supply")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSupply(w, v)
			return wBuf
		}())
	case *Supply:
		codonEncodeByteSlice(int(getMagicNum("Supply")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSupply(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandSupplyI(r RandSrc) SupplyI {
	switch r.GetUint() % 1 {
	case 0:
		return RandSupply(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopySupplyI(x SupplyI) SupplyI {
	switch v := x.(type) {
	case Supply:
		res := DeepCopySupply(v)
		return res
	case *Supply:
		res := DeepCopySupply(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
// Interface
func DecodeToken(bz []byte) (v Token, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 412190080:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp BaseToken
		tmp, n, err = DecodeBaseToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = &tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeToken
func EncodeToken(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case BaseToken:
		codonEncodeByteSlice(int(getMagicNum("BaseToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseToken(w, v)
			return wBuf
		}())
	case *BaseToken:
		codonEncodeByteSlice(int(getMagicNum("BaseToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseToken(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func RandToken(r RandSrc) Token {
	switch r.GetUint() % 1 {
	case 0:
		tmp := RandBaseToken(r)
		return &tmp
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyToken(x Token) Token {
	switch v := x.(type) {
	case *BaseToken:
		res := DeepCopyBaseToken(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func getMagicNum(name string) uint32 {
	switch name {
	case "AccAddress":
		return 45679485
	case "AccAddressList":
		return 156864174
	case "AccountX":
		return 59743196
	case "BaseAccount":
		return 26591049
	case "BaseToken":
		return 412190080
	case "BaseVestingAccount":
		return 351017031
	case "Coin":
		return 206716790
	case "CommentRef":
		return 311958167
	case "CommitInfo":
		return 19375602
	case "CommunityPoolSpendProposal":
		return 23847857
	case "ConsAddress":
		return 373125373
	case "ContinuousVestingAccount":
		return 82892118
	case "DecCoin":
		return 1055185
	case "DelayedVestingAccount":
		return 12374068
	case "FeePool":
		return 130333088
	case "Input":
		return 495419206
	case "LockedCoin":
		return 149648993
	case "MarketInfo":
		return 413425332
	case "ModuleAccount":
		return 321193714
	case "MsgAddTokenWhitelist":
		return 153655625
	case "MsgAliasUpdate":
		return 285022960
	case "MsgBancorCancel":
		return 495620655
	case "MsgBancorInit":
		return 115719891
	case "MsgBancorTrade":
		return 158991245
	case "MsgBeginRedelegate":
		return 241505485
	case "MsgBurnToken":
		return 449833530
	case "MsgCancelOrder":
		return 304885884
	case "MsgCancelTradingPair":
		return 332906025
	case "MsgCommentToken":
		return 56755914
	case "MsgCreateOrder":
		return 377235122
	case "MsgCreateTradingPair":
		return 484769215
	case "MsgCreateValidator":
		return 374726888
	case "MsgDelegate":
		return 106121128
	case "MsgDeposit":
		return 402428987
	case "MsgDonateToCommunityPool":
		return 138691582
	case "MsgEditValidator":
		return 502447547
	case "MsgForbidAddr":
		return 488219826
	case "MsgForbidToken":
		return 286625767
	case "MsgIssueToken":
		return 336682718
	case "MsgMintToken":
		return 135231766
	case "MsgModifyPricePrecision":
		return 380178367
	case "MsgModifyTokenInfo":
		return 405536984
	case "MsgMultiSend":
		return 167957709
	case "MsgMultiSendX":
		return 155708190
	case "MsgRemoveTokenWhitelist":
		return 94212241
	case "MsgSend":
		return 31944689
	case "MsgSendX":
		return 145878995
	case "MsgSetMemoRequired":
		return 376054118
	case "MsgSetWithdrawAddress":
		return 73614168
	case "MsgSubmitProposal":
		return 331897574
	case "MsgSupervisedSend":
		return 87007006
	case "MsgTransferOwnership":
		return 377250605
	case "MsgUnForbidAddr":
		return 92120341
	case "MsgUnForbidToken":
		return 355758274
	case "MsgUndelegate":
		return 479724980
	case "MsgUnjail":
		return 261008494
	case "MsgVerifyInvariant":
		return 121564406
	case "MsgVote":
		return 314113032
	case "MsgWithdrawDelegatorReward":
		return 23904706
	case "MsgWithdrawValidatorCommission":
		return 514264249
	case "Order":
		return 168833459
	case "Output":
		return 337371032
	case "ParamChange":
		return 19646903
	case "ParameterChangeProposal":
		return 84061114
	case "PrivKeyEd25519":
		return 454630091
	case "PrivKeySecp256k1":
		return 252074570
	case "PubKeyEd25519":
		return 32929284
	case "PubKeyMultisigThreshold":
		return 471601077
	case "PubKeySecp256k1":
		return 45700893
	case "SdkDec":
		return 136348929
	case "SdkInt":
		return 497782922
	case "SignedMsgType":
		return 294706048
	case "SoftwareUpgradeProposal":
		return 261665006
	case "State":
		return 25024364
	case "StdSignature":
		return 496135554
	case "StdTx":
		return 420461272
	case "StoreInfo":
		return 415199234
	case "Supply":
		return 288082478
	case "TextProposal":
		return 458846397
	case "Vote":
		return 236500255
	case "VoteOption":
		return 484880195
	case "int64":
		return 344770825
	case "uint64":
		return 318185791
	} // end of switch
	panic("Should not reach here")
	return 0
} // end of getMagicNum
func getMagicNumOfVar(x interface{}) (uint32, bool) {
	switch x.(type) {
	case *AccAddress, AccAddress:
		return 45679485, true
	case *AccAddressList, AccAddressList:
		return 156864174, true
	case *AccountX, AccountX:
		return 59743196, true
	case *BaseAccount, BaseAccount:
		return 26591049, true
	case *BaseToken, BaseToken:
		return 412190080, true
	case *BaseVestingAccount, BaseVestingAccount:
		return 351017031, true
	case *Coin, Coin:
		return 206716790, true
	case *CommentRef, CommentRef:
		return 311958167, true
	case *CommitInfo, CommitInfo:
		return 19375602, true
	case *CommunityPoolSpendProposal, CommunityPoolSpendProposal:
		return 23847857, true
	case *ConsAddress, ConsAddress:
		return 373125373, true
	case *ContinuousVestingAccount, ContinuousVestingAccount:
		return 82892118, true
	case *DecCoin, DecCoin:
		return 1055185, true
	case *DelayedVestingAccount, DelayedVestingAccount:
		return 12374068, true
	case *FeePool, FeePool:
		return 130333088, true
	case *Input, Input:
		return 495419206, true
	case *LockedCoin, LockedCoin:
		return 149648993, true
	case *MarketInfo, MarketInfo:
		return 413425332, true
	case *ModuleAccount, ModuleAccount:
		return 321193714, true
	case *MsgAddTokenWhitelist, MsgAddTokenWhitelist:
		return 153655625, true
	case *MsgAliasUpdate, MsgAliasUpdate:
		return 285022960, true
	case *MsgBancorCancel, MsgBancorCancel:
		return 495620655, true
	case *MsgBancorInit, MsgBancorInit:
		return 115719891, true
	case *MsgBancorTrade, MsgBancorTrade:
		return 158991245, true
	case *MsgBeginRedelegate, MsgBeginRedelegate:
		return 241505485, true
	case *MsgBurnToken, MsgBurnToken:
		return 449833530, true
	case *MsgCancelOrder, MsgCancelOrder:
		return 304885884, true
	case *MsgCancelTradingPair, MsgCancelTradingPair:
		return 332906025, true
	case *MsgCommentToken, MsgCommentToken:
		return 56755914, true
	case *MsgCreateOrder, MsgCreateOrder:
		return 377235122, true
	case *MsgCreateTradingPair, MsgCreateTradingPair:
		return 484769215, true
	case *MsgCreateValidator, MsgCreateValidator:
		return 374726888, true
	case *MsgDelegate, MsgDelegate:
		return 106121128, true
	case *MsgDeposit, MsgDeposit:
		return 402428987, true
	case *MsgDonateToCommunityPool, MsgDonateToCommunityPool:
		return 138691582, true
	case *MsgEditValidator, MsgEditValidator:
		return 502447547, true
	case *MsgForbidAddr, MsgForbidAddr:
		return 488219826, true
	case *MsgForbidToken, MsgForbidToken:
		return 286625767, true
	case *MsgIssueToken, MsgIssueToken:
		return 336682718, true
	case *MsgMintToken, MsgMintToken:
		return 135231766, true
	case *MsgModifyPricePrecision, MsgModifyPricePrecision:
		return 380178367, true
	case *MsgModifyTokenInfo, MsgModifyTokenInfo:
		return 405536984, true
	case *MsgMultiSend, MsgMultiSend:
		return 167957709, true
	case *MsgMultiSendX, MsgMultiSendX:
		return 155708190, true
	case *MsgRemoveTokenWhitelist, MsgRemoveTokenWhitelist:
		return 94212241, true
	case *MsgSend, MsgSend:
		return 31944689, true
	case *MsgSendX, MsgSendX:
		return 145878995, true
	case *MsgSetMemoRequired, MsgSetMemoRequired:
		return 376054118, true
	case *MsgSetWithdrawAddress, MsgSetWithdrawAddress:
		return 73614168, true
	case *MsgSubmitProposal, MsgSubmitProposal:
		return 331897574, true
	case *MsgSupervisedSend, MsgSupervisedSend:
		return 87007006, true
	case *MsgTransferOwnership, MsgTransferOwnership:
		return 377250605, true
	case *MsgUnForbidAddr, MsgUnForbidAddr:
		return 92120341, true
	case *MsgUnForbidToken, MsgUnForbidToken:
		return 355758274, true
	case *MsgUndelegate, MsgUndelegate:
		return 479724980, true
	case *MsgUnjail, MsgUnjail:
		return 261008494, true
	case *MsgVerifyInvariant, MsgVerifyInvariant:
		return 121564406, true
	case *MsgVote, MsgVote:
		return 314113032, true
	case *MsgWithdrawDelegatorReward, MsgWithdrawDelegatorReward:
		return 23904706, true
	case *MsgWithdrawValidatorCommission, MsgWithdrawValidatorCommission:
		return 514264249, true
	case *Order, Order:
		return 168833459, true
	case *Output, Output:
		return 337371032, true
	case *ParamChange, ParamChange:
		return 19646903, true
	case *ParameterChangeProposal, ParameterChangeProposal:
		return 84061114, true
	case *PrivKeyEd25519, PrivKeyEd25519:
		return 454630091, true
	case *PrivKeySecp256k1, PrivKeySecp256k1:
		return 252074570, true
	case *PubKeyEd25519, PubKeyEd25519:
		return 32929284, true
	case *PubKeyMultisigThreshold, PubKeyMultisigThreshold:
		return 471601077, true
	case *PubKeySecp256k1, PubKeySecp256k1:
		return 45700893, true
	case *SdkDec, SdkDec:
		return 136348929, true
	case *SdkInt, SdkInt:
		return 497782922, true
	case *SignedMsgType, SignedMsgType:
		return 294706048, true
	case *SoftwareUpgradeProposal, SoftwareUpgradeProposal:
		return 261665006, true
	case *State, State:
		return 25024364, true
	case *StdSignature, StdSignature:
		return 496135554, true
	case *StdTx, StdTx:
		return 420461272, true
	case *StoreInfo, StoreInfo:
		return 415199234, true
	case *Supply, Supply:
		return 288082478, true
	case *TextProposal, TextProposal:
		return 458846397, true
	case *Vote, Vote:
		return 236500255, true
	case *VoteOption, VoteOption:
		return 484880195, true
	case *int64, int64:
		return 344770825, true
	case *uint64, uint64:
		return 318185791, true
	default:
		return 0, false
	} // end of switch
} // end of func
func EncodeAny(w *[]byte, x interface{}) {
	switch v := x.(type) {
	case AccAddress:
		codonEncodeByteSlice(int(getMagicNum("AccAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeAccAddress(w, v)
			return wBuf
		}())
	case *AccAddress:
		codonEncodeByteSlice(int(getMagicNum("AccAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeAccAddress(w, *v)
			return wBuf
		}())
	case AccAddressList:
		codonEncodeByteSlice(int(getMagicNum("AccAddressList")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeAccAddressList(w, v)
			return wBuf
		}())
	case *AccAddressList:
		codonEncodeByteSlice(int(getMagicNum("AccAddressList")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeAccAddressList(w, *v)
			return wBuf
		}())
	case AccountX:
		codonEncodeByteSlice(int(getMagicNum("AccountX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeAccountX(w, v)
			return wBuf
		}())
	case *AccountX:
		codonEncodeByteSlice(int(getMagicNum("AccountX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeAccountX(w, *v)
			return wBuf
		}())
	case BaseAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseAccount(w, v)
			return wBuf
		}())
	case *BaseAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseAccount(w, *v)
			return wBuf
		}())
	case BaseToken:
		codonEncodeByteSlice(int(getMagicNum("BaseToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseToken(w, v)
			return wBuf
		}())
	case *BaseToken:
		codonEncodeByteSlice(int(getMagicNum("BaseToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseToken(w, *v)
			return wBuf
		}())
	case BaseVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseVestingAccount(w, v)
			return wBuf
		}())
	case *BaseVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("BaseVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeBaseVestingAccount(w, *v)
			return wBuf
		}())
	case Coin:
		codonEncodeByteSlice(int(getMagicNum("Coin")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCoin(w, v)
			return wBuf
		}())
	case *Coin:
		codonEncodeByteSlice(int(getMagicNum("Coin")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCoin(w, *v)
			return wBuf
		}())
	case CommentRef:
		codonEncodeByteSlice(int(getMagicNum("CommentRef")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommentRef(w, v)
			return wBuf
		}())
	case *CommentRef:
		codonEncodeByteSlice(int(getMagicNum("CommentRef")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommentRef(w, *v)
			return wBuf
		}())
	case CommitInfo:
		codonEncodeByteSlice(int(getMagicNum("CommitInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommitInfo(w, v)
			return wBuf
		}())
	case *CommitInfo:
		codonEncodeByteSlice(int(getMagicNum("CommitInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommitInfo(w, *v)
			return wBuf
		}())
	case CommunityPoolSpendProposal:
		codonEncodeByteSlice(int(getMagicNum("CommunityPoolSpendProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommunityPoolSpendProposal(w, v)
			return wBuf
		}())
	case *CommunityPoolSpendProposal:
		codonEncodeByteSlice(int(getMagicNum("CommunityPoolSpendProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeCommunityPoolSpendProposal(w, *v)
			return wBuf
		}())
	case ConsAddress:
		codonEncodeByteSlice(int(getMagicNum("ConsAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeConsAddress(w, v)
			return wBuf
		}())
	case *ConsAddress:
		codonEncodeByteSlice(int(getMagicNum("ConsAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeConsAddress(w, *v)
			return wBuf
		}())
	case ContinuousVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("ContinuousVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeContinuousVestingAccount(w, v)
			return wBuf
		}())
	case *ContinuousVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("ContinuousVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeContinuousVestingAccount(w, *v)
			return wBuf
		}())
	case DecCoin:
		codonEncodeByteSlice(int(getMagicNum("DecCoin")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDecCoin(w, v)
			return wBuf
		}())
	case *DecCoin:
		codonEncodeByteSlice(int(getMagicNum("DecCoin")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDecCoin(w, *v)
			return wBuf
		}())
	case DelayedVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("DelayedVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDelayedVestingAccount(w, v)
			return wBuf
		}())
	case *DelayedVestingAccount:
		codonEncodeByteSlice(int(getMagicNum("DelayedVestingAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeDelayedVestingAccount(w, *v)
			return wBuf
		}())
	case FeePool:
		codonEncodeByteSlice(int(getMagicNum("FeePool")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeFeePool(w, v)
			return wBuf
		}())
	case *FeePool:
		codonEncodeByteSlice(int(getMagicNum("FeePool")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeFeePool(w, *v)
			return wBuf
		}())
	case Input:
		codonEncodeByteSlice(int(getMagicNum("Input")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeInput(w, v)
			return wBuf
		}())
	case *Input:
		codonEncodeByteSlice(int(getMagicNum("Input")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeInput(w, *v)
			return wBuf
		}())
	case LockedCoin:
		codonEncodeByteSlice(int(getMagicNum("LockedCoin")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeLockedCoin(w, v)
			return wBuf
		}())
	case *LockedCoin:
		codonEncodeByteSlice(int(getMagicNum("LockedCoin")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeLockedCoin(w, *v)
			return wBuf
		}())
	case MarketInfo:
		codonEncodeByteSlice(int(getMagicNum("MarketInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMarketInfo(w, v)
			return wBuf
		}())
	case *MarketInfo:
		codonEncodeByteSlice(int(getMagicNum("MarketInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMarketInfo(w, *v)
			return wBuf
		}())
	case ModuleAccount:
		codonEncodeByteSlice(int(getMagicNum("ModuleAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeModuleAccount(w, v)
			return wBuf
		}())
	case *ModuleAccount:
		codonEncodeByteSlice(int(getMagicNum("ModuleAccount")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeModuleAccount(w, *v)
			return wBuf
		}())
	case MsgAddTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgAddTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAddTokenWhitelist(w, v)
			return wBuf
		}())
	case *MsgAddTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgAddTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAddTokenWhitelist(w, *v)
			return wBuf
		}())
	case MsgAliasUpdate:
		codonEncodeByteSlice(int(getMagicNum("MsgAliasUpdate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAliasUpdate(w, v)
			return wBuf
		}())
	case *MsgAliasUpdate:
		codonEncodeByteSlice(int(getMagicNum("MsgAliasUpdate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgAliasUpdate(w, *v)
			return wBuf
		}())
	case MsgBancorCancel:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorCancel")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorCancel(w, v)
			return wBuf
		}())
	case *MsgBancorCancel:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorCancel")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorCancel(w, *v)
			return wBuf
		}())
	case MsgBancorInit:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorInit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorInit(w, v)
			return wBuf
		}())
	case *MsgBancorInit:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorInit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorInit(w, *v)
			return wBuf
		}())
	case MsgBancorTrade:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorTrade")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorTrade(w, v)
			return wBuf
		}())
	case *MsgBancorTrade:
		codonEncodeByteSlice(int(getMagicNum("MsgBancorTrade")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBancorTrade(w, *v)
			return wBuf
		}())
	case MsgBeginRedelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgBeginRedelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBeginRedelegate(w, v)
			return wBuf
		}())
	case *MsgBeginRedelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgBeginRedelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBeginRedelegate(w, *v)
			return wBuf
		}())
	case MsgBurnToken:
		codonEncodeByteSlice(int(getMagicNum("MsgBurnToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBurnToken(w, v)
			return wBuf
		}())
	case *MsgBurnToken:
		codonEncodeByteSlice(int(getMagicNum("MsgBurnToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgBurnToken(w, *v)
			return wBuf
		}())
	case MsgCancelOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelOrder(w, v)
			return wBuf
		}())
	case *MsgCancelOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelOrder(w, *v)
			return wBuf
		}())
	case MsgCancelTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelTradingPair(w, v)
			return wBuf
		}())
	case *MsgCancelTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCancelTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCancelTradingPair(w, *v)
			return wBuf
		}())
	case MsgCommentToken:
		codonEncodeByteSlice(int(getMagicNum("MsgCommentToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCommentToken(w, v)
			return wBuf
		}())
	case *MsgCommentToken:
		codonEncodeByteSlice(int(getMagicNum("MsgCommentToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCommentToken(w, *v)
			return wBuf
		}())
	case MsgCreateOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateOrder(w, v)
			return wBuf
		}())
	case *MsgCreateOrder:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateOrder")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateOrder(w, *v)
			return wBuf
		}())
	case MsgCreateTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateTradingPair(w, v)
			return wBuf
		}())
	case *MsgCreateTradingPair:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateTradingPair")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateTradingPair(w, *v)
			return wBuf
		}())
	case MsgCreateValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateValidator(w, v)
			return wBuf
		}())
	case *MsgCreateValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgCreateValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgCreateValidator(w, *v)
			return wBuf
		}())
	case MsgDelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgDelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDelegate(w, v)
			return wBuf
		}())
	case *MsgDelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgDelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDelegate(w, *v)
			return wBuf
		}())
	case MsgDeposit:
		codonEncodeByteSlice(int(getMagicNum("MsgDeposit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDeposit(w, v)
			return wBuf
		}())
	case *MsgDeposit:
		codonEncodeByteSlice(int(getMagicNum("MsgDeposit")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDeposit(w, *v)
			return wBuf
		}())
	case MsgDonateToCommunityPool:
		codonEncodeByteSlice(int(getMagicNum("MsgDonateToCommunityPool")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDonateToCommunityPool(w, v)
			return wBuf
		}())
	case *MsgDonateToCommunityPool:
		codonEncodeByteSlice(int(getMagicNum("MsgDonateToCommunityPool")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgDonateToCommunityPool(w, *v)
			return wBuf
		}())
	case MsgEditValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgEditValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgEditValidator(w, v)
			return wBuf
		}())
	case *MsgEditValidator:
		codonEncodeByteSlice(int(getMagicNum("MsgEditValidator")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgEditValidator(w, *v)
			return wBuf
		}())
	case MsgForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidAddr(w, v)
			return wBuf
		}())
	case *MsgForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidAddr(w, *v)
			return wBuf
		}())
	case MsgForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidToken(w, v)
			return wBuf
		}())
	case *MsgForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgForbidToken(w, *v)
			return wBuf
		}())
	case MsgIssueToken:
		codonEncodeByteSlice(int(getMagicNum("MsgIssueToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgIssueToken(w, v)
			return wBuf
		}())
	case *MsgIssueToken:
		codonEncodeByteSlice(int(getMagicNum("MsgIssueToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgIssueToken(w, *v)
			return wBuf
		}())
	case MsgMintToken:
		codonEncodeByteSlice(int(getMagicNum("MsgMintToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMintToken(w, v)
			return wBuf
		}())
	case *MsgMintToken:
		codonEncodeByteSlice(int(getMagicNum("MsgMintToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMintToken(w, *v)
			return wBuf
		}())
	case MsgModifyPricePrecision:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyPricePrecision")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyPricePrecision(w, v)
			return wBuf
		}())
	case *MsgModifyPricePrecision:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyPricePrecision")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyPricePrecision(w, *v)
			return wBuf
		}())
	case MsgModifyTokenInfo:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyTokenInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyTokenInfo(w, v)
			return wBuf
		}())
	case *MsgModifyTokenInfo:
		codonEncodeByteSlice(int(getMagicNum("MsgModifyTokenInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgModifyTokenInfo(w, *v)
			return wBuf
		}())
	case MsgMultiSend:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSend(w, v)
			return wBuf
		}())
	case *MsgMultiSend:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSend(w, *v)
			return wBuf
		}())
	case MsgMultiSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSendX(w, v)
			return wBuf
		}())
	case *MsgMultiSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgMultiSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgMultiSendX(w, *v)
			return wBuf
		}())
	case MsgRemoveTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgRemoveTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgRemoveTokenWhitelist(w, v)
			return wBuf
		}())
	case *MsgRemoveTokenWhitelist:
		codonEncodeByteSlice(int(getMagicNum("MsgRemoveTokenWhitelist")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgRemoveTokenWhitelist(w, *v)
			return wBuf
		}())
	case MsgSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSend(w, v)
			return wBuf
		}())
	case *MsgSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSend(w, *v)
			return wBuf
		}())
	case MsgSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSendX(w, v)
			return wBuf
		}())
	case *MsgSendX:
		codonEncodeByteSlice(int(getMagicNum("MsgSendX")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSendX(w, *v)
			return wBuf
		}())
	case MsgSetMemoRequired:
		codonEncodeByteSlice(int(getMagicNum("MsgSetMemoRequired")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetMemoRequired(w, v)
			return wBuf
		}())
	case *MsgSetMemoRequired:
		codonEncodeByteSlice(int(getMagicNum("MsgSetMemoRequired")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetMemoRequired(w, *v)
			return wBuf
		}())
	case MsgSetWithdrawAddress:
		codonEncodeByteSlice(int(getMagicNum("MsgSetWithdrawAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetWithdrawAddress(w, v)
			return wBuf
		}())
	case *MsgSetWithdrawAddress:
		codonEncodeByteSlice(int(getMagicNum("MsgSetWithdrawAddress")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSetWithdrawAddress(w, *v)
			return wBuf
		}())
	case MsgSubmitProposal:
		codonEncodeByteSlice(int(getMagicNum("MsgSubmitProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSubmitProposal(w, v)
			return wBuf
		}())
	case *MsgSubmitProposal:
		codonEncodeByteSlice(int(getMagicNum("MsgSubmitProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSubmitProposal(w, *v)
			return wBuf
		}())
	case MsgSupervisedSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSupervisedSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSupervisedSend(w, v)
			return wBuf
		}())
	case *MsgSupervisedSend:
		codonEncodeByteSlice(int(getMagicNum("MsgSupervisedSend")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgSupervisedSend(w, *v)
			return wBuf
		}())
	case MsgTransferOwnership:
		codonEncodeByteSlice(int(getMagicNum("MsgTransferOwnership")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgTransferOwnership(w, v)
			return wBuf
		}())
	case *MsgTransferOwnership:
		codonEncodeByteSlice(int(getMagicNum("MsgTransferOwnership")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgTransferOwnership(w, *v)
			return wBuf
		}())
	case MsgUnForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidAddr(w, v)
			return wBuf
		}())
	case *MsgUnForbidAddr:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidAddr")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidAddr(w, *v)
			return wBuf
		}())
	case MsgUnForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidToken(w, v)
			return wBuf
		}())
	case *MsgUnForbidToken:
		codonEncodeByteSlice(int(getMagicNum("MsgUnForbidToken")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnForbidToken(w, *v)
			return wBuf
		}())
	case MsgUndelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgUndelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUndelegate(w, v)
			return wBuf
		}())
	case *MsgUndelegate:
		codonEncodeByteSlice(int(getMagicNum("MsgUndelegate")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUndelegate(w, *v)
			return wBuf
		}())
	case MsgUnjail:
		codonEncodeByteSlice(int(getMagicNum("MsgUnjail")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnjail(w, v)
			return wBuf
		}())
	case *MsgUnjail:
		codonEncodeByteSlice(int(getMagicNum("MsgUnjail")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgUnjail(w, *v)
			return wBuf
		}())
	case MsgVerifyInvariant:
		codonEncodeByteSlice(int(getMagicNum("MsgVerifyInvariant")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVerifyInvariant(w, v)
			return wBuf
		}())
	case *MsgVerifyInvariant:
		codonEncodeByteSlice(int(getMagicNum("MsgVerifyInvariant")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVerifyInvariant(w, *v)
			return wBuf
		}())
	case MsgVote:
		codonEncodeByteSlice(int(getMagicNum("MsgVote")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVote(w, v)
			return wBuf
		}())
	case *MsgVote:
		codonEncodeByteSlice(int(getMagicNum("MsgVote")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgVote(w, *v)
			return wBuf
		}())
	case MsgWithdrawDelegatorReward:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawDelegatorReward")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawDelegatorReward(w, v)
			return wBuf
		}())
	case *MsgWithdrawDelegatorReward:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawDelegatorReward")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawDelegatorReward(w, *v)
			return wBuf
		}())
	case MsgWithdrawValidatorCommission:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawValidatorCommission")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawValidatorCommission(w, v)
			return wBuf
		}())
	case *MsgWithdrawValidatorCommission:
		codonEncodeByteSlice(int(getMagicNum("MsgWithdrawValidatorCommission")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeMsgWithdrawValidatorCommission(w, *v)
			return wBuf
		}())
	case Order:
		codonEncodeByteSlice(int(getMagicNum("Order")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeOrder(w, v)
			return wBuf
		}())
	case *Order:
		codonEncodeByteSlice(int(getMagicNum("Order")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeOrder(w, *v)
			return wBuf
		}())
	case Output:
		codonEncodeByteSlice(int(getMagicNum("Output")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeOutput(w, v)
			return wBuf
		}())
	case *Output:
		codonEncodeByteSlice(int(getMagicNum("Output")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeOutput(w, *v)
			return wBuf
		}())
	case ParamChange:
		codonEncodeByteSlice(int(getMagicNum("ParamChange")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeParamChange(w, v)
			return wBuf
		}())
	case *ParamChange:
		codonEncodeByteSlice(int(getMagicNum("ParamChange")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeParamChange(w, *v)
			return wBuf
		}())
	case ParameterChangeProposal:
		codonEncodeByteSlice(int(getMagicNum("ParameterChangeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeParameterChangeProposal(w, v)
			return wBuf
		}())
	case *ParameterChangeProposal:
		codonEncodeByteSlice(int(getMagicNum("ParameterChangeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeParameterChangeProposal(w, *v)
			return wBuf
		}())
	case PrivKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PrivKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeyEd25519(w, v)
			return wBuf
		}())
	case *PrivKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PrivKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeyEd25519(w, *v)
			return wBuf
		}())
	case PrivKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PrivKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeySecp256k1(w, v)
			return wBuf
		}())
	case *PrivKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PrivKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePrivKeySecp256k1(w, *v)
			return wBuf
		}())
	case PubKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PubKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyEd25519(w, v)
			return wBuf
		}())
	case *PubKeyEd25519:
		codonEncodeByteSlice(int(getMagicNum("PubKeyEd25519")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyEd25519(w, *v)
			return wBuf
		}())
	case PubKeyMultisigThreshold:
		codonEncodeByteSlice(int(getMagicNum("PubKeyMultisigThreshold")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyMultisigThreshold(w, v)
			return wBuf
		}())
	case *PubKeyMultisigThreshold:
		codonEncodeByteSlice(int(getMagicNum("PubKeyMultisigThreshold")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeyMultisigThreshold(w, *v)
			return wBuf
		}())
	case PubKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PubKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeySecp256k1(w, v)
			return wBuf
		}())
	case *PubKeySecp256k1:
		codonEncodeByteSlice(int(getMagicNum("PubKeySecp256k1")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodePubKeySecp256k1(w, *v)
			return wBuf
		}())
	case SdkDec:
		codonEncodeByteSlice(int(getMagicNum("SdkDec")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSdkDec(w, v)
			return wBuf
		}())
	case *SdkDec:
		codonEncodeByteSlice(int(getMagicNum("SdkDec")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSdkDec(w, *v)
			return wBuf
		}())
	case SdkInt:
		codonEncodeByteSlice(int(getMagicNum("SdkInt")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSdkInt(w, v)
			return wBuf
		}())
	case *SdkInt:
		codonEncodeByteSlice(int(getMagicNum("SdkInt")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSdkInt(w, *v)
			return wBuf
		}())
	case SignedMsgType:
		codonEncodeByteSlice(int(getMagicNum("SignedMsgType")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSignedMsgType(w, v)
			return wBuf
		}())
	case *SignedMsgType:
		codonEncodeByteSlice(int(getMagicNum("SignedMsgType")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSignedMsgType(w, *v)
			return wBuf
		}())
	case SoftwareUpgradeProposal:
		codonEncodeByteSlice(int(getMagicNum("SoftwareUpgradeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSoftwareUpgradeProposal(w, v)
			return wBuf
		}())
	case *SoftwareUpgradeProposal:
		codonEncodeByteSlice(int(getMagicNum("SoftwareUpgradeProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSoftwareUpgradeProposal(w, *v)
			return wBuf
		}())
	case State:
		codonEncodeByteSlice(int(getMagicNum("State")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeState(w, v)
			return wBuf
		}())
	case *State:
		codonEncodeByteSlice(int(getMagicNum("State")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeState(w, *v)
			return wBuf
		}())
	case StdSignature:
		codonEncodeByteSlice(int(getMagicNum("StdSignature")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdSignature(w, v)
			return wBuf
		}())
	case *StdSignature:
		codonEncodeByteSlice(int(getMagicNum("StdSignature")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdSignature(w, *v)
			return wBuf
		}())
	case StdTx:
		codonEncodeByteSlice(int(getMagicNum("StdTx")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdTx(w, v)
			return wBuf
		}())
	case *StdTx:
		codonEncodeByteSlice(int(getMagicNum("StdTx")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStdTx(w, *v)
			return wBuf
		}())
	case StoreInfo:
		codonEncodeByteSlice(int(getMagicNum("StoreInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStoreInfo(w, v)
			return wBuf
		}())
	case *StoreInfo:
		codonEncodeByteSlice(int(getMagicNum("StoreInfo")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeStoreInfo(w, *v)
			return wBuf
		}())
	case Supply:
		codonEncodeByteSlice(int(getMagicNum("Supply")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSupply(w, v)
			return wBuf
		}())
	case *Supply:
		codonEncodeByteSlice(int(getMagicNum("Supply")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeSupply(w, *v)
			return wBuf
		}())
	case TextProposal:
		codonEncodeByteSlice(int(getMagicNum("TextProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeTextProposal(w, v)
			return wBuf
		}())
	case *TextProposal:
		codonEncodeByteSlice(int(getMagicNum("TextProposal")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeTextProposal(w, *v)
			return wBuf
		}())
	case Vote:
		codonEncodeByteSlice(int(getMagicNum("Vote")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeVote(w, v)
			return wBuf
		}())
	case *Vote:
		codonEncodeByteSlice(int(getMagicNum("Vote")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeVote(w, *v)
			return wBuf
		}())
	case VoteOption:
		codonEncodeByteSlice(int(getMagicNum("VoteOption")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeVoteOption(w, v)
			return wBuf
		}())
	case *VoteOption:
		codonEncodeByteSlice(int(getMagicNum("VoteOption")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			EncodeVoteOption(w, *v)
			return wBuf
		}())
	case int64:
		codonEncodeByteSlice(int(getMagicNum("int64")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			Encodeint64(w, v)
			return wBuf
		}())
	case *int64:
		codonEncodeByteSlice(int(getMagicNum("int64")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			Encodeint64(w, *v)
			return wBuf
		}())
	case uint64:
		codonEncodeByteSlice(int(getMagicNum("uint64")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			Encodeuint64(w, v)
			return wBuf
		}())
	case *uint64:
		codonEncodeByteSlice(int(getMagicNum("uint64")), w, func() []byte {
			wBuf := make([]byte, 0, 64)
			w := &wBuf
			Encodeuint64(w, *v)
			return wBuf
		}())
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func DecodeAny(bz []byte) (v interface{}, total int, err error) {

	var n int
	tag := codonDecodeUint64(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	total += n
	magicNum := uint32(tag >> 3)
	switch magicNum {
	case 45679485:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp AccAddress
		tmp, n, err = DecodeAccAddress(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 156864174:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp AccAddressList
		tmp, n, err = DecodeAccAddressList(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 59743196:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp AccountX
		tmp, n, err = DecodeAccountX(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 26591049:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp BaseAccount
		tmp, n, err = DecodeBaseAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 412190080:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp BaseToken
		tmp, n, err = DecodeBaseToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 351017031:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp BaseVestingAccount
		tmp, n, err = DecodeBaseVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 206716790:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp Coin
		tmp, n, err = DecodeCoin(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 311958167:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp CommentRef
		tmp, n, err = DecodeCommentRef(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 19375602:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp CommitInfo
		tmp, n, err = DecodeCommitInfo(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 23847857:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp CommunityPoolSpendProposal
		tmp, n, err = DecodeCommunityPoolSpendProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 373125373:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ConsAddress
		tmp, n, err = DecodeConsAddress(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 82892118:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ContinuousVestingAccount
		tmp, n, err = DecodeContinuousVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 1055185:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp DecCoin
		tmp, n, err = DecodeDecCoin(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 12374068:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp DelayedVestingAccount
		tmp, n, err = DecodeDelayedVestingAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 130333088:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp FeePool
		tmp, n, err = DecodeFeePool(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 495419206:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp Input
		tmp, n, err = DecodeInput(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 149648993:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp LockedCoin
		tmp, n, err = DecodeLockedCoin(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 413425332:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MarketInfo
		tmp, n, err = DecodeMarketInfo(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 321193714:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ModuleAccount
		tmp, n, err = DecodeModuleAccount(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 153655625:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgAddTokenWhitelist
		tmp, n, err = DecodeMsgAddTokenWhitelist(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 285022960:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgAliasUpdate
		tmp, n, err = DecodeMsgAliasUpdate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 495620655:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBancorCancel
		tmp, n, err = DecodeMsgBancorCancel(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 115719891:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBancorInit
		tmp, n, err = DecodeMsgBancorInit(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 158991245:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBancorTrade
		tmp, n, err = DecodeMsgBancorTrade(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 241505485:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBeginRedelegate
		tmp, n, err = DecodeMsgBeginRedelegate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 449833530:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgBurnToken
		tmp, n, err = DecodeMsgBurnToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 304885884:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCancelOrder
		tmp, n, err = DecodeMsgCancelOrder(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 332906025:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCancelTradingPair
		tmp, n, err = DecodeMsgCancelTradingPair(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 56755914:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCommentToken
		tmp, n, err = DecodeMsgCommentToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 377235122:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCreateOrder
		tmp, n, err = DecodeMsgCreateOrder(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 484769215:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCreateTradingPair
		tmp, n, err = DecodeMsgCreateTradingPair(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 374726888:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgCreateValidator
		tmp, n, err = DecodeMsgCreateValidator(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 106121128:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgDelegate
		tmp, n, err = DecodeMsgDelegate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 402428987:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgDeposit
		tmp, n, err = DecodeMsgDeposit(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 138691582:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgDonateToCommunityPool
		tmp, n, err = DecodeMsgDonateToCommunityPool(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 502447547:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgEditValidator
		tmp, n, err = DecodeMsgEditValidator(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 488219826:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgForbidAddr
		tmp, n, err = DecodeMsgForbidAddr(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 286625767:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgForbidToken
		tmp, n, err = DecodeMsgForbidToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 336682718:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgIssueToken
		tmp, n, err = DecodeMsgIssueToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 135231766:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgMintToken
		tmp, n, err = DecodeMsgMintToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 380178367:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgModifyPricePrecision
		tmp, n, err = DecodeMsgModifyPricePrecision(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 405536984:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgModifyTokenInfo
		tmp, n, err = DecodeMsgModifyTokenInfo(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 167957709:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgMultiSend
		tmp, n, err = DecodeMsgMultiSend(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 155708190:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgMultiSendX
		tmp, n, err = DecodeMsgMultiSendX(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 94212241:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgRemoveTokenWhitelist
		tmp, n, err = DecodeMsgRemoveTokenWhitelist(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 31944689:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSend
		tmp, n, err = DecodeMsgSend(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 145878995:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSendX
		tmp, n, err = DecodeMsgSendX(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 376054118:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSetMemoRequired
		tmp, n, err = DecodeMsgSetMemoRequired(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 73614168:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSetWithdrawAddress
		tmp, n, err = DecodeMsgSetWithdrawAddress(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 331897574:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSubmitProposal
		tmp, n, err = DecodeMsgSubmitProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 87007006:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgSupervisedSend
		tmp, n, err = DecodeMsgSupervisedSend(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 377250605:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgTransferOwnership
		tmp, n, err = DecodeMsgTransferOwnership(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 92120341:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUnForbidAddr
		tmp, n, err = DecodeMsgUnForbidAddr(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 355758274:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUnForbidToken
		tmp, n, err = DecodeMsgUnForbidToken(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 479724980:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUndelegate
		tmp, n, err = DecodeMsgUndelegate(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 261008494:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgUnjail
		tmp, n, err = DecodeMsgUnjail(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 121564406:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgVerifyInvariant
		tmp, n, err = DecodeMsgVerifyInvariant(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 314113032:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgVote
		tmp, n, err = DecodeMsgVote(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 23904706:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgWithdrawDelegatorReward
		tmp, n, err = DecodeMsgWithdrawDelegatorReward(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 514264249:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp MsgWithdrawValidatorCommission
		tmp, n, err = DecodeMsgWithdrawValidatorCommission(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 168833459:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp Order
		tmp, n, err = DecodeOrder(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 337371032:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp Output
		tmp, n, err = DecodeOutput(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 19646903:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ParamChange
		tmp, n, err = DecodeParamChange(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 84061114:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp ParameterChangeProposal
		tmp, n, err = DecodeParameterChangeProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 454630091:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PrivKeyEd25519
		tmp, n, err = DecodePrivKeyEd25519(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 252074570:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PrivKeySecp256k1
		tmp, n, err = DecodePrivKeySecp256k1(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 32929284:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PubKeyEd25519
		tmp, n, err = DecodePubKeyEd25519(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 471601077:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PubKeyMultisigThreshold
		tmp, n, err = DecodePubKeyMultisigThreshold(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 45700893:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp PubKeySecp256k1
		tmp, n, err = DecodePubKeySecp256k1(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 136348929:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp SdkDec
		tmp, n, err = DecodeSdkDec(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 497782922:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp SdkInt
		tmp, n, err = DecodeSdkInt(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 294706048:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp SignedMsgType
		tmp, n, err = DecodeSignedMsgType(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 261665006:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp SoftwareUpgradeProposal
		tmp, n, err = DecodeSoftwareUpgradeProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 25024364:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp State
		tmp, n, err = DecodeState(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 496135554:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp StdSignature
		tmp, n, err = DecodeStdSignature(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 420461272:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp StdTx
		tmp, n, err = DecodeStdTx(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 415199234:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp StoreInfo
		tmp, n, err = DecodeStoreInfo(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 288082478:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp Supply
		tmp, n, err = DecodeSupply(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 458846397:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp TextProposal
		tmp, n, err = DecodeTextProposal(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 236500255:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp Vote
		tmp, n, err = DecodeVote(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 484880195:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp VoteOption
		tmp, n, err = DecodeVoteOption(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 344770825:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp int64
		tmp, n, err = Decodeint64(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	case 318185791:
		l := codonDecodeUint64(bz, &n, &err)
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) > len(bz) {
			err = errors.New("Length Too Large")
			return
		}
		var tmp uint64
		tmp, n, err = Decodeuint64(bz[:l])
		if err != nil {
			return
		}
		bz = bz[n:]
		total += n
		if int(l) != n {
			err = errors.New("Length Mismatch")
			return
		}
		v = tmp
		return
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeAny
func AssignIfcPtrFromStruct(ifcPtrIn interface{}, structObjIn interface{}) {
	switch ifcPtr := ifcPtrIn.(type) {
	case *Msg:
		switch structObj := structObjIn.(type) {
		case MsgWithdrawValidatorCommission:
			*ifcPtr = &structObj
		case MsgBancorCancel:
			*ifcPtr = &structObj
		case MsgMultiSend:
			*ifcPtr = &structObj
		case MsgSupervisedSend:
			*ifcPtr = &structObj
		case MsgDonateToCommunityPool:
			*ifcPtr = &structObj
		case MsgAddTokenWhitelist:
			*ifcPtr = &structObj
		case MsgEditValidator:
			*ifcPtr = &structObj
		case MsgRemoveTokenWhitelist:
			*ifcPtr = &structObj
		case MsgCreateValidator:
			*ifcPtr = &structObj
		case MsgTransferOwnership:
			*ifcPtr = &structObj
		case MsgUnForbidAddr:
			*ifcPtr = &structObj
		case MsgCommentToken:
			*ifcPtr = &structObj
		case MsgModifyTokenInfo:
			*ifcPtr = &structObj
		case MsgModifyPricePrecision:
			*ifcPtr = &structObj
		case MsgIssueToken:
			*ifcPtr = &structObj
		case MsgVerifyInvariant:
			*ifcPtr = &structObj
		case MsgBurnToken:
			*ifcPtr = &structObj
		case MsgSend:
			*ifcPtr = &structObj
		case MsgForbidAddr:
			*ifcPtr = &structObj
		case MsgForbidToken:
			*ifcPtr = &structObj
		case MsgCancelOrder:
			*ifcPtr = &structObj
		case MsgCreateTradingPair:
			*ifcPtr = &structObj
		case MsgUnjail:
			*ifcPtr = &structObj
		case MsgSetMemoRequired:
			*ifcPtr = &structObj
		case MsgUndelegate:
			*ifcPtr = &structObj
		case MsgUnForbidToken:
			*ifcPtr = &structObj
		case MsgBancorTrade:
			*ifcPtr = &structObj
		case MsgCancelTradingPair:
			*ifcPtr = &structObj
		case MsgDeposit:
			*ifcPtr = &structObj
		case MsgVote:
			*ifcPtr = &structObj
		case MsgBeginRedelegate:
			*ifcPtr = &structObj
		case MsgMultiSendX:
			*ifcPtr = &structObj
		case MsgMintToken:
			*ifcPtr = &structObj
		case MsgBancorInit:
			*ifcPtr = &structObj
		case MsgSetWithdrawAddress:
			*ifcPtr = &structObj
		case MsgSendX:
			*ifcPtr = &structObj
		case MsgCreateOrder:
			*ifcPtr = &structObj
		case MsgDelegate:
			*ifcPtr = &structObj
		case MsgAliasUpdate:
			*ifcPtr = &structObj
		case MsgWithdrawDelegatorReward:
			*ifcPtr = &structObj
		case MsgSubmitProposal:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *Account:
		switch structObj := structObjIn.(type) {
		case ContinuousVestingAccount:
			*ifcPtr = &structObj
		case DelayedVestingAccount:
			*ifcPtr = &structObj
		case ModuleAccount:
			*ifcPtr = &structObj
		case BaseAccount:
			*ifcPtr = &structObj
		case BaseVestingAccount:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *Tx:
		switch structObj := structObjIn.(type) {
		case StdTx:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *SupplyI:
		switch structObj := structObjIn.(type) {
		case Supply:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *Token:
		switch structObj := structObjIn.(type) {
		case BaseToken:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *VestingAccount:
		switch structObj := structObjIn.(type) {
		case ContinuousVestingAccount:
			*ifcPtr = &structObj
		case DelayedVestingAccount:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *Content:
		switch structObj := structObjIn.(type) {
		case TextProposal:
			*ifcPtr = &structObj
		case ParameterChangeProposal:
			*ifcPtr = &structObj
		case CommunityPoolSpendProposal:
			*ifcPtr = &structObj
		case SoftwareUpgradeProposal:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *ModuleAccountI:
		switch structObj := structObjIn.(type) {
		case ModuleAccount:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *PubKey:
		switch structObj := structObjIn.(type) {
		case StdSignature:
			*ifcPtr = &structObj
		case PubKeyMultisigThreshold:
			*ifcPtr = &structObj
		case PubKeyEd25519:
			*ifcPtr = &structObj
		case PubKeySecp256k1:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	case *PrivKey:
		switch structObj := structObjIn.(type) {
		case PrivKeyEd25519:
			*ifcPtr = &structObj
		case PrivKeySecp256k1:
			*ifcPtr = &structObj
		default:
			panic(fmt.Sprintf("Type mismatch %v %v\n", reflect.TypeOf(ifcPtr), reflect.TypeOf(structObjIn)))
		} // end switch of structs
	default:
		panic(fmt.Sprintf("Unknown Type %v\n", reflect.TypeOf(ifcPtrIn)))
	} // end switch of interfaces
}
func RandAny(r RandSrc) interface{} {
	switch r.GetUint() % 83 {
	case 0:
		return RandAccAddress(r)
	case 1:
		return RandAccAddressList(r)
	case 2:
		return RandAccountX(r)
	case 3:
		return RandBaseAccount(r)
	case 4:
		return RandBaseToken(r)
	case 5:
		return RandBaseVestingAccount(r)
	case 6:
		return RandCoin(r)
	case 7:
		return RandCommentRef(r)
	case 8:
		return RandCommitInfo(r)
	case 9:
		return RandCommunityPoolSpendProposal(r)
	case 10:
		return RandConsAddress(r)
	case 11:
		return RandContinuousVestingAccount(r)
	case 12:
		return RandDecCoin(r)
	case 13:
		return RandDelayedVestingAccount(r)
	case 14:
		return RandFeePool(r)
	case 15:
		return RandInput(r)
	case 16:
		return RandLockedCoin(r)
	case 17:
		return RandMarketInfo(r)
	case 18:
		return RandModuleAccount(r)
	case 19:
		return RandMsgAddTokenWhitelist(r)
	case 20:
		return RandMsgAliasUpdate(r)
	case 21:
		return RandMsgBancorCancel(r)
	case 22:
		return RandMsgBancorInit(r)
	case 23:
		return RandMsgBancorTrade(r)
	case 24:
		return RandMsgBeginRedelegate(r)
	case 25:
		return RandMsgBurnToken(r)
	case 26:
		return RandMsgCancelOrder(r)
	case 27:
		return RandMsgCancelTradingPair(r)
	case 28:
		return RandMsgCommentToken(r)
	case 29:
		return RandMsgCreateOrder(r)
	case 30:
		return RandMsgCreateTradingPair(r)
	case 31:
		return RandMsgCreateValidator(r)
	case 32:
		return RandMsgDelegate(r)
	case 33:
		return RandMsgDeposit(r)
	case 34:
		return RandMsgDonateToCommunityPool(r)
	case 35:
		return RandMsgEditValidator(r)
	case 36:
		return RandMsgForbidAddr(r)
	case 37:
		return RandMsgForbidToken(r)
	case 38:
		return RandMsgIssueToken(r)
	case 39:
		return RandMsgMintToken(r)
	case 40:
		return RandMsgModifyPricePrecision(r)
	case 41:
		return RandMsgModifyTokenInfo(r)
	case 42:
		return RandMsgMultiSend(r)
	case 43:
		return RandMsgMultiSendX(r)
	case 44:
		return RandMsgRemoveTokenWhitelist(r)
	case 45:
		return RandMsgSend(r)
	case 46:
		return RandMsgSendX(r)
	case 47:
		return RandMsgSetMemoRequired(r)
	case 48:
		return RandMsgSetWithdrawAddress(r)
	case 49:
		return RandMsgSubmitProposal(r)
	case 50:
		return RandMsgSupervisedSend(r)
	case 51:
		return RandMsgTransferOwnership(r)
	case 52:
		return RandMsgUnForbidAddr(r)
	case 53:
		return RandMsgUnForbidToken(r)
	case 54:
		return RandMsgUndelegate(r)
	case 55:
		return RandMsgUnjail(r)
	case 56:
		return RandMsgVerifyInvariant(r)
	case 57:
		return RandMsgVote(r)
	case 58:
		return RandMsgWithdrawDelegatorReward(r)
	case 59:
		return RandMsgWithdrawValidatorCommission(r)
	case 60:
		return RandOrder(r)
	case 61:
		return RandOutput(r)
	case 62:
		return RandParamChange(r)
	case 63:
		return RandParameterChangeProposal(r)
	case 64:
		return RandPrivKeyEd25519(r)
	case 65:
		return RandPrivKeySecp256k1(r)
	case 66:
		return RandPubKeyEd25519(r)
	case 67:
		return RandPubKeyMultisigThreshold(r)
	case 68:
		return RandPubKeySecp256k1(r)
	case 69:
		return RandSdkDec(r)
	case 70:
		return RandSdkInt(r)
	case 71:
		return RandSignedMsgType(r)
	case 72:
		return RandSoftwareUpgradeProposal(r)
	case 73:
		return RandState(r)
	case 74:
		return RandStdSignature(r)
	case 75:
		return RandStdTx(r)
	case 76:
		return RandStoreInfo(r)
	case 77:
		return RandSupply(r)
	case 78:
		return RandTextProposal(r)
	case 79:
		return RandVote(r)
	case 80:
		return RandVoteOption(r)
	case 81:
		return Randint64(r)
	case 82:
		return Randuint64(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DeepCopyAny(x interface{}) interface{} {
	switch v := x.(type) {
	case AccAddress:
		res := DeepCopyAccAddress(v)
		return res
	case *AccAddress:
		res := DeepCopyAccAddress(*v)
		return &res
	case AccAddressList:
		res := DeepCopyAccAddressList(v)
		return res
	case *AccAddressList:
		res := DeepCopyAccAddressList(*v)
		return &res
	case AccountX:
		res := DeepCopyAccountX(v)
		return res
	case *AccountX:
		res := DeepCopyAccountX(*v)
		return &res
	case BaseAccount:
		res := DeepCopyBaseAccount(v)
		return res
	case *BaseAccount:
		res := DeepCopyBaseAccount(*v)
		return &res
	case BaseToken:
		res := DeepCopyBaseToken(v)
		return res
	case *BaseToken:
		res := DeepCopyBaseToken(*v)
		return &res
	case BaseVestingAccount:
		res := DeepCopyBaseVestingAccount(v)
		return res
	case *BaseVestingAccount:
		res := DeepCopyBaseVestingAccount(*v)
		return &res
	case Coin:
		res := DeepCopyCoin(v)
		return res
	case *Coin:
		res := DeepCopyCoin(*v)
		return &res
	case CommentRef:
		res := DeepCopyCommentRef(v)
		return res
	case *CommentRef:
		res := DeepCopyCommentRef(*v)
		return &res
	case CommitInfo:
		res := DeepCopyCommitInfo(v)
		return res
	case *CommitInfo:
		res := DeepCopyCommitInfo(*v)
		return &res
	case CommunityPoolSpendProposal:
		res := DeepCopyCommunityPoolSpendProposal(v)
		return res
	case *CommunityPoolSpendProposal:
		res := DeepCopyCommunityPoolSpendProposal(*v)
		return &res
	case ConsAddress:
		res := DeepCopyConsAddress(v)
		return res
	case *ConsAddress:
		res := DeepCopyConsAddress(*v)
		return &res
	case ContinuousVestingAccount:
		res := DeepCopyContinuousVestingAccount(v)
		return res
	case *ContinuousVestingAccount:
		res := DeepCopyContinuousVestingAccount(*v)
		return &res
	case DecCoin:
		res := DeepCopyDecCoin(v)
		return res
	case *DecCoin:
		res := DeepCopyDecCoin(*v)
		return &res
	case DelayedVestingAccount:
		res := DeepCopyDelayedVestingAccount(v)
		return res
	case *DelayedVestingAccount:
		res := DeepCopyDelayedVestingAccount(*v)
		return &res
	case FeePool:
		res := DeepCopyFeePool(v)
		return res
	case *FeePool:
		res := DeepCopyFeePool(*v)
		return &res
	case Input:
		res := DeepCopyInput(v)
		return res
	case *Input:
		res := DeepCopyInput(*v)
		return &res
	case LockedCoin:
		res := DeepCopyLockedCoin(v)
		return res
	case *LockedCoin:
		res := DeepCopyLockedCoin(*v)
		return &res
	case MarketInfo:
		res := DeepCopyMarketInfo(v)
		return res
	case *MarketInfo:
		res := DeepCopyMarketInfo(*v)
		return &res
	case ModuleAccount:
		res := DeepCopyModuleAccount(v)
		return res
	case *ModuleAccount:
		res := DeepCopyModuleAccount(*v)
		return &res
	case MsgAddTokenWhitelist:
		res := DeepCopyMsgAddTokenWhitelist(v)
		return res
	case *MsgAddTokenWhitelist:
		res := DeepCopyMsgAddTokenWhitelist(*v)
		return &res
	case MsgAliasUpdate:
		res := DeepCopyMsgAliasUpdate(v)
		return res
	case *MsgAliasUpdate:
		res := DeepCopyMsgAliasUpdate(*v)
		return &res
	case MsgBancorCancel:
		res := DeepCopyMsgBancorCancel(v)
		return res
	case *MsgBancorCancel:
		res := DeepCopyMsgBancorCancel(*v)
		return &res
	case MsgBancorInit:
		res := DeepCopyMsgBancorInit(v)
		return res
	case *MsgBancorInit:
		res := DeepCopyMsgBancorInit(*v)
		return &res
	case MsgBancorTrade:
		res := DeepCopyMsgBancorTrade(v)
		return res
	case *MsgBancorTrade:
		res := DeepCopyMsgBancorTrade(*v)
		return &res
	case MsgBeginRedelegate:
		res := DeepCopyMsgBeginRedelegate(v)
		return res
	case *MsgBeginRedelegate:
		res := DeepCopyMsgBeginRedelegate(*v)
		return &res
	case MsgBurnToken:
		res := DeepCopyMsgBurnToken(v)
		return res
	case *MsgBurnToken:
		res := DeepCopyMsgBurnToken(*v)
		return &res
	case MsgCancelOrder:
		res := DeepCopyMsgCancelOrder(v)
		return res
	case *MsgCancelOrder:
		res := DeepCopyMsgCancelOrder(*v)
		return &res
	case MsgCancelTradingPair:
		res := DeepCopyMsgCancelTradingPair(v)
		return res
	case *MsgCancelTradingPair:
		res := DeepCopyMsgCancelTradingPair(*v)
		return &res
	case MsgCommentToken:
		res := DeepCopyMsgCommentToken(v)
		return res
	case *MsgCommentToken:
		res := DeepCopyMsgCommentToken(*v)
		return &res
	case MsgCreateOrder:
		res := DeepCopyMsgCreateOrder(v)
		return res
	case *MsgCreateOrder:
		res := DeepCopyMsgCreateOrder(*v)
		return &res
	case MsgCreateTradingPair:
		res := DeepCopyMsgCreateTradingPair(v)
		return res
	case *MsgCreateTradingPair:
		res := DeepCopyMsgCreateTradingPair(*v)
		return &res
	case MsgCreateValidator:
		res := DeepCopyMsgCreateValidator(v)
		return res
	case *MsgCreateValidator:
		res := DeepCopyMsgCreateValidator(*v)
		return &res
	case MsgDelegate:
		res := DeepCopyMsgDelegate(v)
		return res
	case *MsgDelegate:
		res := DeepCopyMsgDelegate(*v)
		return &res
	case MsgDeposit:
		res := DeepCopyMsgDeposit(v)
		return res
	case *MsgDeposit:
		res := DeepCopyMsgDeposit(*v)
		return &res
	case MsgDonateToCommunityPool:
		res := DeepCopyMsgDonateToCommunityPool(v)
		return res
	case *MsgDonateToCommunityPool:
		res := DeepCopyMsgDonateToCommunityPool(*v)
		return &res
	case MsgEditValidator:
		res := DeepCopyMsgEditValidator(v)
		return res
	case *MsgEditValidator:
		res := DeepCopyMsgEditValidator(*v)
		return &res
	case MsgForbidAddr:
		res := DeepCopyMsgForbidAddr(v)
		return res
	case *MsgForbidAddr:
		res := DeepCopyMsgForbidAddr(*v)
		return &res
	case MsgForbidToken:
		res := DeepCopyMsgForbidToken(v)
		return res
	case *MsgForbidToken:
		res := DeepCopyMsgForbidToken(*v)
		return &res
	case MsgIssueToken:
		res := DeepCopyMsgIssueToken(v)
		return res
	case *MsgIssueToken:
		res := DeepCopyMsgIssueToken(*v)
		return &res
	case MsgMintToken:
		res := DeepCopyMsgMintToken(v)
		return res
	case *MsgMintToken:
		res := DeepCopyMsgMintToken(*v)
		return &res
	case MsgModifyPricePrecision:
		res := DeepCopyMsgModifyPricePrecision(v)
		return res
	case *MsgModifyPricePrecision:
		res := DeepCopyMsgModifyPricePrecision(*v)
		return &res
	case MsgModifyTokenInfo:
		res := DeepCopyMsgModifyTokenInfo(v)
		return res
	case *MsgModifyTokenInfo:
		res := DeepCopyMsgModifyTokenInfo(*v)
		return &res
	case MsgMultiSend:
		res := DeepCopyMsgMultiSend(v)
		return res
	case *MsgMultiSend:
		res := DeepCopyMsgMultiSend(*v)
		return &res
	case MsgMultiSendX:
		res := DeepCopyMsgMultiSendX(v)
		return res
	case *MsgMultiSendX:
		res := DeepCopyMsgMultiSendX(*v)
		return &res
	case MsgRemoveTokenWhitelist:
		res := DeepCopyMsgRemoveTokenWhitelist(v)
		return res
	case *MsgRemoveTokenWhitelist:
		res := DeepCopyMsgRemoveTokenWhitelist(*v)
		return &res
	case MsgSend:
		res := DeepCopyMsgSend(v)
		return res
	case *MsgSend:
		res := DeepCopyMsgSend(*v)
		return &res
	case MsgSendX:
		res := DeepCopyMsgSendX(v)
		return res
	case *MsgSendX:
		res := DeepCopyMsgSendX(*v)
		return &res
	case MsgSetMemoRequired:
		res := DeepCopyMsgSetMemoRequired(v)
		return res
	case *MsgSetMemoRequired:
		res := DeepCopyMsgSetMemoRequired(*v)
		return &res
	case MsgSetWithdrawAddress:
		res := DeepCopyMsgSetWithdrawAddress(v)
		return res
	case *MsgSetWithdrawAddress:
		res := DeepCopyMsgSetWithdrawAddress(*v)
		return &res
	case MsgSubmitProposal:
		res := DeepCopyMsgSubmitProposal(v)
		return res
	case *MsgSubmitProposal:
		res := DeepCopyMsgSubmitProposal(*v)
		return &res
	case MsgSupervisedSend:
		res := DeepCopyMsgSupervisedSend(v)
		return res
	case *MsgSupervisedSend:
		res := DeepCopyMsgSupervisedSend(*v)
		return &res
	case MsgTransferOwnership:
		res := DeepCopyMsgTransferOwnership(v)
		return res
	case *MsgTransferOwnership:
		res := DeepCopyMsgTransferOwnership(*v)
		return &res
	case MsgUnForbidAddr:
		res := DeepCopyMsgUnForbidAddr(v)
		return res
	case *MsgUnForbidAddr:
		res := DeepCopyMsgUnForbidAddr(*v)
		return &res
	case MsgUnForbidToken:
		res := DeepCopyMsgUnForbidToken(v)
		return res
	case *MsgUnForbidToken:
		res := DeepCopyMsgUnForbidToken(*v)
		return &res
	case MsgUndelegate:
		res := DeepCopyMsgUndelegate(v)
		return res
	case *MsgUndelegate:
		res := DeepCopyMsgUndelegate(*v)
		return &res
	case MsgUnjail:
		res := DeepCopyMsgUnjail(v)
		return res
	case *MsgUnjail:
		res := DeepCopyMsgUnjail(*v)
		return &res
	case MsgVerifyInvariant:
		res := DeepCopyMsgVerifyInvariant(v)
		return res
	case *MsgVerifyInvariant:
		res := DeepCopyMsgVerifyInvariant(*v)
		return &res
	case MsgVote:
		res := DeepCopyMsgVote(v)
		return res
	case *MsgVote:
		res := DeepCopyMsgVote(*v)
		return &res
	case MsgWithdrawDelegatorReward:
		res := DeepCopyMsgWithdrawDelegatorReward(v)
		return res
	case *MsgWithdrawDelegatorReward:
		res := DeepCopyMsgWithdrawDelegatorReward(*v)
		return &res
	case MsgWithdrawValidatorCommission:
		res := DeepCopyMsgWithdrawValidatorCommission(v)
		return res
	case *MsgWithdrawValidatorCommission:
		res := DeepCopyMsgWithdrawValidatorCommission(*v)
		return &res
	case Order:
		res := DeepCopyOrder(v)
		return res
	case *Order:
		res := DeepCopyOrder(*v)
		return &res
	case Output:
		res := DeepCopyOutput(v)
		return res
	case *Output:
		res := DeepCopyOutput(*v)
		return &res
	case ParamChange:
		res := DeepCopyParamChange(v)
		return res
	case *ParamChange:
		res := DeepCopyParamChange(*v)
		return &res
	case ParameterChangeProposal:
		res := DeepCopyParameterChangeProposal(v)
		return res
	case *ParameterChangeProposal:
		res := DeepCopyParameterChangeProposal(*v)
		return &res
	case PrivKeyEd25519:
		res := DeepCopyPrivKeyEd25519(v)
		return res
	case *PrivKeyEd25519:
		res := DeepCopyPrivKeyEd25519(*v)
		return &res
	case PrivKeySecp256k1:
		res := DeepCopyPrivKeySecp256k1(v)
		return res
	case *PrivKeySecp256k1:
		res := DeepCopyPrivKeySecp256k1(*v)
		return &res
	case PubKeyEd25519:
		res := DeepCopyPubKeyEd25519(v)
		return res
	case *PubKeyEd25519:
		res := DeepCopyPubKeyEd25519(*v)
		return &res
	case PubKeyMultisigThreshold:
		res := DeepCopyPubKeyMultisigThreshold(v)
		return res
	case *PubKeyMultisigThreshold:
		res := DeepCopyPubKeyMultisigThreshold(*v)
		return &res
	case PubKeySecp256k1:
		res := DeepCopyPubKeySecp256k1(v)
		return res
	case *PubKeySecp256k1:
		res := DeepCopyPubKeySecp256k1(*v)
		return &res
	case SdkDec:
		res := DeepCopySdkDec(v)
		return res
	case *SdkDec:
		res := DeepCopySdkDec(*v)
		return &res
	case SdkInt:
		res := DeepCopySdkInt(v)
		return res
	case *SdkInt:
		res := DeepCopySdkInt(*v)
		return &res
	case SignedMsgType:
		res := DeepCopySignedMsgType(v)
		return res
	case *SignedMsgType:
		res := DeepCopySignedMsgType(*v)
		return &res
	case SoftwareUpgradeProposal:
		res := DeepCopySoftwareUpgradeProposal(v)
		return res
	case *SoftwareUpgradeProposal:
		res := DeepCopySoftwareUpgradeProposal(*v)
		return &res
	case State:
		res := DeepCopyState(v)
		return res
	case *State:
		res := DeepCopyState(*v)
		return &res
	case StdSignature:
		res := DeepCopyStdSignature(v)
		return res
	case *StdSignature:
		res := DeepCopyStdSignature(*v)
		return &res
	case StdTx:
		res := DeepCopyStdTx(v)
		return res
	case *StdTx:
		res := DeepCopyStdTx(*v)
		return &res
	case StoreInfo:
		res := DeepCopyStoreInfo(v)
		return res
	case *StoreInfo:
		res := DeepCopyStoreInfo(*v)
		return &res
	case Supply:
		res := DeepCopySupply(v)
		return res
	case *Supply:
		res := DeepCopySupply(*v)
		return &res
	case TextProposal:
		res := DeepCopyTextProposal(v)
		return res
	case *TextProposal:
		res := DeepCopyTextProposal(*v)
		return &res
	case Vote:
		res := DeepCopyVote(v)
		return res
	case *Vote:
		res := DeepCopyVote(*v)
		return &res
	case VoteOption:
		res := DeepCopyVoteOption(v)
		return res
	case *VoteOption:
		res := DeepCopyVoteOption(*v)
		return &res
	case int64:
		res := DeepCopyint64(v)
		return res
	case *int64:
		res := DeepCopyint64(*v)
		return &res
	case uint64:
		res := DeepCopyuint64(v)
		return res
	case *uint64:
		res := DeepCopyuint64(*v)
		return &res
	default:
		panic(fmt.Sprintf("Unknown Type %v %v\n", x, reflect.TypeOf(x)))
	} // end of switch
} // end of func
func GetSupportList() []string {
	return []string{
		".int64",
		".uint64",
		"AccAddressList",
		"github.com/coinexchain/dex/modules/alias/internal/types.MsgAliasUpdate",
		"github.com/coinexchain/dex/modules/asset/internal/types.BaseToken",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgAddTokenWhitelist",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgBurnToken",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgForbidAddr",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgForbidToken",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgIssueToken",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgMintToken",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgModifyTokenInfo",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgRemoveTokenWhitelist",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgTransferOwnership",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgUnForbidAddr",
		"github.com/coinexchain/dex/modules/asset/internal/types.MsgUnForbidToken",
		"github.com/coinexchain/dex/modules/asset/internal/types.Token",
		"github.com/coinexchain/dex/modules/authx/internal/types.AccountX",
		"github.com/coinexchain/dex/modules/authx/internal/types.LockedCoin",
		"github.com/coinexchain/dex/modules/bancorlite/internal/types.MsgBancorCancel",
		"github.com/coinexchain/dex/modules/bancorlite/internal/types.MsgBancorInit",
		"github.com/coinexchain/dex/modules/bancorlite/internal/types.MsgBancorTrade",
		"github.com/coinexchain/dex/modules/bankx/internal/types.MsgMultiSend",
		"github.com/coinexchain/dex/modules/bankx/internal/types.MsgSend",
		"github.com/coinexchain/dex/modules/bankx/internal/types.MsgSetMemoRequired",
		"github.com/coinexchain/dex/modules/bankx/internal/types.MsgSupervisedSend",
		"github.com/coinexchain/dex/modules/comment/internal/types.CommentRef",
		"github.com/coinexchain/dex/modules/comment/internal/types.MsgCommentToken",
		"github.com/coinexchain/dex/modules/distributionx/types.MsgDonateToCommunityPool",
		"github.com/coinexchain/dex/modules/incentive/internal/types.State",
		"github.com/coinexchain/dex/modules/market/internal/types.MarketInfo",
		"github.com/coinexchain/dex/modules/market/internal/types.MsgCancelOrder",
		"github.com/coinexchain/dex/modules/market/internal/types.MsgCancelTradingPair",
		"github.com/coinexchain/dex/modules/market/internal/types.MsgCreateOrder",
		"github.com/coinexchain/dex/modules/market/internal/types.MsgCreateTradingPair",
		"github.com/coinexchain/dex/modules/market/internal/types.MsgModifyPricePrecision",
		"github.com/coinexchain/dex/modules/market/internal/types.Order",
		"github.com/cosmos/cosmos-sdk/store/rootmulti.commitInfo",
		"github.com/cosmos/cosmos-sdk/store/rootmulti.storeInfo",
		"github.com/cosmos/cosmos-sdk/types.AccAddress",
		"github.com/cosmos/cosmos-sdk/types.Coin",
		"github.com/cosmos/cosmos-sdk/types.ConsAddress",
		"github.com/cosmos/cosmos-sdk/types.Dec",
		"github.com/cosmos/cosmos-sdk/types.DecCoin",
		"github.com/cosmos/cosmos-sdk/types.Int",
		"github.com/cosmos/cosmos-sdk/types.Msg",
		"github.com/cosmos/cosmos-sdk/types.Tx",
		"github.com/cosmos/cosmos-sdk/x/auth/exported.Account",
		"github.com/cosmos/cosmos-sdk/x/auth/exported.VestingAccount",
		"github.com/cosmos/cosmos-sdk/x/auth/types.BaseAccount",
		"github.com/cosmos/cosmos-sdk/x/auth/types.BaseVestingAccount",
		"github.com/cosmos/cosmos-sdk/x/auth/types.ContinuousVestingAccount",
		"github.com/cosmos/cosmos-sdk/x/auth/types.DelayedVestingAccount",
		"github.com/cosmos/cosmos-sdk/x/auth/types.StdSignature",
		"github.com/cosmos/cosmos-sdk/x/auth/types.StdTx",
		"github.com/cosmos/cosmos-sdk/x/bank/internal/types.Input",
		"github.com/cosmos/cosmos-sdk/x/bank/internal/types.MsgMultiSend",
		"github.com/cosmos/cosmos-sdk/x/bank/internal/types.MsgSend",
		"github.com/cosmos/cosmos-sdk/x/bank/internal/types.Output",
		"github.com/cosmos/cosmos-sdk/x/crisis/internal/types.MsgVerifyInvariant",
		"github.com/cosmos/cosmos-sdk/x/distribution/types.CommunityPoolSpendProposal",
		"github.com/cosmos/cosmos-sdk/x/distribution/types.FeePool",
		"github.com/cosmos/cosmos-sdk/x/distribution/types.MsgSetWithdrawAddress",
		"github.com/cosmos/cosmos-sdk/x/distribution/types.MsgWithdrawDelegatorReward",
		"github.com/cosmos/cosmos-sdk/x/distribution/types.MsgWithdrawValidatorCommission",
		"github.com/cosmos/cosmos-sdk/x/gov/types.Content",
		"github.com/cosmos/cosmos-sdk/x/gov/types.MsgDeposit",
		"github.com/cosmos/cosmos-sdk/x/gov/types.MsgSubmitProposal",
		"github.com/cosmos/cosmos-sdk/x/gov/types.MsgVote",
		"github.com/cosmos/cosmos-sdk/x/gov/types.SoftwareUpgradeProposal",
		"github.com/cosmos/cosmos-sdk/x/gov/types.TextProposal",
		"github.com/cosmos/cosmos-sdk/x/gov/types.VoteOption",
		"github.com/cosmos/cosmos-sdk/x/params/types.ParamChange",
		"github.com/cosmos/cosmos-sdk/x/params/types.ParameterChangeProposal",
		"github.com/cosmos/cosmos-sdk/x/slashing/types.MsgUnjail",
		"github.com/cosmos/cosmos-sdk/x/staking/types.MsgBeginRedelegate",
		"github.com/cosmos/cosmos-sdk/x/staking/types.MsgCreateValidator",
		"github.com/cosmos/cosmos-sdk/x/staking/types.MsgDelegate",
		"github.com/cosmos/cosmos-sdk/x/staking/types.MsgEditValidator",
		"github.com/cosmos/cosmos-sdk/x/staking/types.MsgUndelegate",
		"github.com/cosmos/cosmos-sdk/x/supply/exported.ModuleAccountI",
		"github.com/cosmos/cosmos-sdk/x/supply/exported.SupplyI",
		"github.com/cosmos/cosmos-sdk/x/supply/internal/types.ModuleAccount",
		"github.com/cosmos/cosmos-sdk/x/supply/internal/types.Supply",
		"github.com/tendermint/tendermint/crypto.PrivKey",
		"github.com/tendermint/tendermint/crypto.PubKey",
		"github.com/tendermint/tendermint/crypto/ed25519.PrivKeyEd25519",
		"github.com/tendermint/tendermint/crypto/ed25519.PubKeyEd25519",
		"github.com/tendermint/tendermint/crypto/multisig.PubKeyMultisigThreshold",
		"github.com/tendermint/tendermint/crypto/secp256k1.PrivKeySecp256k1",
		"github.com/tendermint/tendermint/crypto/secp256k1.PubKeySecp256k1",
		"github.com/tendermint/tendermint/types.SignedMsgType",
		"github.com/tendermint/tendermint/types.Vote",
	}
} // end of GetSupportList
