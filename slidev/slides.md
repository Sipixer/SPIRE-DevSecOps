---
theme: default
title: Zero-Trust mTLS pour microservices
info: Soutenance DevSecOps — Sylvain Rougié
author: Sylvain Rougié
class: text-left
highlighter: shiki
lineNumbers: false
transition: slide-left
mdc: true
fonts:
  sans: Inter
  mono: JetBrains Mono
css: unocss
aspectRatio: 16/9
canvasWidth: 1280
---

<style>
@import './styles/index.css';
</style>

<div class="dark-slide side-accent absolute inset-0 flex flex-col justify-center px-20">
  <div class="kicker mb-3">Soutenance DevSecOps</div>
  <h1 class="!text-7xl !leading-none !font-extrabold mb-2">Zero-Trust mTLS</h1>
  <h1 class="!text-7xl !leading-none !font-extrabold t-white mb-6">pour microservices</h1>
  <div class="text-2xl" style="color:#C9D2E0">Une histoire en trois temps.</div>

  <div class="mt-16 pt-5" style="border-top:1px solid #3A4457; max-width:520px">
    <div class="text-xl font-bold t-white">Sylvain Rougié</div>
    <div class="text-sm t-mist mt-1">Soutenance DevSecOps</div>
  </div>

  <div class="absolute right-24 top-0 bottom-0 flex flex-col justify-center select-none"
       style="color:#2A3347; font-weight:800; font-size:12rem; line-height:.78">
    <div>1</div><div>2</div><div>3</div>
  </div>
</div>

<!--
Bonjour. Je vais vous présenter un projet DevSecOps : sécuriser les communications internes
d'une application microservices, selon le principe du zero-trust.

Plutôt qu'un catalogue de technologies, je vais le raconter comme une histoire : on part d'une
solution simple, elle marche, puis un besoin réel fait apparaître un nouveau problème — et c'est
ça qui justifie chaque brique, jusqu'au service mesh.
-->

---
layout: default
title: Le fil rouge
---

<Slide kicker="Le fil rouge" title="Trois temps, trois problèmes à résoudre">

<div class="grid grid-cols-3 gap-8">
  <div v-motion :initial="{ opacity: 0, y: 30 }" :enter="{ opacity: 1, y: 0, transition: { delay: 150, duration: 450 } }"
       class="card card-top h-64 flex flex-col" style="border-top-color:var(--ink)">
    <div class="w-16 h-16 rounded-full flex items-center justify-center text-3xl font-extrabold text-white mb-5" style="background:var(--ink)">1</div>
    <div class="kicker" style="color:var(--ink)">Identité</div>
    <div class="text-2xl font-bold mt-2">Prouver qui appelle</div>
    <div class="t-slate mt-auto text-lg">Go + SPIRE + mTLS</div>
  </div>
  <div v-motion :initial="{ opacity: 0, y: 30 }" :enter="{ opacity: 1, y: 0, transition: { delay: 350, duration: 450 } }"
       class="card card-top h-64 flex flex-col">
    <div class="w-16 h-16 rounded-full flex items-center justify-center text-3xl font-extrabold text-white mb-5" style="background:var(--claude)">2</div>
    <div class="kicker">Multi-langage</div>
    <div class="text-2xl font-bold mt-2">Un service en Node</div>
    <div class="t-slate mt-auto text-lg">+ Envoy en sidecar</div>
  </div>
  <div v-motion :initial="{ opacity: 0, y: 30 }" :enter="{ opacity: 1, y: 0, transition: { delay: 550, duration: 450 } }"
       class="card card-top h-64 flex flex-col" style="border-top-color:var(--ok)">
    <div class="w-16 h-16 rounded-full flex items-center justify-center text-3xl font-extrabold text-white mb-5" style="background:var(--ok)">3</div>
    <div class="kicker" style="color:var(--ok)">Échelle</div>
    <div class="text-2xl font-bold mt-2">Le câblage devient lourd</div>
    <div class="t-slate mt-auto text-lg">→ service mesh Istio</div>
  </div>
</div>

</Slide>

<!--
Trois actes.
Acte 1 — l'identité : sur un réseau interne, joindre un service ne prouve pas qui appelle.
SPIRE donne une identité à chaque service, mTLS en Go.
Acte 2 — le multi-langage : on ajoute un service Node. Réimplémenter mTLS en Node serait
coûteux → on délègue à un proxy, Envoy.
Acte 3 — l'échelle : on recâble la même plomberie partout, c'est lourd. C'est le problème que
résout un service mesh : Istio. Chaque acte résout un problème ET révèle le suivant.
-->

---
layout: default
title: Le décor
clicks: 4
---

<Slide kicker="Le décor" title="Une boutique en microservices, une règle à protéger">

<div class="card mb-6" style="background:var(--claude-lt); border-color:var(--claude); text-align:center; padding:.75rem">
  <span class="font-bold t-claude-dk">La règle :</span>
  <span>Gateway peut lire <b class="t-ok">/history</b>, mais seul <b>Orders</b> peut appeler <b class="t-danger">/pay</b></span>
</div>

<div style="height: 430px">
  <FlowDiagram
    :vb="[1000, 480]"
    :nodes="[
      { id:'nav', label:'Navigateur', x:40,  y:210, w:160, variant:'ink', at:0 },
      { id:'gw',  label:'Gateway',    x:280, y:210, at:0 },
      { id:'cat', label:'Catalog',    x:520, y:70,  at:1 },
      { id:'ord', label:'Orders',     x:520, y:350, at:1 },
      { id:'pay', label:'Payments',   x:800, y:210, variant:'accent', at:2 },
    ]"
    :edges="[
      { from:'nav', to:'gw', at:0 },
      { from:'gw', to:'cat', at:1, curve:true },
      { from:'gw', to:'ord', at:1, curve:true },
      { from:'gw', to:'pay', label:'/history ✓', variant:'ok', at:3, bend:-150 },
      { from:'gw', to:'pay', label:'/pay 403', variant:'danger', dashed:true, at:4, bend:90 },
      { from:'ord', to:'pay', label:'/pay ✓', variant:'ok', at:3, curve:true },
    ]"
  />
</div>

</Slide>

<!--
Le terrain. Gateway en entrée appelle Catalog et Orders ; Orders appelle Payments pour encaisser.
Le point clé, regardez Gateway → Payments : DEUX appels possibles, même source, même destination.
Gateway a le droit de lire l'historique — /history, en vert. Mais elle ne doit PAS déclencher un
paiement — /pay, en rouge, 403. Seul Orders passe sur /pay.
C'est ÇA le cœur du problème : on ne peut pas trancher au niveau réseau. Une NetworkPolicy dirait
juste « Gateway peut joindre Payments : oui » — elle ne voit ni la route, ni l'identité applicative.
Il faut décider par route ET par identité. C'est exactement ce que le mTLS + SPIFFE va permettre.
-->

