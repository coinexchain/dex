#!/bin/bash

LADDR=$(ifconfig | awk '/inet/&&!/127.0.0.1/{print substr($2,length("addr:")+1)}')
CHAIN_ID=`sed -En 's/.*chain_id": ?"([^"]*)".*/\1 /p' /cetd/node0/cetd/config/genesis.json`

export _RR_TRACE_DIR="$1"
./cetcli --home=$1 rest-server --laddr=tcp://${LADDR}:27000 --node tcp://localhost:26657 --trust-node=true --chain-id=${CHAIN_ID} > $1/rest.log
