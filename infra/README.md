# Infrastructure

```
infra/
├── terraform/   serveur Hetzner (firewall, SSH, IP)
└── ansible/     k3s + SPIRE + ArgoCD
```

Terraform provisionne le serveur. Ansible le configure. ArgoCD tire ensuite les
manifests de [`k8s/`](../k8s) et réconcilie le cluster.

## Provisionner

```bash
cd infra/terraform
terraform init
terraform apply
```

Variables (`hcloud_token`, `ssh_public_key`) fournies par le workspace Terraform
Cloud. En sortie : `server_ip`.

## Configurer

Inventaire inline depuis l'IP Terraform, clé SSH en argument :

```bash
cd infra
ansible-galaxy collection install kubernetes.core
ansible-playbook ansible/site.yml \
  -i "$(terraform -chdir=terraform output -raw server_ip)," \
  -u root --private-key ~/.ssh/spire-devsecops
```

## Détruire

```bash
cd infra/terraform && terraform destroy
```
