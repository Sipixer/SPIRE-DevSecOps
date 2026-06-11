output "server_ip" {
  description = "IP publique du serveur (cible d'Ansible)"
  value       = hcloud_server.main.ipv4_address
}

output "server_status" {
  description = "État du serveur"
  value       = hcloud_server.main.status
}

output "ssh_command" {
  description = "Commande pour se connecter au serveur"
  value       = "ssh root@${hcloud_server.main.ipv4_address}"
}

# Inventaire Ansible prêt à l'emploi, écrit à la racine d'infra/ansible.
resource "local_file" "ansible_inventory" {
  filename = "${path.module}/../ansible/inventory.ini"
  content  = <<-EOT
    [k3s]
    ${hcloud_server.main.ipv4_address}

    [k3s:vars]
    ansible_user=root
    ansible_ssh_common_args='-o StrictHostKeyChecking=accept-new'
  EOT
}
