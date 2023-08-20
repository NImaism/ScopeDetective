# Use a specific version of Go, such as Go 1.20 (if available) or Go 1.17
FROM golang:1.20 AS builder

WORKDIR /app

# Copy the entire application source code into the container
COPY . .

# Download dependencies using go modules
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app

# Start a new stage with a base image that includes SSL certificates
FROM debian:latest

# Install ca-certificates to provide SSL certificate support
RUN apt-get update && apt-get install -y ca-certificates

# Copy the compiled binary from the builder stage into the new image
COPY --from=builder /go/bin/app /go/bin/app

# Expose the port that your application listens on
EXPOSE 9090

# Set the command to run your application
CMD ["/go/bin/app"]
