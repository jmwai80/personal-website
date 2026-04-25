# Engineering Knowledge Base — Personal Website Project

> A living document. Updated as we build and learn.

## What We're Building
A production-grade deployment pipeline for 2 Go web apps, deployed on a free Oracle Cloud VM, served via Docker containers, with automatic deploys from GitHub.

- [[01-our-architecture]] — The full picture of our specific setup
- [[02-progress-log]] — What's been done, what's next

## Core Concepts (read in order)

### Infrastructure
1. [[concepts/virtual-machines]] — What is a VM, why it exists, when to use it
2. [[concepts/containers-docker]] — Containers vs VMs, Docker, Docker Compose
3. [[concepts/reverse-proxy]] — What Caddy/Nginx do and why you need one
4. [[concepts/dns-tls]] — How domains, DNS, and HTTPS work end to end

### Deployment
5. [[concepts/cicd]] — CI/CD pipelines, GitHub Actions
6. [[concepts/linux-server]] — Ubuntu server basics, systemd, file system

### Scaling (future)
7. [[concepts/kubernetes]] — Pods, services, orchestration — when you outgrow Docker Compose
8. [[concepts/cloud-providers]] — Oracle, AWS, GCP, Hetzner — how to compare

## Tech Stack in This Project
| Layer | Tool | Why |
|-------|------|-----|
| Language | Go | Fast, compiled binary, low memory use |
| Containerization | Docker | Consistent builds, easy deploy |
| Orchestration (local) | Docker Compose | Multi-container setup without Kubernetes complexity |
| Reverse Proxy | Caddy | Auto HTTPS, minimal config |
| CI/CD | GitHub Actions | Free, integrates with GitHub |
| Cloud VM | Oracle Cloud Free Tier | 4 ARM cores + 24GB RAM, free forever |
| Container Registry | GitHub Container Registry (ghcr.io) | Free, integrated with GitHub |
| DNS | GoDaddy | Existing domain |
