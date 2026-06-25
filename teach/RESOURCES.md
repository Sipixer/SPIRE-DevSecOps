# Ressources

Sources de confiance pour ancrer les leçons. Priorité : specs officielles, NIST,
CNCF, papiers originaux > blogs. Vérifié = URL testée et citation extraite du PDF/page.

## Zero-trust (fondations)

- **NIST SP 800-207 — Zero Trust Architecture** ✅ *(LA citation canonique)*
  - PDF : https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-207.pdf
  - Landing : https://csrc.nist.gov/pubs/sp/800/207/final
  - Tenet 2 (verbatim, §2.1) : « All communication is secured regardless of network
    location. Network location alone does not imply trust. »
  - §1 : « once attackers breach the perimeter, further lateral movement is unhindered. »
  - Standard fédéral US, août 2020. Autorité maximale.
  - ⚠️ NIST ne dit PAS « never trust, always verify » → dit « no implicit trust ».

- **Forrester / John Kindervag — "No More Chewy Centers" (Zero Trust Model)** ✅
  - PDF (copie publique) : https://media.paloaltonetworks.com/documents/Forrester-No-More-Chewy-Centers.pdf
  - Origine du slogan (2010). Wording original : « verify and never trust ».
  - C'est Forrester, PAS NIST, qui a forgé le mantra. Attribution précise = point jury.

- **Google BeyondCorp (Ward & Beyer, ;login: déc. 2014)** ✅ *(zero-trust réel à l'échelle)*
  - Google Research : https://research.google/pubs/beyondcorp-a-new-approach-to-enterprise-security/
  - USENIX PDF : https://www.usenix.org/system/files/login/articles/login_dec14_02_ward.pdf
  - Thèse : déplacer le contrôle d'accès du périmètre réseau vers l'utilisateur+device.

- **Jericho Forum / déperimétrisation (Spencer & Pizio, 2023, open access)** ✅ *(historique)*
  - https://pmc.ncbi.nlm.nih.gov/articles/PMC11528882/
  - Terme « de-perimeterised » forgé en 2001 (Jon Measham, Royal Mail) ; Jericho Forum 2003.

## SPIFFE / SPIRE (identité de workload)

- **spiffe.io — Concepts** ✅ *(meilleure source unique pour SPIFFE/SVID/trust domain/Workload API)*
  - https://spiffe.io/docs/latest/spiffe-about/spiffe-concepts/
- **spiffe.io — Overview** ✅
  - https://spiffe.io/docs/latest/spiffe-about/overview/
  - « SPIFFE … is a set of open-source standards for securely identifying software
    systems in dynamic and heterogeneous environments. »
  - « network policies that only allow traffic between particular IP addresses
    struggle to scale under this complexity. » (= le pourquoi de l'identité de workload)
- **spiffe.io — SPIFFE ID spec** (format `spiffe://trust-domain/path`) ✅
  - https://spiffe.io/docs/latest/spiffe-specs/spiffe-id/
- **spiffe.io — Workload API spec** ✅
  - https://spiffe.io/docs/latest/spiffe-specs/spiffe_workload_api/
  - « the Workload API does not require that a calling workload … possess any
    authentication token. » (= pas de secret à co-déployer = bootstrap trust)
- **CNCF — SPIFFE & SPIRE : Graduated depuis août 2022** ✅
  - https://www.cncf.io/projects/spiffe/ · https://www.cncf.io/projects/spire/
  - ⚠️ « Graduated », PAS « incubating ». Dire « incubating » = obsolète.

## Service mesh (Istio / Linkerd / Cilium) ✅ vérifié juin 2026

- **Istio — architecture (data/control plane, Envoy, istiod)** ✅
  - https://istio.io/latest/docs/ops/deployment/architecture/
- **Istio — Ambient mode GA (Istio 1.24, nov. 2024)** ✅ — ztunnel (L4/node) + waypoints (L7)
  - https://istio.io/latest/blog/2024/ambient-reaches-ga/ · https://istio.io/latest/docs/ambient/overview/
- **Istio ↔ SPIFFE/SPIRE (identité SVID native)** ✅ *(argument-clé du projet)*
  - https://istio.io/latest/docs/ops/integrations/spire/
  - « Istio and SPIFFE share the same identity document: SVID. »
- **Kiali (dashboard riche + validation config)** ✅ *(la vraie raison du choix Istio)*
  - https://istio.io/latest/docs/ops/integrations/kiali/ · https://kiali.io/docs/features/topology/
- **Linkerd — micro-proxy Rust (pourquoi pas Envoy)** ✅
  - https://linkerd.io/2020/12/03/why-linkerd-doesnt-use-envoy/
- **Linkerd 2.15 — SPIFFE pour workloads hors-cluster seulement** ✅
  - https://linkerd.io/2024/02/21/announcing-linkerd-2.15/
- **⚠️ Friction Linkerd : stables payantes (>50 salariés, ~2000$/cluster/mois, mai 2024)** ✅
  - https://thenewstack.io/some-linkerd-users-must-pay-fear-and-anger-explained/
  - Argument pro-Istio (pas de péage éditeur). Le code + edge restent libres.
- **Cilium — eBPF sans sidecar, mutual auth = Beta, SPIFFE/SPIRE** ✅
  - https://docs.cilium.io/en/stable/network/servicemesh/
- **CNCF Graduated** : Linkerd juil. 2021 (1er), Istio juil. 2023, Cilium oct. 2023 ✅
- **⚠️ Benchmarks = TOUS vendor-biaisés** (Buoyant favorise Linkerd, Istio favorise Istio).
  Pas de benchmark CNCF neutre récent. Affirmer le fait architectural (Linkerd + léger par
  conception), JAMAIS un chiffre précis comme vérité neutre.
  - https://linkerd.io/2021/11/29/linkerd-vs-istio-benchmarks-2021/ (biaisé, à citer comme tel)

## À compléter (axes restants — recherche à faire avant les leçons concernées)
- [ ] mTLS : RFC 8446 (TLS 1.3), explication du handshake mutuel.
- [ ] Supply chain : SLSA framework (slsa.dev), Sigstore/Cosign (docs), OWASP
      (SAST/DAST/SCA definitions), Semgrep/Trivy/ZAP/Checkov docs officielles.
- [ ] GitOps : OpenGitOps principles (opengitops.dev), doc ArgoCD.
- [ ] Infra : doc Terraform vs Ansible (HashiCorp/Red Hat), k3s (Rancher), Cloudflare Tunnel.

## Communautés (acquérir la sagesse — tester ses connaissances en vrai)

- **CNCF Slack** — canaux `#spiffe`, `#spire`, `#istio` : les mainteneurs y répondent.
  https://slack.cncf.io/
- **r/devops**, **r/kubernetes** : retours d'expérience, débats d'archi.
- *(à proposer selon les besoins de Sylvain)*
