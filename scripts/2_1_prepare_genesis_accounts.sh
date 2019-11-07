#!/bin/bash
#(2/3) generate account keys

set -eux

cetcli keys add circulation
cetcli keys add coinex_foundation
cetcli keys add vesting2020
cetcli keys add vesting2021
cetcli keys add vesting2022
cetcli keys add vesting2023
cetcli keys add vesting2024


echo "export circulation=$(cetcli keys show -a circulation) \
coinex_foundation=$(cetcli keys show -a coinex_foundation) \
vesting2020=$(cetcli keys show -a vesting2020) \
vesting2021=$(cetcli keys show -a vesting2021) \
vesting2022=$(cetcli keys show -a vesting2022) \
vesting2023=$(cetcli keys show -a vesting2023) \
vesting2024=$(cetcli keys show -a vesting2024)"
