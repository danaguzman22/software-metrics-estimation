# Guía del proyecto: Software Metrics & Estimation (UTN FRSR - ICS 2026)

Lenguaje: **Go**. App web (Go + templates HTML). Persistencia: archivos JSON.
Equipo: Dana Guzmán (Agile Enabler / Product Builder), Victoria Caro Boldrini (Product Builder).

## Reglas que TODOS los agentes respetan
1. Todo el código es Go. No usar Python ni otros lenguajes (salvo HTML/CSS/JS de la interfaz).
2. Orden obligatorio por historia: **Spec (SDD) -> Escenarios (BDD) -> Tests (TDD, en RED) -> Código (GREEN) -> Refactor**.
3. Trazabilidad: Historia -> specs/US-XX.md -> criterios de aceptación -> features/US-XX.feature -> tests -> código.
4. La IA **nunca** hace push, nunca hace merge, nunca toca `main` ni `develop`. Eso lo hace una persona.
5. Una rama por issue: `feature/US-XX-nombre-corto`, creada desde `develop` actualizada.
6. Nada se da por bueno sin que una persona lo lea y lo valide. Cada agente termina diciendo qué debe revisar la persona.
7. La lógica de negocio vive en `internal/domain`; la capa web no duplica reglas.

## Estructura
cmd/server, internal/{domain,service,storage,web}, specs/, features/, docs/

## Agentes (5)
orquestador (coordina) · gestor-backlog (historias y GitHub) · especificador (SDD + BDD) · desarrollador-tdd (RED y GREEN por separado) · revisor (pre-revisión).

## Flujo de una tarea
1. `gestor-backlog` lista mis tareas -> elijo una.
2. Persona: `git checkout develop && git pull && git checkout -b feature/US-XX-...`
3. `especificador` -> spec + escenarios. Persona revisa.
4. `desarrollador-tdd` fase RED -> persona corre `go test`, confirma que FALLAN. Commit `test:`.
5. `desarrollador-tdd` fase GREEN -> persona corre `go test`, confirma que pasan. Commit `feat:`. Refactor y commit `refactor:`.
6. Persona: push, Pull Request a `develop`, tarjeta a "En revisión".
7. `revisor` pre-revisa; el Scrum Master revisa y hace el merge.
(O pedirle todo al `orquestador` con el número de issue.)

## Commits
Prefijos: `test:`, `feat:`, `refactor:`, `docs:`, `chore:`. Mensaje en español, referenciando el issue (ej. `test: US-01 validación de nombre vacío`).
