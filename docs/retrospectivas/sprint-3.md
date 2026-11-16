# Retrospectiva del Sprint 3

Sprint 3 — Funcionalidad y Calidad. Fechas del Sprint: (completar con las fechas reales del Sprint).

## Resultado

Historias entregadas y mergeadas en `develop`, cada una con su spec, escenarios BDD y evidencias RED, GREEN y REFACTOR:

| Historia | Responsable | Resultado |
|---|---|---|
| US-19 Dashboard con gráficos (servicio y web, dos gráficos SVG dibujados en el servidor) | Dana | Entregada |
| US-20 Reporte de proyecto o Sprint (páginas imprimibles) | Dana | Entregada |
| US-15 Gestión de defectos | Victoria | Entregada (servicio, persistencia e interfaz web) |
| US-18 Métricas de defectos | Victoria | Entregada |

No quedó trabajo del Sprint sin completar. El proyecto ya tiene el flujo completo de gestión: proyectos, backlog, Sprints, Planning Poker, defectos, dashboard y reportes. El registro de esfuerzo y las horas siguen sin pantalla de carga.

## Qué salió bien

- El flujo SDD + BDD + TDD se sostuvo también en las historias nuevas, con evidencias GREEN con PostgreSQL real y sin tests omitidos.
- La revisión con `@revisor` volvió a encontrar defectos reales antes del merge: en la US-19, la desviación por Sprint se calculaba en la capa web (contra la regla de que la web no recalcula métricas), había `template.HTML` innecesario y faltaba un test de composición del servidor; en la US-20, las reglas de impresión (`@media print`) no estaban acotadas al reporte. Todos se corrigieron con commits propios (`test:`, `feat:`, `docs:`) antes del merge.
- Los agentes informaron los tests contradictorios (por ejemplo, una condición invertida y un total de Story Points mal calculado en el seed de la US-20) en lugar de modificarlos en silencio; la persona decidió y la corrección quedó en un commit `test:` aparte.
- La lección de la US-19 se aplicó desde el principio en la US-20: la web no deriva cifras y no usa `template.HTML`, y no se repitió el hallazgo.
- Las decisiones se registraron como DECIDIDAS y DERIVADAS, y se enmendaron las specs cuando la revisión las contradijo.
- Los gráficos se dibujan en el servidor como SVG, sin librerías JavaScript ni CDN, así que la aplicación funciona dentro de Docker sin internet.

## Qué mejorar

- El primer push de la US-15 llegó con spec, tests, código y evidencias en un solo commit, y la evidencia GREEN omitía los tests de PostgreSQL; el CI detectó un error en un test de integración que no se había ejecutado en local. (Confirmar si el historial se reconstruyó antes del merge y ajustar esta línea.)
- Sigue sin existir pantalla para cargar integrantes, horas estimadas y registros de esfuerzo: el dashboard y el reporte muestran horas solo si los datos se cargan por SQL.
- No hay inicio de sesión ni usuarios; el sistema se usa sin autenticación.
- El aspecto visual de la interfaz es básico.
- Los conflictos entre ramas en archivos compartidos (`render.go`, `main.go`, plantillas del proyecto) hubo que resolverlos a mano en dos oportunidades; conviene traer `develop` a la rama antes de empezar cada fase.
- Algunas evidencias se generaron antes de aplicar las correcciones y hubo que regenerarlas.

## Acciones para el Sprint 4

- Agregar pantallas para integrantes, horas estimadas y registro de esfuerzo, y un script de datos de demostración documentado.
- Exportación del informe a PDF (a partir de los reportes imprimibles).
- Mejorar el aspecto de la interfaz y el dashboard.
- Mantener el orden de commits (`docs:` → `test:` → `feat:` → `docs:`) revisando el primer push de cada historia, y ejecutar siempre el GREEN con PostgreSQL real antes de abrir el PR.
- Traer `develop` a la rama de trabajo antes de cada fase para reducir conflictos.
- Regenerar el informe de cobertura al cierre del Sprint y mantener la documentación continua actualizada.
- Dejar declarado como trabajo futuro: autenticación y usuarios.
