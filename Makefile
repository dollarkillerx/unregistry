.PHONY: server docker client-linux clean

# Build server
server:
	go build -o bin/server ./cmd/server

# Docker build
docker:
	docker build -t dollarkiller/unregistry:latest .

# Build client for Linux
client-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/unrg-linux ./cmd/client

# Clean
clean:
	rm -rf bin/