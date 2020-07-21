### "Atlantis" Hard Fork Guide

#### 1. Overview

CoinEx Chain will have a hard fork on March 30, 2020, with the code name "Atlantis". This hard fork will use a genesis.json file for data migration. The data directory of the old chain, which stores historical blocks and latest state, is no longer needed by the new chain. The new chain will use a new data directory and have a new Chain-ID.

The overall flow of this hard fork is:

1. When the old chain reaches a pre-defined height (or when it exceeds this height by several blocks), we stop the cetd process of the old chain.
2. Use the cetd of the old chain to export the latest state at the pre-defined height, which will be stored in a file named genesis.exported.json 
3. Use the cetd of the new chain to process this genesis.exported.json file, and we'll get the genesis.json for the new chain.
4. Use this genesis.json file to start cetd of the new chain.

This article introduces the whole hard fork flow. You can find some linux commands here, which can be used for "copy&paste". But, please note this article may be updated any time before the final hard fork. So, before you copy&paste, please make sure this article is up-to-date, and please reload your browser page if necessary.

#### 2. Install dependencies

```bash
git clone https://github.com/facebook/rocksdb.git && cd rocksdb
git checkout v6.6.4
make static_lib && sudo make install
```


#### 3. Use genesis.json to start a new chain

##### 3.1. Download the binary file of the new chain

First define some environmenal variables:


>  export ARTIFACTS_URL=https://github.com/coinexchain/dex/releases/download/v0.2.17/linux_x86_64.tar.gz <br/>
>  export PUBLIC_IP=~~123.36.28.137~~ <br/>
>  export RUN_DIR=~~/path/to/work-dir~~ <br/>

Please not the above `PUBLIC_IP` and `RUN_DIR` have example values, which may not be suitable for your case, so you must change them according to your needs.

Then you can download the files:
```bash
mkdir ${RUN_DIR} && cd ${RUN_DIR}
wget ${ARTIFACTS_URL} && tar -zxvf linux_x86_64.tar.gz 
```

#### 3.2. Make new data directory

1. `${RUN_DIR}/cetd init moniker --chain-id=coinexdex2 --home=${RUN_DIR}/.cetd`
2. Copy the downloaded genesis.json file to the data directory: `cp genesis.json ${RUN_DIR}/.cetd/config`
3. If you are a validator

    *   Copy the ED25519 private key file `priv_validator_key.json` of the old chain to the data direction of the new chain. It should be at `${RUN_DIR}/.cetd/config`
4. If you aren't a validator, configure the external IP of this node

   *   `sed -i "/external_address/cexternal_address = \"tcp://${PUBLIC_IP}:26656\"" ${RUN_DIR}/.cetd/config/config.toml`
5. Verify the execuatbles and genesis.json file, etc:
   *  `bash dex2_check.sh`


#### 3.3. Use nohup to start a new node

First define a list of seed nodes. Seed nodes are the several nodes who start firstly in the hard fork, which are not known until the hard fork happens.


```bash
export CHAIN_SEEDS=903458cf236851ccf8604689c3f391c528191f47@47.75.37.80:26656,9be765dffed72adcd27ebb37c79bf8ac501f43e8@47.52.155.115:26656,cd79d6c2b3b6b561c91b61b8e3a706249b532ca4@47.56.215.151:26656,cf34ba278ce69be1240f1dabad9b57ffecae206a@47.75.60.29:26656,c70feea1a4f8ea2fd55c366fdcb7ca4d53f1c775@18.144.85.87:26656,94b718f31dedf4afee4c04d768343166625cf961@47.52.70.137:26656,2cbef50b8c996745b9c8a0059fe32a1fbfef8b46@47.52.129.186:26656,17ec2dcfd7c72fabcb7c7cfe2d71006fc39c85c9@18.180.56.174:26656
```

Then use nohup to start cetd:

```bash
nohup ${RUN_DIR}/cetd start --home=${RUN_DIR}/.cetd --minimum-gas-prices=20.0cet --p2p.seeds=${CHAIN_SEEDS} &> cetd.log &
```

#### 3.4. Configure cetd to start automatically using systemctl or supervisor

When cetd has been running correctly for about 1~2 hours, you can kill the cetd process and switch to systemctl or supervisor, which can be configure cetd to start automatically.

First add the seed list into config.toml:

```bash
sed -i "s/seeds = \"\"/seeds = \"${CHAIN_SEEDS}\"/" ${RUN_DIR}/.cetd/config/config.toml
```

Then you can use your favorite tool such as systemctl or supervisor, to make cetd as an automatically-starting deamon.

