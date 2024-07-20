FROM golang:1.19 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY main.go .
COPY local/config/config.go ./local/config/
COPY local/splitter/splitter.go ./local/splitter/

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o ecowitt-proxy main.go

# Start a runtime stage from scratch
FROM debian:bookworm-slim

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/ecowitt-proxy .

# Ensure the binary has execute permissions
RUN chmod +x ecowitt-proxy

EXPOSE 8123

# Command to run the executable
CMD ["./ecowitt-proxy"]