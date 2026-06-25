<!-- Indicateur de chapitres (les 3 actes), l'actif surligné -->
<template>
  <div class="chapters">
    <div v-for="(c, i) in items" :key="i"
         class="chip" :class="{ active: i === active }">
      <span class="num">{{ i + 1 }}</span>
      <span v-if="i === active" class="lbl">{{ c }}</span>
    </div>
  </div>
</template>

<script setup>
defineProps({
  active: { type: Number, default: 0 },
  items: { type: Array, default: () => ['Identité', 'Multi-langage', 'Échelle'] },
})
</script>

<style scoped>
.chapters {
  position: absolute;
  top: 26px; right: 40px;
  display: flex; gap: 10px;
}
.chip {
  display: flex; align-items: center; gap: 7px;
  padding: 5px 13px;
  border-radius: 999px;
  background: var(--cloud);
  border: 1px solid var(--line);
  color: var(--mist);
  font-size: 0.82rem; font-weight: 600;
}
.chip.active {
  background: var(--claude);
  border-color: var(--claude);
  color: #fff;
}
.chip .num { font-weight: 800; }
</style>
