# 0001 — Mission posée et point de départ

**Date :** 2026-06-24
**Type :** mise en place + première leçon

## Contexte

Sylvain (étudiant DevSecOps) a invoqué `/teach` avec « je veux tout comprendre :
pourquoi telle techno, avantages/inconvénients vs concurrents, c'est quoi SPIRE, les
alternatives mTLS, la pipeline, ArgoCD, tout ».

Il a **déjà construit** le projet SPIRE/DevSecOps (code Go, pipeline CI/CD, infra
Terraform/Ansible, migration v1 SPIRE → v2 Istio, tout est sur `main` et a tourné). Il
connaît le *quoi*. Le besoin réel = blinder le *pourquoi* pour une **soutenance**.

## Décisions

- **Mission** : maîtrise de fond **calée sur la soutenance** (réponse de l'apprenant :
  « les deux, soutenance d'abord »). Chaque leçon vise 3 questions de jury : c'est quoi /
  pourquoi lui / pourquoi pas le concurrent.
- **Couverture** : il veut **tout voir** (les 4 axes : identité/mTLS, service mesh,
  pipeline/supply-chain, infra/zero-trust réseau). Ne pas sur-filtrer.
- **Ordre pédagogique** = son propre **récit en 3 actes** (Identité → Multi-langage →
  Échelle) puis transverse, parce que c'est aussi le plan de sa soutenance (deck slides/).
- **Langue** : français.

## Zone de développement proximal

Point de départ posé à **Leçon 1 (zero-trust : le problème)** — la fondation conceptuelle
qui justifie l'existence de tout le projet, et l'ouverture de sa soutenance. Volontairement
en-deçà de son niveau technique réel pour ancrer le *vocabulaire* (AuthN/AuthZ,
périmètre, lateral movement) avant d'attaquer SPIFFE/SVID en leçon 2. Il *fait* déjà du
mTLS ; ici on muscle sa capacité à l'**expliquer** et le **défendre**.

## Atelier mis en place

- `MISSION.md`, `NOTES.md`, `RESOURCES.md` (sources vérifiées : NIST 800-207, Forrester,
  BeyondCorp, spiffe.io, CNCF).
- `assets/teach.css` (style Tufte partagé) + `assets/quiz.js` (composant quiz réutilisable).
- `reference/glossaire.html` (autorité vocabulaire, ~30 termes, exemples tirés de SON code).
- `lessons/0001-zero-trust-le-probleme.html` (livrée).

## 3 corrections d'exactitude à retenir (issues de la recherche)

1. « never trust, always verify » = **Forrester/Kindervag (2010)**, pas NIST (qui dit « no
   implicit trust »). Wording d'origine : « verify and never trust ».
2. SPIFFE/SPIRE sont **Graduated** CNCF depuis **août 2022** — surtout pas dire « incubating ».
3. Expansion officielle : « Secure Production Identity Framework **for** Everyone » (for minuscule).

## Suite

Leçon 2 — SPIFFE & SVID (l'identité de workload, vs token/API key). Puis 3 (SPIRE/
attestation), 4 (mTLS/handshake), 5 (la policy `/pay` = authentifié ≠ autorisé).
