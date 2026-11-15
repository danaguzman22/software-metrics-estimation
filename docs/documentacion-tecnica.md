# Documentación técnica

Software Metrics & Estimation. Este documento describe cómo está construido el sistema. El uso de la aplicación está en [`manual-usuario.md`](manual-usuario.md).

## 1. Visión general

Aplicación web en Go que gestiona proyectos de software: Product Backlog, Sprints, estimación, esfuerzo, defectos y métricas. Se ejecuta con Docker Compose (aplicación + PostgreSQL 16) y se desarrolla con SDD + BDD + TDD.

## 2. Arquitectura

```text
cmd/server/            Punto de entrada: lee DATABASE_URL, aplica migraciones, arma repositorios -> servicios -> handlers
internal/domain/       Reglas de negocio puras. No depende de HTTP ni de la base
internal/service/      Casos de uso y contratos de repositorio (definidos del lado del consumidor)
internal/storage/      Adaptadores: jsonrepo (agregados originales) y postgres (todos)
internal/web/          Handlers HTTP, plantillas html/template y estilos, embebidos con go:embed
migrations/            Migraciones SQL con goose, embebidas y aplicadas al arrancar
specs/ features/       Especificaciones SDD y escenarios BDD por historia
docs/                  ADR, evidencias TDD, retrospectivas y documentación
```

Reglas de dependencia, verificadas por tests de arquitectura:

- `internal/domain` no importa nada del resto del proyecto.
- `internal/service` define las interfaces de repositorio (`RepositorioProyectos`, `RepositorioBacklog`, `RepositorioSprints`, `RepositorioAsignaciones`, `RepositorioIntegrantes`, `RepositorioEsfuerzo`, `RepositorioRondas`) y orquesta el dominio. No importa `pgx`.
- `internal/web` no importa `internal/storage` ni `pgx` y no contiene reglas de negocio: solo traduce HTTP a llamadas de servicio y errores de dominio a códigos de estado.
- `cmd/server` es el único lugar que conoce a la vez servicios y adaptadores.

## 3. Persistencia

Decisión registrada en [`adr/ADR-001-persistencia-postgresql.md`](adr/ADR-001-persistencia-postgresql.md): PostgreSQL con migraciones goose y driver `pgx v5`.

| Migración | Contenido |
|---|---|
| `00001_esquema_inicial.sql` | Proyectos, integrantes e ítems de backlog |
| `00003_sprints.sql` | Sprints (un solo Sprint activo por proyecto, índice parcial) |
| `00004_sprint_items.sql` | Asignación de historias a Sprints |
| `00005_cierre_sprint.sql` | Estados Completada/Cerrado y Story Points planificados al cierre |
| `00006_planning_poker.sql` | Rondas y votos de Planning Poker (una ronda abierta por historia) |
| `00007_esfuerzo.sql` | Registros de esfuerzo (horas por integrante, fecha y actividad) |
| `00008_metricas_esfuerzo.sql` | Horas estimadas por historia, Sprint de cada registro de esfuerzo y fotos de esfuerzo al cierre del Sprint |

Los agregados de Sprint, asignaciones, esfuerzo y Planning Poker se persisten solo en PostgreSQL, como excepción documentada en las enmiendas (d) a (g) del ADR-001. Los tests de integración con PostgreSQL se omiten si no hay `TEST_DATABASE_URL`; en la evidencia GREEN de cada historia se ejecutan con una base real y sin tests omitidos.

## 4. Capa web

Plantillas HTML embebidas (`internal/web/plantillas`) y una hoja de estilos (`estaticos/estilos.css`). Convenciones:

- Los formularios usan **POST-Redirect-GET**: éxito responde 303 al listado.
- Errores de validación y de dominio: **422** volviendo a renderizar el formulario con los datos ingresados y el mensaje exacto del dominio.
- Recurso inexistente: **404**. Método no permitido: **405**. Error de infraestructura: **500** con mensaje genérico (el detalle solo va al log).

Rutas principales:

