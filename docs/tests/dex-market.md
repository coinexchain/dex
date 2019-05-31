# dex-market

## 启动节点

`./cetd start`

## 发行token 

发行两种token: eth, cet

1. 发行 eth

`./cetcli tx asset issue-token --name="eth" --symbol="eth" --total-supply=2100000000000000 --mintable=false --burnable=true --addr-forbiddable=false --token-forbiddable=false --from bob --chain-id=coinexdex`

2. 发行 cet 

`./cetcli tx asset issue-token --name="cet" --symbol="cet" --total-supply=2100000000000000 --mintable=false --burnable=true --addr-forbiddable=false --token-forbiddable=false --from bob --chain-id=coinexdex`


## 查询token

1. 查询cet
`./cetcli query asset token cet --chain-id=coinexdex`

2. 查询eth
`./cetcli query asset token eth --chain-id=coinexdex`

## 创建交易对市场
创建 eth/cet 对交易市场

`./cetcli tx market createmarket --from bob --chain-id=coinexdex  --gas 20000 --stock=eth --money=cet --price-precision=8`
   
## 创建订单

`./cetcli tx market creategteoreder --symbol="eth/cet" --order-type=2 --price=520 --quantity=10000000 --side=1 --time-in-force=1000 --from bob --price-precision=9 --chain-id=coinexdex  `   

## 查询指定市场信息

`./cetcli query market marketinfo eth/cet --trust-node=true`

## 查询指定订单信息

`./cetcli query market orderinfo --symbol=eth/cet --orderid=cosmos16gvnhynu7veexyyaadk60k28cn5s9k7p7p5v9p-13 --trust-node=true`

