# Pas de dépendances : on copie juste le source. Image Alpine minimale, non-root.
# npm est retiré (inutile au runtime) pour ne pas traîner ses CVE transitives.
# L'app écoute sur 127.0.0.1:9090 ; Envoy est devant.
FROM node:22-alpine
RUN apk -U upgrade --no-cache && rm -rf /usr/local/lib/node_modules/npm /usr/local/bin/npm /usr/local/bin/npx
WORKDIR /app
COPY main.js ./
USER node
EXPOSE 9090
CMD ["node", "main.js"]