---
layout: default
title: La stack en un coup d'œil
clicks: 3
---

<Slide kicker="Le décor" title="La stack, en un coup d'œil" accent="claude"
       sub="Trois plans empilés — c'est la carte à garder en tête pour la suite">

<div class="grid grid-cols-3 gap-6">

  <div v-click="1" class="card h-72 flex flex-col" style="border-top:5px solid var(--ink)">
    <div class="kicker mb-3" style="color:var(--ink)">① Les workloads</div>
    <div class="font-bold text-lg mb-2">5 services</div>
    <ul class="bullets space-y-1 text-base t-slate">
      <li>Gateway · Orders · Catalog · Payments <span class="t-mist">(Go)</span></li>
      <li>Analytics <span class="t-mist">(Node.js)</span></li>
      <li>Loadgen <span class="t-mist">(trafic continu)</span></li>
    </ul>
    <div class="t-slate mt-auto text-sm">La règle <b class="t-danger">/pay → Orders seul</b> vit ici.</div>
  </div>

  <div v-click="2" class="card card-top h-72 flex flex-col">
    <div class="kicker mb-3">② La sécurité runtime</div>
    <div class="font-bold text-lg mb-2">Identité + mTLS</div>
    <ul class="bullets space-y-1 text-base t-slate">
      <li><b>v1</b> — SPIRE + mTLS dans le code Go</li>
      <li><b>v2</b> — service mesh <b class="t-claude-dk">Istio</b></li>
      <li>Observé dans <b>Kiali</b></li>
    </ul>
    <div class="t-slate mt-auto text-sm">C'est l'histoire en 3 actes qui vient.</div>
  </div>

  <div v-click="3" class="card h-72 flex flex-col" style="border-top:5px solid var(--ok)">
    <div class="kicker mb-3" style="color:var(--ok)">③ La chaîne & l'infra</div>
    <div class="font-bold text-lg mb-2">CI/CD + GitOps</div>
    <ul class="bullets space-y-1 text-base t-slate">
      <li>Checkov · Semgrep · Trivy</li>
      <li>Cosign → GHCR → ArgoCD</li>
      <li>Terraform + Ansible · k3s · SSH JIT</li>
    </ul>
    <div class="t-slate mt-auto text-sm">Tout est <b>reproductible from scratch</b>.</div>
  </div>

</div>

</Slide>

<!--
Avant de plonger, la carte d'ensemble — trois plans empilés, à garder en tête.
Un : les workloads. Cinq services applicatifs — quatre en Go, un en Node, plus un générateur de
trafic. C'est ici que vit la règle métier : seul Orders peut déclencher un paiement.
Deux : la sécurité au runtime. C'est le cœur de l'exposé — comment on donne une identité prouvée à
chaque service et comment on impose le mTLS. Ça a évolué : v1 dans le code avec SPIRE, v2 dans la
plateforme avec Istio, le tout observable dans Kiali.
Trois : la chaîne de livraison et l'infra. La CI scanne et signe, ArgoCD déploie en GitOps, et tout
le serveur est monté par Terraform et Ansible — reproductible de zéro, sans étape manuelle.
Les actes qui suivent zooment sur le plan deux. Gardez cette carte : on y revient à la fin.
-->

---
layout: default
title: Acte 1 — Identité
clicks: 1
---

<Slide :chapter="0" kicker="Acte 1 · Identité" title="Donner une identité à chaque service"
       sub="SPIRE émet un certificat par workload · les services parlent en mTLS" accent="ink">

<div style="height: 330px">
  <FlowDiagram
    auto
    :vb="[1000, 360]"
    :nodes="[
      { id:'spire', label:'SPIRE', x:40,  y:140, w:150, variant:'accent', at:0 },
      { id:'gw',  label:'Gateway',  x:340, y:140, at:1 },
      { id:'ord', label:'Orders',   x:580, y:140, at:1 },
      { id:'pay', label:'Payments', x:820, y:140, at:1 },
    ]"
    :edges="[
      { from:'spire', to:'gw', label:'SVID', variant:'soft', dashed:true, at:2 },
      { from:'gw', to:'ord', label:'mTLS', variant:'ok', at:3 },
      { from:'ord', to:'pay', label:'mTLS', variant:'ok', at:3 },
    ]"
  />
</div>

<div v-click="1" class="flex items-center gap-3 justify-center mt-6">
  <div class="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold text-lg" style="background:var(--ok)">✓</div>
  <span class="pill text-base" style="background:var(--ok-lt); color:var(--ok); padding:.3em 1em">Payments sait, de façon prouvée, que c'est Orders</span>
</div>

</Slide>

<!--
Acte 1, la solution. SPIRE donne à chaque service une identité cryptographique — le SVID,
un certificat X.509 court — dérivée de son ServiceAccount Kubernetes.
Les services Go récupèrent ce certificat et font du mTLS : chiffré ET authentification mutuelle.
Résultat : Payments sait, de façon prouvée, que c'est bien Orders qui appelle — par la clé, pas
par un header ni une IP. Le réseau n'a plus besoin d'être de confiance.
Et ça soulève une vraie question, qu'on traite tout de suite : pourquoi ce mTLS, et pas juste
une NetworkPolicy ?
-->

---
layout: default
title: Pourquoi pas une NetworkPolicy
clicks: 4
---

<Slide :chapter="0" kicker="Acte 1 · La bonne couche"
       title="Pourquoi mTLS et pas une NetworkPolicy ?" accent="ink">

<div class="np-table">
  <div class="np-row np-head">
    <div>Couche</div><div>Outil</div><div>Identité basée sur</div><div>Granularité</div><div>« /pay vs /history » ?</div>
  </div>
  <div v-click="1" class="np-row">
    <div><span class="np-tag" style="background:var(--cloud); color:var(--ink)">L3 / L4</span></div>
    <div><b>NetworkPolicy</b> native</div>
    <div>label → pod <span class="t-mist">(contextuel)</span></div>
    <div>pod + port</div>
    <div class="t-danger font-bold">inexprimable</div>
  </div>
  <div v-click="2" class="np-row">
    <div><span class="np-tag" style="background:var(--cloud); color:var(--ink)">L7 réseau</span></div>
    <div>Cilium <span class="t-mist">(CNNP, eBPF)</span></div>
    <div>label → pod, + HTTP</div>
    <div>route / méthode</div>
    <div class="t-slate">oui, mais identité <b>réseau</b></div>
  </div>
  <div v-click="3" class="np-row np-hot">
    <div><span class="np-tag" style="background:var(--ok); color:#fff">L7 appli</span></div>
    <div><b>mTLS + SPIFFE</b> <span class="t-mist">(mon choix)</span></div>
    <div>clé privée / certificat</div>
    <div>route / méthode</div>
    <div class="t-ok font-bold">oui, par identité prouvée</div>
  </div>
