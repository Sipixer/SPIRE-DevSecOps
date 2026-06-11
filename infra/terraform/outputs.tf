output "server_ip" {
  description = "IP publique du serveur"
  value       = hcloud_server.main.ipv4_address
}

output "ssh_command" {
  value = "ssh -i ~/.ssh/spire-devsecops root@${hcloud_server.main.ipv4_address}"
}

# Inventaire Ansible généré pour la cible.
resource "local_file" "ansible_inventory" {
  filename        = "${path.module}/../ansible/inventory.ini"
  file_permission = "0644"
  content         = <<-EOT
    [k3s]
    ${hcloud_server.main.ipv4_address}
  EOT
}
