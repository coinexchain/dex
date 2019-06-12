module github.com/coinexchain/dex

require (
	github.com/btcsuite/btcutil v0.0.0-20180706230648-ab6388e0c60a
	github.com/cosmos/cosmos-sdk v0.35.0
	github.com/golangci/golangci-lint v1.17.1 // indirect
	github.com/gorilla/mux v1.7.0
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.6
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/tendermint v0.31.5
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
