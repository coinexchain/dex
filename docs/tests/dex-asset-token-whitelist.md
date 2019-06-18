# Asset-Forbid-Whitelist

## Add/Remove Whitelist

- 添加Asset模块Token全局禁止白名单。
  - 只有token的owner可以进行add/Remove whitelist
  - 不能对未发行token进行此操作
  - 只有具备可全局禁止能力的token才能进行此操作
  - 不能添加和删除空地址

> add/remove扣除fee，暂时没有在asset 模块实现，待评估ante-Handler统一收取fee
>
> add/removefee未确认，需要和coinex对齐

## Whitelist CLI & API

- CLI命令
  - `$ cetcli tx asset add-whitelist [flags]` 
  - `$ cetcli tx asset remove-whitelist [flags]` 
  - `$ cetcli q asset whitelist abc [flags]` 
- Rest-curl命令
  - `curl -X POST http://localhost:1317/asset/tokens/coin2/forbidden/whitelist --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"11","account_number":"0"},"whitelist":["coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"]}''`
  - `curl -X POST http://localhost:1317/asset/tokens/coin2/unforbidden/whitelist --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"11","account_number":"0"},"whitelist":["coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"]}'`
  - `curl -X GET http://localhost:1317/asset/tokens/coin2/forbidden/whitelist`

## Whitelist CLI Example

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

3. 通过cli添加coin1的whitelist

```bash
$ cetcli tx asset add-whitelist --symbol="coin1" \
        --whitelist=cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt,cosmos167w96tdvmazakdwkw2u57227eduula2cy572lf,cosmos1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz \
    --from $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: C66A5DFB5CCAAB8F9A2BE039DAC9E3DFDDEACF044E9943760E9E71730B3B88A1
```

4. 此时可以查看到coin1白名单

```bash
$ cetcli q asset whitelist coin1 --chain-id=coinexdex
```

本地返回：

```bash
[
  "coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz",
  "coinex16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt",
  "coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"
]
```

5. remove coin1白名单中的地址

```bash
$ cetcli tx asset remove-whitelist --symbol="coin1" \
        --whitelist=cosmos16gdxm24ht2mxtpz9cma6tr6a6d47x63hlq4pxt,cosmos167w96tdvmazakdwkw2u57227eduula2cy572lf \
    --from $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: C66A5DFB5CCAAB8F9A2BE039DAC9E3DFDDEACF044E9943760E9E71730B3B88A1
```

6. 此时查看到coin1的白名单

```bash
$ cetcli q asset whitelist coin1 --chain-id=coinexdex
```

本地返回：

```bash
[
  "coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz"
]
```



## Whitelist Rest Example

1. 查询本地AccountNumber和Sequence

```bash
$ cetcli query account $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回：

```bash
Account:
  Address:       cosmos1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4
  Pubkey:        cosmospub1addwnpepq03r5ud4j4yx3yzqnzz4yyj8r0r9ysf7pqm92s3at3r8s7rt93ay778g4la
  Coins:         9997000000000000cet,2100000000000000coin1,2100000000000000coin2
  AccountNumber: 0
  Sequence:      11
  LockedCoins:
  FrozenCoins:
  MemoRequired:  false
```

2. 启动rest-server.  参考[本地rest-server中访问swagger-ui的方法](https://github.com/coinexchain/dex/blob/df3c59704ed32917af9e9e47cd203efbfbbc4227/docs/tests/dex-rest-api-swagger.md)

```bash
$ cetcli rest-server --chain-id=coinexdex \ --laddr=tcp://localhost:1317 \ --node tcp://localhost:26657 --trust-node=false

```

3. 通过Rest AP添加coin2白名单，填写本地from/amount/sequence/account_number等信息

```bash
$ curl -X POST http://localhost:1317/asset/tokens/coin2/forbidden/whitelist --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"11","account_number":"0"},"whitelist":["coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz","coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"]}' > unsigned.json
```

返回未签名交易存入unsigned.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgAddTokenWhitelist",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
          "whitelist": [
            "coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz",
            "coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"
          ]
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

5. 本地对交易进行签名

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
        "type": "asset/MsgAddTokenWhitelist",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
          "whitelist": [
            "coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz",
            "coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"
          ]
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
        "signature": "oAioMDCfKwRsuwZJMHMPRa1M+zDAku94ET4b3t/Yc+JroUnJ/hnic2poWCpeavzw7+lf7q4Kvf2JxNw3+43WZg=="
      }
    ],
    "memo": ""
  }
}
```

6. 广播交易

```bash
$ cetcli tx broadcast signed.json
```

本地返回交易Hash

```bash
Response:
  TxHash: 0900D4A88B4D4137168B20C756A77673CD00BB2486B13553A1A0A7100CB70FA5
```

7. 查询coin2的白名单已经添加

```bash
$ curl -X GET http://localhost:1317/asset/tokens/coin2/forbidden/whitelist
```

返回信息：

```bash
[
  "coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz",
  "coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"
]%
```

8. 通过Rest API删除coin2的白名单，填写本地from/amount/sequence/account_number等信息

```bash
$ curl -X POST http://localhost:1317/asset/tokens/coin2/unforbidden/whitelist --data-binary '{"base_req":{"from":"coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4","chain_id":"coinexdex","sequence":"11","account_number":"0"},"whitelist":["coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"]}' > unsigned.json
```

返回未签名交易存入unsigned.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgRemoveTokenWhitelist",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
          "whitelist": [
            "coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"
          ]
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

9. 本地对交易进行签名

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
        "type": "asset/MsgRemoveTokenWhitelist",
        "value": {
          "symbol": "coin2",
          "owner_address": "coinex1x75pqkqaju8eauejjn0kq6pkx907qydusl0ua4",
          "whitelist": [
            "coinex167w96tdvmazakdwkw2u57227eduula2cy572lf"
          ]
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
        "signature": "2Vtfd+44BKmPcT1tz1HyqOtOo9G19je8Q76qF7eSx2pFJ4D55HSFO6knbx/026EmAH5D6WFXDBzL3yuCNI6r3A=="
      }
    ],
    "memo": ""
  }
}
```

10. 广播交易

```bash
$ cetcli tx broadcast signed.json
```

本地返回交易Hash

```bash
Response:
  TxHash: 0900D4A88B4D4137168B20C756A77673CD00BB2486B13553A1A0A7100CB70FA5
```

11. 此时查询coin2的白名单已经更新

```bash
$ curl -X GET http://localhost:1317/asset/tokens/coin2/whitelist
```

返回信息：

```bash
[
  "coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz"
]%
```

