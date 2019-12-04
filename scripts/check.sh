#!/usr/bin/env bash


set -ex

function cleanup {
    rc=$?
    [ $rc -ne 0 ] && echo "==============script/check.sh failed=================="
    exit $rc
}
trap cleanup EXIT

if [ ! -x "$(type -p glide)" ]; then
    echo "glide not installed ?"
    exit 1
fi

if [ ! -x "$(type -p golangci-lint)" ]; then
    echo "golangci-lint not installed ?"
    exit 1
fi



find . -name "*.go" -not -path "./vendor/*" -not -path "./git/*" | xargs gofmt -w

if [ "$RUN_IN_TRAVIS" != "true" ]; then
    find . -name "*.go" -not -path "./vendor/*" -not -path "./git/*" | xargs goimports -w -local github.com/coinexchain
fi

linter_targets=$(glide novendor)

test -z "$(golangci-lint  run -j 4 --disable-all \
--enable=gofmt \
--enable=golint \
--enable=gosimple \
--enable=ineffassign \
--enable=vet \
--enable=misspell \
--enable=unconvert \
--exclude='should have comment' \
--exclude='and that stutters;' \
 $linter_targets 2>&1 | grep -v 'ALL_CAPS\|OP_' 2>&1 | tee /dev/stderr)"

time go test -covermode=atomic -coverprofile=coverage.out -race -tags rpctest $linter_targets



