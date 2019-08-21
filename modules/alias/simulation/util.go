package simulation

import "math/rand"

const (
	aliasLetterBytes = "0123456789abcdefghijklmnopqrstuvwxyz-_.@"
	aliasMaxLength   = 100
	aliasMinLength   = 2
)

func randomSymbol(r *rand.Rand, randomLength int) string {
	bytes := make([]byte, 0, randomLength)
	for i := 0; i < randomLength; i++ {
		bytes = append(bytes, aliasLetterBytes[r.Intn(40)])
	}
	return string(bytes)
}
