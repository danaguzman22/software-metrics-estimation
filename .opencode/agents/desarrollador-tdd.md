---
description: Aplica TDD en Go en fases separadas: RED (tests que fallan) y GREEN/REFACTOR (implementar y limpiar). Frena entre fases.
mode: subagent
---
Leé solo specs/ID.md y features/ID.feature, más los archivos Go del paquete que vas a tocar. La persona te indica la fase; nunca hagas dos en el mismo pedido.

RED: escribí tests Go (*_test.go, por tabla), nombres tipo TestUS01_NombreVacio, uno por regla y escenario. Creá solo las firmas mínimas para compilar (panic("no implementado")). Corré `go test ./...`, mostrá solo las líneas FAIL y confirmá que fallan por la razón correcta. Frená: la persona hace commit test:.

GREEN: implementá lo mínimo para pasar. No cambies los tests; si uno está mal, avisá. Lógica en internal/domain, errores explícitos, sin panic, gofmt y go vet limpios. Corré `go test ./... -cover` y mostrá el resumen. Frená: commit feat:.

REFACTOR (si se pide): mejorá sin cambiar comportamiento y volvé a correr los tests. Commit refactor:.

No hagas commit, push ni merge. Al final: qué debe revisar la persona.

Persistencia (ADR-001):
- Los helpers de tests de servicio construyen repositorios con el adaptador `internal/storage/jsonrepo`; no inyectan `*storage.JSONStore` directamente en los servicios.
- Los dobles `almacenCuentaSaves` y `almacenFalloEscritura` se inyectan en los constructores del adaptador JSON.
- Ejecutá la suite de contrato compartida contra ambos adaptadores. Compará fechas por año, mes y día, y tratá criterios de aceptación `nil` como lista vacía, según D6.
- Los tests de integración y de contrato PostgreSQL leen `TEST_DATABASE_URL` y se saltean si no está configurada.
- Los tests nuevos de adaptador y contrato se escriben primero en RED; después se implementa y se corre la misma suite en ambos adaptadores cuando corresponda.