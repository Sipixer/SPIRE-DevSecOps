# Provisionne le serveur Hetzner. La configuration logicielle est faite par Ansible.

resource "hcloud_ssh_key" "main" {
  name       = "${var.server_name}-key"
  public_key = var.ssh_public_key
}

resource "hcloud_firewall" "main" {
  name = "${var.server_name}-fw"

  # SSH juste-à-temps : règle 22 présente seulement quand allowed_ssh_cidr est
  # renseigné (la CI y met l'IP du runner le temps d'Ansible). "" = port fermé.
  dynamic "rule" {
    for_each = var.allowed_ssh_cidr == "" ? [] : [var.allowed_ssh_cidr]
    content {
      direction  = "in"
      protocol   = "tcp"
      port       = "22"
      source_ips = [rule.value]
    }
  }

  # Aucun port web entrant : tout le HTTP passe par le tunnel Cloudflare (sortant).
}

resource "hcloud_server" "main" {
  name        = var.server_name
  server_type = var.server_type
  image       = var.image
  location    = var.location

  ssh_keys     = [hcloud_ssh_key.main.id]
  firewall_ids = [hcloud_firewall.main.id]

  labels = {
    projet = "spire-devsecops"
  }
}
