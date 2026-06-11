# Build statique (le front est embarqué via go:embed) puis image distroless.
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /gateway .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /gateway /gateway
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/gateway"]
