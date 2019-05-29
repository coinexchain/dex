## Docker build basic

```bash
 docker build  -t coinexchain/basic:tag -f Dockerfile.basic .
```
Push docker image to docker hub
```bash
docker login
docker push coinexchain/basic:tag
```