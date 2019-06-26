# Asset-Forbid-Whitelist

## Add/Remove Whitelist

- 只有token的owner可以进行add/Remove whitelist
- 不能对未发行token进行此操作
- 只有具备可全局禁止能力的token才能进行此操作
- 不能添加和删除空地址

## Whitelist CLI

- `$ cetcli tx asset add-whitelist [flags]`
- `$ cetcli tx asset remove-whitelist [flags]`
- `$ cetcli q asset whitelist abc [flags]`

## Whitelist CLI Example

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

2. 通过cli添加token2的whitelist

```bash
$ cetcli tx asset add-whitelist --symbol="token2" \
        --whitelist=coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47,coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g,coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，添加白名单成功

```bash
  Height: 586
  TxHash: 77EE969B535827CEDF11DC0D083BD18DDC4FA0DFBA5B1979AE2CC51DAE3B2A66
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 22971
  Tags:
    - action = add_token_whitelist
    - category = asset
    - token = token2
    - add-whitelist = coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47,coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g,coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn,

  Timestamp: 2019-06-26T01:50:42Z
```

3. 此时可以查看token2白名单

```bash
$ cetcli q asset whitelist token2 --chain-id=coinexdex
```

本地返回：

```bash
[
  "coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g",
  "coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn",
  "coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47"
]
```

4. remove token2白名单中的地址

```bash
$ cetcli tx asset remove-whitelist --symbol="token2" \
        --whitelist=coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47,coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

本地返回TxHash：

```bash
Response:
  Height: 675
  TxHash: 62FEF4E5C92F70685340DE90672D9BFB923E16176C9F6241311FAE096F4B7CB2
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 18751
  Tags:
    - action = remove_token_whitelist
    - category = asset
    - token = token2
    - remove-whitelist = coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g

  Timestamp: 2019-06-26T01:58:10Z
```

5. 此时查看token2的白名单

```bash
$ cetcli q asset whitelist  --chain-id=coinexdex
```

本地返回：

```bash
[
  "coinex1xl6453f6q6dv5770c9ue6hspdc0vxfuqtudkhz"
]
```
