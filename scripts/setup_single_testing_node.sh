#!/bin/bash

set -e

if [ ! -f "$GOPATH/bin/cetd" ]; then
    echo "Make sure cetd compiled by (make tools install) and currently can be find in PATH"
    exit 1
fi

if [ -d "${HOME}/.cetd" ]; then
    echo "Please backup and delete ~/.cetd, before run this script. Exiting..."
    exit 1
fi

if [ -d "${HOME}/.cetcli" ]; then
    echo "Please backup and delete ~/.cetcli, before run this script. Exiting..."
    exit 1
fi

cetd init moniker0 --chain-id=coinexdex-test1
cetcli keys add bob <<<$'12345678\n12345678\n'
cetd add-genesis-account $(cetcli keys show bob -a) 100000000000000000cet
cetd add-genesis-token --name="CoinEx Chain Native Token" \
	--symbol="cet" \
	--owner=$(cetcli keys show bob -a)  \
	--total-supply=100000000000000000 \
	--mintable=false \
	--burnable=true \
	--addr-forbiddable=false \
	--token-forbiddable=false \
	--total-burn=411211452994260000 \
	--total-mint=0 \
	--is-forbidden=false \
	--url="www.coinex.org" \
	--description="A public chain built for the decentralized exchange" \
        --identity=""
cetd gentx --amount=100000000000000cet --min-self-delegation=100000000000000 --name bob <<<$'12345678\n12345678\n'
cetd collect-gentxs

echo DONE!
