# Terraform Hetzner : crée uniquement le serveur vide et son périmètre réseau.
# Toute la configuration logicielle (k3s, ArgoCD, SPIRE, apps) est faite par
# Ansible dans une étape distincte. Terraform ne touche pas au contenu du cluster.

resource "hcloud_ssh_key" "main" {
  name       = "${var.server_name}-key"
  public_key = var.ssh_public_key
}

# Firewall : on n'ouvre que le strict nécessaire. L'API k3s (6443) reste fermée
# au monde — le déploiement se fait en pull (ArgoCD) depuis l'intérieur.
resource "hcloud_firewall" "main" {
  name = "${var.server_name}-fw"

  rule {
    description = "SSH (administration + Ansible)"
    direction   = "in"
    protocol    = "tcp"
    port        = "22"
    source_ips  = [var.allowed_ssh_cidr]
  }

  rule {
    description = "HTTP (ingress, redirige vers HTTPS)"
    direction   = "in"
    protocol    = "tcp"
    port        = "80"
    source_ips  = ["0.0.0.0/0", "::/0"]
  }

  rule {
    description = "HTTPS (front public via ingress)"
    direction   = "in"
    protocol    = "tcp"
    port        = "443"
    source_ips  = ["0.0.0.0/0", "::/0"]
  }
}

resource "hcloud_server" "main" {
  name        = var.server_name
  server_type = var.server_type
  image       = var.image
  location    = var.location

  ssh_keys     = [hcloud_ssh_key.main.id]
  firewall_ids = [hcloud_firewall.main.id]

  public_net {
    ipv4_enabled = true
    ipv6_enabled = true
  }

  labels = {
    projet = "spire-devsecops"
    role   = "k3s"
  }
}
