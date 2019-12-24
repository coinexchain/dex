module github.com/coinexchain/dex

go 1.13

require (
	github.com/coinexchain/cet-sdk v0.0.0-20191224084723-c730f48d882b
	github.com/coinexchain/codon v0.0.0-20191012070227-3ee72dde596c
	github.com/coinexchain/randsrc v0.0.0-20191012073615-acfab7318ec6
	github.com/cosmos/cosmos-sdk v0.37.4
	github.com/gorilla/mux v1.7.3
	github.com/olekukonko/tablewriter v0.0.1
	github.com/rakyll/statik v0.1.6
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/tendermint v0.32.7
	github.com/tendermint/tm-db v0.2.0
)

replace github.com/cosmos/cosmos-sdk => github.com/coinexchain/cosmos-sdk v0.0.0-20191210021926-99ec1332fbaa

replace github.com/tendermint/tendermint => github.com/coinexchain/tendermint v0.0.0-20191108024645-d56dafa4d3cd
