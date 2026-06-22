# Provisionne le serveur Hetzner. La configuration logicielle est faite par Ansible.

resource "hcloud_ssh_key" "main" {
  name       = "${var.server_name}-key"
  public_key = var.ssh_public_key
}

resource "hcloud_firewall" "main" {
  name = "${var.server_name}-fw"

  # Règle SSH JUSTE-À-TEMPS : présente uniquement quand allowed_ssh_cidr est
  # renseigné (la CI y met l'IP du runner le temps d'Ansible). Au repos la
  # variable vaut "" → aucune règle 22 → port fermé. C'est le cœur du zero-trust :
  # SSH n'existe au pare-feu que pour l'IP qui en a besoin, le temps du job.
  dynamic "rule" {
    for_each = var.allowed_ssh_cidr == "" ? [] : [var.allowed_ssh_cidr]
    content {
      direction  = "in"
      protocol   = "tcp"
      port       = "22"
      source_ips = [rule.value]
    }
  }

  # Zero-trust total : AUCUN port web entrant. Tout le trafic HTTP (ArgoCD,
  # Grafana, boutique) passe par le tunnel Cloudflare, qui établit une connexion
  # SORTANTE — donc rien à ouvrir en entrée. Au repos, ce firewall ne laisse
  # passer que SSH juste-à-temps (et seulement pendant un provisioning).
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
