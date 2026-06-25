# 0002 — Pause active « mode jury » sur les leçons 1→3

**Date :** 2026-06-24
**Type :** évaluation par récupération (retrieval practice)

## Format

5 questions orales type soutenance, réponses de mémoire dans le chat, correction immédiate.

## Résultat : ~4,5 / 5

| Q | Sujet | Verdict | Note |
|---|-------|---------|------|
| 1 | Définition zero-trust + origine | ✅ bon fond | manque « emplacement réseau » + « déplacement latéral » comme justification |
| 2 | Firewall/réseau privé vs SPIRE/mTLS | ✅ très bon | a nommé identité-du-processus + L3/path ; manque « falsifiable » + « complémentaires » |
| 3 | Vocabulaire SPIFFE/SPIFFE ID/SVID/SPIRE | ✅✅ parfait | analogie passeport maîtrisée, à garder telle quelle |
| 4 | Attestation (comment SPIRE est sûr) | 🟡 PARTIEL | **n'a retenu que la workload attestation ; a OUBLIÉ la node attestation (phase 1)** |
| 5 | Piège d'attribution « never trust… » | ✅ aucun piège | a repéré que c'est Forrester (pas NIST) ET que le wording diffère |

## Trou identifié (à re-tester — espacement)

**L'attestation en DEUX phases.** Sylvain a solidement la phase 2 (workload : kubelet
observe PID→pod→ServiceAccount) mais a oublié la **phase 1 (node attestation, `k8s_psat`
+ TokenReview)**. Et la phrase-clé « on ne livre pas un secret, on vérifie des faits / le
workload ne se déclare pas, il est observé » n'est pas encore réflexe.

→ **À refaire** dans une prochaine pause active (dans 2-3 leçons, pas tout de suite =
spacing). La leçon 4 (mTLS) réutilise l'attestation, ce qui va déjà la consolider en
passant.

## Forces confirmées

- Distinction SPIFFE/SVID/SPIRE = acquise (storage strength, pas juste fluence).
- Vigilance sur l'attribution NIST vs Forrester = acquise.
- Sait que firewall et mTLS opèrent à des niveaux différents (L3 vs identité/L7).

## Style à coacher

Bon réflexe de fond ; travailler la **formulation « qui fait mouche »** (phrases-chocs
prêtes à dire). Les encarts « Formule à l'oral » des leçons servent à ça.
