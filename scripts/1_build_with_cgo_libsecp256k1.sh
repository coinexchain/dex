#!/bin/bash
#(1/3) build cetd and cetcli

set -ex

if [ ! -f "./go.mod" ]; then
    echo "make sure ./scripts/build.sh executed under dex code root dir."
    exit 1
fi

DEX_SRC_DIR=`pwd`


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

TMP_DIR=/tmp/libsecp256k1
rm -rdf $TMP_DIR && mkdir -p $TMP_DIR
cp -r $SECP256K1_PATH/* $TMP_DIR
chmod a+w -R $TMP_DIR
cd $TMP_DIR

# build cgo libsecp256k1
./autogen.sh
./configure --with-bignum=gmp --enable-endomorphism --prefix=$TMP_DIR/output

make -j9
make install


# build cetd and cetcli with cgo libsecp256k1
cd $DEX_SRC_DIR
make tools install BUILD_TAGS=libsecp256k1


# check build result contains cgo libsecp256k1
nm `which cetd` | grep func_secp256k1_context_create_sign_verify
RET=$?
if [ $RET -ne 0 ]; then
    echo "FAILED: compiled binary do not contains cgo libsecp256k1!!!"
else
    echo "=====build cetd with cgo libsecp256k1 succeeded====="
fi

md5sum `which cetd`
md5sum `which cetcli`
md5sum `which cetdev`

cetd version --long
cetcli version --long

exit $RET
