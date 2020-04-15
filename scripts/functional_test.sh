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
echo $PWD
mkdir func_test

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USER" --password-stdin
echo "begin pull walle"
date +%s
docker pull coinexchain/walle
echo "end pull walle"
date +%s
docker run --name walle coinexchain/walle /data/script/cp_data.sh
docker cp walle:/test func_test

mkdir func_test/run
pushd func_test
mv test/* . && rm -R test
echo $PWD
ls -R
ls script
bash script/init.sh
popd

echo "Test begin"
date +%s

echo "$(pwd)"
cd func_test
# bash script/run_ft_all.sh 
bash script/ft.sh ./features/ --tags=~@wip --no-capture --no-capture-stderr --no-logcapture --no-skipped -D TEST_KAFKA=true

date +%s
echo "Test end"
