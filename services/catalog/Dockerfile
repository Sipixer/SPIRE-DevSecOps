# Build statique puis image distroless : pas de shell, surface d'attaque minimale.
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /catalog .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /catalog /catalog
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/catalog"]
