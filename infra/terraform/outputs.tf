output "server_ip" {
  description = "IP publique du serveur"
  value       = hcloud_server.main.ipv4_address
}

output "ssh_command" {
  value = "ssh -i ~/.ssh/spire-devsecops root@${hcloud_server.main.ipv4_address}"
}

output "argocd_url" {
  value = "https://${var.argocd_hostname}"
}

output "grafana_url" {
  value = "https://${var.grafana_hostname}"
}

output "shop_url" {
  value = "https://${var.shop_hostname}"
}

output "kiali_url" {
  value = "https://${var.kiali_hostname}"
}

output "cloudflared_tunnel_token" {
  value     = data.cloudflare_zero_trust_tunnel_cloudflared_token.main.token
  sensitive = true
}
