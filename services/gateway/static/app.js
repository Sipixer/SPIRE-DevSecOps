// Front de démo. Déclenche les appels via la Gateway et affiche, pour chacun,
// l'identité de l'appelant (lue côté service dans le certificat), la décision,
// l'horodatage et le chemin d'appel. Journaux des services en auto-refresh.

const SA = id => (id || '').replace('spiffe://example.org/ns/shop/sa/', '') || '(inconnu)';

// Chemin d'appel par scénario, pour l'afficher en chips.
const PATHS = {
  order:     ['gateway', 'orders', 'payments'],
  products:  ['gateway', 'catalog'],
  forbidden: ['gateway', 'payments'],
};

let filter = 'all';
let autoTimer = null;

function now() {
  const d = new Date();
  const p = n => String(n).padStart(2, '0');
  return `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`;
}

async function run(btn, url, method, label) {
  const scn = btn.dataset.scn;
  btn.disabled = true;

  let data, allowed;
  try {
    const r = await fetch(url, { method });
    data = await r.json();
    allowed = r.ok && data.allowed !== false && data.error === undefined &&
              (data.payments_status === undefined || data.payments_status === 200);
  } catch (e) {
    data = { error: String(e) }; allowed = false;
  } finally {
    btn.disabled = false;
  }

  addEntry({ scn, label, allowed, data, time: now() });
  flashTopology(scn, allowed);
  loadCalls();
}

function addEntry({ scn, label, allowed, data, time }) {
  document.getElementById('history-empty').style.display = 'none';

  const caller = data.caller ? SA(data.caller) : null;
  const path = (PATHS[scn] || []).map((s, i, arr) =>
    `<span class="chip">${s}</span>` + (i < arr.length - 1 ? '<span class="a">→</span>' : '')
  ).join('');

  const el = document.createElement('div');
  el.className = `entry ${allowed ? 'ok' : 'ko'}`;
  el.dataset.verdict = allowed ? 'ok' : 'ko';
  el.innerHTML = `
    <div class="row1">
      <span class="label">${label}</span>
      <span class="time">${time}</span>
    </div>
    <span class="verdict">${allowed ? '● autorisé' : '● refusé'}</span>
    <div class="path">${path}</div>
    ${caller ? `<div class="id">identité vue par le service : <code>${caller}</code></div>` : ''}
    <details><summary>réponse brute</summary><pre>${esc(JSON.stringify(data, null, 2))}</pre></details>`;

  const list = document.getElementById('history');
  list.prepend(el);
  applyFilter();
}

function setFilter(f, btn) {
  filter = f;
  document.querySelectorAll('.filter button').forEach(b => b.classList.toggle('active', b === btn));
  applyFilter();
}

function applyFilter() {
  document.querySelectorAll('#history .entry').forEach(e => {
    e.style.display = (filter === 'all' || e.dataset.verdict === filter) ? '' : 'none';
  });
}

function clearHistory() {
  document.getElementById('history').innerHTML = '';
  document.getElementById('history-empty').style.display = '';
}

function flashTopology(scn, allowed) {
  document.querySelectorAll('.node').forEach(n => n.classList.remove('flash-ok', 'flash-ko'));
  const cls = allowed ? 'flash-ok' : 'flash-ko';
  (PATHS[scn] || []).forEach(id =>
    document.querySelector(`.node[data-id="${id}"]`)?.classList.add(cls));
  setTimeout(() => document.querySelectorAll('.node')
    .forEach(n => n.classList.remove('flash-ok', 'flash-ko')), 2200);
}

async function loadCalls() {
  const box = document.getElementById('calls');
  let data;
  try {
    data = await (await fetch('/api/calls')).json();
  } catch {
    box.innerHTML = '<div class="placeholder">journaux indisponibles</div>'; return;
  }

  // Chaque service récepteur, avec les appels qu'il a reçus :
  //   <appelant>  →  <route>     <heure>  <verdict>
  const services = ['orders', 'catalog', 'payments'];
  let html = '';
  let total = 0;
  for (const svc of services) {
    const entries = Array.isArray(data[svc]) ? data[svc] : [];
    if (!entries.length) continue;
    total += entries.length;
    html += `<div class="svc"><h3>${svc}<span class="count">${entries.length}</span></h3>`;
    entries.slice(-6).reverse().forEach(e => {
      const ok = e.allowed;
      html += `<div class="logline ${ok ? '' : 'ko'}">
        <span class="who"><code>${SA(e.caller)}</code><span class="a">→</span><span class="route">${e.route || ''}</span></span>
        <span class="right">
          <span class="lt">${fmtTime(e.timestamp)}</span>
          <span class="v ${ok ? 'ok' : 'ko'}">${ok ? '✓ autorisé' : '✕ refusé'}</span>
        </span>
      </div>`;
    });
    html += `</div>`;
  }
  box.innerHTML = total
    ? html
    : '<div class="placeholder">Aucun appel pour l\'instant.<br>Lance un scénario.</div>';
}

function fmtTime(ts) {
  if (!ts) return '';
  const d = new Date(ts);
  if (isNaN(d)) return '';
  const p = n => String(n).padStart(2, '0');
  return `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`;
}

function toggleAuto(cb) {
  if (cb.checked) startAuto(); else stopAuto();
}
function startAuto() {
  stopAuto();
  autoTimer = setInterval(loadCalls, 3000);
}
function stopAuto() {
  if (autoTimer) { clearInterval(autoTimer); autoTimer = null; }
}

function esc(s) {
  return s.replace(/[&<>]/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;' }[c]));
}

// Init
loadCalls();
startAuto();
