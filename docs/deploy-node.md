# 安装节点

1. 安装rocksdb依赖
2. 下载新版本节点、genesis.json 文件；
3. 初始化新数据目录
4. 启动节点


## 安装Rocksdb依赖

> git clone https://github.com/facebook/rocksdb.git && cd rocksdb </br> 
> git checkout v6.6.4 </br>
> mkdir build && cd build && cmake .. && make -j2 </br>
> sudo make install </br>


## 下载新版本

> wget https://github.com/coinexchain/dex/releases/download/v0.2.17/linux_x86_64.tar.gz </br>
> tar -zxvf linux_x86_64.tar.gz </br>
> </br>
> </br>
> export RUN_DIR=~~/path/to/node~~ </br>
> export PUBLIC_IP=~~123.23.42.22~~ </br>
> export CHAIN_SEEDS=903458cf236851ccf8604689c3f391c528191f47@47.75.37.80:26656,9be765dffed72adcd27ebb37c79bf8ac501f43e8@47.52.155.115:26656,cd79d6c2b3b6b561c91b61b8e3a706249b532ca4@47.56.215.151:26656,cf34ba278ce69be1240f1dabad9b57ffecae206a@47.75.60.29:26656,c70feea1a4f8ea2fd55c366fdcb7ca4d53f1c775@18.144.85.87:26656,94b718f31dedf4afee4c04d768343166625cf961@47.52.70.137:26656,2cbef50b8c996745b9c8a0059fe32a1fbfef8b46@47.52.129.186:26656,17ec2dcfd7c72fabcb7c7cfe2d71006fc39c85c9@18.180.56.174:26656 </br>
> mkdir ${RUN_DIR}  && mv linux_x86_64/* ${RUN_DIR}  && cd ${RUN_DIR} </br>


## 创建新的数据目录

1. `${RUN_DIR}/cetd init moniker  --home=${RUN_DIR}/.cetd`
2. 拷贝下载的genesis.json 到数据目录: cp ${RUN_DIR}/genesis.json ${RUN_DIR}/.cetd/config </br>
3. 如果是验证者节点：

    *   拷贝原节点数据目录的`priv_validator_key.json` 至新数据目录，该文件所在的位置：`${RUN_DIR}/.cetd/config`
4. 配置节点seeds

   *    `ansible localhost -m ini_file -a "path=${RUN_DIR}/.cetd/config/config.toml section=p2p option=seeds value='\"${CHAIN_SEEDS}\"' backup=true"`
   *   [ansible安装文档](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#installing-ansible-on-ubuntu)
5. 设置节点的对外IP

	*	`ansible localhost -m ini_file -a "path=${RUN_DIR}/.cetd/config/config.toml section=p2p option=external_address value='\"tcp://${PUBLIC_IP}:26656\"' backup=true"`

    
## 启动新节点    

接下来可按照您的习惯，选择使用systemctl或supervisor等工具来配置新链cetd的自动运行。这里的配置方式各不相同，本文不再一一介绍。

或者，先用下述简单方式启动节点，看是否连接到主链，是否出块

>  ${RUN_DIR}/cetd start --home=${RUN_DIR}/.cetd --minimum-gas-prices=20.0cet   <br/>



