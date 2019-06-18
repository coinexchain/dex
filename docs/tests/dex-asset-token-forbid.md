# Asset-Forbid-Token

## ForbidToken

- 添加Asset模块Token的全局功能，支持添加白名单。
  - 只有token的owner可以进行token forbid
  - 不能对未发行token进行此操作
  - 只有具备可全局禁止能力的token才能进行此操作
  - 已经被forbid的token需要解除禁止后才能再次forbid

> forbid 扣除fee，暂时没有在asset 模块实现，待评估ante-Handler统一收取fee
>
> forbid  fee未确认，需要和coinex对齐

## ForbidToken CLI & API

- CLI命令
  - `$ cetcli tx asset forbid-token [flags]` 
  - `$ cetcli tx asset unforbid-token [flags]` 
- Rest-curl命令
  - `curl -X POST http://localhost:1317/asset/tokens/coin2/forbids --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"6","account_number":"0"}}'`
  - `curl -X POST http://localhost:1317/asset/tokens/coin2/unforbids --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"6","account_number":"0"}}'`

## ForbidToken CLI Example

参考[single_node_test](https://github.com/coinexchain/dex/blob/master/docs/tests/single-node-test.md)搭建节点，也可以从genesis.json中导入状态，节点启动后

1. 本地创建token，可参考[dex-asset-iusse](https://github.com/coinexchain/dex/blob/master/docs/tests/dex-asset-issue.md) ,查询本地所有token信息：

```bash
$ cetcli query asset tokens --chain-id=coinexdex
```

本地返回所有的token：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": " 1' Token",
      "symbol": "coin1",
      "total_supply": "2100000000000000",
      "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": " 2' Token",
      "symbol": "coin2",
      "total_supply": "2100000000000000",
      "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  }
]
```

2. 通过cli对coin1进行forbid

```bash
$ cetcli tx asset forbid-token --symbol="coin1" \
    --from $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: C66A5DFB5CCAAB8F9A2BE039DAC9E3DFDDEACF044E9943760E9E71730B3B88A1
```

3. 此时查看coin1信息，is_forbidden已经更新

```bash
$ cetcli q asset tokens --chain-id=coinexdex
```

本地返回：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": " 1' Token",
      "symbol": "coin1",
      "total_supply": "2100000000000000",
      "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": true
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": " 2' Token",
      "symbol": "coin2",
      "total_supply": "2100000000000000",
      "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  }
]
```

4. 解除coin1的禁止状态

```bash
$ cetcli tx asset unforbid-token --symbol="coin1" \
    --from $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: C66A5DFB5CCAAB8F9A2BE039DAC9E3DFDDEACF044E9943760E9E71730B3B88A1
```

5. 此时查看coin1信息，is_forbidden已经更新

```bash
$ cetcli q asset tokens --chain-id=coinexdex
```

本地返回：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": " 1' Token",
      "symbol": "coin1",
      "total_supply": "2100000000000000",
      "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": " 2' Token",
      "symbol": "coin2",
      "total_supply": "2100000000000000",
      "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  }
]
```



## ForbidToken Rest Example

1. 查询本地AccountNumber和Sequence

```bash
$ cetcli query account $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回：

```bash
Account:
  Address:       coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4
  Pubkey:        coinexpub1addwnpepq03r5ud4j4yx3yzqnzz4yyj8r0r9ysf7pqm92s3at3r8s7rt93ay778g4la
  Coins:         9997000000000000cet,2100000000000000coin1,2100000000000000coin2
  AccountNumber: 0
  Sequence:      5
  LockedCoins:
  FrozenCoins:
  MemoRequired:  false
```

