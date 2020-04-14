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
docker pull ludetewill/walle
echo "end pull walle"
date +%s
docker run --rm -v $(pwd)/func_test:/test:Z ludetewill/walle /data/script/cp_data.sh

echo $PWD
ls -R
mkdir func_test/run
pushd func_test
echo $PWD
ls -R
cd ..
echo $PWD
ls -R
ls script
bash script/init.sh
popd

echo "Test begin"
date +%s

echo "$(pwd)"
cd func_test
if [[ $1 -lt 0 ]]; then
    bash script/run_ft_non_cli_mode.sh $(abs $1) $2
else
    bash script/run_ft_in_parallel.sh $1 $2
fi

date +%s
echo "Test end"
