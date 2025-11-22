# Multi-stage build for optimal image size and security
FROM rust:1.83-slim-bullseye AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    liblmdb-dev \
    && rm -rf /var/lib/apt/lists/*

# Create app user for security
RUN useradd -m -u 1001 appuser

# Set working directory
WORKDIR /app

# Copy manifests first for better layer caching
COPY Cargo.toml Cargo.lock ./

# Create a dummy main.rs to build dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs

# Build dependencies (this layer will be cached unless Cargo files change)
RUN cargo build --release && rm -rf src target/release/deps/immortal*

# Copy source code
COPY src/ ./src/

# Build the application
RUN cargo build --release

# Runtime stage - use minimal base image
FROM debian:bullseye-slim AS runtime

# Install runtime dependencies only
RUN apt-get update && apt-get install -y \
    liblmdb0 \
    libssl1.1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Create app user (same UID as builder for consistency)
RUN useradd -m -u 1001 appuser

# Create working directory and data directory
WORKDIR /app
RUN mkdir -p /app/immo/database && chown -R appuser:appuser /app

# Copy the built binary from builder stage
COPY --from=builder /app/target/release/immortal /usr/local/bin/immortal

# Copy configuration file
COPY config.toml /app/config.toml

# Ensure binary is executable and owned by appuser
RUN chmod +x /usr/local/bin/immortal && chown appuser:appuser /usr/local/bin/immortal /app/config.toml

# Switch to non-root user
USER appuser

# Expose the default port (configurable via config.toml)
EXPOSE 7777

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:7777/ || exit 1

# Use exec form for proper signal handling
ENTRYPOINT ["/usr/local/bin/immortal"]

# Add labels for better maintainability
LABEL maintainer="dezh-tech" \
      version="0.1.0" \
      description="Immortal Nostr Relay" \
      org.opencontainers.image.source="https://github.com/dezh-tech/immortal"