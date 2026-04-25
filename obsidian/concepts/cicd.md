# CI/CD — Continuous Integration & Continuous Deployment

## What Is CI/CD?

**CI (Continuous Integration)**: Every code push automatically runs tests, linting, builds — catching problems immediately.

**CD (Continuous Deployment)**: Every code push that passes CI automatically deploys to production.

Together: push code → tests run → if green → live in production. No manual steps.

## Why It Matters

Without CI/CD:
1. Write code locally
2. Manually build Docker image
3. Manually push image to registry
4. SSH into VM
5. Manually pull image
6. Manually restart containers
7. Manually check it's working

With CI/CD: `git push` — done in ~2 minutes automatically.

## GitHub Actions

What we use. YAML workflows defined in `.github/workflows/`.

### Key Concepts

**Workflow**: A YAML file defining when and what to run.

**Trigger (on)**: What starts the workflow. E.g. push to `main`, pull request, manual trigger.

**Job**: A unit of work that runs on a machine (called a **runner**). Jobs can run in parallel or sequentially.

**Step**: A single command or action within a job.

**Action**: A reusable unit of automation. E.g. `actions/checkout@v4` checks out your code.

**Runner**: The machine GitHub provides to run your jobs. GitHub offers free Linux/macOS/Windows runners.

**Secrets**: Encrypted values stored in GitHub Settings. Referenced as `${{ secrets.MY_SECRET }}`.

### Our Workflow

```yaml
name: Deploy

on:
  push:
    branches: [main]       # trigger on push to main

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      # 1. Get the code
      - uses: actions/checkout@v4

      # 2. Log into GitHub Container Registry
      - name: Login to ghcr.io
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}  # auto-provided, no setup needed

      # 3. Build Docker image for ARM64 (our Oracle VM is ARM)
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: ./personal-site
          platforms: linux/arm64
          push: true
          tags: ghcr.io/${{ github.repository_owner }}/personal-site:latest

      # 4. SSH into VM and pull new image
      - name: Deploy to VM
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.VM_HOST }}
          username: ${{ secrets.VM_USER }}
          key: ${{ secrets.VM_SSH_KEY }}
          script: |
            cd ~/app
            docker compose pull personal-site
            docker compose up -d personal-site
```

### Secrets We Need

| Secret | What it is |
|--------|-----------|
| `VM_HOST` | Oracle VM public IP |
| `VM_USER` | `ubuntu` |
| `VM_SSH_KEY` | Private SSH key (the one matching the public key on the VM) |
| `GITHUB_TOKEN` | Auto-provided by GitHub — no setup needed |

### GitHub Container Registry (ghcr.io)

Free Docker image registry built into GitHub.

- Images are tied to your GitHub account/org
- `ghcr.io/<username>/<image-name>:<tag>`
- Private by default (need to set package to public or add auth)
- Free for public repos; private repos have storage limits

## The Full Deploy Flow

```
git push origin main
        │
        ▼
GitHub Actions runner starts (Ubuntu VM, GitHub's infrastructure)
        │
        ├─ Checkout code
        ├─ Authenticate to ghcr.io
        ├─ docker buildx build --platform linux/arm64  ← cross-compile for ARM
        ├─ docker push ghcr.io/user/personal-site:latest
        │
        ├─ SSH into Oracle VM
        │   ├─ docker compose pull personal-site  ← download new image
        │   └─ docker compose up -d personal-site ← restart with new image
        │
        ▼
New version live (~2 min total)
```

## Cross-Platform Build (QEMU/buildx)

Our VM is **ARM64** but GitHub runners are **AMD64** (x86). We need to cross-compile the Docker image.

`docker buildx` with QEMU emulation handles this automatically. The `platforms: linux/arm64` line in the workflow handles it.

Alternatively, we can cross-compile Go directly in the Dockerfile:
```dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o server .
```

## Rollback

Since images are tagged `:latest`, to rollback:
1. Find the previous image digest in GitHub Packages
2. SSH into VM: `docker compose pull personal-site@sha256:<previous-digest>`
3. `docker compose up -d`

Better practice: tag with git commit SHA too:
```
ghcr.io/user/personal-site:latest
ghcr.io/user/personal-site:abc1234   ← commit SHA, permanent
```

## Related Concepts
- [[containers-docker]] — What gets built and deployed
- [[virtual-machines]] — Where it deploys to
- [[kubernetes]] — What replaces Docker Compose when you need more (rolling deploys, auto-scaling)
