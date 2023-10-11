FROM golang:1.21-alpine as builder

RUN apk add make build-base

COPY . /builder
RUN cd /builder && make build

# The runtime image
FROM alpine:3

WORKDIR /work

RUN apk add --no-cache bash curl

COPY --from=builder /builder/xinfin-monitor /usr/bin

ENTRYPOINT ["bash","xinfin-monitor"]