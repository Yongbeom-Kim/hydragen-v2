# 2026-02-09 - CI/CD + Infrastructure + Deployment

## Server Setup

New Hetzner Server: https://console.hetzner.com/projects/13384053/servers/120476441/overview
- CX23: 2vCPU, 4GB RAM, 40GB Disk local
- USD 3.80/month ≈ SGD 4.82/month
- IPv4: `5.75.177.76`
- IPv6: `2a01:4f8:1c1a:eda1::/64`

## DNS
- Cloud DNS Console: [link](https://console.cloud.google.com/net-services/dns/zones/senpailearn-com/details?authuser=1&project=senpailearn-global)
- Implemented with OpenTofu + Encrypted local state (committed to git)

## Reverse Proxy & TLS Termination
- We use Caddy for reverse proxy and TLS Termination.
	- Dev: http://localhost:3111
	- Prod: https://hydragen.senpailearn.com

[Caddyfile](../../caddy/Caddyfile):
```Caddyfile
{$SITE_ADDR} {
  @root path /
  redir @root /v2 permanent

  @api path /v2/api*
  handle @api {
    uri strip_prefix /v2/api
    reverse_proxy server:8080
  }

  handle /v2* {
    reverse_proxy client:3000
  }

  handle {
    respond "Not Found" 404
  }
}
```
## Tanstack Start
Basic Tanstack Start App hosted on hydragen.senpailearn.com/v2

## Go Backend
Basic Go backend hosted on hydragen.senpailearn.com/v2/api
Healthcheck endpoint on hydragen.senpailearn.com/v2/api/health