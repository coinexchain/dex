## Compile cetd and cetcli

```bash
git clone https://github.com/coinexchain/dex.git
cd dex

```

Compile
```bash
go build github.com/coinexchain/dex/cmd/cetd
go build github.com/coinexchain/dex/cmd/cetcli
```
or 
```bash
./scripts/build.sh
```
Generating configuration files

> $HOME/.cetd/config/genesis.json


## Start cetd

```bash
./cetd start
```
