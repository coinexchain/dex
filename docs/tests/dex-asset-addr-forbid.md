# Asset-Forbid-Addr

## Forbid Addr

- 只有token的owner可以进行forbid addr
- 不能对未发行token进行此操作
- 只有具备可addr禁止能力的token才能进行此操作
- 不能添加和删除空地址

## Forbid Addr CLI

- `$ cetcli tx asset forbid-addr [flags]`
- `$ cetcli tx asset unforbid-addr [flags]`
- `$ cetcli q asset forbid-addr abc [flags]`

## Forbid Addr CLI Example

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

2. 通过cli添加token2的禁止Address

```bash
$ cetcli tx asset forbid-addr --symbol="token2" \
        --addresses=coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47,coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g,coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

本地返回TxHash：

```bash
Response:
  Height: 796
  TxHash: 9E060C07E2A1F4C87FE635B4832B66272CB38C2E0AF4B5BE712FD492604E2676
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 22971
  Tags:
    - action = forbid_addr
    - category = asset
    - token = token2
    - addresses = coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47,coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g,coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn,

  Timestamp: 2019-06-26T02:08:20Z
```

3. 此时可以查看到token2被禁止的addr

```bash
$ cetcli q asset forbid-addr token2 --chain-id=coinexdex
```

本地返回：

```bash
[
  "coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g",
  "coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn",
  "coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47"
]
```

4. 解除被禁止的addr

```bash
$ cetcli tx asset unforbid-addr --symbol="token2" \
        --addresses=coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g,coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

本地返回TxHash：

```bash
Response:
  Height: 826
  TxHash: 1BCEFC7BF0B6591D60EA1280FA40792409E397140144A3E58DF461E308101859
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 18751
  Tags:
    - action = unforbid_addr
    - category = asset
    - token = token2
    - addresses = coinex1p9ek7d3r9z4l288v4lrkwwrnh9k5htezk2q68g,coinex1sxdg68j29l057a7utz7hy9pztdv94a3gsw98hn,

  Timestamp: 2019-06-26T02:10:52Z
```

5. 此时查看token2被禁止的addr

```bash
$ cetcli q asset forbid-addr token2 --chain-id=coinexdex

```

本地返回：

```bash
[
  "coinex1ekevrsx6s853fqjt6rln9r84u8cwuft7e4wp47"
]
```

