# Etapa de build: compila el servidor con la versión de Go de go.mod.
FROM golang:1.24-alpine AS build

WORKDIR /src

# Se copian primero los manifiestos para aprovechar la caché de dependencias.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# pgx es puro Go: no hace falta CGO para el binario estático.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Etapa final: base mínima con una sonda (wget) para el healthcheck.
FROM alpine:3.20

# Usuario no root con UID/GID numéricos.
RUN addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app

# Las migraciones viajan embebidas en el binario; no se copian archivos .sql.
COPY --from=build /out/server /app/server

WORKDIR /app

USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["/app/server"]
