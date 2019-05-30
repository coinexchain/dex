# 本地rest-server中访问swagger-ui的方法：

- 添加新增的API： cmd/cetcli/swagger-ui/swagger.yaml
- 使用 ./script/build.sh编译
- statik会将swagger-yaml相关文件编译打包成FS，并加入可执行文件
- 使用以下命令启动rest-server:
> ./cetcli rest-server --chain-id=coinexdex  --laddr=tcp://localhost:1317  --node tcp://localhost:26657 --trust-node=false

- 本地访问路径：http://localhost:1317/swagger-ui/ 

> 注：如果更新yaml，执行完上面的流程后，浏览器没有反映最新的信息，可能是浏览器缓存的问题



安装statik：

```bash
go get github.com/rakyll/statik
```

> 注：如果想要在任意目录下直接执行statik命令，需要把Go安装目录/bin/添加到path下





see also:

- https://github.com/rakyll/statik



