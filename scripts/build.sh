#!/bin/bash

if [ ! -f "./go.mod" ]; then
    echo "make sure ./scripts/build.sh executed under dex code root dir."
    exit 1
fi

statik -src=./cmd/cetcli/swagger -dest=./cmd/cetcli -f -m

(go build github.com/coinexchain/dex/cmd/cetd && go build github.com/coinexchain/dex/cmd/cetcli) && echo "---------- build OK" || echo "---------- build Failed"
