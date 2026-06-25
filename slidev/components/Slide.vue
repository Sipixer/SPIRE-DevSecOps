<!--
  Slide — coquille standard : header (kicker + titre + chapitres) en haut,
  corps CENTRÉ verticalement dans l'espace restant, footer en bas.
  Élimine le "tout colle en haut, vide en bas".

  Props :
    kicker, title, sub : en-tête
    chapter : index d'acte (0/1/2) → affiche l'indicateur de chapitres
    accent : 'ink' | 'claude' | 'ok' | 'danger' (couleur du filet + kicker)
    pad : padding latéral du corps (def true)
  Slot par défaut = le corps, centré.
-->
<template>
  <div class="slide-shell">
    <Chapters v-if="chapter !== undefined && chapter !== null" :active="chapter" />

    <header class="slide-head">
      <div v-if="kicker" class="kicker" :style="kickerStyle">{{ kicker }}</div>
      <h1 v-if="title" class="accent-title slide-title" :style="accentStyle">{{ title }}</h1>
      <div v-if="sub" class="slide-sub">{{ sub }}</div>
    </header>

    <main class="slide-body" :class="{ 'no-pad': pad === false }">
      <div class="slide-body-inner">
        <slot />
      </div>
    </main>

    <footer class="slide-foot">
      <span>SPIRE Zero-Trust</span>
      <span>Sylvain Rougié</span>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({
  kicker: String,
  title: String,
  sub: String,
  chapter: { type: Number, default: undefined },
  accent: { type: String, default: 'claude' },
  pad: { type: Boolean, default: true },
})
const COLORS = { ink: 'var(--ink)', claude: 'var(--claude)', ok: 'var(--ok)', danger: 'var(--danger)' }
const accentStyle = computed(() => ({ '--accent': COLORS[props.accent] || 'var(--claude)' }))
const kickerStyle = computed(() => ({ color: COLORS[props.accent] || 'var(--claude)' }))
</script>

<style scoped>
.slide-shell {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  padding: 2.4rem 3rem 2.2rem;
}
.slide-head { flex: none; }
.slide-title {
  font-size: 2.3rem;
  line-height: 1.1;
  margin: 0.15rem 0 0;
}
.slide-title::before { background: var(--accent); }
.slide-sub {
  color: var(--slate);
  font-size: 1.05rem;
  margin-top: 0.5rem;
}
/* le corps prend tout l'espace restant et centre son contenu */
.slide-body {
  flex: 1 1 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 0;
  padding: 1.2rem 0 0.4rem;
}
.slide-body.no-pad { padding: 0.6rem 0 0; }
.slide-body-inner {
  width: 100%;
}
.slide-foot {
  flex: none;
  display: flex;
  justify-content: space-between;
  font-size: 0.72rem;
  color: var(--mist);
  border-top: 1px solid var(--line);
  padding-top: 8px;
}
</style>
