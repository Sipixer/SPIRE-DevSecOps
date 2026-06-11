output "server_ip" {
  description = "IP publique du serveur"
  value       = hcloud_server.main.ipv4_address
}

output "ssh_command" {
  value = "ssh -i ~/.ssh/spire-devsecops root@${hcloud_server.main.ipv4_address}"
}
