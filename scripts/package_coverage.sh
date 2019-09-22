#!/usr/bin/env bash
currdir=`pwd`
cd $1
echo "" > coverage.txt
for d in $(go list ./... | grep -v simulation); do
    echo "processing $d"
    go test -coverprofile=profile.out $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
cd $currdir
go run coverutil/main.go $1/coverage.txt
