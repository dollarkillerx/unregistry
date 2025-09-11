# Unregistry

A lightweight private file and Docker image storage system with client-server architecture.

## Features

- **File Storage**: Upload, download, list, and delete files
- **Docker Image Storage**: Push, pull, list, and delete Docker images (as tar.gz)
- **Token-based Authentication**: Secure API access with Bearer tokens
- **Simple CLI Client**: Easy-to-use command-line interface
- **Docker Support**: Ready-to-deploy with Docker and docker-compose

## Quick Start

### cli

mac: 
```
curl -L https://github.com/dollarkillerx/unregistry/releases/download/v0.0.2/unrg-mac > /usr/local/bin/unrg
chmod +x /usr/local/bin/unrg
```
linux:
```
curl -L https://github.com/dollarkillerx/unregistry/releases/download/v0.0.2/unrg-linux > /usr/local/bin/unrg
chmod +x /usr/local/bin/unrg
```

setconfig
```
unrg config set-url http://127.0.0.0:8888
unrg config set-token xxxxxxxxxxx
```

### Server Setup

1. **Build Docker image**:
```bash
make docker
```

2. **Start with docker-compose**:
```bash
docker-compose up -d
```

3. **Or run directly**:
```bash
docker run -d \
  -p 8080:8080 \
  -e TOKEN=123456 \
  -v ./data:/data \
  unregistry:latest
```

### Client Setup

1. **Build Linux client**:
```bash
make client-linux
```

2. **Configure the client**:
```bash
# Set authentication token
./bin/unrg-linux config set-token 123456

# Set server URL (if not localhost:8080)
./bin/unrg-linux config set-url http://your-server:8080
```

## Usage

### File Operations

```bash
# Upload a file
./bin/unrg-linux file push ./document.pdf

# List all files
./bin/unrg-linux file list

# Download a file
./bin/unrg-linux file pull document.pdf ./downloads/

# Delete a file
./bin/unrg-linux file delete document.pdf
```

### Docker Image Operations

```bash
# Push a Docker image
./bin/unrg-linux img push nginx:latest

# List all images
./bin/unrg-linux img list

# Pull a Docker image
./bin/unrg-linux img pull nginx:latest

# Delete an image
./bin/unrg-linux img delete nginx:latest
```

## Build Commands

```bash
# Build server
make server

# Build Docker image
make docker

# Build Linux client
make client-linux

# Clean build files
make clean
```

## Configuration

### Server Environment Variables

- `TOKEN`: Authentication token (default: "123456")
- `LISTEN_ADDR`: Server listen address (default: "0.0.0.0:8080")
- `DATA_PATH`: Data storage path (default: "/data")

### Client Configuration

The client stores configuration in `~/.unrg/config.json`:

```json
{
  "token": "123456",
  "base_url": "http://localhost:8080"
}
```

## API Endpoints

### File Operations
- `POST /api/file/upload` - Upload a file
- `GET /api/file/download/:filename` - Download a file
- `GET /api/file/list` - List all files
- `DELETE /api/file/:filename` - Delete a file

### Image Operations
- `POST /api/img/upload` - Upload an image (tar.gz)
- `GET /api/img/download/:name` - Download an image
- `GET /api/img/list` - List all images
- `DELETE /api/img/:name` - Delete an image

### Health Check
- `GET /health` - Health check (no auth required)

All API endpoints (except `/health`) require authentication via the `Authorization: Bearer <token>` header.

## Deployment

1. **Build and start server**:
```bash
make docker
docker-compose up -d
```

2. **Build client for Linux**:
```bash
make client-linux
```

3. **Configure client**:
```bash
./bin/unrg-linux config set-token 123456
./bin/unrg-linux config set-url http://your-server:8080
```

## Architecture

```
┌─────────────────┐    HTTP/API    ┌─────────────────┐
│                 │◄──────────────►│                 │
│   Client (unrg) │                │     Server      │
│                 │                │                 │
└─────────────────┘                └─────────────────┘
                                            │
                                            ▼
                                   ┌─────────────────┐
                                   │   File System   │
                                   │                 │
                                   │  /data/files/   │
                                   │  /data/images/  │
                                   └─────────────────┘
```

- **Server**: REST API with token authentication
- **Client**: CLI tool for file and image operations  
- **Storage**: File system storage (/data/files, /data/images)
- **Authentication**: Bearer token authentication