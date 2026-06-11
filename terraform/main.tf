provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = var.kube_context
}

provider "helm" {
  kubernetes {
    config_path    = "~/.kube/config"
    config_context = var.kube_context
  }
}

resource "kubernetes_namespace" "spire" {
  metadata {
    name = var.namespace
    labels = {
      "pod-security.kubernetes.io/enforce" = "privileged"
    }
  }
}

# Les CRDs SPIRE (ClusterSPIFFEID, etc.) — doivent exister avant le chart principal.
resource "helm_release" "spire_crds" {
  name       = "spire-crds"
  repository = "https://spiffe.github.io/helm-charts-hardened/"
  chart      = "spire-crds"
  version    = var.spire_crds_chart_version
  namespace  = kubernetes_namespace.spire.metadata[0].name
}

# Le chart SPIRE hardened : Server + Agent (DaemonSet) + controller manager.
resource "helm_release" "spire" {
  name       = "spire"
  repository = "https://spiffe.github.io/helm-charts-hardened/"
  chart      = "spire"
  version    = var.spire_chart_version
  namespace  = kubernetes_namespace.spire.metadata[0].name

  values = [
    templatefile("${path.module}/values/spire.yaml.tftpl", {
      trust_domain = var.trust_domain
    })
  ]

  depends_on = [helm_release.spire_crds]
}
