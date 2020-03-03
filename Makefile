#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif


ifeq ($(IS_TESTNET),true)
	ldflags += -X github.com/coinexchain/cet-sdk/types.Bech32MainPrefix=cettest
	build_tags += testnet
endif


ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
  build_tags += cleveldb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags += -X "github.com/cosmos/cosmos-sdk/version.Name=CoinEx Chain" \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=cetd \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=cetcli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X "github.com/cosmos/cosmos-sdk/types.reDnmString=[a-z][a-z0-9]{1,13}([a-z0-9]{1,2}|(\.[a-z]))?"


ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

ifeq ($(DEBUG), true)
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

# The below include contains the tools target.
include contrib/devtools/Makefile

all: install lint check

build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/cetd.exe ./cmd/cetd
	go build -mod=readonly $(BUILD_FLAGS) -o build/cetcli.exe ./cmd/cetcli
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/cetd ./cmd/cetd
	go build -mod=readonly $(BUILD_FLAGS) -o build/cetcli ./cmd/cetcli
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

install: go.sum check-ledger update-cet-lite-docs
	go install  $(BUILD_FLAGS) ./cmd/cetd
	go install  $(BUILD_FLAGS) ./cmd/cetcli
	go install  $(BUILD_FLAGS) ./cmd/cetdev
	go install  $(BUILD_FLAGS) -buildmode=plugin app/plugin/simple_plugin/simple_plugin.go

install-debug: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/cetdebug


########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/cetd -d 2 | dot -Tpng -o dependency-graph.png

update-cet-lite-docs:
	@statik -src=./cmd/cetcli/swagger -dest=./cmd/cetcli/ -f

clean:
	rm -rf snapcraft-local.yaml build/

distclean: clean
	rm -rf vendor/

########################################
### Testing


check: check-unit check-build
check-all: check check-race check-cover

check-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./...

check-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

check-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

check-build: build
	@go test -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_test


lint: ci-lint
ci-lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/cosmos/cosmos-sdk

benchmark:
	@go test -mod=readonly -bench=. ./...


########################################
### Local validator nodes using docker and docker-compose

build-docker-cetdnode:
	$(MAKE) -C networks/local

statik-swagger:
	statik -src=./cmd/cetcli/swagger -dest=./cmd/cetcli -f -m

build-test-docker: clean statik-swagger build-linux
	cp build/cetd networks/test/cetdnode/
	cp build/cetcli networks/test/cetdnode/
	$(MAKE) -C networks/test
	rm networks/test/cetdnode/cetd
	rm networks/test/cetdnode/cetcli

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/cetd/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/cetd:Z coinexchain/cetdnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

# include simulations
# include sims.mk

.PHONY: all build-linux install install-debug \
	go-mod-cache draw-deps clean \
	check check-all check-build check-cover check-ledger check-unit check-race
