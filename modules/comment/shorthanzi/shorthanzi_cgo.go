// +build lz4usecgo

package shorthanzi

// #cgo CFLAGS: -O3
// #include "lz4.h"
import "C"

func CompressText(data string) ([]byte, bool) {
	buf := make([]byte, len(data)*2)

	n, err := CompressDefault([]byte(data), buf)
	if err != nil {
		return nil, false
	}
	if n >= len(data) || n == 0 {
		return nil, false
	}
	return buf[:n], true // compressed data
}

func DecompressText(buf []byte) (string, bool) {
	outSize := 10 * len(buf)
	out := make([]byte, outSize)
	n, err := DecompressSafe(buf, out)
	if err == nil {
		return string(out[:n]), true // uncompressed data
	}
	out = make([]byte, maxOutputSize)
	n, err = DecompressSafe(buf, out)
	if err == nil {
		return string(out[:n]), true // uncompressed data
	}
	fmt.Printf("DeErr: %d %s\n", n, err.Error())
	return "", false
}

func byteSliceToCharPointer(b []byte) *C.char {
	if len(b) == 0 {
		return (*C.char)(unsafe.Pointer(nil))
	}
	return (*C.char)(unsafe.Pointer(&b[0]))
}

// CompressDefault compresses buffer "source" into already allocated "dest" buffer.
// Compression is guaranteed to succeed if size of "dest" >= CompressBound(size of "src")
// The function returns the number of bytes written into buffer "dest".
// If the function cannot compress "source" into a more limited "dest" budget,
// compression stops immediately, and the function result is zero.
func CompressDefault(source, dest []byte) (int, error) {
	ret := int(C.LZ4_compress_default(byteSliceToCharPointer(source),
		byteSliceToCharPointer(dest), C.int(len(source)), C.int(len(dest))))
	if ret == 0 {
		return ret, errors.New("Insufficient destination buffer")
	}

	return ret, nil
}

// DecompressSafe decompresses buffer "source" into already allocated "dest" buffer.
// The function returns the number of bytes written into buffer "dest".
// If destination buffer is not large enough, decoding will stop and output an error code (<0).
// If the source stream is detected malformed, the function will stop decoding and return a negative result.
func DecompressSafe(source, dest []byte) (int, error) {
	ret := int(C.LZ4_decompress_safe(byteSliceToCharPointer(source),
		byteSliceToCharPointer(dest), C.int(len(source)), C.int(len(dest))))
	if ret < 0 {
		return ret, errors.New("Malformed LZ4 source or insufficient destination buffer")
	}

	return ret, nil
}
