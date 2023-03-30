FROM golang:1.18-alpine3.15 as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev bash

ARG GITUSER=bamboo
ARG GITPASS=FS_Q5LmxwExwK6hFN9Fs
ARG GOPRIVATE=gitlab.bianjie.ai
ARG GOPROXY=http://192.168.0.60:8081/repository/go-bianjie/,http://nexus.bianjie.ai/repository/golang-group,https://goproxy.cn,direct
ARG APKPROXY=http://mirrors.ustc.edu.cn/alpine

COPY  . $GOPATH/src
WORKDIR $GOPATH/src

# Install minimum necessary dependencies, build binary
RUN sed -i "s+http://dl-cdn.alpinelinux.org/alpine+${APKPROXY}+g" /etc/apk/repositories && \
    apk add --no-cache $PACKAGES && \
    git config --global url."https://${GITUSER}:${GITPASS}@gitlab.bianjie.ai".insteadOf "https://gitlab.bianjie.ai" && make all

FROM alpine:3.15

COPY --from=builder /go/src/iobscan-ibc-openapi /usr/local/bin/

CMD ["iobscan-ibc-openapi", "start"]