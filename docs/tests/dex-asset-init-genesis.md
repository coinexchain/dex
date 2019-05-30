# genesis export of asset module
> 导出配置， 重点关注asset.tokens字段
> ./cetd export

```
BJ00609 ~/lab/dex (master) $ ./cetd export
{
...
    "asset": {
    "params": {
        "issue_token_fee": [
        {
            "denom": "cet",
            "amount": "1000000000000"
        }
        ],
        "forbid_address_fee": [
        {
            "denom": "cet",
            "amount": "1000000000"
        }
        ],
        "unforbid_address_fee": [
        {
            "denom": "cet",
            "amount": "1000000000"
        }
        ],
        "forbid_token_fee": [
        {
            "denom": "cet",
            "amount": "100000000000"
        }
        ],
        "unforbid_token_fee": [
        {
            "denom": "cet",
            "amount": "100000000000"
        }
        ],
        "token_forbid_whitelist_add_fee": [
        {
            "denom": "cet",
            "amount": "10000000000"
        }
        ],
        "token_forbid_whitelist_remove_fee": [
        {
            "denom": "cet",
            "amount": "10000000000"
        }
        ],
        "burn_fee": [
        {
            "denom": "cet",
            "amount": "100000000000"
        }
        ],
        "mint_fee": [
        {
            "denom": "cet",
            "amount": "100000000000"
        }
        ]
    },
    "tokens": [
        {
        "type": "asset/Token",
        "value": {
            "name": "CoinEx Chain Native Token",
            "symbol": "cet",
            "total_supply": "588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
        }
        }
    ]
    },
    "gentxs": null
}
}
```

# init from genesis
```
##删除旧的配置及数据， 以便使程序可以从第0块启动时调用initChain, 进而调用到initGenesis逻辑
rm -rdf ~/.cetd ~/.cetcli   
~/lab/dex/cetd init moniker0 --chain-id=coindex
~/lab/dex/cetcli keys add validator0 <<<$'12345678\n12345678\n'
~/lab/dex/cetd add-genesis-account $(~/lab/dex/cetcli keys show validator0 -a) 10000000000000000cet
~/lab/dex/cetd gentx --name validator0 <<<$'12345678\n12345678\n'
~/lab/dex/cetd collect-gentxs
```

在默认生成的配置中， 加入Tokens信息，测试asset模块能够正常加载tokens信息
```
      "tokens": [
        {
          "type": "asset/Token",
          "value": {
            "name": "CoinEx Chain Native Token",
            "symbol": "cet",
            "total_supply": "588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
          }
        }
      ]
```

最终形成的配置如下：
```
{
  "genesis_time": "2019-05-22T14:48:42.022569Z",
  "chain_id": "coindex",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age": "100000"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    }
  },
  "app_hash": "",
  "app_state": {
    "accounts": [
      {
        "address": "cosmos1mqesskxjc6xq2r8tg7sdh8mlkq5tl3nj0tuyfj",
        "coins": [
          {
            "denom": "cet",
            "amount": "10000000000000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": null,
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "0"
      }
    ],
    "auth": {
      "collected_fees": null,
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
        "activated_fee": "1"
      }
    },
    "staking": {
      "pool": {
        "not_bonded_tokens": "10000000000000000",
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
        "community_pool": null
      },
      "community_tax": "0.020000000000000000",
      "base_proposer_reward": "0.010000000000000000",
      "bonus_proposer_reward": "0.040000000000000000",
      "withdraw_addr_enabled": true,
      "delegator_withdraw_infos": null,
      "previous_proposer": "",
      "outstanding_rewards": null,
      "validator_accumulated_commissions": null,
      "validator_historical_rewards": null,
      "validator_current_rewards": null,
      "delegator_starting_infos": null,
      "validator_slash_events": null
    },
    "gov": {
      "starting_proposal_id": "1",
      "deposits": null,
      "votes": null,
      "proposals": null,
      "deposit_params": {
        "min_deposit": [
          {
            "denom": "cet",
            "amount": "10000000"
          }
        ],
        "max_deposit_period": "172800000000000"
      },
      "voting_params": {
        "voting_period": "172800000000000"
      },
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto": "0.334000000000000000"
      }
    },
    "crisis": {
      "constant_fee": {
        "denom": "cet",
        "amount": "1000"
      }
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
        "issue_token_fee": [
          {
            "denom": "cet",
            "amount": "1000000000000"
          }
        ],
        "forbid_address_fee": [
          {
            "denom": "cet",
            "amount": "1000000000"
          }
        ],
        "unforbid_address_fee": [
          {
            "denom": "cet",
            "amount": "1000000000"
          }
        ],
        "forbid_token_fee": [
          {
            "denom": "cet",
            "amount": "100000000000"
          }
        ],
        "unforbid_token_fee": [
          {
            "denom": "cet",
            "amount": "100000000000"
          }
        ],
        "token_forbid_whitelist_add_fee": [
          {
            "denom": "cet",
            "amount": "10000000000"
          }
        ],
        "token_forbid_whitelist_remove_fee": [
          {
            "denom": "cet",
            "amount": "10000000000"
          }
        ],
        "burn_fee": [
          {
            "denom": "cet",
            "amount": "100000000000"
          }
        ],
        "mint_fee": [
          {
            "denom": "cet",
            "amount": "100000000000"
          }
        ]
      },
      "tokens": [
        {
          "type": "asset/Token",
          "value": {
            "name": "CoinEx Chain Native Token",
            "symbol": "cet",
            "total_supply": "588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
          }
        }
      ]
    },
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
                "delegator_address": "cosmos1mqesskxjc6xq2r8tg7sdh8mlkq5tl3nj0tuyfj",
                "validator_address": "cosmosvaloper1mqesskxjc6xq2r8tg7sdh8mlkq5tl3nj2lg39p",
                "pubkey": "cosmosvalconspub1zcjduepqzlumatluxcvaveyn88l4767hp84ltzuxpp8zgdxuyjeev4942f9sj4jj4x",
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
                "value": "AqeQFpc/t28oemObU7KJ1jZe3dVsdcPmV0XeMJ9Fs7up"
              },
              "signature": "0ruJ48Hlf3RdLCQZi+gtJXeNCW/o9YV3RuudN1NtS3Fy484rJxvpfX6xyYgFLxxWUhglsh57mcToYxMsqG2+7A=="
            }
          ],
          "memo": "fd4c585655c92ae5aedee3be89de859d9c2ba5e6@192.168.43.135:26656"
        }
      }
    ]
  }
}

```