</div>

<div v-click="4" class="text-center text-2xl mt-8">
  La NetworkPolicy prouve <b>d'où vient un paquet</b>.<br/>
  Le mTLS prouve <b class="t-claude-dk">qui l'a signé</b>.
</div>

</Slide>

<!--
On traite la question tout de suite, parce qu'un jury va la poser.

Première ligne — la NetworkPolicy native de Kubernetes est L3/L4 : c'est la spec elle-même. Elle
sélectionne des pods par label, et autorise du trafic par port et protocole. Les mots « route »,
« path », « méthode HTTP » n'existent nulle part dans l'objet. Donc « /history oui, /pay non » est
littéralement inexprimable. Et attention : ce n'est pas qu'elle est ratée — c'est un choix
d'architecture en couches. Le réseau fait du réseau ; le L7 est délégué à la couche au-dessus.

Deuxième ligne — anticipons l'objection. Certains CNI comme Cilium, en eBPF, font du L7 : ils
savent filtrer GET /order vs POST /pay. Mais même eux décident par identité RÉSEAU — le label, le
pod. Pas par identité cryptographique prouvée.

Troisième ligne — mon choix. Le SVID prouve la possession d'une clé privée. C'est une identité
intrinsèque et portable : elle prouve « je suis orders » même à travers un proxy, même hors du
cluster, indépendamment de la topologie réseau et du CNI.

Nuance honnête à connaître : l'IP n'est pas trivialement spoofable, un bon CNI rejette un paquet
dont l'IP source ne correspond pas au pod. La vraie faiblesse, c'est que l'IP prouve d'OÙ vient le
paquet, pas QUI l'a émis : si un attaquant compromet le pod Orders, il hérite de son IP et de tous
ses droits réseau. Le SVID a la même limite si le pod est pris — sauf qu'il est à courte durée de
vie, révocable, et prouve la possession d'une clé, pas juste une position dans le réseau.

Et ce n'est pas opposé : c'est de la défense en profondeur. NetworkPolicy réduit la surface au
niveau réseau, le mesh prouve l'identité au niveau applicatif. La phrase à retenir : la
NetworkPolicy prouve d'où vient un paquet ; le mTLS prouve qui l'a signé.
-->

---
layout: default
title: Acte 1 — Mécanisme
clicks: 1
---

<Slide :chapter="0" kicker="Acte 1 · Mécanisme" title="Comment ça marche, concrètement" accent="ink">

<div class="grid grid-cols-4 gap-6">
  <div v-motion :initial="{ opacity: 0, y: 24 }" :enter="{ opacity: 1, y: 0, transition: { delay: 100, duration: 400 } }"
       class="card text-center h-52 flex flex-col justify-center">
    <div class="text-5xl font-extrabold t-ink mb-3">①</div>
    <div class="font-bold text-lg">SPIRE émet</div>
    <div class="t-slate text-sm mt-2">un X.509-SVID</div>
  </div>
  <div v-motion :initial="{ opacity: 0, y: 24 }" :enter="{ opacity: 1, y: 0, transition: { delay: 350, duration: 400 } }"
       class="card text-center h-52 flex flex-col justify-center">
    <div class="text-5xl font-extrabold t-ink mb-3">②</div>
    <div class="font-bold text-lg">Handshake mTLS</div>
    <div class="t-slate text-sm mt-2">les deux pairs s'authentifient</div>
  </div>
  <div v-motion :initial="{ opacity: 0, y: 24 }" :enter="{ opacity: 1, y: 0, transition: { delay: 600, duration: 400 } }"
       class="card card-top text-center h-52 flex flex-col justify-center">
    <div class="text-5xl font-extrabold t-claude mb-3">③</div>
    <div class="font-bold text-lg">Identité lue</div>
    <div class="t-slate text-sm mt-2">dans le certificat</div>
  </div>
  <div v-motion :initial="{ opacity: 0, y: 24 }" :enter="{ opacity: 1, y: 0, transition: { delay: 850, duration: 400 } }"
       class="card card-top text-center h-52 flex flex-col justify-center">
    <div class="text-5xl font-extrabold t-claude mb-3">④</div>
    <div class="font-bold text-lg">Décision</div>
    <div class="t-slate text-sm mt-2">policy par route</div>
  </div>
</div>

<div v-click="1" class="text-center text-2xl mt-12">
  L'identité est <b class="t-claude-dk">prouvée pendant le handshake</b>, jamais déclarée dans un header.
</div>

</Slide>

<!--
Quatre temps : SPIRE émet le SVID. Orders et Payments font un handshake mTLS, chacun vérifie
le certificat de l'autre. Côté Payments, on lit l'identité DANS le certificat. Puis on applique
la policy de la route.
Le point clé : l'identité vient du certificat vérifié au handshake — pas d'un header, donc
non-forgeable. On le voit dans le code.
-->

---
layout: default
title: Acte 1 — L'attestation
clicks: 5
---

<Slide :chapter="0" kicker="Acte 1 · Sous le capot" title="Mais comment SPIRE sait que ce pod est Orders ?"
       sub="Le SVID n'est pas donné sur parole — il est mérité par une double attestation" accent="ink">

<div style="height: 340px">
  <FlowDiagram
    :vb="[1000, 440]"
    :nodes="[
      { id:'kube',  label:'API Kubernetes', sub:'preuve du nœud', x:360, y:20,  w:210, h:88, at:1 },
      { id:'pod',   label:'Pod Orders', sub:'ServiceAccount', x:360, y:332, w:210, h:88, at:2 },
      { id:'node',  label:'Agent SPIRE', sub:'sur le nœud k3s', x:40,  y:176, w:210, h:96, variant:'ink', at:0 },
      { id:'server',label:'Serveur SPIRE', sub:'décide + signe', x:740, y:176, w:210, h:96, variant:'accent', at:3 },
    ]"
    :edges="[
      { from:'node', to:'kube', label:'1 · qui suis-je ?', variant:'soft', dashed:true, at:1, curve:true },
      { from:'pod', to:'node', label:'2 · quel workload ?', variant:'soft', dashed:true, at:2, curve:true },
      { from:'node', to:'server', label:'3 · attestation', variant:'default', at:3, bend:-70 },
      { from:'server', to:'node', label:'4 · SVID signé', variant:'ok', at:4, bend:70 },
    ]"
  />
</div>

<div v-click="5" class="text-center text-xl mt-4">
  L'identité découle de <b>ce que le pod EST</b> (son ServiceAccount), pas d'un secret qu'on lui confie.
  <br/><span class="t-mist text-base">Pas de mot de passe à voler, pas de clé à distribuer.</span>
</div>

</Slide>

