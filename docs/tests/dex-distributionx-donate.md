# DonateToCommunityPool

新增DonateToCommunityPool的交易，以支持普通用户对CommunityPool的捐赠行为。

## CLI Example

1.  初始化链

   ```bash
   ./cetd init node0 --chain-id=coinexdex
   ./cetd add-genesis-token --name="CoinEx Chain Native Token" \
   	--symbol="cet" \
   	--owner=coinex1628t2zxa9antj3qtkg7xj2m4t68uljqvyjqrup \
   	--total-supply=10000000000000000 \
   	--mintable=false \
   	--burnable=true \
   	--addr-forbiddable=false \
   	--token-forbiddable=false \
   	--total-burn=0 \
   	--total-mint=0 \
   	--is-forbidden=false 
   ./cetcli keys add bob
   ./cetd add-genesis-account $(./cetcli keys show bob -a) 10000000000000000cet
   ./cetd gentx --name bob
   ./cetd collect-gentxs
   ./cetd start
   ```

   

2. 查询ommunityPool的余额：

   ```bash
   ./cetcli query distr community-pool --trust-node
   ```

   返回结果：

   ```bash
   0.000000000000000000cet
   ```

3. 发起向CommunityPool 捐赠的交易：

   ```bash
   ./cetcli tx donate 100000000cet --from bob --fees 4000000cet --chain-id=coinexdex --gas 40000
   ```

4. 查询CommunityPool的余额：

   ```bash
   ./cetcli query distr community-pool --trust-node
   ```

   返回结果：

   ```bash
   100080000.000000000000000000cet
   ```

   可以看到，除了捐赠的100000000cet，交易手续费4000000cet的2%也进入了CommunityPool。

## REST Example

1. 启动Rest-Sever:

   ```bash
   ./cetcli rest-server --chain-id=coinexdex  --laddr=tcp://localhost:1317  --node tcp://localhost:26657 --trust-node=false
   ```

   

2. 发起rest请求：

   ```
   curl -X POST "http://localhost:1317/distribution/coinex1pvcp883r3wjv9u79d5ja2ka6dr6sc7ahtmqzza/donates" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"base_req\": { \"from\": \"coinex1pvcp883r3wjv9u79d5ja2ka6dr6sc7ahtmqzza\", \"memo\": \"from bob to ddd\", \"chain_id\": \"coinexdex\", \"account_number\": \"0\", \"sequence\": \"5\", \"gas\": \"200000\", \"gas_adjustment\": \"1.2\", \"fees\": [ { \"denom\": \"cet\", \"amount\": \"200000000\" } ], \"simulate\": false }, \"amount\": [ { \"denom\": \"cet\", \"amount\": \"1300000000\" } ]}" > unsignedSendTx.json
   ```

   

3. 对生成的unsignedSendTx.json进行签名：

   ```bash
   ./cetcli tx sign --chain-id=coinexdex --from=$(cetcli keys show bob -a)  unsignedSendTx.json > signedSendTx.json
   ```

   

4. 广播交易：

   ```bash
   ./cetcli tx broadcast signedSendTx.json
   ```

   

5. 查询CommunityPool的余额：

   ```bash
   ./cetcli query distr community-pool --trust-node
   ```

   返回结果：

   ```bash
   1404080000.000000000000000000cet
   ```

   可以看到，除了捐赠的1300000000，手续费200000000的2%也进入了CommunityPool。



