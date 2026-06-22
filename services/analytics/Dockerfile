# Pas de dépendances : on copie juste le source. Image Alpine minimale, non-root.
# npm est retiré (inutile au runtime) pour ne pas traîner ses CVE transitives.
# L'app écoute sur :8080 ; le mTLS est assuré par le sidecar Istio (plus d'Envoy maison).
FROM node:22-alpine
RUN apk -U upgrade --no-cache && rm -rf /usr/local/lib/node_modules/npm /usr/local/bin/npm /usr/local/bin/npx
WORKDIR /app
COPY main.js ./
USER node
EXPOSE 8080
CMD ["node", "main.js"]