<!--
La question que le jury va poser : « d'accord, SPIRE émet un certificat — mais comment sait-il qu'il
parle bien à Orders, et pas à un imposteur ? » C'est l'attestation, et c'est en deux temps.

Un : node-attestation. Un agent SPIRE tourne sur chaque nœud. Il prouve d'ABORD l'identité du nœud
lui-même — ici via l'API Kubernetes, qui confirme « ce nœud fait bien partie du cluster ».
Deux : workload-attestation. Quand un pod demande son identité, l'agent l'interroge localement —
il lit son ServiceAccount, son namespace — via le noyau, des infos que le pod ne peut pas falsifier.
Trois : l'agent présente ces deux preuves au serveur SPIRE.
Quatre : le serveur compare aux RegistrationEntries — « le SA orders dans le ns shop = spiffe://…/orders »
— et signe le SVID seulement si ça correspond.

Le point « ah oui » : l'identité n'est jamais un secret qu'on a copié dans le pod. Elle découle de ce
que le pod EST, vérifié par la plateforme. Donc rien à voler, rien à distribuer, rien qui fuit dans
un fichier de conf. C'est ça qui rend SPIFFE supérieur à un token statique.
-->

---
layout: default
title: Acte 1 — Le code
clicks: 4
---

<Slide :chapter="0" kicker="Acte 1 · Le code" title="Le seul bout de code qui compte" accent="ink">

<div class="grid grid-cols-5 gap-8 items-center">
<div class="col-span-3">

<div class="code-tight">

````md magic-move {at:1}
```go
// payments — l'appel arrive
func handlePay(r *Request) {

}
```
```go
// l'identité vient du certificat, pas d'un header
func handlePay(r *Request) {
    id := r.TLS.PeerCertificates[0].URIs[0]
}
```
```go
// authentifié ≠ autorisé : on vérifie la policy
func handlePay(r *Request) {
    id := r.TLS.PeerCertificates[0].URIs[0]

    if !policy["/pay"].allows(id) {
        return Error(403)   // refusé
    }
}
```
````

</div>

</div>
<div class="col-span-2">
  <div v-click="4" class="card card-top">
    <div class="text-2xl font-bold mb-4">Non forgeable</div>
    <div class="t-slate leading-relaxed">
      Se faire passer pour Orders exigerait sa <b>clé privée</b>, renouvelée en continu par SPIRE.
    </div>
    <div class="t-slate leading-relaxed mt-5">
      Gateway a une identité valide… mais reste <b class="t-danger">refusée</b> sur <code>/pay</code>.
    </div>
  </div>
</div>
</div>

<style>
.code-tight :deep(.slidev-code),
.code-tight pre code { font-size: 0.82rem !important; line-height: 1.5 !important; }
</style>

</Slide>

<!--
Un seul extrait, le plus parlant. (Cliquer pour dérouler.)
On lit l'identité de l'appelant dans son certificat TLS — PeerCertificates. C'est ça qui la rend
non-forgeable : vérifiée au handshake, on ne fait que la lire.
Puis on applique la policy. Même authentifié, si pas autorisé → 403.
Authentifié ≠ autorisé : Gateway a une identité valide, mais sur /pay elle est refusée.
-->

---
layout: center
class: dark-slide side-accent
title: Attendez (1)
---

<div class="px-16">
  <div class="kicker mb-5" style="letter-spacing:.3em">Attendez…</div>
  <div class="!text-6xl !font-extrabold !leading-tight t-white">Et si on ajoute un service…</div>
  <div class="!text-6xl !font-extrabold !leading-tight t-claude">dans un autre langage ?</div>
  <div class="text-xl mt-10" style="color:#C9D2E0">
    Analytics sera écrit en Node.js. Pas de go-spiffe en Node — on réimplémente tout ?
  </div>
  <div class="absolute select-none" style="right:7rem; bottom:2rem; color:#2A3347; font-size:18rem; font-weight:800; line-height:1">?</div>
</div>

<!--
Premier rebond. Tout marche… tant que tout est en Go, grâce à go-spiffe.
Mais on ajoute un service d'analytics en Node.js. Pas d'équivalent direct : faut-il réimplémenter
toute la logique SPIFFE/mTLS dans chaque langage ? Ça ne tient pas. C'est le problème de l'acte 2.
-->

---
layout: default
title: Acte 2 — Envoy
clicks: 1
---

<Slide :chapter="1" kicker="Acte 2 · Multi-langage" title="Déléguer le mTLS à un proxy"
       sub="Analytics reste du Node pur · Envoy en sidecar porte SPIFFE/mTLS" accent="claude">

<div style="height: 340px">
  <FlowDiagram
    auto
    :vb="[1000, 460]"
    :nodes="[
      { id:'gw',    label:'Gateway · Go', sub:'mTLS dans le code', x:40,  y:120, w:210, h:90, at:0 },
      { id:'envoy', label:'Envoy',  sub:'sidecar',  x:490, y:120, w:170, h:90, variant:'accent', at:2 },
      { id:'node',  label:'Analytics', sub:'Node.js', x:770, y:120, w:190, h:90, at:1 },
      { id:'spire', label:'SPIRE', x:40,  y:320, w:170, variant:'ink', at:3 },
    ]"
    :edges="[
      { from:'envoy', to:'node', label:'HTTP local', variant:'soft', at:2 },
      { from:'gw', to:'envoy', label:'mTLS', variant:'ok', at:3 },
      { from:'spire', to:'envoy', label:'SVID (SDS)', variant:'soft', dashed:true, at:3, curve:true },
    ]"
  />
</div>

<div v-click="1" class="flex items-center gap-3 justify-center mt-2">
  <div class="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold text-lg" style="background:var(--ok)">✓</div>
  <span class="pill text-base" style="background:var(--ok-lt); color:var(--ok); padding:.3em 1em">ça marche aussi</span>
</div>

</Slide>

<!--
Acte 2, la solution. Plutôt que réécrire le mTLS en Node, on place un proxy à côté du service :
Envoy, en sidecar dans le même pod.
Envoy récupère le SVID auprès de SPIRE, termine le mTLS, vérifie l'identité du pair, puis
transmet la requête à l'app Node en HTTP local, dans le pod.
L'app Node ne connaît rien à SPIFFE. Même garantie que les services Go, sans toucher au code.
-->

---
layout: center
class: dark-slide side-accent
title: Attendez (2)
---

<div class="px-16">
  <div class="kicker mb-5" style="letter-spacing:.3em; color:var(--danger)">Attendez…</div>
  <div class="!text-6xl !font-extrabold !leading-tight t-white">Vous voyez le problème ?</div>
  <div class="!text-6xl !font-extrabold !leading-tight" style="color:var(--danger)">On recâble la même chose. Partout.</div>
  <div class="text-xl mt-10" style="color:#C9D2E0">
    Sidecar à configurer à la main, mTLS répété service par service, policy éparpillée.
  </div>
  <div class="absolute select-none" style="right:7rem; bottom:2rem; color:#2A3347; font-size:18rem; font-weight:800; line-height:1">!</div>
