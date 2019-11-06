#!/bin/bash
#(2/3) generate account keys

set -eux

cetcli keys add genesis_node

echo "export genesis_node=$(cetcli keys show -a genesis_node)"
