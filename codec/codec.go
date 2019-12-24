//nolint
package codec

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func codonEncodeBool(w io.Writer, v bool) error {
	slice := []byte{0}
	if v {
		slice = []byte{1}
	}
	_, err := w.Write(slice)
	return err
}
func codonEncodeVarint(w io.Writer, v int64) error {
	var buf [10]byte
	n := binary.PutVarint(buf[:], v)
	_, err := w.Write(buf[0:n])
	return err
}
func codonEncodeInt8(w io.Writer, v int8) error {
	_, err := w.Write([]byte{byte(v)})
	return err
}
func codonEncodeInt16(w io.Writer, v int16) error {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(v))
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeUvarint(w io.Writer, v uint64) error {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], v)
	_, err := w.Write(buf[0:n])
	return err
}
func codonEncodeUint8(w io.Writer, v uint8) error {
	_, err := w.Write([]byte{byte(v)})
	return err
}
func codonEncodeUint16(w io.Writer, v uint16) error {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeFloat32(w io.Writer, v float32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(v))
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeFloat64(w io.Writer, v float64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(v))
	_, err := w.Write(buf[:])
	return err
}
func codonEncodeByteSlice(w io.Writer, v []byte) error {
	err := codonEncodeVarint(w, int64(len(v)))
	if err != nil {
		return err
	}
	_, err = w.Write(v)
	return err
}
func codonEncodeString(w io.Writer, v string) error {
	return codonEncodeByteSlice(w, []byte(v))
}
func codonDecodeBool(bz []byte, n *int, err *error) bool {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return false
	}
	*n = 1
	*err = nil
	return bz[0] != 0
}
func codonDecodeInt(bz []byte, m *int, err *error) int {
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
	return int(i)
}
func codonDecodeInt8(bz []byte, n *int, err *error) int8 {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*err = nil
	*n = 1
	return int8(bz[0])
}
func codonDecodeInt16(bz []byte, n *int, err *error) int16 {
	if len(bz) < 2 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 2
	*err = nil
	return int16(binary.LittleEndian.Uint16(bz[:2]))
}
func codonDecodeInt32(bz []byte, n *int, err *error) int32 {
	i := codonDecodeInt64(bz, n, err)
	return int32(i)
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
	i := codonDecodeUint64(bz, n, err)
	return uint(i)
}
func codonDecodeUint8(bz []byte, n *int, err *error) uint8 {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 1
	*err = nil
	return uint8(bz[0])
}
func codonDecodeUint16(bz []byte, n *int, err *error) uint16 {
	if len(bz) < 2 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 2
	*err = nil
	return uint16(binary.LittleEndian.Uint16(bz[:2]))
}
func codonDecodeUint32(bz []byte, n *int, err *error) uint32 {
	i := codonDecodeUint64(bz, n, err)
	return uint32(i)
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
func codonDecodeFloat64(bz []byte, n *int, err *error) float64 {
	if len(bz) < 8 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 8
	*err = nil
	i := binary.LittleEndian.Uint64(bz[:8])
	return math.Float64frombits(i)
}
func codonDecodeFloat32(bz []byte, n *int, err *error) float32 {
	if len(bz) < 4 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 4
	*err = nil
	i := binary.LittleEndian.Uint32(bz[:4])
	return math.Float32frombits(i)
}
func codonGetByteSlice(bz []byte, length int) ([]byte, int, error) {
	if len(bz) < length {
		return nil, 0, errors.New("Not enough bytes to read")
	}
	return bz[:length], length, nil
}
func codonDecodeString(bz []byte, n *int, err *error) string {
	var m int
	length := codonDecodeInt64(bz, &m, err)
	if *err != nil {
		return ""
	}
	var bs []byte
	var l int
	bs, l, *err = codonGetByteSlice(bz[m:], int(length))
	*n = m + l
	return string(bs)
}

func EncodeTime(w io.Writer, t time.Time) error {
	t = t.UTC()
	sec := t.Unix()
	var buf [10]byte
	n := binary.PutVarint(buf[:], sec)
	_, err := w.Write(buf[0:n])
	if err != nil {
		return err
	}

	nanosec := t.Nanosecond()
	n = binary.PutVarint(buf[:], int64(nanosec))
	_, err = w.Write(buf[0:n])
	if err != nil {
		return err
	}
	return nil
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

func RandTime(r RandSrc) time.Time {
	return time.Unix(r.GetInt64(), r.GetInt64()).UTC()
}

func EncodeInt(w io.Writer, v sdk.Int) error {
	s, err := v.MarshalAmino()
	if err != nil {
		return err
	}
	return codonEncodeString(w, s)
}

func DecodeInt(bz []byte) (sdk.Int, int, error) {
	v := sdk.ZeroInt()
	var n int
	var err error
	s := codonDecodeString(bz, &n, &err)
	if err != nil {
		return v, n, err
	}

	err = (&v).UnmarshalAmino(s)
	if err != nil {
		return v, n, err
	}

	return v, n, nil
}

func RandInt(r RandSrc) sdk.Int {
	res := sdk.NewInt(r.GetInt64())
	count := int(r.GetInt64() % 3)
	for i := 0; i < count; i++ {
		res = res.MulRaw(r.GetInt64())
	}
	return res
}

func EncodeDec(w io.Writer, v sdk.Dec) error {
	s, err := v.MarshalAmino()
	if err != nil {
		return err
	}
	return codonEncodeString(w, s)
}

func DecodeDec(bz []byte) (sdk.Dec, int, error) {
	v := sdk.ZeroDec()
	var n int
	var err error
	s := codonDecodeString(bz, &n, &err)
	if err != nil {
		return v, n, err
	}

	err = (&v).UnmarshalAmino(s)
	if err != nil {
		return v, n, err
	}

	return v, n, nil
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

// Non-Interface
func EncodeDuplicateVoteEvidence(w io.Writer, v DuplicateVoteEvidence) error {
	// codon version: 1
	var err error
	err = EncodePubKey(w, v.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeUint8(w, uint8(v.VoteA.Type))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteA.Height))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteA.Round))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteA.BlockID.Hash[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteA.BlockID.PartsHeader.Total))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteA.BlockID.PartsHeader.Hash[:])
	if err != nil {
		return err
	}
	// end of v.VoteA.BlockID.PartsHeader
	// end of v.VoteA.BlockID
	err = EncodeTime(w, v.VoteA.Timestamp)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteA.ValidatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteA.ValidatorIndex))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteA.Signature[:])
	if err != nil {
		return err
	}
	// end of v.VoteA
	err = codonEncodeUint8(w, uint8(v.VoteB.Type))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteB.Height))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteB.Round))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteB.BlockID.Hash[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteB.BlockID.PartsHeader.Total))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteB.BlockID.PartsHeader.Hash[:])
	if err != nil {
		return err
	}
	// end of v.VoteB.BlockID.PartsHeader
	// end of v.VoteB.BlockID
	err = EncodeTime(w, v.VoteB.Timestamp)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteB.ValidatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.VoteB.ValidatorIndex))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.VoteB.Signature[:])
	if err != nil {
		return err
	}
	// end of v.VoteB
	return nil
} //End of EncodeDuplicateVoteEvidence

