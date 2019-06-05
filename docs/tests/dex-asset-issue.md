# Asset-IssueToken

## IssueToken

- 参考[cosmos tutorial](https://cosmos.network/docs/intro/#sdk-application-architecture)添加Asset模块，支持token的发行。
  - Token Name 限制在32 unicode characters
  - Token Symbol format `^[a-z][a-z0-9]{1,7}$`
  - Token Symbol 不能重名，创建已有的token
  - Token TotalSupply 在boosting前不能超过90 billion
  - Token issue fee扣除1000000000000cet

> 发行token扣除的fee，暂时在asset 模块Handler中做了实现，待ante-Handler统一收取fee

## IssueToken CLI & API

- CLI命令
  - `$ cetcli tx asset issue-token [flags]` 
  - `$ cetcli query asset token [symbol] [flags]`
  - `$ cetcli query asset tokens  [flags]`

- Rest-curl命令
  - `$ curl -X GET http://localhost:1317/asset/tokens/symbol`
  - `$ curl -X GET http://localhost:1317/asset/tokens`
  - `$ curl -X POST http://localhost:1317/asset/tokens --data-binary '{"base_req":{"from":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","chain_id":"coinexdex","sequence":"4","account_number":"0"},"name":"my first rest coin","symbol":"coin1","total_supply":"10000000000","mintable":"false","burnable":"true","addr_forbiddable":"false","token_forbiddable":"true"}'`

## IssueToken CLI Example

节点的搭建参考[single_node_test](https://github.com/coinexchain/dex/blob/df3c59704ed32917af9e9e47cd203efbfbbc4227/docs/tests/single-node-test.md)，也可以从genesis.json中导入状态，节点启动后

1. 尝试查询当前token-list

```bash
$ cetcli query asset tokens --chain-id=coinexdex
```

当前并没有新创建的token：

```bash
ERROR: {"codespace":"asset","code":205,"message":"can not query any token"}
```

2. 查询当前validator地址

```bash
$ cetcli keys show bob -a
```

本地返回：

```bash
cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd
```

3. 创建token，通过--Flag指定所要创建的token信息

```bash
$ cetcli tx asset issue-token --name="my first token" \
        --symbol="token1" \
        --total-supply=2100000000000000 \
        --mintable=false \
        --burnable=true \
        --addr-forbiddable=0 \
        --token-forbiddable=1 \
        --from $(cetcli keys show bob -a) \
        --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: CF7686265FB4DA2AA74B20CF92DEDF5792D22DA3385C9C91964D849B26C1C5FA
```

4. 如果再次重复创建symbol:"token1"，会创建失败

```bash
$ cetcli tx asset issue-token --name="my duplicat token" \
        --symbol="token1" \
        --total-supply=2100000000000000 \
        --mintable=false \
        --burnable=true \
        --addr-forbiddable=false \
        --token-forbiddable=true \
        --from $(cetcli keys show bob -a) \
        --chain-id=coinexdex
```

本地返回提示：

```bash
ERROR: token symbol already exists，pls query tokens and issue another symbol
```

5. 此时可以传入symbol参数`token1`查询到token-info

```bash
$ cetcli query asset token token1 --chain-id=coinexdex
```

本地返回token的信息：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": "my first token",
    "symbol": "token1",
    "total_supply": "2100000000000000",
    "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
    "mintable": false,
    "burnable": true,
    "addr_forbiddable": false,
    "token_forbiddable": true,
    "total_burn": "0",
    "total_mint": "0",
    "is_forbidden": false
  }
}
```

6. 如上可以创建token2，token3等token，可以查询所有的token

```bash
$ cetcli query asset tokens --chain-id=coinexdex
```

本地返回所有的token：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": "my first token",
      "symbol": "token1",
      "total_supply": "2100000000000000",
      "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
      "mintable": false,
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
      "name": "my sec token",
      "symbol": "token2",
      "total_supply": "2100000000000000",
      "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
      "mintable": false,
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
      "name": "my th token",
      "symbol": "token3",
      "total_supply": "2100000000000000",
      "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  }
]
```

## IssueToken Rest Example

1. 查询本地AccountNumber和Sequence

