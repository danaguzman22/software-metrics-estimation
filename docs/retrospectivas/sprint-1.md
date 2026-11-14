# Retrospectiva del Sprint 1

## Resultado
- Historias planificadas completadas: US-09 y US-16 (las últimas, mergeadas el 14/11) y las anteriores del Sprint.
- Se cerró con evidencias TDD por historia en docs/evidencias/.

## Qué pasó
- El Sprint desbordó: además de las historias planificadas (39 SP) hubo ~35 SP de trabajo técnico no planificado (T-03 a T-13) por la decisión de pasar la persistencia a PostgreSQL (ADR-001).
- Las tareas técnicas T-03 y T-04 se adelantaron para poder correr el sistema con Docker.

## Qué salió bien
- Flujo SDD + BDD + TDD sostenido: spec -> test -> código -> evidencias en cada historia.
- Revisión cruzada entre las dos que detectó defectos reales (evidencias RED, errores sin envolver, fixtures).

## Qué mejorar
- Estimar y planificar el trabajo de infraestructura antes de empezar el Sprint.
- Dejar Docker y el entorno de pruebas listos en las dos máquinas desde el inicio.
- Documentar en el momento, no al final.

## Acciones para el Sprint 2
- Definir fechas y Sprint Goal en el planning.
- Cada PR con CI en verde y evidencias sin tests omitidos antes de mergear.