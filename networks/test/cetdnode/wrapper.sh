#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/cetd/${BINARY:-cetd}
BINARY_CLI=/cetd/${BINARY_CLI:-cetcli}
ID=${ID:-0}
LOG=${LOG:-cetd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	cp /usr/bin/cetd ${BINARY}
	cp /usr/bin/cetcli ${BINARY_CLI}
	echo "Copy binary to work dirctory."
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

