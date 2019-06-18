# Memo测试

## Memo ClI
1、跑通[单节点测试](single_node_test.md)

```
$ ./cetd start
```

2、创建新账号

```bash
$ ./cetcli keys add charles <<<$'12345678\n'
```

3、激活新账号

```bash
$ ./cetcli tx send $(./cetcli keys show charles -a) 1000000000cet --from bob --chain-id=coinexdex --gas 200000 <<<$'Y\n12345678\n'
{"chain_id":"coinexdex","account_number":"0","sequence":"11","fee":{"amount":null,"gas":"200000"},"msgs":[{"type":"cosmos-sdk/MsgSend","value":{"from_address":"coinex1fu24efn6syyvlem6t5zw7432t9tpa4nyedvw78","to_address":"coinex1jkcwep7zkvgdwg3nfe5q637qz6n75tv20uvktk","amount":[{"denom":"cet","amount":"10"}]}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: Y
Password to sign with 'bob':
Response:
  TxHash: D0258701B92C07206DEEAA6D8521BAF47EBDB6A84C7B7382ECA44B036C29E4FB
  
$ ./cetcli query tx D0258701B92C07206DEEAA6D8521BAF47EBDB6A84C7B7382ECA44B036C29E4FB --chain-id=coinexdex
Response:
  Height: 1113
  TxHash: D0258701B92C07206DEEAA6D8521BAF47EBDB6A84C7B7382ECA44B036C29E4FB
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 41021
  Tags: 
    - action = send
    - sender = cosmos1fu24efn6syyvlem6t5zw7432t9tpa4nyedvw78
    - recipient = cosmos1jkcwep7zkvgdwg3nfe5q637qz6n75tv20uvktk

  Timestamp: 2019-05-22T05:56:44Z

```

4、设置memo为必须

```bash
$ ./cetcli tx require-memo true --from charles --chain-id=coinexdex --gas 100000 <<<$'Y\n12345678\n'
{"chain_id":"coinexdex","account_number":"1","sequence":"0","fee":{"amount":null,"gas":"100000"},"msgs":[{"type":"cet-chain/MsgSetMemoRequired","value":{"address":"coinex1jkcwep7zkvgdwg3nfe5q637qz6n75tv20uvktk","required":true}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: Y
Password to sign with 'charles':
Response:
  TxHash: 2B11A19F375CE4ACAC8E55B8698D0E33866CC1DBBF1D5EF6E9AC1B2FAC2CEB03

$ ./cetcli query tx 2B11A19F375CE4ACAC8E55B8698D0E33866CC1DBBF1D5EF6E9AC1B2FAC2CEB03 --chain-id=coinexdex
Response:
  Height: 1208
  TxHash: 2B11A19F375CE4ACAC8E55B8698D0E33866CC1DBBF1D5EF6E9AC1B2FAC2CEB03
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 100000
  GasUsed: 12031
  Tags: 
    - action = set_memo_required

  Timestamp: 2019-05-22T06:04:43Z
```

5、转账失败（没有memo）

```bash
$ ./cetcli tx send $(./cetcli keys show charles -a) 1000000000cet --from bob --chain-id=coinexdex --gas 200000 <<<$'Y\n12345678\n'
{"chain_id":"coinexdex","account_number":"0","sequence":"12","fee":{"amount":null,"gas":"200000"},"msgs":[{"type":"cosmos-sdk/MsgSend","value":{"from_address":"coinex1fu24efn6syyvlem6t5zw7432t9tpa4nyedvw78","to_address":"coinex1jkcwep7zkvgdwg3nfe5q637qz6n75tv20uvktk","amount":[{"denom":"cet","amount":"10"}]}}],"memo":""}

confirm transaction before signing and broadcasting [Y/n]: Y
Password to sign with 'bob':
Response:
  TxHash: 2D14C8BFD59F6FE6E98B437A82F4D747859E6228E03CD82763A1159977F5ED2C
  Code: 112
  Raw Log: {"codespace":"bankx","code":112,"message":"memo is empty"}
```

6、转账成功（有memo）