| Ruta | Función |
|---|---|
| `GET /healthz` | Estado de la aplicación y de la base |
| `GET /proyectos`, `/proyectos/nuevo`, `POST /proyectos` | Listar y crear proyectos |
| `GET /proyectos/{id}` | Detalle y estado del proyecto (calculado por fechas) |
| `GET/POST /proyectos/{id}/backlog` (+ `/nuevo`, `/{itemID}/editar`, `/{itemID}/estimar`) | Product Backlog |
| `GET/POST /proyectos/{id}/sprints` (+ `/nuevo`) | Listar y crear Sprints |
| `POST /proyectos/{id}/sprints/{numero}/asignar` | Asignar una historia al Sprint |
| `POST .../sprints/{numero}/items/{itemID}/completar` | Marcar una historia como completada |
| `POST .../sprints/{numero}/cerrar` | Cerrar el Sprint |
| `GET /proyectos/{id}/sprints/historico` | Histórico de Sprints cerrados con Sprint Goal, SP planificados y completados y velocidad |
| `GET /proyectos/{id}/backlog/{itemID}/poker` (+ `POST .../poker/rondas`, `.../rondas/{n}/votos`, `.../rondas/{n}/revelar`, `.../rondas/{n}/acordar`) | Planning Poker: rondas, votos ocultos, revelación, diferencias, historial y estimación acordada |

## 5. Reglas de negocio destacadas

- **Un solo Sprint activo** por proyecto, garantizado por la lógica y por un índice parcial en la base.
- **Cierre de Sprint** en una única transacción: guarda los Story Points planificados, marca el Sprint como Cerrado y devuelve al backlog las historias no completadas.
- **Métricas**: Story Points planificados y completados por Sprint, velocidad del Sprint, velocidad promedio (solo Sprints cerrados; 0 si no hay) y porcentaje de historias completadas sobre todo el backlog. La web las muestra tal como las calcula el servicio, sin recalcular.
- **Planning Poker**: un voto por integrante por ronda, cambiable mientras la ronda está abierta; con la ronda abierta solo se informa quiénes votaron, nunca los valores; revelar es irreversible. Con la ronda revelada se calculan las diferencias (mínimo, máximo y quién votó cada valor), se conserva el historial y se puede acordar una estimación, que se guarda como Story Points de la historia. Una nueva ronda solo se abre si la anterior fue revelada.
- **Métricas de esfuerzo**: horas estimadas por historia, horas reales (suma de registros de esfuerzo) y desviación porcentual, sin división por cero; calculadas por historia, por Sprint y por proyecto.
- **Esfuerzo**: horas de 0 a 24 con dos decimales, fecha no futura, integrante del proyecto e historia que no pertenezca a un Sprint cerrado.

## 6. Proceso de desarrollo y trazabilidad

Historia de usuario -> `specs/US-XX.md` -> criterios de aceptación -> `features/US-XX.feature` -> tests (RED -> GREEN -> REFACTOR) -> código Go. La evidencia de cada etapa está en `docs/evidencias/`. Ver [`trazabilidad-y-evidencias.md`](trazabilidad-y-evidencias.md).

Agentes de IA (OpenCode): `@especificador` redacta specs y escenarios, `@desarrollador-tdd` escribe tests y código por etapas y `@revisor` revisa antes del Pull Request. Los agentes no hacen commit ni push; las personas revisan cada resultado. Las reglas están en `AGENTS.md`.

Flujo de ramas: `develop` es la rama de integración; cada historia vive en `feature/US-XX-...` y se integra por Pull Request con merge commit y CI en verde (`go vet` y `go test` con servicio PostgreSQL).

## 7. Ejecución y pruebas

```powershell
docker compose up --build -d --wait     # aplicación en http://localhost:8080
go vet ./...
go test ./... -count=1                   # sin TEST_DATABASE_URL se omiten los tests de PostgreSQL
```

Para los tests de integración completos ver la sección "Pruebas automatizadas" del [README](../README.md).

## 8. Limitaciones conocidas

- Registro de esfuerzo y métricas de horas: servicio y persistencia implementados (US-14 y US-17), sin interfaz web todavía.
- La interfaz web no permite cargar integrantes; para Planning Poker se cargan en la base (US-03 existe en el servicio).
- Dashboard con gráficos, reportes y exportación a PDF: Sprints 3 y 4.
