FROM frolvlad/alpine-glibc
LABEL maintainer="dev@coinex.org"

RUN apk --no-cache add curl tzdata

VOLUME [ /cetd ]
WORKDIR /cetd
EXPOSE 26656 26657 27000
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY networks/test/cetdnode/wrapper.sh networks/test/cetdnode/rest_start.sh build/cetcli build/cetd /usr/bin/
RUN ["chmod", "+x", "/usr/bin/wrapper.sh"]
