# builder image
FROM golang:1.18-alpine3.16 as builder

RUN apk --no-cache add alpine-sdk
WORKDIR /go/src/github.com/solarlabsteam/coingecko-exporter
COPY . .
RUN go build -o /usr/local/bin/coingecko-exporter -v -ldflags "-w -s"
RUN /usr/local/bin/coingecko-exporter --help

# final image
FROM alpine:3.16.2

RUN apk --no-cache add ca-certificates dumb-init
COPY --from=builder /usr/local/bin/coingecko-exporter /usr/local/bin/coingecko-exporter

USER 65534
ENTRYPOINT ["dumb-init", "--", "/usr/local/bin/coingecko-exporter"]
