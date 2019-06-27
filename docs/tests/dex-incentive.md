## Dex-incentive测试文档

#### 需求

按照default plans去分发incentive pool中的cet到各个块，升级后要对现有链高度进行修正后去匹配plan表中的奖励规则。

#### 测试步骤

1. 初始化链

```
./cetd init node0 --chain-id=coinexdex
./cetcli keys add bob
./cetd add-genesis-account $(./cetcli keys show bob -a) 10000000000000000cet
./cetd gentx --name bob
./cetd collect-gentxs
./cetd start
```

2. 停止链并导出genesis.json，并按原有高度继续启动节点

```
./cetd export > genesis.json
helldealer:dex bogon$ 
helldealer:dex bogon$ 
helldealer:dex bogon$ 
helldealer:dex bogon$ cp genesis.json ~/.cetd/config/
helldealer:dex bogon$ ./cetd start
I[2019-06-27|18:03:25.163] Starting ABCI with Tendermint                module=main 
E[2019-06-27|18:03:25.275] Couldn't connect to any seeds                module=p2p 
0
height 9
I[2019-06-27|18:03:27.411] Executed block                               module=state height=9 validTxs=0 invalidTxs=0
I[2019-06-27|18:03:27.419] Committed state                              module=state height=9 txs=0 appHash=AA87E87BDB66A11B830DF8C5E99B1124B9A1AABBCC2A60240D854399ACC7C955
0
height 10
I[2019-06-27|18:03:29.536] Executed block                               module=state height=10 validTxs=0 invalidTxs=0
I[2019-06-27|18:03:29.540] Committed state                              module=state height=10 txs=0 appHash=E90CB079CFA6D489A13705C5F6908BD006C730AE3B819637BAE6F10BB3F6EB2D
0
height 11
I[2019-06-27|18:03:31.672] Executed block                               module=state height=11 validTxs=0 invalidTxs=0
I[2019-06-27|18:03:31.678] Committed state                              module=state height=11 txs=0 appHash=3E5EF07A1EEAE1E57E1DB05610A4E463A73A12041CCEEBF9E77C5657ED80464C
0
height 12
I[2019-06-27|18:03:33.795] Executed block                               module=state height=12 validTxs=0 invalidTxs=0
I[2019-06-27|18:03:33.801] Committed state                              module=state height=12 txs=0 appHash=66C6B7345F874857AA2BCA3C4A3A630A7B76DAFD2A833BF0F6B859901EFCF1A0
```

3. 停止链并导出genesis.json，并从高度0继续启动节点

```
./cetd export --for-zero-height > genesis.json
helldealer:dex bogon$ ./cetd unsafe-reset-all
I[2019-06-27|18:05:04.600] Removed all blockchain history               module=main dir=/Users/helldealer/.cetd/data
I[2019-06-27|18:05:04.602] Reset private validator file to genesis state module=main keyFile=/Users/helldealer/.cetd/config/priv_validator_key.json stateFile=/Users/helldealer/.cetd/data/priv_validator_state.json
helldealer:dex bogon$ 
helldealer:dex bogon$ 
helldealer:dex bogon$ cp genesis.json ~/.cetd/config/
helldealer:dex bogon$ 
helldealer:dex bogon$ 
helldealer:dex bogon$ ./cetd start
I[2019-06-27|18:05:27.432] Starting ABCI with Tendermint                module=main 
E[2019-06-27|18:05:27.533] Couldn't connect to any seeds                module=p2p 
12
height 13
I[2019-06-27|18:05:29.671] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
I[2019-06-27|18:05:29.676] Committed state                              module=state height=1 txs=0 appHash=F282FC7664339E187706E762BFD343D04352FD2E6A6334D756215CFEE70D1167
12
height 14
I[2019-06-27|18:05:31.824] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
I[2019-06-27|18:05:31.830] Committed state                              module=state height=2 txs=0 appHash=D4CAA0F810E619FBF57393A7C75E3236A9DEC7F80D1B30CBB1F8359B9BED1BFB
12
height 15
I[2019-06-27|18:05:33.936] Executed block                               module=state height=3 validTxs=0 invalidTxs=0
I[2019-06-27|18:05:33.940] Committed state                              module=state height=3 txs=0 appHash=37037898E2D44787DC22A7FE6116EB7D1B91A0C1EFA08146F1EF6C300E6544F8
12
height 16
I[2019-06-27|18:05:36.068] Executed block                               module=state height=4 validTxs=0 invalidTxs=0
I[2019-06-27|18:05:36.072] Committed state                              module=state height=4 txs=0 appHash=8BC09331559E823FBD286DBB2CCFFA4E7F66C533D4B61F51AD5A5490C7D8BA85
```

可以看到两次升级都能按照修正后的正确高度分发预定奖励，预定奖励如下

```
"plans": [
          {
            "start_height": "0",
            "end_height": "10512000",
            "reward_per_block": "10",
            "total_incentive": "105120000"
          },
          {
            "start_height": "10512000",
            "end_height": "21024000",
            "reward_per_block": "8",
            "total_incentive": "84096000"
          },
          {
            "start_height": "21024000",
            "end_height": "31536000",
            "reward_per_block": "6",
            "total_incentive": "63072000"
          },
          {
            "start_height": "31536000",
            "end_height": "42048000",
            "reward_per_block": "4",
            "total_incentive": "42048000"
          },
          {
            "start_height": "42048000",
            "end_height": "52560000",
            "reward_per_block": "2",
            "total_incentive": "21024000"
          }
        ]
```

