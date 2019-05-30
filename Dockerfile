FROM coinexchain/basic:latest

COPY . /dex
WORKDIR /dex

RUN echo '{ "allow_root": true }' > /root/.bowerrc
RUN go mod tidy
RUN go mod vendor
RUN go install github.com/rakyll/statik

RUN ./scripts/build.sh
