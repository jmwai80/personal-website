# Linux Server Basics

## Essential Commands

### Navigation & Files
```bash
pwd                    # where am I?
ls -la                 # list files (including hidden)
cd /home/ubuntu        # change directory
mkdir -p /app/data     # create directory (and parents)
cp file1 file2         # copy
mv file1 file2         # move/rename
rm -rf /path           # delete (careful — no undo!)
cat file               # print file content
less file              # paginate file (q to quit)
tail -f /var/log/syslog  # follow log file live
```

### Users & Permissions
```bash
whoami                 # current user
sudo command           # run as root
sudo su                # switch to root shell
chmod 755 file         # rwxr-xr-x
chown ubuntu:ubuntu file  # change owner
```

File permissions: `rwxrwxrwx` = owner/group/others
- r=4, w=2, x=1
- 755 = owner:rwx, group:rx, others:rx

### Processes
```bash
ps aux                 # all running processes
htop                   # interactive process viewer
kill -9 <pid>          # force kill a process
lsof -i :8080          # what's using port 8080?
```

### Networking
```bash
curl http://localhost:8080         # test local service
wget https://example.com/file.zip  # download file
netstat -tlnp                      # listening ports
ss -tlnp                           # same, more modern
ping google.com                    # test connectivity
```

### Disk & Memory
```bash
df -h              # disk usage
du -sh /app        # size of directory
free -h            # RAM usage
```

## systemd — Service Management

systemd is the init system on Ubuntu. It manages long-running services (start on boot, auto-restart on crash).

```bash
systemctl status docker        # is Docker running?
systemctl start docker         # start
systemctl stop docker          # stop
systemctl restart docker       # restart
systemctl enable docker        # start on boot
systemctl disable docker       # don't start on boot
journalctl -u docker -f        # view service logs
```

## SSH

### Connect to VM
```bash
ssh ubuntu@<VM_IP>
ssh -i ~/.ssh/id_ed25519 ubuntu@<VM_IP>  # specify key explicitly
```

### SSH Key Auth (how it works)
```
Your Machine                    VM
~/.ssh/id_ed25519 (private)    ~/.ssh/authorized_keys (public key stored here)

On connect:
1. VM sends a challenge
2. Your machine signs it with the private key
3. VM verifies the signature using the public key
4. Match → connected (no password needed)
```

### Copy files to VM
```bash
scp ./file ubuntu@<VM_IP>:/home/ubuntu/
scp -r ./folder ubuntu@<VM_IP>:/home/ubuntu/
```

### SSH config shortcut (~/.ssh/config)
```
Host myvm
    HostName 123.45.67.89
    User ubuntu
    IdentityFile ~/.ssh/id_ed25519
```
Then just: `ssh myvm`

## UFW — Firewall

Ubuntu's simple firewall.

```bash
sudo ufw status
sudo ufw allow 22     # SSH (always keep this open!)
sudo ufw allow 80     # HTTP
sudo ufw allow 443    # HTTPS
sudo ufw enable
```

**Note**: Oracle Cloud has a *second* firewall layer (Security Lists in OCI Console). Both must allow the port.

## File Locations

| Path | What's There |
|------|-------------|
| `/home/ubuntu/` | Your home directory |
| `/etc/` | Config files |
| `/var/log/` | System logs |
| `/tmp/` | Temporary files (cleared on reboot) |
| `/usr/local/bin/` | Manually installed binaries |
| `/opt/` | Optional software |

## Our App on the VM

```bash
/home/ubuntu/app/
├── docker-compose.yml
├── Caddyfile
└── .env
```

Common tasks:
```bash
# See what's running
docker compose ps

# View app logs
docker compose logs -f personal-site

# Restart just one service
docker compose restart personal-site

# Deploy new version manually
docker compose pull && docker compose up -d

# Get a shell inside a container
docker compose exec personal-site sh
```

## Related Concepts
- [[virtual-machines]] — The VM running this OS
- [[containers-docker]] — What runs on top of the OS