</div>

<!--
Deuxième rebond, le plus important. Prenons du recul.
Pour le Go, on a écrit la plomberie mTLS dans chaque service. Pour le Node, un Envoy configuré
à la main. À chaque nouveau service, on refait la même chose.
Vous voyez le problème ? Ça marche, mais ça ne passe pas à l'échelle. La sécurité est dispersée
dans le code et dans des configs manuelles. Lourd, fragile, dur à auditer. Il nous faut autre chose.
-->

---
layout: default
title: Le vrai problème
clicks: 3
---

<Slide kicker="Le vrai problème" title="Trois douleurs qui ne passent pas à l'échelle" accent="danger">

<div class="grid grid-cols-3 gap-8">
  <div v-click="1" class="card text-center h-64 flex flex-col justify-center" style="border-top:5px solid var(--danger)">
    <div class="text-6xl mb-4">⚠️</div>
    <div class="text-2xl font-bold">Répétition</div>
    <div class="t-slate mt-3">La même plomberie mTLS recâblée dans chaque service</div>
  </div>
  <div v-click="2" class="card text-center h-64 flex flex-col justify-center" style="border-top:5px solid var(--danger)">
    <div class="text-6xl mb-4">⚠️</div>
    <div class="text-2xl font-bold">Friction multi-langage</div>
    <div class="t-slate mt-3">Chaque langage = une nouvelle intégration à porter à la main</div>
  </div>
  <div v-click="3" class="card text-center h-64 flex flex-col justify-center" style="border-top:5px solid var(--danger)">
    <div class="text-6xl mb-4">⚠️</div>
    <div class="text-2xl font-bold">Policy dispersée</div>
    <div class="t-slate mt-3">Changer une règle = modifier, recompiler, redéployer</div>
  </div>
</div>

</Slide>

<!--
Trois douleurs concrètes.
Un : la répétition. La même logique mTLS réécrite à chaque service Go.
Deux : la friction multi-langage. Chaque langage demande sa propre intégration — on l'a vu avec Envoy.
Trois : la policy dispersée. Les règles vivent dans le code ; les changer = recompiler-redéployer.
La conclusion s'impose : SORTIR la sécurité du code et en faire une propriété de la plateforme.
C'est exactement ce qu'est un service mesh.
-->

---
layout: default
title: Acte 3 — Mesh
clicks: 1
---

<Slide :chapter="2" kicker="Acte 3 · Échelle" title="Le service mesh fait tout, partout"
       sub="Istio injecte un proxy par pod · mTLS STRICT automatique, policy déclarative" accent="ok">

<div style="height: 400px">
  <FlowDiagram
    auto
    :vb="[1000, 520]"
    :nodes="[
      { id:'gw',  label:'Gateway',  sub:'+ proxy', x:400, y:30,  w:180, h:84, at:0 },
      { id:'ord', label:'Orders',   sub:'+ proxy', x:110, y:220, w:180, h:84, at:1 },
      { id:'cat', label:'Catalog',  sub:'+ proxy', x:400, y:220, w:180, h:84, at:1 },
      { id:'an',  label:'Analytics',sub:'+ proxy', x:690, y:220, w:180, h:84, at:1 },
      { id:'pay', label:'Payments', sub:'+ proxy', x:400, y:415, w:180, h:84, at:2 },
    ]"
    :edges="[
      { from:'gw', to:'ord', label:'mTLS auto', variant:'ok', at:3, curve:true },
      { from:'gw', to:'cat', label:'mTLS auto', variant:'ok', at:3 },
      { from:'gw', to:'an',  label:'mTLS auto', variant:'ok', at:3, curve:true },
      { from:'ord', to:'pay', variant:'ok', at:4, curve:true },
      { from:'gw', to:'pay', label:'/pay : 403', variant:'danger', dashed:true, at:5, curve:true },
    ]"
  />
</div>

<div v-click="1" class="text-center text-xl mt-2">
  Un seul mécanisme, <b class="t-ok">identique pour Go et Node</b> — et le 403 zero-trust, en config.
</div>

</Slide>

<!--
Acte 3, la bascule. On adopte un service mesh : Istio.
Istio injecte automatiquement un proxy à côté de CHAQUE pod — Go comme Node, sans distinction.
Le mTLS devient automatique et imposé dans tout le namespace : une ligne, PeerAuthentication STRICT.
L'autorisation devient déclarative : des AuthorizationPolicy en YAML, versionnées dans Git, qui
rejouent notre règle — Gateway vers /pay, c'est 403, en rouge.
Plus de plomberie dans le code, plus d'Envoy à la main. La sécurité est passée dans la plateforme.
-->

---
layout: default
title: Acte 3 — Qui émet l'identité
clicks: 5
---

<Slide :chapter="2" kicker="Acte 3 · Sous le capot" title="Du coup, qui émet l'identité en v2 ?"
       sub="La question piège : SPIRE disparaît, mais le mécanisme reste le même" accent="ok">

<div class="grid grid-cols-2 gap-10">

  <div class="card h-80 flex flex-col" style="border-top:5px solid var(--ink)">
    <div class="kicker mb-4" style="color:var(--ink)">v1 — l'app fait le travail</div>
    <ul class="bullets space-y-3 text-lg t-ink-soft">
      <li><b>SPIRE</b> émet le SVID</li>
      <li>Le code Go <b>lit le certificat</b> du pair<br/><code class="text-sm">r.TLS.PeerCertificates</code></li>
      <li>L'app décide elle-même</li>
    </ul>
    <div class="t-slate mt-auto text-base">Identité <b>dans le code</b>.</div>
  </div>

  <div v-click="1" class="card h-80 flex flex-col" style="border-top:5px solid var(--ok); background:var(--ok-lt)">
    <div class="kicker mb-4" style="color:var(--ok)">v2 — le proxy fait le travail</div>
    <ul class="bullets space-y-3 text-lg t-ink-soft">
      <li v-click="2"><b>istiod</b> (sa CA) émet le certificat</li>
      <li v-click="3">Le proxy termine le mTLS et injecte l'identité dans <code class="text-sm">X-Forwarded-Client-Cert</code></li>
      <li v-click="4">L'<b>AuthorizationPolicy</b> décide, hors du code</li>
    </ul>
    <div class="t-slate mt-auto text-base">Identité <b>dans la plateforme</b>.</div>
  </div>

</div>

<div v-click="5" class="text-center text-lg mt-6">
  Même principe — <b>une identité prouvée au handshake</b>. Seul <b class="t-claude-dk">l'émetteur</b> change : l'app cède la place au mesh.
</div>

