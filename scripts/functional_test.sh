#!/usr/bin/env bash

abs() {
    [[ $[ $@ ] -lt 0 ]] && echo "$[ ($@) * -1 ]" || echo "$[ $@ ]"
}

image_check=$(docker images | grep coinexchain/cetdtest)
if [[ ! ${image_check} ]]
then
    echo "Docker image does NOT exist."
    exit 1
fi

set -e

mkdir func_test

echo "$DPW" | docker login -u "$DUN" --password-stdin
echo "begin pull walle"
date +%s
docker pull coinexchain/walle
echo "end pull walle"
date +%s
docker run --rm -v $(pwd)/func_test:/test:Z coinexchain/walle /data/script/cp_data.sh

mkdir func_test/run
pushd func_test
bash script/init.sh
popd

echo "Test begin"
date +%s

echo "$(pwd)"
cd func_test
# bash script/run_ft_all.sh 
bash script/ft.sh ./features/ --tags=~@wip --no-capture --no-capture-stderr --no-logcapture --no-skipped

date +%s
echo "Test end"
