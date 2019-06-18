## Account-Query

查询账户时增加锁定金额，memorequired等字段的显示。

## IssueToken CLI & API

- CLI命令
  - `$ cetcli query account addr [flags]`
- Rest-curl命令
  - `$ curl -X GET http://localhost:1317/auth/accounts/{address}`

## Cli Example

- 查询账户

```
./cetcli query account coinex17w3sdtvtddqqw2ctydk36ua2xu7qj3wz4a5xzm --chain-id=coinexdex

Account:
  Address:      coinex7w3sdtvtddqqw2ctydk36ua2xu7qj3wz4a5xzm
  Pubkey:        
  Coins:         100000000cet
  AccountNumber: 1
  Sequence:      0
  LockedCoins:   coin: 200000000eth, unlocked_time: 1658670138
                 coin: 100000000eth, unlocked_time: 1658670138
  FrozenCoins:   
  MemoRequired:  false
```



## Rest Example

- 查询账户

```
curl -X GET http://localhost:1317/auth/accounts/cosmos17w3sdtvtddqqw2ctydk36ua2xu7qj3wz4a5xzm

{
    "address": "coinex17w3sdtvtddqqw2ctydk36ua2xu7qj3wz4a5xzm",
    "coins": [
        {
            "denom": "cet",
            "amount": "100000000"
        }
    ],
    "locked_coins": [
        {
            "coin": {
                "denom": "eth",
                "amount": "200000000"
            },
            "unlock_time": "1658670138"
        },
        {
            "coin": {
                "denom": "eth",
                "amount": "100000000"
            },
            "unlock_time": "1658670138"
        }
    ],
    "frozen_coins": null,
    "public_key": null,
    "account_number": "1",
    "sequence": "0",
    "memo_required": false
}
```

- 查询余额

```
curl -X GET http://localhost:1317/bank/balances/cosmos17w3sdtvtddqqw2ctydk36ua2xu7qj3wz4a5xzm

{
    "coins": [
        {
            "denom": "cet",
            "amount": "100000000"
        }
    ],
    "locked_coins": [
        {
            "coin": {
                "denom": "eth",
                "amount": "200000000"
            },
            "unlock_time": "1658670138"
        },
        {
            "coin": {
                "denom": "eth",
                "amount": "100000000"
            },
            "unlock_time": "1658670138"
        }
    ]
}
```