2. 启动rest-server.  参考[本地rest-server中访问swagger-ui的方法](https://github.com/coinexchain/dex/blob/df3c59704ed32917af9e9e47cd203efbfbbc4227/docs/tests/dex-rest-api-swagger.md)

```bash
$ cetcli rest-server --chain-id=coinexdex \ --laddr=tcp://localhost:1317 \ --node tcp://localhost:26657 --trust-node=false
```

3. 通过Rest API禁止coin2，填写本地from/amount/sequence/account_number等信息

```bash
$ curl -X POST http://localhost:1317/asset/tokens/coin2/forbids --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"5","account_number":"0"}}' > unsigned.json
```

返回未签名交易存入unsigned.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgForbidToken",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4"
        }
      }
    ],
    "fee": {
      "amount": null,
      "gas": "200000"
    },
    "signatures": null,
    "memo": ""
  }
}
```

4. 本地对交易进行签名

```bash
$ cetcli tx sign \
  --chain-id=coinexdex \
  --from $(cetcli keys show bob -a)  unsigned.json > signed.json
```

本地签名后将已签名交易存入signed.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgForbidToken",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4"
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
          "value": "A+I6cbWVSGiQQJiFUhJHG8ZSQT4INlVCPVxGeHhrLHpP"
        },
        "signature": "2062VgdJ4bYURWpmTj4uNLmeCGHWaD/HGJFN5QzK/PZ/bVs8bj4mdKGB37qKWLZ3/MHmMnD7U05hCXoK2tM7SQ=="
      }
    ],
    "memo": ""
  }
}
```

5. 广播交易

```bash
$ cetcli tx broadcast signed.json
```

本地返回交易Hash

```bash
Response:
  TxHash: 0900D4A88B4D4137168B20C756A77673CD00BB2486B13553A1A0A7100CB70FA5
```

6. 此时查询coin2的is_forbidden已经更新

```bash
$ curl -X GET http://localhost:1317/asset/tokens/coin2
```

返回信息：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": " 2' Token",
    "symbol": "coin2",
    "total_supply": "2100000000000000",
    "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
    "mintable": false,
    "burnable": true,
    "addr_forbiddable": true,
    "token_forbiddable": true,
    "total_burn": "0",
    "total_mint": "0",
    "is_forbidden": true
  }
}
```

7. 通过Rest API取消禁止coin2，填写本地from/amount/sequence/account_number等信息

```bash
$ curl -X POST http://localhost:1317/asset/tokens/coin2/unforbids --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"6","account_number":"0"}}' > unsigned.json
```

返回未签名交易存入unsigned.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgUnForbidToken",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4"
        }
      }
    ],
    "fee": {
      "amount": null,
      "gas": "200000"
    },
    "signatures": null,
    "memo": ""
  }
}
```

8. 本地对交易进行签名

```bash
$ cetcli tx sign \
  --chain-id=coinexdex \
  --from $(cetcli keys show bob -a)  unsigned.json > signed.json
```

本地签名后将已签名交易存入signed.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgUnForbidToken",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4"
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
          "value": "A+I6cbWVSGiQQJiFUhJHG8ZSQT4INlVCPVxGeHhrLHpP"
        },
        "signature": "/Bgf9qUXoREw35mX07ZALU/yzefnycyOU7OLIGsFzfssz1FCpPYLTrbuVXSoP5BrlYWs7/KdNL6Pme0QE24UVg=="
      }
    ],
    "memo": ""
  }
}
```

9. 广播交易

```bash
$ cetcli tx broadcast signed.json
```

本地返回交易Hash

```bash
Response:
  TxHash: 0900D4A88B4D4137168B20C756A77673CD00BB2486B13553A1A0A7100CB70FA5
```

10. 此时查询coin2的is_forbidden已经更新

```bash
$ curl -X GET http://localhost:1317/asset/tokens/coin2
```

返回信息：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": " 2' Token",
    "symbol": "coin2",
    "total_supply": "2100000000000000",
    "owner": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
    "mintable": false,
    "burnable": true,
    "addr_forbiddable": true,
    "token_forbiddable": true,
    "total_burn": "0",
    "total_mint": "0",
    "is_forbidden": false
  }
}
```

