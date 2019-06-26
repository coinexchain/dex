# Asset-Mint-Token

## MintToken

- 只有token的owner可以进行token增发
- 不能对未发行token进行此操作
- 只有具备可增发能力的token才能进行此操作
- 增发后的token总量不能超过90billion

## Mint Token CLI

- `$ cetcli tx asset mint-token [flags]`

## Mint Token CLI Example

参考[Asset-IssueToken](https://github.com/coinexchain/dex/blob/master/docs/tests/dex-asset-issue.md)创建测试token

节点启动后

1. 查询本地所有token

```bash
$ cetcli query asset tokens --chain-id=coinexdex
```

本地返回所有的token：

```bash
[
  {
    "type": "asset/BaseToken",
    "value": {
      "name": "ABC Token",
      "symbol": "abc",
      "total_supply": "2100000000000000",
      "owner": "coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765",
      "mintable": true,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false,
      "url": "www.abc.org",
      "description": "token abc is a example token"
    }
  },
  {
    "type": "asset/BaseToken",
    "value": {
      "name": "First Token",
      "symbol": "token1",
      "total_supply": "2100000000000000",
      "owner": "coinex1h6jte3avry5q5fnn6cyzfh65r74tf7tmxdfxu6",
      "mintable": true,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false,
      "url": "www.token1.org",
      "description": "token1 is a example token"
    }
  },
  {
    "type": "asset/BaseToken",
    "value": {
      "name": "Second Token",
      "symbol": "token2",
      "total_supply": "2100000000000000",
      "owner": "coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765",
      "mintable": true,
      "burnable": true,
      "addr_forbiddable": true,
      "token_forbiddable": true,
      "total_burn": "0",
      "total_mint": "0",
      "is_forbidden": false,
      "url": "www.token2.org",
      "description": "token2 is a example token"
    }
  }
]
```

3. 通过cli进行token2的增发

```bash
$ cetcli tx asset mint-token --symbol="token2" \
        --amount=100 \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，增发成功

```bash
Response:
  Height: 331
  TxHash: D7533A7DBE91539C55EB3C71EAB04208240984D93FF363C95367445B8312F021
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 21619
  Tags:
    - action = mint_token
    - category = asset
    - token = token2
    - amount = 100

  Timestamp: 2019-06-26T01:29:16Z
```

4. 此时查看toekn2信息，totalsupply已经增发100

```bash
$ cetcli q asset token token2 --chain-id=coinexdex
```

本地返回：

```bash
{
  "type": "asset/BaseToken",
  "value": {
    "name": "Second Token",
    "symbol": "token2",
    "total_supply": "2100000000000100",
    "owner": "coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765",
    "mintable": true,
    "burnable": true,
    "addr_forbiddable": true,
    "token_forbiddable": true,
    "total_burn": "0",
    "total_mint": "100",
    "is_forbidden": false,
    "url": "www.token2.org",
    "description": "token2 is a example token"
  }
}
```
