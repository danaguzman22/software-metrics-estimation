# Software Metrics & Estimation

## Trabajo Práctico Integrador 2026

**Universidad Tecnológica Nacional**  
**Facultad Regional San Rafael**  
**Ingeniería en Sistemas de Información**  
**Ingeniería y Calidad de Software**

---

## Descripción del proyecto

**Software Metrics & Estimation** es una aplicación web para la gestión de proyectos de software: permite administrar proyectos, Product Backlog y Sprints, estimar historias, registrar esfuerzo y defectos, y obtener métricas de avance como Story Points planificados y completados, velocidad del equipo y porcentaje de historias completadas.

El proyecto corresponde al Trabajo Práctico Integrador de la asignatura Ingeniería y Calidad de Software 2026 y se desarrolla aplicando Scrum, SDD, BDD y TDD con herramientas de Inteligencia Artificial como soporte. Todo resultado generado con IA es comprendido, revisado y validado por el equipo.

---

## Equipo

| Integrante | Rol |
|---|---|
| **Dana Guzmán** | Agile Enabler / Product Builder |
| **Victoria Caro Boldrini** | Product Builder |

Los profesores de la asignatura cumplen el rol de **Product Architects**.

---

## Tecnologías

- **Go 1.24** (`net/http`, `html/template`, `embed`).
- **PostgreSQL 16** con migraciones [goose](https://github.com/pressly/goose) (`migrations/`) y driver `pgx v5`.
- **Docker / Docker Compose** para ejecutar la aplicación y la base de datos.
- **GitHub Actions** (CI con servicio PostgreSQL), **GitHub Projects** (tablero Scrum).
- **SDD + BDD + TDD** con agentes de OpenCode (`@especificador`, `@desarrollador-tdd`, `@revisor`).

---

## Ejecución con Docker

La aplicación y PostgreSQL se levantan con Docker Compose. Solo se publica el puerto 8080;
la base de datos queda accesible únicamente dentro de la red de Compose.

### Levantar todo

```powershell
docker compose up --build -d --wait
```

`--build` construye la imagen de la app, `--wait` espera a que los healthchecks de
`postgres` y `app` queden en `healthy` y `-d` deja los contenedores en segundo plano.
Verificar el estado:

```powershell
docker compose ps
curl.exe -s -o - -w "%{http_code}" http://localhost:8080/healthz
```

Luego abrir <http://localhost:8080> en el navegador.

### Despliegue de demostración

La aplicación también está desplegada en Render (plan gratuito, con base PostgreSQL gestionada y datos de demostración): <https://software-metrics-estimation.onrender.com>. Si estuvo inactiva, la primera carga puede tardar cerca de un minuto.

### Configuración por entorno

`compose.yml` funciona sin archivo `.env` porque usa valores por defecto de desarrollo
(`POSTGRES_USER=postgres`, `POSTGRES_DB=metricas`). Para personalizarlos, copiar
`.env.example` a `.env` y ajustar las variables. `.env` no se versiona y no debe contener
credenciales reales.

### Persistencia de datos

Los datos viven en el volumen nombrado `datos_postgres`, que no se borra al detener los
contenedores. Para reiniciar conservando los datos:

```powershell
docker compose down
docker compose up -d --wait
```

### Limpieza

Para detener y eliminar el volumen con los datos de la base:

```powershell
docker compose down -v
```

El siguiente `docker compose up` arranca con la base vacía y las migraciones se reaplican.

---

## Cómo usar la aplicación

Desde la página de inicio se accede a los proyectos. Por cada proyecto se puede:

1. **Crear el proyecto** con nombre y fechas de inicio y fin; el estado se calcula por fechas.
2. **Cargar los integrantes** con su rol (necesarios para Planning Poker).
3. **Administrar el Product Backlog**: alta de historias con título, descripción, prioridad, Story Points y criterios de aceptación; edición y priorización.
4. **Registrar y seguir defectos** asociados a historias y Sprints.
5. **Gestionar Sprints**: crear un Sprint con su objetivo y fechas (un solo Sprint activo a la vez), asignarle historias, marcar historias como completadas y cerrar el Sprint. Las historias no completadas vuelven al backlog.
6. **Ver métricas**: Story Points planificados, completados y velocidad de cada Sprint; velocidad promedio, porcentaje de historias completadas y horas estimadas vs. reales.
7. **Ver el dashboard** del proyecto, con tarjetas de métricas y gráficos de velocidad y horas por Sprint.
8. **Generar reportes** del proyecto completo o de un Sprint, para imprimir o descargar en PDF.

La barra superior muestra solo **Proyectos** fuera de un proyecto y, dentro de uno, las secciones de ese proyecto. Las migas de pan y el botón **Volver** indican dónde está la persona.

---

## Pruebas automatizadas

```powershell
go vet ./...
go test ./... -count=1
```

Sin la variable `TEST_DATABASE_URL`, los tests de integración con PostgreSQL se omiten (SKIP). Para ejecutarlos completos:

```powershell
docker run --name pg-test -e POSTGRES_PASSWORD=dev_password -e POSTGRES_DB=metricas -p 5432:5432 -d postgres:16
$env:TEST_DATABASE_URL = "postgres://postgres:dev_password@localhost:5432/metricas?sslmode=disable"
go test ./... -count=1
docker rm -f pg-test
```

La integración continua (`.github/workflows/ci.yml`) ejecuta `go vet` y `go test` con un servicio PostgreSQL en cada push y pull request.

---

## Arquitectura

```text
cmd/server/            Punto de entrada y composición de dependencias
internal/domain/       Reglas de negocio puras (sin infraestructura)
internal/service/      Casos de uso y contratos de repositorio
internal/storage/      Adaptadores: jsonrepo y postgres
internal/web/          Handlers HTTP, plantillas HTML y estilos embebidos
migrations/            Migraciones SQL (goose)
specs/                 Especificaciones SDD por historia de usuario
features/              Escenarios BDD (Gherkin) por historia de usuario
docs/adr/              Decisiones de arquitectura (ADR-001: persistencia)
docs/evidencias/       Evidencias TDD (RED, GREEN, REFACTOR) por historia
docs/retrospectivas/   Actas de retrospectiva y cobertura por Sprint
AGENTS.md              Reglas de trabajo con los agentes de IA
```

La capa web no contiene reglas de negocio ni depende de la persistencia; los contratos de repositorio se definen del lado del consumidor (`internal/service`). Los agregados de Sprint, asignaciones, esfuerzo y Planning Poker se persisten solo en PostgreSQL (enmiendas a ADR-001).

---

## Metodología de trabajo

El proyecto se gestiona con **Scrum** sobre GitHub Projects. Cada historia de usuario sigue el mismo ciclo:

**Historia de usuario → Especificación SDD → Criterios de aceptación → Escenarios BDD → Tests (RED → GREEN → REFACTOR) → Código Go → Revisión → Pull request**

- Cada historia tiene su `specs/US-XX.md` y su `features/US-XX.feature`.
- La evidencia de TDD queda en `docs/evidencias/US-XX-red.txt`, `US-XX-green.txt` (con PostgreSQL real y 0 tests omitidos) y `US-XX-refactor.txt`.
- Orden de commits por historia: `docs:` (spec) → `test:` (RED) → `feat:` (GREEN) → `docs:` (evidencias).
- Las ramas parten de `develop` y se integran por Pull Request con merge commit. Ningún agente hace commit ni push: lo hacen las personas.

### Sprints

| Sprint | Objetivo | Estado |
|---|---|---|
| Sprint 0 — La Preparación | Equipo, roles, entorno, repositorio, GitHub Projects y backlog inicial | Completado |
| Sprint 1 — El MVP | Versión funcional básica del sistema (dominio, servicios y persistencia PostgreSQL) | Cerrado. Retrospectiva en `docs/retrospectivas/sprint-1.md` |
| Sprint 2 — La Interfaz | Dotar al software de una interfaz usable | Cerrado. Retrospectiva en docs/retrospectivas/sprint-2.md |
| Sprint 3 — Funcionalidad y Calidad | Funcionalidades clave, robustez y visualización | Cerrado. Retrospectiva en docs/retrospectivas/sprint-3.md. |
| Sprint 4 — El Cierre y la Entrega Final | Pulido, documentación y exportación de informe en PDF | Cerrado. Retrospectiva en `docs/retrospectivas/sprint-4.md`. |

---

## Estado de las funcionalidades

| Funcionalidad | Estado |
|---|---|
| Gestión de proyectos (crear, listar, consultar estado calculado por fechas) | Implementada, con interfaz web |
| Integrantes del proyecto (alta y listado con rol) | Implementada, con interfaz web (T-14) |
| Product Backlog (alta, edición, priorización, estimación) | Implementado, con interfaz web |
| Sprints: crear, asignar historias, completar, cerrar, consultar anteriores | Implementados, con interfaz web (incluye histórico) |
| Métricas de Sprint y de proyecto: SP planificados y completados, velocidad, velocidad promedio, % de historias completadas | Implementadas y visibles en la interfaz |
| Planning Poker: votos ocultos, revelación, diferencias, nuevas rondas y estimación acordada | Implementado, con interfaz web por historia (US-12 y US-13) |
| Registro de esfuerzo (integrante, fecha, actividad, horas) | Implementado en servicio y persistencia (US-14); sin pantalla de carga |
| Métricas de horas y desviación estimado vs. real | Implementadas (US-17); se muestran en el dashboard y los reportes; las horas se cargan por SQL |
| Gestión de defectos (registro, listado, cambio de estado y Sprint de resolución) | Implementada, con interfaz web (US-15) |
| Métricas de defectos (detectados y resueltos por Sprint y por proyecto) | Implementadas (US-18) |
| Dashboard con gráficos (velocidad y horas por Sprint) | Implementado, con interfaz web (US-19); diseño compacto en T-17 |
| Reportes de proyecto y de Sprint (imprimibles) | Implementados, con interfaz web (US-20) |
| Exportación de los reportes a PDF, generada en el servidor | Implementada (US-21) |
| Navegación: migas de pan y botón Volver en todas las pantallas | Implementada (T-15) |
| Barra de navegación contextual (solo Proyectos fuera de un proyecto; secciones del proyecto dentro) | Implementada (T-17) |
| Mejora visual de la interfaz (solo CSS) | Implementada (T-16) |

### Limitaciones conocidas

- No hay inicio de sesión ni usuarios, ni protección CSRF.
- Las horas estimadas y los registros de esfuerzo no tienen pantalla de carga; se cargan por SQL.
- Un proyecto en estado Finalizado todavía puede modificarse (issue #96, registrada como trabajo futuro).
- No hay filtros por período ni búsqueda dentro del proyecto.
- La ruta global `/defectos` se conserva solo por compatibilidad; los defectos se consultan dentro de cada proyecto.

---

## Entregables

1. Repositorio Git con historial de contribuciones.
2. Tablero Scrum (GitHub Projects).
3. Product Backlog y Sprint Backlogs.
4. Especificaciones SDD (`specs/`).
5. Escenarios BDD (`features/`).
6. Pruebas automatizadas.
7. Código fuente en Go.
8. Evidencias de aplicación de TDD (`docs/evidencias/`).
9. Software funcional (Docker Compose).
10. Informe de métricas y cobertura de pruebas (`docs/metricas-y-cobertura.md`).
11. Actas de retrospectivas (`docs/retrospectivas/sprint-1.md` a `docs/retrospectivas/sprint-4.md`) y trazabilidad con evidencias de TDD (`docs/trazabilidad-y-evidencias.md`).
12. Documentación técnica (`docs/documentacion-tecnica.md`, ADR y specs) y manual breve de usuario (`docs/manual-usuario.md`).
13. Presentación y demostración final.

---

## Institución

**Universidad Tecnológica Nacional — Facultad Regional San Rafael**

**Carrera:** Ingeniería en Sistemas de Información  
**Asignatura:** Ingeniería y Calidad de Software  
**Año:** 2026