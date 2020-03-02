#!/bin/bash

set -eux

if [ ! -f "./cetd" ]; then
  echo "Make sure cetd currently can be find in pwd"
  exit 1
fi

if [ ! -f "./cetcli" ]; then
  echo "Make sure cetcli currently can be find in pwd"
  exit 1
fi

if [ ! -f "./cetd2" ]; then
  echo "Make sure cetd2 currently can be find in pwd"
  exit 1
fi

if [ ! -f "./cetcli2" ]; then
  echo "Make sure cetcli2 currently can be find in pwd"
  exit 1
fi

GENESIS_FILE=genesis.json

./cetd export --for-zero-height=false >${GENESIS_FILE}
./cetd2 migrate ${GENESIS_FILE} --genesis-block-height="${GENESIS_BLOCK_HEIGHT:-0}" --output ${GENESIS_FILE}
./cetd unsafe-reset-all
cp ${GENESIS_FILE} "${CHAIN_DIR:-${HOME}/.cetd}"/config/genesis.json

rm -rf ${GENESIS_FILE}

echo Done!
