1. 启动节点

```
$ ./cetd start
```

2. 创建新账号

```
$ ./cetcli keys add bear

NAME:	TYPE:	ADDRESS:					PUBKEY:
bear	local	cosmos126y5g4ugyck6c5agjcv8jhpuwppctmprgn3x6x	cosmospub1addwnpepq0g03sad7jn437ehdlgh6x70ymnsjynmus4vsvcqa8glmzn9txmtqhha328
```

3. 激活新账号

```
$ ./cetcli tx send $(./cetcli keys show bear -a) 200000000cet --unlock-time=0 --from bob --chain-id=coinexdex --gas 200000

{"chain_id":"coinexdex","account_number":"0","sequence":"1","fee":{"amount":null,"gas":"200000"},"msgs":[{"type":"cet-chain/MsgSend","value":{"from_address":"cosmos19pnqnk82skwqeask05smw7rwdwss9gw4ej4e6c","to_address":"cosmos126y5g4ugyck6c5agjcv8jhpuwppctmprgn3x6x","amount":[{"denom":"cet","amount":"200000000"}],"unlock_time":"0"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: AFBD191C7FF8EC3223D5A8C699FEEB50746DE0EE97132607D7DF47CF453D822E
```

4. 解锁时间设置正确，但是交易币种是cet，交易失败

```
$ ./cetcli tx send $(./cetcli keys show bear -a) 200000000cet --unlock-time=1658670138 --from bob --chain-id=coinexdex --gas 200000

{"chain_id":"coinexdex","account_number":"0","sequence":"2","fee":{"amount":null,"gas":"200000"},"msgs":[{"type":"cet-chain/MsgSend","value":{"from_address":"cosmos19pnqnk82skwqeask05smw7rwdwss9gw4ej4e6c","to_address":"cosmos126y5g4ugyck6c5agjcv8jhpuwppctmprgn3x6x","amount":[{"denom":"cet","amount":"200000000"}],"unlock_time":"1658670138"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: F2452BE34B3CFB826B44291ECE78CE5B694EDE0A2CA7EA3FA484DD3A85E5D84F
  Code: 116
  Raw Log: {"codespace":"bankx","code":116,"message":"Cet cannot be locked"}
```

5. 转账eth，解锁时间设置错误，即不为0且小于当前unix时间戳，交易失败

```
$ ./cetcli tx send $(./cetcli keys show bear -a) 200000000eth --unlock-time=12 --from bob --chain-id=coinexdex --gas 200000

{"chain_id":"coinexdex","account_number":"0","sequence":"2","fee":{"amount":null,"gas":"200000"},"msgs":[{"type":"cet-chain/MsgSend","value":{"from_address":"cosmos19pnqnk82skwqeask05smw7rwdwss9gw4ej4e6c","to_address":"cosmos126y5g4ugyck6c5agjcv8jhpuwppctmprgn3x6x","amount":[{"denom":"eth","amount":"200000000"}],"unlock_time":"12"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: 5AAEDA98458E67D876B84E5A38E8B8B94EACCA93F201679F1DEDDA5801B6B1F7
  Code: 115
  Raw Log: {"codespace":"bankx","code":115,"message":"Invalid Unlock Time"}
```

6. 转账eth，解锁时间正确，交易成功

```
$ ./cetcli tx send $(./cetcli keys show bear -a) 200000000eth --unlock-time=1658670138 --from bob --chain-id=coinexdex --gas 200000

{"chain_id":"coinexdex","account_number":"0","sequence":"2","fee":{"amount":null,"gas":"200000"},"msgs":[{"type":"cet-chain/MsgSend","value":{"from_address":"cosmos19pnqnk82skwqeask05smw7rwdwss9gw4ej4e6c","to_address":"cosmos126y5g4ugyck6c5agjcv8jhpuwppctmprgn3x6x","amount":[{"denom":"eth","amount":"200000000"}],"unlock_time":"1658670138"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: 1E06716ED35B707A63B9CE576F3B7AE331D85D8CB1802128BDD7E409B40C4CE0
```

6. 导出状态，关注accounts和accountsx字段

```
$ ./cetd export
```

```
"app_state": {
    "accounts": [
      {
        "address": "cosmos19pnqnk82skwqeask05smw7rwdwss9gw4ej4e6c",
        "coins": [
          {
            "denom": "cet",
            "amount": "9999999700000000"
          },
          {
            "denom": "eth",
            "amount": "9999999800000000"
          }
        ],
        "sequence_number": "3",
        "account_number": "0",
        "original_vesting": null,
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "0"
      },
      {
        "address": "cosmos126y5g4ugyck6c5agjcv8jhpuwppctmprgn3x6x",
        "coins": [
          {
            "denom": "cet",
            "amount": "100000000"
          }
        ],
        "sequence_number": "0",
        "account_number": "1",
        "original_vesting": null,
        "delegated_free": null,
        "delegated_vesting": null,
        "start_time": "0",
        "end_time": "0"
      }
    ],
    "accountsx": [
      {
        "address": "cosmos19pnqnk82skwqeask05smw7rwdwss9gw4ej4e6c",
        "activated": true,
        "memo_required": false,
        "locked_coins": null
      },
      {
        "address": "cosmos126y5g4ugyck6c5agjcv8jhpuwppctmprgn3x6x",
        "activated": true,
        "memo_required": false,
        "locked_coins": [
          {
            "coin": {
              "denom": "eth",
              "amount": "200000000"
            },
            "unlock_time": "1658670138"
          }
        ]
      }
    ],
```

