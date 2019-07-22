#!/usr/bin/env bash

set -e

mkdir func_test

echo "$DPW" | docker login -u "$DUN" --password-stdin
docker pull coinexchain/walle
docker run --rm -v $(pwd)/func_test:/test:Z coinexchain/walle /data/script/cp_data.sh

mkdir func_test/run
cd func_test
bash script/init.sh

echo "Test begin"

pipenv run behave ./features/ --tags=~@wip

echo "Test end"
