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
COPY ./billing-service/ ./billing-service

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /bin/billing-service ./billing-service/cmd/server

# Runtime stage
FROM alpine:latest

# Configure Alpine mirrors for China
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /bin/billing-service /app/billing-service

# Copy config
COPY --from=builder /app/billing-service/configs /app/configs

# Create data directory
RUN mkdir -p /data

EXPOSE 50054

CMD ["/app/billing-service"]
