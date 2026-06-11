output "namespace" {
  description = "Namespace où SPIRE est déployé"
  value       = kubernetes_namespace.spire.metadata[0].name
}

output "trust_domain" {
  description = "Trust domain SPIFFE configuré"
  value       = var.trust_domain
}

output "verify_command" {
  description = "Commande pour vérifier que SPIRE tourne"
  value       = "kubectl get pods -n ${var.namespace}"
}
