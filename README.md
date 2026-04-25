# Personal Website

Personal website built with Go, deployed via Docker on Oracle Cloud.

## Stack
- **Language**: Go
- **Container**: Docker
- **Reverse Proxy**: Caddy (auto HTTPS)
- **CI/CD**: GitHub Actions
- **Hosting**: Oracle Cloud ARM VM

## Development

```bash
go run ./...
```

## Deployment

Pushing to `main` triggers an automatic deploy via GitHub Actions.
