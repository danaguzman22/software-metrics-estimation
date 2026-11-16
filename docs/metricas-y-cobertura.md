# Informe de métricas y cobertura de pruebas

Software Metrics & Estimation. Estado al cierre del Sprint 4 (16/11/2026).

## 1. Cobertura de pruebas

**Cobertura total de sentencias: 80.7 %** (`go test ./... -coverprofile`), medida con PostgreSQL 16 real y la variable `TEST_DATABASE_URL`, de modo que los tests de integración se ejecutan y no se omiten. El detalle por función está en [`evidencias/cobertura-sprint-4.txt`](evidencias/cobertura-sprint-4.txt).

Cobertura por paquete:

| Paquete | Cobertura |
|---|---|
| `cmd/server` | 80.6 % |
| `internal/domain` | 99.5 % |
| `internal/reportepdf` | 93.3 % |
| `internal/service` | 93.4 % |
| `internal/storage` | 76.7 % |
| `internal/storage/jsonrepo` | 95.4 % |
| `internal/storage/postgres` | 75.9 % |
| `internal/web` | 84.8 % |
| `migrations` | sin sentencias |
| `internal/service/repotest` | 0.0 % (utilidades de la suite de contrato, sin tests propios) |

Evolución de la cobertura total: 77.1 % (cierre del Sprint 2, `evidencias/cobertura-sprint-2.txt`), 79.5 % (Sprint 3, `evidencias/cobertura-sprint-3.txt`) y 80.7 % (Sprint 4), con una base de código mayor en cada Sprint (PDF, integrantes y navegación en el último).

Cómo se reproduce (en PowerShell, los argumentos con punto van entre comillas):

```powershell
docker run --name pg-cov -e POSTGRES_PASSWORD=dev_password -e POSTGRES_DB=metricas -p 5440:5432 -d postgres:16
$env:TEST_DATABASE_URL = "postgres://postgres:dev_password@localhost:5440/metricas?sslmode=disable"
go test ./... -count=1 "-coverprofile=cover.out"
go tool cover "-func=cover.out"
docker rm -f pg-cov
```

Qué cubre la suite:

- **Dominio**: reglas puras (validaciones de ítems, Sprints, esfuerzo, Planning Poker) con tests unitarios.
- **Servicio**: casos de uso con dobles en memoria de los repositorios, incluyendo orden de validaciones y propagación de errores.
- **Persistencia**: tests de integración contra PostgreSQL (transacciones, índices parciales, claves foráneas, concurrencia de votos y migraciones Up/Down).
- **Web**: handlers con `httptest` (códigos 200, 303, 404, 405, 422 y 500, mensajes exactos, renderizado, navegación contextual y migas).
- **PDF**: tests del generador sobre los datos de los reportes y de los handlers que los descargan.
- **Arquitectura**: tests que verifican que `internal/web` no depende de `internal/storage` ni de `pgx`, que `internal/domain` no depende de infraestructura y que `internal/reportepdf` no depende de HTTP ni de la persistencia.

Integración continua: cada push y Pull Request ejecuta `go vet` y `go test` con un servicio PostgreSQL (`.github/workflows/ci.yml`).

## 2. Métricas del proceso

| Métrica | Valor | Fuente |
|---|---|---|
| Historias planificadas del Sprint 1 | 39 Story Points | [`retrospectivas/sprint-1.md`](retrospectivas/sprint-1.md) |
| Trabajo técnico no planificado en el Sprint 1 | ~35 Story Points (T-03 a T-13, migración a PostgreSQL) | idem |
| Historias del Sprint 3 | 21 Story Points (US-19: 8, US-20: 5, US-15: 5, US-18: 3), todas entregadas | [`retrospectivas/sprint-3.md`](retrospectivas/sprint-3.md) |
| Historia del Sprint 4 | US-21 Exportar informe a PDF (5 Story Points), entregada | [`retrospectivas/sprint-4.md`](retrospectivas/sprint-4.md) |
| Tareas técnicas del Sprint 4 | T-14 (3 Story Points), T-15, T-16 y T-17, todas entregadas | idem |
| Historias con spec, escenarios BDD y evidencias RED/GREEN | Todas las historias y tareas funcionales; T-16 sin RED por ser solo CSS | [`trazabilidad-y-evidencias.md`](trazabilidad-y-evidencias.md) |
| Evidencias GREEN con tests omitidos | 0 (PostgreSQL real en cada evidencia) | `docs/evidencias/*-green.txt` |
| Pull Requests integrados | Con CI en verde antes de mergear (regla del equipo desde el Sprint 2) | GitHub |
| Hallazgos de revisión corregidos antes del merge | Sprint 3: desviación calculada en la web (US-19) y reglas de impresión sin acotar (US-20); Sprint 4: validación del formato HU-NN y error de ruta al asignar con Sprint inválido (T-15) | retrospectivas |

## 3. Métricas del producto (calculadas por la aplicación)

La aplicación calcula y muestra, por Sprint y por proyecto:

- Story Points planificados (fijados al cerrar el Sprint) y completados.
- Velocidad del Sprint (Story Points completados) y velocidad promedio de los Sprints cerrados.
- Porcentaje de historias completadas sobre el backlog del proyecto.
- Horas estimadas, horas reales y desviación porcentual entre ambas, por historia, por Sprint y por proyecto (US-17), en el dashboard y en los reportes. La carga de horas todavía no tiene pantalla.
- Defectos detectados y resueltos por Sprint y por proyecto, con su distribución por severidad y estado (US-15 y US-18).
- Dashboard del proyecto con tarjetas de métricas y dos gráficos (velocidad por Sprint y horas estimadas vs. reales, US-19), y reportes de proyecto y de Sprint, imprimibles y descargables en PDF (US-20 y US-21).

Pendientes (trabajo futuro): pantallas de carga de horas y esfuerzo, autenticación y bloqueo de proyectos finalizados (issue #96).

## 4. Lectura y acciones

- La cobertura total de 80.7 % es razonable para el estado del proyecto. Los paquetes de dominio, servicio, PDF y repositorio JSON superan el 93 %; los más bajos son los adaptadores de persistencia (`internal/storage` y `internal/storage/postgres`, cerca del 76 %), donde faltan sobre todo rutas de error de infraestructura y de arranque.
- Acción para un próximo Sprint: reforzar las rutas de error de `internal/storage/postgres` y del arranque de `cmd/server`.
- Se mantiene la práctica de regenerar la cobertura al cierre de cada Sprint y guardarla en `docs/evidencias/`.
