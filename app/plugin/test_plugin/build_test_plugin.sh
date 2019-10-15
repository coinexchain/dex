#!/usr/bin/env bash

cd `dirname $0`
go build -race --buildmode=plugin -o ./data/plugin.so test_plugin.go