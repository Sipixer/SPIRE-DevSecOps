# Accès à l'UI ArgoCD via Cloudflare Tunnel + Access (zero-trust).
#
# Principe : aucun port entrant n'est ouvert sur le serveur. Un pod cloudflared
# (déployé en GitOps par ArgoCD) établit une connexion SORTANTE vers Cloudflare.
# Cloudflare Access met une page de login (one-time PIN par email) DEVANT ArgoCD :
# seul access_allowed_email peut atteindre l'UI, tout le reste est refusé.
#
#   Navigateur ─▶ https://argocd.sylvainrougie.fr
#                  └─ Cloudflare Access (login email) ─▶ tunnel ─▶ argocd-server (ClusterIP)

# --- Le tunnel (remote-managed : la config d'ingress vit côté Cloudflare) ---
resource "cloudflare_zero_trust_tunnel_cloudflared" "argocd" {
  account_id = var.cloudflare_account_id
  name       = "spire-argocd"
  config_src = "cloudflare"
}

# Token utilisé par le pod cloudflared pour se connecter (injecté en Secret k8s).
data "cloudflare_zero_trust_tunnel_cloudflared_token" "argocd" {
  account_id = var.cloudflare_account_id
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.argocd.id
}

# Routage : le hostname public -> le service ArgoCD interne au cluster.
# argocd-server écoute en HTTPS avec un cert auto-signé, d'où le no-TLS-verify.
resource "cloudflare_zero_trust_tunnel_cloudflared_config" "argocd" {
  account_id = var.cloudflare_account_id
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.argocd.id
  config = {
    ingress = [
      {
        # argocd-server tourne en --insecure (TLS terminé par Cloudflare en
        # amont), donc on l'atteint en HTTP interne. Pas de boucle de redirection.
        hostname = var.argocd_hostname
        service  = "http://argocd-server.argocd.svc.cluster.local:80"
      },
      {
        service = "http_status:404"
      }
    ]
  }
}

# --- DNS : CNAME proxifié vers le tunnel ---
resource "cloudflare_dns_record" "argocd" {
  zone_id = data.cloudflare_zone.main.zone_id
  name    = var.argocd_hostname
  content = "${cloudflare_zero_trust_tunnel_cloudflared.argocd.id}.cfargotunnel.com"
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

data "cloudflare_zone" "main" {
  filter = {
    name = var.cloudflare_zone
  }
}

# --- Cloudflare Access : login obligatoire devant ArgoCD ---
resource "cloudflare_zero_trust_access_application" "argocd" {
  account_id       = var.cloudflare_account_id
  name             = "ArgoCD (SPIRE DevSecOps)"
  domain           = var.argocd_hostname
  type             = "self_hosted"
  session_duration = "24h"
  policies = [{
    id         = cloudflare_zero_trust_access_policy.argocd.id
    precedence = 1
  }]
}

# Policy deny-by-default : seul l'email autorisé passe (code OTP par email).
resource "cloudflare_zero_trust_access_policy" "argocd" {
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

# --- Sorties utiles ---
output "argocd_url" {
  value = "https://${var.argocd_hostname}"
}

# Token du tunnel : sert au Secret k8s du pod cloudflared. Sensible.
output "cloudflared_tunnel_token" {
  value     = data.cloudflare_zero_trust_tunnel_cloudflared_token.argocd.token
  sensitive = true
}
