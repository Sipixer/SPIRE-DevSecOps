#!/usr/bin/env python3
"""Génère slides/SPIRE-Zero-Trust.pptx — deck sobre, schémas natifs, peu de texte.
L'oral porte le détail ; les slides ne sont que des repères visuels."""
from pptx import Presentation
from pptx.util import Inches, Pt, Emu
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN, MSO_ANCHOR
from pptx.enum.shapes import MSO_SHAPE
from pptx.enum.shapes import MSO_CONNECTOR
from pptx.oxml.ns import qn

# --- palette sobre ---
INK   = RGBColor(0x1A, 0x1F, 0x2B)   # texte principal
DIM   = RGBColor(0x6B, 0x72, 0x80)   # texte secondaire
BRAND = RGBColor(0x2F, 0x6F, 0xED)   # bleu accent
OK    = RGBColor(0x1F, 0xA9, 0x55)   # vert
KO    = RGBColor(0xD9, 0x3A, 0x3A)   # rouge
BG    = RGBColor(0xFF, 0xFF, 0xFF)   # fond blanc
SOFT  = RGBColor(0xF2, 0xF4, 0xF8)   # gris très clair (cartes)
LINE  = RGBColor(0xD7, 0xDC, 0xE5)   # bordures

prs = Presentation()
prs.slide_width  = Inches(13.333)
prs.slide_height = Inches(7.5)
BLANK = prs.slide_layouts[6]
W, H = prs.slide_width, prs.slide_height


def slide():
    s = prs.slides.add_slide(BLANK)
    bg = s.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, 0, W, H)
    bg.fill.solid(); bg.fill.fore_color.rgb = BG
    bg.line.fill.background()
    bg.shadow.inherit = False
    s.shapes._spTree.remove(bg._element)
    s.shapes._spTree.insert(2, bg._element)
    return s


def textbox(s, x, y, w, h, text, size=18, color=INK, bold=False,
            align=PP_ALIGN.LEFT, anchor=MSO_ANCHOR.TOP, font="Calibri", spacing=1.0):
    tb = s.shapes.add_textbox(x, y, w, h)
    tf = tb.text_frame; tf.word_wrap = True
    tf.vertical_anchor = anchor
    lines = text.split("\n")
    for i, line in enumerate(lines):
        p = tf.paragraphs[0] if i == 0 else tf.add_paragraph()
        p.alignment = align; p.line_spacing = spacing
        r = p.add_run(); r.text = line
        r.font.size = Pt(size); r.font.bold = bold
        r.font.color.rgb = color; r.font.name = font
    return tb


