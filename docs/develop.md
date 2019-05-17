#### go dependencies
set `go module`:
```bash
export GO111MODULE=on
```
make vendored copy of dependencies

```bash
go mod vendor
```

Add missing and remove unused modules
```
 go mod tidy 
```

verify dependencies have expected content
```bash
go mod verify
```
get dependencies
```bash
go get -u
```

### GolangCI-Lint

```bash
go get  github.com/golangci/golangci-lint/cmd/golangci-lint
```