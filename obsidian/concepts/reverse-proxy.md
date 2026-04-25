# Reverse Proxy (Caddy / Nginx)

## What is a Reverse Proxy?

A reverse proxy sits in front of your apps and forwards incoming requests to the right backend service.

```
Internet
    │
    ▼
Caddy (port 443, HTTPS)         ← reverse proxy
    │
    ├── yourdomain.com      → personal-site:8080
    └── apartments.com      → apartments-site:8081
```

Without it, you'd have to expose each app on a different port (`:8080`, `:8081`) — ugly URLs, no HTTPS, no central place to handle TLS.

## Why You Need One

| Problem | Reverse Proxy Solution |
|---------|----------------------|
| Multiple apps, one server | Routes by domain/path |
| HTTPS (TLS) termination | Handles certs, your app only speaks HTTP internally |
| HTTP → HTTPS redirect | Done automatically |
| Load balancing | Can spread traffic across multiple app instances |
| Rate limiting, auth | Central enforcement point |
| Static file serving | Serve assets without hitting your Go app |

## Forward Proxy vs Reverse Proxy

- **Forward proxy**: sits in front of *clients* (e.g. VPN, corporate firewall — the client knows about it)
- **Reverse proxy**: sits in front of *servers* (the client doesn't know — they just hit your domain)

## Caddy

What we use. Written in Go. Key advantage: **automatic HTTPS**.

### Caddyfile
```
yourdomain.com {
    reverse_proxy personal-site:8080
}

apartments.yourdomain.com {
    reverse_proxy apartments-site:8081
}
```

That's it. Caddy automatically:
- Gets a TLS cert from Let's Encrypt
- Renews it before expiry
- Redirects HTTP to HTTPS
- Handles the ACME challenge

### How Caddy gets HTTPS working
1. Caddy sees a new domain in its config
2. Makes an HTTP-01 challenge request to Let's Encrypt
3. Let's Encrypt verifies you control the domain by hitting `yourdomain.com/.well-known/acme-challenge/<token>`
4. Let's Encrypt issues a cert valid for 90 days
5. Caddy stores it in its data volume and auto-renews at 30 days remaining

**Requirement**: DNS must point to your server IP *before* Caddy starts, or the cert request fails.

## Nginx (alternative)

More widely used, more config required. Same capabilities.

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://personal-site:8080;
    }
}
```

Plus you'd need Certbot separately to manage certs. Caddy does all this automatically.

## TLS Termination

Your Go app listens on plain HTTP (`:8080`). Caddy handles the encryption/decryption. This is called **TLS termination**.

```
Browser ──HTTPS──▶ Caddy ──HTTP──▶ Go App
                   (decrypts)      (plain text internally)
```

This is fine — the internal network (inside the VM) is trusted. The encrypted channel is between the user's browser and Caddy.

## Related Concepts
- [[dns-tls]] — How DNS, TLS certificates, and HTTPS fit together
- [[containers-docker]] — Caddy runs as a container in our setup
