<!--
  FlowDiagram — schéma de flux animé piloté par les clics Slidev.

  Props :
    nodes : [{ id, label, sub?, x, y, w?, h?, variant?, at? }]
            x/y/w/h en unités SVG (viewBox 0..vbW × 0..vbH).
            variant : 'svc' | 'accent' | 'ink' | 'ok' | 'danger'
            at : index de clic à partir duquel le nœud apparaît (0 = visible d'emblée)
    edges : [{ from, to, label?, variant?, at?, dashed?, curve? }]
            variant : 'default' | 'ok' | 'danger' | 'soft'
    vb    : [largeur, hauteur] du viewBox (def [1000, 560])
    auto  : si vrai, le schéma se construit tout seul à l'arrivée sur la slide
            (cascade en stagger basée sur l'ordre des `at:`), sans consommer de clic.

  Sans `auto`, le composant lit $clicks (injecté par Slidev) pour révéler progressivement.
-->
<template>
  <svg :viewBox="`0 0 ${vbW} ${vbH}`" :width="vbW" :height="vbH"
       class="flow" preserveAspectRatio="xMidYMid meet"
       style="width:100%; height:100%; display:block">
    <defs>
      <marker v-for="m in markers" :key="m.id" :id="m.id" markerWidth="9" markerHeight="9"
              refX="7" refY="4.5" orient="auto" markerUnits="userSpaceOnUse">
        <path d="M0,0 L9,4.5 L0,9 Z" :fill="m.color" />
      </marker>
    </defs>

    <!-- arêtes -->
    <g>
      <g v-for="(e, i) in resolvedEdges" :key="'e'+i"
         class="edge" :class="{ shown: isShown(e.at) }"
         :style="auto ? { 'transition-delay': delayFor(e.at) } : null">
        <path :d="e.d" fill="none"
              :stroke="e.color" :stroke-width="e.width"
              :marker-end="`url(#${e.marker})`"
              :class="['edge-path', e.dashed ? 'is-dashed' : 'is-solid']"
              :style="{ '--len': e.len, ...(auto ? { 'animation-delay': delayFor(e.at) } : {}) }" />
        <g v-if="e.label">
          <rect :x="e.lx - e.lw/2" :y="e.ly - 15" :width="e.lw" height="26" rx="6"
                fill="#FBEDE7" />
          <text :x="e.lx" :y="e.ly + 3" text-anchor="middle"
                :fill="e.labelColor" class="t-edge">{{ e.label }}</text>
        </g>
      </g>
    </g>

    <!-- nœuds -->
    <g>
      <g v-for="n in resolvedNodes" :key="n.id"
         class="node" :class="{ shown: isShown(n.at) }"
         :style="auto ? { 'transition-delay': delayFor(n.at) } : null">
        <rect :x="n.x" :y="n.y" :width="n.w" :height="n.h" :rx="n.rx"
              :fill="n.fill" :stroke="n.stroke" :stroke-width="n.sw" />
        <text :x="n.x + n.w/2" :y="n.sub ? n.y + n.h/2 - 7 : n.y + n.h/2 + 6"
              text-anchor="middle" :fill="n.text" class="t-label">{{ n.label }}</text>
        <text v-if="n.sub" :x="n.x + n.w/2" :y="n.y + n.h/2 + 15"
              text-anchor="middle" :fill="n.subColor" class="t-sub">{{ n.sub }}</text>
      </g>
    </g>
  </svg>
</template>

<script setup>
import { computed } from 'vue'
import { useSlideContext } from '@slidev/client'

const props = defineProps({
  nodes: { type: Array, required: true },
  edges: { type: Array, default: () => [] },
  vb: { type: Array, default: () => [1000, 560] },
  // cascade automatique à l'arrivée sur la slide, sans consommer de clic
  auto: { type: Boolean, default: false },
  // délai (ms) entre deux paliers `at:` successifs en mode auto
  step: { type: Number, default: 300 },
})

const { $clicks } = useSlideContext()
const clicks = computed(() => $clicks?.value ?? 0)

// en mode auto : tout est révélé d'emblée, le décalage vient du transition-delay.
// sinon : on respecte le seuil de clic `at:`.
function isShown(at) {
  return props.auto ? true : clicks.value >= (at ?? 0)
}
function delayFor(at) {
  return `${(at ?? 0) * props.step}ms`
}

const vbW = computed(() => props.vb[0])
const vbH = computed(() => props.vb[1])

const PALETTE = {
  svc:    { fill: '#EEF1F6', stroke: '#1C2333', sw: 1.5, text: '#1C2333', subColor: '#5B6678' },
  accent: { fill: '#FBEDE7', stroke: '#D97757', sw: 2.5, text: '#1C2333', subColor: '#B75A3D' },
  ink:    { fill: '#1C2333', stroke: '#1C2333', sw: 1.5, text: '#FFFFFF', subColor: '#9AA4B2' },
  ok:     { fill: '#E7F1EC', stroke: '#2E7D5B', sw: 2,   text: '#1C2333', subColor: '#2E7D5B' },
  danger: { fill: '#FBE9E7', stroke: '#C0392B', sw: 2,   text: '#1C2333', subColor: '#C0392B' },
}

const EDGE = {
  default: { color: '#5B6678', width: 2,   labelColor: '#1C2333' },
  soft:    { color: '#9AA4B2', width: 1.8, labelColor: '#5B6678' },
  ok:      { color: '#2E7D5B', width: 3,   labelColor: '#2E7D5B' },
  danger:  { color: '#C0392B', width: 2.5, labelColor: '#C0392B' },
}

const markers = [
  { id: 'mk-default', color: '#5B6678' },
  { id: 'mk-soft', color: '#9AA4B2' },
  { id: 'mk-ok', color: '#2E7D5B' },
  { id: 'mk-danger', color: '#C0392B' },
]

const resolvedNodes = computed(() => props.nodes.map(n => {
  const p = PALETTE[n.variant || 'svc']
  const w = n.w ?? 170, h = n.h ?? 72
  return { ...n, w, h, rx: n.rx ?? 10, ...p }
}))

function nodeById(id) {
  const n = resolvedNodes.value.find(x => x.id === id)
  return n
}

const resolvedEdges = computed(() => props.edges.map(e => {
  const a = nodeById(e.from), b = nodeById(e.to)
  if (!a || !b) return null
  const ax = a.x + a.w / 2, ay = a.y + a.h / 2
  const bx = b.x + b.w / 2, by = b.y + b.h / 2
  // point de sortie/entrée sur le bord le plus proche (horizontal priorisé)
  const sx = bx > ax ? a.x + a.w : (bx < ax ? a.x : ax)
  const ex = bx > ax ? b.x : (bx < ax ? b.x + b.w : bx)
  const sy = ay, ey = by
  const bend = e.bend || 0   // décalage vertical des points de contrôle (sépare 2 arêtes parallèles)
  const mx = (sx + ex) / 2
  const cy = (sy + ey) / 2 + bend
  let d
  if (bend) {
    // courbe en cloche : passe par un point haut/bas pour écarter les arêtes parallèles
    d = `M ${sx} ${sy} Q ${mx} ${cy}, ${ex} ${ey}`
  } else if (e.curve) {
    d = `M ${sx} ${sy} C ${mx} ${sy}, ${mx} ${ey}, ${ex} ${ey}`
  } else if (Math.abs(sy - ey) < 2) {
    d = `M ${sx} ${sy} L ${ex} ${ey}`
  } else {
    d = `M ${sx} ${sy} C ${mx} ${sy}, ${mx} ${ey}, ${ex} ${ey}`
  }
  const st = EDGE[e.variant || 'default']
  const len = Math.hypot(ex - sx, ey - sy) + Math.abs(ey - sy) + Math.abs(bend)
  const lw = e.label ? e.label.length * 10 + 16 : 0
  // label : positionné sur le sommet de la cloche si bend, sinon au milieu
  const ly = bend ? (sy + ey) / 2 + bend * 0.62 : (sy + ey) / 2
  return {
    d, len,
    color: st.color, width: st.width, labelColor: st.labelColor,
    marker: 'mk-' + (e.variant || 'default'),
    dashed: e.dashed, at: e.at, label: e.label,
    lx: (sx + ex) / 2, ly, lw,
  }
}).filter(Boolean))
</script>

<style scoped>
.flow { display: block; }

/* font-size en unités viewBox — !important pour battre les règles globales Slidev/UnoCSS */
.flow .t-label { font-size: 20px !important; font-weight: 700; }
.flow .t-sub   { font-size: 14px !important; font-weight: 400; }
.flow .t-edge  { font-size: 17px !important; font-weight: 600; }

.node {
  opacity: 0;
  transform: translateY(8px) scale(0.97);
  transform-box: fill-box;
  transform-origin: center;
  transition: opacity .45s ease, transform .45s cubic-bezier(.2,.8,.2,1);
}
.node.shown { opacity: 1; transform: none; }

.edge { opacity: 0; transition: opacity .4s ease; }
.edge.shown { opacity: 1; }

/* tracé progressif des arêtes pleines */
.edge-path.is-solid {
  stroke-dashoffset: var(--len);
  stroke-dasharray: var(--len);
}
.edge.shown .edge-path.is-solid {
  animation: draw .6s ease forwards;
}
@keyframes draw {
  to { stroke-dashoffset: 0; }
}
/* arêtes en pointillés : dash fixe, fondu simple via .edge */
.edge-path.is-dashed { stroke-dasharray: 8 8; }
</style>
