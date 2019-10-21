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


## [Unreleased]

### State Machine Breaking

### API Breaking Changes

### Client Breaking Changes

### Features

### Improvements

### Bug Fixes

## [v0.0.18]

### State Machine Breaking
### API Breaking Changes
* [\#1](https://github.com/coinexchain/dex/issues/1) Add a new type of transaction to support OTC.
* [\#2](https://github.com/coinexchain/dex/issues/2) Limit max address length to 45, so UI display will be easier.

### Client Breaking Changes
### Features
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

