#!/usr/bin/env bash

image_check=$(docker images | grep coinexchain/cetdtest)
if [[ ! ${image_check} ]]
then
    echo "Docker image does NOT exist."
    exit 1
fi

set -e

mkdir func_test

echo "$DPW" | docker login -u "$DUN" --password-stdin
docker pull coinexchain/walle
docker run --rm -v $(pwd)/func_test:/test:Z coinexchain/walle /data/script/cp_data.sh

mkdir func_test/run
pushd func_test
bash script/init.sh
popd

echo "Test begin"
echo "$(pwd)"
tree .
cd func_test
bash script/ft.sh
echo "Test end"
