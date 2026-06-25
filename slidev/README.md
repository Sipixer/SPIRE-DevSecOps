# Deck Slidev — SPIRE Zero-Trust (version animée)

Support de soutenance **web animé** (Slidev + Vue), thème clair & corporate.
Récit en 3 actes « solution → … attendez, le problème ? », schémas animés
(nœuds + flèches qui se tracent), code Go en *magic-move*, transitions.

Tout le discours est dans les **notes speaker** (mode présentateur).

## Lancer (dev, hot-reload)

```bash
cd slidev
bun install        # première fois
bun run dev        # → http://localhost:3030/
```

- **Présentation** : http://localhost:3030/
- **Mode présentateur** (slides + notes) : http://localhost:3030/presenter/
- **Vue d'ensemble** : http://localhost:3030/overview/

Navigation : flèches ←/→ ou espace. Chaque clic révèle une étape d'animation.
`f` = plein écran, `o` = overview, `g` = aller à une slide.

## Exporter un PDF de secours

```bash
bun run export     # → slidev/spire-zerotrust.pdf
```

Le PDF capture chaque slide à son état final (filet de sécurité si le wifi de la
salle lâche, ou pour la remise du dossier).

## Structure

| Fichier | Rôle |
|---|---|
| `slides.md` | le deck (16 slides) + notes speaker |
| `styles/index.css` | thème : palette corporate, cartes, pills, footer |
| `components/FlowDiagram.vue` | schéma animé (nœuds + flèches révélés au clic) — sert pour le décor, l'identité, Envoy, le mesh, la CI/CD |
| `components/Chapters.vue` | indicateur des 3 actes (l'actif surligné) |
| `components/Footer.vue` | pied de page |
| `public/diagrams/` | PNG Mermaid (réutilisés depuis le PowerPoint, en réserve) |

## Le composant FlowDiagram

Schéma de flux piloté par les clics. Exemple :

```vue
<FlowDiagram
  :vb="[1000, 480]"
  :nodes="[{ id:'gw', label:'Gateway', x:280, y:200, at:0 }, ...]"
  :edges="[{ from:'gw', to:'pay', label:'mTLS', variant:'ok', at:3 }]"
/>
```

- `at` = à partir de quel clic le nœud / l'arête apparaît
- `variant` nœud : `svc` (def) · `accent` · `ink` · `ok` · `danger`
- `variant` arête : `default` · `soft` · `ok` · `danger` ; `dashed:true`, `curve:true`
- les arêtes pleines se **tracent** progressivement

## À compléter

- La slide **Avant / Après** annonce la démo live des dashboards (Kiali / Grafana).
  Pas de slide captures : tu fais la démo en direct.
- Relire les notes speaker et les adapter à ton phrasé.

> Le PowerPoint statique reste disponible dans `../slides/` comme alternative
> hors-ligne sans dépendance technique.
