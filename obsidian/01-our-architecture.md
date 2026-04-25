# Our Architecture

## The Full Flow

```
Developer (you)
    │
    │  git push origin main
    ▼
GitHub Repository
    │
    │  triggers
    ▼
GitHub Actions (CI/CD pipeline)
    │  1. Checkout code
    │  2. Build Docker image (linux/arm64)
    │  3. Push image to ghcr.io (GitHub Container Registry)
    │  4. SSH into Oracle VM
    │  5. docker compose pull + up -d
    ▼
Oracle Cloud VM (ARM Ubuntu 22.04)
    │
    │  Running containers:
    ├── Caddy container (ports 80, 443)
    │       │  routes by domain name
    │       ├── yourdomain.com → personal-site:8080
    │       └── apartments.yourdomain.com → apartments-site:8081
    │
    ├── personal-site container (:8080)
    └── apartments-site container (:8081)
                                    │
                              GoDaddy DNS
                              A record → VM IP
                              Caddy handles TLS (Let's Encrypt)
```

## Why Each Piece Exists

| Component | Role | Alternative |
|-----------|------|-------------|
| Docker | Packages the Go app + dependencies into a portable unit | Run binary directly on VM |
| Docker Compose | Manages multiple containers as one unit | Manual `docker run` commands |
| Caddy | Terminates HTTPS, routes traffic to the right container | Nginx, Traefik |
| GitHub Actions | Automates build + deploy on every push | Jenkins, manual SSH deploy |
| ghcr.io | Stores Docker images between build and deploy | Docker Hub, self-hosted registry |
| Oracle Free Tier | The actual server running everything | Hetzner, DigitalOcean, AWS |

## Key Design Decisions

### Why Docker instead of running the binary directly?
- **Consistency**: The binary built on GitHub Actions runs identically on the VM
- **Portability**: Move to any VM by just installing Docker — no environment setup
- **Isolation**: Each app is sandboxed, can't interfere with each other
- **Rollback**: Pull a previous image tag to revert

### Why Caddy instead of Nginx?
- Auto-renews TLS certs via Let's Encrypt with zero config
- Single `Caddyfile` is human-readable vs Nginx's verbose config
- Handles HTTP → HTTPS redirect automatically

### Why ARM (Oracle A1 Flex)?
- It's the Always Free shape with the most resources (4 cores, 24GB)
- Go compiles natively to ARM64 — no performance penalty
- We explicitly target `GOARCH=arm64` in the Dockerfile

## Resource Usage Estimate
| Resource | Available | Expected Use |
|----------|-----------|-------------|
| CPU | 4 ARM cores | ~2-5% idle |
| RAM | 24 GB | ~200MB for both Go apps + Caddy |
| Disk | 200 GB | ~2GB for OS + images + data |
| Network | 10 TB/month | Well within limits for personal site |

## File Structure on VM
```
/home/ubuntu/
└── app/
    ├── docker-compose.yml
    ├── Caddyfile
    └── .env               # secrets (not in git)
```

## File Structure in Repo
```
personal-website/
├── SETUP.md
├── obsidian/              # this knowledge base
├── docker-compose.yml     # mirrors what's on VM
├── Caddyfile
├── personal-site/
│   ├── Dockerfile
│   ├── main.go
│   └── go.mod
├── apartments-site/
│   ├── Dockerfile
│   ├── main.go
│   └── go.mod
└── .github/
    └── workflows/
        └── deploy.yml
```
