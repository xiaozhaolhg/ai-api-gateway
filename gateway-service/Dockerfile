# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# Install dependencies
RUN apk add --no-cache git

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# Copy go mod files and api directory
COPY ./gateway-service/go.mod ./
COPY ./api/ ./api/

# Copy source code
COPY ./gateway-service/ ./

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway-service ./cmd/server

# Runtime stage
FROM alpine:latest

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gateway-service .

# Copy config
COPY --from=builder /app/configs ./configs

# Run as non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser

EXPOSE 8080

CMD ["./gateway-service"]
