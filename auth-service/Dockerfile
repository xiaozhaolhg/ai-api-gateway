# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# Copy go mod files
COPY ./auth-service/go.mod ./
COPY ./auth-service/go.sum ./

# Copy api directory
COPY ./api/ ./api/

# Copy source code
COPY ./auth-service/ ./auth-service

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /bin/auth-service ./auth-service/cmd/server

# Runtime stage
FROM alpine:latest

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /bin/auth-service .

# Copy config
COPY --from=builder /app/auth-service/configs ./configs

# Create data directory and startup script before switching user
RUN mkdir -p /root/data
RUN printf '#!/bin/sh\necho "Starting auth service..."\n./auth-service\n' > /root/start.sh
RUN chmod +x /root/start.sh

# Run as non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /root
USER appuser

EXPOSE 50051

CMD ["/root/start.sh"]
