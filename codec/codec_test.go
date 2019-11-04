package codec

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	amino "github.com/tendermint/go-amino"
)

func TestBasicOp(t *testing.T) {
	var s CodonStub

	u := uint64(1000)
	assert.Equal(t, amino.UvarintSize(u), s.UvarintSize(u))
	u = uint64(math.MaxInt64)
	assert.Equal(t, amino.UvarintSize(u), s.UvarintSize(u))
	i := int64(1000)
	assert.Equal(t, amino.VarintSize(i), s.VarintSize(i))
	i = int64(-1000)
	assert.Equal(t, amino.VarintSize(i), s.VarintSize(i))
	i = math.MaxInt64
	assert.Equal(t, amino.VarintSize(i), s.VarintSize(i))
	i = math.MinInt64
	assert.Equal(t, amino.VarintSize(i), s.VarintSize(i))

	w := new(bytes.Buffer)
	bz := []byte("Long long ago, there is a...")
	err := s.EncodeByteSlice(w, bz)
	assert.Equal(t, nil, err)
	bz2, _, err := s.DecodeByteSlice(w.Bytes())
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bytes.Equal(bz, bz2))

	assert.Equal(t, amino.ByteSliceSize(bz), s.ByteSliceSize(bz))

	w = new(bytes.Buffer)
	str := "Long long ago, there is a..."
	err = s.EncodeString(w, str)
	assert.Equal(t, nil, err)
	str2, _, err := s.DecodeString(w.Bytes())
	assert.Equal(t, nil, err)
	assert.Equal(t, str, str2)

	iList := []int64{1000, 10000000, 398837838343, math.MaxInt64, math.MinInt64, -10000000}
	for _, i := range iList {
		w = new(bytes.Buffer)
		err := s.EncodeVarint(w, i)
		assert.Equal(t, nil, err)
		i2, _, err := s.DecodeVarint(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, i, i2)
	}

	for _, i := range iList {
		w = new(bytes.Buffer)
		err := s.EncodeInt64(w, i)
		assert.Equal(t, nil, err)
		i2, _, err := s.DecodeInt64(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, i, i2)
	}

	for _, ii := range iList {
		i := int8(ii)
		w = new(bytes.Buffer)
		err := s.EncodeInt8(w, i)
		assert.Equal(t, nil, err)
		i2, _, err := s.DecodeInt8(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, i, i2)
	}

	for _, ii := range iList {
		i := int16(ii)
		w = new(bytes.Buffer)
		err := s.EncodeInt16(w, i)
		assert.Equal(t, nil, err)
		i2, _, err := s.DecodeInt16(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, i, i2)
	}

	for _, ii := range iList {
		i := int32(ii)
		w = new(bytes.Buffer)
		err := s.EncodeInt32(w, i)
		assert.Equal(t, nil, err)
		i2, _, err := s.DecodeInt32(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, i, i2)
	}

	uList := []uint64{1000, 10000000, 398837838343, uint64(math.MaxInt64), math.MaxUint64, math.MaxUint64 - 10000000}
	for _, u := range uList {
		w = new(bytes.Buffer)
		err := s.EncodeUvarint(w, u)
		assert.Equal(t, nil, err)
		u2, _, err := s.DecodeUvarint(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, u, u2)
	}

	for _, u := range uList {
		w = new(bytes.Buffer)
		err := s.EncodeUint64(w, u)
		assert.Equal(t, nil, err)
		u2, _, err := s.DecodeUint64(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, u, u2)
	}

	for _, uu := range uList {
		u := byte(uu)
		w = new(bytes.Buffer)
		err := s.EncodeByte(w, u)
		assert.Equal(t, nil, err)
		u2, _, err := s.DecodeByte(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, u, u2)
	}

	for _, uu := range uList {
		u := uint8(uu)
		w = new(bytes.Buffer)
		err := s.EncodeUint8(w, u)
		assert.Equal(t, nil, err)
		u2, _, err := s.DecodeUint8(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, u, u2)
	}

	for _, uu := range uList {
		u := uint16(uu)
		w = new(bytes.Buffer)
		err := s.EncodeUint16(w, u)
		assert.Equal(t, nil, err)
		u2, _, err := s.DecodeUint16(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, u, u2)
	}

	for _, uu := range uList {
		u := uint32(uu)
		w = new(bytes.Buffer)
		err := s.EncodeUint32(w, u)
		assert.Equal(t, nil, err)
		u2, _, err := s.DecodeUint32(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, u, u2)
	}

	for _, b := range []bool{true, false} {
		w = new(bytes.Buffer)
		err := s.EncodeBool(w, b)
		assert.Equal(t, nil, err)
		b2, _, err := s.DecodeBool(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, b, b2)
	}

	for _, uu := range uList {
		u := uint32(uu)
		f := math.Float32frombits(u)
		w = new(bytes.Buffer)
		err := s.EncodeFloat32(w, f)
		assert.Equal(t, nil, err)
		f2, _, err := s.DecodeFloat32(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, math.Float32bits(f), math.Float32bits(f2))
	}

	for _, u := range uList {
		f := math.Float64frombits(u)
		w = new(bytes.Buffer)
		err := s.EncodeFloat64(w, f)
		assert.Equal(t, nil, err)
		f2, _, err := s.DecodeFloat64(w.Bytes())
		assert.Equal(t, nil, err)
		assert.Equal(t, math.Float64bits(f), math.Float64bits(f2))
	}
}
