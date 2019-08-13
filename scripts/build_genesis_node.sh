#!/bin/bash
#(2/2) generate genesis.json

set -eux;

CHAIN_ID=coinexdex-test2000
GENESIS_NODE_MONIKER=GenesisNode
TOKEN_IDENTITY=CF1FAAA36A78BE02
TOKEN_NAME="CoinEx Chain Native Token"
TOKEN_SYMBOL=cet

OUTPUT_DIR=/tmp/build

#prepare output dir
rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

#assure cetd and cetcli exists
which cetd
which cetcli

# generate initial genesis.json
cd ${OUTPUT_DIR}

cetd init ${GENESIS_NODE_MONIKER} --chain-id=${CHAIN_ID} --home ${OUTPUT_DIR}/.cetd

INCENTIVE_POOL_ADDR=coinex1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97
cetd add-genesis-account ${INCENTIVE_POOL_ADDR}                    31536000000000000cet --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account $(cetcli keys show circulation -a)       287767527061317189cet --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account $(cetcli keys show coinex_foundation -a)  88464000000000000cet --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account $(cetcli keys show vesting2020 -a)        36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1577836800  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account $(cetcli keys show vesting2021 -a)        36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1609459200  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account $(cetcli keys show vesting2022 -a)        36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1640995200  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account $(cetcli keys show vesting2023 -a)        36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1672531200  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account $(cetcli keys show vesting2024 -a)        36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1704067200  --home ${OUTPUT_DIR}/.cetd

cetd add-genesis-token --name=${TOKEN_NAME}                             \
    --symbol="${TOKEN_SYMBOL}"                                          \
    --owner=$(cetcli keys show coinex_foundation -a)                    \
    --total-supply=587767527061317189                                   \
    --mintable=false                                                    \
    --burnable=true                                                     \
    --addr-forbiddable=false                                            \
    --token-forbiddable=false                                           \
    --total-burn=412232472938682811                                     \
    --total-mint=0                                                      \
    --is-forbidden=false                                                \
    --url="www.coinex.org"                                              \
    --description="A public chain built for the decentralized exchange" \
    --identity="${TOKEN_IDENTITY}"                                      \
    --home ${OUTPUT_DIR}/.cetd


# generate tx to create initial validator
mkdir ${OUTPUT_DIR}/gentx

cetd gentx                                \
--name coinex_foundation                  \
--website www.coinex.org                  \
--details "Initial genesis node."         \
--amount=200000000000000cet               \
--commission-rate=0.2                     \
--commission-max-rate=1                   \
--commission-max-change-rate=0.01         \
--min-self-delegation=100000000000000     \
--home ${OUTPUT_DIR}/.cetd                \
--output-document ${OUTPUT_DIR}/gentx/gentx.json

# collect gentx
cetd collect-gentxs --gentx-dir ${OUTPUT_DIR}/gentx  --home ${OUTPUT_DIR}/.cetd

#clean up
rm -rdf ${OUTPUT_DIR}/gentx

#make data dir tarball
cd ${OUTPUT_DIR}
tar cvf ./package.tar ./.cetd
cp ./.cetd/config/genesis.json .
cp `which cetcli` ${OUTPUT_DIR}
cp `which cetd` ${OUTPUT_DIR}
md5sum * > md5

echo "prepare testnet release package succeeded: $(pwd)"

