#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/coinexchain/dex/app

sim-dex-nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -SimulationEnabled=true -v -timeout 10m

sim-dex-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.cetd/config/genesis.json will be used."
	@go test -mod=readonly github.com/coinexchain/dex/app -run TestFullAppSimulation -SimulationGenesis=${HOME}/.cetd/config/genesis.json \
		-SimulationEnabled=true -SimulationNumBlocks=100 -SimulationBlockSize=200 -SimulationCommit=true -SimulationSeed=99 -SimulationPeriod=5 -v -timeout 24h

sim-dex-fast:
	@echo "Running quick CoinEx Chain simulation. This may take several minutes..."
	@go test -mod=readonly github.com/coinexchain/dex/app -run TestFullAppSimulation -SimulationEnabled=true -SimulationNumBlocks=100 -SimulationBlockSize=200 -SimulationCommit=true -SimulationSeed=99 -SimulationPeriod=5 -v -timeout 24h

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
	-SimulationEnabled=true -SimulationNumBlocks=1000 -SimulationBlockSize=200 \
	-SimulationCommit=true -SimulationSeed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-dex-benchmark:
	@echo "Running CoinEx Chain benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/coinexchain/dex/app -bench ^BenchmarkFullAppSimulation$$  \
		-SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h

sim-dex-profile:
	@echo "Running CoinEx Chain benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/coinexchain/dex/app -bench ^BenchmarkFullAppSimulation$$ \
		-SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-dex-nondeterminism sim-dex-custom-genesis-fast sim-dex-fast sim-dex-import-export \
	sim-dex-simulation-after-import sim-dex-custom-genesis-multi-seed sim-dex-multi-seed \
	sim-benchmark-invariants sim-dex-benchmark sim-dex-profile