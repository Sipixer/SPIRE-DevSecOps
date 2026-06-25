# Alpine non-root, sans dépendances. npm est retiré pour ne pas traîner ses CVE.
FROM node:22-alpine
RUN apk -U upgrade --no-cache && rm -rf /usr/local/lib/node_modules/npm /usr/local/bin/npm /usr/local/bin/npx
WORKDIR /app
COPY main.js ./
USER node
EXPOSE 8080
CMD ["node", "main.js"]
