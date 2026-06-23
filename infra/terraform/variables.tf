variable "hcloud_token" {
  description = "Token API Hetzner Cloud (secret, fourni par Terraform Cloud)"
  type        = string
  sensitive   = true
}

variable "server_name" {
  description = "Nom du serveur"
  type        = string
  default     = "spire-devsecops"
}

variable "server_type" {
  description = "Type de serveur Hetzner (cpx32 = 4 vCPU / 8 Go, AMD x86)"
  type        = string
  default     = "cpx32"
}

variable "location" {
  description = "Datacenter Hetzner (fsn1 Falkenstein, nbg1 Nuremberg, hel1 Helsinki)"
  type        = string
  default     = "fsn1"
}

variable "image" {
  description = "Image de base du serveur"
  type        = string
  default     = "ubuntu-24.04"
}

variable "ssh_public_key" {
  description = "Clé SSH publique autorisée à se connecter au serveur"
  type        = string
}

variable "cloudflare_api_token" {
  description = "Token API Cloudflare (Tunnel + Access + DNS). Variable sensitive du workspace Terraform Cloud."
  type        = string
  sensitive   = true
}

variable "cloudflare_account_id" {
  description = "ID du compte Cloudflare propriétaire du tunnel et de la zone"
  type        = string
  default     = "37da22730527e0b95523e6f5e3792309"
}

variable "cloudflare_zone" {
  description = "Domaine géré par Cloudflare"
  type        = string
  default     = "sylvainrougie.fr"
}

variable "argocd_hostname" {
  description = "Sous-domaine exposant l'UI ArgoCD via le tunnel Cloudflare"
  type        = string
  default     = "argocd.sylvainrougie.fr"
}

variable "grafana_hostname" {
  description = "Sous-domaine exposant Grafana via le tunnel (protégé par Access)"
  type        = string
  default     = "grafana.sylvainrougie.fr"
}

variable "shop_hostname" {
  description = "Sous-domaine exposant la boutique de démo (Gateway) via le tunnel, public (sans Access)"
  type        = string
  default     = "shop.sylvainrougie.fr"
}

variable "kiali_hostname" {
  description = "Sous-domaine exposant le dashboard Kiali (mesh Istio) via le tunnel (protégé par Access)"
  type        = string
  default     = "kiali.sylvainrougie.fr"
}

variable "access_allowed_email" {
  description = "Email autorisé à passer Cloudflare Access pour atteindre ArgoCD"
  type        = string
  default     = "sipixer@gmail.com"
}

variable "allowed_ssh_cidr" {
  description = <<-EOT
    CIDR autorisé pour SSH. Vide ("") = port 22 FERMÉ (état au repos, zero-trust).
    La CI le passe à l'IP du runner (<ip>/32) le temps du provisioning Ansible,
    puis le remet à "" pour refermer. SSH n'est donc jamais ouvert en permanence.
  EOT
  type        = string
  default     = ""
}
