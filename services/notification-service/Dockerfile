# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install protoc and plugins for gRPC code generation
RUN apk add --no-cache protobuf protobuf-dev git make

# Install Go protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy proto files and generate Go code
COPY proto/ ./proto/
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/common.proto proto/notification.proto

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates wget

COPY --from=builder /app/notification-service .

EXPOSE 8083 50053

HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8083/health || exit 1

CMD ["./notification-service"]
