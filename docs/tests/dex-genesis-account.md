# Genesis Account For Coinex



针对Coinex团队的[通告]([https://announcement.coinex.com/hc/zh-cn/articles/360027205411-%E5%85%B3%E4%BA%8ECET%E6%B5%81%E9%80%9A%E5%8F%8A%E6%9C%AA%E6%9D%A5%E5%BA%94%E7%94%A8%E7%9A%84%E4%BF%A1%E6%81%AF%E6%8A%AB%E9%9C%B2?from=groupmessage&isappinstalled=0](https://announcement.coinex.com/hc/zh-cn/articles/360027205411-关于CET流通及未来应用的信息披露?from=groupmessage&isappinstalled=0))，在初始时需要为其创建七个账户，已将相应信息添加到[example_genesis.json](https://github.com/coinexchain/dex/blob/master/docs/example_genesis.json)中。



```json
"accounts": [
      {
        "address": "cosmos1c79cqwzah604v0pqg0h88g99p5zg08hgf0cspy",
        "coins": [
          {
            "denom": "cet",
            "amount": "288788547005740000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": null,
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "0",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      },
      {
        "address": "cosmos1n3n5w8mqjf339xse0rwvl0u7nqgp8e5d0nwt20",
        "coins": [
          {
            "denom": "cet",
            "amount": "120000000000000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": null,
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "0",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      },
      {
        "address": "cosmos1xtpex9x7yq8n9d7f8dpgu5mfajrv2thvr6u34q",
        "coins": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "1577836800",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      },
      {
        "address": "cosmos1966f22al7r23h3melq8yt8tnglhweunrxkcezl",
        "coins": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "1609459200",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      },
      {
        "address": "cosmos12kt3yq0kdvu3zm0pq65dkd83hy3j9wgd2m9hfv",
        "coins": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "1640995200",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      },
      {
        "address": "cosmos1r0z8lf82euwlxx0fuvny3jfl0jj2tmdxwuutxj",
        "coins": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "1672531200",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      },
      {
        "address": "cosmos1wezn7xuu5ha39t089mwfeypx0rxvxsutnr0h9p",
        "coins": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "0",
        "original_vesting": [
          {
            "denom": "cet",
            "amount": "36000000000000000"
          }
        ],
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "1704067200",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      }
  ]
```

可以看到，目前genesis.json中，account和accountx两个结构已经合并起来，在account中增加了activated、memo——required、locked_coins三个字段，对于添加的genesis_accout，其默认值分别为true、false、nil。

> 为简单起见，目前暂不支持通过add-genesis-account命令指定activated, memoRequired, lockedCoins三个字段，其中activated的默认值为true，memoRequired的默认值为false，lockedCoins默认为空。后两个字段可以后续通过发送相应交易进行设置。
>
> 相应的导出功能已经实现。

## 测试

1. 参考[单节点启动](signle_node_test.md)，在启动cetd之前，将genesis.json中的accounts下添加以上七个账户，并将genesis.json中staking模块的pool参数中"not_bonded_tokens"的值修改为598800000000000000，随后启动cetd。

2. 通过cetcli查询genesis账户，以第二个账户为例：

   ```bash
   ./cetcli query account cosmos1n3n5w8mqjf339xse0rwvl0u7nqgp8e5d0nwt20 --trust-node
   ```

   返回结果：

   ```bash
   Account:
     Address:       cosmos1n3n5w8mqjf339xse0rwvl0u7nqgp8e5d0nwt20
     Pubkey:        
     Coins:         120000000000000000cet
     AccountNumber: 2
     Sequence:      0
   
   ```

3. 节点启动一段时间之后停止该进程。随后，测试export命令：

   ```bash
   ./cetd export
   ```

   返回结果为：

   ```bash
   {
     "genesis_time": "2019-05-24T08:19:43.781847Z",
     "chain_id": "coinexdex",
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
     "validators": [
       {
         "address": "",
         "pub_key": {
           "type": "tendermint/PubKeyEd25519",
           "value": "Iqp7VjvWTQMSajwXe8s/+GJxlQCtlYnIIGRl9O5JJGY="
         },
         "power": "100",
         "name": "moniker0"
       }
     ],
     "app_hash": "",
     "app_state": {
       "accounts": [
         {
           "address": "cosmos1r0z8lf82euwlxx0fuvny3jfl0jj2tmdxwuutxj",
           "coins": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "sequence_number": "0",
           "account_number": "6",
           "original_vesting": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "1672531200",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
         },
         {
           "address": "cosmos1966f22al7r23h3melq8yt8tnglhweunrxkcezl",
           "coins": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "sequence_number": "0",
           "account_number": "4",
           "original_vesting": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "1609459200",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
         },
         {
           "address": "cosmos1xtpex9x7yq8n9d7f8dpgu5mfajrv2thvr6u34q",
           "coins": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "sequence_number": "0",
           "account_number": "3",
           "original_vesting": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "1577836800",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
         },
         {
           "address": "cosmos12kt3yq0kdvu3zm0pq65dkd83hy3j9wgd2m9hfv",
           "coins": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "sequence_number": "0",
           "account_number": "5",
           "original_vesting": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "1640995200",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
         },
         {
           "address": "cosmos1wezn7xuu5ha39t089mwfeypx0rxvxsutnr0h9p",
           "coins": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "sequence_number": "0",
           "account_number": "7",
           "original_vesting": [
             {
               "denom": "cet",
               "amount": "36000000000000000"
             }
           ],
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "1704067200",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
         },
         {
           "address": "cosmos1n3n5w8mqjf339xse0rwvl0u7nqgp8e5d0nwt20",
           "coins": [
             {
               "denom": "cet",
               "amount": "120000000000000000"
             }
           ],
           "sequence_number": "0",
           "account_number": "2",
           "original_vesting": null,
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "0",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
         },
         {
           "address": "cosmos1c79cqwzah604v0pqg0h88g99p5zg08hgf0cspy",
           "coins": [
             {
               "denom": "cet",
               "amount": "288800000000000000"
             }
           ],
           "sequence_number": "0",
           "account_number": "1",
           "original_vesting": null,
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "0",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
         },
         {
           "address": "cosmos1lmy0as3slwrkh6g56vx2rghcdslmclu5gnphz5",
           "coins": [
             {
               "denom": "cet",
               "amount": "9999999900000000"
             }
           ],
           "sequence_number": "1",
           "account_number": "0",
           "original_vesting": null,
           "delegated_free": null,
           "delegated_vesting": null,
           "start_time": "0",
           "end_time": "0",
           "activated": true,
           "memo_required": false,
           "locked_coins": null
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
           "activated_fee": "100000000"
         }
       },
       "staking": {
         "pool": {
           "not_bonded_tokens": "598799999900000000",
           "bonded_tokens": "100000000"
         },
         "params": {
           "unbonding_time": "259200000000000",
           "max_validators": 100,
           "max_entries": 7,
           "bond_denom": "cet"
         },
         "last_total_power": "100",
         "last_validator_powers": [
           {
             "Address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "Power": "100"
           }
         ],
         "validators": [
           {
             "operator_address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "consensus_pubkey": "cosmosvalconspub1zcjduepqy248k43m6exsxyn28sthhjellp38r9gq4k2cnjpqv3jlfmjfy3nqy88u7j",
             "jailed": false,
             "status": 2,
             "tokens": "100000000",
             "delegator_shares": "100000000.000000000000000000",
             "description": {
               "moniker": "moniker0",
               "identity": "",
               "website": "",
               "details": ""
             },
             "unbonding_height": "0",
             "unbonding_time": "1970-01-01T00:00:00Z",
             "commission": {
               "rate": "0.100000000000000000",
               "max_rate": "0.200000000000000000",
               "max_change_rate": "0.010000000000000000",
               "update_time": "2019-05-24T08:19:43.781847Z"
             },
             "min_self_delegation": "1"
           }
         ],
         "delegations": [
           {
             "delegator_address": "cosmos1lmy0as3slwrkh6g56vx2rghcdslmclu5gnphz5",
             "validator_address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "shares": "100000000.000000000000000000"
           }
         ],
         "unbonding_delegations": null,
         "redelegations": null,
         "exported": true
       },
       "distr": {
         "fee_pool": {
           "community_pool": [
             {
               "denom": "cet",
               "amount": "2.000000000000000000"
             }
           ]
         },
         "community_tax": "0.020000000000000000",
         "base_proposer_reward": "0.010000000000000000",
         "bonus_proposer_reward": "0.040000000000000000",
         "withdraw_addr_enabled": true,
         "delegator_withdraw_infos": [],
         "previous_proposer": "cosmosvalcons1kkc4rlvkns38t398vtksmdk2yma2l0p463r96c",
         "outstanding_rewards": [
           {
             "validator_address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "outstanding_rewards": [
               {
                 "denom": "cet",
                 "amount": "98.000000000000000000"
               }
             ]
           }
         ],
         "validator_accumulated_commissions": [
           {
             "validator_address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "accumulated": [
               {
                 "denom": "cet",
                 "amount": "9.800000000000000000"
               }
             ]
           }
         ],
         "validator_historical_rewards": [
           {
             "validator_address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "period": "1",
             "rewards": {
               "cumulative_reward_ratio": null,
               "reference_count": 2
             }
           }
         ],
         "validator_current_rewards": [
           {
             "validator_address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "rewards": {
               "rewards": [
                 {
                   "denom": "cet",
                   "amount": "88.200000000000000000"
                 }
               ],
               "period": "2"
             }
           }
         ],
         "delegator_starting_infos": [
           {
             "delegator_address": "cosmos1lmy0as3slwrkh6g56vx2rghcdslmclu5gnphz5",
             "validator_address": "cosmosvaloper1lmy0as3slwrkh6g56vx2rghcdslmclu5d84zw8",
             "starting_info": {
               "previous_period": "1",
               "stake": "100000000.000000000000000000",
               "height": "0"
             }
           }
         ],
         "validator_slash_events": []
       },
       "gov": {
         "starting_proposal_id": "1",
         "deposits": null,
         "votes": null,
         "proposals": [],
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
         "signing_infos": {
           "cosmosvalcons1kkc4rlvkns38t398vtksmdk2yma2l0p463r96c": {
             "start_height": "0",
             "index_offset": "1",
             "jailed_until": "1970-01-01T00:00:00Z",
             "tombstoned": false,
             "missed_blocks_counter": "0"
           }
         },
         "missed_blocks": {
           "cosmosvalcons1kkc4rlvkns38t398vtksmdk2yma2l0p463r96c": []
         }
       },
       "asset": {
         "params": {
           "issue_token_fee": [
             {
               "denom": "cet",
               "amount": "1000000000000"
             }
           ],
           "freeze_address_fee": [
             {
               "denom": "cet",
               "amount": "1000000000"
             }
           ],
           "unfreeze_address_fee": [
             {
               "denom": "cet",
               "amount": "1000000000"
             }
           ],
           "freeze_token_fee": [
             {
               "denom": "cet",
               "amount": "100000000000"
             }
           ],
           "unfreeze_token_fee": [
             {
               "denom": "cet",
               "amount": "100000000000"
             }
           ],
           "token_freeze_whitelist_add_fee": [
             {
               "denom": "cet",
               "amount": "10000000000"
             }
           ],
           "token_freeze_whitelist_remove_fee": [
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
         "tokens": null
       },
       "gentxs": null
     }
   }
   
   ```

可以看到，genesis.json中添加的七个账户都成功被导出。

并且，尝试通过发送SetMemoRequiredTx、SendLockedCoinsTx来设置memoRequired字段和lockedcoin之后。导出的state也正常。

