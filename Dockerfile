FROM golang:1.25

WORKDIR /app

# Copy everything
COPY . .

# Build server
RUN go build -o server ./cmd/server

# Create data directory
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variables
ENV LISTEN_ADDR=0.0.0.0:8080
ENV DATA_PATH=/data
ENV TOKEN=123456

# Run the server
CMD ["./server"]