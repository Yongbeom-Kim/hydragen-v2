variable "passphrase" {
  description = "Passphrase for encrypting local Terraform/OpenTofu state. Use at least 16 characters. Prefer setting via TF_VAR_passphrase or -var instead of default."
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "GCP project ID where the Cloud DNS managed zone lives."
  type        = string
}

variable "base_domain" {
  description = "Base domain of the hosted zone (e.g. example.com). Used to look up the zone and to scope records."
  type        = string
}

variable "actual_domain" {
  description = "Domain name for the A and AAAA records (e.g. app.example.com or sub.example.com). Must be within the base domain."
  type        = string
}

variable "ipv4_address" {
  description = "IPv4 address for the A record."
  type        = string
}

variable "ipv6_address" {
  description = "IPv6 address for the AAAA record."
  type        = string
}

variable "zone_name" {
  description = "Name of the GCP Cloud DNS managed zone (resource name in GCP). If unset, derived from base_domain (e.g. example.com -> example-com)."
  type        = string
  default     = ""
}

variable "ttl" {
  description = "TTL in seconds for the A and AAAA records."
  type        = number
  default     = 300
}
