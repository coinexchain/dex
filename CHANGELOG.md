<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog


## [v0.2.8-alpha] - 2020.03-23

### Bug Fixes

*   [#16](https://github.com/coinexchain/dex/issues/16) Rebate amount

### API Breaking

*   [#1](https://github.com/coinexchain/cet-sdk/issues/1)  Add market creator in push msg
*   [#15](https://github.com/coinexchain/dex/issues/15) Uniform timestamp to unix.
*   [#17](https://github.com/coinexchain/dex/issues/17) Uniform bancor message format
*   [#18](https://github.com/coinexchain/dex/issues/18) Modify name of the rebate field

## [v0.2.0] \(WIP\)

### State Machine Breaking

* [CIP0003](https://github.com/coinexchain/CIPs/blob/master/cip-0003.md)
* [CIP0004](https://github.com/coinexchain/CIPs/blob/master/cip-0004.md)
* [CIP0005](https://github.com/coinexchain/CIPs/blob/master/cip-0005.md)

### API Breaking Changes

### Client Breaking Changes

Parameter changes: 

| REST Endpoint                | Change       | Detail                                                       |
| ---------------------------- | ------------ | ------------------------------------------------------------ |
| /asset/parameters            | Response     | added new field: **issue_3char_token_fee**<br />added new field: **issue_4char_token_fee**<br />added new field: **issue_5char_token_fee**<br />added new field: **issue_6char_token_fee** |
| /asset/tokens/{symbol}/infos | Request Body | added new field: **name**<br />added new field: **total_supply**<br />added new field: **mintable**<br />added new field: **burnable**<br />added new field: **addr_forbiddable**<br />added new field: **token_forbiddable** |
|                              |              |                                                              |


### Features

### Improvements

### Bug Fixes



## [v0.0.20]

### State Machine Breaking

### API Breaking Changes
*   [\#4](https://github.com/coinexchain/dex/issues/4) Modify the json name of the field 
*   [\#3](https://github.com/coinexchain/dex/issues/3) Modify time unit 
*   [\#6](https://github.com/coinexchain/dex/issues/6) Modify modules emit events.
*   [\#7](https://github.com/coinexchain/dex/issues/7) Modify swagger.

### Client Breaking Changes
* [\#8](https://github.com/coinexchain/dex/issues/8) Parameter changes

| REST Endpoint       | Response Field                  | Change                                 |
| ------------------- | ------------------------------- | -------------------------------------- |
| /asset/parameters   | issue_token_fee                 | format changed from sdk.Coins to int64 |
| /asset/parameters   | issue_rare_token_fee            | format changed from sdk.Coins to int64 |
| /market/parameters  | gte_order_lifetime              | format changed from int to int64       |
| /market/parameters  | max_executed_price_change_ratio | format changed from int to int64       |
| /staking/parameters | min_self_delegation             | format changed from sdk.Int to int64   |

### Features

### Improvements
* (sdk) Bump SDK version to [v0.37.4](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.37.4).
* (tendermint) Bump Tendermint version to [v0.32.7](https://github.com/tendermint/tendermint/releases/tag/v0.32.7).
* [\#4990](https://github.com/cosmos/cosmos-sdk/issues/4990) Add `Events` to the `ABCIMessageLog` to
provide context and grouping of events based on the messages they correspond to. The `Events` field
in `TxResponse` is deprecated and will be removed in the next major release.

*   [\#5](https://github.com/coinexchain/dex/issues/5) The function of modify the price precision is adjusted.

### Bug Fixes

## [v0.0.18]

### State Machine Breaking
### API Breaking Changes
* [\#1](https://github.com/coinexchain/dex/issues/1) Add a new type of transaction to support OTC.
* [\#2](https://github.com/coinexchain/dex/issues/2) Limit max address length to 45, so UI display will be easier.

### Client Breaking Changes
* [\#1](https://github.com/coinexchain/dex/issues/1) Add a new type of transaction to support OTC.

### Features
* [\#1](https://github.com/coinexchain/dex/issues/1) Add a new type of transaction to support OTC.

### Improvements
### Bug Fixes

## [v0.0.1] - example change log entries

### State Machine Breaking
* [\#4979](https://github.com/cosmos/cosmos-sdk/issues/4979) Introduce a new `halt-time` config and
CLI option to the `start` command. When provided, an application will halt during `Commit` when the
block time is >= the `halt-time`.

### API Breaking Changes
* [\#4979](https://github.com/cosmos/cosmos-sdk/issues/4979) Introduce a new `halt-time` config and
CLI option to the `start` command. When provided, an application will halt during `Commit` when the
block time is >= the `halt-time`.

### Client Breaking Changes
* [\#4979](https://github.com/cosmos/cosmos-sdk/issues/4979) Introduce a new `halt-time` config and
CLI option to the `start` command. When provided, an application will halt during `Commit` when the
block time is >= the `halt-time`.

### Features

* (cli) [\#4973](https://github.com/cosmos/cosmos-sdk/pull/4973) Enable application CPU profiling
via the `--cpu-profile` flag.
* [\#4979](https://github.com/cosmos/cosmos-sdk/issues/4979) Introduce a new `halt-time` config and
CLI option to the `start` command. When provided, an application will halt during `Commit` when the
block time is >= the `halt-time`.

### Improvements

* [\#4990](https://github.com/cosmos/cosmos-sdk/issues/4990) Add `Events` to the `ABCIMessageLog` to
provide context and grouping of events based on the messages they correspond to. The `Events` field
in `TxResponse` is deprecated and will be removed in the next major release.

### Bug Fixes

* [\#4979](https://github.com/cosmos/cosmos-sdk/issues/4979) Use `Signal(os.Interrupt)` over
`os.Exit(0)` during configured halting to allow any `defer` calls to be executed.
* [\#5034](https://github.com/cosmos/cosmos-sdk/issues/5034) Binary search in NFT Module wasn't working on larger sets.

