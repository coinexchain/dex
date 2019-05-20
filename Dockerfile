FROM golang:latest

COPY . /dex
WORKDIR /dex/

RUN go mod tidy
RUN go mod vendor
RUN go build ./...

ENTRYPOINT ["/sh/start.sh"]