def box(s, x, y, w, h, label, sub=None, fill=SOFT, border=LINE, ink=INK, accent=None):
    sh = s.shapes.add_shape(MSO_SHAPE.ROUNDED_RECTANGLE, x, y, w, h)
    sh.fill.solid(); sh.fill.fore_color.rgb = fill
    sh.line.color.rgb = border; sh.line.width = Pt(1)
    sh.shadow.inherit = False
    if accent:
        bar = s.shapes.add_shape(MSO_SHAPE.RECTANGLE, x, y, Emu(int(w)//28), h)
        bar.fill.solid(); bar.fill.fore_color.rgb = accent
        bar.line.fill.background(); bar.shadow.inherit = False
    tf = sh.text_frame; tf.word_wrap = True
    tf.vertical_anchor = MSO_ANCHOR.MIDDLE
    if accent:
        tf.margin_left = Inches(0.25)
    p = tf.paragraphs[0]; p.alignment = PP_ALIGN.CENTER
    r = p.add_run(); r.text = label
    r.font.size = Pt(15); r.font.bold = True; r.font.color.rgb = ink; r.font.name = "Calibri"
    if sub:
        p2 = tf.add_paragraph(); p2.alignment = PP_ALIGN.CENTER
        r2 = p2.add_run(); r2.text = sub
        r2.font.size = Pt(10); r2.font.color.rgb = DIM; r2.font.name = "Calibri"
    return sh


def arrow(s, x1, y1, x2, y2, color=DIM, label=None, dashed=False):
    cn = s.shapes.add_connector(MSO_CONNECTOR.STRAIGHT, x1, y1, x2, y2)
    cn.line.color.rgb = color; cn.line.width = Pt(1.75)
    cn.shadow.inherit = False
    le = cn.line._get_or_add_ln()
    tail = le.makeelement(qn('a:tailEnd'), {'type': 'triangle', 'w': 'med', 'len': 'med'})
    le.append(tail)
    if dashed:
        d = le.makeelement(qn('a:prstDash'), {'val': 'dash'}); le.append(d)
    if label:
        mx, my = (x1 + x2)//2, (y1 + y2)//2
        textbox(s, mx - Inches(0.7), my - Inches(0.28), Inches(1.4), Inches(0.3),
                label, size=9, color=color, align=PP_ALIGN.CENTER)
    return cn


def kicker(s, text):
    textbox(s, Inches(0.7), Inches(0.45), Inches(8), Inches(0.4),
            text.upper(), size=12, color=BRAND, bold=True)


def title(s, text):
    textbox(s, Inches(0.7), Inches(0.8), Inches(12), Inches(0.9),
            text, size=30, color=INK, bold=True)


def footer(s, n):
    textbox(s, Inches(11.8), Inches(7.0), Inches(1.3), Inches(0.3),
            f"{n:02d}", size=10, color=DIM, align=PP_ALIGN.RIGHT)


# =====================================================================
# 1 — Titre
# =====================================================================
s = slide()
band = s.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, Inches(2.6), Inches(0.22), Inches(2.3))
band.fill.solid(); band.fill.fore_color.rgb = BRAND; band.line.fill.background(); band.shadow.inherit = False
textbox(s, Inches(0.9), Inches(2.55), Inches(11), Inches(0.5),
        "SPIFFE / SPIRE · mTLS · GitOps", size=14, color=BRAND, bold=True)
textbox(s, Inches(0.9), Inches(3.05), Inches(11.5), Inches(1.5),
        "Zero-Trust mTLS\npour microservices", size=44, color=INK, bold=True, spacing=1.0)
textbox(s, Inches(0.9), Inches(4.9), Inches(11), Inches(0.6),
        "Chaque service prouve son identité avant de parler à un autre.\nAucun secret partagé.",
        size=16, color=DIM, spacing=1.1)

# =====================================================================
# 2 — Le problème
# =====================================================================
s = slide(); kicker(s, "Le problème"); title(s, "Le réseau ne suffit pas")
textbox(s, Inches(0.7), Inches(1.9), Inches(11.9), Inches(0.8),
        "Dans un cluster, un service peut en joindre un autre par le réseau.\nMais le réseau ne prouve pas QUI appelle.",
        size=18, color=INK, spacing=1.2)
box(s, Inches(1.0), Inches(3.4), Inches(2.6), Inches(1.0), "Catalog", "compromis", fill=SOFT, accent=KO)
box(s, Inches(8.7), Inches(3.4), Inches(2.6), Inches(1.0), "Payments", "sensible", fill=SOFT, accent=KO)
arrow(s, Inches(3.6), Inches(3.9), Inches(8.7), Inches(3.9), color=KO, label="peut joindre ?")
textbox(s, Inches(1.0), Inches(5.0), Inches(11), Inches(1.0),
        "Si un attaquant compromet Catalog, rien au niveau réseau ne l'empêche\nd'atteindre Payments. Il faut une identité applicative vérifiable.",
        size=15, color=DIM, spacing=1.2)
footer(s, 2)

# =====================================================================
# 3 — Les services
# =====================================================================
s = slide(); kicker(s, "L'application"); title(s, "Une mini-boutique, 5 services")
data = [
    ("Gateway", "porte d'entrée", BRAND, 0.7),
    ("Orders", "commandes", INK, 3.18),
    ("Catalog", "produits", INK, 5.66),
    ("Payments", "paiements", KO, 8.14),
    ("Analytics", "Node.js / Envoy", DIM, 10.62),
]
for name, sub, acc, x in data:
    box(s, Inches(x), Inches(2.5), Inches(2.3), Inches(1.2), name, sub, accent=acc)
textbox(s, Inches(0.7), Inches(4.3), Inches(12), Inches(1.2),
        "Règle sensible :", size=16, color=INK, bold=True)
box(s, Inches(0.7), Inches(4.85), Inches(11.9), Inches(0.95),
    "Seul Orders peut déclencher un paiement.",
    "Gateway orchestre, Catalog expose — mais aucun ne peut appeler /pay directement.",
    fill=SOFT, accent=OK)
footer(s, 3)

# =====================================================================
# 4 — SPIRE : l'identité
# =====================================================================
s = slide(); kicker(s, "Identité"); title(s, "SPIRE attribue une identité cryptographique")
box(s, Inches(0.9), Inches(2.4), Inches(2.6), Inches(1.0), "SPIRE Server", "autorité / CA", accent=BRAND)
box(s, Inches(0.9), Inches(4.0), Inches(2.6), Inches(1.0), "SPIRE Agent", "DaemonSet")
arrow(s, Inches(2.2), Inches(3.4), Inches(2.2), Inches(4.0), color=BRAND, label="atteste")
for i, (svc, y) in enumerate([("gateway", 2.4), ("orders", 3.5), ("payments", 4.6)]):
    box(s, Inches(5.3), Inches(y), Inches(2.4), Inches(0.9), svc)
    arrow(s, Inches(3.5), Inches(4.5), Inches(5.3), Inches(y + 0.45), color=DIM)
textbox(s, Inches(8.1), Inches(2.5), Inches(4.5), Inches(2.5),
        "SVID X.509", size=18, color=INK, bold=True)
box(s, Inches(8.1), Inches(3.1), Inches(4.4), Inches(1.5),
    "spiffe://example.org/\nns/shop/sa/orders",
    "identité courte durée, rotative", fill=SOFT, accent=BRAND)
textbox(s, Inches(0.7), Inches(5.6), Inches(12), Inches(1),
        "L'identité dérive du ServiceAccount Kubernetes. L'Agent ne distribue un SVID\nqu'après avoir attesté le nœud (PSAT) et le workload.",
        size=14, color=DIM, spacing=1.2)
footer(s, 4)

# =====================================================================
# 5 — mTLS + autorisation (le cœur)
# =====================================================================
s = slide(); kicker(s, "Le cœur"); title(s, "mTLS transporte l'identité, le code l'autorise")
box(s, Inches(0.8), Inches(2.6), Inches(2.5), Inches(1.1), "Orders", accent=BRAND)
box(s, Inches(9.9), Inches(2.6), Inches(2.5), Inches(1.1), "Payments", accent=KO)
arrow(s, Inches(3.3), Inches(3.15), Inches(9.9), Inches(3.15), color=OK, label="mTLS  ·  /pay")
textbox(s, Inches(3.4), Inches(3.5), Inches(6.4), Inches(0.5),
        "présente son SVID — Payments lit l'identité dans le certificat",
        size=11, color=DIM, align=PP_ALIGN.CENTER)
# policy card
box(s, Inches(3.2), Inches(4.5), Inches(6.9), Inches(1.7),
    'var policy = map[string][]string{\n    "/pay": {ordersID},\n}',
    "Payments : seule l'identité Orders est autorisée sur /pay", fill=SOFT, accent=OK)
textbox(s, Inches(0.7), Inches(6.45), Inches(12), Inches(0.6),
        "Authentifié ≠ autorisé : être dans le trust domain ne suffit pas, il faut être dans la liste de la route.",
        size=14, color=INK, bold=True, align=PP_ALIGN.CENTER)
footer(s, 5)

# =====================================================================
# 6 — Démo : autorisé vs refusé
# =====================================================================
s = slide(); kicker(s, "Démonstration"); title(s, "Autorisé vs refusé, en direct")
# allowed
box(s, Inches(0.8), Inches(2.5), Inches(11.7), Inches(0.95),
    "Gateway → Orders → Payments     ✓ autorisé", fill=RGBColor(0xEC,0xF8,0xF0), border=OK, ink=OK)
# refused 1
box(s, Inches(0.8), Inches(3.7), Inches(11.7), Inches(0.95),
    "Gateway → Payments /pay     ✗ refusé", fill=RGBColor(0xFC,0xED,0xED), border=KO, ink=KO)
# refused 2
box(s, Inches(0.8), Inches(4.9), Inches(11.7), Inches(0.95),
    "Catalog → Payments /pay     ✗ refusé", fill=RGBColor(0xFC,0xED,0xED), border=KO, ink=KO)
textbox(s, Inches(0.7), Inches(6.1), Inches(12), Inches(0.6),
        "Même cluster, même réseau. Seule l'identité change la décision.",
        size=15, color=DIM, align=PP_ALIGN.CENTER)
footer(s, 6)

# =====================================================================
# 7 — Analytics derrière Envoy
# =====================================================================
s = slide(); kicker(s, "Multi-langage"); title(s, "Analytics (Node.js) délègue le mTLS à Envoy")
box(s, Inches(0.8), Inches(3.0), Inches(2.4), Inches(1.1), "Gateway", accent=BRAND)
big = box(s, Inches(5.2), Inches(2.4), Inches(6.5), Inches(2.4), "", fill=RGBColor(0xF7,0xF9,0xFC), border=LINE)
box(s, Inches(5.6), Inches(2.9), Inches(2.4), Inches(1.0), "Envoy", "SDS · mTLS", accent=BRAND)
box(s, Inches(8.7), Inches(2.9), Inches(2.6), Inches(1.0), "Analytics", "Node.js", accent=DIM)
arrow(s, Inches(3.2), Inches(3.55), Inches(5.6), Inches(3.4), color=OK, label="mTLS")
arrow(s, Inches(8.0), Inches(3.4), Inches(8.7), Inches(3.4), color=DIM, label="HTTP local")
textbox(s, Inches(0.7), Inches(5.3), Inches(12), Inches(1),
        "Envoy récupère le SVID via SDS auprès de l'Agent, termine le mTLS, et transmet\nl'identité vérifiée à l'app. Même sécurité sans réécrire SPIFFE dans chaque langage.",
        size=14, color=DIM, spacing=1.2)
footer(s, 7)

# =====================================================================
# 8 — La plateforme (infra)
# =====================================================================
s = slide(); kicker(s, "Plateforme"); title(s, "Infra reproductible : Terraform → Ansible")
steps = [("Terraform", "serveur Hetzner", 0.8),
         ("Ansible", "k3s · SPIRE · ArgoCD", 4.0),
         ("k3s", "cluster", 7.2),
         ("ArgoCD", "GitOps", 10.0)]
for name, sub, x in steps:
    box(s, Inches(x), Inches(3.0), Inches(2.7), Inches(1.2), name, sub, accent=BRAND)
for x in (3.5, 6.7, 9.5):
    arrow(s, Inches(x), Inches(3.6), Inches(x+0.5), Inches(3.6), color=DIM)
textbox(s, Inches(0.7), Inches(4.8), Inches(12), Inches(1),
        "Terraform crée le serveur (state dans Terraform Cloud). Ansible installe tout.\nArgoCD réconcilie le cluster depuis Git — modèle pull, rien n'entre de l'extérieur.",
        size=14, color=DIM, spacing=1.2)
footer(s, 8)

# =====================================================================
# 9 — CI/CD GitOps
# =====================================================================
s = slide(); kicker(s, "Livraison"); title(s, "Pipeline sécurisée + GitOps versionné")
ci = ["push", "lint · SAST", "Trivy (CVE)", "Cosign", "GHCR @sha", "deploy/prod", "ArgoCD sync"]
x = 0.55
for i, step in enumerate(ci):
    acc = OK if step in ("Trivy (CVE)", "Cosign") else BRAND
    box(s, Inches(x), Inches(2.7), Inches(1.62), Inches(1.0), step, fill=SOFT, accent=acc)
    if i < len(ci)-1:
        arrow(s, Inches(x+1.62), Inches(3.2), Inches(x+1.75), Inches(3.2), color=DIM)
    x += 1.78
box(s, Inches(0.8), Inches(4.4), Inches(11.7), Inches(1.0),
    "ArgoCD surveille deploy/prod, pas main.",
    "Les manifests y portent l'image épinglée au SHA du commit. Git = source de vérité, rollback = revert.",
    fill=SOFT, accent=OK)
textbox(s, Inches(0.7), Inches(5.7), Inches(12), Inches(0.6),
        "Une CVE haute non corrigée casse le build. Pas de tag latest en prod.",
        size=14, color=DIM, align=PP_ALIGN.CENTER)
footer(s, 9)

# =====================================================================
# 10 — Observabilité
# =====================================================================
s = slide(); kicker(s, "Observabilité"); title(s, "Prometheus + Grafana : qui, quoi, refusé")
cards = [
    ("SVID émis", "rotation visible", BRAND),
    ("Trafic HTTP", "par service / route", BRAND),
    ("Décisions", "allowed / denied", OK),
    ("Analytics", "via Envoy", DIM),
]
x = 0.8
for name, sub, acc in cards:
    box(s, Inches(x), Inches(2.8), Inches(2.8), Inches(1.4), name, sub, accent=acc)
    x += 3.0
textbox(s, Inches(0.7), Inches(4.7), Inches(12), Inches(1),
        "Le projet ne montre pas seulement que les appels fonctionnent.\nIl montre qui appelle, qui est autorisé, ce qui est refusé — et le rend observable.",
        size=15, color=INK, spacing=1.2, bold=False)
footer(s, 10)

# =====================================================================
# 11 — Récap / clôture
# =====================================================================
s = slide()
band = s.shapes.add_shape(MSO_SHAPE.RECTANGLE, 0, Inches(2.4), Inches(0.22), Inches(2.7))
band.fill.solid(); band.fill.fore_color.rgb = BRAND; band.line.fill.background(); band.shadow.inherit = False
textbox(s, Inches(0.9), Inches(2.4), Inches(11), Inches(0.6),
        "En une phrase", size=14, color=BRAND, bold=True)
textbox(s, Inches(0.9), Inches(2.95), Inches(11.6), Inches(2.0),
        "Des services qui se font confiance sans secret partagé :\ndes identités attestées, vérifiées en mTLS, autorisées par route,\ndéployées en GitOps et observables.",
        size=26, color=INK, bold=True, spacing=1.1)
textbox(s, Inches(0.9), Inches(5.4), Inches(11), Inches(0.5),
        "SPIFFE/SPIRE · Go + Envoy · Terraform/Ansible · ArgoCD · Prometheus/Grafana",
        size=14, color=DIM)

prs.save("slides/SPIRE-Zero-Trust.pptx")
print(f"OK — {len(prs.slides.__iter__.__self__._sldIdLst)} slides")
