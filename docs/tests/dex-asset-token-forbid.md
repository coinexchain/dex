# Asset-Forbid-Token

## ForbidToken

- 只有token的owner可以进行token forbid
- 不能对未发行token进行此操作
- 只有具备可全局禁止能力的token才能进行此操作
- 已经被forbid的token需要解除禁止后才能再次forbid

## ForbidToken CLI

- `$ cetcli tx asset forbid-token [flags]`
- `$ cetcli tx asset unforbid-token [flags]`

## ForbidToken CLI Example

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
      "total_burn": "100",
      "total_mint": "100",
      "is_forbidden": false,
      "url": "www.token2.org",
      "description": "token2 is a example token"
    }
  }
]
```

2. 通过cli对token2进行forbid

```bash
$ cetcli tx asset forbid-token --symbol="token2" \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，token禁止成功

```bash
Response:
  Height: 438
  TxHash: 3E89659CCEE63E8BE7CCF68E8D56C2456A36E7F00EDE97091B2BF71155BE7885
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 21731
  Tags:
    - action = forbid_token
    - category = asset
    - token = token2

  Timestamp: 2019-06-26T01:38:15Z
```

3. 此时查看token2信息，is_forbidden已经更新

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
    "is_forbidden": true,
    "url": "www.token2.org",
    "description": "token2 is a example token"
  }
}
```

4. 解除token2的禁止状态

```bash
$ cetcli tx asset unforbid-token --symbol="token2" \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，解除禁止成功

```bash
Response:
  Height: 463
  TxHash: B016CDD53EEF7576C0B29BCF3230F2867BE29E4C3CA272DE10BDD473CBE38809
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 21677
  Tags:
    - action = unforbid_token
    - category = asset
    - token = token2

  Timestamp: 2019-06-26T01:40:21Z
```

5. 此时查看token2信息，is_forbidden已经更新

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
