#!/usr/bin/env bash

cd `dirname $0`

RACE=$(ps -f -p $(ps -f -p $PPID | awk '!/PID/{print $3}') | grep ' -race ')

if [[ -z ${RACE} ]]
then
    go build --buildmode=plugin -o ./data/plugin.so test_plugin.go
else
    go build -race --buildmode=plugin -o ./data/plugin.so test_plugin.go
fi