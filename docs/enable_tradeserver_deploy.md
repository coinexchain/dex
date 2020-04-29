# 启动cetd 中的trade-server功能

### 安装Rocksdb依赖

> git clone https://github.com/facebook/rocksdb.git && cd rocksdb </br> 
> git checkout v6.6.4 </br>
> mkdir build && cd build && cmake .. && make -j2 </br>
> sudo make install </br>

### 初始化节点配置

下载可执行程序、初始化数据目录

> wget https://github.com/coinexchain/dex/releases/download/v0.2.17/linux_x86_64.tar.gz </br>
> tar -zxvf linux_x86_64.tar.gz </br>
> </br>
> </br>
> </br>
> export RUN_DIR=~~/path/to/node~~ </br>
> mkdir ${RUN_DIR}  && mv linux_x86_64/* ${RUN_DIR}  && cd ${RUN_DIR} </br>
> ${RUN_DIR}/cetd init moniker --home=${RUN_DIR}/.cetd </br>
> cp ${RUN_DIR}/genesis.json ${RUN_DIR}/.cetd/config </br>


### 修改节点配置

修改`cetd`配置文件 `${RUN_DIR}/.cetd/config/app.toml`; 在文件末尾添加下述内容


> feature-toggle = true </br>
>  subscribe-modules = "comment,authx,bankx,market,bancorlite" </br>
>  brokers = [ </br>
>      "prune:/path/to/dex_data"                # 指定节点吐数据的存储目录 </br>
>  ] </br>
>

修改`cetd`配置文件 `${RUN_DIR}/.cetd/config/config.toml`; 修改文件中的`seeds`字段，替换为如下内容

`seeds = "903458cf236851ccf8604689c3f391c528191f47@47.75.37.80:26656,9be765dffed72adcd27ebb37c79bf8ac501f43e8@47.52.155.115:26656,cd79d6c2b3b6b561c91b61b8e3a706249b532ca4@47.56.215.151:26656,cf34ba278ce69be1240f1dabad9b57ffecae206a@47.75.60.29:26656,c70feea1a4f8ea2fd55c366fdcb7ca4d53f1c775@18.144.85.87:26656,94b718f31dedf4afee4c04d768343166625cf961@47.52.70.137:26656,2cbef50b8c996745b9c8a0059fe32a1fbfef8b46@47.52.129.186:26656,17ec2dcfd7c72fabcb7c7cfe2d71006fc39c85c9@18.180.56.174:26656"`

### 修改trade-server 配置

##### 设置cetd推送数据的目录

拷贝项目目录下的[trade-server.toml.default](https://github.com/coinexchain/dex/blob/master/trade-server.toml.default)` 至 `${RUN_DIR}/.cetd/config/trade-server.toml`; 

修改该配置文件中`dir`的配置与`cetd`配置文件`app.toml`中`brokers`下`prune`模式配置的路径一致;


##### 设置trade-server功能的数据目录
 
下载cetd的历史数据: `wget https://github.com/coinexchain/artifacts/raw/master/coinexdex-v0.2/history_data.tar.gz`.

    *   历史数据中存储了cetd `coinexdex`链上的订单信息，与当前`coinexdex2`链上的订单信息一起组织，形成用户所需要的订单深度、ticker等数据；
 
解压该历史数据：`tar -zxvf history_data.tar.gz`

修改`${RUN_DIR}/.cetd/config/trade-server.toml`配置文件中`data-dir`字段，指向存储历史数据的目录；


`${RUN_DIR}/.cetd/config/trade-server.toml` 配置文件中[各字段含义](https://github.com/coinexchain/trade-server/blob/master/docs/trade-server-deploy.md#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E)

### 启动节点

接下来可按照您的习惯，选择使用systemctl或supervisor等工具来配置新链cetd的自动运行。这里的配置方式各不相同，本文不再一一介绍。

或者，先用下述简单方式启动节点，看是否连接到主链，是否出块

>  ${RUN_DIR}/cetd start --home=${RUN_DIR}/.cetd --minimum-gas-prices=20.0cet   <br/>

