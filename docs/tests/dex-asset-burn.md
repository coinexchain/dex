# Asset-Burn-Token

## BurnToken

- 添加Asset模块的BurnToken功能，支持token的燃烧。
  - token的owner可以进行token燃烧
  - 非owner不能进行此操作
  - 不能对未发行token进行此操作
  - 只有具备可燃烧能力的token才能进行此操作
  - 燃烧后的token总量不能为负

> burn 扣除fee，暂时没有在asset 模块实现，待评估ante-Handler统一收取fee
>
> burn fee未确认，需要和coinex对齐

## BurnToken CLI & API

- CLI命令
  - `$ cetcli tx asset burn-token [flags]` 
- Rest-curl命令
  - `curl -X POST http://localhost:1317/asset/tokens/coin3/burns --data-binary '{"base_req":{"from":"cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3","chain_id":"coinexdex","sequence":"8","account_number":"0"},"amount":"2000"}'`

## TransferOwnership CLI Example

参考[single_node_test](https://gitlab.com/cetchain/docs/blob/master/dex/tests/single_node_test.md)搭建节点，也可以从genesis.json中导入状态，节点启动后

1. 本地创建token，可参考[dex-asset-iusse](https://gitlab.com/cetchain/docs/blob/master/dex/tests/dex-asset-issue.md) ,查询本地所有token信息：

```bash
$ cetcli query asset tokens --chain-id=coinexdex
```

本地返回所有的token：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": "first token",
      "symbol": "coin1",
      "total_supply": "2100000000000000",
      "owner": "cosmos16cyga47yh3cv6pzemy0fjtkeqjtrjjukgngey6",
      "mintable": true,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": "sec token",
      "symbol": "coin2",
      "total_supply": "2100000000000100",
      "owner": "cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3",
      "mintable": true,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "100",
      "is_forbidden": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": "th token",
      "symbol": "coin3",
      "total_supply": "2100000000002000",
      "owner": "cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3",
      "mintable": true,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "2000",
      "is_forbidden": false
    }
  }
]
```

3. 我们通过cli进行coin2的燃烧

```bash
$ cetcli tx asset burn-token --symbol="coin2" \
        --amount=100 \
    --from $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: C66A5DFB5CCAAB8F9A2BE039DAC9E3DFDDEACF044E9943760E9E71730B3B88A1
```

6. 此时查看coin2信息，totalsupply已经减少100

```bash
$ cetcli q asset token coin2 --chain-id=coinexdex
```

本地返回：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": "sec token",
    "symbol": "coin2",
    "total_supply": "2100000000000000",
    "owner": "cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3",
    "mintable": true,
    "burnable": true,
    "addr_forbiddable": false,
    "token_forbiddable": true,
    "total_burn": "100",
    "total_mint": "100",
    "is_forbidden": false
  }
}
```



## BurnToken Rest Example

1. 查询本地AccountNumber和Sequence

```bash
$ cetcli query account $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回：

```bash
Account:
  Address:       cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3
  Pubkey:        cosmospub1addwnpepq2uns08x3873dhp0q722pf8yunudlhl3j4s6uxhe0zglusr7p64swxxjts5
  Coins:         9996997900000000cet,2100000000000000coin1,2100000000000000coin2,2100000000000000coin3
  AccountNumber: 0
  Sequence:      11
```

2. 首先需要启动rest-server.  参考[本地rest-server中访问swagger-ui的方法](https://gitlab.com/cetchain/docs/blob/master/dex/tests/dex_rest_api_swagger.md)

```bash
$ cetcli rest-server --chain-id=coinexdex \ --laddr=tcp://localhost:1317 \ --node tcp://localhost:26657 --trust-node=false

```

3. 通过Rest API燃烧coin3，填写本地from/amount/sequence/account_number等信息

```bash
curl -X POST http://localhost:1317/asset/tokens/coin3/burns --data-binary '{"base_req":{"from":"cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3","chain_id":"coinexdex","sequence":"11","account_number":"0"},"amount":"2000"}' > unsigned.json
```

返回未签名交易存入unsigned.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgBurnToken",
        "value": {
          "Symbol": "coin3",
          "Amount": "2000",
          "OwnerAddress": "cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3"
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
        "type": "asset/MsgBurnToken",
        "value": {
          "Symbol": "coin3",
          "Amount": "2000",
          "OwnerAddress": "cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3"
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
          "value": "Ark4POaJ/RbcLweUoKTk5Pjf3/GVYa4a+XiR/kB+DqsH"
        },
        "signature": "lgy8m/+LDhPhLQmyZBrcxXdJRpClSQzl7Ntx0WvUJaAJun2AJIi/uvuAmnTU02CEkUpL+ip43vXVoEOEEM+0Qg=="
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

7. 此时查询coin3已经减少

```bash
$ curl -X GET http://localhost:1317/asset/tokens/coin3
```

返回此token信息：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": "th token",
    "symbol": "coin3",
    "total_supply": "2100000000000000",
    "owner": "cosmos1u0nlxpfsngsyefpa4vjgnng8m8qn3el4cy3ut3",
    "mintable": true,
    "burnable": true,
    "addr_forbiddable": false,
    "token_forbiddable": true,
    "total_burn": "2000",
    "total_mint": "2000",
    "is_forbidden": false
  }
}
```



