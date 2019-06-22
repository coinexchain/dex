#!/bin/bash

set -e

if [ ! -f "$GOPATH/bin/cetd" ]; then
    echo "Make sure cetd compiled by scripts/build.sh and currently can be accessed by ./cetd"
    exit 1
fi

rm -rdf ~/.cetd ~/.cetcli
cetd init moniker0 --chain-id=coinexdex
cetcli keys add bob <<<$'12345678\n12345678\n'
cetd add-genesis-account $(cetcli keys show bob -a) 10000000000000000cet
cetd gentx --name bob <<<$'12345678\n12345678\n'
cetd collect-gentxs

echo DONE!
