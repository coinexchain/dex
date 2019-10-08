FROM frolvlad/alpine-glibc
LABEL maintainer="dev@coinex.org"

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file tzdata

VOLUME [ /cetd ]
WORKDIR /cetd
EXPOSE 26656 26657 27000
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh
COPY rest_start.sh /usr/bin/rest_start.sh
COPY cetd /usr/bin/cetd
COPY cetcli /usr/bin/cetcli

RUN ["chmod", "+x", "/usr/bin/wrapper.sh"]
