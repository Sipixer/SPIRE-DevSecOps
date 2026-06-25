# Mission

## En une phrase

Maîtriser **en profondeur le « pourquoi »** de chaque choix technologique du projet
SPIRE/DevSecOps — au point de pouvoir le **défendre face à un jury de soutenance**
qui va challenger chaque décision (« pourquoi Istio et pas Linkerd ? », « SPIRE vs
cert-manager ? », « pourquoi pas juste un token ? »).

## Pour qui / contexte

- **Apprenant** : Sylvain, étudiant DevSecOps. A *construit* le projet (le code, la
  pipeline, l'infra existent et tournent). Sait *quoi* il a fait ; veut maintenant
  blinder le *pourquoi*.
- **Échéance** : soutenance à venir. Priorité = être inattaquable à l'oral.
- **Au-delà** : la compréhension de fond sert aussi la carrière DevSecOps, donc on
  vise des modèles mentaux solides, pas du par-cœur.

## Le projet à défendre (résumé)

Une mini-boutique microservices **zero-trust** qui a **évolué en 2 versions** — et
c'est cette évolution qui est le cœur du récit :

- **v1 — Identité dans le code** : SPIFFE/**SPIRE** émet des identités de workload
  (X.509-SVID) ; chaque service Go câble le **mTLS à la main** (`go-spiffe/v2`) et
  applique l'autorisation dans le code. Analytics (Node) est mis derrière un
  **Envoy** manuel.
- **v2 — Identité dans la plateforme** : un **service mesh (Istio)** porte le mTLS,
  l'identité et l'autorisation (`AuthorizationPolicy`). Le code applicatif redevient
  du HTTP clair. Dashboard **Kiali**.
- **Transverse** : pipeline **CI/CD sécurisée** (tflint/Checkov → Semgrep → Trivy →
  Cosign → GHCR → branche `deploy/prod` → **ArgoCD** GitOps) et **infra zero-trust**
  (Terraform/Hetzner + Ansible/k3s, **SSH juste-à-temps**, tunnel Cloudflare).

La règle métier-fil-rouge : **seul Orders peut déclencher un paiement** (`/pay`).
Gateway est authentifiée mais **pas autorisée** → 403. C'est la démonstration vivante
du zero-trust (« authentifié ≠ autorisé »).

## Ce qui compte pour juger les leçons

Chaque leçon doit, à terme, permettre de répondre à 3 questions de jury :
1. **C'est quoi ?** (définition juste, vocabulaire maîtrisé)
2. **Pourquoi celui-là ?** (le problème qu'il résout, vs ne rien faire)
3. **Pourquoi pas le concurrent ?** (alternatives + tradeoff assumé)

## Plan de couverture (4 axes, tout doit être vu)

1. **Identité & mTLS** : SPIFFE, SPIRE, SVID, attestation, trust domain ; mTLS vs
   tokens/API keys ; alternatives à SPIRE (cert-manager, Vault, mesh natif).
2. **Service mesh** : sidecar/Envoy, plan de contrôle/données ; v1→v2 (pourquoi) ;
   Istio vs Linkerd vs Cilium ; ambient mesh.
3. **Pipeline & supply chain** : SAST/DAST/SCA/IaC-scan, qui fait quoi (Semgrep,
   Trivy, ZAP, Checkov) ; signature/SLSA (Cosign) ; GitOps (ArgoCD) + `deploy/prod`.
4. **Infra & zero-trust réseau** : Terraform vs Ansible ; k3s vs k8s ; SSH JIT ;
   firewall/tunnel Cloudflare ; Terraform Cloud (remote state, secrets).

## Statut

Mission posée le 2026-06-24. Objectif : « les deux, soutenance d'abord ».
