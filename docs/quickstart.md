## Compile cetd and cetcli

```bash
git clone https://github.com/coinexchain/dex.git
cd dex
```

Compile
```bash
./scripts/build.sh
```

## Bootstrap single testing node
```bash
./scripts/setup_single_testing_node.sh
```

> The generated genesis file's location is $HOME/.cetd/config/genesis.json

## Start cetd

```bash
./cetd start
```
