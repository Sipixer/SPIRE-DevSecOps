# Pas de dépendances externes : on copie juste le source. Image Alpine minimale
# (openssl patché rapidement) en non-root. Écoute sur 127.0.0.1:9090, Envoy devant.
FROM node:22-alpine
WORKDIR /app
COPY package.json main.js ./
RUN apk -U upgrade --no-cache
USER node
EXPOSE 9090
CMD ["node", "main.js"]
