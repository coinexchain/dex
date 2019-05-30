# 单节点测试

Reference: https://github.com/cosmos/gaia/blob/master/docs/deploy-testnet.md#single-node-local-manual-testnet



1、编译cetd和cetcli

```
git clone https://github.com/coinexchain/dex.git
cd dex
go build github.com/coinexchain/dex/cmd/cetd
go build github.com/coinexchain/dex/cmd/cetcli
```

执行完这两个命令后，当前目录下会出现可执行文件`cetd`和`cetcli`。此外，也可以用`go run`命令直接运行，比如：

```bash
go run github.com/coinexchain/dex/cmd/cetd
go run github.com/coinexchain/dex/cmd/cetcli
```

1.1 GoLand直接打开工程就可以
  - 需要启用`vgo：GoLand > Preferences > Go > Go Modules (vgo) -> Enable Go Modules (vgo) integration）`

2、生成genesis.json

```bash
# cetd init [moniker] [flags]
./cetd init moniker0 --chain-id=coinexdex # --overwrite
```

运行完上面这条命令后，会生成配置文件`$HOME/.cetd/config/genesis.json`，内容看起来是下面这样：

```json
{
  "genesis_time": "2019-05-15T10:17:56.803706Z",
  "chain_id": "coinexdex",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": { "max_age": "100000" },
    "validator": {
      "pub_key_types": ["ed25519"]
    }
  },
  "app_hash": "",
  "app_state": {
    "accounts": null,
    "auth": {
      "collected_fees": [],
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      }
    },
    "bank": {
      "send_enabled": true
    },
    "bankx": {
      "param": {
        "ActivatedFee": "1"
      }
    },
    "staking": {
      "pool": {
        "not_bonded_tokens": "0",
        "bonded_tokens": "0"
      },
      "params": {
        "unbonding_time": "259200000000000",
        "max_validators": 100,
        "max_entries": 7,
        "bond_denom": "cet"
      },
      "last_total_power": "0",
      "last_validator_powers": null,
      "validators": null,
      "delegations": null,
      "unbonding_delegations": null,
      "redelegations": null,
      "exported": false
    },
    "distr": {
      "fee_pool": {
        "community_pool": []
      },
      "community_tax": "0.020000000000000000",
      "base_proposer_reward": "0.010000000000000000",
      "bonus_proposer_reward": "0.040000000000000000",
      "withdraw_addr_enabled": true,
      "delegator_withdraw_infos": [],
      "previous_proposer": "",
      "outstanding_rewards": [],
      "validator_accumulated_commissions": [],
      "validator_historical_rewards": [],
      "validator_current_rewards": [],
      "delegator_starting_infos": [],
      "validator_slash_events": []
    },
    "gov": {
      "starting_proposal_id": "1",
      "deposits": null,
      "votes": null,
      "proposals": null,
      "deposit_params": {
        "min_deposit": [{"denom": "cet", "amount": "10000000"}],
        "max_deposit_period": "172800000000000"
      },
      "voting_params": { "voting_period": "172800000000000" },
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto": "0.334000000000000000"
      }
    },
    "crisis": {
      "constant_fee": {"denom": "cet", "amount": "1000"}
    },
    "slashing": {
      "params": {
        "max_evidence_age": "120000000000",
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "600000000000",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "signing_infos": {},
      "missed_blocks": {}
    },
    "asset": {
      "params": {
        "issue_token_fee": [{"denom": "cet", "amount": "10000"}],
        "freeze_address_fee": [{"denom": "cet", "amount": "10"}],
        "unfreeze_address_fee": [{"denom": "cet", "amount": "10"}],
        "freeze_token_fee": [{"denom": "cet", "amount": "1000"}],
        "unfreeze_token_fee": [{"denom": "cet", "amount": "1000"}],
        "token_freeze_whitelist_add_fee": [{"denom": "cet", "amount": "100"}],
        "token_freeze_whitelist_remove_fee": [{"denom": "cet", "amount": "100"}],
        "burn_fee": [{"denom": "cet", "amount": "1000"}],
        "mint_fee": [{"denom": "cet", "amount": "1000"}]
      }
    },
    "gentxs": null
  }
}
```



3、生成validator地址

```bash
# cetcli keys add <name> [flags]
./cetcli keys add bob
```



4、添加validator地址

```
# cetcli keys show [name [name...]] [flags]
# cetd add-genesis-account [address_or_key_name] [coin][,[coin]] [flags]
./cetd add-genesis-account $(./cetcli keys show bob -a) 10000000000000000cet
```

这个命令运行完毕之后，genesis.json文件更新了，主要是accounts字段发生了变化：

```json
{
  "genesis_time": "2019-05-14T10:34:52.817032Z",
  "chain_id": "coinexdex",
  "consensus_params": { ... },
  "app_hash": "",
  "app_state": {
    "accounts": [
      {
        "address": "cosmos1qfxpc6hful5hz8p03dk0cy5ygjgzc75jusl9fn",
        "coins": [{"denom": "cet", "amount": "10000000000000000"}],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": null,
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "0"
      }
    ],
    "auth": { ... },
    "bank": { ... },
    "bankx":{...}.
    "staking": { ... },
    "distr": { ... },
    "gov": { ... },
    "crisis": { ... },
    "slashing": { ... },
    "asset": { ... },
    "gentxs": null
  }
}
```



