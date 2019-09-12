FROM golang:1.12-alpine
LABEL maintainer="dev@coinex.org"

COPY networks/test/cetdnode/wrapper.sh    /usr/bin/
COPY networks/test/cetdnode/rest_start.sh /usr/bin/
COPY build/cetcli /usr/bin/
COPY build/cetd /usr/bin/

VOLUME [ /cetd ]
WORKDIR /cetd
EXPOSE 26656 26657 27000
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

RUN ["chmod", "+x", "/usr/bin/wrapper.sh"]