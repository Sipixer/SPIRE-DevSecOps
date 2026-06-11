# Runbook

Procédures opérationnelles. Le README décrit l'architecture, ce fichier décrit
comment lancer, démontrer et arrêter.

## Démarrer la plateforme

```bash
cd infra/terraform && terraform init && terraform apply   # serveur Hetzner
cd ..
ansible-playbook ansible/site.yml \
  -i "$(terraform -chdir=terraform output -raw server_ip)," \
  -u root --private-key ~/.ssh/spire-devsecops
```

Ansible installe k3s, SPIRE, ArgoCD (qui surveille `deploy/prod`) et
l'observabilité. ArgoCD déploie ensuite les services depuis GHCR.

Alternative : workflow GitHub `infra` (Actions → infra → Run), action `apply`.
Requiert les secrets `TF_API_TOKEN`, `HCLOUD_TOKEN`, `SSH_PUBLIC_KEY`,
`SSH_PRIVATE_KEY`.

## Démo

| Branche | État | Bascule |
|---------|------|---------|
| `demo/base` | autorisation désactivée (réseau ouvert) | `AUTHZ_DISABLED=true` |
| `main` | mTLS + autorisation par route | défaut |

Scénario depuis la Gateway :

```bash
# autorisé
curl -X POST http://<gateway>/api/order
# refusé (Gateway n'a pas le droit sur /pay)
curl -X POST http://<gateway>/api/forbidden
# génère un trafic complet (commandes, refus, analytics via Envoy)
curl -X POST http://<gateway>/api/demo
```

Grafana : dashboard "SPIRE / Zero-trust" (décisions allowed/denied, trafic,
SVID). Identifiants démo : admin / admin.

## Vérifier l'état

```bash
ssh -i ~/.ssh/spire-devsecops root@<ip>
k3s kubectl get pods -A
k3s kubectl get application shop -n argocd
```

## Arrêter (ne pas payer le VPS)

```bash
cd infra/terraform && terraform destroy
```

Hetzner facture à l'heure : `destroy` après chaque session.
