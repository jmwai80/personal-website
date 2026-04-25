# Kubernetes (K8s)

> We are NOT using Kubernetes yet. This is reference for when you outgrow Docker Compose.

## What Is Kubernetes?

Kubernetes is a **container orchestration system**. It manages running containers across multiple machines, handles failures, scales automatically, and does rolling deployments.

Think of Docker Compose as managing containers on *one* machine. Kubernetes manages containers across *many* machines as if they were one.

## When to Use Kubernetes (vs Docker Compose)

| Situation | Use |
|-----------|-----|
| 1-2 servers, <100k req/day | Docker Compose ✓ (our current setup) |
| Need zero-downtime rolling deploys | Kubernetes |
| Need auto-scaling (traffic spikes) | Kubernetes |
| Multiple teams, many services | Kubernetes |
| Self-healing (auto-restart failed nodes) | Kubernetes |
| 10+ services | Kubernetes |

**For a personal site + apartment site: Docker Compose is the right choice.** Don't add K8s complexity until you need it.

## Core Concepts

### Cluster
A set of machines (nodes) that Kubernetes manages as one unit.

```
Kubernetes Cluster
├── Control Plane (master) — makes decisions
└── Worker Nodes — run your containers
    ├── Node 1 (VM)
    ├── Node 2 (VM)
    └── Node 3 (VM)
```

### Pod
The smallest deployable unit in Kubernetes. A pod is one or more containers that share network and storage.

```yaml
# Pod definition
apiVersion: v1
kind: Pod
metadata:
  name: personal-site
spec:
  containers:
    - name: app
      image: ghcr.io/user/personal-site:latest
      ports:
        - containerPort: 8080
```

- Usually 1 container per pod
- Pods are **ephemeral** — they can be killed and replaced at any time
- Each pod gets its own IP within the cluster

### Deployment
Manages a set of identical pods. Ensures N replicas are always running. Handles rolling updates.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: personal-site
spec:
  replicas: 3                    # always keep 3 pods running
  selector:
    matchLabels:
      app: personal-site
  template:
    spec:
      containers:
        - name: app
          image: ghcr.io/user/personal-site:latest
```

When you update the image, Kubernetes does a rolling update: starts new pods, waits for them to be healthy, then kills old ones. **Zero downtime.**

### Service
Pods come and go with random IPs. A **Service** gives them a stable IP/DNS name and load balances across all matching pods.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: personal-site
spec:
  selector:
    app: personal-site     # routes to pods with this label
  ports:
    - port: 8080
```

### Ingress
Like Caddy/Nginx, but for Kubernetes. Routes external traffic to the right Service by domain/path.

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
spec:
  rules:
    - host: yourdomain.com
      http:
        paths:
          - path: /
            backend:
              service:
                name: personal-site
                port:
                  number: 8080
```

### Namespace
Virtual cluster within a cluster. Used to separate environments (dev/staging/prod) or teams.

### ConfigMap & Secret
- **ConfigMap**: Non-sensitive config (env vars, config files)
- **Secret**: Sensitive data (passwords, API keys) — base64 encoded

### Persistent Volume (PV) / Persistent Volume Claim (PVC)
Like Docker volumes, but cluster-aware. Data persists even if the pod moves to a different node.

## Docker Compose → Kubernetes Mapping

| Docker Compose | Kubernetes |
|----------------|------------|
| `services` | Deployment + Service |
| `ports` | Service + Ingress |
| `volumes` | PersistentVolume |
| `environment` | ConfigMap / Secret |
| `depends_on` | Init containers / readiness probes |
| `docker compose up` | `kubectl apply -f` |
| `docker compose logs` | `kubectl logs` |
| `docker ps` | `kubectl get pods` |

## Managed Kubernetes Services

You don't run K8s yourself — you use a managed service:

| Provider | Service | Notes |
|----------|---------|-------|
| Google | GKE | Best K8s experience, most features |
| AWS | EKS | Most common in enterprise |
| DigitalOcean | DOKS | Cheapest managed K8s (~$12/month) |
| Oracle | OKE | Has free control plane |

## When You'd Migrate From Docker Compose

Signs you need Kubernetes:
- You have >3 services and deployments are getting complex
- You need horizontal auto-scaling (handle traffic spikes)
- You need zero-downtime rolling deploys (Docker Compose restarts cause brief downtime)
- You have multiple VMs and need to coordinate containers across them
- You want self-healing (auto-replace crashed nodes)

Migration path: keep the same Docker images, write K8s manifests or use Helm charts.

## Related Concepts
- [[containers-docker]] — Kubernetes orchestrates Docker containers
- [[cicd]] — CI/CD pipelines deploy to K8s via `kubectl apply`
- [[virtual-machines]] — K8s worker nodes are VMs