</Slide>

<!--
LA question de jury sur la v2 : « vous avez retiré SPIRE — du coup qui donne l'identité maintenant,
et est-ce que SPIRE servait à quelque chose ? » Réponse en miroir, gauche/droite.

À gauche, la v1. SPIRE émettait le SVID, et c'est l'APP qui faisait tout le reste : le code Go lisait
le certificat du pair dans r.TLS.PeerCertificates et décidait lui-même. L'identité vivait dans le code.

À droite, la v2. L'émetteur change : ce n'est plus SPIRE, c'est istiod — le control plane d'Istio a sa
propre autorité de certification, dérivée elle aussi du ServiceAccount Kubernetes. Le proxy sidecar
termine le mTLS à la place de l'app, puis injecte l'identité vérifiée dans un header,
X-Forwarded-Client-Cert — c'est d'ailleurs ce que mes services lisent, juste pour le journal, plus
pour la sécu. Et la décision sort du code : c'est une AuthorizationPolicy en YAML.

Le point clé à marteler : le MÉCANISME est identique — une identité cryptographique prouvée au
handshake, dérivée de ce que le workload est. SPIRE n'était pas inutile : il a prouvé le concept SPIFFE
en v1, et Istio applique exactement le même standard SPIFFE en interne. Ce qui change, c'est QUI
l'émet et QUI l'applique — l'application cède la place à la plateforme. C'est tout le sujet de l'exposé.
-->

---
layout: default
title: Avant / Après
clicks: 1
---

<Slide kicker="Avant / Après" title="La sécurité change de couche" accent="ok">

<div class="grid grid-cols-2 gap-10">
  <div class="card h-80 flex flex-col" style="border-top:5px solid var(--ink)">
    <div class="kicker mb-5" style="color:var(--ink)">Avant — dans le code</div>
    <ul class="bullets space-y-4 text-lg t-ink-soft">
      <li>mTLS écrit en Go, service par service</li>
      <li>Envoy configuré à la main pour Node</li>
      <li>policy d'autorisation compilée dans l'app</li>
      <li>à refaire pour chaque service</li>
    </ul>
  </div>
  <div v-click="1" class="card h-80 flex flex-col" style="border-top:5px solid var(--ok); background:var(--ok-lt)">
    <div class="kicker mb-5" style="color:var(--ok)">Après — dans la plateforme</div>
    <ul class="bullets space-y-4 text-lg t-ink-soft">
      <li>mTLS automatique (1 ligne de config)</li>
      <li>même proxy injecté pour Go ET Node</li>
      <li>AuthorizationPolicy déclaratives (YAML)</li>
      <li>gratuit pour chaque nouveau service</li>
    </ul>
  </div>
</div>

<div class="text-center mt-8">
  <span class="pill text-lg" style="background:var(--ink); color:#fff; padding:.55em 1.6em">La sécurité a changé de couche</span>
</div>

</Slide>

<!--
Le résultat, en un coup d'œil.
Avant, à gauche : la sécurité vivait dans le code et des configs manuelles — mTLS Go, Envoy à la
main, policy compilée, tout à refaire à chaque service.
Après, à droite : tout devient une propriété de la plateforme. mTLS en une ligne, même proxy pour
tous les langages, règles en YAML déclaratif versionné dans Git.
La sécurité n'a pas disparu : elle a changé de couche, lisible et réutilisable. Et on va la voir
tourner en direct dans la démo.
-->

---
layout: default
title: Supply chain
clicks: 7
---

<Slide kicker="Au-delà du runtime" title="La sécurité commence dans la CI/CD"
       sub="Chaque merge franchit deux temps : on porte (scans bloquants), puis on livre (GitOps)" accent="claude">

<div class="ci-flow">

  <!-- TEMPS 1 — LES PORTES -->
  <div class="ci-phase" style="--ph:var(--danger)">
    <div class="ci-tag" style="background:var(--danger)">① Porter — bloquant</div>
    <div class="ci-row">
      <div v-click="1" class="ci-box">Checkov<span>+ tflint · IaC</span></div>
      <div v-click="1" class="ci-box">Semgrep<span>SAST · code Go</span></div>
      <div v-click="1" class="ci-box">go vet + build<span>4 services</span></div>
    </div>
    <div v-click="2" class="ci-sub">tout en parallèle, puis ↓</div>
    <div v-click="2" class="ci-row">
      <div class="ci-box ci-wide">Trivy<span>scan CVE HIGH/CRITICAL — bloquant</span></div>
    </div>
    <div v-click="3" class="ci-note">Un seul échec → <b>exit ≠ 0 → le build s'arrête.</b> C'est une porte, pas un avertissement.</div>
  </div>

  <div v-click="4" class="ci-arrow">↓ images vertes</div>

  <!-- TEMPS 2 — LA LIVRAISON -->
  <div v-click="4" class="ci-phase" style="--ph:var(--ok)">
    <div class="ci-tag" style="background:var(--ok)">② Livrer — GitOps</div>
    <div class="ci-row">
      <div v-click="5" class="ci-box ci-ok">Cosign<span>signe (keyless)</span></div>
      <div v-click="5" class="ci-arrow-h">→</div>
      <div v-click="5" class="ci-box ci-ok">GHCR<span>image : sha</span></div>
      <div v-click="6" class="ci-arrow-h">→</div>
      <div v-click="6" class="ci-box ci-ok">deploy/prod<span>tag épinglé</span></div>
      <div v-click="6" class="ci-arrow-h">→</div>
      <div v-click="6" class="ci-box ci-ok">ArgoCD<span>sync → k3s</span></div>
    </div>
    <div v-click="7" class="ci-note">Pas de <code>latest</code> en prod. <b>Git fait foi</b> — un rollback = restaurer un commit.</div>
  </div>

</div>

