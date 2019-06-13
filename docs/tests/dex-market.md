# dex-market

## 启动节点

`./cetd start`

## 发行token 

发行两种token: seth, cet

1. 发行 seth

`./cetcli tx asset issue-token --name="eth" --symbol="seth" --total-supply=2100000000000000 --mintable=false --burnable=true --addr-forbiddable=false --token-forbiddable=false --from bob --chain-id=coinexdex`

2. 发行 cet 

`./cetcli tx asset issue-token --name="cet" --symbol="cet" --total-supply=2100000000000000 --mintable=false --burnable=true --addr-forbiddable=false --token-forbiddable=false --from bob --chain-id=coinexdex`


## 查询token

1. 查询cet
`./cetcli query asset token cet --chain-id=coinexdex`

2. 查询eth
`./cetcli query asset token seth --chain-id=coinexdex`

## 创建交易对市场
创建 eth/cet 对交易市场

`./cetcli tx market createmarket --from bob --chain-id=coinexdex  --gas 60000 --fees=1000cet --stock=seth --money=cet --price-precision=8`


## 查询指定市场信息

`./cetcli query market marketinfo seth/cet --trust-node=true`

## 创建订单

1. 创建**GTE**类型的订单
  `./cetcli tx market creategteoreder --symbol="seth/cet" --order-type=2 --price=520 --quantity=10000000 --side=1 --from bob --price-precision=8 --chain-id=coinexdex --gas=60000 --fees=1000cet` 

  > 这条命令即创建一个seth的买单，money为cet。实际的价格是price sato.cet ，实际买入量是$quantity/10^{price-precision}$，因此，最终该交易会需要订单发起者支付52 sato.cet。

2. 创建**IOC**类型的订单

`./cetcli tx market createiocorder --symbol="seth/cet" --order-type=2 --price=520 --quantity=10000000 --side=1 --from bob --price-precision=8 --chain-id=coinexdex  --gas=60000 --fees=1000cet` `

## 查询指定订单信息

`./cetcli query market orderinfo  --order-id=cosmos16gvnhynu7veexyyaadk60k28cn5s9k7p7p5v9p-13 --trust-node=true`

## 查询指定地址的所有订单列表

`./cetcli query market userorderlist --address=cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z --trust-node=true`

## 取消区块链上的指定订单

`./cetcli tx market cancelorder --order-id=cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z-6 --trust-node=true --from=bob --chain-id=coinexdex`


## Rest API 创建订单

### 使用REST 接口创建未签名的交易

#### 创建未签名的GTE 订单交易

`curl -X POST http://localhost:1317/market/create-gte-order  --data-binary '{"base_req":{"from":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z", "chain_id":"coinexdex"}, "order_type":"2", "symbol":"eth/cet", "price_precision":"8", "price":"32123", "quantity":"1267632", "side":"1"}'  > unsignedSendTx.json`


#### 创建未签名的IOC 订单交易

`curl -X POST http://localhost:1317/market/create-ioc-order  --data-binary '{"base_req":{"from":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z", "chain_id":"coinexdex"}, "order_type":"2", "symbol":"eth/cet", "price_precision":"8", "price":"32782123", "quantity":"77563632", "side":"1"}'  > unsignedSendTx.json`

### 对未签名的交易进行处理

#### 对上步未签名的交易签名

`./cetcli tx sign --chain-id=coinexdex   --from=$(./cetcli keys show bob -a)  unsignedSendTx.json > signedSendTx.json`

#### 发送签名后的交易

`./cetcli tx broadcast signedSendTx.json`


### 查询市场的交易对

`curl -X GET http://localhost:1317/market/market-info --data-binary '{"base_req":{"from":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z", "chain_id":"coinexdex"}, "symbol":"dash/cet" }'`    

### 查询订单

`curl -X GET http://localhost:1317/market/order-info --data-binary '{"base_req":{"from":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z", "chain_id":"coinexdex"}, "order_id":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z-16" }'`

### 取消订单

`curl -X POST http://localhost:1317/market/cancel-order --data-binary '{"base_req":{"from":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z", "chain_id":"coinexdex"}, "order_id":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z-16" }'`

### 查询用户订单列表

 `curl -X GET http://localhost:1317/market/user-order-list --data-binary '{"base_req":{"from":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z", "chain_id":"coinexdex"}, "address":"cosmos1wdzsu25mwlen0twt7vlar76af84mnsjtul4d9z" }'`
