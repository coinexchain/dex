#!/bin/bash

set -e

if [ ! -f "$GOPATH/bin/cetd" ]; then
    echo "Make sure cetd compiled by (make tools install) and currently can be find in PATH"
    exit 1
fi

rm -rdf ~/.cetd ~/.cetcli
cetd init moniker0 --chain-id=coinexdex-test1
cetcli keys add bob <<<$'12345678\n12345678\n'
cetd add-genesis-account $(cetcli keys show bob -a) 10000000000000000cet
cetd add-genesis-token --name="CoinEx Chain Native Token" \
	--symbol="cet" \
	--owner=$(cetcli keys show bob -a)  \
	--total-supply=10000000000000000 \
	--mintable=false \
	--burnable=true \
	--addr-forbiddable=false \
	--token-forbiddable=false \
	--total-burn=411211452994260000 \
	--total-mint=0 \
	--is-forbidden=false \
	--url="www.coinex.org" \
	--description="A public chain built for the decentralized exchange"
cetd gentx --name bob <<<$'12345678\n12345678\n'
cetd collect-gentxs

echo DONE!
