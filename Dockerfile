# Start from the official Go image
FROM golang:1.19-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage with a minimal alpine image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
#RUN apk install bash
#RUN apt install bash
RUN apk add --no-cache bash

# Copy the static files
COPY static ./static

# Expose the port the app runs on
EXPOSE 8088

# Run the binary
CMD ["./main"]
