FROM golang:1.21-alpine as builder

RUN apk add make build-base

COPY . /builder
RUN cd /builder && go build

# The runtime image
FROM alpine:3

WORKDIR /work

RUN apk add --no-cache bash curl

COPY .env /work/
COPY *.json /work/
COPY --from=builder /builder/XDC-blockchain-monitor /usr/bin
