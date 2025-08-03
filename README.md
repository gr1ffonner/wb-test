# WB Test Application

A Go-based microservice application with PostgreSQL, Redis, and NATS messaging.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)


# Server
SERVER_PORT=8080

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# NATS
NATS_URL=nats://nats-streaming:4222

# Logging
LOG_LEVEL=info
```

### Running with Docker Compose

Start all services:
```bash
make up
```

Start only infrastructure (without app):
```bash
make up-dev
```

Stop all services:
```bash
make down
```

### Local Development

1. Install dependencies:
```bash
go mod tidy
```

2. Start only infrastructure (without app):
```bash
make up-dev
```

3. Run the application:
```bash
make run
```

## 📋 Services & Ports

| Service | Port | Description |
|---------|------|-------------|
| **App** | `8080` | Main Go application |
| **PostgreSQL** | `5432` | Database |
| **Redis** | `6379` | Cache |
| **NATS Streaming** | `4222` | Message broker |
| **NATS Monitoring** | `8222` | NATS monitoring UI |



## 📝 Makefile Commands

```bash
make up          # Start all services with test profile
make down        # Stop all services
make build       # Build Go application
make run         # Run application locally
make test        # Run tests
make clean       # Clean build artifacts
```
