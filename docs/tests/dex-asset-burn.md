# Asset-Burn-Token

## BurnToken

- 只有token的owner可以进行token燃烧
- 不能对未发行token进行此操作
- 只有具备可燃烧能力的token才能进行此操作
- 燃烧后的token总量不能为负

## Burn Token CLI 

- `$ cetcli tx asset burn-token [flags]`

## Burn Token CLI Example

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
]
```

3. 通过cli进行token2的燃烧

```bash
$ cetcli tx asset burn-token --symbol="token2" \
        --amount=100 \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，燃烧成功

```bash
Response:
  Height: 380
  TxHash: C6A804A939FF8473FC3EF3B818D73771169D91FC82F5EB4BEF23BF7FFE7FD58B
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 21685
  Tags:
    - action = burn_token
    - category = asset
    - token = token2
    - amount = 100

  Timestamp: 2019-06-26T01:33:23Z
```

4. 此时查看token2信息，totalsupply已经减少100

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
    "total_supply": "2100000000000000",
    "owner": "coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765",
    "mintable": true,
    "burnable": true,
    "addr_forbiddable": true,
    "token_forbiddable": true,
    "total_burn": "100",
    "total_mint": "100",
    "is_forbidden": false,
    "url": "www.token2.org",
    "description": "token2 is a example token"
  }
}
```
