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
Cloud. En sortie : l'IP du serveur et `infra/ansible/inventory.ini`.

## Configurer

```bash
cd infra/ansible
ansible-galaxy collection install kubernetes.core
ansible-playbook site.yml
```

## Détruire

```bash
cd infra/terraform && terraform destroy
```
