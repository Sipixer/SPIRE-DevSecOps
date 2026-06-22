# Runbook

Procédures opérationnelles. Le README décrit l'architecture, ce fichier décrit
comment lancer, démontrer et arrêter.

## Démarrer la plateforme

```bash
cd infra/terraform && terraform init
terraform apply                                                    # serveur, SSH fermé
terraform apply -var "allowed_ssh_cidr=$(curl -fsS https://api.ipify.org)/32"  # ouvre 22 pour ton IP
cd ..
ansible-galaxy collection install -r ansible/requirements.yml
ansible-playbook ansible/site.yml \
  -i "$(terraform -chdir=terraform output -raw server_ip)," \
  -u root --private-key ~/.ssh/spire-devsecops
terraform -chdir=terraform apply -var "allowed_ssh_cidr="          # referme 22
```

Ansible durcit le serveur (sshd, fail2ban, unattended-upgrades), puis installe
k3s, SPIRE, ArgoCD (qui surveille
`deploy/prod`) et l'observabilité. ArgoCD déploie ensuite les services depuis
GHCR. En pratique, préférer le workflow GitHub `infra` qui gère l'ouverture et la
fermeture SSH automatiquement.

Alternative : workflow GitHub `infra` (Actions → infra → Run), action `apply`.
Le workspace Terraform Cloud doit être en mode d'exécution **Remote** : les
variables `hcloud_token` et `ssh_public_key` y sont stockées, donc GitHub n'a
besoin que des secrets `TF_API_TOKEN` (auth Terraform Cloud) et
`SSH_PRIVATE_KEY` (pour qu'Ansible se connecte au serveur).

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

Le port 22 est **fermé par défaut** (accès SSH juste-à-temps, ouvert par la CI le
temps d'Ansible puis refermé). Pour une vérification ponctuelle, deux options :

```bash
# Option 1 : depuis le serveur, via la console Hetzner Cloud (hors SSH, zéro port).
k3s kubectl get pods -A
k3s kubectl get application shop -n argocd

# Option 2 : ouvrir SSH temporairement pour ton IP, puis le refermer.
cd infra/terraform
terraform apply -auto-approve -var "allowed_ssh_cidr=$(curl -fsS https://api.ipify.org)/32"
ssh -i ~/.ssh/spire-devsecops root@<ip> "k3s kubectl get pods -A"
terraform apply -auto-approve -var "allowed_ssh_cidr="   # referme
```

## Arrêter (ne pas payer le VPS)

```bash
cd infra/terraform && terraform destroy
```

Hetzner facture à l'heure : `destroy` après chaque session.
