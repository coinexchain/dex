# 启动cetd 中的trade-server功能

### 安装Rocksdb依赖

> git clone https://github.com/facebook/rocksdb.git && cd rocksdb </br> 
> git checkout v6.6.4 </br>
> mkdir build && cd build && cmake .. && make -j2 </br>
> sudo make install </br>

### 节点初始化配置

参照[节点的通用步骤步骤](https://github.com/coinexchain/artifacts/blob/master/coinexdex-v0.2/Validator+%E5%93%A8%E5%85%B5%E8%8A%82%E7%82%B9-%E9%83%A8%E7%BD%B2%E6%96%B9%E6%A1%88.md#%E8%8A%82%E7%82%B9%E7%9A%84%E9%80%9A%E7%94%A8%E9%83%A8%E7%BD%B2%E6%AD%A5%E9%AA%A4)初始化节点配置；

### 修改节点配置

修改`cetd`自身配置：配置文件路径 `RUN_DIR/.cetd/config/app.toml`


> feature-toggle = true </br>
>  subscribe-modules = "comment,authx,bankx,market,bancorlite" </br>
>  brokers = [ </br>
>      "prune:/path/to/dex_data"                # 指定节点吐数据的存储目录 </br>
>  ] </br>
>

### 修改trade-server 配置

拷贝项目目录下的`trade-server.toml.default` 至 `RUN_DIR/.cetd/config/trade-server.toml`; 

修改 `dir`的配置与`cetd`配置文件中`brokers`下`prune`模式配置的路径一致;


`trade-server.toml` 配置文件中[各字段含义](https://github.com/coinexchain/trade-server/blob/master/docs/trade-server-deploy.md#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E)

### 放置历史数据
 
 从[地址](todo)下载cetd的历史数据，历史数据中存储了cetd `coinexdex`链上的订单信息，与当前`coinexdex2`链上的订单信息一起组织，
 形成用户所需要的订单深度、ticker等数据；
 
 修改`trade-server.toml`配置文件中`data-dir`字段，指向存储 `cetd`历史数据的目录；

### 启动节点

接下来可按照您的习惯，选择使用systemctl或supervisor等工具来配置新链cetd的自动运行。这里的配置方式各不相同，本文不再一一介绍。
