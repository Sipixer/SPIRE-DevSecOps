# Notes de l'enseignant

## Préférences de l'apprenant

- **Langue** : français (toutes les leçons, le glossaire, les questions).
- **Veut TOUT voir** : ne pas sur-filtrer. Couvrir les 4 axes en entier, mais en
  respectant la mémoire de travail (1 leçon = 1 idée tenable).
- **Objectif déclaré** : « les deux, soutenance d'abord » → maîtrise de fond cadencée
  par l'oral. Donc chaque leçon vise les 3 questions de jury (c'est quoi / pourquoi
  lui / pourquoi pas le concurrent).
- A **construit** le projet lui-même. Il connaît le *quoi*. Ne pas ré-expliquer son
  propre code comme s'il ne l'avait jamais vu — partir de SON code et creuser le
  *pourquoi* derrière.

## Angle pédagogique

- Toujours ancrer dans les **vrais fichiers** du repo (chemins, extraits de SON code)
  plutôt que des exemples génériques. C'est ce qui rend la révision utile avant
  l'oral.
- Chaque techno = un **tableau comparatif vs concurrents** (le jury adore ça).
- Anticiper les **questions-pièges** du jury et les mettre en encart.

## Format

- Leçons HTML, courtes, belles (style Tufte), feuille de style partagée
  `assets/teach.css`.
- Quiz à choix : réponses de **longueur identique** (pas d'indice par la forme).
- Glossaire = `reference/glossaire.html`, autorité unique sur le vocabulaire.

## Idées de leçons (backlog, ordre = récit de soutenance)

ACTE 1 — Identité
- [x] 0001 — Zero-trust : le problème (pourquoi le réseau ne suffit pas) ✅ LIVRÉE 2026-06-24
- [x] 0002 — SPIFFE & SVID : l'identité de workload (vs token/API key)  ✅ LIVRÉE 2026-06-24
- [x] 0003 — SPIRE : attestation, trust domain, comment l'identité est *délivrée* ✅ LIVRÉE 2026-06-24
- [x] 0004 — mTLS : comment l'identité voyage et est vérifiée (le handshake) ✅ LIVRÉE 2026-06-24
- [x] 0005 — Authentifié ≠ autorisé : la policy `/pay` (le fil rouge) ✅ LIVRÉE 2026-06-24 — CLÔT L'ACTE 1

ACTE 2 — Multi-langage  ← PROCHAIN ACTE

- [x] 0006 — Le problème du multi-langage : Envoy en sidecar (Analytics/Node) ✅ LIVRÉE 2026-06-24
- [x] 0007 — Sidecar pattern : plan de données vs plan de contrôle ✅ LIVRÉE 2026-06-24

ACTE 3 — Échelle  ← ACTE EN COURS
- [x] 0008 — Service mesh : ce qu'il remplace dans le code (v1→v2, le diff) ✅ LIVRÉE 2026-06-24 (chiffres réels : +321/-1300)
- [x] 0009 — Istio vs Linkerd vs Cilium : le choix assumé (+ ambient mesh) ✅ LIVRÉE 2026-06-24
- [x] 0010 — SPIRE vs le mesh : qui émet l'identité en v2 ? (piège de jury) ✅ LIVRÉE 2026-06-24 — CLÔT L'ACTE 3

TRANSVERSE — Supply chain & pipeline  ← ACTE EN COURS

- [x] 0011 — SAST/DAST/SCA/IaC-scan : la taxonomie (qui trouve quoi, quand) ✅ LIVRÉE 2026-06-24
- [x] 0012 — Semgrep / Trivy / ZAP / Checkov : outil par outil, vs concurrents ✅ LIVRÉE 2026-06-24
- [x] 0013 — Signature & provenance : Cosign keyless, SLSA, pourquoi signer ✅ LIVRÉE 2026-06-24
- [x] 0014 — GitOps & ArgoCD : la branche deploy/prod, pull vs push, rollback ✅ LIVRÉE 2026-06-24

TRANSVERSE — Infra & zero-trust réseau
- [x] 0015 — Terraform vs Ansible : provisionnement vs configuration (qui fait quoi) ✅ LIVRÉE 2026-06-24
- [x] 0016 — k3s vs k8s : pourquoi la distrib légère ✅ LIVRÉE 2026-06-24
- [x] 0017 — SSH juste-à-temps : le firewall dynamique, break-glass, if:always() ✅ LIVRÉE 2026-06-24
- [x] 0018 — Tunnel Cloudflare + Terraform Cloud : zéro port entrant, secrets hors GitHub ✅ LIVRÉE 2026-06-24

🎓 **LES 18 LEÇONS SONT LIVRÉES (2026-06-24).** Sylvain va tout relire proprement, puis
on fera un QUIZ FINAL de synthèse (mode jury, transversal sur les 3 actes + transverse).

Trous à RE-TESTER au quiz final (espacement) :
- Q4 (jury acte 1) : « v2 ≠ plus sécurisée, juste déplacée » + Analytics avait DÉJÀ du mTLS (Envoy). [[learning-records/0003]]
- Q5 (jury acte 1) : défense en profondeur — articuler « enfermé dans l'identité du pod compromis (clé privée infalsifiable) ».
- Node attestation (phase 1) : a été comblé une fois, re-confirmer la rétention.

(Backlog vivant — réordonner selon les questions de Sylvain.)
