module github.com/coinexchain/dex

require (
	github.com/btcsuite/btcutil v0.0.0-20180706230648-ab6388e0c60a
	github.com/cosmos/cosmos-sdk v0.34.7
	github.com/gorilla/mux v1.7.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.6
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/tendermint v0.31.5
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a // indirect
	golang.org/x/net v0.0.0-20190313220215-9f648a60d977 // indirect
	golang.org/x/sys v0.0.0-20190312061237-fead79001313 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
