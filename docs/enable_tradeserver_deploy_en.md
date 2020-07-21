# Enable the trade-server function in cetd



Follow the [document[2,  3.2]](docs/AtlantisHardForkGuide.en.md) to configure the common settings of the `cetd` node, and then configure the special settings to `trade-server function` as described below.


#### 1. Customize node configuration

In the configuration file of `cetd`, which is located in `${RUN_DIR}/.cetd/config/app.toml`, need to be appended the following content:

```bash
feature-toggle = true 
subscribe-modules = "comment,authx,bankx,market,bancorlite" 
brokers = [
   "prune:/path/to/dex_data"            # Directory for storing specified node data </br>
]
```
In the configuration file of `cetd`, which is located in `${RUN_DIR}/.cetd/config/app.toml`

#### 2. Modify the configuration of trade-server 

##### 2.1 Set the push data directory of cetd

Copy the file [trade-server.toml.default](https://github.com/coinexchain/dex/blob/master/trade-server.toml.default) to `${RUN_DIR}/.cetd/config/trade-server.toml`; 

Then modify the `dir`  field of the configuration(`{RUN_DIR}/.cetd/config/trade-server.toml`) , which will be consistent with the path of` prune` mode configuration under `brokers` in` cetd` configuration file `${RUN_DIR}/.cetd/config/app.toml`.

##### 2.2 Set the data directory of the trade-server function

First download the history data of cetd:
 `wget https://github.com/coinexchain/artifacts/raw/master/coinexdex-v0.2/history_data.tar.gz`.

    The history data stores the order information on the `coinexdex` chain, which is organized together with the current order information on the` coinexdex2` chain to form the data such as market depths and tickers.

Then unzip the history data: `tar -zxvf history_data.tar.gz`

Finally, in the configuration file `${RUN_DIR}/.cetd/config/trade-server.toml`, modify the `data-dir` field, which is consistent with the directory where history data is stored.

`${RUN_DIR}/.cetd/config/trade-server.toml` [The meaning of fields in this file](https://github.com/coinexchain/trade-server/blob/master/docs/trade-server-deploy.md#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E)

#### 3. Start cetd node

Follow the remaining description of the [document[3.3,  3.4]](docs/AtlantisHardForkGuide.en.md) and start the node

