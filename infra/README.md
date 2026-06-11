# Infrastructure

Provisionnement et configuration de la plateforme, en deux étages séparés.

```
infra/
├── terraform/   crée le serveur Hetzner vide (serveur, firewall, SSH, IP)
└── ansible/     configure le serveur : k3s + SPIRE + ArgoCD
```

Terraform ne fait que l'infrastructure « physique ». Toute la configuration
logicielle est faite par Ansible. ArgoCD prend ensuite le relais : il tire les
manifests du dossier [`k8s/`](../k8s) et réconcilie le cluster (modèle pull).

## 1. Provisionner le serveur (Terraform Cloud, org `Sipixer`)

Prérequis : un token API Hetzner Cloud et la clé SSH publique, configurés comme
variables du workspace `spire-devsecops` sur Terraform Cloud.

```bash
cd infra/terraform
terraform login      # une seule fois
terraform init
terraform apply      # crée le serveur, écrit infra/ansible/inventory.ini
```

En sortie : l'IP publique du serveur et un inventaire Ansible prêt à l'emploi.

## 2. Configurer le serveur (Ansible)

Prérequis : `ansible` + la collection `kubernetes.core`.

```bash
cd infra/ansible
ansible-galaxy collection install kubernetes.core
ansible-playbook site.yml
```

Le playbook installe k3s, déploie SPIRE (chart Helm hardened) puis ArgoCD, et
déclare l'Application ArgoCD qui surveille `k8s/` sur la branche `main`.

## 3. Détruire (pour ne pas payer entre les sessions)

```bash
cd infra/terraform
terraform destroy
```

Hetzner facture à l'heure : `destroy` après chaque session, `apply` au début de
la suivante. Tout est reproductible.
