FROM alpine

RUN apk add -U --no-cache gnupg curl jq && \
    gpg --no-tty --recv-keys 94E3083A18042697

COPY in check /opt/resource/
