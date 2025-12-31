# renovate-exporter

**renovate-exporter** is a Prometheus metrics exporter that processes [Renovate](https://github.com/renovatebot/renovate) logs in real-time, extracting valuable insights about your dependency management workflow.

## Features

- **Real-time log processing** - Watches directories for new log files and processes them as they appear
- **Prometheus metrics** - Exports metrics in Prometheus format for monitoring and alerting
- **Docker-ready** - Designed to run as a sidecar container alongside Renovate
- **Rich metrics** - Track repository status, job duration, dependency age (libyears), and PR statistics

## Installation

### Docker (Recommended)

Pull the latest image from Docker Hub:

```bash
docker pull martinvigg/renovate-exporter:latest
```

### Binary Releases

Download pre-built binaries from the [Releases](https://github.com/martin-viggiano/renovate-exporter/releases) page.

### Build from Source

```bash
git clone https://github.com/martin-viggiano/renovate-exporter.git
cd renovate-exporter
go build -o renovate-exporter .
```

## Quick Start

### Docker Compose with Renovate CE

Deploy renovate-exporter as a sidecar alongside Renovate:

```yaml
services:
  renovate:
    image: ghcr.io/mend/renovate-ce:latest
    container_name: renovate-ce
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data/logs:/logs

  renovate-exporter:
    image: martinvigg/renovate-exporter:latest
    container_name: renovate-exporter
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./data/logs:/logs
    command: ["--path", "/logs"]
```

Start the services:

```bash
docker-compose up -d
```

Metrics will be available at `http://localhost:9090/metrics`.

### Standalone Docker

```bash
docker run -d \
  --name renovate-exporter \
  -p 9090:9090 \
  -v /path/to/renovate/logs:/logs \
  martinvigg/renovate-exporter:latest \
  --path /logs
```

### CLI Usage

```bash
renovate-exporter --path /var/log/renovate --address :9090
```

## Configuration

renovate-exporter can be configured via **CLI flags** or **environment variables**.

### Configuration Options

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--path`, `-p` | `RENOVATE_EXPORTER_PATH` | *(required)* | Directory to watch for Renovate log files |
| `--address` | `RENOVATE_EXPORTER_ADDRESS` | `:9090` | Metrics server address and port |

### Configuration Priority

When multiple configuration sources are provided, the following precedence order applies (highest to lowest):

1. **CLI flags** - Explicit command-line arguments
2. **Environment variables** - Variables prefixed with `RENOVATE_EXPORTER_`
3. **Defaults** - Built-in default values

### Examples

#### Using CLI Flags

```bash
renovate-exporter --path /logs --address :8080
```

#### Using Environment Variables

```bash
export RENOVATE_EXPORTER_PATH=/logs
export RENOVATE_EXPORTER_ADDRESS=:8080
renovate-exporter
```

#### Docker with Environment Variables

```bash
docker run -d \
  -e RENOVATE_EXPORTER_PATH=/logs \
  -e RENOVATE_EXPORTER_ADDRESS=:9090 \
  -v /path/to/logs:/logs \
  -p 9090:9090 \
  martinvigg/renovate-exporter:latest
```

## Metrics

renovate-exporter exposes the following Prometheus metrics:

### `renovate_repositories`

**Type:** Gauge  
**Description:** Tracks the status of repositories being processed by Renovate.

**Labels:**

- `repository` - Repository name (e.g., `owner/repo`)
- `status` - Repository processing status

**Example PromQL:**

```promql
# Count repositories by status
sum by (status) (renovate_repositories)

# Check if a specific repository completed successfully
renovate_repositories{repository="myorg/myrepo", status="done"}
```

---

### `renovate_repositories_duration_seconds`

**Type:** Gauge  
**Description:** Duration (in seconds) of the Renovate job for each repository.

**Labels:**

- `repository` - Repository name

**Example PromQL:**

```promql
# Average job duration across all repositories
avg(renovate_repositories_duration_seconds)

# Repositories with jobs taking longer than 5 minutes
renovate_repositories_duration_seconds > 300

# Top 5 slowest repositories
topk(5, renovate_repositories_duration_seconds)
```

---

### `renovate_repositories_libyears`

**Type:** Gauge  
**Description:** [Libyear](https://libyear.com/) metric representing how outdated dependencies are for each repository and package manager.

**Labels:**

- `repository` - Repository name
- `manager` - Package manager (e.g., `npm`, `pip`, `docker`, `go`)

**Example PromQL:**

```promql
# Total libyears across all repositories
sum(renovate_repositories_libyears)

# Libyears by package manager
sum by (manager) (renovate_repositories_libyears)

# Repositories with high npm dependency age
renovate_repositories_libyears{manager="npm"} > 10
```

---

### `renovate_pull_requests_total`

**Type:** Gauge  
**Description:** Number of pull requests managed by Renovate, broken down by state.

**Labels:**

- `repository` - Repository name
- `state` - PR state (`open`, `merged`, `closed`)

**Example PromQL:**

```promql
# Total open PRs across all repositories
sum(renovate_pull_requests_total{state="open"})

# Repositories with the most open PRs
topk(10, renovate_pull_requests_total{state="open"})

# PR merge rate
sum(renovate_pull_requests_total{state="merged"}) / 
  (sum(renovate_pull_requests_total{state="merged"}) + 
   sum(renovate_pull_requests_total{state="closed"}))
```

---

### `renovate_dependencies_total`

**Type:** Gauge  
**Description:** Total number of dependencies per repository and package manager discovered by Renovate.

**Labels:**

- `repository` - Repository name
- `manager` - Package manager (e.g., `npm`, `gomod`, `docker`, `pip`)

**Example PromQL:**

```promql
# Total dependencies across all repositories
sum(renovate_dependencies_total)

# Dependencies by package manager
sum by (manager) (renovate_dependencies_total)

# Repositories with the most dependencies
topk(10, sum by (repository) (renovate_dependencies_total))

# Average dependencies per repository
avg(sum by (repository) (renovate_dependencies_total))
```

---

### `renovate_dependency_files_total`

**Type:** Gauge  
**Description:** Number of dependency files per repository and package manager discovered by Renovate.

**Labels:**

- `repository` - Repository name
- `manager` - Package manager (e.g., `npm`, `gomod`, `docker`, `pip`)

**Example PromQL:**

```promql
# Total dependency files across all repositories
sum(renovate_dependency_files_total)

# Dependency files by package manager
sum by (manager) (renovate_dependency_files_total)

# Repositories with multiple package.json files
renovate_dependency_files_total{manager="npm"} > 1
```

---

### `renovate_dependency_outdated_total`

**Type:** Gauge  
**Description:** Total number of outdated dependencies per repository and package manager that have available updates.

**Labels:**

- `repository` - Repository name
- `manager` - Package manager (e.g., `npm`, `gomod`, `docker`, `pip`)

**Example PromQL:**

```promql
# Total outdated dependencies across all repositories
sum(renovate_dependency_outdated_total)

# Outdated dependencies by package manager
sum by (manager) (renovate_dependency_outdated_total)

# Repositories with the most outdated dependencies
topk(10, sum by (repository) (renovate_dependency_outdated_total))

# Percentage of outdated dependencies per repository
(sum by (repository) (renovate_dependency_outdated_total) / 
 sum by (repository) (renovate_dependencies_total)) * 100
```

## Usage Examples

### Prometheus Scrape Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'renovate-exporter'
    static_configs:
      - targets: ['renovate-exporter:9090']
```

### Alerting Examples

Example Prometheus alert rules:

```yaml
groups:
  - name: renovate
    rules:
      - alert: RenovateJobTooSlow
        expr: renovate_repositories_duration_seconds > 600
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Renovate job taking too long"
          description: "Repository {{ $labels.repository }} took {{ $value }}s to process"

      - alert: HighDependencyAge
        expr: renovate_repositories_libyears > 20
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "Dependencies are significantly outdated"
          description: "Repository {{ $labels.repository }} has {{ $value }} libyears of dependency age for {{ $labels.manager }}"

      - alert: TooManyOpenPRs
        expr: renovate_pull_requests_total{state="open"} > 50
        for: 30m
        labels:
          severity: info
        annotations:
          summary: "High number of open Renovate PRs"
          description: "Repository {{ $labels.repository }} has {{ $value }} open PRs"
```

### Grafana Dashboard

A Grafana dashboard for visualizing these metrics will be available soon. Stay tuned!

## How It Works

renovate-exporter follows a simple pipeline architecture:

1. **File System Watcher** (`fswatch`) - Monitors the specified directory for new `.log` files
2. **Log Reader** (`reader`) - Tails new files and streams log lines
3. **Log Analyzer** (`analyzer`) - Parses JSON log entries and matches patterns
4. **Metrics Extraction** (`matchers`) - Extracts relevant data and updates Prometheus metrics
5. **HTTP Server** - Exposes metrics on `/metrics` endpoint for Prometheus scraping

### Log Format

renovate-exporter expects Renovate logs in JSON format with structured fields. The following log messages are currently parsed:

- `"Repository finished"` - Extracts repository status and duration
- `"Renovate repository PR statistics"` - Extracts PR counts
- `"Repository libYears"` - Extracts libyear metrics per package manager

## Development

### Prerequisites

- Go 1.25 or later
- Docker (optional, for container testing)

### Build

```bash
go build -o renovate-exporter .
```

### Run Tests

```bash
go test ./...
```

### Run with Race Detection

```bash
go test -race ./...
```

### Run Locally

```bash
go run main.go --path /path/to/logs
```

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 Martin Viggiano

## Acknowledgments

- [Renovate](https://github.com/renovatebot/renovate) - The awesome dependency update tool
- [Prometheus](https://prometheus.io/) - Monitoring and alerting toolkit
- All contributors who help improve this project

---

**Made with ❤️ for better dependency management observability**
