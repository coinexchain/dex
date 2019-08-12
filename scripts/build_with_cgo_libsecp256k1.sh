#!/bin/bash

set -ex


if [ ! -f "./go.mod" ]; then
    echo "make sure ./scripts/build.sh executed under dex code root dir."
    exit 1
fi

SRC_DIR=`pwd`


# delete old binaries
CETD_PATH=`which cetd` || true
rm $CETD_PATH || true

CETCLI_PATH=`which cetcli` || true
rm $CETCLI_PATH || true


# install make tools
#sudo apt-get install -y autoconf
#sudo apt-get install -y libtool
#sudo apt-get install -y libgmp-dev


# find correct library version by go.mod
TENDERMINT_VERSION=`grep tendermint/tendermint go.mod | sed -r -e 's/(.*) (v[^ ]*)/\2/g'`
SECP256K1_PATH="$GOPATH/pkg/mod/github.com/tendermint/tendermint@$TENDERMINT_VERSION/crypto/secp256k1/internal/secp256k1/libsecp256k1"
cd $SECP256K1_PATH


# build cgo libsecp256k1
./autogen.sh
./configure --with-bignum=gmp --enable-endomorphism

make -j9
sudo make install


# build cetd and cetcli with cgo libsecp256k1
cd $SRC_DIR
make tools install BUILD_TAGS=libsecp256k1


# check build result contains cgo libsecp256k1
nm `which cetd` | grep func_secp256k1_context_create_sign_verify
RET=$?
if [ $RET -ne 0 ]; then
    echo "FAILED: compiled binary do not contains cgo libsecp256k1!!!"
else
    echo "=====build cetd with cgo libsecp256k1 succeeded====="
fi

exit $RET
