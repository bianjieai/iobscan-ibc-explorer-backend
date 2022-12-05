FROM golang:1.18-alpine3.15 as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev bash

COPY  . $GOPATH/src
WORKDIR $GOPATH/src

# Install minimum necessary dependencies, build binary
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk add --no-cache $PACKAGES && make all

FROM alpine:3.15

COPY --from=builder /go/src/iobscan-ibc-openapi /usr/local/bin/

CMD ["iobscan-ibc-openapi", "start"]