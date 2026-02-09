output "zone_name" {
  description = "Name of the managed zone."
  value       = data.google_dns_managed_zone.zone.name
}

output "zone_dns_name" {
  description = "DNS name of the managed zone (e.g. example.com.)."
  value       = data.google_dns_managed_zone.zone.dns_name
}

output "a_record_fqdn" {
  description = "FQDN of the A record."
  value       = google_dns_record_set.a.name
}

output "aaaa_record_fqdn" {
  description = "FQDN of the AAAA record."
  value       = google_dns_record_set.aaaa.name
}
