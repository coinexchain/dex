#!/bin/sh

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
if ! [[ -f "${BINARY}" ]]; then
	cp /usr/bin/cetd ${BINARY}
	cp /usr/bin/cetcli ${BINARY_CLI}
	echo "Copy binary to work dirctory."
fi

##
## Run binary with all parameters
##
export CETDHOME="/cetd/node${ID}/cetd"

if [[ -d "`dirname ${CETDHOME}/${LOG}`" ]]; then
  "$BINARY" --home "$CETDHOME" "$@" | tee "${CETDHOME}/${LOG}"
else
  "$BINARY" --home "$CETDHOME" "$@"
fi

chmod 777 -R /cetd
