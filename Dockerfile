# Dockerfile

# Stage 1: Build the Go application
FROM golang:1.23.4 AS builder

# Set the working directory inside the builder
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code from cmd/opggvisualizer
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./grafana ./grafana

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o opggvisualizer ./cmd/opggvisualizer

# Stage 2: Create the final image
FROM golang:1.23.4

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/opggvisualizer .

RUN chmod +x /app/opggvisualizer

# Create a volume for the database
VOLUME /opggvisualizer_data

# Expose the API server port
EXPOSE 8080

# Set the entrypoint to the built executable
ENTRYPOINT ["./opggvisualizer"]

# Default command to run the API server
CMD ["server", "start"]
