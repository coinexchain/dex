## install statik

```bash
go get -d github.com/rakyll/statik
go install github.com/rakyll/statik
```

## start rest-server

```bash
./cetcli rest-server --chain-id=coinexdex  --laddr=tcp://localhost:1317  --node tcp://localhost:26657 --trust-node=false

```

access path ：
[http://localhost:1317/swagger/](http://localhost:1317/swagger/)