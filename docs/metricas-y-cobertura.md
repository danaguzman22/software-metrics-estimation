# Informe de métricas y cobertura de pruebas

Software Metrics & Estimation. Estado al 15/11/2026 (Sprint 2 en curso).

## 1. Cobertura de pruebas

**Cobertura total de sentencias: 75.9 %** (`go test ./... -coverprofile`), medida con PostgreSQL 16 real y la variable `TEST_DATABASE_URL`, de modo que los tests de integración se ejecutan y no se omiten. El detalle por función está en [`evidencias/cobertura-actual.txt`](evidencias/cobertura-actual.txt).

Cómo se reproduce:

```powershell
docker run --name pg-cov -e POSTGRES_PASSWORD=dev_password -e POSTGRES_DB=metricas -p 5432:5432 -d postgres:16
$env:TEST_DATABASE_URL = "postgres://postgres:dev_password@localhost:5432/metricas?sslmode=disable"
go test ./... -count=1 -coverprofile=cover.out
go tool cover -func=cover.out
docker rm -f pg-cov
```

Qué cubre la suite:

- **Dominio**: reglas puras (validaciones de ítems, Sprints, esfuerzo, Planning Poker) con tests unitarios.
- **Servicio**: casos de uso con dobles en memoria de los repositorios, incluyendo orden de validaciones y propagación de errores.
- **Persistencia**: tests de integración contra PostgreSQL (transacciones, índices parciales, claves foráneas, concurrencia de votos y migraciones Up/Down).
- **Web**: handlers con `httptest` (códigos 200, 303, 404, 405, 422 y 500, mensajes exactos y renderizado).
- **Arquitectura**: tests que verifican que `internal/web` no depende de `internal/storage` ni de `pgx`, y que `internal/domain` no depende de infraestructura.

Integración continua: cada push y Pull Request ejecuta `go vet` y `go test` con un servicio PostgreSQL (`.github/workflows/ci.yml`).

## 2. Métricas del proceso

| Métrica | Valor | Fuente |
|---|---|---|
| Historias planificadas del Sprint 1 | 39 Story Points | [`retrospectivas/sprint-1.md`](retrospectivas/sprint-1.md) |
| Trabajo técnico no planificado en el Sprint 1 | ~35 Story Points (T-03 a T-13, migración a PostgreSQL) | idem |
| Historias con spec, escenarios BDD y evidencias RED/GREEN | Todas las implementadas | [`trazabilidad-y-evidencias.md`](trazabilidad-y-evidencias.md) |
| Evidencias GREEN con tests omitidos | 0 (PostgreSQL real en cada evidencia) | `docs/evidencias/*-green.txt` |
| Pull Requests integrados | Con CI en verde antes de mergear (regla del equipo desde el Sprint 2) | GitHub |

## 3. Métricas del producto (calculadas por la aplicación)

La aplicación calcula y muestra, por Sprint y por proyecto:

- Story Points planificados (fijados al cerrar el Sprint) y completados.
- Velocidad del Sprint (Story Points completados) y velocidad promedio de los Sprints cerrados.
- Porcentaje de historias completadas sobre el backlog del proyecto.

Pendientes para los próximos Sprints: Las horas estimadas, reales y la desviación (por historia, Sprint y proyecto) están calculadas en servicio y persistencia (US-17); falta su pantalla.

## 4. Lectura y acciones

- La cobertura total de 77.1 % es razonable para el estado del proyecto. El detalle por función está en `evidencias/cobertura-sprint-2.txt` y permite identificar los puntos con menor cobertura.
- Acción: revisar las funciones con menor cobertura (típicamente rutas de error de infraestructura y arranque) y reforzarlas en el Sprint 3.
- Acción: calcular y mostrar la cobertura por paquete en cada cierre de Sprint (el reporte del Sprint 1 había quedado sin datos por un error de montaje de Docker y se regeneró).
