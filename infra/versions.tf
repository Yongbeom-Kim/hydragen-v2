terraform {
  required_version = ">= 1.5.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }

  encryption {
    key_provider "pbkdf2" "state_key" {
      passphrase = var.passphrase
    }
    method "aes_gcm" "state_method" {
      keys = key_provider.pbkdf2.state_key
    }
    state {
      method = method.aes_gcm.state_method
    }
  }
}
