# Build stage
FROM golang:1.23-alpine AS builder

# Install required packages
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    make \
    sqlite-dev \
    flex \
    bison \
    linux-headers \
    autoconf \
    automake \
    libtool

# Set up Go environment for cross-compilation
ENV GOARCH=arm64
ENV GOOS=linux
ENV CGO_ENABLED=1
ENV CC=aarch64-linux-musl-gcc

# Install aarch64 cross-compiler
RUN wget https://musl.cc/aarch64-linux-musl-cross.tgz \
    && tar -xzf aarch64-linux-musl-cross.tgz \
    && rm aarch64-linux-musl-cross.tgz \
    && mv aarch64-linux-musl-cross /opt/cross

# Set up the PATH to include the cross-compiler
ENV PATH="/opt/cross/bin:${PATH}"

# Download and cross-compile libpcap
RUN wget http://www.tcpdump.org/release/libpcap-1.10.4.tar.gz \
    && tar -xzf libpcap-1.10.4.tar.gz \
    && cd libpcap-1.10.4 \
    && CC=aarch64-linux-musl-gcc \
       AC_CV_HEADER_LINUX_IF_PACKET_H=yes \
       ac_cv_linux_vers=2 \
       ./configure --host=aarch64-linux-musl --with-pcap=linux \
    && make \
    && make install DESTDIR=/opt/cross-root \
    && cd .. \
    && rm -rf libpcap-1.10.4*

# Set environment variables for linking
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/opt/cross-root/usr/local/include"
ENV CGO_LDFLAGS="-L/opt/cross-root/usr/local/lib"
ENV PKG_CONFIG_PATH="/opt/cross-root/usr/local/lib/pkgconfig"

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
    go build -ldflags="-linkmode external -extldflags '-static'" -o /app/backend

# Final stage
FROM scratch

# Copy the binary from the build stage
COPY --from=builder /app/backend /backend

# Set the entrypoint
ENTRYPOINT ["/backend"]