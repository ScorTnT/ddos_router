# Build stage
FROM golang:1.22-alpine AS builder

# Install required packages
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    make \
    sqlite-dev

# Set up Go environment for cross-compilation
ENV GOARCH=arm64
ENV GOOS=linux
ENV CGO_ENABLED=1
ENV CC=aarch64-linux-musl-gcc

# Install aarch64 cross-compiler for SQLite
RUN wget https://musl.cc/aarch64-linux-musl-cross.tgz \
    && tar -xzf aarch64-linux-musl-cross.tgz \
    && rm aarch64-linux-musl-cross.tgz \
    && mv aarch64-linux-musl-cross /opt/cross

# Set up the PATH to include the cross-compiler
ENV PATH="/opt/cross/bin:${PATH}"

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-musl-gcc \
    go build -ldflags="-linkmode external -extldflags '-static'" -o /app/router_api

# Final stage
FROM scratch

# Copy the binary from the build stage
COPY --from=builder /app/router_api /router_api

# Set the entrypoint
ENTRYPOINT ["/router_api"]