5、GenTx

```bash
# cetd gentx [flags]
./cetd gentx --name bob
```

上面的命令执行成功之后，在`$HOME/.cetd/config/gentx/`目录下面生成了一个json文件，格式化后长这样：

```json
//gentx-dc5db59c0a9cec35fcf416b6176eeb8165c9030e.json
{  
   "type":"auth/StdTx",
   "value":{  
      "msg":[  
         {  
            "type":"cosmos-sdk/MsgCreateValidator",
            "value":{  
               "description":{  
                  "moniker":"moniker0",
                  "identity":"",
                  "website":"",
                  "details":""
               },
               "commission":{  
                  "rate":"0.100000000000000000",
                  "max_rate":"0.200000000000000000",
                  "max_change_rate":"0.010000000000000000"
               },
               "min_self_delegation":"1",
               "delegator_address":"cosmos1qfxpc6hful5hz8p03dk0cy5ygjgzc75jusl9fn",
               "validator_address":"cosmosvaloper1qfxpc6hful5hz8p03dk0cy5ygjgzc75jeyts9q",
               "pubkey":"cosmosvalconspub1zcjduepqp24c3zske6hh3g0u99l6yv6qeaw4pz0ct658q6hwy9plvjhqepkql7rduv",
               "value":{  
                  "denom":"cet",
                  "amount":"100000000"
               }
            }
         }
      ],
      "fee":{  
         "amount":null,
         "gas":"200000"
      },
      "signatures":[  
         {  
            "pub_key":{  
               "type":"tendermint/PubKeySecp256k1",
               "value":"AoNxi/CGrv58uQMuB1BxWRkJswiRUSinLG0aYmwZtJa4"
            },
            "signature":"Fy9t/5w3WMIaiMQ2DzXwURkkB4A3IHGzFptxDvIwoTtyHCJv/RTuiL8+LQ1Dhv09BAsdNgB/47udwj70WjnPgQ=="
         }
      ],
      "memo":"dc5db59c0a9cec35fcf416b6176eeb8165c9030e@192.168.16.9:26656"
   }
}
```



6、 collect-gentxs

```bash
# cetd collect-gentxs [flags]
./cetd collect-gentxs
```

命令执行完毕后，genesis.json发生变化（主要是gentxs字段）：

```json
{
  "genesis_time": "2019-05-14T10:34:52.817032Z",
  "chain_id": "coinexdex",
  "consensus_params": { ... },
  "app_hash": "",
  "app_state": {
    "accounts": [ ... ],
    "auth": { ... },
    "bank": { ... },
    "bankx":{...},
    "staking": { ... },
    "distr": { ... },
    "gov": { ... },
    "crisis": { ... },
    "slashing": { ... },
     "asset":{...},
    "gentxs": [
      {
        "type": "auth/StdTx",
        "value": {
          "msg": [
            {
              "type": "cosmos-sdk/MsgCreateValidator",
              "value": {
                "description": {
                  "moniker": "moniker0",
                  "identity": "",
                  "website": "",
                  "details": ""
                },
                "commission": {
                  "rate": "0.100000000000000000",
                  "max_rate": "0.200000000000000000",
                  "max_change_rate": "0.010000000000000000"
                },
                "min_self_delegation": "1",
                "delegator_address": "cosmos1qfxpc6hful5hz8p03dk0cy5ygjgzc75jusl9fn",
                "validator_address": "cosmosvaloper1qfxpc6hful5hz8p03dk0cy5ygjgzc75jeyts9q",
                "pubkey": "cosmosvalconspub1zcjduepqp24c3zske6hh3g0u99l6yv6qeaw4pz0ct658q6hwy9plvjhqepkql7rduv",
                "value": {
                  "denom": "cet",
                  "amount": "100000000"
                }
              }
            }
          ],
          "fee": {
            "amount": null,
            "gas": "200000"
          },
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "AoNxi/CGrv58uQMuB1BxWRkJswiRUSinLG0aYmwZtJa4"
              },
              "signature": "Fy9t/5w3WMIaiMQ2DzXwURkkB4A3IHGzFptxDvIwoTtyHCJv/RTuiL8+LQ1Dhv09BAsdNgB/47udwj70WjnPgQ=="
            }
          ],
          "memo": "dc5db59c0a9cec35fcf416b6176eeb8165c9030e@192.168.16.9:26656"
        }
      }
    ]
  }
}
```



7、cetd start

```bash
./cetd start
```



8、sendtx

```bash
./cetcli keys add alice
./cetcli keys list
./cetcli tx send $(./cetcli keys show -a alice) 200000000cet 0 \
	--from bob --chain-id=coinexdex --gas 9800
```



9、reset

```
./cetd unsafe-reset-all
```

