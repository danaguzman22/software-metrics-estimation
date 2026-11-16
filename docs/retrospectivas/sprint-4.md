# Retrospectiva del Sprint 4

Sprint 4 — El Cierre y la Entrega Final. Fechas del Sprint: 16 de noviembre de 2026 (las tareas del Sprint se integraron y se cerró el Sprint ese día; Pull Requests #88 a #97).

## Resultado

| Elemento | Responsable | Resultado |
|---|---|---|
| US-21 Exportar informe a PDF (reporte del proyecto y del Sprint, generado en el servidor) | Dana | Entregada |
| T-14 Pantalla de integrantes del proyecto (alta y listado, aviso en Planning Poker sin integrantes) | Dana | Entregada |
| T-15 Navegación jerárquica: migas de pan y botón Volver en todas las pantallas | Dana | Entregada |
| T-16 Mejora visual de la interfaz (solo CSS) | Victoria | Entregada (PR #93, con una corrección posterior de las celdas de Estado) |
| T-17 Barra de navegación contextual y dashboard compacto | Dana | Entregada (issue #95) |
| Despliegue de la aplicación | Dana y Victoria | (completar: URL y plataforma, o "no realizado") |
| D-01 Documentación técnica y manual de usuario | Dana | Entregada (actualizados al cierre) |
| D-02 Informe de métricas y cobertura de pruebas | Dana | Entregada (cobertura regenerada al cierre) |
| D-03 Retrospectivas, trazabilidad y evidencias de TDD | Dana | Entregada |

El proyecto cuenta con el flujo completo de gestión: proyectos, integrantes, backlog, Sprints, Planning Poker, defectos, dashboard, reportes y exportación a PDF. Todo se ejecuta con Docker Compose.

## Qué salió bien

- El flujo SDD + BDD + TDD se sostuvo hasta el final, con evidencias RED, GREEN y REFACTOR por cada historia y tarea técnica funcional, y GREEN con PostgreSQL real y sin tests omitidos.
- La revisión con `@revisor` antes de cada Pull Request volvió a encontrar problemas antes del merge.
- Las necesidades que aparecieron al usar la aplicación (cargar integrantes desde la interfaz, poder volver al lugar de origen) se convirtieron en tareas técnicas (T-14 y T-15) con su spec, escenarios y tests, en lugar de resolverse sin registro.
- La exportación a PDF se resolvió en Go puro (biblioteca `go-pdf/fpdf`), sin navegador ni servicios externos, y un test de arquitectura impide que el generador dependa de HTTP o de la persistencia.
- La integración continua detectó un error que no se veía en local (un BOM dentro de archivos `.go`), y se incorporó un chequeo previo a cada commit.

## Qué mejorar

- Se generó evidencia RED antes de que el agente terminara su trabajo (todos los tests aparecían en verde), por lo que hubo que repetirla; conviene verificar que existan fallos reales (`--- FAIL`) y que no haya errores de compilación antes de commitear la evidencia.
- Dos veces el agente trabajó en una carpeta distinta de la rama (worktree) y hubo que mover archivos con `git stash`.
- Un commit de spec quedó después del de tests y hubo que reordenar la historia con cherry-pick y `--force-with-lease` sobre la propia rama.
- Los archivos con BOM (UTF-8 con marca) rompieron la compilación en el CI; se corrigió y se agregó un chequeo de BOM antes de cada commit.
- La tarea de mejora visual y el despliegue se planificaron al final del Sprint, con poco margen.
- La spec de la T-16 quedó en `develop` con marcas de conflicto de merge sin resolver (dos versiones pegadas); se detectó al armar la documentación de cierre y se corrigió. Conviene revisar los archivos con `git grep` de marcas de conflicto antes de cada merge.
- Recién al usar la aplicación en el último día aparecieron problemas de uso: la barra de navegación mostraba enlaces que no llevaban a ningún lado fuera de un proyecto y Defectos apuntaba a un marcador (se resolvió en la T-17). Convendría hacer un recorrido de uso completo antes del cierre de cada Sprint.
- Un proyecto Finalizado todavía permite modificar sus datos (integrantes, backlog, Sprints, defectos). Se detectó en la verificación manual del último día; no figuraba en los requisitos y quedó registrado en la issue #96 con la regla propuesta.
- Siguen sin existir pantallas para cargar horas estimadas y registros de esfuerzo: el dashboard y los reportes muestran horas solo si los datos se cargan por SQL.
- No hay inicio de sesión ni usuarios.

## Trabajo futuro

- Pantallas para registrar esfuerzo y horas estimadas, y un script de datos de demostración documentado.
- Autenticación y usuarios.
- Bloquear la edición de proyectos finalizados (solo lectura, mensaje claro y estado 422 en todas las acciones que modifican datos).
- Pulido visual adicional de la interfaz.
- Filtros por período y búsqueda en el proyecto.

## Conclusión

El equipo entrega una aplicación web completa para gestionar proyectos de software (proyectos, integrantes, backlog, Sprints, Planning Poker, defectos, métricas, dashboard, reportes y exportación a PDF), ejecutable con Docker Compose y desarrollada con SDD, BDD y TDD, con evidencias, trazabilidad y cobertura al cierre. Quedan fuera, y declaradas como trabajo futuro, la carga de esfuerzo y horas desde la interfaz, la autenticación y el bloqueo de proyectos finalizados, porque no figuraban en el alcance del Sprint o se detectaron al final.
