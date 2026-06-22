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

variable "allowed_ssh_cidr" {
  description = <<-EOT
    CIDR autorisé pour SSH. Vide ("") = port 22 FERMÉ (état au repos, zero-trust).
    La CI le passe à l'IP du runner (<ip>/32) le temps du provisioning Ansible,
    puis le remet à "" pour refermer. SSH n'est donc jamais ouvert en permanence.
  EOT
  type        = string
  default     = ""
}