func DecodeDuplicateVoteEvidence(bz []byte) (DuplicateVoteEvidence, int, error) {
	// codon version: 1
	var err error
	var length int
	var v DuplicateVoteEvidence
	var n int
	var total int
	v.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	v.VoteA = &Vote{}
	v.VoteA.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.Height = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.Round = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.VoteA.BlockID.PartsHeader
	// end of v.VoteA.BlockID
	v.VoteA.Timestamp, n, err = DecodeTime(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteA.Signature, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.VoteA
	v.VoteB = &Vote{}
	v.VoteB.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.Height = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.Round = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.VoteB.BlockID.PartsHeader
	// end of v.VoteB.BlockID
	v.VoteB.Timestamp, n, err = DecodeTime(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.VoteB.Signature, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.VoteB
	return v, total, nil
} //End of DecodeDuplicateVoteEvidence

func RandDuplicateVoteEvidence(r RandSrc) DuplicateVoteEvidence {
	// codon version: 1
	var length int
	var v DuplicateVoteEvidence
	v.PubKey = RandPubKey(r) // interface_decode
	v.VoteA = &Vote{}
	v.VoteA.Type = SignedMsgType(r.GetUint8())
	v.VoteA.Height = r.GetInt64()
	v.VoteA.Round = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteA.BlockID.Hash = r.GetBytes(length)
	v.VoteA.BlockID.PartsHeader.Total = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteA.BlockID.PartsHeader.Hash = r.GetBytes(length)
	// end of v.VoteA.BlockID.PartsHeader
	// end of v.VoteA.BlockID
	v.VoteA.Timestamp = RandTime(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteA.ValidatorAddress = r.GetBytes(length)
	v.VoteA.ValidatorIndex = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteA.Signature = r.GetBytes(length)
	// end of v.VoteA
	v.VoteB = &Vote{}
	v.VoteB.Type = SignedMsgType(r.GetUint8())
	v.VoteB.Height = r.GetInt64()
	v.VoteB.Round = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteB.BlockID.Hash = r.GetBytes(length)
	v.VoteB.BlockID.PartsHeader.Total = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteB.BlockID.PartsHeader.Hash = r.GetBytes(length)
	// end of v.VoteB.BlockID.PartsHeader
	// end of v.VoteB.BlockID
	v.VoteB.Timestamp = RandTime(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteB.ValidatorAddress = r.GetBytes(length)
	v.VoteB.ValidatorIndex = r.GetInt()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.VoteB.Signature = r.GetBytes(length)
	// end of v.VoteB
	return v
} //End of RandDuplicateVoteEvidence

// Non-Interface
func EncodePrivKeyEd25519(w io.Writer, v PrivKeyEd25519) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodePrivKeyEd25519

func DecodePrivKeyEd25519(bz []byte) (PrivKeyEd25519, int, error) {
	// codon version: 1
	var err error
	var length int
	var v PrivKeyEd25519
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodePrivKeyEd25519

func RandPrivKeyEd25519(r RandSrc) PrivKeyEd25519 {
	// codon version: 1
	var length int
	var v PrivKeyEd25519
	length = 64
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPrivKeyEd25519

// Non-Interface
func EncodePrivKeySecp256k1(w io.Writer, v PrivKeySecp256k1) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodePrivKeySecp256k1

func DecodePrivKeySecp256k1(bz []byte) (PrivKeySecp256k1, int, error) {
	// codon version: 1
	var err error
	var length int
	var v PrivKeySecp256k1
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodePrivKeySecp256k1

func RandPrivKeySecp256k1(r RandSrc) PrivKeySecp256k1 {
	// codon version: 1
	var length int
	var v PrivKeySecp256k1
	length = 32
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPrivKeySecp256k1

// Non-Interface
func EncodePubKeyEd25519(w io.Writer, v PubKeyEd25519) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodePubKeyEd25519

func DecodePubKeyEd25519(bz []byte) (PubKeyEd25519, int, error) {
	// codon version: 1
	var err error
	var length int
	var v PubKeyEd25519
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodePubKeyEd25519

func RandPubKeyEd25519(r RandSrc) PubKeyEd25519 {
	// codon version: 1
	var length int
	var v PubKeyEd25519
	length = 32
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPubKeyEd25519

// Non-Interface
func EncodePubKeySecp256k1(w io.Writer, v PubKeySecp256k1) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodePubKeySecp256k1

func DecodePubKeySecp256k1(bz []byte) (PubKeySecp256k1, int, error) {
	// codon version: 1
	var err error
	var length int
	var v PubKeySecp256k1
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodePubKeySecp256k1

func RandPubKeySecp256k1(r RandSrc) PubKeySecp256k1 {
	// codon version: 1
	var length int
	var v PubKeySecp256k1
	length = 33
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //array of uint8
		v[_0] = r.GetUint8()
	}
	return v
} //End of RandPubKeySecp256k1

// Non-Interface
func EncodePubKeyMultisigThreshold(w io.Writer, v PubKeyMultisigThreshold) error {
	// codon version: 1
	var err error
	err = codonEncodeUvarint(w, uint64(v.K))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.PubKeys)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.PubKeys); _0++ {
		err = EncodePubKey(w, v.PubKeys[_0])
		if err != nil {
			return err
		} // interface_encode
	}
	return nil
} //End of EncodePubKeyMultisigThreshold

func DecodePubKeyMultisigThreshold(bz []byte) (PubKeyMultisigThreshold, int, error) {
	// codon version: 1
	var err error
	var length int
	var v PubKeyMultisigThreshold
	var n int
	var total int
	v.K = uint(codonDecodeUint(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.PubKeys = make([]PubKey, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		v.PubKeys[_0], n, err = DecodePubKey(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodePubKeyMultisigThreshold

func RandPubKeyMultisigThreshold(r RandSrc) PubKeyMultisigThreshold {
	// codon version: 1
	var length int
	var v PubKeyMultisigThreshold
	v.K = r.GetUint()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.PubKeys = make([]PubKey, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		v.PubKeys[_0] = RandPubKey(r)
	}
	return v
} //End of RandPubKeyMultisigThreshold

// Non-Interface
func EncodeSignedMsgType(w io.Writer, v SignedMsgType) error {
	// codon version: 1
	var err error
	err = codonEncodeUint8(w, uint8(v))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeSignedMsgType

func DecodeSignedMsgType(bz []byte) (SignedMsgType, int, error) {
	// codon version: 1
	var err error
	var v SignedMsgType
	var n int
	var total int
	v = SignedMsgType(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeSignedMsgType

func RandSignedMsgType(r RandSrc) SignedMsgType {
	// codon version: 1
	var v SignedMsgType
	v = SignedMsgType(r.GetUint8())
	return v
} //End of RandSignedMsgType

// Non-Interface
func EncodeVoteOption(w io.Writer, v VoteOption) error {
	// codon version: 1
	var err error
	err = codonEncodeUint8(w, uint8(v))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeVoteOption

func DecodeVoteOption(bz []byte) (VoteOption, int, error) {
	// codon version: 1
	var err error
	var v VoteOption
	var n int
	var total int
	v = VoteOption(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeVoteOption

func RandVoteOption(r RandSrc) VoteOption {
	// codon version: 1
	var v VoteOption
	v = VoteOption(r.GetUint8())
	return v
} //End of RandVoteOption

// Non-Interface
func EncodeVote(w io.Writer, v Vote) error {
	// codon version: 1
	var err error
	err = codonEncodeUint8(w, uint8(v.Type))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Height))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Round))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.BlockID.Hash[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.BlockID.PartsHeader.Total))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.BlockID.PartsHeader.Hash[:])
	if err != nil {
		return err
	}
	// end of v.BlockID.PartsHeader
	// end of v.BlockID
	err = EncodeTime(w, v.Timestamp)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ValidatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.ValidatorIndex))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Signature[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeVote

func DecodeVote(bz []byte) (Vote, int, error) {
	// codon version: 1
	var err error
	var length int
	var v Vote
	var n int
	var total int
	v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Height = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Round = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.BlockID.PartsHeader
	// end of v.BlockID
	v.Timestamp, n, err = DecodeTime(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Signature, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeVote

func RandVote(r RandSrc) Vote {
	// codon version: 1
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

// Non-Interface
func EncodeCoin(w io.Writer, v Coin) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Denom)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Amount)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeCoin

func DecodeCoin(bz []byte) (Coin, int, error) {
	// codon version: 1
	var err error
	var v Coin
	var n int
	var total int
	v.Denom = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeCoin

func RandCoin(r RandSrc) Coin {
	// codon version: 1
	var v Coin
	v.Denom = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = RandInt(r)
	return v
} //End of RandCoin

// Non-Interface
func EncodeLockedCoin(w io.Writer, v LockedCoin) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Coin.Denom)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Coin.Amount)
	if err != nil {
		return err
	}
	// end of v.Coin
	err = codonEncodeVarint(w, int64(v.UnlockTime))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.FromAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Supervisor[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Reward))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeLockedCoin

func DecodeLockedCoin(bz []byte) (LockedCoin, int, error) {
	// codon version: 1
	var err error
	var length int
	var v LockedCoin
	var n int
	var total int
	v.Coin.Denom = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Coin.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Coin
	v.UnlockTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.FromAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Supervisor, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Reward = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeLockedCoin

func RandLockedCoin(r RandSrc) LockedCoin {
	// codon version: 1
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

// Non-Interface
func EncodeStdSignature(w io.Writer, v StdSignature) error {
	// codon version: 1
	var err error
	err = EncodePubKey(w, v.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeByteSlice(w, v.Signature[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeStdSignature

func DecodeStdSignature(bz []byte) (StdSignature, int, error) {
	// codon version: 1
	var err error
	var length int
	var v StdSignature
	var n int
	var total int
	v.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Signature, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeStdSignature

func RandStdSignature(r RandSrc) StdSignature {
	// codon version: 1
	var length int
	var v StdSignature
	v.PubKey = RandPubKey(r) // interface_decode
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Signature = r.GetBytes(length)
	return v
} //End of RandStdSignature

// Non-Interface
func EncodeParamChange(w io.Writer, v ParamChange) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Subspace)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Key)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Subkey)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Value)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeParamChange

func DecodeParamChange(bz []byte) (ParamChange, int, error) {
	// codon version: 1
	var err error
	var v ParamChange
	var n int
	var total int
	v.Subspace = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Key = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Subkey = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Value = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeParamChange

func RandParamChange(r RandSrc) ParamChange {
	// codon version: 1
	var v ParamChange
	v.Subspace = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Key = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Subkey = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Value = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandParamChange

// Non-Interface
func EncodeInput(w io.Writer, v Input) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Coins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Coins); _0++ {
		err = codonEncodeString(w, v.Coins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Coins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Coins[_0]
	}
	return nil
} //End of EncodeInput

func DecodeInput(bz []byte) (Input, int, error) {
	// codon version: 1
	var err error
	var length int
	var v Input
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeInput

func RandInput(r RandSrc) Input {
	// codon version: 1
	var length int
	var v Input
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0] = RandCoin(r)
	}
	return v
} //End of RandInput

// Non-Interface
func EncodeOutput(w io.Writer, v Output) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Coins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Coins); _0++ {
		err = codonEncodeString(w, v.Coins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Coins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Coins[_0]
	}
	return nil
} //End of EncodeOutput

func DecodeOutput(bz []byte) (Output, int, error) {
	// codon version: 1
	var err error
	var length int
	var v Output
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeOutput

func RandOutput(r RandSrc) Output {
	// codon version: 1
	var length int
	var v Output
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0] = RandCoin(r)
	}
	return v
} //End of RandOutput

// Non-Interface
func EncodeAccAddress(w io.Writer, v AccAddress) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeAccAddress

func DecodeAccAddress(bz []byte) (AccAddress, int, error) {
	// codon version: 1
	var err error
	var length int
	var v AccAddress
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeAccAddress

func RandAccAddress(r RandSrc) AccAddress {
	// codon version: 1
	var length int
	var v AccAddress
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v = r.GetBytes(length)
	return v
} //End of RandAccAddress

// Non-Interface
func EncodeCommentRef(w io.Writer, v CommentRef) error {
	// codon version: 1
	var err error
	err = codonEncodeUvarint(w, uint64(v.ID))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.RewardTarget[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.RewardToken)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.RewardAmount))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Attitudes)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Attitudes); _0++ {
		err = codonEncodeVarint(w, int64(v.Attitudes[_0]))
		if err != nil {
			return err
		}
	}
	return nil
} //End of EncodeCommentRef

func DecodeCommentRef(bz []byte) (CommentRef, int, error) {
	// codon version: 1
	var err error
	var length int
	var v CommentRef
	var n int
	var total int
	v.ID = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.RewardTarget, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.RewardToken = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.RewardAmount = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Attitudes = make([]int32, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of int32
		v.Attitudes[_0] = int32(codonDecodeInt32(bz, &n, &err))
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeCommentRef

func RandCommentRef(r RandSrc) CommentRef {
	// codon version: 1
	var length int
	var v CommentRef
	v.ID = r.GetUint64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.RewardTarget = r.GetBytes(length)
	v.RewardToken = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.RewardAmount = r.GetInt64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Attitudes = make([]int32, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of int32
		v.Attitudes[_0] = r.GetInt32()
	}
	return v
} //End of RandCommentRef

// Non-Interface
func EncodeBaseAccount(w io.Writer, v BaseAccount) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Coins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Coins); _0++ {
		err = codonEncodeString(w, v.Coins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Coins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Coins[_0]
	}
	err = EncodePubKey(w, v.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeUvarint(w, uint64(v.AccountNumber))
	if err != nil {
		return err
	}
	err = codonEncodeUvarint(w, uint64(v.Sequence))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeBaseAccount

func DecodeBaseAccount(bz []byte) (BaseAccount, int, error) {
	// codon version: 1
	var err error
	var length int
	var v BaseAccount
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	v.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeBaseAccount

func RandBaseAccount(r RandSrc) BaseAccount {
	// codon version: 1
	var length int
	var v BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Coins[_0] = RandCoin(r)
	}
	v.PubKey = RandPubKey(r) // interface_decode
	v.AccountNumber = r.GetUint64()
	v.Sequence = r.GetUint64()
	return v
} //End of RandBaseAccount

// Non-Interface
func EncodeBaseVestingAccount(w io.Writer, v BaseVestingAccount) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.BaseAccount.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.BaseAccount.Coins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseAccount.Coins); _0++ {
		err = codonEncodeString(w, v.BaseAccount.Coins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseAccount.Coins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseAccount.Coins[_0]
	}
	err = EncodePubKey(w, v.BaseAccount.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeUvarint(w, uint64(v.BaseAccount.AccountNumber))
	if err != nil {
		return err
	}
	err = codonEncodeUvarint(w, uint64(v.BaseAccount.Sequence))
	if err != nil {
		return err
	}
	// end of v.BaseAccount
	err = codonEncodeVarint(w, int64(len(v.OriginalVesting)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.OriginalVesting); _0++ {
		err = codonEncodeString(w, v.OriginalVesting[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.OriginalVesting[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.OriginalVesting[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.DelegatedFree)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.DelegatedFree); _0++ {
		err = codonEncodeString(w, v.DelegatedFree[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.DelegatedFree[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.DelegatedFree[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.DelegatedVesting)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.DelegatedVesting); _0++ {
		err = codonEncodeString(w, v.DelegatedVesting[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.DelegatedVesting[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.DelegatedVesting[_0]
	}
	err = codonEncodeVarint(w, int64(v.EndTime))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeBaseVestingAccount

func DecodeBaseVestingAccount(bz []byte) (BaseVestingAccount, int, error) {
	// codon version: 1
	var err error
	var length int
	var v BaseVestingAccount
	var n int
	var total int
	v.BaseAccount = &BaseAccount{}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseAccount.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseAccount.Coins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.BaseAccount.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	v.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.BaseAccount
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OriginalVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.OriginalVesting[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatedFree = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.DelegatedFree[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatedVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.DelegatedVesting[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.EndTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeBaseVestingAccount

func RandBaseVestingAccount(r RandSrc) BaseVestingAccount {
	// codon version: 1
	var length int
	var v BaseVestingAccount
	v.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OriginalVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.OriginalVesting[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatedFree = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.DelegatedFree[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatedVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.DelegatedVesting[_0] = RandCoin(r)
	}
	v.EndTime = r.GetInt64()
	return v
} //End of RandBaseVestingAccount

// Non-Interface
func EncodeContinuousVestingAccount(w io.Writer, v ContinuousVestingAccount) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.BaseVestingAccount.BaseAccount.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.BaseAccount.Coins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.BaseAccount.Coins); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.BaseAccount.Coins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.BaseAccount.Coins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.BaseAccount.Coins[_0]
	}
	err = EncodePubKey(w, v.BaseVestingAccount.BaseAccount.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeUvarint(w, uint64(v.BaseVestingAccount.BaseAccount.AccountNumber))
	if err != nil {
		return err
	}
	err = codonEncodeUvarint(w, uint64(v.BaseVestingAccount.BaseAccount.Sequence))
	if err != nil {
		return err
	}
	// end of v.BaseVestingAccount.BaseAccount
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.OriginalVesting)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.OriginalVesting); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.OriginalVesting[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.OriginalVesting[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.OriginalVesting[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.DelegatedFree)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedFree); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.DelegatedFree[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.DelegatedFree[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.DelegatedFree[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.DelegatedVesting)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedVesting); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.DelegatedVesting[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.DelegatedVesting[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.DelegatedVesting[_0]
	}
	err = codonEncodeVarint(w, int64(v.BaseVestingAccount.EndTime))
	if err != nil {
		return err
	}
	// end of v.BaseVestingAccount
	err = codonEncodeVarint(w, int64(v.StartTime))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeContinuousVestingAccount

func DecodeContinuousVestingAccount(bz []byte) (ContinuousVestingAccount, int, error) {
	// codon version: 1
	var err error
	var length int
	var v ContinuousVestingAccount
	var n int
	var total int
	v.BaseVestingAccount = &BaseVestingAccount{}
	v.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.BaseAccount.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.BaseAccount.Coins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.BaseVestingAccount.BaseAccount.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	v.BaseVestingAccount.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.BaseVestingAccount.BaseAccount
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.OriginalVesting[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedFree[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedVesting[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.BaseVestingAccount.EndTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.BaseVestingAccount
	v.StartTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeContinuousVestingAccount

func RandContinuousVestingAccount(r RandSrc) ContinuousVestingAccount {
	// codon version: 1
	var length int
	var v ContinuousVestingAccount
	v.BaseVestingAccount = &BaseVestingAccount{}
	v.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseVestingAccount.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseVestingAccount.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseVestingAccount.BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.OriginalVesting[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedFree[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedVesting[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.EndTime = r.GetInt64()
	// end of v.BaseVestingAccount
	v.StartTime = r.GetInt64()
	return v
} //End of RandContinuousVestingAccount

// Non-Interface
func EncodeDelayedVestingAccount(w io.Writer, v DelayedVestingAccount) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.BaseVestingAccount.BaseAccount.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.BaseAccount.Coins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.BaseAccount.Coins); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.BaseAccount.Coins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.BaseAccount.Coins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.BaseAccount.Coins[_0]
	}
	err = EncodePubKey(w, v.BaseVestingAccount.BaseAccount.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeUvarint(w, uint64(v.BaseVestingAccount.BaseAccount.AccountNumber))
	if err != nil {
		return err
	}
	err = codonEncodeUvarint(w, uint64(v.BaseVestingAccount.BaseAccount.Sequence))
	if err != nil {
		return err
	}
	// end of v.BaseVestingAccount.BaseAccount
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.OriginalVesting)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.OriginalVesting); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.OriginalVesting[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.OriginalVesting[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.OriginalVesting[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.DelegatedFree)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedFree); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.DelegatedFree[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.DelegatedFree[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.DelegatedFree[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.BaseVestingAccount.DelegatedVesting)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseVestingAccount.DelegatedVesting); _0++ {
		err = codonEncodeString(w, v.BaseVestingAccount.DelegatedVesting[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseVestingAccount.DelegatedVesting[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseVestingAccount.DelegatedVesting[_0]
	}
	err = codonEncodeVarint(w, int64(v.BaseVestingAccount.EndTime))
	if err != nil {
		return err
	}
	// end of v.BaseVestingAccount
	return nil
} //End of EncodeDelayedVestingAccount

func DecodeDelayedVestingAccount(bz []byte) (DelayedVestingAccount, int, error) {
	// codon version: 1
	var err error
	var length int
	var v DelayedVestingAccount
	var n int
	var total int
	v.BaseVestingAccount = &BaseVestingAccount{}
	v.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.BaseAccount.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.BaseAccount.Coins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.BaseVestingAccount.BaseAccount.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	v.BaseVestingAccount.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.BaseVestingAccount.BaseAccount
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.OriginalVesting[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedFree[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedVesting[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.BaseVestingAccount.EndTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.BaseVestingAccount
	return v, total, nil
} //End of DecodeDelayedVestingAccount

func RandDelayedVestingAccount(r RandSrc) DelayedVestingAccount {
	// codon version: 1
	var length int
	var v DelayedVestingAccount
	v.BaseVestingAccount = &BaseVestingAccount{}
	v.BaseVestingAccount.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseVestingAccount.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseVestingAccount.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseVestingAccount.BaseAccount
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.OriginalVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.OriginalVesting[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.DelegatedFree = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedFree[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseVestingAccount.DelegatedVesting = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseVestingAccount.DelegatedVesting[_0] = RandCoin(r)
	}
	v.BaseVestingAccount.EndTime = r.GetInt64()
	// end of v.BaseVestingAccount
	return v
} //End of RandDelayedVestingAccount

// Non-Interface
func EncodeModuleAccount(w io.Writer, v ModuleAccount) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.BaseAccount.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.BaseAccount.Coins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.BaseAccount.Coins); _0++ {
		err = codonEncodeString(w, v.BaseAccount.Coins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.BaseAccount.Coins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.BaseAccount.Coins[_0]
	}
	err = EncodePubKey(w, v.BaseAccount.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeUvarint(w, uint64(v.BaseAccount.AccountNumber))
	if err != nil {
		return err
	}
	err = codonEncodeUvarint(w, uint64(v.BaseAccount.Sequence))
	if err != nil {
		return err
	}
	// end of v.BaseAccount
	err = codonEncodeString(w, v.Name)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Permissions)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Permissions); _0++ {
		err = codonEncodeString(w, v.Permissions[_0])
		if err != nil {
			return err
		}
	}
	return nil
} //End of EncodeModuleAccount

func DecodeModuleAccount(bz []byte) (ModuleAccount, int, error) {
	// codon version: 1
	var err error
	var length int
	var v ModuleAccount
	var n int
	var total int
	v.BaseAccount = &BaseAccount{}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseAccount.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseAccount.Coins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.BaseAccount.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	v.BaseAccount.AccountNumber = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.BaseAccount.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.BaseAccount
	v.Name = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Permissions = make([]string, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of string
		v.Permissions[_0] = string(codonDecodeString(bz, &n, &err))
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeModuleAccount

func RandModuleAccount(r RandSrc) ModuleAccount {
	// codon version: 1
	var length int
	var v ModuleAccount
	v.BaseAccount = &BaseAccount{}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseAccount.Address = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.BaseAccount.Coins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.BaseAccount.Coins[_0] = RandCoin(r)
	}
	v.BaseAccount.PubKey = RandPubKey(r) // interface_decode
	v.BaseAccount.AccountNumber = r.GetUint64()
	v.BaseAccount.Sequence = r.GetUint64()
	// end of v.BaseAccount
	v.Name = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Permissions = make([]string, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of string
		v.Permissions[_0] = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	}
	return v
} //End of RandModuleAccount

// Non-Interface
func EncodeStdTx(w io.Writer, v StdTx) error {
	// codon version: 1
	var err error
	err = codonEncodeVarint(w, int64(len(v.Msgs)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Msgs); _0++ {
		err = EncodeMsg(w, v.Msgs[_0])
		if err != nil {
			return err
		} // interface_encode
	}
	err = codonEncodeVarint(w, int64(len(v.Fee.Amount)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Fee.Amount); _0++ {
		err = codonEncodeString(w, v.Fee.Amount[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Fee.Amount[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Fee.Amount[_0]
	}
	err = codonEncodeUvarint(w, uint64(v.Fee.Gas))
	if err != nil {
		return err
	}
	// end of v.Fee
	err = codonEncodeVarint(w, int64(len(v.Signatures)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Signatures); _0++ {
		err = EncodePubKey(w, v.Signatures[_0].PubKey)
		if err != nil {
			return err
		} // interface_encode
		err = codonEncodeByteSlice(w, v.Signatures[_0].Signature[:])
		if err != nil {
			return err
		}
		// end of v.Signatures[_0]
	}
	err = codonEncodeString(w, v.Memo)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeStdTx

func DecodeStdTx(bz []byte) (StdTx, int, error) {
	// codon version: 1
	var err error
	var length int
	var v StdTx
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Msgs = make([]Msg, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		v.Msgs[_0], n, err = DecodeMsg(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Fee.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Fee.Amount[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.Fee.Gas = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Fee
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Signatures = make([]StdSignature, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Signatures[_0], n, err = DecodeStdSignature(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.Memo = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeStdTx

func RandStdTx(r RandSrc) StdTx {
	// codon version: 1
	var length int
	var v StdTx
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Msgs = make([]Msg, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of interface
		v.Msgs[_0] = RandMsg(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Fee.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Fee.Amount[_0] = RandCoin(r)
	}
	v.Fee.Gas = r.GetUint64()
	// end of v.Fee
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Signatures = make([]StdSignature, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Signatures[_0] = RandStdSignature(r)
	}
	v.Memo = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandStdTx

// Non-Interface
func EncodeMsgBeginRedelegate(w io.Writer, v MsgBeginRedelegate) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.DelegatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ValidatorSrcAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ValidatorDstAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Amount.Denom)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Amount.Amount)
	if err != nil {
		return err
	}
	// end of v.Amount
	return nil
} //End of EncodeMsgBeginRedelegate

func DecodeMsgBeginRedelegate(bz []byte) (MsgBeginRedelegate, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgBeginRedelegate
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorSrcAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorDstAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount.Denom = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Amount
	return v, total, nil
} //End of DecodeMsgBeginRedelegate

func RandMsgBeginRedelegate(r RandSrc) MsgBeginRedelegate {
	// codon version: 1
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

// Non-Interface
func EncodeMsgCreateValidator(w io.Writer, v MsgCreateValidator) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Description.Moniker)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description.Identity)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description.Website)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description.Details)
	if err != nil {
		return err
	}
	// end of v.Description
	err = EncodeDec(w, v.Commission.Rate)
	if err != nil {
		return err
	}
	err = EncodeDec(w, v.Commission.MaxRate)
	if err != nil {
		return err
	}
	err = EncodeDec(w, v.Commission.MaxChangeRate)
	if err != nil {
		return err
	}
	// end of v.Commission
	err = EncodeInt(w, v.MinSelfDelegation)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.DelegatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ValidatorAddress[:])
	if err != nil {
		return err
	}
	err = EncodePubKey(w, v.PubKey)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeString(w, v.Value.Denom)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Value.Amount)
	if err != nil {
		return err
	}
	// end of v.Value
	return nil
} //End of EncodeMsgCreateValidator

func DecodeMsgCreateValidator(bz []byte) (MsgCreateValidator, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgCreateValidator
	var n int
	var total int
	v.Description.Moniker = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description.Identity = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description.Website = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description.Details = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Description
	v.Commission.Rate, n, err = DecodeDec(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Commission.MaxRate, n, err = DecodeDec(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Commission.MaxChangeRate, n, err = DecodeDec(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Commission
	v.MinSelfDelegation, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.PubKey, n, err = DecodePubKey(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	v.Value.Denom = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Value.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Value
	return v, total, nil
} //End of DecodeMsgCreateValidator

func RandMsgCreateValidator(r RandSrc) MsgCreateValidator {
	// codon version: 1
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

// Non-Interface
func EncodeMsgDelegate(w io.Writer, v MsgDelegate) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.DelegatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ValidatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Amount.Denom)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Amount.Amount)
	if err != nil {
		return err
	}
	// end of v.Amount
	return nil
} //End of EncodeMsgDelegate

func DecodeMsgDelegate(bz []byte) (MsgDelegate, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgDelegate
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount.Denom = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Amount
	return v, total, nil
} //End of DecodeMsgDelegate

func RandMsgDelegate(r RandSrc) MsgDelegate {
	// codon version: 1
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

// Non-Interface
func EncodeMsgEditValidator(w io.Writer, v MsgEditValidator) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Description.Moniker)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description.Identity)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description.Website)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description.Details)
	if err != nil {
		return err
	}
	// end of v.Description
	err = codonEncodeByteSlice(w, v.ValidatorAddress[:])
	if err != nil {
		return err
	}
	err = EncodeDec(w, *(v.CommissionRate))
	if err != nil {
		return err
	}
	err = EncodeInt(w, *(v.MinSelfDelegation))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgEditValidator

func DecodeMsgEditValidator(bz []byte) (MsgEditValidator, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgEditValidator
	var n int
	var total int
	v.Description.Moniker = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description.Identity = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description.Website = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description.Details = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Description
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.CommissionRate = &sdk.Dec{}
	*(v.CommissionRate), n, err = DecodeDec(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.MinSelfDelegation = &sdk.Int{}
	*(v.MinSelfDelegation), n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgEditValidator

func RandMsgEditValidator(r RandSrc) MsgEditValidator {
	// codon version: 1
	var length int
	var v MsgEditValidator
	v.Description.Moniker = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Identity = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Website = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description.Details = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	// end of v.Description
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	v.CommissionRate = &sdk.Dec{}
	*(v.CommissionRate) = RandDec(r)
	v.MinSelfDelegation = &sdk.Int{}
	*(v.MinSelfDelegation) = RandInt(r)
	return v
} //End of RandMsgEditValidator

// Non-Interface
func EncodeMsgSetWithdrawAddress(w io.Writer, v MsgSetWithdrawAddress) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.DelegatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.WithdrawAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgSetWithdrawAddress

func DecodeMsgSetWithdrawAddress(bz []byte) (MsgSetWithdrawAddress, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgSetWithdrawAddress
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.WithdrawAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgSetWithdrawAddress

func RandMsgSetWithdrawAddress(r RandSrc) MsgSetWithdrawAddress {
	// codon version: 1
	var length int
	var v MsgSetWithdrawAddress
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.WithdrawAddress = r.GetBytes(length)
	return v
} //End of RandMsgSetWithdrawAddress

// Non-Interface
func EncodeMsgUndelegate(w io.Writer, v MsgUndelegate) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.DelegatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ValidatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Amount.Denom)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Amount.Amount)
	if err != nil {
		return err
	}
	// end of v.Amount
	return nil
} //End of EncodeMsgUndelegate

func DecodeMsgUndelegate(bz []byte) (MsgUndelegate, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgUndelegate
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount.Denom = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	// end of v.Amount
	return v, total, nil
} //End of DecodeMsgUndelegate

func RandMsgUndelegate(r RandSrc) MsgUndelegate {
	// codon version: 1
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

// Non-Interface
func EncodeMsgUnjail(w io.Writer, v MsgUnjail) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.ValidatorAddr[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgUnjail

func DecodeMsgUnjail(bz []byte) (MsgUnjail, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgUnjail
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddr, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgUnjail

func RandMsgUnjail(r RandSrc) MsgUnjail {
	// codon version: 1
	var length int
	var v MsgUnjail
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddr = r.GetBytes(length)
	return v
} //End of RandMsgUnjail

// Non-Interface
func EncodeMsgWithdrawDelegatorReward(w io.Writer, v MsgWithdrawDelegatorReward) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.DelegatorAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ValidatorAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgWithdrawDelegatorReward

func DecodeMsgWithdrawDelegatorReward(bz []byte) (MsgWithdrawDelegatorReward, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgWithdrawDelegatorReward
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DelegatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgWithdrawDelegatorReward

func RandMsgWithdrawDelegatorReward(r RandSrc) MsgWithdrawDelegatorReward {
	// codon version: 1
	var length int
	var v MsgWithdrawDelegatorReward
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.DelegatorAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	return v
} //End of RandMsgWithdrawDelegatorReward

// Non-Interface
func EncodeMsgWithdrawValidatorCommission(w io.Writer, v MsgWithdrawValidatorCommission) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.ValidatorAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgWithdrawValidatorCommission

func DecodeMsgWithdrawValidatorCommission(bz []byte) (MsgWithdrawValidatorCommission, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgWithdrawValidatorCommission
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgWithdrawValidatorCommission

func RandMsgWithdrawValidatorCommission(r RandSrc) MsgWithdrawValidatorCommission {
	// codon version: 1
	var length int
	var v MsgWithdrawValidatorCommission
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ValidatorAddress = r.GetBytes(length)
	return v
} //End of RandMsgWithdrawValidatorCommission

// Non-Interface
func EncodeMsgDeposit(w io.Writer, v MsgDeposit) error {
	// codon version: 1
	var err error
	err = codonEncodeUvarint(w, uint64(v.ProposalID))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Depositor[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Amount)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Amount); _0++ {
		err = codonEncodeString(w, v.Amount[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Amount[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Amount[_0]
	}
	return nil
} //End of EncodeMsgDeposit

func DecodeMsgDeposit(bz []byte) (MsgDeposit, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgDeposit
	var n int
	var total int
	v.ProposalID = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Depositor, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgDeposit

func RandMsgDeposit(r RandSrc) MsgDeposit {
	// codon version: 1
	var length int
	var v MsgDeposit
	v.ProposalID = r.GetUint64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Depositor = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandMsgDeposit

// Non-Interface
func EncodeMsgSubmitProposal(w io.Writer, v MsgSubmitProposal) error {
	// codon version: 1
	var err error
	err = EncodeContent(w, v.Content)
	if err != nil {
		return err
	} // interface_encode
	err = codonEncodeVarint(w, int64(len(v.InitialDeposit)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.InitialDeposit); _0++ {
		err = codonEncodeString(w, v.InitialDeposit[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.InitialDeposit[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.InitialDeposit[_0]
	}
	err = codonEncodeByteSlice(w, v.Proposer[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgSubmitProposal

func DecodeMsgSubmitProposal(bz []byte) (MsgSubmitProposal, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgSubmitProposal
	var n int
	var total int
	v.Content, n, err = DecodeContent(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n // interface_decode
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.InitialDeposit = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.InitialDeposit[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Proposer, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgSubmitProposal

func RandMsgSubmitProposal(r RandSrc) MsgSubmitProposal {
	// codon version: 1
	var length int
	var v MsgSubmitProposal
	v.Content = RandContent(r) // interface_decode
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.InitialDeposit = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.InitialDeposit[_0] = RandCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Proposer = r.GetBytes(length)
	return v
} //End of RandMsgSubmitProposal

// Non-Interface
func EncodeMsgVote(w io.Writer, v MsgVote) error {
	// codon version: 1
	var err error
	err = codonEncodeUvarint(w, uint64(v.ProposalID))
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Voter[:])
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, uint8(v.Option))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgVote

func DecodeMsgVote(bz []byte) (MsgVote, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgVote
	var n int
	var total int
	v.ProposalID = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Voter, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Option = VoteOption(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgVote

func RandMsgVote(r RandSrc) MsgVote {
	// codon version: 1
	var length int
	var v MsgVote
	v.ProposalID = r.GetUint64()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Voter = r.GetBytes(length)
	v.Option = VoteOption(r.GetUint8())
	return v
} //End of RandMsgVote

// Non-Interface
func EncodeParameterChangeProposal(w io.Writer, v ParameterChangeProposal) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Title)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Changes)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Changes); _0++ {
		err = codonEncodeString(w, v.Changes[_0].Subspace)
		if err != nil {
			return err
		}
		err = codonEncodeString(w, v.Changes[_0].Key)
		if err != nil {
			return err
		}
		err = codonEncodeString(w, v.Changes[_0].Subkey)
		if err != nil {
			return err
		}
		err = codonEncodeString(w, v.Changes[_0].Value)
		if err != nil {
			return err
		}
		// end of v.Changes[_0]
	}
	return nil
} //End of EncodeParameterChangeProposal

func DecodeParameterChangeProposal(bz []byte) (ParameterChangeProposal, int, error) {
	// codon version: 1
	var err error
	var length int
	var v ParameterChangeProposal
	var n int
	var total int
	v.Title = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Changes = make([]ParamChange, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Changes[_0], n, err = DecodeParamChange(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeParameterChangeProposal

func RandParameterChangeProposal(r RandSrc) ParameterChangeProposal {
	// codon version: 1
	var length int
	var v ParameterChangeProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Changes = make([]ParamChange, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Changes[_0] = RandParamChange(r)
	}
	return v
} //End of RandParameterChangeProposal

// Non-Interface
func EncodeSoftwareUpgradeProposal(w io.Writer, v SoftwareUpgradeProposal) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Title)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeSoftwareUpgradeProposal

func DecodeSoftwareUpgradeProposal(bz []byte) (SoftwareUpgradeProposal, int, error) {
	// codon version: 1
	var err error
	var v SoftwareUpgradeProposal
	var n int
	var total int
	v.Title = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeSoftwareUpgradeProposal

func RandSoftwareUpgradeProposal(r RandSrc) SoftwareUpgradeProposal {
	// codon version: 1
	var v SoftwareUpgradeProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandSoftwareUpgradeProposal

// Non-Interface
func EncodeTextProposal(w io.Writer, v TextProposal) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Title)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeTextProposal

func DecodeTextProposal(bz []byte) (TextProposal, int, error) {
	// codon version: 1
	var err error
	var v TextProposal
	var n int
	var total int
	v.Title = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeTextProposal

func RandTextProposal(r RandSrc) TextProposal {
	// codon version: 1
	var v TextProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandTextProposal

// Non-Interface
func EncodeCommunityPoolSpendProposal(w io.Writer, v CommunityPoolSpendProposal) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Title)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Recipient[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Amount)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Amount); _0++ {
		err = codonEncodeString(w, v.Amount[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Amount[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Amount[_0]
	}
	return nil
} //End of EncodeCommunityPoolSpendProposal

func DecodeCommunityPoolSpendProposal(bz []byte) (CommunityPoolSpendProposal, int, error) {
	// codon version: 1
	var err error
	var length int
	var v CommunityPoolSpendProposal
	var n int
	var total int
	v.Title = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Recipient, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeCommunityPoolSpendProposal

func RandCommunityPoolSpendProposal(r RandSrc) CommunityPoolSpendProposal {
	// codon version: 1
	var length int
	var v CommunityPoolSpendProposal
	v.Title = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Description = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Recipient = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandCommunityPoolSpendProposal

// Non-Interface
func EncodeMsgMultiSend(w io.Writer, v MsgMultiSend) error {
	// codon version: 1
	var err error
	err = codonEncodeVarint(w, int64(len(v.Inputs)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Inputs); _0++ {
		err = codonEncodeByteSlice(w, v.Inputs[_0].Address[:])
		if err != nil {
			return err
		}
		err = codonEncodeVarint(w, int64(len(v.Inputs[_0].Coins)))
		if err != nil {
			return err
		}
		for _1 := 0; _1 < len(v.Inputs[_0].Coins); _1++ {
			err = codonEncodeString(w, v.Inputs[_0].Coins[_1].Denom)
			if err != nil {
				return err
			}
			err = EncodeInt(w, v.Inputs[_0].Coins[_1].Amount)
			if err != nil {
				return err
			}
			// end of v.Inputs[_0].Coins[_1]
		}
		// end of v.Inputs[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.Outputs)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Outputs); _0++ {
		err = codonEncodeByteSlice(w, v.Outputs[_0].Address[:])
		if err != nil {
			return err
		}
		err = codonEncodeVarint(w, int64(len(v.Outputs[_0].Coins)))
		if err != nil {
			return err
		}
		for _1 := 0; _1 < len(v.Outputs[_0].Coins); _1++ {
			err = codonEncodeString(w, v.Outputs[_0].Coins[_1].Denom)
			if err != nil {
				return err
			}
			err = EncodeInt(w, v.Outputs[_0].Coins[_1].Amount)
			if err != nil {
				return err
			}
			// end of v.Outputs[_0].Coins[_1]
		}
		// end of v.Outputs[_0]
	}
	return nil
} //End of EncodeMsgMultiSend

func DecodeMsgMultiSend(bz []byte) (MsgMultiSend, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgMultiSend
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Inputs = make([]Input, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Inputs[_0], n, err = DecodeInput(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Outputs = make([]Output, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Outputs[_0], n, err = DecodeOutput(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgMultiSend

func RandMsgMultiSend(r RandSrc) MsgMultiSend {
	// codon version: 1
	var length int
	var v MsgMultiSend
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Inputs = make([]Input, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Inputs[_0] = RandInput(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Outputs = make([]Output, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Outputs[_0] = RandOutput(r)
	}
	return v
} //End of RandMsgMultiSend

// Non-Interface
func EncodeMsgSend(w io.Writer, v MsgSend) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.FromAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ToAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Amount)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Amount); _0++ {
		err = codonEncodeString(w, v.Amount[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Amount[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Amount[_0]
	}
	return nil
} //End of EncodeMsgSend

func DecodeMsgSend(bz []byte) (MsgSend, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgSend
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.FromAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ToAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgSend

func RandMsgSend(r RandSrc) MsgSend {
	// codon version: 1
	var length int
	var v MsgSend
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ToAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandMsgSend

// Non-Interface
func EncodeMsgVerifyInvariant(w io.Writer, v MsgVerifyInvariant) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.InvariantModuleName)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.InvariantRoute)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgVerifyInvariant

func DecodeMsgVerifyInvariant(bz []byte) (MsgVerifyInvariant, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgVerifyInvariant
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.InvariantModuleName = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.InvariantRoute = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgVerifyInvariant

func RandMsgVerifyInvariant(r RandSrc) MsgVerifyInvariant {
	// codon version: 1
	var length int
	var v MsgVerifyInvariant
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.InvariantModuleName = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.InvariantRoute = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandMsgVerifyInvariant

// Non-Interface
func EncodeSupply(w io.Writer, v Supply) error {
	// codon version: 1
	var err error
	err = codonEncodeVarint(w, int64(len(v.Total)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Total); _0++ {
		err = codonEncodeString(w, v.Total[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Total[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Total[_0]
	}
	return nil
} //End of EncodeSupply

func DecodeSupply(bz []byte) (Supply, int, error) {
	// codon version: 1
	var err error
	var length int
	var v Supply
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Total = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Total[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeSupply

func RandSupply(r RandSrc) Supply {
	// codon version: 1
	var length int
	var v Supply
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Total = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Total[_0] = RandCoin(r)
	}
	return v
} //End of RandSupply

// Non-Interface
func EncodeAccountX(w io.Writer, v AccountX) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.MemoRequired)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.LockedCoins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.LockedCoins); _0++ {
		err = codonEncodeString(w, v.LockedCoins[_0].Coin.Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.LockedCoins[_0].Coin.Amount)
		if err != nil {
			return err
		}
		// end of v.LockedCoins[_0].Coin
		err = codonEncodeVarint(w, int64(v.LockedCoins[_0].UnlockTime))
		if err != nil {
			return err
		}
		err = codonEncodeByteSlice(w, v.LockedCoins[_0].FromAddress[:])
		if err != nil {
			return err
		}
		err = codonEncodeByteSlice(w, v.LockedCoins[_0].Supervisor[:])
		if err != nil {
			return err
		}
		err = codonEncodeVarint(w, int64(v.LockedCoins[_0].Reward))
		if err != nil {
			return err
		}
		// end of v.LockedCoins[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.FrozenCoins)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.FrozenCoins); _0++ {
		err = codonEncodeString(w, v.FrozenCoins[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.FrozenCoins[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.FrozenCoins[_0]
	}
	return nil
} //End of EncodeAccountX

func DecodeAccountX(bz []byte) (AccountX, int, error) {
	// codon version: 1
	var err error
	var length int
	var v AccountX
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.MemoRequired = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.LockedCoins = make([]LockedCoin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.LockedCoins[_0], n, err = DecodeLockedCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.FrozenCoins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.FrozenCoins[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeAccountX

func RandAccountX(r RandSrc) AccountX {
	// codon version: 1
	var length int
	var v AccountX
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	v.MemoRequired = r.GetBool()
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.LockedCoins = make([]LockedCoin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.LockedCoins[_0] = RandLockedCoin(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FrozenCoins = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.FrozenCoins[_0] = RandCoin(r)
	}
	return v
} //End of RandAccountX

// Non-Interface
func EncodeMsgMultiSendX(w io.Writer, v MsgMultiSendX) error {
	// codon version: 1
	var err error
	err = codonEncodeVarint(w, int64(len(v.Inputs)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Inputs); _0++ {
		err = codonEncodeByteSlice(w, v.Inputs[_0].Address[:])
		if err != nil {
			return err
		}
		err = codonEncodeVarint(w, int64(len(v.Inputs[_0].Coins)))
		if err != nil {
			return err
		}
		for _1 := 0; _1 < len(v.Inputs[_0].Coins); _1++ {
			err = codonEncodeString(w, v.Inputs[_0].Coins[_1].Denom)
			if err != nil {
				return err
			}
			err = EncodeInt(w, v.Inputs[_0].Coins[_1].Amount)
			if err != nil {
				return err
			}
			// end of v.Inputs[_0].Coins[_1]
		}
		// end of v.Inputs[_0]
	}
	err = codonEncodeVarint(w, int64(len(v.Outputs)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Outputs); _0++ {
		err = codonEncodeByteSlice(w, v.Outputs[_0].Address[:])
		if err != nil {
			return err
		}
		err = codonEncodeVarint(w, int64(len(v.Outputs[_0].Coins)))
		if err != nil {
			return err
		}
		for _1 := 0; _1 < len(v.Outputs[_0].Coins); _1++ {
			err = codonEncodeString(w, v.Outputs[_0].Coins[_1].Denom)
			if err != nil {
				return err
			}
			err = EncodeInt(w, v.Outputs[_0].Coins[_1].Amount)
			if err != nil {
				return err
			}
			// end of v.Outputs[_0].Coins[_1]
		}
		// end of v.Outputs[_0]
	}
	return nil
} //End of EncodeMsgMultiSendX

func DecodeMsgMultiSendX(bz []byte) (MsgMultiSendX, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgMultiSendX
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Inputs = make([]Input, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Inputs[_0], n, err = DecodeInput(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Outputs = make([]Output, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Outputs[_0], n, err = DecodeOutput(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgMultiSendX

func RandMsgMultiSendX(r RandSrc) MsgMultiSendX {
	// codon version: 1
	var length int
	var v MsgMultiSendX
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Inputs = make([]Input, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Inputs[_0] = RandInput(r)
	}
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Outputs = make([]Output, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Outputs[_0] = RandOutput(r)
	}
	return v
} //End of RandMsgMultiSendX

// Non-Interface
func EncodeMsgSendX(w io.Writer, v MsgSendX) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.FromAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.ToAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Amount)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Amount); _0++ {
		err = codonEncodeString(w, v.Amount[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Amount[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Amount[_0]
	}
	err = codonEncodeVarint(w, int64(v.UnlockTime))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgSendX

func DecodeMsgSendX(bz []byte) (MsgSendX, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgSendX
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.FromAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ToAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	v.UnlockTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgSendX

func RandMsgSendX(r RandSrc) MsgSendX {
	// codon version: 1
	var length int
	var v MsgSendX
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.ToAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	v.UnlockTime = r.GetInt64()
	return v
} //End of RandMsgSendX

// Non-Interface
func EncodeMsgSetMemoRequired(w io.Writer, v MsgSetMemoRequired) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Address[:])
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.Required)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgSetMemoRequired

func DecodeMsgSetMemoRequired(bz []byte) (MsgSetMemoRequired, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgSetMemoRequired
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Address, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Required = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgSetMemoRequired

func RandMsgSetMemoRequired(r RandSrc) MsgSetMemoRequired {
	// codon version: 1
	var length int
	var v MsgSetMemoRequired
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Address = r.GetBytes(length)
	v.Required = r.GetBool()
	return v
} //End of RandMsgSetMemoRequired

// Non-Interface
func EncodeBaseToken(w io.Writer, v BaseToken) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Name)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.TotalSupply)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.SendLock)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Owner[:])
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.Mintable)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.Burnable)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.AddrForbiddable)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.TokenForbiddable)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.TotalBurn)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.TotalMint)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.IsForbidden)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.URL)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Identity)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeBaseToken

func DecodeBaseToken(bz []byte) (BaseToken, int, error) {
	// codon version: 1
	var err error
	var length int
	var v BaseToken
	var n int
	var total int
	v.Name = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TotalSupply, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.SendLock, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Owner, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Mintable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Burnable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.AddrForbiddable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TokenForbiddable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TotalBurn, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TotalMint, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.IsForbidden = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.URL = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Identity = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeBaseToken

func RandBaseToken(r RandSrc) BaseToken {
	// codon version: 1
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

// Non-Interface
func EncodeMsgAddTokenWhitelist(w io.Writer, v MsgAddTokenWhitelist) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Whitelist)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Whitelist); _0++ {
		err = codonEncodeByteSlice(w, v.Whitelist[_0][:])
		if err != nil {
			return err
		}
	}
	return nil
} //End of EncodeMsgAddTokenWhitelist

func DecodeMsgAddTokenWhitelist(bz []byte) (MsgAddTokenWhitelist, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgAddTokenWhitelist
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Whitelist = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = codonDecodeInt(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		v.Whitelist[_0], n, err = codonGetByteSlice(bz, length)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgAddTokenWhitelist

func RandMsgAddTokenWhitelist(r RandSrc) MsgAddTokenWhitelist {
	// codon version: 1
	var length int
	var v MsgAddTokenWhitelist
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Whitelist = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Whitelist[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgAddTokenWhitelist

// Non-Interface
func EncodeMsgBurnToken(w io.Writer, v MsgBurnToken) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Amount)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgBurnToken

func DecodeMsgBurnToken(bz []byte) (MsgBurnToken, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgBurnToken
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgBurnToken

func RandMsgBurnToken(r RandSrc) MsgBurnToken {
	// codon version: 1
	var length int
	var v MsgBurnToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = RandInt(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgBurnToken

// Non-Interface
func EncodeMsgForbidAddr(w io.Writer, v MsgForbidAddr) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddr[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Addresses)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Addresses); _0++ {
		err = codonEncodeByteSlice(w, v.Addresses[_0][:])
		if err != nil {
			return err
		}
	}
	return nil
} //End of EncodeMsgForbidAddr

func DecodeMsgForbidAddr(bz []byte) (MsgForbidAddr, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgForbidAddr
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddr, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Addresses = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = codonDecodeInt(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		v.Addresses[_0], n, err = codonGetByteSlice(bz, length)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgForbidAddr

func RandMsgForbidAddr(r RandSrc) MsgForbidAddr {
	// codon version: 1
	var length int
	var v MsgForbidAddr
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddr = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Addresses = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Addresses[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgForbidAddr

// Non-Interface
func EncodeMsgForbidToken(w io.Writer, v MsgForbidToken) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgForbidToken

func DecodeMsgForbidToken(bz []byte) (MsgForbidToken, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgForbidToken
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgForbidToken

func RandMsgForbidToken(r RandSrc) MsgForbidToken {
	// codon version: 1
	var length int
	var v MsgForbidToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgForbidToken

// Non-Interface
func EncodeMsgIssueToken(w io.Writer, v MsgIssueToken) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Name)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.TotalSupply)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Owner[:])
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.Mintable)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.Burnable)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.AddrForbiddable)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.TokenForbiddable)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.URL)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Identity)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgIssueToken

func DecodeMsgIssueToken(bz []byte) (MsgIssueToken, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgIssueToken
	var n int
	var total int
	v.Name = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TotalSupply, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Owner, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Mintable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Burnable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.AddrForbiddable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TokenForbiddable = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.URL = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Identity = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgIssueToken

func RandMsgIssueToken(r RandSrc) MsgIssueToken {
	// codon version: 1
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

// Non-Interface
func EncodeMsgMintToken(w io.Writer, v MsgMintToken) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.Amount)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgMintToken

func DecodeMsgMintToken(bz []byte) (MsgMintToken, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgMintToken
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgMintToken

func RandMsgMintToken(r RandSrc) MsgMintToken {
	// codon version: 1
	var length int
	var v MsgMintToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Amount = RandInt(r)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgMintToken

// Non-Interface
func EncodeMsgModifyTokenInfo(w io.Writer, v MsgModifyTokenInfo) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.URL)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Description)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Identity)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgModifyTokenInfo

func DecodeMsgModifyTokenInfo(bz []byte) (MsgModifyTokenInfo, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgModifyTokenInfo
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.URL = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Description = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Identity = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgModifyTokenInfo

func RandMsgModifyTokenInfo(r RandSrc) MsgModifyTokenInfo {
	// codon version: 1
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

// Non-Interface
func EncodeMsgRemoveTokenWhitelist(w io.Writer, v MsgRemoveTokenWhitelist) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddress[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Whitelist)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Whitelist); _0++ {
		err = codonEncodeByteSlice(w, v.Whitelist[_0][:])
		if err != nil {
			return err
		}
	}
	return nil
} //End of EncodeMsgRemoveTokenWhitelist

func DecodeMsgRemoveTokenWhitelist(bz []byte) (MsgRemoveTokenWhitelist, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgRemoveTokenWhitelist
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Whitelist = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = codonDecodeInt(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		v.Whitelist[_0], n, err = codonGetByteSlice(bz, length)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgRemoveTokenWhitelist

func RandMsgRemoveTokenWhitelist(r RandSrc) MsgRemoveTokenWhitelist {
	// codon version: 1
	var length int
	var v MsgRemoveTokenWhitelist
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Whitelist = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Whitelist[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgRemoveTokenWhitelist

// Non-Interface
func EncodeMsgTransferOwnership(w io.Writer, v MsgTransferOwnership) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OriginalOwner[:])
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.NewOwner[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgTransferOwnership

func DecodeMsgTransferOwnership(bz []byte) (MsgTransferOwnership, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgTransferOwnership
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OriginalOwner, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.NewOwner, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgTransferOwnership

func RandMsgTransferOwnership(r RandSrc) MsgTransferOwnership {
	// codon version: 1
	var length int
	var v MsgTransferOwnership
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OriginalOwner = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.NewOwner = r.GetBytes(length)
	return v
} //End of RandMsgTransferOwnership

// Non-Interface
func EncodeMsgUnForbidAddr(w io.Writer, v MsgUnForbidAddr) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddr[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Addresses)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Addresses); _0++ {
		err = codonEncodeByteSlice(w, v.Addresses[_0][:])
		if err != nil {
			return err
		}
	}
	return nil
} //End of EncodeMsgUnForbidAddr

func DecodeMsgUnForbidAddr(bz []byte) (MsgUnForbidAddr, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgUnForbidAddr
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddr, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Addresses = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = codonDecodeInt(bz, &n, &err)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
		v.Addresses[_0], n, err = codonGetByteSlice(bz, length)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgUnForbidAddr

func RandMsgUnForbidAddr(r RandSrc) MsgUnForbidAddr {
	// codon version: 1
	var length int
	var v MsgUnForbidAddr
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddr = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Addresses = make([]AccAddress, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of slice
		length = 1 + int(r.GetUint()%(MaxSliceLength-1))
		v.Addresses[_0] = r.GetBytes(length)
	}
	return v
} //End of RandMsgUnForbidAddr

// Non-Interface
func EncodeMsgUnForbidToken(w io.Writer, v MsgUnForbidToken) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Symbol)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.OwnerAddress[:])
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgUnForbidToken

func DecodeMsgUnForbidToken(bz []byte) (MsgUnForbidToken, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgUnForbidToken
	var n int
	var total int
	v.Symbol = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OwnerAddress, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgUnForbidToken

func RandMsgUnForbidToken(r RandSrc) MsgUnForbidToken {
	// codon version: 1
	var length int
	var v MsgUnForbidToken
	v.Symbol = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.OwnerAddress = r.GetBytes(length)
	return v
} //End of RandMsgUnForbidToken

// Non-Interface
func EncodeMsgBancorCancel(w io.Writer, v MsgBancorCancel) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Owner[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Stock)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Money)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgBancorCancel

func DecodeMsgBancorCancel(bz []byte) (MsgBancorCancel, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgBancorCancel
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Owner, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Stock = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Money = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgBancorCancel

func RandMsgBancorCancel(r RandSrc) MsgBancorCancel {
	// codon version: 1
	var length int
	var v MsgBancorCancel
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Owner = r.GetBytes(length)
	v.Stock = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Money = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandMsgBancorCancel

// Non-Interface
func EncodeMsgBancorInit(w io.Writer, v MsgBancorInit) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Owner[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Stock)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Money)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.InitPrice)
	if err != nil {
		return err
	}
	err = EncodeInt(w, v.MaxSupply)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.MaxPrice)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.StockPrecision)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.EarliestCancelTime))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgBancorInit

func DecodeMsgBancorInit(bz []byte) (MsgBancorInit, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgBancorInit
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Owner, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Stock = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Money = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.InitPrice = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.MaxSupply, n, err = DecodeInt(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.MaxPrice = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.StockPrecision = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.EarliestCancelTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgBancorInit

func RandMsgBancorInit(r RandSrc) MsgBancorInit {
	// codon version: 1
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

// Non-Interface
func EncodeMsgBancorTrade(w io.Writer, v MsgBancorTrade) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Stock)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Money)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Amount))
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.IsBuy)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.MoneyLimit))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgBancorTrade

func DecodeMsgBancorTrade(bz []byte) (MsgBancorTrade, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgBancorTrade
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Stock = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Money = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.IsBuy = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.MoneyLimit = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgBancorTrade

func RandMsgBancorTrade(r RandSrc) MsgBancorTrade {
	// codon version: 1
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

// Non-Interface
func EncodeMsgCancelOrder(w io.Writer, v MsgCancelOrder) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.OrderID)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgCancelOrder

func DecodeMsgCancelOrder(bz []byte) (MsgCancelOrder, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgCancelOrder
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OrderID = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgCancelOrder

func RandMsgCancelOrder(r RandSrc) MsgCancelOrder {
	// codon version: 1
	var length int
	var v MsgCancelOrder
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.OrderID = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	return v
} //End of RandMsgCancelOrder

// Non-Interface
func EncodeMsgCancelTradingPair(w io.Writer, v MsgCancelTradingPair) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.TradingPair)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.EffectiveTime))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgCancelTradingPair

func DecodeMsgCancelTradingPair(bz []byte) (MsgCancelTradingPair, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgCancelTradingPair
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TradingPair = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.EffectiveTime = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgCancelTradingPair

func RandMsgCancelTradingPair(r RandSrc) MsgCancelTradingPair {
	// codon version: 1
	var length int
	var v MsgCancelTradingPair
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.TradingPair = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.EffectiveTime = r.GetInt64()
	return v
} //End of RandMsgCancelTradingPair

// Non-Interface
func EncodeMsgCreateOrder(w io.Writer, v MsgCreateOrder) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.Identify)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.TradingPair)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.OrderType)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.PricePrecision)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Price))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Quantity))
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.Side)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.TimeInForce))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.ExistBlocks))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgCreateOrder

func DecodeMsgCreateOrder(bz []byte) (MsgCreateOrder, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgCreateOrder
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Identify = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TradingPair = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OrderType = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Price = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Quantity = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Side = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TimeInForce = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ExistBlocks = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgCreateOrder

func RandMsgCreateOrder(r RandSrc) MsgCreateOrder {
	// codon version: 1
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

// Non-Interface
func EncodeMsgCreateTradingPair(w io.Writer, v MsgCreateTradingPair) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Stock)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Money)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Creator[:])
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.PricePrecision)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.OrderPrecision)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgCreateTradingPair

func DecodeMsgCreateTradingPair(bz []byte) (MsgCreateTradingPair, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgCreateTradingPair
	var n int
	var total int
	v.Stock = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Money = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Creator, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OrderPrecision = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgCreateTradingPair

func RandMsgCreateTradingPair(r RandSrc) MsgCreateTradingPair {
	// codon version: 1
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

// Non-Interface
func EncodeMsgModifyPricePrecision(w io.Writer, v MsgModifyPricePrecision) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.TradingPair)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.PricePrecision)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgModifyPricePrecision

func DecodeMsgModifyPricePrecision(bz []byte) (MsgModifyPricePrecision, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgModifyPricePrecision
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TradingPair = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgModifyPricePrecision

func RandMsgModifyPricePrecision(r RandSrc) MsgModifyPricePrecision {
	// codon version: 1
	var length int
	var v MsgModifyPricePrecision
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Sender = r.GetBytes(length)
	v.TradingPair = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.PricePrecision = r.GetUint8()
	return v
} //End of RandMsgModifyPricePrecision

// Non-Interface
func EncodeOrder(w io.Writer, v Order) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeUvarint(w, uint64(v.Sequence))
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.Identify)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.TradingPair)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.OrderType)
	if err != nil {
		return err
	}
	err = EncodeDec(w, v.Price)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Quantity))
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.Side)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.TimeInForce))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Height))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.FrozenCommission))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.ExistBlocks))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.FrozenFeatureFee))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.LeftStock))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Freeze))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.DealStock))
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.DealMoney))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeOrder

func DecodeOrder(bz []byte) (Order, int, error) {
	// codon version: 1
	var err error
	var length int
	var v Order
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sequence = uint64(codonDecodeUint64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Identify = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TradingPair = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OrderType = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Price, n, err = DecodeDec(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Quantity = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Side = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.TimeInForce = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Height = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.FrozenCommission = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ExistBlocks = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.FrozenFeatureFee = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.LeftStock = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Freeze = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DealStock = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.DealMoney = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeOrder

func RandOrder(r RandSrc) Order {
	// codon version: 1
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
	v.FrozenCommission = r.GetInt64()
	v.ExistBlocks = r.GetInt64()
	v.FrozenFeatureFee = r.GetInt64()
	v.LeftStock = r.GetInt64()
	v.Freeze = r.GetInt64()
	v.DealStock = r.GetInt64()
	v.DealMoney = r.GetInt64()
	return v
} //End of RandOrder

// Non-Interface
func EncodeMarketInfo(w io.Writer, v MarketInfo) error {
	// codon version: 1
	var err error
	err = codonEncodeString(w, v.Stock)
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Money)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.PricePrecision)
	if err != nil {
		return err
	}
	err = EncodeDec(w, v.LastExecutedPrice)
	if err != nil {
		return err
	}
	err = codonEncodeUint8(w, v.OrderPrecision)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMarketInfo

func DecodeMarketInfo(bz []byte) (MarketInfo, int, error) {
	// codon version: 1
	var err error
	var v MarketInfo
	var n int
	var total int
	v.Stock = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Money = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.PricePrecision = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.LastExecutedPrice, n, err = DecodeDec(bz)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.OrderPrecision = uint8(codonDecodeUint8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMarketInfo

func RandMarketInfo(r RandSrc) MarketInfo {
	// codon version: 1
	var v MarketInfo
	v.Stock = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.Money = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.PricePrecision = r.GetUint8()
	v.LastExecutedPrice = RandDec(r)
	v.OrderPrecision = r.GetUint8()
	return v
} //End of RandMarketInfo

// Non-Interface
func EncodeMsgDonateToCommunityPool(w io.Writer, v MsgDonateToCommunityPool) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.FromAddr[:])
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.Amount)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.Amount); _0++ {
		err = codonEncodeString(w, v.Amount[_0].Denom)
		if err != nil {
			return err
		}
		err = EncodeInt(w, v.Amount[_0].Amount)
		if err != nil {
			return err
		}
		// end of v.Amount[_0]
	}
	return nil
} //End of EncodeMsgDonateToCommunityPool

func DecodeMsgDonateToCommunityPool(bz []byte) (MsgDonateToCommunityPool, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgDonateToCommunityPool
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.FromAddr, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0], n, err = DecodeCoin(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgDonateToCommunityPool

func RandMsgDonateToCommunityPool(r RandSrc) MsgDonateToCommunityPool {
	// codon version: 1
	var length int
	var v MsgDonateToCommunityPool
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.FromAddr = r.GetBytes(length)
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Amount = make([]Coin, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.Amount[_0] = RandCoin(r)
	}
	return v
} //End of RandMsgDonateToCommunityPool

// Non-Interface
func EncodeMsgCommentToken(w io.Writer, v MsgCommentToken) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Sender[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Token)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(v.Donation))
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Title)
	if err != nil {
		return err
	}
	err = codonEncodeByteSlice(w, v.Content[:])
	if err != nil {
		return err
	}
	err = codonEncodeInt8(w, v.ContentType)
	if err != nil {
		return err
	}
	err = codonEncodeVarint(w, int64(len(v.References)))
	if err != nil {
		return err
	}
	for _0 := 0; _0 < len(v.References); _0++ {
		err = codonEncodeUvarint(w, uint64(v.References[_0].ID))
		if err != nil {
			return err
		}
		err = codonEncodeByteSlice(w, v.References[_0].RewardTarget[:])
		if err != nil {
			return err
		}
		err = codonEncodeString(w, v.References[_0].RewardToken)
		if err != nil {
			return err
		}
		err = codonEncodeVarint(w, int64(v.References[_0].RewardAmount))
		if err != nil {
			return err
		}
		err = codonEncodeVarint(w, int64(len(v.References[_0].Attitudes)))
		if err != nil {
			return err
		}
		for _1 := 0; _1 < len(v.References[_0].Attitudes); _1++ {
			err = codonEncodeVarint(w, int64(v.References[_0].Attitudes[_1]))
			if err != nil {
				return err
			}
		}
		// end of v.References[_0]
	}
	return nil
} //End of EncodeMsgCommentToken

func DecodeMsgCommentToken(bz []byte) (MsgCommentToken, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgCommentToken
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Sender, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Token = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Donation = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Title = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Content, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.ContentType = int8(codonDecodeInt8(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.References = make([]CommentRef, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.References[_0], n, err = DecodeCommentRef(bz)
		if err != nil {
			return v, total, err
		}
		bz = bz[n:]
		total += n
	}
	return v, total, nil
} //End of DecodeMsgCommentToken

func RandMsgCommentToken(r RandSrc) MsgCommentToken {
	// codon version: 1
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
	v.References = make([]CommentRef, length)
	for _0, length_0 := 0, length; _0 < length_0; _0++ { //slice of struct
		v.References[_0] = RandCommentRef(r)
	}
	return v
} //End of RandMsgCommentToken

// Non-Interface
func EncodeState(w io.Writer, v State) error {
	// codon version: 1
	var err error
	err = codonEncodeVarint(w, int64(v.HeightAdjustment))
	if err != nil {
		return err
	}
	return nil
} //End of EncodeState

func DecodeState(bz []byte) (State, int, error) {
	// codon version: 1
	var err error
	var v State
	var n int
	var total int
	v.HeightAdjustment = int64(codonDecodeInt64(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeState

func RandState(r RandSrc) State {
	// codon version: 1
	var v State
	v.HeightAdjustment = r.GetInt64()
	return v
} //End of RandState

// Non-Interface
func EncodeMsgAliasUpdate(w io.Writer, v MsgAliasUpdate) error {
	// codon version: 1
	var err error
	err = codonEncodeByteSlice(w, v.Owner[:])
	if err != nil {
		return err
	}
	err = codonEncodeString(w, v.Alias)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.IsAdd)
	if err != nil {
		return err
	}
	err = codonEncodeBool(w, v.AsDefault)
	if err != nil {
		return err
	}
	return nil
} //End of EncodeMsgAliasUpdate

func DecodeMsgAliasUpdate(bz []byte) (MsgAliasUpdate, int, error) {
	// codon version: 1
	var err error
	var length int
	var v MsgAliasUpdate
	var n int
	var total int
	length = codonDecodeInt(bz, &n, &err)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Owner, n, err = codonGetByteSlice(bz, length)
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.Alias = string(codonDecodeString(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.IsAdd = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	v.AsDefault = bool(codonDecodeBool(bz, &n, &err))
	if err != nil {
		return v, total, err
	}
	bz = bz[n:]
	total += n
	return v, total, nil
} //End of DecodeMsgAliasUpdate

func RandMsgAliasUpdate(r RandSrc) MsgAliasUpdate {
	// codon version: 1
	var length int
	var v MsgAliasUpdate
	length = 1 + int(r.GetUint()%(MaxSliceLength-1))
	v.Owner = r.GetBytes(length)
	v.Alias = r.GetString(1 + int(r.GetUint()%(MaxStringLength-1)))
	v.IsAdd = r.GetBool()
	v.AsDefault = r.GetBool()
	return v
} //End of RandMsgAliasUpdate

// Interface
func EncodePubKey(w io.Writer, x interface{}) error {
	switch v := x.(type) {
	case PubKeyEd25519:
		w.Write(getMagicBytes("PubKeyEd25519"))
		return EncodePubKeyEd25519(w, v)
	case *PubKeyEd25519:
		w.Write(getMagicBytes("PubKeyEd25519"))
		return EncodePubKeyEd25519(w, *v)
	case PubKeyMultisigThreshold:
		w.Write(getMagicBytes("PubKeyMultisigThreshold"))
		return EncodePubKeyMultisigThreshold(w, v)
	case *PubKeyMultisigThreshold:
		w.Write(getMagicBytes("PubKeyMultisigThreshold"))
		return EncodePubKeyMultisigThreshold(w, *v)
	case PubKeySecp256k1:
		w.Write(getMagicBytes("PubKeySecp256k1"))
		return EncodePubKeySecp256k1(w, v)
	case *PubKeySecp256k1:
		w.Write(getMagicBytes("PubKeySecp256k1"))
		return EncodePubKeySecp256k1(w, *v)
	case StdSignature:
		w.Write(getMagicBytes("StdSignature"))
		return EncodeStdSignature(w, v)
	case *StdSignature:
		w.Write(getMagicBytes("StdSignature"))
		return EncodeStdSignature(w, *v)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DecodePubKey(bz []byte) (PubKey, int, error) {
	var v PubKey
	var magicBytes [4]byte
	var n int
	for i := 0; i < 4; i++ {
		magicBytes[i] = bz[i]
	}
	switch magicBytes {
	case [4]byte{108, 143, 2, 48}:
		v, n, err := DecodePubKeyEd25519(bz[4:])
		return v, n + 4, err
	case [4]byte{131, 227, 102, 173}:
		v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
		return v, n + 4, err
	case [4]byte{10, 126, 85, 105}:
		v, n, err := DecodePubKeySecp256k1(bz[4:])
		return v, n + 4, err
	case [4]byte{88, 244, 106, 18}:
		v, n, err := DecodeStdSignature(bz[4:])
		return v, n + 4, err
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodePubKey
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
// Interface
func EncodeMsg(w io.Writer, x interface{}) error {
	switch v := x.(type) {
	case MsgAddTokenWhitelist:
		w.Write(getMagicBytes("MsgAddTokenWhitelist"))
		return EncodeMsgAddTokenWhitelist(w, v)
	case *MsgAddTokenWhitelist:
		w.Write(getMagicBytes("MsgAddTokenWhitelist"))
		return EncodeMsgAddTokenWhitelist(w, *v)
	case MsgAliasUpdate:
		w.Write(getMagicBytes("MsgAliasUpdate"))
		return EncodeMsgAliasUpdate(w, v)
	case *MsgAliasUpdate:
		w.Write(getMagicBytes("MsgAliasUpdate"))
		return EncodeMsgAliasUpdate(w, *v)
	case MsgBancorCancel:
		w.Write(getMagicBytes("MsgBancorCancel"))
		return EncodeMsgBancorCancel(w, v)
	case *MsgBancorCancel:
		w.Write(getMagicBytes("MsgBancorCancel"))
		return EncodeMsgBancorCancel(w, *v)
	case MsgBancorInit:
		w.Write(getMagicBytes("MsgBancorInit"))
		return EncodeMsgBancorInit(w, v)
	case *MsgBancorInit:
		w.Write(getMagicBytes("MsgBancorInit"))
		return EncodeMsgBancorInit(w, *v)
	case MsgBancorTrade:
		w.Write(getMagicBytes("MsgBancorTrade"))
		return EncodeMsgBancorTrade(w, v)
	case *MsgBancorTrade:
		w.Write(getMagicBytes("MsgBancorTrade"))
		return EncodeMsgBancorTrade(w, *v)
	case MsgBeginRedelegate:
		w.Write(getMagicBytes("MsgBeginRedelegate"))
		return EncodeMsgBeginRedelegate(w, v)
	case *MsgBeginRedelegate:
		w.Write(getMagicBytes("MsgBeginRedelegate"))
		return EncodeMsgBeginRedelegate(w, *v)
	case MsgBurnToken:
		w.Write(getMagicBytes("MsgBurnToken"))
		return EncodeMsgBurnToken(w, v)
	case *MsgBurnToken:
		w.Write(getMagicBytes("MsgBurnToken"))
		return EncodeMsgBurnToken(w, *v)
	case MsgCancelOrder:
		w.Write(getMagicBytes("MsgCancelOrder"))
		return EncodeMsgCancelOrder(w, v)
	case *MsgCancelOrder:
		w.Write(getMagicBytes("MsgCancelOrder"))
		return EncodeMsgCancelOrder(w, *v)
	case MsgCancelTradingPair:
		w.Write(getMagicBytes("MsgCancelTradingPair"))
		return EncodeMsgCancelTradingPair(w, v)
	case *MsgCancelTradingPair:
		w.Write(getMagicBytes("MsgCancelTradingPair"))
		return EncodeMsgCancelTradingPair(w, *v)
	case MsgCommentToken:
		w.Write(getMagicBytes("MsgCommentToken"))
		return EncodeMsgCommentToken(w, v)
	case *MsgCommentToken:
		w.Write(getMagicBytes("MsgCommentToken"))
		return EncodeMsgCommentToken(w, *v)
	case MsgCreateOrder:
		w.Write(getMagicBytes("MsgCreateOrder"))
		return EncodeMsgCreateOrder(w, v)
	case *MsgCreateOrder:
		w.Write(getMagicBytes("MsgCreateOrder"))
		return EncodeMsgCreateOrder(w, *v)
	case MsgCreateTradingPair:
		w.Write(getMagicBytes("MsgCreateTradingPair"))
		return EncodeMsgCreateTradingPair(w, v)
	case *MsgCreateTradingPair:
		w.Write(getMagicBytes("MsgCreateTradingPair"))
		return EncodeMsgCreateTradingPair(w, *v)
	case MsgCreateValidator:
		w.Write(getMagicBytes("MsgCreateValidator"))
		return EncodeMsgCreateValidator(w, v)
	case *MsgCreateValidator:
		w.Write(getMagicBytes("MsgCreateValidator"))
		return EncodeMsgCreateValidator(w, *v)
	case MsgDelegate:
		w.Write(getMagicBytes("MsgDelegate"))
		return EncodeMsgDelegate(w, v)
	case *MsgDelegate:
		w.Write(getMagicBytes("MsgDelegate"))
		return EncodeMsgDelegate(w, *v)
	case MsgDeposit:
		w.Write(getMagicBytes("MsgDeposit"))
		return EncodeMsgDeposit(w, v)
	case *MsgDeposit:
		w.Write(getMagicBytes("MsgDeposit"))
		return EncodeMsgDeposit(w, *v)
	case MsgDonateToCommunityPool:
		w.Write(getMagicBytes("MsgDonateToCommunityPool"))
		return EncodeMsgDonateToCommunityPool(w, v)
	case *MsgDonateToCommunityPool:
		w.Write(getMagicBytes("MsgDonateToCommunityPool"))
		return EncodeMsgDonateToCommunityPool(w, *v)
	case MsgEditValidator:
		w.Write(getMagicBytes("MsgEditValidator"))
		return EncodeMsgEditValidator(w, v)
	case *MsgEditValidator:
		w.Write(getMagicBytes("MsgEditValidator"))
		return EncodeMsgEditValidator(w, *v)
	case MsgForbidAddr:
		w.Write(getMagicBytes("MsgForbidAddr"))
		return EncodeMsgForbidAddr(w, v)
	case *MsgForbidAddr:
		w.Write(getMagicBytes("MsgForbidAddr"))
		return EncodeMsgForbidAddr(w, *v)
	case MsgForbidToken:
		w.Write(getMagicBytes("MsgForbidToken"))
		return EncodeMsgForbidToken(w, v)
	case *MsgForbidToken:
		w.Write(getMagicBytes("MsgForbidToken"))
		return EncodeMsgForbidToken(w, *v)
	case MsgIssueToken:
		w.Write(getMagicBytes("MsgIssueToken"))
		return EncodeMsgIssueToken(w, v)
	case *MsgIssueToken:
		w.Write(getMagicBytes("MsgIssueToken"))
		return EncodeMsgIssueToken(w, *v)
	case MsgMintToken:
		w.Write(getMagicBytes("MsgMintToken"))
		return EncodeMsgMintToken(w, v)
	case *MsgMintToken:
		w.Write(getMagicBytes("MsgMintToken"))
		return EncodeMsgMintToken(w, *v)
	case MsgModifyPricePrecision:
		w.Write(getMagicBytes("MsgModifyPricePrecision"))
		return EncodeMsgModifyPricePrecision(w, v)
	case *MsgModifyPricePrecision:
		w.Write(getMagicBytes("MsgModifyPricePrecision"))
		return EncodeMsgModifyPricePrecision(w, *v)
	case MsgModifyTokenInfo:
		w.Write(getMagicBytes("MsgModifyTokenInfo"))
		return EncodeMsgModifyTokenInfo(w, v)
	case *MsgModifyTokenInfo:
		w.Write(getMagicBytes("MsgModifyTokenInfo"))
		return EncodeMsgModifyTokenInfo(w, *v)
	case MsgMultiSend:
		w.Write(getMagicBytes("MsgMultiSend"))
		return EncodeMsgMultiSend(w, v)
	case *MsgMultiSend:
		w.Write(getMagicBytes("MsgMultiSend"))
		return EncodeMsgMultiSend(w, *v)
	case MsgMultiSendX:
		w.Write(getMagicBytes("MsgMultiSendX"))
		return EncodeMsgMultiSendX(w, v)
	case *MsgMultiSendX:
		w.Write(getMagicBytes("MsgMultiSendX"))
		return EncodeMsgMultiSendX(w, *v)
	case MsgRemoveTokenWhitelist:
		w.Write(getMagicBytes("MsgRemoveTokenWhitelist"))
		return EncodeMsgRemoveTokenWhitelist(w, v)
	case *MsgRemoveTokenWhitelist:
		w.Write(getMagicBytes("MsgRemoveTokenWhitelist"))
		return EncodeMsgRemoveTokenWhitelist(w, *v)
	case MsgSend:
		w.Write(getMagicBytes("MsgSend"))
		return EncodeMsgSend(w, v)
	case *MsgSend:
		w.Write(getMagicBytes("MsgSend"))
		return EncodeMsgSend(w, *v)
	case MsgSendX:
		w.Write(getMagicBytes("MsgSendX"))
		return EncodeMsgSendX(w, v)
	case *MsgSendX:
		w.Write(getMagicBytes("MsgSendX"))
		return EncodeMsgSendX(w, *v)
	case MsgSetMemoRequired:
		w.Write(getMagicBytes("MsgSetMemoRequired"))
		return EncodeMsgSetMemoRequired(w, v)
	case *MsgSetMemoRequired:
		w.Write(getMagicBytes("MsgSetMemoRequired"))
		return EncodeMsgSetMemoRequired(w, *v)
	case MsgSetWithdrawAddress:
		w.Write(getMagicBytes("MsgSetWithdrawAddress"))
		return EncodeMsgSetWithdrawAddress(w, v)
	case *MsgSetWithdrawAddress:
		w.Write(getMagicBytes("MsgSetWithdrawAddress"))
		return EncodeMsgSetWithdrawAddress(w, *v)
	case MsgSubmitProposal:
		w.Write(getMagicBytes("MsgSubmitProposal"))
		return EncodeMsgSubmitProposal(w, v)
	case *MsgSubmitProposal:
		w.Write(getMagicBytes("MsgSubmitProposal"))
		return EncodeMsgSubmitProposal(w, *v)
	case MsgTransferOwnership:
		w.Write(getMagicBytes("MsgTransferOwnership"))
		return EncodeMsgTransferOwnership(w, v)
	case *MsgTransferOwnership:
		w.Write(getMagicBytes("MsgTransferOwnership"))
		return EncodeMsgTransferOwnership(w, *v)
	case MsgUnForbidAddr:
		w.Write(getMagicBytes("MsgUnForbidAddr"))
		return EncodeMsgUnForbidAddr(w, v)
	case *MsgUnForbidAddr:
		w.Write(getMagicBytes("MsgUnForbidAddr"))
		return EncodeMsgUnForbidAddr(w, *v)
	case MsgUnForbidToken:
		w.Write(getMagicBytes("MsgUnForbidToken"))
		return EncodeMsgUnForbidToken(w, v)
	case *MsgUnForbidToken:
		w.Write(getMagicBytes("MsgUnForbidToken"))
		return EncodeMsgUnForbidToken(w, *v)
	case MsgUndelegate:
		w.Write(getMagicBytes("MsgUndelegate"))
		return EncodeMsgUndelegate(w, v)
	case *MsgUndelegate:
		w.Write(getMagicBytes("MsgUndelegate"))
		return EncodeMsgUndelegate(w, *v)
	case MsgUnjail:
		w.Write(getMagicBytes("MsgUnjail"))
		return EncodeMsgUnjail(w, v)
	case *MsgUnjail:
		w.Write(getMagicBytes("MsgUnjail"))
		return EncodeMsgUnjail(w, *v)
	case MsgVerifyInvariant:
		w.Write(getMagicBytes("MsgVerifyInvariant"))
		return EncodeMsgVerifyInvariant(w, v)
	case *MsgVerifyInvariant:
		w.Write(getMagicBytes("MsgVerifyInvariant"))
		return EncodeMsgVerifyInvariant(w, *v)
	case MsgVote:
		w.Write(getMagicBytes("MsgVote"))
		return EncodeMsgVote(w, v)
	case *MsgVote:
		w.Write(getMagicBytes("MsgVote"))
		return EncodeMsgVote(w, *v)
	case MsgWithdrawDelegatorReward:
		w.Write(getMagicBytes("MsgWithdrawDelegatorReward"))
		return EncodeMsgWithdrawDelegatorReward(w, v)
	case *MsgWithdrawDelegatorReward:
		w.Write(getMagicBytes("MsgWithdrawDelegatorReward"))
		return EncodeMsgWithdrawDelegatorReward(w, *v)
	case MsgWithdrawValidatorCommission:
		w.Write(getMagicBytes("MsgWithdrawValidatorCommission"))
		return EncodeMsgWithdrawValidatorCommission(w, v)
	case *MsgWithdrawValidatorCommission:
		w.Write(getMagicBytes("MsgWithdrawValidatorCommission"))
		return EncodeMsgWithdrawValidatorCommission(w, *v)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DecodeMsg(bz []byte) (Msg, int, error) {
	var v Msg
	var magicBytes [4]byte
	var n int
	for i := 0; i < 4; i++ {
		magicBytes[i] = bz[i]
	}
	switch magicBytes {
	case [4]byte{147, 136, 220, 215}:
		v, n, err := DecodeMsgAddTokenWhitelist(bz[4:])
		return v, n + 4, err
	case [4]byte{173, 181, 17, 162}:
		v, n, err := DecodeMsgAliasUpdate(bz[4:])
		return v, n + 4, err
	case [4]byte{187, 190, 104, 91}:
		v, n, err := DecodeMsgBancorCancel(bz[4:])
		return v, n + 4, err
	case [4]byte{177, 152, 232, 60}:
		v, n, err := DecodeMsgBancorInit(bz[4:])
		return v, n + 4, err
	case [4]byte{225, 122, 18, 80}:
		v, n, err := DecodeMsgBancorTrade(bz[4:])
		return v, n + 4, err
	case [4]byte{247, 3, 0, 105}:
		v, n, err := DecodeMsgBeginRedelegate(bz[4:])
		return v, n + 4, err
	case [4]byte{228, 0, 236, 212}:
		v, n, err := DecodeMsgBurnToken(bz[4:])
		return v, n + 4, err
	case [4]byte{106, 229, 80, 141}:
		v, n, err := DecodeMsgCancelOrder(bz[4:])
		return v, n + 4, err
	case [4]byte{13, 177, 95, 127}:
		v, n, err := DecodeMsgCancelTradingPair(bz[4:])
		return v, n + 4, err
	case [4]byte{79, 125, 235, 121}:
		v, n, err := DecodeMsgCommentToken(bz[4:])
		return v, n + 4, err
	case [4]byte{246, 238, 7, 164}:
		v, n, err := DecodeMsgCreateOrder(bz[4:])
		return v, n + 4, err
	case [4]byte{130, 221, 55, 57}:
		v, n, err := DecodeMsgCreateTradingPair(bz[4:])
		return v, n + 4, err
	case [4]byte{58, 78, 252, 114}:
		v, n, err := DecodeMsgCreateValidator(bz[4:])
		return v, n + 4, err
	case [4]byte{1, 82, 140, 71}:
		v, n, err := DecodeMsgDelegate(bz[4:])
		return v, n + 4, err
	case [4]byte{205, 134, 140, 190}:
		v, n, err := DecodeMsgDeposit(bz[4:])
		return v, n + 4, err
	case [4]byte{20, 250, 115, 197}:
		v, n, err := DecodeMsgDonateToCommunityPool(bz[4:])
		return v, n + 4, err
	case [4]byte{202, 62, 140, 8}:
		v, n, err := DecodeMsgEditValidator(bz[4:])
		return v, n + 4, err
	case [4]byte{105, 235, 112, 10}:
		v, n, err := DecodeMsgForbidAddr(bz[4:])
		return v, n + 4, err
	case [4]byte{36, 174, 203, 238}:
		v, n, err := DecodeMsgForbidToken(bz[4:])
		return v, n + 4, err
	case [4]byte{233, 180, 92, 129}:
		v, n, err := DecodeMsgIssueToken(bz[4:])
		return v, n + 4, err
	case [4]byte{66, 148, 56, 203}:
		v, n, err := DecodeMsgMintToken(bz[4:])
		return v, n + 4, err
	case [4]byte{76, 91, 156, 199}:
		v, n, err := DecodeMsgModifyPricePrecision(bz[4:])
		return v, n + 4, err
	case [4]byte{40, 158, 40, 18}:
		v, n, err := DecodeMsgModifyTokenInfo(bz[4:])
		return v, n + 4, err
	case [4]byte{207, 152, 156, 90}:
		v, n, err := DecodeMsgMultiSend(bz[4:])
		return v, n + 4, err
	case [4]byte{61, 117, 88, 200}:
		v, n, err := DecodeMsgMultiSendX(bz[4:])
		return v, n + 4, err
	case [4]byte{44, 154, 68, 83}:
		v, n, err := DecodeMsgRemoveTokenWhitelist(bz[4:])
		return v, n + 4, err
	case [4]byte{100, 168, 39, 140}:
		v, n, err := DecodeMsgSend(bz[4:])
		return v, n + 4, err
	case [4]byte{198, 76, 8, 81}:
		v, n, err := DecodeMsgSendX(bz[4:])
		return v, n + 4, err
	case [4]byte{184, 238, 253, 154}:
		v, n, err := DecodeMsgSetMemoRequired(bz[4:])
		return v, n + 4, err
	case [4]byte{190, 178, 173, 144}:
		v, n, err := DecodeMsgSetWithdrawAddress(bz[4:])
		return v, n + 4, err
	case [4]byte{115, 119, 137, 48}:
		v, n, err := DecodeMsgSubmitProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{200, 224, 118, 175}:
		v, n, err := DecodeMsgTransferOwnership(bz[4:])
		return v, n + 4, err
	case [4]byte{167, 165, 166, 227}:
		v, n, err := DecodeMsgUnForbidAddr(bz[4:])
		return v, n + 4, err
	case [4]byte{78, 83, 156, 139}:
		v, n, err := DecodeMsgUnForbidToken(bz[4:])
		return v, n + 4, err
	case [4]byte{122, 66, 160, 76}:
		v, n, err := DecodeMsgUndelegate(bz[4:])
		return v, n + 4, err
	case [4]byte{216, 247, 180, 46}:
		v, n, err := DecodeMsgUnjail(bz[4:])
		return v, n + 4, err
	case [4]byte{84, 44, 219, 65}:
		v, n, err := DecodeMsgVerifyInvariant(bz[4:])
		return v, n + 4, err
	case [4]byte{238, 246, 67, 141}:
		v, n, err := DecodeMsgVote(bz[4:])
		return v, n + 4, err
	case [4]byte{94, 251, 176, 152}:
		v, n, err := DecodeMsgWithdrawDelegatorReward(bz[4:])
		return v, n + 4, err
	case [4]byte{18, 172, 190, 152}:
		v, n, err := DecodeMsgWithdrawValidatorCommission(bz[4:])
		return v, n + 4, err
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeMsg
func RandMsg(r RandSrc) Msg {
	switch r.GetUint() % 40 {
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
		return RandMsgTransferOwnership(r)
	case 32:
		return RandMsgUnForbidAddr(r)
	case 33:
		return RandMsgUnForbidToken(r)
	case 34:
		return RandMsgUndelegate(r)
	case 35:
		return RandMsgUnjail(r)
	case 36:
		return RandMsgVerifyInvariant(r)
	case 37:
		return RandMsgVote(r)
	case 38:
		return RandMsgWithdrawDelegatorReward(r)
	case 39:
		return RandMsgWithdrawValidatorCommission(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
// Interface
func EncodeAccount(w io.Writer, x interface{}) error {
	switch v := x.(type) {
	case BaseVestingAccount:
		w.Write(getMagicBytes("BaseVestingAccount"))
		return EncodeBaseVestingAccount(w, v)
	case *BaseVestingAccount:
		w.Write(getMagicBytes("BaseVestingAccount"))
		return EncodeBaseVestingAccount(w, *v)
	case ContinuousVestingAccount:
		w.Write(getMagicBytes("ContinuousVestingAccount"))
		return EncodeContinuousVestingAccount(w, v)
	case *ContinuousVestingAccount:
		w.Write(getMagicBytes("ContinuousVestingAccount"))
		return EncodeContinuousVestingAccount(w, *v)
	case DelayedVestingAccount:
		w.Write(getMagicBytes("DelayedVestingAccount"))
		return EncodeDelayedVestingAccount(w, v)
	case *DelayedVestingAccount:
		w.Write(getMagicBytes("DelayedVestingAccount"))
		return EncodeDelayedVestingAccount(w, *v)
	case ModuleAccount:
		w.Write(getMagicBytes("ModuleAccount"))
		return EncodeModuleAccount(w, v)
	case *ModuleAccount:
		w.Write(getMagicBytes("ModuleAccount"))
		return EncodeModuleAccount(w, *v)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DecodeAccount(bz []byte) (Account, int, error) {
	var v Account
	var magicBytes [4]byte
	var n int
	for i := 0; i < 4; i++ {
		magicBytes[i] = bz[i]
	}
	switch magicBytes {
	case [4]byte{178, 47, 121, 129}:
		v, n, err := DecodeBaseVestingAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{95, 96, 75, 0}:
		v, n, err := DecodeContinuousVestingAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{187, 71, 224, 1}:
		v, n, err := DecodeDelayedVestingAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{190, 107, 1, 124}:
		v, n, err := DecodeModuleAccount(bz[4:])
		return v, n + 4, err
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeAccount
func RandAccount(r RandSrc) Account {
	switch r.GetUint() % 4 {
	case 0:
		return RandBaseVestingAccount(r)
	case 1:
		return RandContinuousVestingAccount(r)
	case 2:
		return RandDelayedVestingAccount(r)
	case 3:
		return RandModuleAccount(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
// Interface
func EncodeContent(w io.Writer, x interface{}) error {
	switch v := x.(type) {
	case CommunityPoolSpendProposal:
		w.Write(getMagicBytes("CommunityPoolSpendProposal"))
		return EncodeCommunityPoolSpendProposal(w, v)
	case *CommunityPoolSpendProposal:
		w.Write(getMagicBytes("CommunityPoolSpendProposal"))
		return EncodeCommunityPoolSpendProposal(w, *v)
	case ParameterChangeProposal:
		w.Write(getMagicBytes("ParameterChangeProposal"))
		return EncodeParameterChangeProposal(w, v)
	case *ParameterChangeProposal:
		w.Write(getMagicBytes("ParameterChangeProposal"))
		return EncodeParameterChangeProposal(w, *v)
	case SoftwareUpgradeProposal:
		w.Write(getMagicBytes("SoftwareUpgradeProposal"))
		return EncodeSoftwareUpgradeProposal(w, v)
	case *SoftwareUpgradeProposal:
		w.Write(getMagicBytes("SoftwareUpgradeProposal"))
		return EncodeSoftwareUpgradeProposal(w, *v)
	case TextProposal:
		w.Write(getMagicBytes("TextProposal"))
		return EncodeTextProposal(w, v)
	case *TextProposal:
		w.Write(getMagicBytes("TextProposal"))
		return EncodeTextProposal(w, *v)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DecodeContent(bz []byte) (Content, int, error) {
	var v Content
	var magicBytes [4]byte
	var n int
	for i := 0; i < 4; i++ {
		magicBytes[i] = bz[i]
	}
	switch magicBytes {
	case [4]byte{37, 214, 119, 170}:
		v, n, err := DecodeCommunityPoolSpendProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{166, 63, 172, 210}:
		v, n, err := DecodeParameterChangeProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{37, 100, 208, 251}:
		v, n, err := DecodeSoftwareUpgradeProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{169, 32, 176, 245}:
		v, n, err := DecodeTextProposal(bz[4:])
		return v, n + 4, err
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeContent
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
func getMagicBytes(name string) []byte {
	switch name {
	case "AccAddress":
		return []byte{37, 50, 37, 208}
	case "AccountX":
		return []byte{52, 40, 250, 245}
	case "BaseAccount":
		return []byte{100, 94, 81, 72}
	case "BaseToken":
		return []byte{34, 178, 244, 51}
	case "BaseVestingAccount":
		return []byte{178, 47, 121, 129}
	case "Coin":
		return []byte{141, 140, 97, 80}
	case "CommentRef":
		return []byte{17, 162, 164, 235}
	case "CommunityPoolSpendProposal":
		return []byte{37, 214, 119, 170}
	case "ContinuousVestingAccount":
		return []byte{95, 96, 75, 0}
	case "DelayedVestingAccount":
		return []byte{187, 71, 224, 1}
	case "DuplicateVoteEvidence":
		return []byte{130, 76, 198, 17}
	case "Input":
		return []byte{165, 152, 189, 47}
	case "LockedCoin":
		return []byte{227, 236, 168, 93}
	case "MarketInfo":
		return []byte{174, 117, 167, 230}
	case "ModuleAccount":
		return []byte{190, 107, 1, 124}
	case "MsgAddTokenWhitelist":
		return []byte{147, 136, 220, 215}
	case "MsgAliasUpdate":
		return []byte{173, 181, 17, 162}
	case "MsgBancorCancel":
		return []byte{187, 190, 104, 91}
	case "MsgBancorInit":
		return []byte{177, 152, 232, 60}
	case "MsgBancorTrade":
		return []byte{225, 122, 18, 80}
	case "MsgBeginRedelegate":
		return []byte{247, 3, 0, 105}
	case "MsgBurnToken":
		return []byte{228, 0, 236, 212}
	case "MsgCancelOrder":
		return []byte{106, 229, 80, 141}
	case "MsgCancelTradingPair":
		return []byte{13, 177, 95, 127}
	case "MsgCommentToken":
		return []byte{79, 125, 235, 121}
	case "MsgCreateOrder":
		return []byte{246, 238, 7, 164}
	case "MsgCreateTradingPair":
		return []byte{130, 221, 55, 57}
	case "MsgCreateValidator":
		return []byte{58, 78, 252, 114}
	case "MsgDelegate":
		return []byte{1, 82, 140, 71}
	case "MsgDeposit":
		return []byte{205, 134, 140, 190}
	case "MsgDonateToCommunityPool":
		return []byte{20, 250, 115, 197}
	case "MsgEditValidator":
		return []byte{202, 62, 140, 8}
	case "MsgForbidAddr":
		return []byte{105, 235, 112, 10}
	case "MsgForbidToken":
		return []byte{36, 174, 203, 238}
	case "MsgIssueToken":
		return []byte{233, 180, 92, 129}
	case "MsgMintToken":
		return []byte{66, 148, 56, 203}
	case "MsgModifyPricePrecision":
		return []byte{76, 91, 156, 199}
	case "MsgModifyTokenInfo":
		return []byte{40, 158, 40, 18}
	case "MsgMultiSend":
		return []byte{207, 152, 156, 90}
	case "MsgMultiSendX":
		return []byte{61, 117, 88, 200}
	case "MsgRemoveTokenWhitelist":
		return []byte{44, 154, 68, 83}
	case "MsgSend":
		return []byte{100, 168, 39, 140}
	case "MsgSendX":
		return []byte{198, 76, 8, 81}
	case "MsgSetMemoRequired":
		return []byte{184, 238, 253, 154}
	case "MsgSetWithdrawAddress":
		return []byte{190, 178, 173, 144}
	case "MsgSubmitProposal":
		return []byte{115, 119, 137, 48}
	case "MsgTransferOwnership":
		return []byte{200, 224, 118, 175}
	case "MsgUnForbidAddr":
		return []byte{167, 165, 166, 227}
	case "MsgUnForbidToken":
		return []byte{78, 83, 156, 139}
	case "MsgUndelegate":
		return []byte{122, 66, 160, 76}
	case "MsgUnjail":
		return []byte{216, 247, 180, 46}
	case "MsgVerifyInvariant":
		return []byte{84, 44, 219, 65}
	case "MsgVote":
		return []byte{238, 246, 67, 141}
	case "MsgWithdrawDelegatorReward":
		return []byte{94, 251, 176, 152}
	case "MsgWithdrawValidatorCommission":
		return []byte{18, 172, 190, 152}
	case "Order":
		return []byte{57, 115, 63, 27}

	case "Output":
		return []byte{251, 0, 54, 127}
	case "ParamChange":
		return []byte{234, 101, 49, 27}
	case "ParameterChangeProposal":
		return []byte{166, 63, 172, 210}
	case "PrivKeyEd25519":
		return []byte{93, 160, 108, 51}
	case "PrivKeySecp256k1":
		return []byte{209, 107, 141, 98}
	case "PubKeyEd25519":
		return []byte{108, 143, 2, 48}
	case "PubKeyMultisigThreshold":
		return []byte{131, 227, 102, 173}
	case "PubKeySecp256k1":
		return []byte{10, 126, 85, 105}
	case "SignedMsgType":
		return []byte{169, 174, 252, 87}
	case "SoftwareUpgradeProposal":
		return []byte{37, 100, 208, 251}
	case "State":
		return []byte{186, 223, 120, 4}
	case "StdSignature":
		return []byte{88, 244, 106, 18}
	case "StdTx":
		return []byte{71, 175, 179, 184}
	case "Supply":
		return []byte{233, 209, 209, 86}
	case "TextProposal":
		return []byte{169, 32, 176, 245}
	case "Vote":
		return []byte{113, 227, 24, 224}
	case "VoteOption":
		return []byte{56, 159, 20, 227}
	} // end of switch
	panic("Should not reach here")
	return []byte{}
} // end of getMagicBytes
func EncodeAny(w io.Writer, x interface{}) error {
	switch v := x.(type) {
	case AccAddress:
		w.Write(getMagicBytes("AccAddress"))
		return EncodeAccAddress(w, v)
	case *AccAddress:
		w.Write(getMagicBytes("AccAddress"))
		return EncodeAccAddress(w, *v)
	case AccountX:
		w.Write(getMagicBytes("AccountX"))
		return EncodeAccountX(w, v)
	case *AccountX:
		w.Write(getMagicBytes("AccountX"))
		return EncodeAccountX(w, *v)
	case BaseAccount:
		w.Write(getMagicBytes("BaseAccount"))
		return EncodeBaseAccount(w, v)
	case *BaseAccount:
		w.Write(getMagicBytes("BaseAccount"))
		return EncodeBaseAccount(w, *v)
	case BaseToken:
		w.Write(getMagicBytes("BaseToken"))
		return EncodeBaseToken(w, v)
	case *BaseToken:
		w.Write(getMagicBytes("BaseToken"))
		return EncodeBaseToken(w, *v)
	case BaseVestingAccount:
		w.Write(getMagicBytes("BaseVestingAccount"))
		return EncodeBaseVestingAccount(w, v)
	case *BaseVestingAccount:
		w.Write(getMagicBytes("BaseVestingAccount"))
		return EncodeBaseVestingAccount(w, *v)
	case Coin:
		w.Write(getMagicBytes("Coin"))
		return EncodeCoin(w, v)
	case *Coin:
		w.Write(getMagicBytes("Coin"))
		return EncodeCoin(w, *v)
	case CommentRef:
		w.Write(getMagicBytes("CommentRef"))
		return EncodeCommentRef(w, v)
	case *CommentRef:
		w.Write(getMagicBytes("CommentRef"))
		return EncodeCommentRef(w, *v)
	case CommunityPoolSpendProposal:
		w.Write(getMagicBytes("CommunityPoolSpendProposal"))
		return EncodeCommunityPoolSpendProposal(w, v)
	case *CommunityPoolSpendProposal:
		w.Write(getMagicBytes("CommunityPoolSpendProposal"))
		return EncodeCommunityPoolSpendProposal(w, *v)
	case ContinuousVestingAccount:
		w.Write(getMagicBytes("ContinuousVestingAccount"))
		return EncodeContinuousVestingAccount(w, v)
	case *ContinuousVestingAccount:
		w.Write(getMagicBytes("ContinuousVestingAccount"))
		return EncodeContinuousVestingAccount(w, *v)
	case DelayedVestingAccount:
		w.Write(getMagicBytes("DelayedVestingAccount"))
		return EncodeDelayedVestingAccount(w, v)
	case *DelayedVestingAccount:
		w.Write(getMagicBytes("DelayedVestingAccount"))
		return EncodeDelayedVestingAccount(w, *v)
	case DuplicateVoteEvidence:
		w.Write(getMagicBytes("DuplicateVoteEvidence"))
		return EncodeDuplicateVoteEvidence(w, v)
	case *DuplicateVoteEvidence:
		w.Write(getMagicBytes("DuplicateVoteEvidence"))
		return EncodeDuplicateVoteEvidence(w, *v)
	case Input:
		w.Write(getMagicBytes("Input"))
		return EncodeInput(w, v)
	case *Input:
		w.Write(getMagicBytes("Input"))
		return EncodeInput(w, *v)
	case LockedCoin:
		w.Write(getMagicBytes("LockedCoin"))
		return EncodeLockedCoin(w, v)
	case *LockedCoin:
		w.Write(getMagicBytes("LockedCoin"))
		return EncodeLockedCoin(w, *v)
	case MarketInfo:
		w.Write(getMagicBytes("MarketInfo"))
		return EncodeMarketInfo(w, v)
	case *MarketInfo:
		w.Write(getMagicBytes("MarketInfo"))
		return EncodeMarketInfo(w, *v)
	case ModuleAccount:
		w.Write(getMagicBytes("ModuleAccount"))
		return EncodeModuleAccount(w, v)
	case *ModuleAccount:
		w.Write(getMagicBytes("ModuleAccount"))
		return EncodeModuleAccount(w, *v)
	case MsgAddTokenWhitelist:
		w.Write(getMagicBytes("MsgAddTokenWhitelist"))
		return EncodeMsgAddTokenWhitelist(w, v)
	case *MsgAddTokenWhitelist:
		w.Write(getMagicBytes("MsgAddTokenWhitelist"))
		return EncodeMsgAddTokenWhitelist(w, *v)
	case MsgAliasUpdate:
		w.Write(getMagicBytes("MsgAliasUpdate"))
		return EncodeMsgAliasUpdate(w, v)
	case *MsgAliasUpdate:
		w.Write(getMagicBytes("MsgAliasUpdate"))
		return EncodeMsgAliasUpdate(w, *v)
	case MsgBancorCancel:
		w.Write(getMagicBytes("MsgBancorCancel"))
		return EncodeMsgBancorCancel(w, v)
	case *MsgBancorCancel:
		w.Write(getMagicBytes("MsgBancorCancel"))
		return EncodeMsgBancorCancel(w, *v)
	case MsgBancorInit:
		w.Write(getMagicBytes("MsgBancorInit"))
		return EncodeMsgBancorInit(w, v)
	case *MsgBancorInit:
		w.Write(getMagicBytes("MsgBancorInit"))
		return EncodeMsgBancorInit(w, *v)
	case MsgBancorTrade:
		w.Write(getMagicBytes("MsgBancorTrade"))
		return EncodeMsgBancorTrade(w, v)
	case *MsgBancorTrade:
		w.Write(getMagicBytes("MsgBancorTrade"))
		return EncodeMsgBancorTrade(w, *v)
	case MsgBeginRedelegate:
		w.Write(getMagicBytes("MsgBeginRedelegate"))
		return EncodeMsgBeginRedelegate(w, v)
	case *MsgBeginRedelegate:
		w.Write(getMagicBytes("MsgBeginRedelegate"))
		return EncodeMsgBeginRedelegate(w, *v)
	case MsgBurnToken:
		w.Write(getMagicBytes("MsgBurnToken"))
		return EncodeMsgBurnToken(w, v)
	case *MsgBurnToken:
		w.Write(getMagicBytes("MsgBurnToken"))
		return EncodeMsgBurnToken(w, *v)
	case MsgCancelOrder:
		w.Write(getMagicBytes("MsgCancelOrder"))
		return EncodeMsgCancelOrder(w, v)
	case *MsgCancelOrder:
		w.Write(getMagicBytes("MsgCancelOrder"))
		return EncodeMsgCancelOrder(w, *v)
	case MsgCancelTradingPair:
		w.Write(getMagicBytes("MsgCancelTradingPair"))
		return EncodeMsgCancelTradingPair(w, v)
	case *MsgCancelTradingPair:
		w.Write(getMagicBytes("MsgCancelTradingPair"))
		return EncodeMsgCancelTradingPair(w, *v)
	case MsgCommentToken:
		w.Write(getMagicBytes("MsgCommentToken"))
		return EncodeMsgCommentToken(w, v)
	case *MsgCommentToken:
		w.Write(getMagicBytes("MsgCommentToken"))
		return EncodeMsgCommentToken(w, *v)
	case MsgCreateOrder:
		w.Write(getMagicBytes("MsgCreateOrder"))
		return EncodeMsgCreateOrder(w, v)
	case *MsgCreateOrder:
		w.Write(getMagicBytes("MsgCreateOrder"))
		return EncodeMsgCreateOrder(w, *v)
	case MsgCreateTradingPair:
		w.Write(getMagicBytes("MsgCreateTradingPair"))
		return EncodeMsgCreateTradingPair(w, v)
	case *MsgCreateTradingPair:
		w.Write(getMagicBytes("MsgCreateTradingPair"))
		return EncodeMsgCreateTradingPair(w, *v)
	case MsgCreateValidator:
		w.Write(getMagicBytes("MsgCreateValidator"))
		return EncodeMsgCreateValidator(w, v)
	case *MsgCreateValidator:
		w.Write(getMagicBytes("MsgCreateValidator"))
		return EncodeMsgCreateValidator(w, *v)
	case MsgDelegate:
		w.Write(getMagicBytes("MsgDelegate"))
		return EncodeMsgDelegate(w, v)
	case *MsgDelegate:
		w.Write(getMagicBytes("MsgDelegate"))
		return EncodeMsgDelegate(w, *v)
	case MsgDeposit:
		w.Write(getMagicBytes("MsgDeposit"))
		return EncodeMsgDeposit(w, v)
	case *MsgDeposit:
		w.Write(getMagicBytes("MsgDeposit"))
		return EncodeMsgDeposit(w, *v)
	case MsgDonateToCommunityPool:
		w.Write(getMagicBytes("MsgDonateToCommunityPool"))
		return EncodeMsgDonateToCommunityPool(w, v)
	case *MsgDonateToCommunityPool:
		w.Write(getMagicBytes("MsgDonateToCommunityPool"))
		return EncodeMsgDonateToCommunityPool(w, *v)
	case MsgEditValidator:
		w.Write(getMagicBytes("MsgEditValidator"))
		return EncodeMsgEditValidator(w, v)
	case *MsgEditValidator:
		w.Write(getMagicBytes("MsgEditValidator"))
		return EncodeMsgEditValidator(w, *v)
	case MsgForbidAddr:
		w.Write(getMagicBytes("MsgForbidAddr"))
		return EncodeMsgForbidAddr(w, v)
	case *MsgForbidAddr:
		w.Write(getMagicBytes("MsgForbidAddr"))
		return EncodeMsgForbidAddr(w, *v)
	case MsgForbidToken:
		w.Write(getMagicBytes("MsgForbidToken"))
		return EncodeMsgForbidToken(w, v)
	case *MsgForbidToken:
		w.Write(getMagicBytes("MsgForbidToken"))
		return EncodeMsgForbidToken(w, *v)
	case MsgIssueToken:
		w.Write(getMagicBytes("MsgIssueToken"))
		return EncodeMsgIssueToken(w, v)
	case *MsgIssueToken:
		w.Write(getMagicBytes("MsgIssueToken"))
		return EncodeMsgIssueToken(w, *v)
	case MsgMintToken:
		w.Write(getMagicBytes("MsgMintToken"))
		return EncodeMsgMintToken(w, v)
	case *MsgMintToken:
		w.Write(getMagicBytes("MsgMintToken"))
		return EncodeMsgMintToken(w, *v)
	case MsgModifyPricePrecision:
		w.Write(getMagicBytes("MsgModifyPricePrecision"))
		return EncodeMsgModifyPricePrecision(w, v)
	case *MsgModifyPricePrecision:
		w.Write(getMagicBytes("MsgModifyPricePrecision"))
		return EncodeMsgModifyPricePrecision(w, *v)
	case MsgModifyTokenInfo:
		w.Write(getMagicBytes("MsgModifyTokenInfo"))
		return EncodeMsgModifyTokenInfo(w, v)
	case *MsgModifyTokenInfo:
		w.Write(getMagicBytes("MsgModifyTokenInfo"))
		return EncodeMsgModifyTokenInfo(w, *v)
	case MsgMultiSend:
		w.Write(getMagicBytes("MsgMultiSend"))
		return EncodeMsgMultiSend(w, v)
	case *MsgMultiSend:
		w.Write(getMagicBytes("MsgMultiSend"))
		return EncodeMsgMultiSend(w, *v)
	case MsgMultiSendX:
		w.Write(getMagicBytes("MsgMultiSendX"))
		return EncodeMsgMultiSendX(w, v)
	case *MsgMultiSendX:
		w.Write(getMagicBytes("MsgMultiSendX"))
		return EncodeMsgMultiSendX(w, *v)
	case MsgRemoveTokenWhitelist:
		w.Write(getMagicBytes("MsgRemoveTokenWhitelist"))
		return EncodeMsgRemoveTokenWhitelist(w, v)
	case *MsgRemoveTokenWhitelist:
		w.Write(getMagicBytes("MsgRemoveTokenWhitelist"))
		return EncodeMsgRemoveTokenWhitelist(w, *v)
	case MsgSend:
		w.Write(getMagicBytes("MsgSend"))
		return EncodeMsgSend(w, v)
	case *MsgSend:
		w.Write(getMagicBytes("MsgSend"))
		return EncodeMsgSend(w, *v)
	case MsgSendX:
		w.Write(getMagicBytes("MsgSendX"))
		return EncodeMsgSendX(w, v)
	case *MsgSendX:
		w.Write(getMagicBytes("MsgSendX"))
		return EncodeMsgSendX(w, *v)
	case MsgSetMemoRequired:
		w.Write(getMagicBytes("MsgSetMemoRequired"))
		return EncodeMsgSetMemoRequired(w, v)
	case *MsgSetMemoRequired:
		w.Write(getMagicBytes("MsgSetMemoRequired"))
		return EncodeMsgSetMemoRequired(w, *v)
	case MsgSetWithdrawAddress:
		w.Write(getMagicBytes("MsgSetWithdrawAddress"))
		return EncodeMsgSetWithdrawAddress(w, v)
	case *MsgSetWithdrawAddress:
		w.Write(getMagicBytes("MsgSetWithdrawAddress"))
		return EncodeMsgSetWithdrawAddress(w, *v)
	case MsgSubmitProposal:
		w.Write(getMagicBytes("MsgSubmitProposal"))
		return EncodeMsgSubmitProposal(w, v)
	case *MsgSubmitProposal:
		w.Write(getMagicBytes("MsgSubmitProposal"))
		return EncodeMsgSubmitProposal(w, *v)
	case MsgTransferOwnership:
		w.Write(getMagicBytes("MsgTransferOwnership"))
		return EncodeMsgTransferOwnership(w, v)
	case *MsgTransferOwnership:
		w.Write(getMagicBytes("MsgTransferOwnership"))
		return EncodeMsgTransferOwnership(w, *v)
	case MsgUnForbidAddr:
		w.Write(getMagicBytes("MsgUnForbidAddr"))
		return EncodeMsgUnForbidAddr(w, v)
	case *MsgUnForbidAddr:
		w.Write(getMagicBytes("MsgUnForbidAddr"))
		return EncodeMsgUnForbidAddr(w, *v)
	case MsgUnForbidToken:
		w.Write(getMagicBytes("MsgUnForbidToken"))
		return EncodeMsgUnForbidToken(w, v)
	case *MsgUnForbidToken:
		w.Write(getMagicBytes("MsgUnForbidToken"))
		return EncodeMsgUnForbidToken(w, *v)
	case MsgUndelegate:
		w.Write(getMagicBytes("MsgUndelegate"))
		return EncodeMsgUndelegate(w, v)
	case *MsgUndelegate:
		w.Write(getMagicBytes("MsgUndelegate"))
		return EncodeMsgUndelegate(w, *v)
	case MsgUnjail:
		w.Write(getMagicBytes("MsgUnjail"))
		return EncodeMsgUnjail(w, v)
	case *MsgUnjail:
		w.Write(getMagicBytes("MsgUnjail"))
		return EncodeMsgUnjail(w, *v)
	case MsgVerifyInvariant:
		w.Write(getMagicBytes("MsgVerifyInvariant"))
		return EncodeMsgVerifyInvariant(w, v)
	case *MsgVerifyInvariant:
		w.Write(getMagicBytes("MsgVerifyInvariant"))
		return EncodeMsgVerifyInvariant(w, *v)
	case MsgVote:
		w.Write(getMagicBytes("MsgVote"))
		return EncodeMsgVote(w, v)
	case *MsgVote:
		w.Write(getMagicBytes("MsgVote"))
		return EncodeMsgVote(w, *v)
	case MsgWithdrawDelegatorReward:
		w.Write(getMagicBytes("MsgWithdrawDelegatorReward"))
		return EncodeMsgWithdrawDelegatorReward(w, v)
	case *MsgWithdrawDelegatorReward:
		w.Write(getMagicBytes("MsgWithdrawDelegatorReward"))
		return EncodeMsgWithdrawDelegatorReward(w, *v)
	case MsgWithdrawValidatorCommission:
		w.Write(getMagicBytes("MsgWithdrawValidatorCommission"))
		return EncodeMsgWithdrawValidatorCommission(w, v)
	case *MsgWithdrawValidatorCommission:
		w.Write(getMagicBytes("MsgWithdrawValidatorCommission"))
		return EncodeMsgWithdrawValidatorCommission(w, *v)
	case Order:
		w.Write(getMagicBytes("Order"))
		return EncodeOrder(w, v)
	case *Order:
		w.Write(getMagicBytes("Order"))
		return EncodeOrder(w, *v)
	case Output:
		w.Write(getMagicBytes("Output"))
		return EncodeOutput(w, v)
	case *Output:
		w.Write(getMagicBytes("Output"))
		return EncodeOutput(w, *v)
	case ParamChange:
		w.Write(getMagicBytes("ParamChange"))
		return EncodeParamChange(w, v)
	case *ParamChange:
		w.Write(getMagicBytes("ParamChange"))
		return EncodeParamChange(w, *v)
	case ParameterChangeProposal:
		w.Write(getMagicBytes("ParameterChangeProposal"))
		return EncodeParameterChangeProposal(w, v)
	case *ParameterChangeProposal:
		w.Write(getMagicBytes("ParameterChangeProposal"))
		return EncodeParameterChangeProposal(w, *v)
	case PrivKeyEd25519:
		w.Write(getMagicBytes("PrivKeyEd25519"))
		return EncodePrivKeyEd25519(w, v)
	case *PrivKeyEd25519:
		w.Write(getMagicBytes("PrivKeyEd25519"))
		return EncodePrivKeyEd25519(w, *v)
	case PrivKeySecp256k1:
		w.Write(getMagicBytes("PrivKeySecp256k1"))
		return EncodePrivKeySecp256k1(w, v)
	case *PrivKeySecp256k1:
		w.Write(getMagicBytes("PrivKeySecp256k1"))
		return EncodePrivKeySecp256k1(w, *v)
	case PubKeyEd25519:
		w.Write(getMagicBytes("PubKeyEd25519"))
		return EncodePubKeyEd25519(w, v)
	case *PubKeyEd25519:
		w.Write(getMagicBytes("PubKeyEd25519"))
		return EncodePubKeyEd25519(w, *v)
	case PubKeyMultisigThreshold:
		w.Write(getMagicBytes("PubKeyMultisigThreshold"))
		return EncodePubKeyMultisigThreshold(w, v)
	case *PubKeyMultisigThreshold:
		w.Write(getMagicBytes("PubKeyMultisigThreshold"))
		return EncodePubKeyMultisigThreshold(w, *v)
	case PubKeySecp256k1:
		w.Write(getMagicBytes("PubKeySecp256k1"))
		return EncodePubKeySecp256k1(w, v)
	case *PubKeySecp256k1:
		w.Write(getMagicBytes("PubKeySecp256k1"))
		return EncodePubKeySecp256k1(w, *v)
	case SignedMsgType:
		w.Write(getMagicBytes("SignedMsgType"))
		return EncodeSignedMsgType(w, v)
	case *SignedMsgType:
		w.Write(getMagicBytes("SignedMsgType"))
		return EncodeSignedMsgType(w, *v)
	case SoftwareUpgradeProposal:
		w.Write(getMagicBytes("SoftwareUpgradeProposal"))
		return EncodeSoftwareUpgradeProposal(w, v)
	case *SoftwareUpgradeProposal:
		w.Write(getMagicBytes("SoftwareUpgradeProposal"))
		return EncodeSoftwareUpgradeProposal(w, *v)
	case State:
		w.Write(getMagicBytes("State"))
		return EncodeState(w, v)
	case *State:
		w.Write(getMagicBytes("State"))
		return EncodeState(w, *v)
	case StdSignature:
		w.Write(getMagicBytes("StdSignature"))
		return EncodeStdSignature(w, v)
	case *StdSignature:
		w.Write(getMagicBytes("StdSignature"))
		return EncodeStdSignature(w, *v)
	case StdTx:
		w.Write(getMagicBytes("StdTx"))
		return EncodeStdTx(w, v)
	case *StdTx:
		w.Write(getMagicBytes("StdTx"))
		return EncodeStdTx(w, *v)
	case Supply:
		w.Write(getMagicBytes("Supply"))
		return EncodeSupply(w, v)
	case *Supply:
		w.Write(getMagicBytes("Supply"))
		return EncodeSupply(w, *v)
	case TextProposal:
		w.Write(getMagicBytes("TextProposal"))
		return EncodeTextProposal(w, v)
	case *TextProposal:
		w.Write(getMagicBytes("TextProposal"))
		return EncodeTextProposal(w, *v)
	case Vote:
		w.Write(getMagicBytes("Vote"))
		return EncodeVote(w, v)
	case *Vote:
		w.Write(getMagicBytes("Vote"))
		return EncodeVote(w, *v)
	case VoteOption:
		w.Write(getMagicBytes("VoteOption"))
		return EncodeVoteOption(w, v)
	case *VoteOption:
		w.Write(getMagicBytes("VoteOption"))
		return EncodeVoteOption(w, *v)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func BareEncodeAny(w io.Writer, x interface{}) error {
	switch v := x.(type) {
	case AccAddress:
		return EncodeAccAddress(w, v)
	case *AccAddress:
		return EncodeAccAddress(w, *v)
	case AccountX:
		return EncodeAccountX(w, v)
	case *AccountX:
		return EncodeAccountX(w, *v)
	case BaseAccount:
		return EncodeBaseAccount(w, v)
	case *BaseAccount:
		return EncodeBaseAccount(w, *v)
	case BaseToken:
		return EncodeBaseToken(w, v)
	case *BaseToken:
		return EncodeBaseToken(w, *v)
	case BaseVestingAccount:
		return EncodeBaseVestingAccount(w, v)
	case *BaseVestingAccount:
		return EncodeBaseVestingAccount(w, *v)
	case Coin:
		return EncodeCoin(w, v)
	case *Coin:
		return EncodeCoin(w, *v)
	case CommentRef:
		return EncodeCommentRef(w, v)
	case *CommentRef:
		return EncodeCommentRef(w, *v)
	case CommunityPoolSpendProposal:
		return EncodeCommunityPoolSpendProposal(w, v)
	case *CommunityPoolSpendProposal:
		return EncodeCommunityPoolSpendProposal(w, *v)
	case ContinuousVestingAccount:
		return EncodeContinuousVestingAccount(w, v)
	case *ContinuousVestingAccount:
		return EncodeContinuousVestingAccount(w, *v)
	case DelayedVestingAccount:
		return EncodeDelayedVestingAccount(w, v)
	case *DelayedVestingAccount:
		return EncodeDelayedVestingAccount(w, *v)
	case DuplicateVoteEvidence:
		return EncodeDuplicateVoteEvidence(w, v)
	case *DuplicateVoteEvidence:
		return EncodeDuplicateVoteEvidence(w, *v)
	case Input:
		return EncodeInput(w, v)
	case *Input:
		return EncodeInput(w, *v)
	case LockedCoin:
		return EncodeLockedCoin(w, v)
	case *LockedCoin:
		return EncodeLockedCoin(w, *v)
	case MarketInfo:
		return EncodeMarketInfo(w, v)
	case *MarketInfo:
		return EncodeMarketInfo(w, *v)
	case ModuleAccount:
		return EncodeModuleAccount(w, v)
	case *ModuleAccount:
		return EncodeModuleAccount(w, *v)
	case MsgAddTokenWhitelist:
		return EncodeMsgAddTokenWhitelist(w, v)
	case *MsgAddTokenWhitelist:
		return EncodeMsgAddTokenWhitelist(w, *v)
	case MsgAliasUpdate:
		return EncodeMsgAliasUpdate(w, v)
	case *MsgAliasUpdate:
		return EncodeMsgAliasUpdate(w, *v)
	case MsgBancorCancel:
		return EncodeMsgBancorCancel(w, v)
	case *MsgBancorCancel:
		return EncodeMsgBancorCancel(w, *v)
	case MsgBancorInit:
		return EncodeMsgBancorInit(w, v)
	case *MsgBancorInit:
		return EncodeMsgBancorInit(w, *v)
	case MsgBancorTrade:
		return EncodeMsgBancorTrade(w, v)
	case *MsgBancorTrade:
		return EncodeMsgBancorTrade(w, *v)
	case MsgBeginRedelegate:
		return EncodeMsgBeginRedelegate(w, v)
	case *MsgBeginRedelegate:
		return EncodeMsgBeginRedelegate(w, *v)
	case MsgBurnToken:
		return EncodeMsgBurnToken(w, v)
	case *MsgBurnToken:
		return EncodeMsgBurnToken(w, *v)
	case MsgCancelOrder:
		return EncodeMsgCancelOrder(w, v)
	case *MsgCancelOrder:
		return EncodeMsgCancelOrder(w, *v)
	case MsgCancelTradingPair:
		return EncodeMsgCancelTradingPair(w, v)
	case *MsgCancelTradingPair:
		return EncodeMsgCancelTradingPair(w, *v)
	case MsgCommentToken:
		return EncodeMsgCommentToken(w, v)
	case *MsgCommentToken:
		return EncodeMsgCommentToken(w, *v)
	case MsgCreateOrder:
		return EncodeMsgCreateOrder(w, v)
	case *MsgCreateOrder:
		return EncodeMsgCreateOrder(w, *v)
	case MsgCreateTradingPair:
		return EncodeMsgCreateTradingPair(w, v)
	case *MsgCreateTradingPair:
		return EncodeMsgCreateTradingPair(w, *v)
	case MsgCreateValidator:
		return EncodeMsgCreateValidator(w, v)
	case *MsgCreateValidator:
		return EncodeMsgCreateValidator(w, *v)
	case MsgDelegate:
		return EncodeMsgDelegate(w, v)
	case *MsgDelegate:
		return EncodeMsgDelegate(w, *v)
	case MsgDeposit:
		return EncodeMsgDeposit(w, v)
	case *MsgDeposit:
		return EncodeMsgDeposit(w, *v)
	case MsgDonateToCommunityPool:
		return EncodeMsgDonateToCommunityPool(w, v)
	case *MsgDonateToCommunityPool:
		return EncodeMsgDonateToCommunityPool(w, *v)
	case MsgEditValidator:
		return EncodeMsgEditValidator(w, v)
	case *MsgEditValidator:
		return EncodeMsgEditValidator(w, *v)
	case MsgForbidAddr:
		return EncodeMsgForbidAddr(w, v)
	case *MsgForbidAddr:
		return EncodeMsgForbidAddr(w, *v)
	case MsgForbidToken:
		return EncodeMsgForbidToken(w, v)
	case *MsgForbidToken:
		return EncodeMsgForbidToken(w, *v)
	case MsgIssueToken:
		return EncodeMsgIssueToken(w, v)
	case *MsgIssueToken:
		return EncodeMsgIssueToken(w, *v)
	case MsgMintToken:
		return EncodeMsgMintToken(w, v)
	case *MsgMintToken:
		return EncodeMsgMintToken(w, *v)
	case MsgModifyPricePrecision:
		return EncodeMsgModifyPricePrecision(w, v)
	case *MsgModifyPricePrecision:
		return EncodeMsgModifyPricePrecision(w, *v)
	case MsgModifyTokenInfo:
		return EncodeMsgModifyTokenInfo(w, v)
	case *MsgModifyTokenInfo:
		return EncodeMsgModifyTokenInfo(w, *v)
	case MsgMultiSend:
		return EncodeMsgMultiSend(w, v)
	case *MsgMultiSend:
		return EncodeMsgMultiSend(w, *v)
	case MsgMultiSendX:
		return EncodeMsgMultiSendX(w, v)
	case *MsgMultiSendX:
		return EncodeMsgMultiSendX(w, *v)
	case MsgRemoveTokenWhitelist:
		return EncodeMsgRemoveTokenWhitelist(w, v)
	case *MsgRemoveTokenWhitelist:
		return EncodeMsgRemoveTokenWhitelist(w, *v)
	case MsgSend:
		return EncodeMsgSend(w, v)
	case *MsgSend:
		return EncodeMsgSend(w, *v)
	case MsgSendX:
		return EncodeMsgSendX(w, v)
	case *MsgSendX:
		return EncodeMsgSendX(w, *v)
	case MsgSetMemoRequired:
		return EncodeMsgSetMemoRequired(w, v)
	case *MsgSetMemoRequired:
		return EncodeMsgSetMemoRequired(w, *v)
	case MsgSetWithdrawAddress:
		return EncodeMsgSetWithdrawAddress(w, v)
	case *MsgSetWithdrawAddress:
		return EncodeMsgSetWithdrawAddress(w, *v)
	case MsgSubmitProposal:
		return EncodeMsgSubmitProposal(w, v)
	case *MsgSubmitProposal:
		return EncodeMsgSubmitProposal(w, *v)
	case MsgTransferOwnership:
		return EncodeMsgTransferOwnership(w, v)
	case *MsgTransferOwnership:
		return EncodeMsgTransferOwnership(w, *v)
	case MsgUnForbidAddr:
		return EncodeMsgUnForbidAddr(w, v)
	case *MsgUnForbidAddr:
		return EncodeMsgUnForbidAddr(w, *v)
	case MsgUnForbidToken:
		return EncodeMsgUnForbidToken(w, v)
	case *MsgUnForbidToken:
		return EncodeMsgUnForbidToken(w, *v)
	case MsgUndelegate:
		return EncodeMsgUndelegate(w, v)
	case *MsgUndelegate:
		return EncodeMsgUndelegate(w, *v)
	case MsgUnjail:
		return EncodeMsgUnjail(w, v)
	case *MsgUnjail:
		return EncodeMsgUnjail(w, *v)
	case MsgVerifyInvariant:
		return EncodeMsgVerifyInvariant(w, v)
	case *MsgVerifyInvariant:
		return EncodeMsgVerifyInvariant(w, *v)
	case MsgVote:
		return EncodeMsgVote(w, v)
	case *MsgVote:
		return EncodeMsgVote(w, *v)
	case MsgWithdrawDelegatorReward:
		return EncodeMsgWithdrawDelegatorReward(w, v)
	case *MsgWithdrawDelegatorReward:
		return EncodeMsgWithdrawDelegatorReward(w, *v)
	case MsgWithdrawValidatorCommission:
		return EncodeMsgWithdrawValidatorCommission(w, v)
	case *MsgWithdrawValidatorCommission:
		return EncodeMsgWithdrawValidatorCommission(w, *v)
	case Order:
		return EncodeOrder(w, v)
	case *Order:
		return EncodeOrder(w, *v)
	case Output:
		return EncodeOutput(w, v)
	case *Output:
		return EncodeOutput(w, *v)
	case ParamChange:
		return EncodeParamChange(w, v)
	case *ParamChange:
		return EncodeParamChange(w, *v)
	case ParameterChangeProposal:
		return EncodeParameterChangeProposal(w, v)
	case *ParameterChangeProposal:
		return EncodeParameterChangeProposal(w, *v)
	case PrivKeyEd25519:
		return EncodePrivKeyEd25519(w, v)
	case *PrivKeyEd25519:
		return EncodePrivKeyEd25519(w, *v)
	case PrivKeySecp256k1:
		return EncodePrivKeySecp256k1(w, v)
	case *PrivKeySecp256k1:
		return EncodePrivKeySecp256k1(w, *v)
	case PubKeyEd25519:
		return EncodePubKeyEd25519(w, v)
	case *PubKeyEd25519:
		return EncodePubKeyEd25519(w, *v)
	case PubKeyMultisigThreshold:
		return EncodePubKeyMultisigThreshold(w, v)
	case *PubKeyMultisigThreshold:
		return EncodePubKeyMultisigThreshold(w, *v)
	case PubKeySecp256k1:
		return EncodePubKeySecp256k1(w, v)
	case *PubKeySecp256k1:
		return EncodePubKeySecp256k1(w, *v)
	case SignedMsgType:
		return EncodeSignedMsgType(w, v)
	case *SignedMsgType:
		return EncodeSignedMsgType(w, *v)
	case SoftwareUpgradeProposal:
		return EncodeSoftwareUpgradeProposal(w, v)
	case *SoftwareUpgradeProposal:
		return EncodeSoftwareUpgradeProposal(w, *v)
	case State:
		return EncodeState(w, v)
	case *State:
		return EncodeState(w, *v)
	case StdSignature:
		return EncodeStdSignature(w, v)
	case *StdSignature:
		return EncodeStdSignature(w, *v)
	case StdTx:
		return EncodeStdTx(w, v)
	case *StdTx:
		return EncodeStdTx(w, *v)
	case Supply:
		return EncodeSupply(w, v)
	case *Supply:
		return EncodeSupply(w, *v)
	case TextProposal:
		return EncodeTextProposal(w, v)
	case *TextProposal:
		return EncodeTextProposal(w, *v)
	case Vote:
		return EncodeVote(w, v)
	case *Vote:
		return EncodeVote(w, *v)
	case VoteOption:
		return EncodeVoteOption(w, v)
	case *VoteOption:
		return EncodeVoteOption(w, *v)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func DecodeAny(bz []byte) (interface{}, int, error) {
	var v interface{}
	var magicBytes [4]byte
	var n int
	for i := 0; i < 4; i++ {
		magicBytes[i] = bz[i]
	}
	switch magicBytes {
	case [4]byte{37, 50, 37, 208}:
		v, n, err := DecodeAccAddress(bz[4:])
		return v, n + 4, err
	case [4]byte{52, 40, 250, 245}:
		v, n, err := DecodeAccountX(bz[4:])
		return v, n + 4, err
	case [4]byte{100, 94, 81, 72}:
		v, n, err := DecodeBaseAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{34, 178, 244, 51}:
		v, n, err := DecodeBaseToken(bz[4:])
		return v, n + 4, err
	case [4]byte{178, 47, 121, 129}:
		v, n, err := DecodeBaseVestingAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{141, 140, 97, 80}:
		v, n, err := DecodeCoin(bz[4:])
		return v, n + 4, err
	case [4]byte{17, 162, 164, 235}:
		v, n, err := DecodeCommentRef(bz[4:])
		return v, n + 4, err
	case [4]byte{37, 214, 119, 170}:
		v, n, err := DecodeCommunityPoolSpendProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{95, 96, 75, 0}:
		v, n, err := DecodeContinuousVestingAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{187, 71, 224, 1}:
		v, n, err := DecodeDelayedVestingAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{130, 76, 198, 17}:
		v, n, err := DecodeDuplicateVoteEvidence(bz[4:])
		return v, n + 4, err
	case [4]byte{165, 152, 189, 47}:
		v, n, err := DecodeInput(bz[4:])
		return v, n + 4, err
	case [4]byte{227, 236, 168, 93}:
		v, n, err := DecodeLockedCoin(bz[4:])
		return v, n + 4, err
	case [4]byte{174, 117, 167, 230}:
		v, n, err := DecodeMarketInfo(bz[4:])
		return v, n + 4, err
	case [4]byte{190, 107, 1, 124}:
		v, n, err := DecodeModuleAccount(bz[4:])
		return v, n + 4, err
	case [4]byte{147, 136, 220, 215}:
		v, n, err := DecodeMsgAddTokenWhitelist(bz[4:])
		return v, n + 4, err
	case [4]byte{173, 181, 17, 162}:
		v, n, err := DecodeMsgAliasUpdate(bz[4:])
		return v, n + 4, err
	case [4]byte{187, 190, 104, 91}:
		v, n, err := DecodeMsgBancorCancel(bz[4:])
		return v, n + 4, err
	case [4]byte{177, 152, 232, 60}:
		v, n, err := DecodeMsgBancorInit(bz[4:])
		return v, n + 4, err
	case [4]byte{225, 122, 18, 80}:
		v, n, err := DecodeMsgBancorTrade(bz[4:])
		return v, n + 4, err
	case [4]byte{247, 3, 0, 105}:
		v, n, err := DecodeMsgBeginRedelegate(bz[4:])
		return v, n + 4, err
	case [4]byte{228, 0, 236, 212}:
		v, n, err := DecodeMsgBurnToken(bz[4:])
		return v, n + 4, err
	case [4]byte{106, 229, 80, 141}:
		v, n, err := DecodeMsgCancelOrder(bz[4:])
		return v, n + 4, err
	case [4]byte{13, 177, 95, 127}:
		v, n, err := DecodeMsgCancelTradingPair(bz[4:])
		return v, n + 4, err
	case [4]byte{79, 125, 235, 121}:
		v, n, err := DecodeMsgCommentToken(bz[4:])
		return v, n + 4, err
	case [4]byte{246, 238, 7, 164}:
		v, n, err := DecodeMsgCreateOrder(bz[4:])
		return v, n + 4, err
	case [4]byte{130, 221, 55, 57}:
		v, n, err := DecodeMsgCreateTradingPair(bz[4:])
		return v, n + 4, err
	case [4]byte{58, 78, 252, 114}:
		v, n, err := DecodeMsgCreateValidator(bz[4:])
		return v, n + 4, err
	case [4]byte{1, 82, 140, 71}:
		v, n, err := DecodeMsgDelegate(bz[4:])
		return v, n + 4, err
	case [4]byte{205, 134, 140, 190}:
		v, n, err := DecodeMsgDeposit(bz[4:])
		return v, n + 4, err
	case [4]byte{20, 250, 115, 197}:
		v, n, err := DecodeMsgDonateToCommunityPool(bz[4:])
		return v, n + 4, err
	case [4]byte{202, 62, 140, 8}:
		v, n, err := DecodeMsgEditValidator(bz[4:])
		return v, n + 4, err
	case [4]byte{105, 235, 112, 10}:
		v, n, err := DecodeMsgForbidAddr(bz[4:])
		return v, n + 4, err
	case [4]byte{36, 174, 203, 238}:
		v, n, err := DecodeMsgForbidToken(bz[4:])
		return v, n + 4, err
	case [4]byte{233, 180, 92, 129}:
		v, n, err := DecodeMsgIssueToken(bz[4:])
		return v, n + 4, err
	case [4]byte{66, 148, 56, 203}:
		v, n, err := DecodeMsgMintToken(bz[4:])
		return v, n + 4, err
	case [4]byte{76, 91, 156, 199}:
		v, n, err := DecodeMsgModifyPricePrecision(bz[4:])
		return v, n + 4, err
	case [4]byte{40, 158, 40, 18}:
		v, n, err := DecodeMsgModifyTokenInfo(bz[4:])
		return v, n + 4, err
	case [4]byte{207, 152, 156, 90}:
		v, n, err := DecodeMsgMultiSend(bz[4:])
		return v, n + 4, err
	case [4]byte{61, 117, 88, 200}:
		v, n, err := DecodeMsgMultiSendX(bz[4:])
		return v, n + 4, err
	case [4]byte{44, 154, 68, 83}:
		v, n, err := DecodeMsgRemoveTokenWhitelist(bz[4:])
		return v, n + 4, err
	case [4]byte{100, 168, 39, 140}:
		v, n, err := DecodeMsgSend(bz[4:])
		return v, n + 4, err
	case [4]byte{198, 76, 8, 81}:
		v, n, err := DecodeMsgSendX(bz[4:])
		return v, n + 4, err
	case [4]byte{184, 238, 253, 154}:
		v, n, err := DecodeMsgSetMemoRequired(bz[4:])
		return v, n + 4, err
	case [4]byte{190, 178, 173, 144}:
		v, n, err := DecodeMsgSetWithdrawAddress(bz[4:])
		return v, n + 4, err
	case [4]byte{115, 119, 137, 48}:
		v, n, err := DecodeMsgSubmitProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{200, 224, 118, 175}:
		v, n, err := DecodeMsgTransferOwnership(bz[4:])
		return v, n + 4, err
	case [4]byte{167, 165, 166, 227}:
		v, n, err := DecodeMsgUnForbidAddr(bz[4:])
		return v, n + 4, err
	case [4]byte{78, 83, 156, 139}:
		v, n, err := DecodeMsgUnForbidToken(bz[4:])
		return v, n + 4, err
	case [4]byte{122, 66, 160, 76}:
		v, n, err := DecodeMsgUndelegate(bz[4:])
		return v, n + 4, err
	case [4]byte{216, 247, 180, 46}:
		v, n, err := DecodeMsgUnjail(bz[4:])
		return v, n + 4, err
	case [4]byte{84, 44, 219, 65}:
		v, n, err := DecodeMsgVerifyInvariant(bz[4:])
		return v, n + 4, err
	case [4]byte{238, 246, 67, 141}:
		v, n, err := DecodeMsgVote(bz[4:])
		return v, n + 4, err
	case [4]byte{94, 251, 176, 152}:
		v, n, err := DecodeMsgWithdrawDelegatorReward(bz[4:])
		return v, n + 4, err
	case [4]byte{18, 172, 190, 152}:
		v, n, err := DecodeMsgWithdrawValidatorCommission(bz[4:])
		return v, n + 4, err
	case [4]byte{57, 115, 63, 27}:
		v, n, err := DecodeOrder(bz[4:])
		return v, n + 4, err
	case [4]byte{251, 0, 54, 127}:
		v, n, err := DecodeOutput(bz[4:])
		return v, n + 4, err
	case [4]byte{234, 101, 49, 27}:
		v, n, err := DecodeParamChange(bz[4:])
		return v, n + 4, err
	case [4]byte{166, 63, 172, 210}:
		v, n, err := DecodeParameterChangeProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{93, 160, 108, 51}:
		v, n, err := DecodePrivKeyEd25519(bz[4:])
		return v, n + 4, err
	case [4]byte{209, 107, 141, 98}:
		v, n, err := DecodePrivKeySecp256k1(bz[4:])
		return v, n + 4, err
	case [4]byte{108, 143, 2, 48}:
		v, n, err := DecodePubKeyEd25519(bz[4:])
		return v, n + 4, err
	case [4]byte{131, 227, 102, 173}:
		v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
		return v, n + 4, err
	case [4]byte{10, 126, 85, 105}:
		v, n, err := DecodePubKeySecp256k1(bz[4:])
		return v, n + 4, err
	case [4]byte{169, 174, 252, 87}:
		v, n, err := DecodeSignedMsgType(bz[4:])
		return v, n + 4, err
	case [4]byte{37, 100, 208, 251}:
		v, n, err := DecodeSoftwareUpgradeProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{186, 223, 120, 4}:
		v, n, err := DecodeState(bz[4:])
		return v, n + 4, err
	case [4]byte{88, 244, 106, 18}:
		v, n, err := DecodeStdSignature(bz[4:])
		return v, n + 4, err
	case [4]byte{71, 175, 179, 184}:
		v, n, err := DecodeStdTx(bz[4:])
		return v, n + 4, err
	case [4]byte{233, 209, 209, 86}:
		v, n, err := DecodeSupply(bz[4:])
		return v, n + 4, err
	case [4]byte{169, 32, 176, 245}:
		v, n, err := DecodeTextProposal(bz[4:])
		return v, n + 4, err
	case [4]byte{113, 227, 24, 224}:
		v, n, err := DecodeVote(bz[4:])
		return v, n + 4, err
	case [4]byte{56, 159, 20, 227}:
		v, n, err := DecodeVoteOption(bz[4:])
		return v, n + 4, err
	default:
		panic("Unknown type")
	} // end of switch
	return v, n, nil
} // end of DecodeAny
func BareDecodeAny(bz []byte, x interface{}) (n int, err error) {
	switch v := x.(type) {
	case *AccAddress:
		*v, n, err = DecodeAccAddress(bz)
	case *AccountX:
		*v, n, err = DecodeAccountX(bz)
	case *BaseAccount:
		*v, n, err = DecodeBaseAccount(bz)
	case *BaseToken:
		*v, n, err = DecodeBaseToken(bz)
	case *BaseVestingAccount:
		*v, n, err = DecodeBaseVestingAccount(bz)
	case *Coin:
		*v, n, err = DecodeCoin(bz)
	case *CommentRef:
		*v, n, err = DecodeCommentRef(bz)
	case *CommunityPoolSpendProposal:
		*v, n, err = DecodeCommunityPoolSpendProposal(bz)
	case *ContinuousVestingAccount:
		*v, n, err = DecodeContinuousVestingAccount(bz)
	case *DelayedVestingAccount:
		*v, n, err = DecodeDelayedVestingAccount(bz)
	case *DuplicateVoteEvidence:
		*v, n, err = DecodeDuplicateVoteEvidence(bz)
	case *Input:
		*v, n, err = DecodeInput(bz)
	case *LockedCoin:
		*v, n, err = DecodeLockedCoin(bz)
	case *MarketInfo:
		*v, n, err = DecodeMarketInfo(bz)
	case *ModuleAccount:
		*v, n, err = DecodeModuleAccount(bz)
	case *MsgAddTokenWhitelist:
		*v, n, err = DecodeMsgAddTokenWhitelist(bz)
	case *MsgAliasUpdate:
		*v, n, err = DecodeMsgAliasUpdate(bz)
	case *MsgBancorCancel:
		*v, n, err = DecodeMsgBancorCancel(bz)
	case *MsgBancorInit:
		*v, n, err = DecodeMsgBancorInit(bz)
	case *MsgBancorTrade:
		*v, n, err = DecodeMsgBancorTrade(bz)
	case *MsgBeginRedelegate:
		*v, n, err = DecodeMsgBeginRedelegate(bz)
	case *MsgBurnToken:
		*v, n, err = DecodeMsgBurnToken(bz)
	case *MsgCancelOrder:
		*v, n, err = DecodeMsgCancelOrder(bz)
	case *MsgCancelTradingPair:
		*v, n, err = DecodeMsgCancelTradingPair(bz)
	case *MsgCommentToken:
		*v, n, err = DecodeMsgCommentToken(bz)
	case *MsgCreateOrder:
		*v, n, err = DecodeMsgCreateOrder(bz)
	case *MsgCreateTradingPair:
		*v, n, err = DecodeMsgCreateTradingPair(bz)
	case *MsgCreateValidator:
		*v, n, err = DecodeMsgCreateValidator(bz)
	case *MsgDelegate:
		*v, n, err = DecodeMsgDelegate(bz)
	case *MsgDeposit:
		*v, n, err = DecodeMsgDeposit(bz)
	case *MsgDonateToCommunityPool:
		*v, n, err = DecodeMsgDonateToCommunityPool(bz)
	case *MsgEditValidator:
		*v, n, err = DecodeMsgEditValidator(bz)
	case *MsgForbidAddr:
		*v, n, err = DecodeMsgForbidAddr(bz)
	case *MsgForbidToken:
		*v, n, err = DecodeMsgForbidToken(bz)
	case *MsgIssueToken:
		*v, n, err = DecodeMsgIssueToken(bz)
	case *MsgMintToken:
		*v, n, err = DecodeMsgMintToken(bz)
	case *MsgModifyPricePrecision:
		*v, n, err = DecodeMsgModifyPricePrecision(bz)
	case *MsgModifyTokenInfo:
		*v, n, err = DecodeMsgModifyTokenInfo(bz)
	case *MsgMultiSend:
		*v, n, err = DecodeMsgMultiSend(bz)
	case *MsgMultiSendX:
		*v, n, err = DecodeMsgMultiSendX(bz)
	case *MsgRemoveTokenWhitelist:
		*v, n, err = DecodeMsgRemoveTokenWhitelist(bz)
	case *MsgSend:
		*v, n, err = DecodeMsgSend(bz)
	case *MsgSendX:
		*v, n, err = DecodeMsgSendX(bz)
	case *MsgSetMemoRequired:
		*v, n, err = DecodeMsgSetMemoRequired(bz)
	case *MsgSetWithdrawAddress:
		*v, n, err = DecodeMsgSetWithdrawAddress(bz)
	case *MsgSubmitProposal:
		*v, n, err = DecodeMsgSubmitProposal(bz)
	case *MsgTransferOwnership:
		*v, n, err = DecodeMsgTransferOwnership(bz)
	case *MsgUnForbidAddr:
		*v, n, err = DecodeMsgUnForbidAddr(bz)
	case *MsgUnForbidToken:
		*v, n, err = DecodeMsgUnForbidToken(bz)
	case *MsgUndelegate:
		*v, n, err = DecodeMsgUndelegate(bz)
	case *MsgUnjail:
		*v, n, err = DecodeMsgUnjail(bz)
	case *MsgVerifyInvariant:
		*v, n, err = DecodeMsgVerifyInvariant(bz)
	case *MsgVote:
		*v, n, err = DecodeMsgVote(bz)
	case *MsgWithdrawDelegatorReward:
		*v, n, err = DecodeMsgWithdrawDelegatorReward(bz)
	case *MsgWithdrawValidatorCommission:
		*v, n, err = DecodeMsgWithdrawValidatorCommission(bz)
	case *Order:
		*v, n, err = DecodeOrder(bz)
	case *Output:
		*v, n, err = DecodeOutput(bz)
	case *ParamChange:
		*v, n, err = DecodeParamChange(bz)
	case *ParameterChangeProposal:
		*v, n, err = DecodeParameterChangeProposal(bz)
	case *PrivKeyEd25519:
		*v, n, err = DecodePrivKeyEd25519(bz)
	case *PrivKeySecp256k1:
		*v, n, err = DecodePrivKeySecp256k1(bz)
	case *PubKeyEd25519:
		*v, n, err = DecodePubKeyEd25519(bz)
	case *PubKeyMultisigThreshold:
		*v, n, err = DecodePubKeyMultisigThreshold(bz)
	case *PubKeySecp256k1:
		*v, n, err = DecodePubKeySecp256k1(bz)
	case *SignedMsgType:
		*v, n, err = DecodeSignedMsgType(bz)
	case *SoftwareUpgradeProposal:
		*v, n, err = DecodeSoftwareUpgradeProposal(bz)
	case *State:
		*v, n, err = DecodeState(bz)
	case *StdSignature:
		*v, n, err = DecodeStdSignature(bz)
	case *StdTx:
		*v, n, err = DecodeStdTx(bz)
	case *Supply:
		*v, n, err = DecodeSupply(bz)
	case *TextProposal:
		*v, n, err = DecodeTextProposal(bz)
	case *Vote:
		*v, n, err = DecodeVote(bz)
	case *VoteOption:
		*v, n, err = DecodeVoteOption(bz)
	default:
		panic("Unknown type")
	} // end of switch
	return
} // end of DecodeVar
func RandAny(r RandSrc) interface{} {
	switch r.GetUint() % 73 {
	case 0:
		return RandAccAddress(r)
	case 1:
		return RandAccountX(r)
	case 2:
		return RandBaseAccount(r)
	case 3:
		return RandBaseToken(r)
	case 4:
		return RandBaseVestingAccount(r)
	case 5:
		return RandCoin(r)
	case 6:
		return RandCommentRef(r)
	case 7:
		return RandCommunityPoolSpendProposal(r)
	case 8:
		return RandContinuousVestingAccount(r)
	case 9:
		return RandDelayedVestingAccount(r)
	case 10:
		return RandDuplicateVoteEvidence(r)
	case 11:
		return RandInput(r)
	case 12:
		return RandLockedCoin(r)
	case 13:
		return RandMarketInfo(r)
	case 14:
		return RandModuleAccount(r)
	case 15:
		return RandMsgAddTokenWhitelist(r)
	case 16:
		return RandMsgAliasUpdate(r)
	case 17:
		return RandMsgBancorCancel(r)
	case 18:
		return RandMsgBancorInit(r)
	case 19:
		return RandMsgBancorTrade(r)
	case 20:
		return RandMsgBeginRedelegate(r)
	case 21:
		return RandMsgBurnToken(r)
	case 22:
		return RandMsgCancelOrder(r)
	case 23:
		return RandMsgCancelTradingPair(r)
	case 24:
		return RandMsgCommentToken(r)
	case 25:
		return RandMsgCreateOrder(r)
	case 26:
		return RandMsgCreateTradingPair(r)
	case 27:
		return RandMsgCreateValidator(r)
	case 28:
		return RandMsgDelegate(r)
	case 29:
		return RandMsgDeposit(r)
	case 30:
		return RandMsgDonateToCommunityPool(r)
	case 31:
		return RandMsgEditValidator(r)
	case 32:
		return RandMsgForbidAddr(r)
	case 33:
		return RandMsgForbidToken(r)
	case 34:
		return RandMsgIssueToken(r)
	case 35:
		return RandMsgMintToken(r)
	case 36:
		return RandMsgModifyPricePrecision(r)
	case 37:
		return RandMsgModifyTokenInfo(r)
	case 38:
		return RandMsgMultiSend(r)
	case 39:
		return RandMsgMultiSendX(r)
	case 40:
		return RandMsgRemoveTokenWhitelist(r)
	case 41:
		return RandMsgSend(r)
	case 42:
		return RandMsgSendX(r)
	case 43:
		return RandMsgSetMemoRequired(r)
	case 44:
		return RandMsgSetWithdrawAddress(r)
	case 45:
		return RandMsgSubmitProposal(r)
	case 46:
		return RandMsgTransferOwnership(r)
	case 47:
		return RandMsgUnForbidAddr(r)
	case 48:
		return RandMsgUnForbidToken(r)
	case 49:
		return RandMsgUndelegate(r)
	case 50:
		return RandMsgUnjail(r)
	case 51:
		return RandMsgVerifyInvariant(r)
	case 52:
		return RandMsgVote(r)
	case 53:
		return RandMsgWithdrawDelegatorReward(r)
	case 54:
		return RandMsgWithdrawValidatorCommission(r)
	case 55:
		return RandOrder(r)
	case 56:
		return RandOutput(r)
	case 57:
		return RandParamChange(r)
	case 58:
		return RandParameterChangeProposal(r)
	case 59:
		return RandPrivKeyEd25519(r)
	case 60:
		return RandPrivKeySecp256k1(r)
	case 61:
		return RandPubKeyEd25519(r)
	case 62:
		return RandPubKeyMultisigThreshold(r)
	case 63:
		return RandPubKeySecp256k1(r)
	case 64:
		return RandSignedMsgType(r)
	case 65:
		return RandSoftwareUpgradeProposal(r)
	case 66:
		return RandState(r)
	case 67:
		return RandStdSignature(r)
	case 68:
		return RandStdTx(r)
	case 69:
		return RandSupply(r)
	case 70:
		return RandTextProposal(r)
	case 71:
		return RandVote(r)
	case 72:
		return RandVoteOption(r)
	default:
		panic("Unknown Type.")
	} // end of switch
} // end of func
func GetSupportList() []string {
	return []string{
		"github.com/coinexchain/cet-sdk/modules/alias/internal/types.MsgAliasUpdate",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.BaseToken",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgAddTokenWhitelist",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgBurnToken",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgForbidAddr",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgForbidToken",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgIssueToken",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgMintToken",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgModifyTokenInfo",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgRemoveTokenWhitelist",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgTransferOwnership",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgUnForbidAddr",
		"github.com/coinexchain/cet-sdk/modules/asset/internal/types.MsgUnForbidToken",
		"github.com/coinexchain/cet-sdk/modules/authx/internal/types.AccountX",
		"github.com/coinexchain/cet-sdk/modules/authx/internal/types.LockedCoin",
		"github.com/coinexchain/cet-sdk/modules/bancorlite/internal/types.MsgBancorCancel",
		"github.com/coinexchain/cet-sdk/modules/bancorlite/internal/types.MsgBancorInit",
		"github.com/coinexchain/cet-sdk/modules/bancorlite/internal/types.MsgBancorTrade",
		"github.com/coinexchain/cet-sdk/modules/bankx/internal/types.MsgMultiSend",
		"github.com/coinexchain/cet-sdk/modules/bankx/internal/types.MsgSend",
		"github.com/coinexchain/cet-sdk/modules/bankx/internal/types.MsgSetMemoRequired",
		"github.com/coinexchain/cet-sdk/modules/comment/internal/types.CommentRef",
		"github.com/coinexchain/cet-sdk/modules/comment/internal/types.MsgCommentToken",
		"github.com/coinexchain/cet-sdk/modules/distributionx/types.MsgDonateToCommunityPool",
		"github.com/coinexchain/cet-sdk/modules/incentive/internal/types.State",
		"github.com/coinexchain/cet-sdk/modules/market/internal/types.MarketInfo",
		"github.com/coinexchain/cet-sdk/modules/market/internal/types.MsgCancelOrder",
		"github.com/coinexchain/cet-sdk/modules/market/internal/types.MsgCancelTradingPair",
		"github.com/coinexchain/cet-sdk/modules/market/internal/types.MsgCreateOrder",
		"github.com/coinexchain/cet-sdk/modules/market/internal/types.MsgCreateTradingPair",
		"github.com/coinexchain/cet-sdk/modules/market/internal/types.MsgModifyPricePrecision",
		"github.com/coinexchain/cet-sdk/modules/market/internal/types.Order",
		"github.com/cosmos/cosmos-sdk/types.AccAddress",
		"github.com/cosmos/cosmos-sdk/types.Coin",
		"github.com/cosmos/cosmos-sdk/types.Msg",
		"github.com/cosmos/cosmos-sdk/x/auth/exported.Account",
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
		"github.com/cosmos/cosmos-sdk/x/supply/internal/types.ModuleAccount",
		"github.com/cosmos/cosmos-sdk/x/supply/internal/types.Supply",
		"github.com/tendermint/tendermint/crypto.PubKey",
		"github.com/tendermint/tendermint/crypto/ed25519.PrivKeyEd25519",
		"github.com/tendermint/tendermint/crypto/ed25519.PubKeyEd25519",
		"github.com/tendermint/tendermint/crypto/multisig.PubKeyMultisigThreshold",
		"github.com/tendermint/tendermint/crypto/secp256k1.PrivKeySecp256k1",
		"github.com/tendermint/tendermint/crypto/secp256k1.PubKeySecp256k1",
		"github.com/tendermint/tendermint/types.DuplicateVoteEvidence",
		"github.com/tendermint/tendermint/types.SignedMsgType",
		"github.com/tendermint/tendermint/types.Vote",
	}
} // end of GetSupportList
