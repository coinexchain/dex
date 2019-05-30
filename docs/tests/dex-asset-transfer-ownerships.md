# Asset-TransferOwnership

## TransferOwnership

- 添加Asset模块的TransferOwnership功能，支持token ownership的转移。
  - token的owner可以进行ownership的转移
  - 非owner不能进行此操作
  - 不能对未发行token进行此操作
  - new_owner不能为空

> transferOwnership扣除fee，暂时没有在asset 模块实现，待评估ante-Handler统一收取fee
>
> transferOwnership fee未确认，需要和coinex对齐

## TransferOwnership CLI & API

- CLI命令
  - `$ cetcli tx asset transfer-ownership [flags]` 
- Rest-curl命令
  - `curl -X POST http://localhost:1317/asset/tokens/coin2/ownerships --data-binary '{"base_req":{"from":"cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp","chain_id":"coinexdex","sequence":"5","account_number":"0"},"new_owner":"cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t"}'`

## TransferOwnership CLI Example

参考[single_node_test](https://gitlab.com/cetchain/docs/blob/master/dex/tests/single_node_test.md)搭建节点，也可以从genesis.json中导入状态，节点启动后

1. 查询本地bob地址

```bash
cetcli keys show bob -a
```

本地返回：

```bash
cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp
```

2. 查询本地alice地址

```bash
$ cetcli keys show alice -a
```

本地返回：

```bash
cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t
```

3. 本地创建token，可参考[dex-asset-iusse](https://gitlab.com/cetchain/docs/blob/master/dex/tests/dex-asset-issue.md) 

```bash
$ cetcli tx asset issue-token --name="bob first token" \
        --symbol="coin1" \
        --total-supply=2100000000000000 \
        --mintable=false \
        --burnable=true \
        --addr-forbiddable=0 \
        --token-forbiddable=1 \
        --from $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: DA1EC4886B2469A58A0E3713DF8EA6760CAEC4A9F42B8EE11710F99AA44BF92A
```

4. 如上可以创建coin2，coin3等token，查询下所有token信息

```bash
$ cetcli query asset tokens --chain-id=coinexdex
```

本地返回所有的token：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": "bob first token",
      "symbol": "coin1",
      "total_supply": "2100000000000000",
      "owner": "cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_frozen": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": "bob sec token",
      "symbol": "coin2",
      "total_supply": "2100000000000000",
      "owner": "cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_frozen": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": "bob th token",
      "symbol": "coin3",
      "total_supply": "2100000000000000",
      "owner": "cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_frozen": false
    }
  }
]
```

5. 如上3个token的owner都是bob，现在通过cli转移ownership

```bash
$ cetcli tx asset transfer-ownership --symbol="coin1" \
	--new-owner $(cetcli keys show alice -a) \
    --from $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回TxHash：

```bash
Response:
  TxHash: C66A5DFB5CCAAB8F9A2BE039DAC9E3DFDDEACF044E9943760E9E71730B3B88A1
```

6. 此时查看coin1信息，owner已经变成alice

```bash
$ cetcli q asset token coin1 --chain-id=coinexdex
```

本地返回：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": "bob first token",
    "symbol": "coin1",
    "total_supply": "2100000000000000",
    "owner": "cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t",
    "mintable": false,
    "burnable": true,
    "addr_forbiddable": false,
    "token_forbiddable": true,
    "total_burn": "0",
    "total_mint": "0",
    "is_frozen": false
  }
}
```



## TransferOwnership Rest Example

1. 查询本地AccountNumber和Sequence

```bash
$ cetcli query account $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回：

```bash
Account:
  Address:       cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp
  Pubkey:        cosmospub1addwnpepqgrp6tj3j8507jveu2jmgcp6adz8t95gpfuxfeun6032a6emu2s2g23q55r
  Coins:         9996999900000000cet,2100000000000000coin1,2100000000000000coin2,2100000000000000coin3
  AccountNumber: 0
  Sequence:      5
```

2. 首先需要启动rest-server.  参考[本地rest-server中访问swagger-ui的方法](https://gitlab.com/cetchain/docs/blob/master/dex/tests/dex_rest_api_swagger.md)

```bash
$ cetcli rest-server --chain-id=coinexdex \ --laddr=tcp://localhost:1317 \ --node tcp://localhost:26657 --trust-node=false
```

3. 可以通过GET查询token信息

```bash
$ curl -X GET http://localhost:1317/asset/tokens | jq
```

返回token1的信息：

```bash
[
  {
    "type": "asset/Token",
    "value": {
      "name": "bob first token",
      "symbol": "coin1",
      "total_supply": "2100000000000000",
      "owner": "cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_frozen": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": "bob sec token",
      "symbol": "coin2",
      "total_supply": "2100000000000000",
      "owner": "cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_frozen": false
    }
  },
  {
    "type": "asset/Token",
    "value": {
      "name": "bob th token",
      "symbol": "coin3",
      "total_supply": "2100000000000000",
      "owner": "cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp",
      "mintable": false,
      "burnable": true,
      "addr_forbiddable": false,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_frozen": false
    }
  }
]
```

4. 通过Rest API转移coin2的ownership，填写本地from/new_owner/sequence/account_number等信息

```bash
$ curl -X POST http://localhost:1317/asset/tokens/coin2/ownerships --data-binary '{"base_req":{"from":"cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp","chain_id":"coinexdex","sequence":"5","account_number":"0"},"new_owner":"cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t"}' > unsigned.json
```

返回未签名交易存入unsigned.json

```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "asset/MsgTransferOwnership",
        "value": {
          "Symbol": "coin2",
          "OriginalOwner": "cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp",
          "NewOwner": "cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t"
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
        "type": "asset/MsgTransferOwnership",
        "value": {
          "Symbol": "coin2",
          "OriginalOwner": "cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp",
          "NewOwner": "cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t"
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
          "value": "AgYdLlGR6P9JmeKltGA660R1logKeGTnk9Pirus74qCk"
        },
        "signature": "Zss8X61FaPqW550Jxqqj7xctIUZq2h78PElavFsOKSZhgnC9SRno0V4DkFVDaGjtb+t8W94sPlW4ArJxvi2PnA=="
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

7. 此时查询coin2的owner已经更改

```bash
$ curl -X GET http://localhost:1317/asset/tokens/coin2
```

返回此token信息：

```bash
{
  "type": "asset/Token",
  "value": {
    "name": "bob sec token",
    "symbol": "coin2",
    "total_supply": "2100000000000000",
    "owner": "cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t",
    "mintable": false,
    "burnable": true,
    "addr_forbiddable": false,
    "token_forbiddable": true,
    "total_burn": "0",
    "total_mint": "0",
    "is_frozen": false
  }
}
```



