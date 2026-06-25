# 0003 — Pause active « mode jury » sur l'Acte 1 complet (leçons 1→5)

**Date :** 2026-06-24
**Type :** évaluation par récupération, difficulté montée (questions combinant plusieurs leçons)

## Résultat : ~3,5 / 4 sur le fond

| Q | Sujet | Verdict |
|---|-------|---------|
| 1 | Le voyage complet d'une requête Orders→Payments | ✅ handshake parfait ; a zappé le DÉBUT (demande SVID via Workload API) et la fin (AuthZ avant le 200) |
| 2 | Attestation en 2 phases | ✅✅ **TROU COMBLÉ** — a redonné node PUIS workload dans le bon ordre (échec précédent rattrapé) |
| 3 | Piège « mTLS réussi = sécurité OK ? » | ✅✅ parfait, n'est pas tombé dedans |
| 4 | Code vs plateforme (v2 plus fort ?) | 🟠 **2 erreurs** |
| 5 | Défense en profondeur (Gateway compromise) | 🟡 bonne intuition, justification floue |

## Acquis confirmés (storage strength)

- **Attestation 2 phases** : désormais solide (était le trou du jury précédent → comblé).
- Mécanique du handshake mTLS (challenge crypto, possession clé privée) : maîtrisée.
- Distinction AuthN/AuthZ : réflexe acquis (n'est pas tombé dans le piège Q3).

## NOUVEAUX trous à re-tester (espacement)

**Q4 — le piège « v2 est-elle plus sécurisée ? » (PRIORITAIRE).** Deux erreurs :
1. A dit que v2 « permet à des containers sans mTLS de le supporter » → FAUX : Analytics
   (Node) avait DÉJÀ du mTLS en v1 via Envoy sidecar. Le mesh généralise/automatise un mTLS
   existant, il ne l'ajoute pas à des services qui en étaient dépourvus.
2. N'a pas dit que la règle est IDENTIQUE en force ; v2 change l'ENDROIT (code→YAML Git),
   l'auditabilité, l'agilité — pas le niveau de sécurité.
   → Mantra à ancrer : « même règle, même sécurité ; ce qui change c'est où elle vit ».
   → **L'Acte 2 (Envoy en v1) va corriger l'erreur #1 directement.** Re-tester Q4 après l'Acte 2/3.

**Q5 — défense en profondeur.** Bonne intuition (« coincé comme s'il était gateway ») mais
n'articule pas le POURQUOI : attaquant a la clé privée de la Gateway, pas celle d'orders,
et ne peut pas la fabriquer → enfermé dans l'identité du pod compromis → pas de déplacement
latéral. Re-tester en lui demandant explicitement « pourquoi ne peut-il pas se faire passer
pour orders ? ».

## Style

Fond très correct, formulation à muscler. Continuer à fournir les encarts « Formule à
l'oral ». Insister sur les mantras anti-pièges (« v2 ≠ plus sécurisé, juste déplacé »).
