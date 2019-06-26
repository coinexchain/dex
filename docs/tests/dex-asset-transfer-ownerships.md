# Asset-TransferOwnership

## TransferOwnership

- 只有token的owner可以进行ownership的转移
- 不能对未发行token进行此操作
- new_owner不能为空
- 不能向自己transfer ownership

## TransferOwnership CLI

- `$ cetcli tx asset transfer-ownership [flags]`

## TransferOwnership CLI Example

参考[Asset-IssueToken](https://github.com/coinexchain/dex/blob/master/docs/tests/dex-asset-issue.md)创建测试token

节点启动后

1. 添加alice地址

```bash
$  cetcli keys add alice
```

2. 查询本地账户

```bash
$  cetcli	keys list
```

本地返回：

```bash
NAME:	TYPE:	ADDRESS:					PUBKEY:
alice	local	coinex1h6jte3avry5q5fnn6cyzfh65r74tf7tmxdfxu6	coinexpub1addwnpepqguh9czz68trlt00fu5zvnlzghqzujudmagc563uxfuzdpl29hphy5mlqlq
bob	local	coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765	coinexpub1addwnpepqvfkkum5w78h52xr8dzy347rp4sa79dh2h74tqpgndmvnn5tpcl4x99yzg4
```

3. 查询本地所有token

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
      "owner": "coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765",
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

4. 如上3个token的owner都是bob，现在通过cli转移ownership

```bash
$ cetcli tx asset transfer-ownership --symbol="token1" \
	--new-owner $(cetcli keys show alice -a) \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 50000 --gas-prices 20.0cet
```

解析返回TxHash，transfer ownership成功

```bash
Response:
  Height: 256
  TxHash: D87A922A5AE08CFD12336F03F4C242993DA20305D1B0936A8664B2A9239ABAD5
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 50000
  GasUsed: 21726
  Tags:
    - action = transfer_ownership
    - category = asset
    - token = token1
    - original-owner = coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765
    - new-owner = coinex1h6jte3avry5q5fnn6cyzfh65r74tf7tmxdfxu6

  Timestamp: 2019-06-26T01:22:58Z
```

5. 此时查看token1信息，owner已经变成alice

```bash
$ cetcli q asset token token1 --chain-id=coinexdex
```

本地返回：

```bash
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
}
```
