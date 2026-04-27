# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# Copy go mod files and api directory
COPY ./go.mod ./
COPY ./api/ ./api/

# Copy source code
COPY ./monitor-service/ ./monitor-service

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /bin/monitor-service ./monitor-service/cmd/server

# Runtime stage
FROM alpine:latest

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /bin/monitor-service .

# Copy config
COPY --from=builder /app/monitor-service/configs ./configs

# Create data directory
RUN mkdir -p /root/data

# Run as non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser

EXPOSE 50055

CMD ["./monitor-service"]
