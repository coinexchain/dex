# Enable the trade-server function in cetd

### Install Rocksdb


> git clone https://github.com/facebook/rocksdb.git && cd rocksdb </br> 
> git checkout v6.6.4 </br>
> mkdir build && cd build && cmake .. && make -j2 </br>
> sudo make install </br>

### Initialize node configuration

Download the binary executable program and initialize the data directory


> wget https://github.com/coinexchain/dex/releases/download/v0.2.17/linux_x86_64.tar.gz </br>
> tar -zxvf linux_x86_64.tar.gz </br>
> </br>
> </br>
> </br>
> export RUN_DIR=~~/path/to/node~~ </br>
> mkdir ${RUN_DIR}  && mv linux_x86_64/* ${RUN_DIR}  && cd ${RUN_DIR} </br>
> ${RUN_DIR}/cetd init moniker --home=${RUN_DIR}/.cetd </br>
> cp ${RUN_DIR}/genesis.json ${RUN_DIR}/.cetd/config </br>


### Customize node configuration

In the configuration file of `cetd`, which is located in `${RUN_DIR}/.cetd/config/app.toml`, need to be appended the following content:

> feature-toggle = true </br>
>  subscribe-modules = "comment,authx,bankx,market,bancorlite" </br>
>  brokers = [ </br>
>      "prune:/path/to/dex_data"                # Directory for storing specified node data </br>
>  ] </br>
>

In the configuration file of `cetd`, which is located in `${RUN_DIR}/.cetd/config/app.toml`, need to be replaced the content of `seed` field:

`seeds = "903458cf236851ccf8604689c3f391c528191f47@47.75.37.80:26656,9be765dffed72adcd27ebb37c79bf8ac501f43e8@47.52.155.115:26656,cd79d6c2b3b6b561c91b61b8e3a706249b532ca4@47.56.215.151:26656,cf34ba278ce69be1240f1dabad9b57ffecae206a@47.75.60.29:26656,c70feea1a4f8ea2fd55c366fdcb7ca4d53f1c775@18.144.85.87:26656,94b718f31dedf4afee4c04d768343166625cf961@47.52.70.137:26656,2cbef50b8c996745b9c8a0059fe32a1fbfef8b46@47.52.129.186:26656,17ec2dcfd7c72fabcb7c7cfe2d71006fc39c85c9@18.180.56.174:26656"`

### Modify the configuration of trade-server 

##### Set the push data directory of cetd

Copy the file [trade-server.toml.default](https://github.com/coinexchain/dex/blob/master/trade-server.toml.default) to `${RUN_DIR}/.cetd/config/trade-server.toml`; 

Then modify the configuration of `dir`, which will be consistent with the path of` prune` mode configuration under `brokers` in` cetd` configuration file `app.toml`.

##### Set the data directory of the trade-server function
 
First download the history data of cetd:
 `wget https://github.com/coinexchain/artifacts/raw/master/coinexdex-v0.2/history_data.tar.gz`.

    The history data stores the order information on the `coinexdex` chain, which is organized together with the current order information on the` coinexdex2` chain to form the data such as market depths and tickers.
 
Then unzip the history data: `tar -zxvf history_data.tar.gz`

Finally, in the configuration file `${RUN_DIR}/.cetd/config/trade-server.toml`, modify the `data-dir` field, which is consistent with the directory where history data is stored.

`${RUN_DIR}/.cetd/config/trade-server.toml` [The meaning of fields in this file](https://github.com/coinexchain/trade-server/blob/master/docs/trade-server-deploy.md#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E)

### Start cetd node

Next, you can choose to use tools such as `systemctl` or `supervisor` to configure the automatic operation of the new `cetd` (coinexdex2) according to your habits. The configuration methods are various, and this document will not introduce them one by one.

Or, start the node in the following simple way to see if it can be connected to the main chain and whether the block can be generated.

>  ${RUN_DIR}/cetd start --home=${RUN_DIR}/.cetd --minimum-gas-prices=20.0cet   <br/>

