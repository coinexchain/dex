## Transaction send 

#### 需求罗列

- 转账可以设置该笔转账资金的解锁定时间
- 任何资产都可以被锁定
- 锁定需要收取100cet的功能费
- 区块自动解锁定到达解锁时间的资金

#### 测试过程

1. 初始化并启动链

```
./cetd init node0 --chain-id=coinexdex
./cetcli keys add bob
./cetd add-genesis-account $(./cetcli keys show bob -a) 10000000000000000cet
./cetd gentx --name bob
./cetd collect-gentxs
./cetd start
```

2. 创建新账号

```
$ ./cetcli keys add bear

NAME:	TYPE:	ADDRESS:					PUBKEY:
bear	local	cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9	cosmospub1addwnpepqtq2clys3yl8vvec46ujp7lsjta4akphe3ys3kcpewfukcwnc59rsy23rjj
```

3. 发行cet

```
./cetcli tx asset issue-token --name="Cet" --name="Cet" --symbol="cet" --total-supply=2100000000000000 --mintable=false --burnable=true --addr-forbiddable=0 --token-forbiddable=0 --from $(./cetcli keys show bob -a) --chain-id=coinexdex --fees 200cet --gas 200000
{"chain_id":"coinexdex","account_number":"0","sequence":"1","fee":{"amount":[{"denom":"cet","amount":"200"}],"gas":"200000"},"msgs":[{"type":"asset/MsgIssueToken","value":{"name":"Cet","symbol":"cet","total_supply":"2100000000000000","owner":"cosmos1yah65tmduggvvhpldkmtmc2ldfsw9lrht5l3qj","mintable":false,"burnable":true,"addr_forbiddable":false,"token_forbiddable":false}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: C925838DBDDCAA4ACBF58E66C21D780EB61EA1A6A684102C1562EFC41C92EB2C
```

3. 查询发行结果

```
./cetcli query tx C925838DBDDCAA4ACBF58E66C21D780EB61EA1A6A684102C1562EFC41C92EB2C --chain-id=coinexdex
Response:
  Height: 9
  TxHash: C925838DBDDCAA4ACBF58E66C21D780EB61EA1A6A684102C1562EFC41C92EB2C
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 36263
  Tags: 
    - action = issue_token
    - category = asset
    - token = cet
    - owner = cosmos1yah65tmduggvvhpldkmtmc2ldfsw9lrht5l3qj

  Timestamp: 2019-06-12T09:25:30Z
```

3. 激活账号

```
./cetcli tx send $(./cetcli keys show bear -a) 200000000cet --unlock-time=0 --from bob --chain-id=coinexdex --gas 200000 --fees 100cet
{"chain_id":"coinexdex","account_number":"0","sequence":"2","fee":{"amount":[{"denom":"cet","amount":"100"}],"gas":"200000"},"msgs":[{"type":"bankx/MsgSend","value":{"from_address":"cosmos1yah65tmduggvvhpldkmtmc2ldfsw9lrht5l3qj","to_address":"cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9","amount":[{"denom":"cet","amount":"200000000"}],"unlock_time":"0"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: D2C69136388F9861221FD16AE51841B510ED01472541BBB44BCAD0EFCDDB59D5
```

4. 查询激活结果

```
./cetcli query account cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9 --chain-id=coinexdex
Account:
  Address:       cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9
  Pubkey:        
  Coins:         100000000cet
  AccountNumber: 1
  Sequence:      0
```

4. 锁定转账，时间设为未来3min后

```
./cetcli tx send $(./cetcli keys show bear -a) 200000000cet --unlock-time=1560332036 --from bob --chain-id=coinexdex --gas 200000 --fees 200cet
{"chain_id":"coinexdex","account_number":"0","sequence":"3","fee":{"amount":[{"denom":"cet","amount":"200"}],"gas":"200000"},"msgs":[{"type":"bankx/MsgSend","value":{"from_address":"cosmos1yah65tmduggvvhpldkmtmc2ldfsw9lrht5l3qj","to_address":"cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9","amount":[{"denom":"cet","amount":"200000000"}],"unlock_time":"1560332036"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: 0A74B72F92D56018A7BEEAE60C6D30E513BBEC9E45F9114221E9744C76373AB9
```

5. 锁定转账，时间设置到很远的未来

```
./cetcli tx send $(./cetcli keys show bear -a) 200000000cet --unlock-time=1570332036 --from bob --chain-id=coinexdex --gas 200000 --fees 200cet
{"chain_id":"coinexdex","account_number":"0","sequence":"4","fee":{"amount":[{"denom":"cet","amount":"200"}],"gas":"200000"},"msgs":[{"type":"bankx/MsgSend","value":{"from_address":"cosmos1yah65tmduggvvhpldkmtmc2ldfsw9lrht5l3qj","to_address":"cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9","amount":[{"denom":"cet","amount":"200000000"}],"unlock_time":"1570332036"}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: y
Password to sign with 'bob':
Response:
  TxHash: 9F51F25BF384668EDF4F652AEE038D3639B486F6BA873F276ECE69A5F53A9060
```

5. 查询账户

```
./cetcli query account cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9 --chain-id=coinexdex
Account:
  Address:       cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9
  Pubkey:        
  Coins:         100000000cet
  AccountNumber: 1
  Sequence:      0
  LockedCoins:   coin: 200000000cet, unlocked_time: 1560332036
                 coin: 200000000cet, unlocked_time: 1570332036
  FrozenCoins:   
  MemoRequired:  false
```

6. 3分钟后查询账户

```
./cetcli query account cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9 --chain-id=coinexdex
Account:
  Address:       cosmos1gxplu0twrg7xd503lmzqv76vq9m54ty9u74px9
  Pubkey:        
  Coins:         300000000cet
  AccountNumber: 1
  Sequence:      0
  LockedCoins:   coin: 200000000cet, unlocked_time: 1570332036
  FrozenCoins:   
  MemoRequired:  false
```

验证成功