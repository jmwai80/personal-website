# Cloud Providers

## The Big Picture

Cloud providers rent you computing resources (VMs, storage, databases, networking) instead of buying physical hardware.

## Providers We Evaluated

### Oracle Cloud (what we chose)
- **Always Free ARM**: 4 cores + 24GB RAM — genuinely free forever
- **Catch**: Account signup can fail; idle instances may be reclaimed
- **Best for**: Personal projects where cost matters most

### Hetzner (backup plan)
- **CAX11** (ARM): 2 vCPU, 4GB RAM — €3.29/month
- **CX22** (AMD): 2 vCPU, 4GB RAM — €3.79/month
- **Catch**: European company, datacenters in EU + US
- **Best for**: Cheapest reliable paid option; excellent reputation in developer community

### DigitalOcean
- **Droplet**: $4/month basic
- **Best for**: Beginners — best docs, smoothest UX, lots of tutorials
- **Catch**: More expensive than Hetzner for same specs

### AWS (Amazon)
- **EC2**: Many sizes, complex pricing
- **Free tier**: 750 hours/month t2.micro (1GB RAM) for 12 months (then paid)
- **Best for**: Enterprise, teams, services that need AWS-specific integrations
- **Catch**: Complex console, easy to accidentally incur costs

### Google Cloud (GCP)
- **Free tier**: e2-micro VM always free (limited regions)
- **Best for**: ML workloads, teams already in Google ecosystem

### Fly.io
- **Different model**: Deploys Docker containers directly, manages the VM for you
- **Free tier**: 3 shared VMs
- **Best for**: When you want to skip VM management entirely

## Key Concepts

### Region
Physical location of the datacenter. Choose closest to your users for lower latency.
- Oracle: `eu-frankfurt-1`, `ap-sydney-1`, `us-ashburn-1`
- Hetzner: Nuremberg, Helsinki, Ashburn (VA), Hillsboro (OR)

### OCPU vs vCPU
- **vCPU** (most providers): 1 virtual CPU core
- **OCPU** (Oracle): 1 OCPU = 2 vCPUs (Oracle's definition)
- Oracle's "4 OCPU" = 8 vCPUs in most other providers' terms

### ARM vs AMD64
- **AMD64 (x86_64)**: Traditional servers, most common
- **ARM64 (aarch64)**: Originally for mobile (like your phone). Now in servers — AWS Graviton, Oracle A1, Apple M1/M2
- **ARM advantage**: Better performance per watt, often cheaper
- **ARM consideration**: Some software doesn't support ARM (getting rare)
- Go supports both natively — no issue for us

### Egress Costs
Most providers charge for **outbound** data (data leaving their network). Inbound is usually free.
- Oracle: 10TB/month free egress (very generous)
- Hetzner: 20TB/month included
- AWS: $0.09/GB after first 1GB (can get expensive at scale)

## Migration Path

Since we're Docker-based, migrating between providers is:
1. Spin up new VM
2. Install Docker
3. Copy `docker-compose.yml` and `Caddyfile`
4. Update GitHub Secrets (`VM_HOST`)
5. Run deploy workflow
6. Update DNS A record to new IP
7. Done — ~30 minutes

This is why Docker portability matters.

## Related Concepts
- [[virtual-machines]] — What you're renting from these providers
- [[containers-docker]] — What makes moving between providers easy
