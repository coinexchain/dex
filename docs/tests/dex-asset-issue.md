# Asset-IssueToken

## IssueToken

- Token Name 限制在32 unicode characters
- Token Symbol format `^[a-z][a-z0-9]{1,7}$`
- Token Symbol 不能重名，创建已有的token
- Token TotalSupply 在boosting前不能超过90 billion
- Token issue fee扣除1000000000000cet
- Token url限制在100 unicode characters
- Token description size限制为1k

## IssueToken CLI

- `$ cetcli tx asset issue-token [flags]`
- `$ cetcli query asset token [symbol] [flags]`
- `$ cetcli query asset tokens [flags]`

## IssueToken CLI Example

节点的搭建参考[single_node_test](https://github.com/coinexchain/dex/blob/master/docs/tests/single-node-test.md)，也可以运行`scripts/setup_single_testing_node.sh`，

节点启动后

1. 尝试查询当前token-list

```bash
$ cetcli query asset tokens --chain-id=coinexdex
```

当前没有已发行的token

```bash
[]
```

2. 查询当前用户

```bash
$ cetcli q account $(cetcli keys show bob -a) --chain-id=coinexdex
```

本地返回：

```bash
Account:
  Address:       coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765
  Pubkey:        coinexpub1addwnpepqvfkkum5w78h52xr8dzy347rp4sa79dh2h74tqpgndmvnn5tpcl4x99yzg4
  Coins:         9900000000000000cet
  AccountNumber: 0
  Sequence:      1
  LockedCoins:
  FrozenCoins:
  MemoRequired:  false
```

3. 创建token，通过--Flag指定所要创建的token信息

```bash
$ cetcli tx asset issue-token --name="ABC Token" \
        --symbol="abc" \
        --total-supply=2100000000000000 \
        --mintable=true \
        --burnable=true \
        --addr-forbiddable=true \
        --token-forbiddable=true \
        --url="www.abc.org" \
        --description="token abc is a example token" \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 40000 --gas-prices 20.0cet
```

解析本地返回TxHash，token发行成功

```bash
Response:
  Height: 24
  TxHash: 7549E7DF4FAAB34A6A33F09CB27759E07194323F2B3550FEF71BD883DB04D9C9
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 40000
  GasUsed: 39312
  Tags:
    - action = issue_token
    - category = asset
    - token = abc
    - owner = coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765
```

4. 查询bob账户，已扣除1000000000000sato.cet功能 fee+ 800000sato.cet的gas fee

```bash
Account:
  Address:       coinex1r550scev7m4qg7nv662v5vu0at7kvejagls765
  Pubkey:        coinexpub1addwnpepqvfkkum5w78h52xr8dzy347rp4sa79dh2h74tqpgndmvnn5tpcl4x99yzg4
  Coins:         2100000000000000abc,9898999999200000cet
  AccountNumber: 0
  Sequence:      2
  LockedCoins:
  FrozenCoins:
  MemoRequired:  false
```

5. 如果再次重复创建symbol:"abc"，会创建失败

```bash
$ cetcli tx asset issue-token --name="ABC Token" \
        --symbol="abc" \
        --total-supply=2100000000000000 \
        --mintable=true \
        --burnable=true \
        --addr-forbiddable=true \
        --token-forbiddable=true \
        --url="www.abc.org" \
        --description="token abc is a example token" \
    --from $(cetcli keys show bob -a) \
    --chain-id=coinexdex \
    --gas 40000 --gas-prices 20.0cet
```

本地返回提示：

```bash
ERROR: token symbol already exists，please query tokens and issue another symbol
```

5. 此时可以传入symbol参数`abc`查询到token-info

```bash
$ cetcli query asset token abc --chain-id=coinexdex
```

本地返回token的信息：

```bash
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
}
```

6. 如上可以创建token1，token2等token，可以查询所有的token

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