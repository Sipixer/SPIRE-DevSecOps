# Un Cloudflare Tunnel (sortant, aucun port entrant) route 4 hostnames vers le
# cluster. Access (login) protège argocd/grafana/kiali ; shop reste public.
resource "cloudflare_zero_trust_tunnel_cloudflared" "main" {
  account_id = var.cloudflare_account_id
  name       = "spire-argocd"
  config_src = "cloudflare"
}

data "cloudflare_zero_trust_tunnel_cloudflared_token" "main" {
  account_id = var.cloudflare_account_id
  tunnel_id  = cloudflare_zero_trust_tunnel_cloudflared.main.id
}

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
        hostname = var.kiali_hostname
        service  = "http://kiali.istio-system.svc.cluster.local:20001"
      },
      {
        service = "http_status:404"
      }
    ]
  }
}

# DNS : un CNAME proxifié par hostname, tous vers le tunnel.
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

resource "cloudflare_dns_record" "kiali" {
  zone_id = data.cloudflare_zone.main.zone_id
  name    = var.kiali_hostname
  content = local.tunnel_cname
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

# Cloudflare Access : login obligatoire devant argocd/grafana/kiali (shop reste public).
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

resource "cloudflare_zero_trust_access_application" "kiali" {
  account_id       = var.cloudflare_account_id
  name             = "Kiali (SPIRE DevSecOps)"
  domain           = var.kiali_hostname
  type             = "self_hosted"
  session_duration = "24h"
  policies = [{
    id         = cloudflare_zero_trust_access_policy.allow_admin.id
    precedence = 1
  }]
}

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
