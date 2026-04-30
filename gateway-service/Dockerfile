# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# Install dependencies
RUN apk add --no-cache git

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# Copy go mod files
COPY ./go.mod ./

# Copy api directory
COPY ./api/ ./api/

# Copy pkg directory
COPY ./pkg/ ./pkg/

# Copy source code
COPY ./gateway-service/ ./gateway-service

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/gateway-service ./gateway-service/cmd/server

# Runtime stage
FROM alpine:latest

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /bin/gateway-service .

# Copy config
COPY --from=builder /app/gateway-service/configs ./configs

# Run as non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser

EXPOSE 8080

CMD ["./gateway-service"]
