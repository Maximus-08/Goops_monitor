# goops-monitor

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Build](https://github.com/Maximus-08/goops-monitor/actions/workflows/ci.yml/badge.svg)](https://github.com/Maximus-08/goops-monitor/actions/workflows/ci.yml)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://hub.docker.com)

A lightweight, extensible monitoring and task execution framework written in Go.

## Overview

`goops-monitor` provides real-time health monitoring and automated task execution for distributed services. It's designed to be simple to deploy yet powerful enough to scale with your infrastructure.

## Architecture

```
goops-monitor/
├── monitor/          # Health check and metrics collection
│   ├── main.go       # Entry point for the monitor service
│   ├── config.go     # Configuration management
│   └── metrics.go    # Metrics tracking (up/down counts)
├── runner/           # Task execution engine
│   ├── runner.go     # Command execution logic
│   └── task.go       # Task lifecycle management
└── tutorials/        # Getting started guides
```

## Features

- **Real-time Health Monitoring** – Configurable intervals for endpoint checks
- **Task Runner** – Execute shell commands with status tracking
- **Metrics Collection** – Track uptime and downtime events
- **JSON Configuration** – Easy setup via config files
- **Concurrent Safe** – Thread-safe metrics with mutex locks

## Quick Start

### Prerequisites

- Go 1.22 or higher

### Installation

```bash
git clone https://github.com/Maximus-08/goops-monitor.git
cd goops-monitor
go mod tidy
```

### Running the Monitor

```bash
go run monitor/*.go
```

The monitor will start checking the configured target (default: `http://localhost:8080`) every 5 seconds.

### Using Makefile

```bash
make build    # Build binary
make run      # Build and run
make test     # Run tests
make docker   # Build Docker image
make clean    # Clean build artifacts
```

### Docker

Build and run with Docker:
```bash
docker build -t goops-monitor .
docker run -p 8080:8080 -p 8081:8081 goops-monitor
```

With custom config:
```bash
docker run -p 8080:8080 -p 8081:8081 -v $(pwd)/config.json:/app/config.json goops-monitor
```

### Docker Compose

For easier deployment:
```bash
docker-compose up -d
```

View logs:
```bash
docker-compose logs -f
```

Stop:
```bash
docker-compose down
```

## Configuration

Create a `config.json` file:

```json
{
  "interval": "10s",
  "target": "http://your-service:8080/health"
}
```

## Environment Variables

Override config.json with environment variables (useful for Docker):

| Variable | Description | Example |
|----------|-------------|---------|
| `GOOPS_INTERVAL` | Check interval | `30s` |
| `GOOPS_TARGETS` | Comma-separated URLs | `http://app1:8080,http://app2:8080` |
| `GOOPS_RETRIES` | Retry count | `3` |
| `GOOPS_WEBHOOK_URL` | Alert webhook | `https://hooks.slack.com/...` |
| `GOOPS_ALERT_COOLDOWN` | Alert cooldown | `5m` |
| `GOOPS_ON_FAILURE` | Remediation script | `./restart.sh` |
| `GOOPS_LOG_JSON` | Enable JSON logging | `true` |
| `GOOPS_CONFIG` | Config file path | `/app/config.json` |
| `GOOPS_API_PORT` | API server port | `:9090` |

Example:
```bash
GOOPS_TARGETS="http://myservice:8080" GOOPS_INTERVAL="30s" ./monitor_bin
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `/metrics` | Prometheus metrics |
| `/status` | JSON status with uptime percentage |
| `/ready` | Kubernetes readiness probe |
| `/live` | Kubernetes liveness probe |

## Usage Examples

### Using the Runner

```go
package main

import "goops-monitor/runner"

func main() {
    r := runner.New("echo", "Hello, World!")
    if err := r.Execute(); err != nil {
        panic(err)
    }
}
```

### Task Management

```go
task := runner.NewTask("task-001", "backup.sh")
task.MarkRunning()
// ... execute task ...
task.MarkCompleted()
```