<style>
.ci-flow { display:flex; flex-direction:column; align-items:center; gap:.5rem; width:100%; }
.ci-phase { width:100%; border:1px solid var(--line); border-left:5px solid var(--ph);
  border-radius:12px; padding:.85rem 1.2rem; background:#fff; }
.ci-tag { display:inline-block; color:#fff; font-weight:700; font-size:.8rem; letter-spacing:.05em;
  padding:.22em .9em; border-radius:999px; margin-bottom:.6rem; }
.ci-row { display:flex; gap:.7rem; align-items:stretch; justify-content:center; }
.ci-box { flex:1; background:var(--cloud); border:1px solid var(--line); border-radius:9px;
  padding:.55rem .7rem; text-align:center; font-weight:700; font-size:1.05rem; line-height:1.15;
  display:flex; flex-direction:column; justify-content:center; }
.ci-box span { display:block; font-weight:400; font-size:.78rem; color:var(--slate); margin-top:.18rem; }
.ci-box.ci-ok { background:var(--ok-lt); border-color:var(--ok); }
.ci-box.ci-wide { flex:none; min-width:60%; background:var(--claude-lt); border-color:var(--claude); }
.ci-arrow-h { display:flex; align-items:center; color:var(--mist); font-weight:800; font-size:1.3rem; }
.ci-sub { text-align:center; color:var(--mist); font-size:.8rem; margin:.4rem 0; }
.ci-arrow { color:var(--slate); font-weight:700; font-size:.92rem; }
.ci-note { margin-top:.65rem; font-size:.92rem; color:var(--ink-soft); text-align:center; }
.ci-note code { background:var(--cloud); padding:.05em .35em; border-radius:4px; }
</style>

</Slide>

<!--
Le zero-trust ne s'arrête pas au runtime : la sécurité commence dans la chaîne de livraison. Et cette
chaîne se lit en DEUX temps.

Premier temps, en haut — on PORTE le code, et c'est bloquant. Trois portes tournent en parallèle à
chaque push : Checkov plus tflint scannent l'infra Terraform, Semgrep fait le SAST sur le Go, et on
compile-et-vérifie les quatre services. Si les trois passent, on construit les images et Trivy les
scanne pour les CVE critiques. Le point à marteler : chacune est BLOQUANTE — un seul échec sort en
code non nul et arrête tout. Ce n'est pas un rapport qu'on lit plus tard, c'est une porte fermée.

Deuxième temps, en bas — une fois les images vertes, on LIVRE, et là c'est du GitOps. Cosign signe
l'image en keyless — via l'identité OIDC du runner, sans clé privée à stocker. On pousse sur GHCR
taggée par le SHA du commit. La CI bump alors la branche deploy/prod avec ce tag épinglé, et ArgoCD,
qui surveille cette branche, synchronise le cluster k3s tout seul.

La clé de tout : pas de tag latest en prod. Git fait foi — l'état déployé est exactement ce qui est
écrit dans deploy/prod. Donc un rollback, ce n'est pas une commande risquée : c'est restaurer un commit.
-->

---
layout: default
title: Déploiement infra
clicks: 6
---

<Slide kicker="Au-delà du runtime" title="Une commande monte tout — et referme SSH derrière elle"
       sub="terraform crée le serveur · SSH s'ouvre le temps qu'Ansible pose le cluster · puis se referme" accent="claude">

<div style="height: 410px">
  <FlowDiagram
    :vb="[1040, 560]"
    :nodes="[
      { id:'gha',  label:'GitHub Action', sub:'infra.yml (manuel)', x:30,  y:80,  w:210, h:88, variant:'ink', at:0 },
      { id:'tf',   label:'Terraform', sub:'+ Terraform Cloud', x:320, y:80,  w:200, h:88, variant:'accent', at:1 },
      { id:'srv',  label:'Serveur Hetzner', sub:'port 22 fermé au repos', x:610, y:70,  w:400, h:110, variant:'ink', at:2 },
      { id:'k3s',  label:'k3s', sub:'cluster', x:650, y:330, w:160, h:84, at:4 },
      { id:'mesh', label:'Istio + Kiali', sub:'mesh', x:830, y:330, w:170, h:84, at:4 },
      { id:'argo', label:'ArgoCD', sub:'GitOps', x:650, y:445, w:160, h:84, at:4 },
      { id:'cfd',  label:'cloudflared', sub:'tunnel', x:830, y:445, w:170, h:84, at:4 },
    ]"
    :edges="[
      { from:'gha', to:'tf', label:'apply', variant:'soft', at:1 },
      { from:'tf', to:'srv', label:'crée', variant:'default', at:2 },
      { from:'srv', to:'k3s',  label:'Ansible installe', variant:'ok', at:4, curve:true },
      { from:'srv', to:'mesh', variant:'ok', at:4, curve:true },
      { from:'srv', to:'argo', variant:'ok', at:4, curve:true },
      { from:'srv', to:'cfd',  variant:'ok', at:4, curve:true },
    ]"
  />
</div>

<div class="ssh-band">
  <span v-click="3" class="ssh-step ssh-open">①&nbsp; ouvre le port&nbsp;22 <i>— pour l'IP du runner /32</i></span>
  <span v-click="4" class="ssh-arrow">→</span>
  <span v-click="4" class="ssh-step ssh-work">②&nbsp; Ansible configure en SSH</span>
  <span v-click="5" class="ssh-arrow">→</span>
  <span v-click="5" class="ssh-step ssh-close">③&nbsp; referme le 22 <i>— always(), même en cas d'échec</i></span>
</div>

<div v-click="6" class="ssh-punch">
  Hors de cette fenêtre, <b>aucun port SSH ouvert</b> : un scan de l'IP ne trouve rien.
</div>

<style>
.ssh-band { display:flex; align-items:center; justify-content:center; gap:.7rem; margin-top:1rem; flex-wrap:nowrap; }
.ssh-step { font-size:.95rem; font-weight:700; padding:.5em 1.1em; border-radius:999px; border:1.5px solid; white-space:nowrap; }
.ssh-step i { font-weight:400; font-style:normal; opacity:.85; font-size:.85em; }
.ssh-open  { background:var(--ok-lt);     border-color:var(--ok);     color:var(--ok); }
.ssh-work  { background:var(--claude-lt); border-color:var(--claude); color:var(--claude-dk); }
.ssh-close { background:var(--danger-lt); border-color:var(--danger); color:var(--danger); }
.ssh-arrow { color:var(--mist); font-weight:800; font-size:1.2rem; }
.ssh-punch { text-align:center; margin-top:.8rem; font-size:1.05rem; color:var(--ink-soft); }
</style>

</Slide>

<!--
Comment toute la plateforme se monte, en une commande — et le point fort : SSH se referme tout seul
derrière elle.

Le flux, en haut. On déclenche le workflow infra.yml à la main. Il appelle Terraform — exécuté sur
Terraform Cloud, donc les secrets (token Hetzner, clé SSH) n'apparaissent jamais dans GitHub. Terraform
crée le serveur Hetzner, dont le port 22 est fermé au repos.

Puis Ansible installe tout le DEDANS du serveur — c'est ce qu'on voit dans le cadre : k3s pour le
cluster, Istio et Kiali pour le mesh, ArgoCD pour le GitOps, et cloudflared pour le tunnel. Voilà ce
qui est réellement posé : tout le runtime de la démo.

Mais pour qu'Ansible puisse se connecter, il faut SSH — et c'est là le zero-trust, la bande du bas, en
trois temps. Un : le workflow ouvre le port 22, uniquement pour l'IP du runner GitHub, en /32. Deux :
Ansible fait son travail pendant cette fenêtre. Trois : le workflow referme le 22 — et c'est un step
always(), donc il referme MÊME si Ansible a planté.

