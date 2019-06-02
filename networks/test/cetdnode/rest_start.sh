#!/usr/bin/env bash

LADDR=$(ifconfig | awk '/inet/&&!/127.0.0.1/{print substr($2,length("addr:")+1)}')
CHAINID=$(awk '/chain_id/{print substr($2,2,length($2)-3)}' /cetd/node0/cetd/config/genesis.json)

./cetcli --home=$1 rest-server --laddr=tcp://$LADDR:27000 --node tcp://localhost:26657 --trust-node=false --chain-id=$CHAINID
