# DNS & TLS (Domains, HTTPS)

## How a Domain Request Works End-to-End

When you type `yourdomain.com` in a browser:

```
1. Browser asks: "What IP is yourdomain.com?"
        │
        ▼
2. DNS Resolver (e.g. 8.8.8.8 — Google's DNS)
        │  checks cache, if miss:
        ▼
3. Root DNS servers → TLD servers (.com) → GoDaddy's nameservers
        │
        ▼
4. GoDaddy returns: "yourdomain.com → 123.45.67.89" (your VM's IP)
        │
        ▼
5. Browser connects to 123.45.67.89:443 (HTTPS)
        │
        ▼
6. Caddy on VM receives the connection
        │
        ▼
7. TLS handshake (Caddy proves it owns the domain via cert)
        │
        ▼
8. Encrypted HTTP request: GET / HTTP/1.1 Host: yourdomain.com
        │
        ▼
9. Caddy forwards to personal-site:8080
        │
        ▼
10. Go app responds → Caddy encrypts → browser renders
```

## DNS Record Types

| Type | Purpose | Example |
|------|---------|---------|
| **A** | Maps domain to IPv4 address | `yourdomain.com → 123.45.67.89` |
| **AAAA** | Maps domain to IPv6 address | `yourdomain.com → 2001:db8::1` |
| **CNAME** | Alias to another domain | `www → yourdomain.com` |
| **MX** | Mail server | `yourdomain.com → mail.google.com` |
| **TXT** | Verification, SPF, etc | Various |

### What We Set in GoDaddy

```
Type  Name          Value              TTL
A     @             <VM IP>            600
A     www           <VM IP>            600
A     apartments    <VM IP>            600
```

- `@` means the root domain (`yourdomain.com`)
- TTL = 600 seconds = how long DNS resolvers cache the answer (lower = faster updates)

## TLS / HTTPS

**TLS (Transport Layer Security)** encrypts the connection between browser and server. HTTPS = HTTP over TLS.

### Why it matters
- Data in transit is encrypted (passwords, personal info safe)
- Browser shows padlock icon (trust signal)
- Required for many browser APIs (geolocation, service workers)
- SEO ranking factor

### Certificates
A TLS certificate proves you own the domain. Issued by a **Certificate Authority (CA)**.

- **Let's Encrypt** — free, automated, what Caddy uses
- DigiCert, Comodo — paid, used by enterprises

A cert contains:
- Your domain name
- Your public key
- CA's signature (proving authenticity)
- Expiry date (Let's Encrypt: 90 days, auto-renewed)

### TLS Handshake (simplified)

```
Browser                     Caddy
   │──── "Hello, I want TLS" ──▶│
   │◀── "Here's my certificate" ─│
   │  (browser verifies cert     │
   │   is signed by trusted CA)  │
   │──── "Here's session key" ──▶│
   │◀────── "Acknowledged" ───────│
   │                             │
   │═══ Encrypted channel open ══│
```

## Propagation

When you change a DNS record, it doesn't update everywhere instantly. Each resolver has cached the old value for its TTL duration.

- With TTL 600s: most resolvers updated within 10 minutes
- With TTL 3600s: can take up to 1 hour
- **Tip**: Lower TTL to 600 *before* making changes, then change the record

Check propagation: `dig yourdomain.com` or use whatsmydns.net

## Subdomain vs Path Routing

| Approach | Example | Config |
|----------|---------|--------|
| Subdomain | `apartments.yourdomain.com` | DNS A record + Caddy block |
| Path | `yourdomain.com/apartments` | Single Caddy block, path matching |

We use subdomains — cleaner separation, easier to split to separate servers later.

## Related Concepts
- [[reverse-proxy]] — Caddy handles TLS termination
- [[virtual-machines]] — The VM that holds your server IP
