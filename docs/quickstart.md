## Compile

```bash
git clone https://github.com/coinexchain/dex.git
cd dex
```

Compile with
```bash
make tools install
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
