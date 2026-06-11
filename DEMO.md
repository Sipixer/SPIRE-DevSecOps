# Déroulé de démo

## Idée générale

La démo raconte l'évolution d'une application microservices vers une architecture zero-trust observable.

Le fil conducteur :

```text
réseau interne simple
-> identité de workload avec SPIFFE/SPIRE
-> mTLS et autorisation
-> service Node.js derrière Envoy
-> observabilité Prometheus/Grafana
```

## 1. Démarrer la plateforme

Commencer par lancer le workflow plateforme :

```text
Terraform -> VPS Hetzner -> Ansible -> k3s -> ArgoCD
```

Pendant que le setup tourne, présenter rapidement le projet avec quelques slides :

- les services de la boutique ;
- le problème du réseau interne ;
- la règle sensible : seul Orders peut déclencher un paiement ;
- le rôle de SPIRE, mTLS, Envoy et Grafana.

## 2. Montrer l'état initial

Montrer l'application déployée dans ArgoCD.

Présenter le modèle simple :

```text
Gateway -> Orders
Gateway -> Catalog
Orders -> Payments
```

Expliquer la limite :

```text
Le réseau interne permet de joindre des services,
mais il ne prouve pas finement l'identité applicative de l'appelant.
```

## 3. Merger la PR SPIRE + mTLS

Merger la première PR :

```text
feat: add SPIRE workload identity and Go mTLS authorization
```

Montrer rapidement la CI :

```text
tflint
Checkov
Semgrep
Trivy
Cosign
push GHCR
deploy/prod
```

Puis montrer ArgoCD qui synchronise.

Résultat à montrer :

```text
Gateway -> Orders -> Payments : autorisé
Gateway -> Payments /pay : refusé
Catalog -> Payments /pay : refusé
```

Message :

```text
SPIRE donne une identité.
mTLS transporte cette identité.
L'application lit l'identité et applique la policy.
```

## 4. Montrer le code clé

Montrer seulement les extraits essentiels :

```go
workloadapi.NewX509Source(...)
```

```go
tlsconfig.MTLSClientConfig(...)
tlsconfig.MTLSServerConfig(...)
```

```go
r.TLS.PeerCertificates[0].URIs[0]
```

Et la policy Payments :

```go
var policy = map[string][]string{
    "/pay": {ordersID},
}
```

## 5. Merger la PR Analytics + Envoy

Merger la deuxième PR :

```text
feat: add analytics service behind envoy
```

Montrer le nouveau service :

```text
Analytics Node.js
Envoy sidecar
```

Message :

```text
Les services Go gèrent SPIFFE/mTLS dans le code.
Analytics délègue cette partie à Envoy.
```

## 6. Lancer le mode démo

Lancer le mode démo depuis le front ou l'API.

La Gateway orchestre le scénario. Elle appelle les routes démo des services pour générer un trafic cohérent, proche d'un parcours utilisateur réel.

Scénario principal :

```text
1. Gateway -> Catalog
   récupération des produits

2. Gateway -> Orders
   création de commandes

3. Orders -> Payments
   paiements autorisés

4. Gateway -> Analytics via Envoy
   événements de navigation et de commande
```

Scénario de contrôle :

```text
Gateway -> Payments /pay
Catalog -> Payments /pay
```

Ces appels doivent être refusés. Ils servent à alimenter les métriques `denied` sans créer un trafic incohérent.

Chaque service garde une responsabilité claire :

| Service | Mode démo | Métriques utiles |
|---------|-----------|------------------|
| Gateway | Lance le scénario et agrège les résultats | nombre de scénarios, statuts globaux |
| Catalog | Sert les produits consultés | requêtes catalogue, latence |
| Orders | Crée les commandes | commandes créées, appels Payments |
| Payments | Accepte Orders, refuse les autres | paiements autorisés/refusés |
| Analytics | Compte vues et événements | événements reçus, identité transmise par Envoy |

Les métriques doivent rester lisibles :

```text
service
route
caller
decision = allowed | denied
status_code
```

## 7. Montrer Grafana

Afficher les dashboards :

- état SPIRE ;
- SVIDs émis et expirations ;
- trafic HTTP par service ;
- décisions autorisées/refusées ;
- activité Analytics derrière Envoy.

Message final :

```text
Le projet ne montre pas seulement que les appels fonctionnent.
Il montre qui appelle, qui est autorisé, ce qui est refusé,
et rend tout cela observable.
```
