FROM golang:1.16.10-alpine3.13 as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev bash

COPY  . $GOPATH/src
WORKDIR $GOPATH/src

# Install minimum necessary dependencies, build binary
RUN apk add --no-cache $PACKAGES && make all

FROM alpine:3.13

COPY --from=builder /go/src/iobscan-ibc-explorer-backend /usr/local/bin/

CMD ["iobscan-ibc-explorer-backend", "start"]