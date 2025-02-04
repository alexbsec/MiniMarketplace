FROM ubuntu:22.04

# Install Go, PostgreSQL client and dependencies
RUN apt-get update && apt-get install -y wget curl git gcc postgresql-client && \
    wget https://go.dev/dl/go1.19.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.19.6.linux-amd64.tar.gz && \
    rm go1.19.6.linux-amd64.tar.gz

# Add Go to PATH var
ENV PATH="/usr/local/go/bin:$PATH"

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build
RUN go build -o bin/main ./src

# Expose the application port
EXPOSE 7676

# Run
CMD ["/app/main"]
