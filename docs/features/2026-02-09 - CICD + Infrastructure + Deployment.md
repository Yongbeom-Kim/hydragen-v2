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