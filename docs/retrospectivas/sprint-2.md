# Retrospectiva del Sprint 2

Sprint 2 — La Interfaz. Fechas del Sprint: (completar con las fechas reales del Sprint).

## Resultado

Historias entregadas y mergeadas en `develop`, cada una con su spec, escenarios BDD y evidencias RED, GREEN y REFACTOR:

| Historia | Responsable | Resultado |
|---|---|---|
| US-22 Interfaz web base (3 fases: proyectos y backlog, Sprints, cierre y velocidad) | Dana | Entregada |
| US-12 Planning Poker: votos ocultos y revelación (más la interfaz web, Fase 2) | Dana | Entregada |
| US-13 Planning Poker: diferencias, nuevas rondas y estimación acordada (más la interfaz web, Fase 2) | Dana | Entregada |
| US-14 Registro de esfuerzo | Victoria | Entregada (servicio y persistencia) |
| US-17 Métricas de horas y desviación | Victoria | Entregada (servicio y persistencia) |
| US-10 Consultar Sprints anteriores | Victoria | Entregada |

No quedó trabajo del Sprint sin completar. La interfaz web cubre proyectos, backlog, Sprints, métricas de velocidad y Planning Poker. El registro de esfuerzo y sus métricas existen en servicio y persistencia y quedan sin pantalla para el Sprint 3.

## Qué salió bien

- El flujo SDD + BDD + TDD se sostuvo en todas las historias: spec, tests en RED, implementación en GREEN, evidencias con PostgreSQL real y sin tests omitidos, y revisión antes del Pull Request.
- La revisión con `@revisor` y la revisión cruzada encontraron defectos reales antes del merge: velocidad del Sprint sin dos decimales (US-22), comprobación de errores faltante en los tests de concurrencia (US-12) y evidencias RED ausentes (US-17).
- Dividir las historias grandes en fases (US-22 en tres, US-12 y US-13 en servicio y web) permitió integrar de a poco y con CI en verde.
- Las decisiones se registraron como DECIDIDAS y DERIVADAS en cada spec, lo que evitó discusiones repetidas.
- La documentación se mantuvo continua: README, documentación técnica, manual breve, informe de métricas y cobertura, y trazabilidad.

## Qué mejorar

- Un PR de Victoria (US-17) llegó con un único commit que mezclaba spec, tests, código y evidencias, sin evidencia RED, igual que había pasado con la US-09. Hubo que rehacerlo.
- Los reinicios de la base con `docker compose down -v` borraron los datos de prueba varias veces, y no hay una pantalla para cargar integrantes: para probar Planning Poker hubo que insertarlos por SQL.
- El informe de cobertura del Sprint 1 se había generado con un error de montaje de Docker y quedó sin datos hasta que se regeneró.
- El plazo de entrega se interpretó mal durante una parte del trabajo y generó presión innecesaria.
- Las evidencias de algunas historias son resúmenes y no la salida completa de `go test`.

## Acciones para el Sprint 3

- Acordar desde el principio el orden de commits (`docs:` → `test:` → `feat:` → `docs:`) y generar la evidencia RED antes de implementar; revisar el primer push de cada historia contra esa regla.
- Agregar una pantalla para registrar integrantes y mostrar un aviso claro cuando un proyecto no tiene integrantes.
- Documentar un procedimiento de datos de demostración (script de carga) para no depender de recrear datos a mano.
- Verificar el plazo y el formato de entrega al inicio de cada Sprint y confirmar la rama que se presenta.
- Regenerar el informe de cobertura al cierre de cada Sprint, con el detalle por paquete.
- Priorizar la visualización (dashboard) y la interfaz del registro de esfuerzo con sus métricas de horas.
