#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/cetd/${BINARY:-cetd}
ID=${ID:-0}
LOG=${LOG:-cetd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'cetd' E.g.: -e BINARY=cetd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export CETDHOME="/cetd/node${ID}/cetd"

if [ -d "`dirname ${CETDHOME}/${LOG}`" ]; then
  "$BINARY" --home "$CETDHOME" "$@" | tee "${CETDHOME}/${LOG}"
else
  "$BINARY" --home "$CETDHOME" "$@"
fi

chmod 777 -R /cetd

