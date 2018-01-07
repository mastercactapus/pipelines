FROM alpine

RUN apk add -U --no-cache gnupg curl jq
COPY sign-key.asc /root/
RUN gpg --import </root/sign-key.asc

COPY in check /opt/resource/
