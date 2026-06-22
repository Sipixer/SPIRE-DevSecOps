terraform {
  required_version = ">= 1.6"

  # State et secrets (token Hetzner, clé SSH) centralisés dans Terraform Cloud.
  # Renseigner l'organisation et le workspace dans cloud.auto.tfbackend ou ici.
  # Terraform Cloud crée le workspace au premier `terraform init` s'il n'existe pas.
  cloud {
    organization = "Sipixer"
    workspaces {
      name = "spire-devsecops"
    }
  }

  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.49"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
  }
}

provider "hcloud" {
  token = var.hcloud_token
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}
