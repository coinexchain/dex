package keepers

var (
	MarketIdentifierPrefix = []byte{0x15}
	DelistKey              = []byte{0x40}
)

// Merge several byte slices into one
func concatCopyPreAllocate(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

