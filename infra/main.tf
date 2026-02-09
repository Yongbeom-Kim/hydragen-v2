# Resolve zone name: use explicit zone_name or derive from base_domain (e.g. example.com -> example-com)
# Look up the existing hosted zone by name (data block)
data "google_dns_managed_zone" "zone" {
  name    = var.zone_name
  project = var.project_id
}

# A record for the actual domain
resource "google_dns_record_set" "a" {
  name         = "${var.actual_domain}."
  type         = "A"
  ttl          = var.ttl
  managed_zone = data.google_dns_managed_zone.zone.name
  project      = var.project_id
  rrdatas      = [var.ipv4_address]
}

# AAAA record for the actual domain
resource "google_dns_record_set" "aaaa" {
  name         = "${var.actual_domain}."
  type         = "AAAA"
  ttl          = var.ttl
  managed_zone = data.google_dns_managed_zone.zone.name
  project      = var.project_id
  rrdatas      = [var.ipv6_address]
}
