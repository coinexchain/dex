#!/bin/bash

set -e

if [ ! -f "./go.mod" ]; then
    echo "make sure ./scripts/build.sh executed under dex code root dir."
    exit 1
fi

pushd cmd/cetcli/swagger
bower install
popd

statik -src=./cmd/cetcli/swagger -dest=./cmd/cetcli -f -m

(go build -gcflags='all=-N -l' github.com/coinexchain/dex/cmd/cetd  && go build -gcflags='all=-N -l' github.com/coinexchain/dex/cmd/cetcli ) && echo "---------- build OK" || echo "---------- build Failed"
