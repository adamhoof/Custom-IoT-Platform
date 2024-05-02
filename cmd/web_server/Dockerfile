# Use the official Golang image as a parent image.
# This image includes the Go toolchain.
FROM golang:1.22.2 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace.
# Assuming your go.mod and go.sum are at the project root.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
# Adjust the source paths according to your project's structure
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the Go app as a static binary.
# CGO_ENABLED=0 is used to build a statically-linked binary to run in a minimal base image.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webapp ./cmd/web_server/main.go

# Use a minimal alpine image to run the final application to minimize the container size.
FROM alpine:latest

# Add ca-certificates in case you need HTTPS
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/webapp .

# Expose port 8080 to the outside world
EXPOSE 4444

# Command to run the executable
CMD ["./webapp"]
