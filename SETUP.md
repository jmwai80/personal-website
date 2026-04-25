# Personal Website — Setup & Deployment Guide

> Living document. Updated as we progress through setup.

## Architecture Overview

```
Local (Go + Docker)
    │
    └── git push → GitHub
                      │
                      └── GitHub Actions (build + push Docker image)
                                │
                                └── Oracle Cloud VM (ARM Ubuntu)
                                          │
                                          ├── Caddy (reverse proxy + auto HTTPS)
                                          ├── personal-site container  :8080
                                          └── apartments-site container :8081
                                                    │
                                          GoDaddy DNS → VM public IP
```

## Sites
| Site | Domain | Port | Status |
|------|--------|------|--------|
| Personal | TBD | 8080 | [ ] not started |
| Apartments | TBD | 8081 | [ ] not started |

---

## Phase 1: Prerequisites Checklist

### Local Machine
- [ ] Go installed (`go version` — need 1.21+)
- [ ] Docker Desktop installed and running
- [ ] Git installed
- [ ] SSH key pair generated (`~/.ssh/id_ed25519`)

### Accounts
- [ ] Oracle Cloud account created (free tier)
- [ ] GitHub account (existing)
- [ ] GoDaddy domain (existing) — domain: ___________

---

## Phase 2: Oracle Cloud VM Setup

### 2.1 Create VM
- Region: choose closest to you (or `ap-sydney-1` / `eu-frankfurt-1` for low latency)
- Shape: **VM.Standard.A1.Flex** (ARM, Always Free)
  - OCPUs: 4
  - RAM: 24 GB
- Image: **Ubuntu 22.04**
- Generate or upload SSH key during creation
- Note VM public IP: ___________

### 2.2 Configure VM
```bash
# SSH into VM
ssh ubuntu@<VM_IP>

# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker ubuntu

# Install Docker Compose plugin
sudo apt install -y docker-compose-plugin

# Verify
docker --version
docker compose version
```

### 2.3 Open Firewall Ports (OCI Console)
In OCI Console → Networking → VCN → Security List, add ingress rules:
- Port 22 (SSH) — already open
- Port 80 (HTTP)
- Port 443 (HTTPS)

---

## Phase 3: GitHub Repository Setup

```bash
# On local machine
cd ~/Projects/personal-website
git init
git remote add origin git@github.com:<username>/personal-website.git
```

### GitHub Secrets to add (Settings → Secrets → Actions)
| Secret | Value |
|--------|-------|
| `VM_HOST` | Oracle VM public IP |
| `VM_USER` | `ubuntu` |
| `VM_SSH_KEY` | Contents of `~/.ssh/id_ed25519` (private key) |
| `DOCKER_USERNAME` | GitHub username or Docker Hub username |
| `DOCKER_PASSWORD` | GitHub PAT or Docker Hub token |

---

## Phase 4: Project Structure

```
personal-website/
├── SETUP.md                   # this file
├── docker-compose.yml         # production compose (on VM)
├── Caddyfile                  # reverse proxy config
├── personal-site/
│   ├── Dockerfile
│   ├── main.go
│   └── go.mod
└── .github/
    └── workflows/
        └── deploy.yml         # CI/CD pipeline
```

---

## Phase 5: Docker & App Setup

### personal-site/Dockerfile
Multi-stage build — small final image.
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Caddyfile (reverse proxy)
```
yourdomain.com {
    reverse_proxy personal-site:8080
}

apartments.yourdomain.com {
    reverse_proxy apartments-site:8081
}
```

### docker-compose.yml (on VM)
```yaml
services:
  caddy:
    image: caddy:2-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
    depends_on:
      - personal-site

  personal-site:
    image: ghcr.io/<username>/personal-site:latest
    restart: unless-stopped
    expose:
      - "8080"

volumes:
  caddy_data:
```

---

## Phase 6: CI/CD (GitHub Actions)

On every push to `main`:
1. Build Docker image
2. Push to GitHub Container Registry (ghcr.io) — free
3. SSH into VM
4. Pull new image + `docker compose up -d`

---

## Phase 7: DNS (GoDaddy)

In GoDaddy DNS settings:
| Type | Name | Value | TTL |
|------|------|-------|-----|
| A | @ | `<VM_IP>` | 600 |
| A | www | `<VM_IP>` | 600 |
| A | apartments | `<VM_IP>` | 600 |

Caddy handles TLS automatically via Let's Encrypt once DNS propagates.

---

## Progress Log

| Date | Step | Notes |
|------|------|-------|
| 2026-04-24 | Started planning | Architecture decided: Docker + Caddy + Oracle Free Tier + GitHub Actions |

---

## Useful Commands

```bash
# SSH into VM
ssh ubuntu@<VM_IP>

# View running containers on VM
docker ps

# View logs
docker compose logs -f personal-site

# Deploy manually (if not using CI/CD)
docker compose pull && docker compose up -d

# Caddy reload config without downtime
docker compose exec caddy caddy reload --config /etc/caddy/Caddyfile
```
