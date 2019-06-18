# Bankx-ActivatedFee

在Bankx模块中对基本的转账逻辑进行修改，以保证支持账户首笔交易必须包含（而不是仅有）大于ActivatedFee（默认为1）个CET的资金transfer，这笔激活费用最终被收集到feePool中，随后与交易费一起分发给validator和delegator作为激励。

> 目前实现的逻辑是，在一笔转账的接收账户不存在时，判断这笔转账中是否包含大于ActivatedFee的cet的转账，而不限制第一笔交易必须只有cet的转账。



## Transfer CLI & REST API

基本转账功能涉及的CLI和REST API都在bank模块中提供，无需更改。



## Transfer CLI Example

节点的搭建参考[single_node_test](https://gitlab.com/cetchain/docs/blob/master/dex/single_node_test.md)，在最后生成的genesis.json中，需要确认`bankx`模块下的参数`param`中`ActivatedFee`默认为100000000。

1. 生成新密钥alice：

```bash
./cetcli keys add bob
```

对应地址如下：

```bash
address: cosmos1knpnwr5z24fsc49waqaxw5twxvd33ukwgl9z9u	
```



2. 尝试查询该地址

```bash
/cetcli query account cosmos1knpnwr5z24fsc49waqaxw5twxvd33ukwgl9z9u --trust-node
```

返回结果

```bash
ERROR: {"codespace":"sdk","code":9,"message":"account cosmos1knpnwr5z24fsc49waqaxw5twxvd33ukwgl9z9u does not exist"}
```

说明当前地址对应账户尚未创建。

3. 查询sender账户

```
./cetcli query account $(./cetcli keys show validator -a) --trust-node
```

返回结果

```bash
Account:
  Address:       cosmos1jzvzctg849zesmgq8edy50ufg2yz2hr79rlqkp
  Pubkey:        cosmospub1addwnpepq2mhysu2wpt55wtlknhq8ckqy3qaeeye8vj7gl8pdl9x6xhkxjqgsauxxza
  Coins:         900000000cet
  AccountNumber: 0
  Sequence:      1
```



4. 转账：

```bash
./cetcli tx send cosmos1knpnwr5z24fsc49waqaxw5twxvd33ukwgl9z9u 200000000cet --from validator --chain-id=coinexdex
```

返回response:

```json
TxHash: 397EF98531E4C9BA7A5937A9251F9482786273775888319B342277AD95C14931
```



5. 查询交易

```
./cetcli query tx 397EF98531E4C9BA7A5937A9251F9482786273775888319B342277AD95C14931 --trust-node
```

返回结果：

```bash
Response:
  Height: 10
  TxHash: 397EF98531E4C9BA7A5937A9251F9482786273775888319B342277AD95C14931
  Raw Log: [{"msg_index":"0","success":true,"log":""}]
  Logs: [{"msg_index":0,"success":true,"log":""}]
  GasWanted: 200000
  GasUsed: 39033
  Tags: 
    - action = send
    - sender = cosmos1jzvzctg849zesmgq8edy50ufg2yz2hr79rlqkp
    - recipient = cosmos1knpnwr5z24fsc49waqaxw5twxvd33ukwgl9z9u

  Timestamp: 2019-05-20T07:36:14Z

```

6. 查询bob账户

```bash
./cetcli query account cosmos1knpnwr5z24fsc49waqaxw5twxvd33ukwgl9z9u --trust-node
```

返回结果：

```json
Account:
  Address:       cosmos1knpnwr5z24fsc49waqaxw5twxvd33ukwgl9z9u
  Pubkey:        
  Coins:         100000000cet
  AccountNumber: 1
  Sequence:      0

```

可以看到，bob的账户中只增加了1个CET。

7. 查询sender账户：

```
./cetcli query account $(./cetcli keys show validator -a) --trust-node
```

返回结果：

```bash
Account:
  Address:       cosmos1jzvzctg849zesmgq8edy50ufg2yz2hr79rlqkp
  Pubkey:        cosmospub1addwnpepq2mhysu2wpt55wtlknhq8ckqy3qaeeye8vj7gl8pdl9x6xhkxjqgsauxxza
  Coins:         700000000cet
  AccountNumber: 0
  Sequence:      2
```

可以看到，sender地址中减少了2个CET。





# REST API TEST
1. 创建一个帐户ddd
```bash
./cetcli keys add ddd
```

2. 确认帐户ddd不存在系统中，同样也意味着处于未激活状态
```
BJ00609:~/lab/dex$ ./cetcli query account $(./cetcli keys show ddd -a) --trust-node
ERROR: {"codespace":"sdk","code":9,"message":"account cosmos1aj995p99reua0npfkevkw9v6fv579j449wdk79 does not exist"}
```

3.　转帐发起人bob　的状态查询
```
BJ00609:~/lab/dex$ ./cetcli query account $(./cetcli keys show bob -a) --trust-node
Account:
  Address:       cosmos1esmxfjyanvdevfkrznt00lemhftca7as2vavhk
  Pubkey:        cosmospub1addwnpepqwwrut30t8pawppx2lpl2enayfmar5t59fggkf0k6hrx76empa4ujng6ck2
  Coins:         9999999700000000cet
  AccountNumber: 0
  Sequence:      2
```

4. 启动rest-server
```
./cetcli rest-server --chain-id=coinexdex  --laddr=tcp://localhost:1317  --node tcp://localhost:26657 --trust-node=false
```

5. 访问　http://localhost:1317/swagger/
然后　找　/bank/accounts/{address}/transfers　　这个API
填上　address　　是收款人地址：　`cosmos1aj995p99reua0npfkevkw9v6fv579j449wdk79`
account　是发款人的信息及发款数量，
```

{
  "base_req": {
    "from": "coinex1esmxfjyanvdevfkrznt00lemhftca7as2vavhk",
    "memo": "from bob to ddd",
    "chain_id": "coinexdex",
    "account_number": "0",
    "sequence": "2",
    "gas": "200000",
    "gas_adjustment": "1.2",
    "fees": [
      {
        "denom": "cet",
        "amount": "200000000"
      }
    ],
    "simulate": false
  },
  "amount": [
    {
      "denom": "cet",
      "amount": "1300000000"
    }
  ],
  "unlock_time": "0"
}

```
然后点击Execute.

6. 上一步中对应的CURL请求是：
```
curl -X POST "http://localhost:1317/bank/accounts/cosmos1aj995p99reua0npfkevkw9v6fv579j449wdk79/transfers" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"coinex1esmxfjyanvdevfkrznt00lemhftca7as2vavhk\", \"memo\": \"from bob to ddd\", \"chain_id\": \"coinexdex\", \"account_number\": \"0\", \"sequence\": \"2\", \"gas\": \"200000\", \"gas_adjustment\": \"1.2\", \"fees\": [ { \"denom\": \"cet\", \"amount\": \"200000000\" } ], \"simulate\": false }, \"amount\": [ { \"denom\": \"cet\", \"amount\": \"1300000000\" } ], \"unlock_time\": \"0\"}" > unsignedSendTx.json
```
得到的是未签名的交易：
```
BJ00609:~/lab/dex$ cat unsignedSendTx.json | jq
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "bankx/MsgSend",
        "value": {
          "from_address": "coinex1esmxfjyanvdevfkrznt00lemhftca7as2vavhk",
          "to_address": "coinex1aj995p99reua0npfkevkw9v6fv579j449wdk79",
          "amount": [
            {
              "denom": "cet",
              "amount": "1300000000"
            }
          ],
          "unlock_time": "0"
        }
      }
    ],
    "fee": {
      "amount": [
        {
          "denom": "cet",
          "amount": "200000000"
        }
      ],
      "gas": "200000"
    },
    "signatures": null,
    "memo": "from bob to ddd"
  }
}
```

7. 对其进行签名
```
./cetcli tx sign --chain-id=coinexdex --from=$(cetcli keys show bob -a)  unsignedSendTx.json > signedSendTx.json
```

得到签名后的tx:
```
BJ00609:~/lab/dex$ cat signedSendTx.json | jq
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "bankx/MsgSend",
        "value": {
          "from_address": "coinex1esmxfjyanvdevfkrznt00lemhftca7as2vavhk",
          "to_address": "coinex1aj995p99reua0npfkevkw9v6fv579j449wdk79",
          "amount": [
            {
              "denom": "cet",
              "amount": "1300000000"
            }
          ],
          "unlock_time": "0"
        }
      }
    ],
    "fee": {
      "amount": [
        {
          "denom": "cet",
          "amount": "200000000"
        }
      ],
      "gas": "200000"
    },
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A5w+Li9Zw9cEJlfD9WZ9InfR0XQqUIsl9tXGb2s7D2vJ"
        },
        "signature": "ojy8ymTZxSqPMIZ4DnQNzjDg/YBK2btapyHfv1KTiYZVaZqwD/fvXZVYcLOwnC+Y5YsXJMsUYfREQ5l6sCzCiQ=="
      }
    ],
    "memo": "from bob to ddd"
  }
}
```

8. 发起交易，并查询ddd的状态
```
BJ00609:~/lab/dex$ ./cetcli tx broadcast signedSendTx.json
Response:
  TxHash: 84ECFE1F7A18CD0659890940B592408C490D1DA1916FEFF0333BEEDD7CB1D053
```

可以查到系统中已经有了帐户ddd, 且余额为12$CET. (=转入金额13$CET - 帐户激活费用 1$CET)
```
BJ00609:~/lab/dex$ ./cetcli query account $(./cetcli keys show ddd -a) --trust-node
Account:
  Address:       cosmos1aj995p99reua0npfkevkw9v6fv579j449wdk79
  Pubkey:
  Coins:         1200000000cet
  AccountNumber: 3
  Sequence:      0
```

也可以查看bob的状态： 99999982= 99999987 - 13$CET(转帐金额)　- 2$CET (tx fee)
```
BJ00609:~/lab/dex$ ./cetcli query account $(./cetcli keys show bob -a) --trust-node
Account:
  Address:       cosmos1esmxfjyanvdevfkrznt00lemhftca7as2vavhk
  Pubkey:        cosmospub1addwnpepqwwrut30t8pawppx2lpl2enayfmar5t59fggkf0k6hrx76empa4ujng6ck2
  Coins:         9999998200000000cet
  AccountNumber: 0
  Sequence:      3
BJ00609:~/lab/dex$
```

9. 再次bob转帐给ddd时，　因为ddd已经激活，就不会再次收ddd激活费用了。
```
BJ00609:~/lab/dex$ curl -X POST "http://localhost:1317/bank/accounts/cosmos1aj995p99reua0npfkevkw9v6fv579j449wdk79/transfers" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"coinex1esmxfjyanvdevfkrznt00lemhftca7as2vavhk\", \"memo\": \"from bob to ddd\", \"chain_id\": \"coinexdex\", \"account_number\": \"0\", \"sequence\": \"3\", \"gas\": \"200000\", \"gas_adjustment\": \"1.2\", \"fees\": [ { \"denom\": \"cet\", \"amount\": \"200000000\" } ], \"simulate\": false }, \"amount\": [ { \"denom\": \"cet\", \"amount\": \"1300000000\" } ], \"unlock_time\": \"0\"}" > unsignedSendTx.json
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   741  100   379  100   362  38438  36713 --:--:-- --:--:-- --:--:-- 42111
BJ00609:~/lab/dex$
BJ00609:~/lab/dex$
BJ00609:~/lab/dex$ ./cetcli tx sign --chain-id=coinexdex --from=$(cetcli keys show bob -a)  unsignedSendTx.json > signedSendTx.json
Password to sign with 'bob':
BJ00609:~/lab/dex$ ./cetcli tx broadcast signedSendTx.json
Response:
  TxHash: 3CEF95F2B72B72E3528AE6BB9E52B92A15DC3D5225AD3B6E866FA2C8A85C30C4
BJ00609:~/lab/dex$
BJ00609:~/lab/dex$
BJ00609:~/lab/dex$ ./cetcli query account $(./cetcli keys show ddd -a) --trust-node
Account:
  Address:       cosmos1aj995p99reua0npfkevkw9v6fv579j449wdk79
  Pubkey:
  Coins:         2500000000cet
  AccountNumber: 3
  Sequence:      0
BJ00609:~/lab/dex$
```





