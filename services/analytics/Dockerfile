# Pas de dépendances externes : on copie juste le source dans une image
# distroless Node non-root. L'app écoute sur 127.0.0.1:9090 (Envoy est devant).
FROM node:22-slim AS build
WORKDIR /app
COPY package.json main.js ./

FROM gcr.io/distroless/nodejs22-debian12:nonroot
WORKDIR /app
COPY --from=build /app /app
USER nonroot
EXPOSE 9090
CMD ["main.js"]
