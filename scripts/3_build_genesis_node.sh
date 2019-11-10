#!/bin/bash
#(3/3) generate genesis.json

#for testnet: CHAIN_ID=coinexdex-test1000 ./scripts/3_build_genesis_node.sh
#for mainnet: ./scripts/3_build_genesis_node.sh

set -eux;

if [ "${IS_TESTNET:-false}" == "true" ]; then
    echo "---compile for testnet---"
    TOKEN_IDENTITY=CF1FAAA36A78BE02
    INCENTIVE_POOL_ADDR=cettest1gc5t98jap4zyhmhmyq5af5s7pyv57w566ewmx0
else
    echo "---compile for mainnet---"
    CHAIN_ID=coinexdex
    TOKEN_IDENTITY=C28AB11AA9BB64F0
    INCENTIVE_POOL_ADDR=coinex1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97
fi

# common parameter
TOKEN_SYMBOL=cet
GENESIS_NODE_MONIKER=ViaWallet
OUTPUT_DIR=/tmp/build

# prepare output dir
rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

# assure cetd and cetcli exists
which cetd
which cetcli
which jq || sudo apt-get install -y jq

# generate initial genesis.json
cd ${OUTPUT_DIR}

cetd init ${GENESIS_NODE_MONIKER} --chain-id=${CHAIN_ID} --home ${OUTPUT_DIR}/.cetd

# https://etherscan.io/token/0x081f67afa0ccf8c7b17540767bbe95df2ba8d97f
# date: 2019/10/09 total: 5,868,849,037.61317189
# 586884903761317189 - 300000000000000000 = 286884903761317189

cetd add-genesis-account ${INCENTIVE_POOL_ADDR}   31536000000000000cet --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${circulation}          286384903761317189cet --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${coinex_foundation}     88464000000000000cet --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${genesis_node}            500000000000000cet --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${vesting2020}           36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1577836800  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${vesting2021}           36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1609459200  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${vesting2022}           36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1640995200  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${vesting2023}           36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1672531200  --home ${OUTPUT_DIR}/.cetd
cetd add-genesis-account ${vesting2024}           36000000000000000cet --vesting-amount 36000000000000000cet --vesting-end-time 1704067200  --home ${OUTPUT_DIR}/.cetd

CET_TOKEN_DESCRIPTION="Decentralized public chain ecosystem, Born for financial liberalization"

cetd add-genesis-token --name="CoinEx Token"                            \
    --symbol="${TOKEN_SYMBOL}"                                          \
    --owner=${coinex_foundation}                                        \
    --total-supply=586884903761317189                                   \
    --mintable=false                                                    \
    --burnable=true                                                     \
    --addr-forbiddable=false                                            \
    --token-forbiddable=false                                           \
    --total-burn=413115096238682811                                     \
    --total-mint=0                                                      \
    --is-forbidden=false                                                \
    --url="www.coinex.org"                                              \
    --description="${CET_TOKEN_DESCRIPTION}"                            \
    --identity="${TOKEN_IDENTITY}"                                      \
    --home ${OUTPUT_DIR}/.cetd


# generate tx to create initial validator
mkdir ${OUTPUT_DIR}/gentx

cetd gentx                                \
--name genesis_node                       \
--website wallet.viabtc.com               \
--details "ViaWallet Node"                \
--amount=500000000000000cet               \
--commission-rate=0.1                     \
--commission-max-rate=0.5                 \
--commission-max-change-rate=0.1          \
--min-self-delegation=500000000000000     \
--home ${OUTPUT_DIR}/.cetd                \
--identity 9A30CBDA5872CED8               \
--output-document ${OUTPUT_DIR}/gentx/gentx.json


# add non bondable address
GENESIS_JSON=${OUTPUT_DIR}/.cetd/config/genesis.json

jq ".consensus_params.evidence.max_age = \"1000000\" "                              $GENESIS_JSON  > tmp.$$.json && mv tmp.$$.json $GENESIS_JSON

#if [ "${IS_TESTNET:-false}" == "true" ]; then
#    # adjust testnet parameters
#    jq ".app_state.staking.params.unbonding_time               = \"3600000000000\"  "  $GENESIS_JSON  > tmp.$$.json && mv tmp.$$.json $GENESIS_JSON
#    jq ".app_state.stakingx.params.min_self_delegation         = \"1000000000000\"  "  $GENESIS_JSON  > tmp.$$.json && mv tmp.$$.json $GENESIS_JSON
#    jq ".app_state.gov.deposit_params.max_deposit_period       = \"86400000000000\" "  $GENESIS_JSON  > tmp.$$.json && mv tmp.$$.json $GENESIS_JSON
#    jq ".app_state.gov.voting_params.voting_period             = \"86400000000000\" "  $GENESIS_JSON  > tmp.$$.json && mv tmp.$$.json $GENESIS_JSON
#    jq ".app_state.asset.params.issue_rare_token_fee[0].amount = \"1000000000000\"  "  $GENESIS_JSON  > tmp.$$.json && mv tmp.$$.json $GENESIS_JSON
#    jq ".app_state.asset.params.issue_token_fee[0].amount      = \"100000000000\"   "  $GENESIS_JSON  > tmp.$$.json && mv tmp.$$.json $GENESIS_JSON
#fi


# collect gentx
cetd collect-gentxs --gentx-dir ${OUTPUT_DIR}/gentx  --home ${OUTPUT_DIR}/.cetd

#clean up
rm -rdf ${OUTPUT_DIR}/gentx


# generate secret_connection.key signing.key for tmkms
# tmkms version
# cd ${OUTPUT_DIR} && tmkms softsign keygen secret_connection.key
# cd ${OUTPUT_DIR} && tm-signer-harness extract_key -tmhome ${OUTPUT_DIR}/.cetd -output ./signing.key
# 
# # move the priv_validator_key.json and do not upload it to cloud server
# # cause we'll use tmkms to sign the prevote and precommits
# mv ${OUTPUT_DIR}/.cetd/config/priv_validator_key.json ${OUTPUT_DIR}

#make data dir tarball
cd ${OUTPUT_DIR}
tar cvf ./package.tar ./.cetd
cp ./.cetd/config/genesis.json .
cp `which cetcli` ${OUTPUT_DIR}
cp `which cetd` ${OUTPUT_DIR}


SHA256=sha256sum
if [ "${OSTYPE//[0-9.]/}" == "darwin" ]
then
    SHA256='shasum -a256'
fi

$SHA256 * > sha256.sum

ls ${OUTPUT_DIR}
cat ${OUTPUT_DIR}/sha256.sum
echo "prepare testnet release package succeeded: $(pwd)"


