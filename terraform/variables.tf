variable "kube_context" {
  description = "Contexte kubectl du cluster k3d cible"
  type        = string
  default     = "k3d-spire"
}

variable "namespace" {
  description = "Namespace où déployer SPIRE"
  type        = string
  default     = "spire"
}

variable "trust_domain" {
  description = "Trust domain SPIFFE"
  type        = string
  default     = "example.org"
}

variable "spire_chart_version" {
  description = "Version du chart Helm spiffe/spire"
  type        = string
  default     = "0.29.0"
}

variable "spire_crds_chart_version" {
  description = "Version du chart Helm spiffe/spire-crds"
  type        = string
  default     = "0.5.0"
}
