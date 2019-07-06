FROM coinexchain/go-build-env:latest AS build-env
LABEL maintainer="dev@coinex.org"

ADD . $GOPATH/src/github.com/coinexchain/dex

RUN cd $GOPATH/src/github.com/coinexchain/dex && \
    export GO111MODULE=on && \
    go mod tidy && \
    go mod vendor && \
    make statik-swagger && \
    make build-linux && \
    cp build/cetd /tmp/ && \
    cp build/cetcli /tmp/

FROM alpine:3.7

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

VOLUME [ /cetd ]
WORKDIR /cetd
EXPOSE 26656 26657 27000
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY networks/test/cetdnode/wrapper.sh    /usr/bin/
COPY networks/test/cetdnode/rest_start.sh /usr/bin/
COPY --from=build-env /tmp/cetd           /usr/bin/
COPY --from=build-env /tmp/cetcli         /usr/bin/

RUN ["chmod", "+x", "/usr/bin/wrapper.sh"]
