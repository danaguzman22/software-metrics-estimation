# Trazabilidad, retrospectivas y evidencias de TDD

El equipo decidió mantener estos registros de forma continua en cada Sprint y en cada historia, no solo al final.

## 1. Cadena de trazabilidad

Historia de usuario -> `specs/US-XX.md` (SDD) -> criterios de aceptación -> `features/US-XX.feature` (BDD) -> tests Go (RED -> GREEN -> REFACTOR) -> código Go.

## 2. Orden de commits por historia

| Etapa | Commit | Contenido |
|---|---|---|
| Especificación | `docs:` | `specs/US-XX.md` y `features/US-XX.feature` |
| RED | `test:` | Tests que fallan y `US-XX-red.txt` |
| GREEN | `feat:` | Implementación mínima |
| Evidencias | `docs:` | `US-XX-green.txt` (con PostgreSQL real y 0 tests omitidos) y `US-XX-refactor.txt` si hubo |

Las evidencias GREEN se generan con PostgreSQL 16 en Docker y la variable `TEST_DATABASE_URL`, de modo que los tests de integración se ejecutan y no se omiten.

## 3. Historias con spec, escenarios y evidencias

| Historia | Spec | Escenarios | Evidencias (`docs/evidencias/`) |
|---|---|---|---|
| US-01 | `specs/US-01.md` | `features/US-01.feature` | (ver spec) |
| US-02 | `specs/US-02.md` | `features/US-02.feature` | `US-02-red.txt`, `US-02-green.txt` |
| US-03 | `specs/US-03.md` | `features/US-03.feature` | `US-03-red.txt`, `US-03-green.txt` |
| US-04 | `specs/US-04.md` | `features/US-04.feature` | `US-04-red.txt`, `US-04-green.txt` |
| US-05 | `specs/US-05.md` | `features/US-05.feature` | `US-05-red.txt`, `US-05-green.txt` |
| US-06 | `specs/US-06.md` | `features/US-06.feature` | `US-06-red.txt`, `US-06-green.txt` |
| US-07 | `specs/US-07.md` | `features/US-07.feature` | `US-07-red.txt`, `US-07-green.txt`, `US-07-refactor.txt` |
| US-08 | `specs/US-08.md` | `features/US-08.feature` | `US-08-red.txt`, `US-08-green.txt`, `US-08-refactor.txt` |
| US-09 | `specs/US-09.md` | `features/US-09.feature` | `US-09-red.txt`, `US-09-green.txt` |
| US-11 | `specs/US-11.md` | `features/US-11.feature` | `US-11-red.txt`, `US-11-green.txt` |
| US-12 (Planning Poker) | `specs/US-12.md` | `features/US-12.feature` | `US-12-red.txt`, `US-12-green.txt`, `US-12-refactor.txt` |
| US-14 (esfuerzo) | `specs/US-14.md` | `features/US-14.feature` | `US-14-red.txt`, `US-14-green.txt`, `US-14-refactor.txt` |
| US-16 (velocidad) | `specs/US-16.md` | `features/US-16.feature` | `US-16-red.txt`, `US-16-green.txt`, `US-16-refactor.txt` |
| US-22 (interfaz web) | `specs/US-22.md` | `features/US-22.feature` | Fase 1: `US-22-fase1-red/green.txt`; Fase 2: `US-22-fase2-red/green/refactor.txt`; Fase 3: `US-22-fase3-red/green.txt`; `US-22-refactor.txt` |

Tareas técnicas (migración a PostgreSQL, Docker, CI y documentación de agentes): specs `T-02` a `T-13` en `specs/` y evidencias `T-XX-red/green/refactor.txt` en `docs/evidencias/`.

## 4. Retrospectivas

| Sprint | Documento |
|---|---|
| Sprint 1 — El MVP | [`retrospectivas/sprint-1.md`](retrospectivas/sprint-1.md) |
| Sprint 2 — La Interfaz | Se redactará al cerrar el Sprint |

Aprendizajes del Sprint 1 aplicados en el Sprint 2: definir fechas y Sprint Goal en el planning, exigir CI en verde y evidencias sin tests omitidos antes de mergear, y documentar en el momento.

## 5. Revisión por pares y por agente

Cada historia pasa por `@revisor` antes del Pull Request y por revisión cruzada entre las integrantes. Los hallazgos relevantes (por ejemplo, velocidad del Sprint sin dos decimales en US-22 Fase 3 o la comprobación de errores en los tests de concurrencia de US-12) se corrigieron con commits propios (`test:`, `feat:` o `fix:`) antes del merge.
