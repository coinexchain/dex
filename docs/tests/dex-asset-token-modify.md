# Asset-Modify-Token

## Modify Token Info

- 只有token的owner可以modify token url 和description信息
- 不能对未发行token进行此操作
- token url 限制在100字符以内
- token description大小限制为1k

## Modify Token Info CLI 

- `$ cetcli tx asset modify-token-url [flags]`
- `$ cetcli tx asset modify-token-description [flags]`

## Modify Token Info CLI Example

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

3. 通过cli修改token2的URL

```bash
$ cetcli tx asset modify-token-url --symbol="token2" \
        --url="www.token2.com" \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，修改成功

```bash
Response:
  Height: 1304
  TxHash: F428D53E3016F1893FB83B9795AA1789B973A0A4A8C2F633C1C1A2C084B3B874
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 21831
  Tags:
    - action = modify_token_url
    - category = asset
    - token = token2
    - url = www.token2.com

  Timestamp: 2019-06-26T09:21:25Z
```

4. 此时查看token2信息，url已经变更

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
    "url": "www.token2.com",
    "description": "token2 is a example token"
  }
}
```

5. 通过cli修改token2的description

```bash
$ cetcli tx asset modify-token-description --symbol="token2" \
        --description="token2 example description" \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，修改成功

```bash
Response:
  Height: 1349
  TxHash: 9DFDC3C7FFC15BC9CCCEDB7924BC7D9B9E264328CA9551F2FEAF279D16F06B51
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 21981
  Tags:
    - action = modify_token_description
    - category = asset
    - token = token2
    - description = token2 example description

  Timestamp: 2019-06-26T09:25:12Z
```

6. 此时查看token2信息，description已经变更

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
    "url": "www.token2.com",
    "description": "token2 example description"
  }
}
```

