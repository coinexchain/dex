## Compile

```bash
git clone https://github.com/coinexchain/dex.git
cd dex
```

Using multiple versions of secp256k1 implementations, the time required to verify the signature is as follows.

ECDSA | libsecp256k1 in c | libsecp256k1 in cgo | secp256k1 in go
---------|-----------|----------|---------|
Sign | 46000ns | 92138ns | 81926ns | 
Verify | 69200ns | 151701ns | 236794ns | 

### Compile with go-secp256k1
```bash
make tools install
``` 

### Compile with c-libsecp256k1

If you want the node to run faster, the following command is recommended for cgo compilation。

Compile libsecp256k1
```
cd tendermint@v0.32.1/crypto/secp256k1/internal/secp256k1/libsecp256k1
./autogen.sh
./configure --with-bignum=gmp --enable-endomorphism
make -j2 && make install
```

Compile dex
```
cd dex
make tools install BUILD_TAGS=libsecp256k1
```

- `cetd` and `cetcli` will be install in your GOPATH.
- `cetd` is the CoinEx Chain full node daemon. 
- `cetcli` is the CLI tool to interact with `cetd`.

## Bootstrap single testing node
```bash
./scripts/setup_single_testing_node.sh
```

> The generated genesis file's location is $HOME/.cetd/config/genesis.json

## Start cetd

```bash
./cetd start
```

## Start rest-server

Start with the command：
```bash
./cetcli rest-server --chain-id=coinexdex  --laddr=tcp://localhost:8080  --node tcp://localhost:26657 --trust-node=false
```

Local access :
> http://localhost:1317/swagger/
