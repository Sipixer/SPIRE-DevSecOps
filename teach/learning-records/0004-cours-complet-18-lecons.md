# 0004 — Cours complet livré : les 18 leçons

**Date :** 2026-06-24
**Type :** jalon — fin de la phase de production des leçons

## Contexte

Sylvain a demandé « fait toutes les leçons, je vais tout lire proprement et après on verra
pour un quiz ». J'ai donc produit les 7 leçons restantes (12→18) d'un coup, ancrées dans
son vrai code (ci.yml, cloudflare.tf, l'Application ArgoCD, le firewall dynamique TF).

## État final

**18/18 leçons livrées.** Plus : `index.html` (sommaire navigable), `glossaire.html`
(~30 termes), `RESOURCES.md` (sources vérifiées), `assets/teach.css` + `quiz.js`.

Structure = le récit de soutenance :
- Acte 1 (1-5) : zero-trust → SPIFFE/SVID → SPIRE/attestation → mTLS → AuthZ≠AuthN
- Acte 2 (6-7) : multi-langage/Envoy → data/control plane
- Acte 3 (8-10) : ce que le mesh remplace (+321/-1300) → Istio vs Linkerd vs Cilium → SPIRE vs mesh
- Transverse pipeline (11-14) : taxonomie → outils → Cosign/SLSA → GitOps/ArgoCD
- Transverse infra (15-18) : Terraform vs Ansible → k3s vs k8s → SSH JIT → Cloudflare/TF Cloud

## Faits vérifiés intégrés (corrections clés)

- NIST « no implicit trust » vs Forrester « never trust… » (attribution).
- SPIFFE/SPIRE Graduated CNCF août 2022.
- Attestation 2 phases (k8s_psat + TokenReview / workload via kubelet → selectors).
- ClusterSPIFFEID default template = exactement le SPIFFE ID du projet.
- Istio Ambient GA 1.24 (nov 2024) ; Istio Graduated juil 2023, Linkerd juil 2021, Cilium oct 2023.
- Friction Linkerd : stables payantes >50 salariés (mai 2024).
- Benchmarks meshes TOUS vendor-biaisés (à signaler, pas citer comme neutre).
- **SLSA = v1.2** (v1.0 retirée — correction importante).
- 4 principes OpenGitOps ; pull>push (gitops.tech) ; selfHeal/prune (argo docs) ; éviter :latest (k8s docs).

## Suite prévue

Sylvain relit tout. Puis **QUIZ FINAL mode jury** — transversal, doit couvrir les 3 actes
+ transverse, et RE-TESTER les trous connus :
1. « v2 ≠ plus sécurisée, juste déplacée » + Analytics avait déjà du mTLS (Envoy) [[0003-jury-acte-1]].
2. Défense en profondeur articulée (enfermé dans l'identité du pod compromis).
3. Node attestation (phase 1) — confirmer rétention durable.

Idée : proposer aussi une fiche de révision A4 imprimable (reference/) condensant les
formules-chocs et les questions-pièges, utile la veille de l'oral.
