# 🏗 Stage 1: Build the application
FROM golang:1.23.2-alpine AS builder

# Install dependencies
RUN apk --no-cache add curl

WORKDIR /app

# Install Air for hot reloading
# RUN go install github.com/air-verse/air@latest

# Copy Go modules first (for better caching)
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Install golang-migrate inside the container
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# RUN curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz && \
#     mv migrate /usr/local/bin/migrate

# Build the Go binary
RUN go build -o main .

# 🚀 Stage 2: Run container with hot reloading
FROM golang:1.23.2-alpine
# FROM gcr.io/distroless/static:nonroot


WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
# COPY --from=builder /go/bin/air /usr/local/bin/air
COPY entrypoint.sh .
COPY . .

# Make entrypoint executable
RUN chmod +x entrypoint.sh

# Expose the application port
EXPOSE 8000

# Run the application
# CMD ["air"]
# CMD ["./main"]
ENTRYPOINT ["./entrypoint.sh"]