La conclusion : hors de cette courte fenêtre, aucun port SSH n'est ouvert. Un scan de l'IP ne trouve
rien. SSH n'est jamais une porte ouverte en permanence — il vit quelques minutes, pour une seule IP.
-->

---
layout: default
title: Tunnel & Access
clicks: 5
---

<Slide kicker="Au-delà du runtime" title="Zéro port entrant : le tunnel Cloudflare"
       sub="Le serveur n'ouvre aucun port web · c'est cloudflared qui sort vers Cloudflare · les dashboards exigent un login" accent="ok">

<div style="height: 430px">
  <FlowDiagram
    :vb="[1080, 560]"
    :nodes="[
      { id:'user', label:'Navigateur', sub:'internet', x:30,  y:250, w:175, h:90, variant:'ink', at:0 },
      { id:'cf',   label:'Cloudflare', sub:'edge réseau', x:280, y:246, w:195, h:100, variant:'accent', at:1 },
      { id:'cfd',  label:'cloudflared', sub:'dans le serveur', x:560, y:246, w:195, h:100, variant:'ink', at:2 },
      { id:'shop', label:'shop', sub:'public', x:870, y:90,  w:180, h:84, variant:'ok', at:2 },
      { id:'dash', label:'argocd · grafana · kiali', sub:'dashboards admin', x:830, y:420, w:220, h:96, at:3 },
    ]"
    :edges="[
      { from:'user', to:'cf', label:'HTTPS', variant:'default', at:1 },
      { from:'cfd', to:'cf', label:'tunnel SORTANT — aucun port entrant', variant:'ok', at:2, dashed:true, bend:150 },
      { from:'cfd', to:'shop', label:'libre', variant:'ok', at:2, curve:true },
      { from:'cfd', to:'dash', label:'Access : login', variant:'danger', at:3, curve:true },
    ]"
  />
</div>

<div class="cf-legend">
  <span v-click="4" class="cf-item"><b style="color:var(--ok)">↑ tunnel sortant</b> — le serveur n'expose <b>aucun port web</b> : rien à scanner</span>
  <span v-click="5" class="cf-item"><b style="color:var(--danger)">Access</b> — login email vérifié <b>à l'edge</b>, avant même d'atteindre le cluster</span>
</div>

<style>
.cf-legend { display:flex; gap:2.5rem; justify-content:center; margin-top:.6rem; flex-wrap:wrap; }
.cf-item { font-size:1rem; color:var(--ink-soft); }
</style>

</Slide>

<!--
Et la dernière pièce : comment on atteint le cluster sans ouvrir un seul port web. C'est le tunnel
Cloudflare.

Regardez le sens des flèches. cloudflared tourne DANS le cluster et ouvre une connexion SORTANTE vers
Cloudflare — en pointillé. Il n'y a aucun port entrant sur le serveur : ni 80, ni 443, ni 22. Un scan
de l'IP du serveur ne trouve rien. C'est ça « zéro port entrant » : la surface d'attaque réseau est
quasi nulle.

Le trafic arrive donc par Cloudflare, qui route quatre hostnames vers les bons services du cluster. Et
là, deuxième niveau de zero-trust : Cloudflare Access. La boutique, shop, reste publique — c'est le
produit. Mais les trois dashboards d'admin — ArgoCD, Grafana, Kiali — sont derrière Access : il faut un
login avec un email autorisé, vérifié à l'edge de Cloudflare, AVANT même que la requête touche le
cluster. Tout est déclaré en Terraform, dans le même apply que le reste.

En résumé : pas de SSH permanent, pas de port web ouvert, et les dashboards sensibles derrière un login.
Le zero-trust ne s'arrête pas aux microservices — il va jusqu'à la porte d'entrée de l'infra.
-->

---
layout: default
title: Bilan
clicks: 3
---

<Slide kicker="Bilan" title="Ce que cette histoire démontre" accent="ok">

<div class="text-3xl leading-relaxed mb-12 text-center">
  Chaque brique répond à un <b class="t-claude-dk">problème réel</b>.<br/>
  La sécurité a migré du <b>code</b> vers la <b class="t-ok">plateforme</b>.
</div>

<div class="bilan-grid">
  <div v-click="1" class="bilan-card" style="--c:var(--ink); --lt:var(--cloud)">
    <div class="bilan-n">01</div>
    <div class="bilan-t">Identité prouvée</div>
    <div class="bilan-s">SPIRE / SPIFFE · mTLS</div>
  </div>
  <div v-click="2" class="bilan-card" style="--c:var(--claude); --lt:var(--claude-lt)">
    <div class="bilan-n">02</div>
    <div class="bilan-t">Zero-trust visible</div>
    <div class="bilan-s">403 observable dans Kiali</div>
  </div>
  <div v-click="3" class="bilan-card" style="--c:var(--ok); --lt:var(--ok-lt)">
    <div class="bilan-n">03</div>
    <div class="bilan-t">Reproductible & signé</div>
    <div class="bilan-s">IaC · GitOps · supply chain</div>
  </div>
</div>

<style>
.bilan-grid { display:grid; grid-template-columns:repeat(3,1fr); gap:1.5rem; }
.bilan-card { position:relative; overflow:hidden; height:12rem; border-radius:14px;
  background:var(--lt); border:1px solid var(--c); border-top:6px solid var(--c);
  padding:1.5rem 1.4rem; display:flex; flex-direction:column; justify-content:center; }
.bilan-n { position:absolute; top:-1.2rem; right:.4rem; font-size:6rem; font-weight:800;
  line-height:1; color:var(--c); opacity:.14; }
.bilan-t { font-size:1.5rem; font-weight:800; color:var(--ink); }
.bilan-s { margin-top:.6rem; font-size:1.05rem; font-weight:600; color:var(--c); }
</style>

</Slide>

<!--
Pour conclure. Ce qui compte, ce n'est pas la liste des outils — c'est que CHAQUE brique répond à
un problème concret rencontré en avançant : identité, multi-langage, échelle.
Le fil rouge : la sécurité est partie du code pour devenir une propriété de la plateforme —
déclarative, réutilisable, observable.
Trois acquis : identité prouvée, zero-trust qu'on VOIT fonctionner, livraison reproductible et signée.
Merci, je suis prêt pour vos questions.
-->

---
layout: center
class: dark-slide side-accent
title: Merci
---

<div class="text-center">
  <div class="!text-8xl !font-extrabold t-white mb-6">Merci.</div>
  <div class="w-16 h-1 mx-auto mb-6" style="background:var(--claude); border-radius:3px"></div>
  <div class="text-2xl mb-14" style="color:#C9D2E0">
    Du code à la plateforme — une histoire zero-trust en trois temps
  </div>
  <div class="text-xl">
    <span class="font-bold t-white">Sylvain Rougié</span>
  </div>
</div>

<!--
Disponible pour les questions.
-->