```bash
$ ./cetcli tx send $(./cetcli keys show charles -a) 1000000000cet --from bob --chain-id=coinexdex --gas 200000 --memo hello <<<$'Y\n12345678\n'
{"chain_id":"coinexdex","account_number":"0","sequence":"12","fee":{"amount":null,"gas":"200000"},"msgs":[{"type":"cosmos-sdk/MsgSend","value":{"from_address":"coinex1fu24efn6syyvlem6t5zw7432t9tpa4nyedvw78","to_address":"coinex1jkcwep7zkvgdwg3nfe5q637qz6n75tv20uvktk","amount":[{"denom":"cet","amount":"10"}]}}],"memo":"hello"}

confirm transaction before signing and broadcasting [Y/n]: Y
Password to sign with 'bob':
Response:
  TxHash: F5F77C23C2DAADF362A7618556A60CE2ED8D4BF9740F23A328543BC305DF5BA7

$ ./cetcli query tx F5F77C23C2DAADF362A7618556A60CE2ED8D4BF9740F23A328543BC305DF5BA7 --chain-id=coinexdex
Response:
  Height: 1244
  TxHash: F5F77C23C2DAADF362A7618556A60CE2ED8D4BF9740F23A328543BC305DF5BA7
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 28240
  Tags: 
    - action = send
    - sender = cosmos1fu24efn6syyvlem6t5zw7432t9tpa4nyedvw78
    - recipient = cosmos1jkcwep7zkvgdwg3nfe5q637qz6n75tv20uvktk

  Timestamp: 2019-05-22T06:07:45Z
```



7、导出状态（需要停止cetd）：

```bash
$ /.cetd export
```

主要观察accountsx字段：

```json
{
  "genesis_time": "2019-05-23T10:06:38.504954Z",
  "chain_id": "coinexdex",
  "consensus_params": { /*...*/ },
  "validators": [ /*...*/ ],
  "app_hash": "",
  "app_state": {
    "accounts": [ /*...*/ ],
    "accountsx": [
      {
        "address": "coinex1fjn9htuylwclsflh4cnuhzl0e4jp6fem6dkwym",
        "activated": true,
        "memo_required": true,
        "locked_coins": null
      }
    ],
    "auth": { /*...*/ },
    "bank": { /*...*/ },
    "bankx": { /*...*/ },
    "staking": { /*...*/ },
    "distr": { /*...*/ },
    "gov": { /*...*/ },
    "crisis": { /*...*/ },
    "slashing": { /*...*/ },
    "asset": { /*...*/ }
}
```

## Memo Rest API
1. 通过Rest API设置地址required_memo
```bash
$ curl -X POST "http://localhost:1317/bank/accounts/memo" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"coinex1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t\", \"chain_id\": \"coinexdex\", \"account_number\": \"2\", \"sequence\": \"0\" }, \"memo_required\": true}" >> unsign.json
```
本地生成unsign.json
```bash
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "cet-chain/MsgSetMemoRequired",
        "value": {
          "address": "coinex1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t",
          "required": true
        }
      }
    ],
    "fee": {
      "amount": null,
      "gas": "200000"
    },
    "signatures": null,
    "memo": ""
  }
}
```
2. 本地进行签名
```bash
$ cetcli tx sign \
  --chain-id=coinexdex \
  --from $(cetcli keys show alice -a)  unsigne.json > signedSendTx.json
```
3. 广播已签名交易
```bash
$ cetcli tx broadcast signedSendTx.json
```
4. 测试不带memo转账
```bash
$ cetcli tx send $(cetcli keys show alice -a) 1000000000cet --from bob --chain-id=coinexdex --gas 200000
```
转账失败，memo is empty
```bash
Response:
  TxHash: AA3D128CBD68C526AF05FC08C20EEFBD6A0554FD20717D8707F220A67F1784CA
  Code: 112
  Raw Log: {"codespace":"bankx","code":112,"message":"memo is empty"}
```
5. 测试带memo转账
```bash
$ cetcli tx send $(cetcli keys show alice -a) 1000000000cet --from bob --chain-id=coinexdex --gas 200000 --memo hello
```
转账成功
```bash
Response:
  TxHash: F8CE20B8F7B33431577BBDDA77CF084E9470ABC1C199F97900E9CF3DBDDDA995
```
TX状态success
```bash
Response:
  Height: 13063
  TxHash: F8CE20B8F7B33431577BBDDA77CF084E9470ABC1C199F97900E9CF3DBDDDA995
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 37891
  Tags:
    - action = send
    - sender = cosmos1psmd30v4q47qqgm788mffmx46g49k7afz2nvvp
    - recipient = cosmos1yvnrsxp6cagema97m4uf7vgvh4mcpl9csups2t

  Timestamp: 2019-05-28T10:32:07Z
```














