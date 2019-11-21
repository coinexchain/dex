module github.com/coinexchain/dex

go 1.13

require (
	github.com/DataDog/zstd v1.4.0 // indirect
	github.com/Shopify/sarama v1.23.1
	github.com/coinexchain/codon v0.0.0-20191012070227-3ee72dde596c
	github.com/coinexchain/randsrc v0.1.0
	github.com/coinexchain/shorthanzi v0.1.0
	github.com/cosmos/cosmos-sdk v0.37.4
	github.com/emirpasic/gods v1.12.0
	github.com/gorilla/mux v1.7.3
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.6
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.7
	github.com/tendermint/tm-db v0.2.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894 // indirect
	gopkg.in/jcmturner/goidentity.v3 v3.0.0 // indirect
)

replace github.com/cosmos/cosmos-sdk => ../cosmos-sdk

replace github.com/coinexchain/codon => ../codon

replace github.com/coinexchain/randsrc v0.0.0-20191012073615-acfab7318ec6 => ../randsrc

replace github.com/tendermint/tendermint => github.com/coinexchain/tendermint v0.0.0-20191108024645-d56dafa4d3cd