## start cetd
> 正常启动， 说明initGenesiis流程正常初始化
> ./cetd start

```
BJ00609 ~/lab/dex (master) $ ./cetd start
I[2019-05-22|22:51:00.197] Starting ABCI with Tendermint                module=main
E[2019-05-22|22:51:00.373] Couldn't connect to any seeds                module=p2p
I[2019-05-22|22:51:05.438] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
I[2019-05-22|22:51:05.450] Committed state                              module=state height=1 txs=0 appHash=E58FF2869529F74CB68EAABE26AF90B29346593E28E449249433F01EE08B0C16
I[2019-05-22|22:51:10.495] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
I[2019-05-22|22:51:10.503] Committed state                              module=state height=2 txs=0 appHash=F20399788A201E3384AECA8B375B5B8752BBB3678AF0488E5B8AB92C4E108E00
```

## 查询Token是否已经恢复
> BJ00609 ~/lab/dex (master) $ ./cetcli q asset token-info cet --chain-id=coindex
```
Token Info: [
  Name:            CoinEx Chain Native Token
  Symbol:          cet
  TotalSupply:     588788547005740000
  Owner:           cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v
  Mintable:        false
  Burnable:        true
  AddrForbiddable:  false
  TokenForbiddable: false
  TotalBurn:       411211452994260000
  TotalMint:       0
  IsFrozen:        false ]
```


## initGenesis 重复token symbol检测
> 当genesis.json中包含重复token时， 启动报错
```
      "tokens": [
        {
          "type": "asset/Token",
          "value": {
            "name": "CoinEx Chain Native Token",
            "symbol": "cet",
            "total_supply": "588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
          }
        },
        {
          "type": "asset/Token",
          "value": {
            "name": "CoinEx Chain Native Token2",
            "symbol": "cet",
            "total_supply": "588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
          }
        }
      ]
```

错误信息如下：
```
BJ00609 ~/lab/dex (master) $ ./cetd start
I[2019-05-22|23:46:09.594] Starting ABCI with Tendermint                module=main
panic: Duplicate token symbol found during asset ValidateGenesis
```

## 去掉重复token后可正常这启动
```
      "tokens": [
        {
          "type": "asset/Token",
          "value": {
            "name": "CoinEx Chain Native Token",
            "symbol": "cet",
            "total_supply": "588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
          }
        },
        {
          "type": "asset/Token",
          "value": {
            "name": "CoinEx Chain Native Token2",
            "symbol": "abc",
            "total_supply": "588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
          }
        }
      ]
```

```
BJ00609 ~/lab/dex (master) $ ./cetd start
I[2019-05-23|00:02:05.907] Starting ABCI with Tendermint                module=main
E[2019-05-23|00:02:06.072] Couldn't connect to any seeds                module=p2p
I[2019-05-23|00:02:11.141] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
I[2019-05-23|00:02:11.151] Committed state                              module=state height=1 txs=0 appHash=14ABC33A443AABE13B75302886D85CE3BAAB5121540FE561B79AC53E46120900
```


## Token信息total supply为负数时， 报错启动失败
```
      "tokens": [
        {
          "type": "asset/Token",
          "value": {
            "name": "CoinEx Chain Native Token",
            "symbol": "cet",
            "total_supply": "-588788547005740000",
            "owner": "cosmos1479jkxzl0gdz6jg7x4843z3eqsvlc5me23wn4v",
            "mintable": false,
            "burnable": true,
            "addr_forbiddable": false,
            "token_forbiddable": false,
            "total_burn": "411211452994260000",
            "total_mint": "0",
            "is_frozen": false
          }
        },
```
启动失败报错：
```
panic: ERROR:
Codespace: asset
Code: 2
Message: "token total supply must a positive"
```