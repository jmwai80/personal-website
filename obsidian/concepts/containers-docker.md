# Containers & Docker

## The Problem Containers Solve

> "It works on my machine."

You build an app on macOS with Go 1.23 and a specific library version. You deploy to a VM running Ubuntu with Go 1.19. It breaks. Why? The environment is different.

**Containers** package your app *together with* its entire environment: the runtime, libraries, config, everything. The container runs identically anywhere Docker is installed.

## VM vs Container

```
VM:
┌─────────────────────────────────┐
│  App A    │  App B    │  App C  │
│  Libs A   │  Libs B   │  Libs C │
│  OS A     │  OS B     │  OS C   │ ← 3 full operating systems
├─────────────────────────────────┤
│           Hypervisor            │
├─────────────────────────────────┤
│         Physical Server         │
└─────────────────────────────────┘

Container:
┌─────────────────────────────────┐
│  App A    │  App B    │  App C  │
│  Libs A   │  Libs B   │  Libs C │
├─────────────────────────────────┤
│         Container Runtime       │
│         (shared OS kernel)      │ ← 1 OS kernel, shared
├─────────────────────────────────┤
│         Physical Server         │
└─────────────────────────────────┘
```

| | VM | Container |
|--|----|----|
| Startup time | 30s–2min | <1 second |
| Size | GBs (full OS) | MBs |
| Isolation | Full (separate kernel) | Process-level (shared kernel) |
| Overhead | High | Very low |
| Use case | Full server, different OS | App packaging, microservices |

## Key Concepts

### Image
A **Docker image** is a read-only template. Think of it like a class in OOP — it defines what the container will look like, but isn't running yet.

Built from a `Dockerfile`. Stored in a registry (e.g. ghcr.io, Docker Hub).

### Container
A **container** is a running instance of an image. Like an object instantiated from a class. You can run many containers from the same image.

### Dockerfile
Instructions for building an image.

```dockerfile
# Start from official Go image
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download        # cache dependencies as a layer
COPY . .
RUN CGO_ENABLED=0 go build -o server .

# Final image — tiny, no Go compiler included
FROM alpine:3.19
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

This is a **multi-stage build** — the builder stage compiles, the final stage only contains the binary. Result: ~10MB image instead of ~500MB.

### Registry
Where images are stored and pulled from.
- **ghcr.io** (GitHub Container Registry) — free, what we use
- **Docker Hub** — the public default, free with limits
- **ECR** (AWS), **GCR** (Google) — cloud-specific

### Docker Compose
A tool to define and run **multiple containers together** using a YAML file.

```yaml
services:
  caddy:
    image: caddy:2-alpine
    ports: ["80:80", "443:443"]

  personal-site:
    image: ghcr.io/user/personal-site:latest
    expose: ["8080"]

  apartments-site:
    image: ghcr.io/user/apartments-site:latest
    expose: ["8081"]
```

Run all of them: `docker compose up -d`
Stop all: `docker compose down`
View logs: `docker compose logs -f`

### Volumes
Persistent storage for containers. Containers are ephemeral — when they restart, any data written inside is gone. Volumes persist data outside the container lifecycle.

```yaml
volumes:
  - caddy_data:/data        # named volume
  - ./config:/app/config    # bind mount (maps host dir to container)
```

## Essential Commands

```bash
# Build an image
docker build -t myapp:latest .

# Run a container
docker run -p 8080:8080 myapp:latest

# List running containers
docker ps

# View logs
docker logs <container-id>
docker logs -f <container-id>   # follow (live)

# Shell into a running container
docker exec -it <container-id> sh

# Stop a container
docker stop <container-id>

# Pull an image
docker pull ghcr.io/user/myapp:latest

# Remove unused images (cleanup)
docker image prune
```

## How Our Deploy Works

```
GitHub Actions builds image
    → pushes to ghcr.io/user/personal-site:latest
    → SSHs into VM
    → runs: docker compose pull && docker compose up -d
    → Compose pulls the new image, restarts only that container
    → zero downtime (old container serves while new one starts)
```

## Related Concepts
- [[virtual-machines]] — What runs Docker on the server
- [[kubernetes]] — What you use when Docker Compose isn't enough
- [[cicd]] — How Docker images get built and deployed automatically
