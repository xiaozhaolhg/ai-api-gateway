FROM golang:1.26-alpine AS builder

RUN apk add --no-cache go git

WORKDIR /build

COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn,direct go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o router-service ./cmd/server

FROM alpine:3.23

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 65532 nonroot

COPY --from=builder /build/router-service /router-service

USER nonroot:nonroot

ENTRYPOINT ["/router-service"]