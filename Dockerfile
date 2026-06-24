# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

# Install build dependencies for Alpine
RUN apk add --no-cache \
    protoc \
    protobuf-dev \
    git \
    build-base

WORKDIR /app

# Copy the rest of the source code
COPY . .

# Set path for protoc-gen-go
ENV PATH="$PATH:$(go env GOPATH)/bin"

# Install protobuf generators (if proxy allows)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Generate protobuf code
RUN protoc --go_out=. --go-grpc_out=. proto/rivet.proto

# Build the application using vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /nexus-super-node ./cmd/nexus-super-node

# Stage 2: Create the final image with runtime dependencies
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    nodejs \
    npm \
    git \
    bash

# Copy the static binary from the builder stage
COPY --from=builder /nexus-super-node /nexus-super-node
COPY migrations ./migrations
COPY config/config.yaml ./config.yaml

# Expose the application port
EXPOSE 3000

# Set the entrypoint for the container
CMD [ "/nexus-super-node" ]
