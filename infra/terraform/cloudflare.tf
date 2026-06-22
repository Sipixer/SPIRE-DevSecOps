# Accès aux UI via un seul Cloudflare Tunnel (zero-trust, aucun port entrant).
#
# Un pod cloudflared (connexion SORTANTE) route 3 hostnames vers 3 services
# internes du cluster. Cloudflare Access met un login DEVANT les services
# sensibles (ArgoCD, Grafana) ; la boutique de démo (Gateway) reste publique.
#
#   argocd.sylvainrougie.fr  → Access(login) → argocd-server   (admin cluster)
#   grafana.sylvainrougie.fr → Access(login) → grafana          (observabilité)
#   shop.sylvainrougie.fr    → (public)       → gateway          (démo boutique)

# --- Le tunnel (remote-managed) ---
resource "cloudflare_zero_trust_tunnel_cloudflared" "main" {
  account_id = var.cloudflare_account_id
  name       = "spire-argocd"
  config_src = "cloudflare"
}

data "cloudflare_zero_trust_tunnel_cloudflared_token" "main" {
  account_id = var.cloudflare_account_id
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.main.id
}

# Routage : chaque hostname public -> un service interne (HTTP, TLS terminé par
# Cloudflare). argocd-server et grafana tournent en mode insecure derrière le tunnel.
resource "cloudflare_zero_trust_tunnel_cloudflared_config" "main" {
  account_id = var.cloudflare_account_id
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.main.id
  config = {
    ingress = [
      {
        hostname = var.argocd_hostname
        service  = "http://argocd-server.argocd.svc.cluster.local:80"
      },
      {
        hostname = var.grafana_hostname
        service  = "http://kube-prometheus-stack-grafana.monitoring.svc.cluster.local:80"
      },
      {
        hostname = var.shop_hostname
        service  = "http://gateway.shop.svc.cluster.local:8080"
      },
      {
        service = "http_status:404"
      }
    ]
  }
}

# --- DNS : un CNAME proxifié par hostname, tous vers le tunnel ---
locals {
  tunnel_cname = "${cloudflare_zero_trust_tunnel_cloudflared.main.id}.cfargotunnel.com"
}

data "cloudflare_zone" "main" {
  filter = {
    name = var.cloudflare_zone
  }
}

resource "cloudflare_dns_record" "argocd" {
  zone_id = data.cloudflare_zone.main.zone_id
  name    = var.argocd_hostname
  content = local.tunnel_cname
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_dns_record" "grafana" {
  zone_id = data.cloudflare_zone.main.zone_id
  name    = var.grafana_hostname
  content = local.tunnel_cname
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

resource "cloudflare_dns_record" "shop" {
  zone_id = data.cloudflare_zone.main.zone_id
  name    = var.shop_hostname
  content = local.tunnel_cname
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

# --- Cloudflare Access : login obligatoire devant ArgoCD et Grafana ---
# La boutique (shop) n'a PAS d'application Access => publique.
resource "cloudflare_zero_trust_access_application" "argocd" {
  account_id       = var.cloudflare_account_id
  name             = "ArgoCD (SPIRE DevSecOps)"
  domain           = var.argocd_hostname
  type             = "self_hosted"
  session_duration = "24h"
  policies = [{
    id         = cloudflare_zero_trust_access_policy.allow_admin.id
    precedence = 1
  }]
}

resource "cloudflare_zero_trust_access_application" "grafana" {
  account_id       = var.cloudflare_account_id
  name             = "Grafana (SPIRE DevSecOps)"
  domain           = var.grafana_hostname
  type             = "self_hosted"
  session_duration = "24h"
  policies = [{
    id         = cloudflare_zero_trust_access_policy.allow_admin.id
    precedence = 1
  }]
}

# Policy partagée : seul l'email autorisé passe (code OTP par email), deny par défaut.
resource "cloudflare_zero_trust_access_policy" "allow_admin" {
  account_id       = var.cloudflare_account_id
  name             = "Autoriser l'admin"
  decision         = "allow"
  session_duration = "24h"
  include = [{
    email = {
      email = var.access_allowed_email
    }
  }]
}

# --- Sorties ---
output "argocd_url" {
  value = "https://${var.argocd_hostname}"
}

output "grafana_url" {
  value = "https://${var.grafana_hostname}"
}

output "shop_url" {
  value = "https://${var.shop_hostname}"
}

output "cloudflared_tunnel_token" {
  value     = data.cloudflare_zero_trust_tunnel_cloudflared_token.main.token
  sensitive = true
}
