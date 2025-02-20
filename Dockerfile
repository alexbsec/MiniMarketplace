FROM ubuntu:22.04

# Install Go, PostgreSQL client, MongoDB client, and dependencies
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    make \
    git \
    gcc \
    gnupg \
    postgresql-client && \
    wget https://go.dev/dl/go1.22.7.linux-amd64.tar.gz && \
    curl -sSf https://atlasgo.sh | sh && \
    tar -C /usr/local -xzf go1.22.7.linux-amd64.tar.gz && \
    rm go1.22.7.linux-amd64.tar.gz

    
RUN curl -fsSL https://www.mongodb.org/static/pgp/server-8.0.asc | \
   gpg -o /usr/share/keyrings/mongodb-server-8.0.gpg \
   --dearmor

RUN echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-8.0.gpg ] https://repo.mongodb.org/apt/ubuntu noble/mongodb-org/8.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-8.0.list

RUN apt-get update

# Add Go to PATH
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
RUN make build 

# Expose the application port
EXPOSE 7676

# Run
CMD ["/app/bin/main"]

