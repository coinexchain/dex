#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/coinexchain/dex/app

sim-dex-nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true -v -timeout 10m

sim-dex-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.cetd/config/genesis.json will be used."
	@go test -mod=readonly github.com/coinexchain/dex/app -run TestFullAppSimulation -Genesis=${HOME}/.cetd/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-dex-fast:
	@echo "Running quick CoinEx Chain simulation. This may take several minutes..."
	@go test -mod=readonly github.com/coinexchain/dex/app -run TestFullAppSimulation -Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-dex-import-export: runsim
	@echo "Running CoinEx Chain import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim 25 5 TestAppImportExport

sim-dex-simulation-after-import: runsim
	@echo "Running CoinEx Chain simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim 25 5 TestAppSimulationAfterImport

sim-dex-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.cetd/config/genesis.json will be used."
	$(GOPATH)/bin/runsim -g ${HOME}/.cetd/config/genesis.json 400 5 TestFullAppSimulation

sim-dex-multi-seed: runsim
	@echo "Running multi-seed CoinEx Chain simulation. This may take awhile!"
	$(GOPATH)/bin/runsim 400 5 TestFullAppSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly github.com/coinexchain/dex/app -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Commit=true -Seed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-dex-benchmark:
	@echo "Running CoinEx Chain benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/coinexchain/dex/app -bench ^BenchmarkFullAppSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

sim-dex-profile:
	@echo "Running CoinEx Chain benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/coinexchain/dex/app -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-dex-nondeterminism sim-dex-custom-genesis-fast sim-dex-fast sim-dex-import-export \
	sim-dex-simulation-after-import sim-dex-custom-genesis-multi-seed sim-dex-multi-seed \
	sim-benchmark-invariants sim-dex-benchmark sim-dex-profile