```bash
$ cetcli query account $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回：

```bash
Account:
  Address:       cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd
  Pubkey:        cosmospub1addwnpepq0mdsnxm75k543ruyg7v9gnd9n2s55lwf2t7sp60kl874rchq8vj5w4t7e4
  Coins:         9996999900000000cet,2100000000000000token1,2100000000000000token2,2100000000000000token3
  AccountNumber: 0
  Sequence:      4
```

2. 首先需要启动rest-server.  可参考[本地rest-server中访问swagger-ui的方法](https://github.com/coinexchain/dex/blob/df3c59704ed32917af9e9e47cd203efbfbbc4227/docs/tests/dex-rest-api-swagger.md)

```bash
$ cetcli rest-server --chain-id=coinexdex --laddr=tcp://localhost:1317 --node tcp://localhost:26657 --trust-node=false
```

3. 可以通过GET查询token信息

```bash
$ curl -X GET http://localhost:1317/asset/tokens/token1
```

返回token1的信息：

```bash
{"type":"asset/Token","value":{"name":"my first token","symbol":"token1","total_supply":"2100000000000000","owner":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","mintable":false,"burnable":true,"addr_forbiddable":false,"token_forbiddable":true,"total_burn":"0","total_mint":"0","is_forbidden":false}}%
```

4. 查询所有token信息

```bash
$ curl -X GET http://localhost:1317/asset/tokens
```

返回所有token信息：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": "my first token",
      "symbol": "token1",
      "total_supply": "2100000000000000",
      "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
      "mintable": false,
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
      "name": "my sec token",
      "symbol": "token2",
      "total_supply": "2100000000000000",
      "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
      "mintable": false,
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
      "name": "my th token",
      "symbol": "token3",
      "total_supply": "2100000000000000",
      "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false
    }
  }
]
```

5. 通过Rest API 发行token

```bash
curl -X POST http://localhost:1317/asset/tokens --data-binary '{"base_req":{"from":"cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd","chain_id":"coinexdex","sequence":"3","account_number":"0"},"name":"my first rest coin","symbol":"coin1","total_supply":"10000000000","mintable":false,"burnable":true,"addr_forbiddable":false,"token_forbiddable":true}' > unsigned.json
```

返回未签名交易存入unsigned.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgIssueToken",
        "value": {
          "name": "my first rest coin",
          "symbol": "coin1",
          "total_supply": "10000000000",
          "owner": "cosmos1n9emr2kwt70aajjreklu2w9d3jamm4nwkpnp2l",
          "mintable": false,
          "burnable": true,
          "addr_forbiddable": false,
          "token_forbiddable": true
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

6. 本地对交易进行签名

```bash
$ cetcli tx sign \
  --chain-id=coinexdex \
  --from=$(cetcli keys show bob -a)  unsigned.json > signed.json
```

本地签名后将已签名交易存入signed.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgIssueToken",
        "value": {
          "name": "my first rest coin",
          "symbol": "coin9",
          "total_supply": "10000000000",
          "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
          "mintable": false,
          "burnable": true,
          "addr_forbiddable": false,
          "token_forbiddable": true
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
          "value": "AwCqzZ5e3NhCF4QAtgrz+KT2kDvGt+rzOse87T7HfCS7"
        },
        "signature": "IAf/v0+Z5KV5kCTy5LWy2Y5ySPyaKhBoSbzVIctWSfoGb75XF0HVbE4DABEVnbntWKy7MA9iV7HYTXD+1NBePA=="
      }
    ],
    "memo": ""
  }
}
```

7. 广播交易

```bash
$ cetcli tx broadcast signed.json
```

本地返回交易Hash

```bash
Response:
  TxHash: 623321E41440540D6B980542C4C38BA22D0E6AE7284006188C212AFD801AA270
```

8. 此时查询此coin1已经被创建

```bash
$ curl -X GET http://localhost:1317/asset/tokens/coin1
```

返回此token信息：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": "my first rest coin",
    "symbol": "coin1",
    "total_supply": "10000000000",
    "owner": "cosmos1n9e8krs6dengw6k8ts0xpntyzd27rhj48ve5gd",
    "mintable": false,
    "burnable": false,
    "addr_forbiddable": false,
    "token_forbiddable": false,
    "total_burn": "0",
    "total_mint": "0",
    "is_forbidden": false
  }
}
```



