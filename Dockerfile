# Start from the official Golang image
FROM golang:1.24

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Create app directory
WORKDIR /

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o server .

# Create uploads directory in container
RUN mkdir -p /uploads

# Set environment variable for upload dir (optional, can also be in .env)
ENV UPLOAD_DIR=/uploads

# Expose the port the app runs on
EXPOSE 3007

# Command to run the executable
CMD ["./server"]
