package testutil

import "io"

var _ io.Reader = &BytesReader{}

type BytesReader struct {
	bytes []byte
}

func NewBytesReader(s string) *BytesReader {
	return &BytesReader{bytes: []byte(s)}
}

func (r *BytesReader) Read(p []byte) (n int, err error) {
	want := len(p)
	have := len(r.bytes)

	if want == 0 {
		return 0, nil
	}
	if have == 0 {
		return 0, io.EOF
	}

	n = want
	if have < n {
		n = have
	}

	copy(p, r.bytes[:n])
	r.bytes = r.bytes[n:]
	return
}
