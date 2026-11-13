# Guía del proyecto: Software Metrics & Estimation (UTN FRSR - ICS 2026)

Lenguaje: **Go**. App web (Go + templates HTML). Persistencia: PostgreSQL 13 o superior (pgx v5, migraciones goose embebidas); adaptador JSON (`jsonrepo`) para tests; repositorios por agregado definidos en `internal/service` (ADR-001).
Equipo: Dana Guzmán (Agile Enabler / Product Builder), Victoria Caro Boldrini (Product Builder).

## Reglas que TODOS los agentes respetan
1. Todo el código es Go. No usar Python ni otros lenguajes (salvo HTML/CSS/JS de la interfaz).
2. Orden obligatorio por historia: **Spec (SDD) -> Escenarios (BDD) -> Tests (TDD, en RED) -> Código (GREEN) -> Refactor**.
3. Trazabilidad: Historia -> specs/US-XX.md -> criterios de aceptación -> features/US-XX.feature -> tests -> código.
4. La IA **nunca** hace push, nunca hace merge, nunca toca `main` ni `develop`. Eso lo hace una persona.
5. Una rama por issue: `feature/US-XX-nombre-corto` o `feature/T-XX-nombre-corto`, creada desde la rama base indicada por el issue y actualizada.
6. Nada se da por bueno sin que una persona lo lea y lo valide. Cada agente termina diciendo qué debe revisar la persona.
7. La lógica de negocio vive en `internal/domain`; la capa web no duplica reglas; los contratos de repositorio los define el consumidor (`internal/service`) y los adaptadores viven en `internal/storage/{jsonrepo,postgres}`. El dominio no conoce la persistencia.
8. `DATABASE_URL` configura la conexión del servidor. `TEST_DATABASE_URL` configura la suite de contrato e integración PostgreSQL; esos tests se saltean si la variable no está.
9. Las migraciones ya aplicadas nunca se editan: cualquier cambio de esquema requiere una migración nueva.
10. Los adaptadores JSON y PostgreSQL se verifican con la suite de contrato compartida. Consultá [ADR-001](docs/adr/ADR-001-persistencia-postgresql.md) para las decisiones vigentes de persistencia.
11. Frontera de issues: solo T-10 modifica `.github/workflows/ci.yml`; el `compose.yml` pertenece a T-04.

## Estructura
cmd/server, internal/{domain,service,web}, internal/storage/{jsonrepo,postgres}, migrations/, specs/, features/, docs/

## Agentes (5)
orquestador (coordina) · gestor-backlog (historias y GitHub) · especificador (SDD + BDD) · desarrollador-tdd (RED y GREEN por separado) · revisor (pre-revisión).

## Flujo de una tarea
1. `gestor-backlog` lista mis tareas -> elijo una.
2. Persona: actualiza la rama base indicada por el issue y crea `feature/US-XX-...` o `feature/T-XX-...` desde ella.
3. `especificador` -> spec + escenarios. Persona revisa.
4. `desarrollador-tdd` fase RED -> persona corre `go test`, confirma que FALLAN. Commit `test:`.
5. `desarrollador-tdd` fase GREEN -> persona corre `go test`, confirma que pasan. Commit `feat:`. Refactor y commit `refactor:`.
6. Persona: push, Pull Request a la rama base indicada por el issue, tarjeta a "En revisión".
7. `revisor` pre-revisa; el Scrum Master revisa y hace el merge.
(O pedirle todo al `orquestador` con el número de issue.)

## Commits
Prefijos: `test:`, `feat:`, `refactor:`, `docs:`, `chore:`. Mensaje en español, referenciando el issue (ej. `test: US-01 validación de nombre vacío`).
