.PHONY: server docker client-linux clean

# Build server
server:
	go build -o bin/server ./cmd/server

# Docker build
docker:
	docker build -t dollarkiller/unregistry:latest .

client-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/unrg-linux ./cmd/client
client-mac:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/unrg-mac ./cmd/client

# Clean
clean:
	rm -rf bin/