# goops-monitor

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)

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
go run monitor/main.go
```

The monitor will start checking the configured target (default: `http://localhost:8080`) every 5 seconds.

## Configuration

Create a `config.json` file:

```json
{
  "interval": "10s",
  "target": "http://your-service:8080/health"
}
```

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

