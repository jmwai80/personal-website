# Progress Log

## Status Overview
| Phase | Description | Status |
|-------|-------------|--------|
| 1 | Local prerequisites | 🔄 In Progress |
| 2 | Oracle Cloud VM | ⬜ Not started |
| 3 | GitHub repo + secrets | ⬜ Not started |
| 4 | Go app skeleton + Dockerfiles | ⬜ Not started |
| 5 | Docker Compose + Caddy on VM | ⬜ Not started |
| 6 | GitHub Actions CI/CD | ⬜ Not started |
| 7 | GoDaddy DNS wired up | ⬜ Not started |
| 8 | Personal site live | ⬜ Not started |
| 9 | Apartments site live | ⬜ Not started |

---

## Phase 1 — Local Prerequisites

### Completed
- [x] Homebrew installed
- [x] Go 1.26.2 installed
- [x] GitHub CLI (gh) installed
- [x] Git already present (2.50.1)

### Remaining
- [ ] Docker Desktop installed (needs sudo — run: `brew install --cask docker`)
- [ ] Docker Desktop opened and started (first-run setup)
- [ ] SSH key generated (`~/.ssh/id_ed25519`)
- [ ] GitHub CLI authenticated (`gh auth login`)

---

## Phase 2 — Oracle Cloud VM
- [ ] Account created at cloud.oracle.com
- [ ] ARM VM provisioned (VM.Standard.A1.Flex, 4 OCPU, 24GB RAM, Ubuntu 22.04)
- [ ] VM public IP noted: ___________
- [ ] SSH access confirmed
- [ ] Docker installed on VM
- [ ] Firewall ports opened (80, 443)

---

## Phase 3 — GitHub
- [ ] Repo created: `personal-website`
- [ ] Local repo pushed to GitHub
- [ ] GitHub Secrets added (VM_HOST, VM_USER, VM_SSH_KEY, GHCR credentials)

---

## Decisions Made
| Date | Decision | Reason |
|------|----------|--------|
| 2026-04-24 | Docker over raw binary deploy | Portability, easy rollback, consistent environments |
| 2026-04-24 | Caddy over Nginx | Auto TLS, simpler config |
| 2026-04-24 | Oracle Free Tier as first VM | Free forever, 4 ARM cores + 24GB RAM |
| 2026-04-24 | Hetzner as fallback | In case Oracle account signup fails |
| 2026-04-24 | SQLite for apartments DB (initial) | Zero setup, sufficient for 48 units |
| 2026-04-24 | ghcr.io for container registry | Free, integrated with GitHub Actions |
