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
internal/reportepdf/   Generador de PDF (fpdf) de los reportes; no depende de HTTP ni de la persistencia
internal/web/          Handlers HTTP, plantillas html/template y estilos, embebidos con go:embed
migrations/            Migraciones SQL con goose, embebidas y aplicadas al arrancar
specs/ features/       Especificaciones SDD y escenarios BDD por historia
docs/                  ADR, evidencias TDD, retrospectivas y documentación
```

Reglas de dependencia, verificadas por tests de arquitectura:

- `internal/domain` no importa nada del resto del proyecto.
- `internal/service` define las interfaces de repositorio (`RepositorioProyectos`, `RepositorioBacklog`, `RepositorioSprints`, `RepositorioAsignaciones`, `RepositorioIntegrantes`, `RepositorioEsfuerzo`, `RepositorioRondas`,`RepositorioDefectos`) y orquesta el dominio. No importa `pgx`.
- `internal/web` no importa `internal/storage` ni `pgx` y no contiene reglas de negocio: solo traduce HTTP a llamadas de servicio y errores de dominio a códigos de estado.
- `internal/reportepdf` recibe los datos ya calculados por `ServicioReportes` y arma el PDF en Go puro (`go-pdf/fpdf`), sin navegador ni servicios externos. Un test de arquitectura impide que importe `net/http`, `internal/storage` o `pgx`.
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
| `00009_defectos.sql` | Defectos del proyecto (severidad, estado, historia relacionada, Sprint de detección y de resolución) con restricciones de integridad |

Los agregados de Sprint, asignaciones, esfuerzo y Planning Poker se persisten solo en PostgreSQL, como excepción documentada en las enmiendas (d) a (g) del ADR-001. Los tests de integración con PostgreSQL se omiten si no hay `TEST_DATABASE_URL`; en la evidencia GREEN de cada historia se ejecutan con una base real y sin tests omitidos.

## 4. Capa web

Plantillas HTML embebidas (`internal/web/plantillas`) y una hoja de estilos (`estaticos/estilos.css`). Convenciones:

- Los formularios usan **POST-Redirect-GET**: éxito responde 303 al listado.
- Errores de validación y de dominio: **422** volviendo a renderizar el formulario con los datos ingresados y el mensaje exacto del dominio.
- Recurso inexistente: **404**. Método no permitido: **405**. Error de infraestructura: **500** con mensaje genérico (el detalle solo va al log).

Navegación (T-15 y T-17): la barra superior es contextual. Fuera de un proyecto muestra solo **Proyectos**; dentro de un proyecto muestra Backlog, Sprints, Dashboard, Defectos, Reporte e Integrantes, apuntando a las rutas de ese proyecto. Cada pantalla muestra además migas de pan y un botón **Volver** (se ocultan al imprimir y no aparecen en los PDF). La barra la arma `navContexto` en `internal/web/render.go`, que cada handler completa con `navConProyecto(id)` o `navSinProyecto(...)`.

Interfaz (T-16 y T-17): la mejora visual se hizo solo con CSS (`estilos.css`, con bloques delimitados por tarea) y sin JavaScript nuevo. El dashboard muestra las métricas como tarjetas compactas y los dos gráficos en paneles lado a lado.

Rutas principales:

| Ruta | Función |
|---|---|
| `GET /healthz` | Estado de la aplicación y de la base |
| `GET /proyectos`, `/proyectos/nuevo`, `POST /proyectos` | Listar y crear proyectos |
| `GET /proyectos/{id}` | Detalle y estado del proyecto (calculado por fechas), con accesos a cada sección |
| `GET/POST /proyectos/{id}/integrantes` | Listar y registrar integrantes del proyecto con su rol (T-14) |
| `GET/POST /proyectos/{id}/backlog` (+ `/nuevo`, `/{itemID}/editar`, `/{itemID}/estimar`) | Product Backlog |
| `GET/POST /proyectos/{id}/sprints` (+ `/nuevo`) | Listar y crear Sprints |
| `POST /proyectos/{id}/sprints/{numero}/asignar` | Asignar una historia al Sprint |
| `POST .../sprints/{numero}/items/{itemID}/completar` | Marcar una historia como completada |
| `POST .../sprints/{numero}/cerrar` | Cerrar el Sprint |
| `GET /proyectos/{id}/sprints/historico` | Histórico de Sprints cerrados con Sprint Goal, SP planificados y completados y velocidad |
| `GET /proyectos/{id}/backlog/{itemID}/poker` (+ `POST .../poker/rondas`, `.../rondas/{n}/votos`, `.../rondas/{n}/revelar`, `.../rondas/{n}/acordar`) | Planning Poker: rondas, votos ocultos, revelación, diferencias, historial y estimación acordada |
| `GET/POST /proyectos/{id}/defectos` (+ `POST .../defectos/{defectoID}/estado`) | Registrar defectos, listarlos y cambiar su estado y Sprint de resolución |
| `GET /proyectos/{id}/dashboard` | Dashboard del proyecto: métricas clave y dos gráficos SVG (velocidad por Sprint y horas estimadas vs. reales) con tablas de datos |
| `GET /proyectos/{id}/reporte` | Reporte imprimible del proyecto completo |
| `GET /proyectos/{id}/reporte.pdf` | Reporte del proyecto en PDF, generado en el servidor (US-21) |
| `GET /proyectos/{id}/sprints/{numero}/reporte` | Reporte imprimible de un Sprint |
| `GET /proyectos/{id}/sprints/{numero}/reporte.pdf` | Reporte del Sprint en PDF (US-21) |
| `GET /defectos` | Marcador global conservado por compatibilidad; los defectos se consultan dentro de cada proyecto |

## 5. Reglas de negocio destacadas

- **Un solo Sprint activo** por proyecto, garantizado por la lógica y por un índice parcial en la base.
- **Cierre de Sprint** en una única transacción: guarda los Story Points planificados, marca el Sprint como Cerrado y devuelve al backlog las historias no completadas.
- **Métricas**: Story Points planificados y completados por Sprint, velocidad del Sprint, velocidad promedio (solo Sprints cerrados; 0 si no hay) y porcentaje de historias completadas sobre todo el backlog. La web las muestra tal como las calcula el servicio, sin recalcular.
- **Planning Poker**: un voto por integrante por ronda, cambiable mientras la ronda está abierta; con la ronda abierta solo se informa quiénes votaron, nunca los valores; revelar es irreversible. Con la ronda revelada se calculan las diferencias (mínimo, máximo y quién votó cada valor), se conserva el historial y se puede acordar una estimación, que se guarda como Story Points de la historia. Una nueva ronda solo se abre si la anterior fue revelada.
- **Métricas de esfuerzo**: horas estimadas por historia, horas reales (suma de registros de esfuerzo) y desviación porcentual, sin división por cero; calculadas por historia, por Sprint y por proyecto.
- **Esfuerzo**: horas de 0 a 24 con dos decimales, fecha no futura, integrante del proyecto e historia que no pertenezca a un Sprint cerrado.
- **Defectos**: severidad Alta, Media o Baja; estado Abierto, En curso o Resuelto; un defecto referencia una historia y un Sprint de detección del mismo proyecto; resolverlo exige informar un Sprint de resolución existente y no anterior al de detección.
- **Exportación a PDF (US-21)**: el reporte del proyecto y el de cada Sprint se pueden descargar en PDF. El PDF se arma con los mismos datos del reporte HTML, tomados de `ServicioReportes` sin recalcular.
- **Integrantes (T-14)**: nombre y rol obligatorios, sin duplicados dentro del proyecto (sin distinguir mayúsculas). Planning Poker avisa si el proyecto no tiene integrantes.
- **Dashboard y reportes**: un servicio arma todas las cifras reutilizando los servicios de Sprints, métricas de horas y defectos (`ServicioDashboard`, `ServicioReportes`); la capa web solo las formatea y no recalcula nada. Los gráficos son SVG dibujados en el servidor, sin librerías JavaScript ni CDN. El único JavaScript es el botón de imprimir del reporte (`window.print()`).

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

- Registro de esfuerzo y horas estimadas: servicio y persistencia implementados (US-14 y US-17), sin pantalla de carga; el dashboard y los reportes los muestran si los datos existen (se cargan por SQL, por ejemplo con un script de datos de demostración).
- No hay autenticación, usuarios ni protección CSRF (trabajo futuro). La aplicación se sirve por HTTP; en un despliegue real debería ir detrás de HTTPS.
- Un proyecto Finalizado todavía permite modificar sus datos. Se detectó el último día del Sprint 4; quedó registrado en la issue #96 con la regla propuesta (solo lectura, mensaje claro y estado 422).
- No hay filtros por período ni búsqueda dentro del proyecto.
- El proyecto no tiene un modelo de roles ni de permisos: cualquier persona con acceso puede ejecutar cualquier acción.
