# 多节点本地测试

Reference: https://github.com/cosmos/gaia/blob/master/docs/deploy-testnet.md#multi-node-local-automated-testnet



Summary:

```bash
# Work from the DEX repo
cd path/to/coinexchain/dex

# Build the linux binary in ./build
make build-linux

# Build coinexchain/gaiadnode image
make build-docker-cetdnode

# Start a 4 node testnet run
make localnet-start
```



0、前提条件：请安装[docker](https://www.docker.com/)和[docker-compose](https://docs.docker.com/compose/)



1、编译linux版cetd和cetcli

```bash
$ cd path/to/coinexchain/dex
$ make build-linux
```

命令执行完后，当前目录下会出现build目录，里面有cetd和cetcli两个二进制文件。

```bash
$ tree build/
build/
├── cetcli
└── cetd
```



2、生成docker镜像

```bash
$ make build-docker-cetdnode
```

命令执行完之后，可以看到cetdnode镜像：

```bash
$ docker images
REPOSITORY             TAG                 IMAGE ID            CREATED             SIZE
coinexchain/cetdnode   latest              500d18a69e1b        4 hours ago         13MB
<none>                 <none>              b1d2bf6c141e        4 hours ago         13MB
alpine                 3.7                 6d1ef012b567        2 months ago        4.21MB
```



3、启动4节点测试网

```bash
make localnet-start
```

这条命令实际上会在docker里运行`cetd testnet`命令：

```makefile
# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/cetd/config/genesis.json ]; 
	then 
	  docker run --rm -v $(CURDIR)/build:/cetd:Z coinexchain/cetdnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ;
  fi
	docker-compose up -d
```

执行完毕后，会在build/目录下生成4个节点目录：

```bash
$ tree build/
build/
├── cetcli
├── cetd
├── gentxs
│   ├── node0.json
│   ├── node1.json
│   ├── node2.json
│   └── node3.json
├── node0
│   ├── cetcli
│   │   ├── key_seed.json
│   │   └── keys
│   │       └── keys.db
│   │           ├── 000002.ldb
│   │           ├── 000003.log
│   │           ├── CURRENT
│   │           ├── CURRENT.bak
│   │           ├── LOCK
│   │           ├── LOG
│   │           └── MANIFEST-000004
│   └── cetd
│       ├── cetd.log
│       ├── config
│       │   ├── addrbook.json
│       │   ├── cetd.toml
│       │   ├── config.toml
│       │   ├── genesis.json
│       │   ├── node_key.json
│       │   └── priv_validator_key.json
│       ├── data/
│       └── keys
├── node1/
├── node2/
└── node3/
```

然后会调用`docker-compose up`在容器中启动这4个节点：

```
$ docker container ls -a
CONTAINER ID        IMAGE                  COMMAND                  CREATED             STATUS              PORTS                                                NAMES
0973e3c6c662        coinexchain/cetdnode   "/usr/bin/wrapper.sh…"   2 minutes ago       Up 2 minutes        0.0.0.0:26659->26656/tcp, 0.0.0.0:26660->26657/tcp   cetdnode1
73251b173073        coinexchain/cetdnode   "/usr/bin/wrapper.sh…"   2 minutes ago       Up 2 minutes        0.0.0.0:26663->26656/tcp, 0.0.0.0:26664->26657/tcp   cetdnode3
a27208085f3b        coinexchain/cetdnode   "/usr/bin/wrapper.sh…"   2 minutes ago       Up 2 minutes        0.0.0.0:26656-26657->26656-26657/tcp                 cetdnode0
34164ef5f49a        coinexchain/cetdnode   "/usr/bin/wrapper.sh…"   2 minutes ago       Up 2 minutes        0.0.0.0:26661->26656/tcp, 0.0.0.0:26662->26657/tcp   cetdnode2
```

可以用`docker logs`命令查看node日志：

```
$ docker logs -f cetdnode0
I[2019-05-28|14:04:57.156] Starting ABCI with Tendermint                module=main 
E[2019-05-28|14:04:57.359] Can't add peer's address to addrbook         module=p2p err="Cannot add non-routable address 110ac613962671788e2d82180c3f389b891c2f24@192.168.10.5:26656"
E[2019-05-28|14:04:57.359] Can't add peer's address to addrbook         module=p2p err="Cannot add non-routable address 5060a6dc13d1bc017a233733defa7e397ee89616@192.168.10.3:26656"
E[2019-05-28|14:04:57.359] Can't add peer's address to addrbook         module=p2p err="Cannot add non-routable address 9e190783a58ea4eee4b8130816439cdc388ec42a@192.168.10.4:26656"
E[2019-05-28|14:04:57.403] Couldn't connect to any seeds                module=p2p 
I[2019-05-28|14:05:03.035] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
I[2019-05-28|14:05:03.047] Committed state                              module=state height=1 txs=0 appHash=D23E8BDB5D36D8ED1DEBD18BFEFD19A374673E802595A8EE89C4F35D1836600C
...
```



4、设置addr_book_strict=false

测试网需要修改p2p配置（build/nodeN/cetd/config/config.toml），设置`addr_book_strict = false`，不然会像上面那样报"Cannot add non-routable address"错误。执行下面的命令修改配置：

```bash
sed -i -e 's/addr_book_strict = true/addr_book_strict = false/g' build/node*/cetd/config/config.toml
```

然后执行`make localnet-start`重新启动测试网，重新查看日志：

```
$ docker logs -f cetdnode1
I[2019-05-28|14:58:55.120] Starting ABCI with Tendermint                module=main 
E[2019-05-28|14:58:55.485] Dialing failed                               module=pex addr=86c57bd0bbd5fea89859d58863ba4135bcb598a7@192.168.10.5:26656 err="dial tcp 192.168.10.5:26656: connect: connection refused" attempts=0
E[2019-05-28|14:58:55.485] Dialing failed                               module=pex addr=6483fd7d37e543e8f30c5b1a920684cc89db5ebc@192.168.10.2:26656 err="dial tcp 192.168.10.2:26656: connect: connection refused" attempts=0
I[2019-05-28|14:59:01.030] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
I[2019-05-28|14:59:01.041] Committed state                              module=state height=1 txs=0 appHash=BCEF3CA855A58D58FCC8E892F338C1B30A8D5C97E6EB5D9C71A7A2A2CBE679D0
I[2019-05-28|14:59:06.381] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
I[2019-05-28|14:59:06.357] Committed state                              module=state height=2 txs=0 appHash=A1AB10DB6AA9A81EC5F86497E89677F1125BE831CD6C363F494151C51196FE5F
I[2019-05-28|14:59:11.716] Executed block                               module=state height=3 validTxs=0 invalidTxs=0
I[2019-05-28|14:59:11.724] Committed state                              module=state height=3 txs=0 appHash=E4BA207D29E5BD4D7964F0952D63C99C10CBF10049322F757FE5596328D20F2A
```



5、找到测试网ID和各个node地址

从任一node的genesis.json文件里，可以找到测试网的chain_id，比如"chain-0ze4Qg"。通过执行下面的命令也可以找到chain_id：

```bash
$ docker exec cetdnode0 /cetd/cetcli status
```

从任一node的genesis.json文件里，也可以找到各个node的地址。或者执行下面的命令：

```bash
$ docker exec cetdnode0 /cetd/cetcli keys list --home /cetd/node0/cetcli # cosmos1r309x5f09rwuns2sr8lqmczjgtulkht73hyyew
$ docker exec cetdnode0 /cetd/cetcli keys list --home /cetd/node1/cetcli # cosmos1guuvctmm4fv3psyk43n7gdrrm8zw0r4vnn4s4u
$ docker exec cetdnode0 /cetd/cetcli keys list --home /cetd/node2/cetcli # cosmos13fh95mxxdvuacf3km52t2a0xawv2c700hxf49f
$ docker exec cetdnode0 /cetd/cetcli keys list --home /cetd/node3/cetcli # cosmos1ast383g2ke4g2gjevemwzx5xwhz5jug388tlzp
```



5、在node0发起转账

好了，现在可以发起转账了。在node0里执行下面的命令，给node1转10cet（默认密码12345678）：

```bash
$ docker exec -it cetdnode0 /cetd/cetcli tx send cosmos1guuvctmm4fv3psyk43n7gdrrm8zw0r4vnn4s4u 10cet \
	--from node0 --chain-id=chain-0ze4Qg \
	--gas 50000 --fees 10cet \
	--home /cetd/node0/cetcli
	
{"chain_id":"chain-0ze4Qg","account_number":"0","sequence":"2","fee":{"amount":[{"denom":"cet","amount":"10"}],"gas":"50000"},"msgs":[{"type":"cet-chain/MsgSend","value":{"from_address":"cosmos1r309x5f09rwuns2sr8lqmczjgtulkht73hyyew","to_address":"cosmos1guuvctmm4fv3psyk43n7gdrrm8zw0r4vnn4s4u","amount":[{"denom":"cet","amount":"10"}],"unlock_time":"0"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: Y
Password to sign with 'node0':
Response:
  TxHash: 9B5888F7309C1C965145F5D1C75F884304A325A6CBA1A5CB563ED3FF7B93B480
```

在node0里观察转账结果：

```bash
$ docker exec -it cetdnode0 /cetd/cetcli query tx 9B5888F7309C1C965145F5D1C75F884304A325A6CBA1A5CB563ED3FF7B93B480 \
 --chain-id=chain-0ze4Qg --home /cetd/node0/cetcli
Response:
  Height: 3501
  TxHash: 9B5888F7309C1C965145F5D1C75F884304A325A6CBA1A5CB563ED3FF7B93B480
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 35784
  Tags: 
    - action = send
    - sender = cosmos1r309x5f09rwuns2sr8lqmczjgtulkht73hyyew
    - recipient = cosmos1guuvctmm4fv3psyk43n7gdrrm8zw0r4vnn4s4u

  Timestamp: 2019-05-29T08:18:07
```

在node1里观察转账结果：

```bash
$ docker exec -it cetdnode1 /cetd/cetcli query tx 9B5888F7309C1C965145F5D1C75F884304A325A6CBA1A5CB563ED3FF7B93B480 \
 --chain-id=chain-0ze4Qg --home /cetd/node1/cetcli
 Response:
  Height: 3501
  TxHash: 9B5888F7309C1C965145F5D1C75F884304A325A6CBA1A5CB563ED3FF7B93B480
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 35784
  Tags: 
    - action = send
    - sender = cosmos1r309x5f09rwuns2sr8lqmczjgtulkht73hyyew
    - recipient = cosmos1guuvctmm4fv3psyk43n7gdrrm8zw0r4vnn4s4u

  Timestamp: 2019-05-29T08:18:07Z
```

