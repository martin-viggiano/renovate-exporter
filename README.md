# renovate-exporter

renovate-exporter is a utility to generate metrics based on the processing of Renovate logs.

## Usage

### Docker

renovate-exporter can be deployed as a sidecar container alongside a self-hosted `renovate-ce` container.

```compose
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
