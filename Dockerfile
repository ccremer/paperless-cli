FROM docker.io/library/alpine:3.24 as runtime

ENTRYPOINT ["paperless-cli"]

RUN \
    apk add --update --no-cache \
      bash \
      curl \
      ca-certificates \
      tzdata

COPY paperless-cli /usr/bin/
USER 65536:0
