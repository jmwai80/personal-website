# Virtual Machines (VMs)

## What is a VM?

A Virtual Machine is a software emulation of a physical computer. It runs an operating system and applications just like a real computer, but it exists entirely as software on top of a physical host machine.

```
Physical Server (Host)
├── Hypervisor (manages VMs)
├── VM 1: Ubuntu → your Go app
├── VM 2: Windows → someone else's app
└── VM 3: CentOS → another app
```

The **hypervisor** (e.g. VMware, KVM, Xen) divides the physical hardware and allocates slices to each VM.

## Why VMs Exist

Before VMs, if you wanted to run a web server you had two options:
1. Buy a dedicated physical server (expensive, underutilized)
2. Share a server with others (no isolation — one app crashes, all go down)

VMs solved this: cloud providers buy massive physical servers and slice them into many VMs. You rent a slice. You get isolation (your VM crashes, others are fine), and the provider maximizes hardware utilization.

## Key Properties of a VM

| Property | Description |
|----------|-------------|
| **Isolation** | Full OS isolation — each VM has its own kernel |
| **Dedicated resources** | CPU cores, RAM, disk are reserved for you |
| **Persistent** | Survives reboots; data persists to disk |
| **Slow to start** | Boot time: 30 seconds to 2 minutes |
| **Heavy** | Includes full OS (1–4GB just for the OS) |
| **SSH access** | You connect via SSH and control it like a physical computer |

## VM vs Bare Metal

| | VM | Bare Metal |
|--|----|----|
| Cost | Cheap (share hardware) | Expensive (dedicated) |
| Isolation | Strong (hypervisor) | Complete |
| Performance | ~5-10% overhead | Full |
| Use case | Most workloads | High-performance DB, gaming servers |

## Types of VMs in Cloud

- **Shared CPU (burstable)**: You share CPU cores with others. Cheap. Fine for low-traffic sites.
- **Dedicated CPU**: Cores reserved only for you. More consistent performance.
- **ARM-based**: Uses ARM chips (like your phone). Oracle A1 Flex is ARM. Go compiles natively to it.

## What We're Using

**Oracle VM.Standard.A1.Flex** — ARM-based, 4 OCPU, 24GB RAM. Always Free.

This is a **VPS (Virtual Private Server)** — marketing term for a VM you rent. Same thing.

## Key Commands Once on a VM

```bash
ssh ubuntu@<ip>          # connect
sudo apt update          # update package list
df -h                    # check disk space
free -h                  # check memory
top / htop               # see running processes
systemctl status docker  # check a service
```

## Related Concepts
- [[containers-docker]] — Lighter alternative to VMs for app isolation
- [[linux-server]] — How to operate the Ubuntu server inside the VM
- [[cloud-providers]] — Where to rent VMs
