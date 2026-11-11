# ADR-001: Migrar la persistencia de archivos JSON a PostgreSQL

- **Estado:** Aceptado con correcciones — 2026-11-11. Las correcciones 1 a 7 de la revisión están aprobadas e incorporadas (ver sección 11); el ADR sigue pendiente de aceptación final por la persona.
- **Fecha:** 2026-11-11.
- **Alcance:** capas `internal/service` y `internal/storage` (nuevos paquetes `internal/storage/jsonrepo` e `internal/storage/postgres`), migraciones SQL en `migrations/`, conexión de `cmd/server` y estrategia de tests. **Restricción global: los requisitos funcionales no cambian.** Los criterios de aceptación, las reglas y los escenarios BDD de US-01 a US-05 y US-11 se conservan tal como están. `internal/domain` se conserva íntegro (solo recibe adiciones puras, ver D4) y las firmas públicas de los **métodos** de servicio se conservan; los **constructores** cambian de parámetro (de `almacen` a repositorios): única excepción a las firmas públicas, aceptada como excepción documentada en la revisión del 2026-11-11 (D1; pendiente 1 cerrado).
- **Fuera de alcance:** `internal/web`, nuevas historias de usuario, cambios de requisitos, y los agregados sprints, votos, esfuerzo y defectos (D11).

---

## 1. Encabezado

Ver bloque superior. Resumen en una línea: reemplazar la persistencia actual (colecciones JSON completas cargadas y guardadas por cada operación) por repositorios por agregado con dos adaptadores (JSON para tests, PostgreSQL para el servidor), sin cambiar ningún requisito funcional.

## 2. Contexto

La persistencia actual nació en T-02 y tiene tres limitaciones de fondo:

1. **Interfaz de colección entera, no de registro.** El contrato `almacen { Save(name string, v any) error; Load(name string, v any) error }` está definido (no exportado) en `internal/service/proyectos.go:14` y solo sabe cargar y guardar una colección completa. Cada operación de servicio reescribe el archivo entero: `ModificarProyecto` guarda toda la colección `proyectos` (`internal/service/proyectos.go:81`), `AgregarIntegrante` toda la colección `integrantes` (`internal/service/integrantes.go:73`), `CrearItem` y `EstimarItem` toda la colección `backlog` (`internal/service/backlog.go:78` y `internal/service/estimacion.go:36`). Además, la lógica de "consulta" vive en el servicio: buscar un proyecto por ID recorriendo la colección (`internal/service/proyecto_existente.go:29`), buscar el ítem por (proyecto, ID) recorriendo el backlog (`internal/service/estimacion.go:20`). Cambiar de motor sin cambiar ese contrato obligaría a reescribir los servicios.

2. **Sin transacciones ni concurrencia.** No hay bloqueos entre la lectura y la escritura. El caso crítico es la numeración de ítems: `CrearItem` calcula el siguiente ID escaneando el máximo existente y guardando (`internal/service/backlog.go:63-78`). Dos `CrearItem` simultáneos sobre el mismo proyecto leen el mismo máximo, calculan el mismo `HU-NN` y el segundo pisa al primero. No hay unicidad por base de datos que lo impida, y el "Save atómico" de `*storage.JSONStore` (temporal + rename, `internal/storage/storage.go:43-70`) solo protege contra un crash a mitad de escritura, no contra escritores concurrentes.

3. **Errores y datos sin equivalente SQL ni integridad referencial.** `storage.ErrCorrupt` (JSON inválido, `internal/storage/storage.go:91`) es un concepto de archivo; el servicio dedica código a distinguirlo de `domain.ErrProyectoNoEncontrado` (patrón `errors.Is(err, storage.ErrNotFound)` en `proyectos.go:56-62`, `proyecto_existente.go:21-27`, `integrantes.go:51-57` y `backlog_coleccion.go:15-21`). No hay claves foráneas: un ítem o integrante puede existir para un proyecto inexistente si alguien siembra los archivos a mano. No hay unicidad de integrantes por proyecto garantizada por el motor (solo `domain.EsDuplicado` en memoria). Para T-03/T-04 (adaptador PostgreSQL y `compose.yml`, ya existentes como issues) hace falta un modelo relacional.

Datos verificados de la auditoría previa (tomados como dados):

- Colecciones en disco: `proyectos` = `[]proyectoRegistrado{ID, Proyecto}`, `integrantes` = `[]integranteRegistrado{ProyectoID, Integrante}`, `backlog` = `[]itemRegistrado{ProyectoID, Item}` (tipos no exportados en `internal/service`).
- `internal/service` tiene 63 tests que instancian un `*storage.JSONStore` real (helper sobre `t.TempDir()`); `internal/domain` tiene 45 tests puros; `internal/storage` tiene 9 tests (T-02).
- 6 tests de JSON corrupto exigen `errors.Is(err, storage.ErrCorrupt)`; 18 tests cuentan llamadas a `Save` con el doble `almacenCuentaSaves`; 4 tests de fallo de escritura usan `almacenFalloEscritura` (nombres exactos en D10).
- El ID de proyecto no se genera en ningún lado: no existe `CrearProyecto` en servicio; los tests siembran IDs a mano ("p1", "p2"); `domain.Proyecto` no tiene campo ID (decisión D11 de US-02).
- Specs con sección de persistencia: US-02 (§12, decisiones D11–D17), US-03 (§12), US-05 (§11), US-11 (§11). US-01 y US-04 la excluyen explícitamente.

## 3. Decisión

Se adopta la arquitectura de **repositorios por agregado definidos por el consumidor** (`internal/service`), con dos adaptadores (JSON y PostgreSQL), migraciones SQL versionadas y una suite de contrato compartida. Decisiones D1–D11.

### D1. Repositorios por agregado, contrato del lado del consumidor

- Las interfaces de repositorio se definen **en `internal/service`** (el consumidor define el contrato, puertos y adaptadores). El dominio no conoce la persistencia; la capa web no la duplica.
- Adaptadores:
  - `internal/storage/jsonrepo`: implementa los tres contratos **sobre el `JSONStore` actual, sin modificar `internal/storage`** (T-02 queda intacto con sus 9 tests). Conserva nombres de colecciones (`proyectos`, `integrantes`, `backlog`), formato de archivos y ordenamientos, de modo que los tests actuales siguen sembrando y leyendo archivos exactamente igual.
  - `internal/storage/postgres`: implementa los tres contratos sobre pgx v5 con `pgxpool` (D8).
- Dirección de las dependencias: los adaptadores importan `internal/service` (para satisfacer sus interfaces). No hay ciclo de importación en Go porque `internal/service` deja de importar `internal/storage` (la dependencia actual existe solo por `almacen`, que se elimina). Nota técnica: los tests de `internal/service` pueden importar `internal/storage/jsonrepo` aunque este importa `internal/service`; Go permite ciclos exclusivos de test. Si en T-09 esto presentara algún problema práctico, la alternativa de respaldo es mover las interfaces a un paquete neutro, pero **se decide mantener D1** salvo bloqueo real.
- `internal/storage` (JSONStore) **no se modifica**.
- Constructores de los servicios (consecuencia obligada de D1; única excepción a las firmas públicas — los métodos públicos no cambian —, aceptada como excepción documentada en la revisión del 2026-11-11; pendiente 1 cerrado):

```go
func NuevoServicioProyectos(
	proyectos RepositorioProyectos,
	integrantes RepositorioIntegrantes,
) *ServicioProyectos

func NuevoServicioBacklog(
	proyectos RepositorioProyectos,
	backlog RepositorioBacklog,
) *ServicioBacklog
```

`ServicioProyectos` orquesta `RepositorioProyectos` + `RepositorioIntegrantes` (para `ModificarProyecto` y `AgregarIntegrante`); `ServicioBacklog` orquesta `RepositorioProyectos` + `RepositorioBacklog` (para `CrearItem` y `EstimarItem`, que verifican la existencia del proyecto antes de tocar el backlog). Los métodos públicos existentes conservan firma y comportamiento observables:

- `(*ServicioProyectos) ModificarProyecto(id, nombre string, fechaInicio, fechaFin time.Time) (domain.Proyecto, error)` (US-02).
- `(*ServicioProyectos) AgregarIntegrante(proyectoID, nombre, rol string) (domain.Integrante, error)` (US-03).
- `(*ServicioBacklog) CrearItem(proyectoID, titulo, descripcion, prioridad string, criterios []string) (domain.ItemBacklog, error)` (US-05).
- `(*ServicioBacklog) EstimarItem(proyectoID, itemID string, storyPoints int) (domain.ItemBacklog, error)` (US-11).

### D2. Contratos propuestos (firmas Go exactas)

Se definen en `internal/service`, sobre los tipos reales del dominio (`domain.Proyecto`, `domain.Integrante`, `domain.ItemBacklog`). Ya incluyen las lecturas que necesita T-03 (arranque/web) y US-04, para no cambiar el contrato después: US-04 necesita el proyecto por ID (fechas para `EstadoPorFechas`) y los ítems del proyecto (para `PorcentajeCompletadas`); T-03 necesita listar proyectos.

```go
// ProyectoRegistrado une el ID persistido del proyecto con su agregado de
// dominio. El ID no forma parte de domain.Proyecto (decisión D11 de US-02):
// vive en la capa de servicio y en los contratos que lo devuelven.
// Conserva los json tags actuales para mantener el formato de archivos.
type ProyectoRegistrado struct {
	ID       string          `json:"id"`
	Proyecto domain.Proyecto `json:"proyecto"`
}

// RepositorioProyectos persiste el agregado de proyectos.
type RepositorioProyectos interface {
	// Obtener devuelve el registro {ID, Proyecto}.
	// domain.ErrProyectoNoEncontrado si el ID no existe.
	Obtener(id string) (ProyectoRegistrado, error)
	// Existe informa si hay un proyecto con ese ID, sin cargar el agregado.
	// Un repositorio vacío (sin datos) es Existe == false, no un error.
	Existe(id string) (bool, error)
	// Listar devuelve todos los registros de proyectos.
	Listar() ([]ProyectoRegistrado, error)
	// Actualizar reemplaza el proyecto con ese ID (el ID nunca cambia).
	// domain.ErrProyectoNoEncontrado si el ID no existe.
	Actualizar(id string, proyecto domain.Proyecto) error
}

// RepositorioIntegrantes persiste los integrantes por proyecto.
type RepositorioIntegrantes interface {
	// ListarPorProyecto devuelve los integrantes del proyecto en orden de
	// inserción (D7).
	ListarPorProyecto(proyectoID string) ([]domain.Integrante, error)
	// Agregar registra el integrante en el proyecto.
	// domain.ErrIntegranteDuplicado ante violación de unicidad (nombre del
	// integrante dentro del proyecto, sin distinguir mayúsculas).
	Agregar(proyectoID string, integrante domain.Integrante) error
}

// RepositorioBacklog persiste los ítems del Product Backlog por proyecto.
type RepositorioBacklog interface {
	// AgregarItem asigna el ID "HU-NN" de forma atómica (siguiente al mayor
	// número del proyecto, D4) y devuelve el ítem persistido con su ID.
	AgregarItem(proyectoID string, item domain.ItemBacklog) (domain.ItemBacklog, error)
	// ObtenerItem devuelve el ítem identificado por itemID dentro del proyecto.
	// domain.ErrItemNoEncontrado si no existe en ESE proyecto.
	ObtenerItem(proyectoID, itemID string) (domain.ItemBacklog, error)
	// ListarPorProyecto devuelve los ítems del proyecto ordenados por número
	// (HU-01, HU-02, ..., D7).
	ListarPorProyecto(proyectoID string) ([]domain.ItemBacklog, error)
	// ActualizarItem reemplaza el ítem identificado por item.ID dentro del
	// proyecto. domain.ErrItemNoEncontrado si no existe en ese proyecto.
	ActualizarItem(proyectoID string, item domain.ItemBacklog) error
}
```

Puntos de diseño resueltos y sus alternativas:

- **Dónde vive el ID de proyecto en los contratos que lo devuelven (`Obtener`, `Listar`).** Decidido (revisión 2026-11-11; pendiente 2 cerrado): el ID viaja en `service.ProyectoRegistrado{ID, Proyecto}`, la versión exportada del registro no exportado que ya existe en `internal/service/proyectos.go:21`, con los mismos json tags para conservar el formato de archivos; `domain.Proyecto` no se modifica (no tiene campo ID y no se le agrega, ver alternativas descartadas). Es la evolución natural de la decisión D11 de US-02 ("el ID vive en la capa de servicio") y es el tipo que los adaptadores y los helpers de siembra usan. Alternativa descartada: devolver `(domain.Proyecto, string, error)`.
- **Dónde vive el ID de ítem.** `domain.ItemBacklog` ya tiene `ID string` ("HU-NN"), así que el backlog no necesita registro: `proyectoID` es siempre parámetro de entrada. La PK real en PostgreSQL es `(proyecto_id, numero)` y el `HU-NN` se deriva de `numero` (D4).
- **El chequeo de duplicado de integrantes pasa al repositorio.** Hoy lo hace el servicio con `domain.EsDuplicado` (`internal/service/integrantes.go:59-67`). Con repositorios, la unicidad pasa a ser responsabilidad del adaptador (D3), conservando el mismo error observable (`domain.ErrIntegranteDuplicado`) y sin escribir nada cuando falla.
- **`AgregarItem` recibe el ítem sin ID y devuelve el persistido con ID**: asignar el ID requiere leer el estado del proyecto, que es exactamente lo que el repositorio sabe hacer atómicamente (D4).

### D3. Errores

- Los repositorios devuelven **los errores de dominio existentes**: `domain.ErrProyectoNoEncontrado` y `domain.ErrItemNoEncontrado` ante ausencia, `domain.ErrIntegranteDuplicado` ante violación de unicidad. Los demás errores se propagan **envueltos con `%w`**, sin mapearlos jamás a errores de dominio.
- El servicio devuelve los errores de dominio tal cual y envuelve con `%w` los que vienen de los repositorios (mismo criterio que hoy, p. ej. `internal/service/proyectos.go:58`).
- **El adaptador JSON sigue propagando `storage.ErrCorrupt` envuelto con `%w`** (para que los 6 tests de JSON corrupto actuales sigan en verde) y **traduce `storage.ErrNotFound` a "colección vacía"** (archivo inexistente = colección vacía, hoy decisión D14 de US-02 y equivalentes en US-03/US-05/US-11, implementada en `proyecto_existente.go:24-26` y `backlog_coleccion.go:19-20`). El adaptador PostgreSQL no conoce `storage.ErrCorrupt`; sus errores de infraestructura (conexión, SQL) se propagan envueltos con `%w`.
- El adaptador PostgreSQL traduce la violación de índice único (SQLSTATE 23505) de integrantes a `domain.ErrIntegranteDuplicado` (defensa ante carreras, ver D4/D6).

### D4. IDs

- **Proyectos:** PK `text` con `DEFAULT gen_random_uuid()::text`. Los tests pueden seguir sembrando IDs a mano ("p1", "p2"): el DEFAULT solo aplica cuando el INSERT no aporta ID. Decidido (revisión 2026-11-11): el DEFAULT queda en el esquema inicial aunque todavía no exista `CrearProyecto` ni operación de INSERT de proyectos; es latente hasta la historia de alta (brecha B1/B2).
- **Ítems:** PK compuesta `(proyecto_id, numero)` con ID observable `HU-NN` derivado de `numero`.
- **La regla "siguiente al mayor del proyecto" pasa a una función pura de dominio**, propuesta `domain.SiguienteIDItem`, usada por **ambos** adaptadores (hoy vive implícita en `internal/service/backlog.go:63-72` con `numeroDeItem` en `backlog.go:87`):

```go
// NumeroDeItem extrae el número de un ID con el patrón "HU-NN".
// Devuelve false si el ID no tiene ese patrón.
// (Mudanza de service.numeroDeItem, internal/service/backlog.go:87.)
func NumeroDeItem(id string) (int, bool)

// SiguienteIDItem devuelve el próximo ID de ítem del proyecto:
// "HU-01" si no hay ítems previos (o ninguno con patrón válido) y, si no,
// el mayor número existente más uno, con el formato actual "HU-%02d".
func SiguienteIDItem(idsExistentes []string) string
```

- **El adaptador JSON** aplica `SiguienteIDItem` sobre los IDs de la colección `backlog` del proyecto y luego guarda (comportamiento idéntico al actual).
- **El adaptador PostgreSQL** ejecuta la asignación **dentro de una transacción** con `SELECT ... FOR UPDATE` sobre la fila del proyecto (ver DDL en D6): bloquea la fila del proyecto, lee los IDs existentes del proyecto, aplica `domain.SiguienteIDItem` y hace el INSERT. El `FOR UPDATE` está sobre la fila de `proyectos` y no sobre los ítems: dos `CrearItem` de **distintos** proyectos no se bloquean entre sí (la numeración es por proyecto).
- **Riesgo de concurrencia que esto resuelve:** sin bloqueo, dos `CrearItem` simultáneos sobre el mismo proyecto hacen `SELECT max(numero)` (o leen la colección) al mismo tiempo, ambos obtienen el mismo máximo (p. ej. 4), ambos calculan `HU-05` y ambos hacen INSERT: uno gana y el otro viola la PK `(proyecto_id, numero)` y falla con un error de infraestructura en lugar de devolver el ítem. Con `FOR UPDATE` sobre la fila del proyecto, el segundo espera a que el primero confirme, relee, ve `HU-05` y calcula `HU-06`. La PK compuesta queda además como defensa en profundidad ante cualquier camino que no pase por la transacción.
- **Integrantes:** misma idea para el duplicado: el chequeo autoritativo es `domain.EsDuplicado` (ver D6), y el índice único de PostgreSQL traducido a `ErrIntegranteDuplicado` cubre la carrera entre dos `Agregar` simultáneos que pasan el chequeo con el mismo nombre.

### D5. Transacciones

Cada operación de servicio escribe un solo agregado (una fila de proyecto, o un integrante, o un ítem): no hace falta una unidad de trabajo (UoW) por ahora. Las únicas transacciones son internas del adaptador PostgreSQL: la de `AgregarItem` (asignación de ID atómica, D4) y la que use para la traducción de unicidad de integrantes si hace falta. Si aparece una operación que escriba dos agregados, se introduce UoW con su propia decisión.

### D6. Esquema inicial como DDL de referencia

El DDL siguiente es **de referencia** (no es una migración en sí misma): su materialización —como primera migración goose o como script de inicialización del `compose.yml` de T-04— es un pendiente (ver sección 11). Columnas derivadas de los tipos Go reales, incluidos los valores iniciales de `domain.NuevoItemBacklog` (`internal/domain/item_backlog.go:38-58`): estado siempre `"Pendiente"`, story points siempre `0` al crear; prioridad validada exactamente contra `"Alta" | "Media" | "Baja"`; título con espacios recortados; descripción y criterios tal como se reciben.

```sql
-- DDL de referencia. Pendiente: cómo se materializa (ver sección 11).
CREATE TABLE proyectos (
    id           text PRIMARY KEY DEFAULT gen_random_uuid()::text,
    nombre       text NOT NULL,
    fecha_inicio date NOT NULL,
    fecha_fin    date NOT NULL,
    estado       text NOT NULL CHECK (estado IN ('Planificado'))
);

CREATE TABLE integrantes (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,  -- orden de inserción (D7)
    proyecto_id text NOT NULL REFERENCES proyectos(id),
    nombre      text NOT NULL,
    rol         text NOT NULL
);

CREATE UNIQUE INDEX integrantes_proyecto_nombre_idx
    ON integrantes (proyecto_id, lower(nombre));

CREATE TABLE items_backlog (
    proyecto_id  text NOT NULL REFERENCES proyectos(id),
    numero       integer NOT NULL,          -- fuente de verdad del ID
    titulo       text NOT NULL,
    descripcion  text NOT NULL,             -- '' es válido
    prioridad    text NOT NULL CHECK (prioridad IN ('Alta', 'Media', 'Baja')),
    estado       text NOT NULL CHECK (estado IN ('Pendiente')),
    story_points integer NOT NULL CHECK (story_points IN (0, 1, 2, 3, 5, 8, 13, 21)),
    criterios_aceptacion text[] NOT NULL DEFAULT '{}',
    PRIMARY KEY (proyecto_id, numero)
);

CREATE INDEX items_backlog_proyecto_idx ON items_backlog (proyecto_id, numero);
```

Justificación de cada decisión del DDL:

- **`proyectos.estado` CHECK.** Hoy el único estado persistido es `"Planificado"` (`NuevoProyecto` lo fija, `ModificarProyecto` lo conserva; `internal/domain/proyecto.go:56` y `:83`). `EstadoPorFechas` calcula estados dinámicamente y no se persiste. El CHECK refleja lo persistible hoy y se extiende con una migración nueva cuando una historia persista más estados (las migraciones aplicadas nunca se editan, D8). Pendiente: confirmar el conjunto completo de valores de `domain.Estado` (el archivo `estado.go` no se leyó por restricción de este ADR).
- **`integrantes`: unicidad por proyecto y nombre.** El dominio compara con `strings.EqualFold` sobre nombres ya recortados (`internal/domain/integrante.go:52-59`). Diferencias con `lower()` de PostgreSQL: (1) `EqualFold` hace case-folding Unicode completo (caracteres fuera del ASCII, p. ej. el signo Kelvin U+212A se pliega con "k"), mientras que `lower()` depende del collation de la base; (2) `lower()` no recorta espacios, y el dominio sí recorta antes de comparar. Por eso **el dominio sigue siendo la validación autoritativa** (ambos adaptadores ejecutan `domain.EsDuplicado` antes de escribir, conservando el comportamiento observable exacto) y el índice único es **defensa en profundidad**: cubre la carrera de dos `Agregar` simultáneos que superan el chequeo con el mismo nombre, y su violación se traduce a `domain.ErrIntegranteDuplicado` (D3).
- **`integrantes.id` IDENTITY.** Da el orden de inserción exigido por D7 (`ListarPorProyecto` ordena por `id`).
- **`items_backlog.numero` como fuente de verdad.** El ID observable se deriva: al cargar, el adaptador reconstruye `ID = fmt.Sprintf("HU-%02d", numero)` (mismo formato que `internal/service/backlog.go:72`: ancho mínimo 2, sin truncar números mayores). La PK compuesta `(proyecto_id, numero)` hace que la numeración sea por proyecto y da el orden de D7. No se persiste el string `HU-NN` para evitar dos fuentes de verdad.
- **`story_points` CHECK IN (0, 1, 2, 3, 5, 8, 13, 21).** La serie de Fibonacci válida del dominio más el `0` inicial de `NuevoItemBacklog` (la validación estricta para estimar —`1, 2, 3, 5, 8, 13, 21`, `internal/domain/estimacion.go:10-17`— sigue en el dominio; el CHECK admite el 0 inicial).
- **`prioridad` y `estado` CHECK.** Reflejan los valores exactos del dominio (`internal/domain/item_backlog.go:45-49`): mayúscula inicial exacta, sin variantes.
- **`criterios_aceptacion` como `text[] NOT NULL DEFAULT '{}'`.** `[]string` de Go. Decidido (revisión 2026-11-11; pendiente 8 cerrado): el adaptador convierte `nil` en lista vacía al guardar y al leer; la distinción `null`/`[]` del JSON actual no se preserva. Consecuencia: si algún test dependiera de esa distinción tras un round-trip, se ajusta su expectativa (edición mecánica en T-09); la auditoría de tests no reporta casos.
- **Fechas: `time.Time` en Go, `date` en PostgreSQL.** El dominio interpreta las fechas por día en `time.UTC` (`soloDia` descarta la hora y usa `time.UTC`; las fechas cero no se validan). Decidido (revisión 2026-11-11; pendiente 9 cerrado): los adaptadores **normalizan a fecha UTC al guardar y al leer** — al guardar se convierte a UTC y se conserva solo el día (`date` en PostgreSQL); al leer se devuelve `time.Time` a medianoche UTC. La componente horaria no sobrevive al round-trip en ningún adaptador (el JSON actual sí la preserva). Aceptado porque las validaciones (`fechaFin.Before(fechaInicio)`) se ejecutan en el dominio **antes** de persistir y `EstadoPorFehas` compara por día. Los tests de la suite de contrato comparan fechas solo por año, mes y día; si algún test actual dependiera de la componente horaria tras un round-trip, se ajusta su expectativa (edición mecánica en T-09; la auditoría no reporta casos).
- **`gen_random_uuid()`** requiere **PostgreSQL 13 o superior** (nativo desde esa versión): requisito del ADR para el servidor, los tests de integración y la CI. La versión concreta de la imagen la fija T-04 en el `compose.yml`.

### D7. Orden

- Ítems ordenados por `numero` (ascendente) e integrantes por orden de inserción, en ambos adaptadores. El adaptador JSON conserva el orden del slice (append), igual que hoy. Esto conserva el comportamiento que asumen los tests actuales (p. ej. el caso "ConservaOtrosItemsYOrden" del grupo de US-11) y la numeración legible `HU-01, HU-02, ...`.

### D8. Tecnología

- **pgx v5 nativo con `pgxpool`**, sin ORM. Requiere **PostgreSQL 13 o superior** (requisito del ADR; habilita `gen_random_uuid()` nativo).
- **Migraciones SQL versionadas con goose embebidas** (`go:embed`) en `migrations/`, ejecutadas al arrancar `cmd/server` (antes de servir tráfico).
- Las migraciones ya aplicadas **nunca se editan**: cualquier cambio al esquema es una migración nueva.
- El DDL de D6 es la base de la primera migración (o del script de inicialización; pendiente, sección 11).

### D9. Configuración

- `DATABASE_URL`: cadena de conexión del servidor (pgxpool + migraciones al arrancar).
- `TEST_DATABASE_URL`: cadena de conexión para los tests de integración y la suite de contrato contra PostgreSQL; **los tests se saltean si la variable no está**.
- Sin ORM ni otras variables nuevas. Los tests de servicio y de dominio no requieren ninguna variable (siguen corriendo contra el adaptador JSON sobre `t.TempDir()`).

### D10. Estrategia de tests

- Los **63 tests de servicio siguen usando el adaptador JSON** (rápidos, sin base de datos). Para que los dobles actuales sigan funcionando, `internal/storage/jsonrepo` expone constructores que reciben una interfaz no exportada `almacen { Save; Load }` idéntica a la actual de `internal/service/proyectos.go:14`: `*storage.JSONStore` la satisface y los dobles existentes (`almacenCuentaSaves`, `almacenFalloEscritura`) también, por satisfacción estructural, **sin cambios de métodos** (solo la edición mecánica común: los tipos dejan de referenciar `service.almacen`, que se elimina, y la construcción del sujeto pasa de `almacen` a repositorios).
- **Suite de contrato compartida** para repositorios: los mismos casos corren contra el adaptador JSON y contra el de PostgreSQL (semántica idéntica: errores de dominio, traducción de ErrNotFound/vacío, orden, duplicados, asignación de ID). Decidido (revisión 2026-11-11; pendiente 11 cerrado): los 2 tests reubicados se ejecutan contra **ambos** adaptadores en esta suite, sin reducir cobertura; las comparaciones de fechas de la suite son solo por año, mes y día y las de criterios aceptan lista vacía en lugar de `nil` (normalizaciones de D6). Los helpers de la suite se comparten entre los tests de ambos adaptadores (detalle de implementación de T-08).
- **Tests de integración por adaptador** para lo que solo la base puede garantizar: FK (`proyecto_id` inexistente), CHECKs (prioridad, estado, story_points) y unicidad (índice de integrantes y PK de ítems, incluida la concurrencia de D4).

**Clasificación de los 6 tests de JSON corrupto (nombres exactos): conservados sin cambios.** La propiedad que verifican (un archivo corrupto produce un error de infraestructura que **no** se mapea a "no encontrado") sigue siendo observable a través del servicio con el adaptador JSON, porque D3 manda propagar `storage.ErrCorrupt` envuelto con `%w` y el servicio envuelve con `%w` (los `errors.Is` de los tests siguen encontrándolo):

1. `TestUS02_Modificar_JSONCorruptoNoEsNoEncontrado` — conserva.
2. `TestUS03_AgregarIntegrante_JSONCorruptoNoEsNoEncontrado` — conserva.
3. `TestUS03_AgregarIntegrante_JSONCorruptoIntegrantes` — conserva.
4. `TestUS05_CrearItem_JSONCorruptoNoEsNoEncontrado` — conserva.
5. `TestUS05_CrearItem_BacklogCorruptoNoEsNoEncontrado` — conserva.
6. `TestUS11_EstimarItem_JSONCorruptoPropagaError` — conserva.

Justificación por grupo: los seis siembran un JSON inválido en el directorio del store y aserten `errors.Is(err, storage.ErrCorrupt)` a través del servicio; el mecanismo de siembra y las aserciones quedan intactos porque el adaptador JSON usa los mismos archivos y nombres de colección. Los que corrompen `proyectos.json` (1, 2 y 4) pasan ahora por `RepositorioProyectos.Existe/Obtener`; los que corrompen `integrantes.json` (3) por `RepositorioIntegrantes.Agregar`; los que corrompen `backlog.json` (5 y 6) por `RepositorioBacklog.AgregarItem/ObtenerItem`. Además, `internal/storage/jsonrepo` recibe tests propios de propagación de `ErrCorrupt` y de traducción de `ErrNotFound` (nuevos, no reubicados).

**Clasificación de los 18 tests que cuentan Save (nombres exactos):** 16 conservados sin cambios y 2 reubicados a la suite de contrato.

*Conservados (16) — la propiedad es del servicio (orden de validaciones, proyecto antes que ítems, no escribir en caminos de error):*

1. `TestUS02_Modificar_ValidacionFallidaNoPersiste`
2. `TestUS03_AgregarIntegrante_RolVacioNoAgregaNiPersiste`
3. `TestUS03_AgregarIntegrante_ErrorNoAgregaNiPersiste`
4. `TestUS03_AgregarIntegrante_ProyectoNoEncontrado`
5. `TestUS05_CrearItem_TituloVacioNoAgregaNiPersiste`
6. `TestUS05_CrearItem_PrioridadInvalidaNoAgregaNiPersiste`
7. `TestUS05_CrearItem_ProyectoNoEncontrado`
8. `TestUS05_CrearItem_ErrorNoAgregaNiPersiste`
9. `TestUS11_EstimarItem_ProyectoNoEncontrado`
10. `TestUS11_EstimarItem_ItemNoEncontrado`
11. `TestUS11_EstimarItem_StoryPointsInvalidosNoPersiste`
12. `TestUS11_EstimarItem_ValorInvalidoConservaValorPrevio`
13. `TestUS11_EstimarItem_ProyectoInexistenteAntesQueValorInvalido`
14. `TestUS11_EstimarItem_ItemInexistenteAntesQueValorInvalido`
15. `TestUS11_EstimarItem_ErrorNoEscribe`
16. `TestUS11_EstimarItem_JSONCorruptoPropagaError` (está en ambos grupos: corrupto + conteo de Save; se conserva)

Justificación del grupo: siguen midiendo "no hubo `Save` en el camino de error" inyectando `almacenCuentaSaves` en el adaptador JSON (mecanismo intacto, ver arriba), y los que releen la colección para verificar que no cambió siguen leyendo los mismos archivos. El conteo sigue siendo 0 porque los caminos de error del servicio nunca llaman a las operaciones de escritura del repositorio, y las validaciones del dominio siguen ejecutándose antes de persistir.

*Reubicados (2) — la propiedad pasa a ser del repositorio (contrato), así que se mueven a la suite de contrato compartida, donde corren contra el JSON y contra PostgreSQL (decisión de la revisión 2026-11-11: contra ambos adaptadores, sin reducir cobertura):*

17. `TestUS11_EstimarItem_BacklogSinArchivo` — reubica. Su esencia ("archivo inexistente = colección vacía, no error de infraestructura") pasa a ser el contrato del adaptador (D3: `ErrNotFound` se traduce a colección vacía; en PostgreSQL el equivalente es "tabla vacía"). En la suite de contrato se expresa como "repositorio sin datos → `ObtenerItem` devuelve `domain.ErrItemNoEncontrado` y `ListarPorProyecto` devuelve vacío, sin error", contra ambos adaptadores.
18. `TestUS11_EstimarItem_MismoIDEnOtroProyecto` — reubica. El alcance por proyecto (el `HU-01` de `p2` no es el `HU-01` de `p1`) pasa a estar garantizado por el contrato `ObtenerItem(proyectoID, itemID)` y por la PK compuesta `(proyecto_id, numero)`: la propia firma hace imposible que el servicio pida un ítem "sin proyecto", por lo que el caso deja de tener valor como test de servicio y pasa a la suite de contrato (contra ambos adaptadores).

Nota sobre los 4 tests de fallo de escritura con `almacenFalloEscritura` (`TestUS02_Modificar_FalloDeEscrituraNoInformaExito`, `TestUS03_AgregarIntegrante_FalloDeEscrituraNoInformaExito`, `TestUS05_CrearItem_FalloDeEscrituraNoInformaExito`, `TestUS11_EstimarItem_FalloDeEscrituraNoInformaExito`): se conservan por el mismo mecanismo que el grupo de 16 (el doble se inyecta en el adaptador JSON; el `Save` fallido ocurre dentro del adaptador y el servicio devuelve el error envuelto, nunca éxito). No integran la lista de 18 y se listan aquí para completitud.

### D11. Alcance

Solo los agregados actuales: `proyectos`, `integrantes` e `items_backlog`. Sprints, votos, esfuerzo y defectos traen su propia migración con su propia historia. Fuera de alcance: web, nuevas historias, cambios de requisitos.

## 4. Alternativas descartadas

| Alternativa | Motivo del descarte |
|---|---|
| **Quedarse solo con JSONStore** | No resuelve nada de la sección 2: sin transacciones ni concurrencia (numeración racy), sin unicidad ni FK en el motor, cada operación reescribe el archivo entero, sin consultas para T-03/web. Además T-03/T-04 ya decidieron PostgreSQL (el `compose.yml` lo absorbe T-04). |
| **Implementar `almacen` literal sobre SQL** (cargar/guardar la colección entera por operación) | Mantendría la firma pero semánticamente pésimo: cada `Save` reescribiría todas las filas (DELETE masivo + re-INSERT o upsert de toda la colección), O(n) por operación, sin aprovechar índices, sin atomicidad por registro, y la numeración seguiría siendo racy (habría que serializar transacciones completas solo para calcular un máximo). El servicio seguiría filtrando colecciones en memoria. Es la peor de ambas mundos. |
| **ORM (GORM)** | El dominio ya valida todo en funciones puras; un ORM aportaría reflections, tags de persistencia en structs (empujando la tecnología al dominio o duplicando structs), traducción de errores menos controlable (necesaria para D3) y una curva de aprendizaje para un equipo chico. pgx v5 nativo da control exacto del SQL y de los errores. |
| **sqlc / codegen** | Genera el código de datos del lado del generador: las interfaces quedarían definidas por el storage y no por el consumidor, contradiciendo D1; además agrega un paso de build y obliga a ajustar las queries al contrato a mano de todos modos. Podría reconsiderarse en el futuro para consultas complejas de reportes. |
| **Un único repositorio monolítico** (una interfaz `Repositorio` con todo) | Acopla los tres agregados: todo consumidor dependería de todo, los dobles de test serían más pesados, y mezcla ciclos de vida (p. ej. el backlog crece independientemente de los integrantes). Rompe el diseño por agregado y la trazabilidad por historia. |
| **UUID como parte de `domain.Proyecto`** | El ID es un detalle de persistencia: hoy vive en la capa de servicio (decisión D11 de US-02, registros no exportados en `internal/service`). Meterlo en el dominio cambia US-01 ("se conserva `internal/domain`"), acopla el dominio a la estrategia de IDs y contamina las funciones puras (`NuevoProyecto`, `ModificarProyecto`) y sus 45 tests. Los contratos lo resuelven con `service.ProyectoRegistrado` (D2). |

## 5. Consecuencias

**Positivas:**

- Numeración de ítems atómica por proyecto (D4): sin IDs duplicados bajo concurrencia.
- Unicidad y FK garantizadas por el motor (defensa en profundidad), con traducción a errores de dominio (D3).
- Los servicios dejan de manipular colecciones: cada operación lee/escribe el registro que necesita.
- El dominio queda intacto y sigue siendo la única fuente de reglas; los contratos (D2) ya reservan las lecturas de US-04 y T-03.
- Suite de contrato compartida: paridad verificada entre el adaptador JSON (tests rápidos) y el de PostgreSQL (servidor).
- T-02 queda intacto: `internal/storage` no se modifica y sus 9 tests siguen corriendo.
- Camino despejado para los agregados futuros (mismo patrón, D11).
- Frontera de issues decidida (revisión 2026-11-11): T-03 posee el adaptador y la conexión, T-04 absorbe el `compose.yml` y T-10 es el único que modifica `.github/workflows/ci.yml`.

**Negativas / riesgos:**

- Dos adaptadores que mantener; toda corrección semántica debe duplicarse en ambos (mitigado por la suite de contrato).
- La paridad JSON/PostgreSQL descansa en la suite de contrato: los tests de servicio ejercitan servicio+JSON, no servicio+PostgreSQL.
- Las fechas se normalizan a día UTC en la frontera de persistencia (al guardar y al leer): la componente horaria no sobrevive al round-trip en ningún adaptador (el JSON actual sí la preserva). Decisión de la revisión 2026-11-11 (D6); la suite de contrato compara solo año, mes y día.
- `EqualFold` (dominio) y `lower()` (PostgreSQL) difieren en case-folding Unicode y en el recorte de espacios: el chequeo autoritativo es el dominio; la constraint puede dejar pasar un caso que el dominio rechaza (nunca al revés en la práctica), y la traducción de 23505 mantiene el error observable.
- Los constructores de servicio cambian de firma (`almacen` → repositorios): los 63 tests de servicio requieren una edición mecánica. Única excepción a "firmas públicas conservadas", aceptada como excepción documentada en la revisión del 2026-11-11 (los métodos públicos no cambian).
- La distinción `nil`/lista vacía de los criterios de aceptación no se preserva en el round-trip (D6, se normaliza a lista vacía): ajuste mecánico si algún test dependía de `null` vs `[]`.
- CI más lenta y con un punto de falla más si se corre la suite contra PostgreSQL (mitigado: se saltea sin `TEST_DATABASE_URL`).
- Ejecutar migraciones al arrancar `cmd/server` agrega un paso de arranque que puede fallar (falla rápida, aceptada).
- Más indirección en el código (puerto/adaptador) para un proyecto chico.

## 6. Plan de migración por pasos

Principio: `go test ./...` en verde al terminar cada paso (los pasos 1–5 son refactor puro: el comportamiento no cambia y los tests existentes lo prueban; el paso 6 agrega el adaptador nuevo detrás de una variable de entorno). Cada paso sigue el flujo del proyecto: spec/escenarios ya existen (los requisitos no cambian); los tests nuevos de cada paso se escriben primero en RED.

| Paso | Contenido | Verificación |
|---|---|---|
| 1 (T-05) | TDD: `domain.NumeroDeItem` + `domain.SiguienteIDItem` como funciones puras (mudanza de la lógica de `internal/service/backlog.go:63-96`); `CrearItem` pasa a usarlas. Sin cambio de comportamiento. | `go test ./...` verde; nuevos tests de dominio en verde. |
| 2 (T-06) | Contratos de D2 en `internal/service` (interfaces + `service.ProyectoRegistrado` exportado con los mismos json tags). Los servicios todavía usan `almacen` (compila sin cambiar nada más). | `go build ./...` y `go test ./...` verde (sin cambios funcionales). |
| 3 (T-07) | `internal/storage/jsonrepo`: adaptador JSON de los tres contratos sobre JSONStore (constructores reciben la interfaz `almacen` con Save/Load), conservando nombres de colección, formato y orden. Tests unitarios propios: propagación de `ErrCorrupt` envuelto, traducción de `ErrNotFound` a vacío, duplicados con `EsDuplicado`, asignación de ID con `SiguienteIDItem`. | `go test ./...` verde (tests nuevos en verde; los 63 de servicio intactos). |
| 4 (T-08) | Suite de contrato compartida corriendo contra `jsonrepo` (errores de dominio, vacío, orden, duplicados, numeración, alcance por proyecto). Incluye los 2 casos reubicados de D10. | `go test ./...` verde. |
| 5 (T-09) | Migración de los servicios a repositorios: nuevos constructores (D1), `ModificarProyecto` = `Obtener` → `domain.ModificarProyecto` → `Actualizar`; `AgregarIntegrante` = `Existe` → `domain.NuevoIntegrante` → `Agregar` (duplicado en el repo); `CrearItem` = `Existe` → `domain.NuevoItemBacklog` → `AgregarItem`; `EstimarItem` = `Existe` → `ObtenerItem` → `Estimar` → `ActualizarItem`. Se eliminan `almacen`, los registros no exportados y `verificarProyectoExistente`. Edición mecánica de los 63 tests (constructores + renombres de registros; expectativas de fechas comparadas solo por año/mes/día y de criterios aceptando lista vacía en lugar de `nil`); reubicación de los 2 tests de D10 a la suite de contrato (corren contra ambos adaptadores). | `go test ./...` verde completo; cobertura sin bajar; los 16+6 tests de D10 con aserciones intactas. |
| 6 (T-03, existente) | Adaptador `internal/storage/postgres` (pgx v5 + pgxpool; requiere PostgreSQL 13 o superior), DDL de D6 como primera migración goose embebida en `migrations/`, ejecución al arrancar `cmd/server`, `DATABASE_URL`; suite de contrato contra PostgreSQL y tests de integración (FK, CHECK, unicidad, `FOR UPDATE`) bajo `TEST_DATABASE_URL`. | `go test ./...` verde sin `TEST_DATABASE_URL` (se saltean); con `TEST_DATABASE_URL` la suite de contrato y de integración en verde. |
| 7 (T-04 + T-10) | `compose.yml` con PostgreSQL (T-04, existente; absorbe el compose) y extensión de `.github/workflows/ci.yml` (T-01) con servicio PostgreSQL + `TEST_DATABASE_URL` (T-10, único issue que modifica el workflow). | Pipeline verde con la suite completa. |

Después del paso 5, T-03 depende de este ADR para el adaptador PostgreSQL; los pasos 1–5 no dependen de T-03/T-04 y pueden ejecutarse antes.

## 7. (a) Tabla de impacto por historia

Ninguna spec se re-escribe: US-02, US-03, US-05 y US-11 solo necesitan una **nota** en su sección de persistencia (§12/§11) que referencie este ADR; US-01 y US-04 mantienen su exclusión explícita.

| Historia | Qué se conserva | Qué cambia | Nota en su spec (sección de persistencia) |
|---|---|---|---|
| **US-01** | `domain.NuevoProyecto` y sus validaciones intactos (45 tests de dominio no cambian); CA, reglas y escenarios intactos. | Nada en código hoy: sigue sin existir `CrearProyecto` (brecha B1). Cuando exista la historia de alta, el ID vendrá de `gen_random_uuid()` (D4). | No tiene sección de persistencia (la excluye explícitamente): la exclusión sigue siendo correcta. Nota opcional: la persistencia de proyectos queda cubierta por ADR-001/T-03 y por la historia de alta futura (B1). |
| **US-02** | CA, reglas, escenarios BDD; firma y comportamiento de `ModificarProyecto` (orden validación/existencia, `ErrProyectoNoEncontrado`, errores de dominio sin escribir). | Constructor del servicio (`almacen` → repositorios); internamente pasa de "guardar la colección completa" a `Obtener` + `Actualizar` del registro único; las fechas del proyecto se normalizan a día UTC en la frontera de persistencia (D6). | §12: nota indicando que D13 (colección completa) pasa a ser "registro único vía `RepositorioProyectos`", que D14 ("archivo inexistente = colección vacía") pasa al adaptador JSON (ADR-001 D3) y que D15 se mantiene como contrato (ADR-001 D3). Referencia a ADR-001. |
| **US-03** | CA, reglas, escenarios; firma y comportamiento de `AgregarIntegrante` (mismo `ErrIntegranteDuplicado`, validaciones antes de persistir). | Constructor del servicio; el chequeo de duplicado pasa de `EsDuplicado` en el servicio a los repositorios (mismo error observable, sin escritura al fallar); la unicidad queda además defendida por índice en PostgreSQL. | §12: nota indicando que el duplicado pasa al repositorio (ADR-001 D3/D6) con idéntico error y que el orden de inserción se conserva (D7). Referencia a ADR-001. |
| **US-04** | `EstadoPorFechas` y `PorcentajeCompletadas` intactas (dominio puro, sin persistencia); sus lecturas ya están previstas en los contratos (`Obtener`, `ListarPorProyecto` de backlog). | Nada hoy (sigue sin servicio). Cuando se implemente, usará los contratos existentes sin cambiarlos. | No tiene sección de persistencia (la excluye explícitamente): la exclusión sigue siendo correcta. Nota opcional: ADR-001 D2 ya reserva las lecturas que US-04 necesitará. |
| **US-05** | CA, reglas, escenarios; firma y comportamiento de `CrearItem` (regla "siguiente al mayor del proyecto", `HU-01` inicial). | Constructor del servicio; la asignación del ID pasa al repositorio vía `domain.SiguienteIDItem` (misma regla, atómica en PostgreSQL); los criterios de aceptación `nil` se leen como lista vacía (D6). | §11: nota indicando que la numeración pasa a `domain.SiguienteIDItem` ejecutada por los repositorios (ADR-001 D4) y que el orden por número se conserva (D7). Referencia a ADR-001. |
| **US-11** | CA, reglas, escenarios; firma y comportamiento de `EstimarItem` (validaciones de SP, conservación del valor previo al fallar). | Constructor del servicio; internamente pasa de regrabar la colección a `ObtenerItem` + `ActualizarItem` del registro único. | §11: nota indicando que la escritura pasa a ser por ítem (ADR-001 D2) y que los errores de dominio se conservan (D3). Referencia a ADR-001. |
| **T-01** | Workflow de CI (`go build` + `go test ./... -cover`); los tests de dominio y servicio siguen corriendo sin base de datos. | Extender `.github/workflows/ci.yml` con un servicio PostgreSQL y la variable `TEST_DATABASE_URL` (T-10, único issue que modifica el workflow) para correr la suite de contrato y de integración. | No aplica (issue técnico, no spec): el ajuste es en el propio workflow, referenciando ADR-001. |
| **T-02** | `JSONStore` y sus 9 tests **sin modificar** (`internal/storage` no se toca); escritura atómica con temporal+rename, `ErrNotFound`/`ErrCorrupt`/`ErrInvalidName` intactos. | Su rol: deja de ser la persistencia del servidor y pasa a ser el motor del adaptador `jsonrepo` (tests de servicio y suite de contrato JSON). | No aplica (issue técnico, ya implementado): comentario en el issue referenciando ADR-001 para dejar asentado el nuevo rol. |

## 8. (b) Brechas detectadas que requieren trabajo aparte

1. **No existe `CrearProyecto` ni listados en los servicios.** No hay operación de INSERT de proyectos: por lo tanto `gen_random_uuid()::text` (D4) **aún no tiene camino de ejecución** y los tests siguen sembrando los IDs a mano ("p1", "p2"). Se necesita una **historia de usuario nueva** (alta de proyecto) que agregue `ServicioProyectos.CrearProyecto` (o el nombre que decida la historia) y una operación de creación en `RepositorioProyectos` (p. ej. `Crear(proyecto domain.Proyecto) (ProyectoRegistrado, error)`), que en PostgreSQL haga INSERT sin ID para que aplique el DEFAULT y devuelva el registro con el ID generado. Los listados (`Listar`, `ListarPorProyecto`) tampoco tienen consumidor hoy: se reservan para T-03/web y US-04. Decidido (revisión 2026-11-11): ese diseño pertenece a la historia de alta, no a este ADR; el `DEFAULT gen_random_uuid()::text` queda en el esquema inicial como latente.
2. **Numeración/asignación del ID de proyecto: dónde y cómo se resuelve.** Hoy no se genera en ningún lado (auditoría). Propuesta: cuando exista la historia de alta, el ID lo genere el motor (`DEFAULT gen_random_uuid()::text`), lo devuelva el repositorio dentro de `ProyectoRegistrado` y lo use el servicio; el dominio no genera IDs. Alternativa descartable: ID secuencial tipo "P-01" con `SiguienteID` al estilo de ítems (requiere bloqueo y no aporta nada sobre UUID). Decisión: el diseño de `CrearProyecto` se decide en la historia de alta, no en este ADR; el `DEFAULT gen_random_uuid()::text` queda en el esquema inicial aunque no exista `CrearProyecto` (revisión 2026-11-11; cierra el pendiente 13).
3. **Migración de datos existentes (JSON → PostgreSQL) no contemplada.** No hay operación de importación. Si hay datos reales en los archivos `proyectos.json`/`integrantes.json`/`backlog.json` que deban sobrevivir, hace falta un script one-shot; si se parte de una base vacía, no. **PENDIENTE DE DECIDIR.**
4. **US-04 sigue sin servicio**, pero sus lecturas ya están en los contratos (D2): cuando se implemente no habrá que cambiar el contrato. Es trabajo de esa historia, no de este ADR.
5. **No existe ninguna operación de borrado** de proyectos/integrantes/ítems, así que el comportamiento de las FK ante DELETE (cascada vs. restricción) no tiene consecuencias observables hoy: el DDL usa el NO ACTION por defecto. Se decidirá con la primera historia que borre algo (migración nueva, D8).

## 9. (c) Issues técnicos propuestos (T-05 en adelante)

T-03 (adaptador PostgreSQL + esquema + conexión de `cmd/server`) y T-04 (`compose.yml` con PostgreSQL) **ya existen** y no se duplican ni se redefinen; la lista siguiente cubre el trabajo que este ADR agrega antes y después de ellos. Frontera decidida (revisión 2026-11-11; cierra el pendiente 15): T-04 absorbe el `compose.yml` y T-10 es el único issue que modifica `.github/workflows/ci.yml`; ninguna otra tarea de esta lista toca el workflow ni el compose.

| Issue | Título | Contenido | Dependencias |
|---|---|---|---|
| **T-05** | Regla de numeración de ítems como función pura de dominio | `domain.NumeroDeItem` + `domain.SiguienteIDItem` con TDD (mudanza desde `internal/service/backlog.go`); `CrearItem` pasa a usar la función de dominio; sin cambio de comportamiento. | Ninguna (base T-01). |
| **T-06** | Contratos de repositorio por agregado en `internal/service` | Interfaces `RepositorioProyectos`, `RepositorioIntegrantes`, `RepositorioBacklog` y `service.ProyectoRegistrado` según D2 (firmas exactas de este ADR). | T-05 (la documentación de `AgregarItem` referencia `SiguienteIDItem`). |
| **T-07** | Adaptador JSON de repositorios (`internal/storage/jsonrepo`) | Implementa los tres contratos sobre `JSONStore` sin modificar `internal/storage`; conserva colecciones, formato y orden; constructores reciben la interfaz Save/Load (para dobles); propagación de `ErrCorrupt` y traducción de `ErrNotFound` (D3). Tests unitarios propios. | T-06. |
| **T-08** | Suite de contrato compartida de repositorios | Casos de contrato (errores de dominio, vacío, orden, duplicados, numeración, alcance por proyecto) corriendo contra `jsonrepo`; recibe los 2 tests reubicados de D10, que corren contra ambos adaptadores (comparaciones de fechas solo por año/mes/día y criterios como lista vacía, D6). Diseñada para correr luego contra el adaptador PostgreSQL (helpers de la suite compartidos entre los tests de ambos adaptadores). | T-07. |
| **T-09** | Migración de los servicios a repositorios | Constructores nuevos (D1), eliminación de `almacen`/registros no exportados/`verificarProyectoExistente`, edición mecánica de los 63 tests, reubicación de los 2 tests (D10). `go test ./...` verde. | T-07 y T-08. |
| **T-03** (existente) | Adaptador PostgreSQL | Implementa los tres contratos con pgx v5 + `pgxpool` (requiere PostgreSQL 13 o superior), DDL de D6 como primera migración goose embebida, ejecución al arrancar `cmd/server`, `DATABASE_URL`; suite de contrato contra PostgreSQL y tests de integración (FK, CHECK, unicidad, `FOR UPDATE`). **No se crea: ya existe.** | Depende de T-06 (contratos), T-07 (patrón de adaptador) y T-08 (suite). |
| **T-04** (existente) | `compose.yml` con PostgreSQL | Contiene el servicio PostgreSQL y la versión (PostgreSQL 13 o superior, requisito del ADR, habilita `gen_random_uuid()` nativo); absorbe el `compose.yml` y no toca `.github/workflows/ci.yml`. **No se crea: ya existe.** | Depende de T-03 para la definición del esquema/versión. |
| **T-10** | CI con PostgreSQL para la suite de contrato e integración | Extender `.github/workflows/ci.yml` de T-01 con un servicio PostgreSQL y `TEST_DATABASE_URL` (se saltea sin la variable). **Único issue que modifica el workflow** (frontera decidida, revisión 2026-11-11). | T-03, T-04, T-08. |
| **T-11** | Actualización de documentación de agentes | Aplicar los cambios de la sección 10 (AGENTS.md y `.opencode/agents/*.md`) cuando la arquitectura esté real (después de T-09). | T-09. |

## 10. (d) Cambios necesarios en AGENTS.md y en los agentes de `.opencode/agents/*.md`

Lista, **sin aplicar** (los textos exactos se adaptan al formato de cada archivo cuando se ejecuten):

En `AGENTS.md`:

1. Regla 1 (lenguaje): sin cambios (sigue siendo Go).
2. Regla 7: extenderla — la lógica de negocio vive en `internal/domain`; la capa web no duplica reglas; **los contratos de repositorio los define el consumidor (`internal/service`) y los adaptadores viven en `internal/storage/{jsonrepo,postgres}`**; el dominio no conoce la persistencia.
3. Bloque "Estructura": agregar `internal/storage/{jsonrepo,postgres}` y `migrations/`; actualizar la línea de persistencia de "Persistencia: archivos JSON" a "Persistencia: PostgreSQL 13 o superior (pgx v5, migraciones goose embebidas); adaptador JSON (`jsonrepo`) para tests; repositorios por agregado definidos en `internal/service` (ADR-001)".
4. Agregar a "Reglas" o "Estructura": variables de entorno `DATABASE_URL` (servidor) y `TEST_DATABASE_URL` (tests de integración y suite de contrato; se saltean sin la variable).
5. Agregar a "Reglas": las migraciones ya aplicadas nunca se editan; los cambios de esquema son migraciones nuevas (D8).
6. Agregar a "Reglas" o "Commits"/flujo: los tests de repositorio siguen la suite de contrato compartida (ambos adaptadores); referenciar `docs/adr/ADR-001-persistencia-postgresql.md` como decisión vigente de persistencia. Agregar también que `.github/workflows/ci.yml` solo lo modifica T-10 y que el `compose.yml` lo define T-04 (frontera decidida, revisión 2026-11-11).

En `.opencode/agents/*.md` (por agente; sin leer los archivos aún, se listan los cambios de contenido):

7. `gestor-backlog`: incorporar los issues T-05–T-11 y sus dependencias; en las historias nuevas, recordar que la sección de persistencia referencia ADR-001 y que los agregados nuevos traen su propia migración (D11). Respeta la frontera decidida: el `compose.yml` pertenece a T-04 y `.github/workflows/ci.yml` solo lo modifica T-10.
8. `especificador`: la sección de persistencia de las specs debe describir **qué** se persiste y qué errores de dominio produce el repositorio, no **cómo** (sin colecciones JSON ni DDL); citar ADR-001; para US-02/US-03/US-05/US-11, agregar la nota de la sección 7 de este ADR.
9. `desarrollo-tdd` (desarrollador-tdd): los helpers de servicio construyen `jsonrepo` (no `JSONStore` directo); los dobles `almacenCuentaSaves`/`almacenFalloEscritura` se inyectan en el adaptador JSON; la suite de contrato corre contra ambos adaptadores (comparaciones de fechas solo por año/mes/día y criterios como lista vacía, D6); los tests de integración PostgreSQL se saltean sin `TEST_DATABASE_URL`; los tests de contrato/adapter nuevos van primero en RED.
10. `revisor`: agregar al checklist — el dominio no conoce la persistencia; las reglas de dominio no se duplican en los adaptadores; la suite de contrato corre contra ambos adaptadores; no se editan migraciones aplicadas; los tests clasificados en D10 están donde corresponde.
11. `orquestador`: coordinar el orden T-05 → T-06 → T-07 → T-08 → T-09 → T-03 → T-04/T-10 → T-11 con `go test ./...` verde en cada frontera.

## 11. Pendientes

Revisión del 2026-11-11: las correcciones 1–7 aprobadas por la persona quedaron incorporadas como decisiones (primera subsección, indicando qué pendiente cerró cada una); los que siguen abiertos se renumeran al final. No se supuso nada.

### Decidido en la revisión del 2026-11-11 (cerraron los pendientes 1, 2, 8, 9, 10, 11, 13 y 15 de la numeración original)

1. **Constructores de servicio (cerró el 1):** el cambio de firma de `NuevoServicioProyectos`/`NuevoServicioBacklog` (`almacen` → repositorios) se acepta como excepción documentada a "firmas públicas conservadas"; los métodos públicos no cambian.
2. **ID de proyecto (cerró el 2):** viaja en `service.ProyectoRegistrado{ID, Proyecto}` exportado (mismos json tags); `domain.Proyecto` no se modifica.
3. **`gen_random_uuid()` y versión de PostgreSQL (cerraron el 10 y el 13):** el `DEFAULT gen_random_uuid()::text` queda en el esquema inicial aunque no exista `CrearProyecto` (latente hasta la historia de alta, brechas B1/B2); se exige PostgreSQL 13 o superior y la versión de imagen la fija T-04 en el `compose.yml`. El diseño de `CrearProyecto` pertenece a la historia de alta, no a este ADR.
4. **`criterios_aceptacion` (cerró el 8):** `text[] NOT NULL DEFAULT '{}'`; el adaptador convierte `nil` en lista vacía al guardar y al leer; la distinción `null`/`[]` no se preserva.
5. **Frontera de issues (cerró el 15):** T-04 absorbe el `compose.yml` con PostgreSQL y T-10 es el único issue que modifica `.github/workflows/ci.yml`; ninguna otra tarea toca el workflow.
6. **Fechas (cerró el 9):** se normalizan a fecha UTC al guardar y al leer (día UTC); los tests de la suite de contrato comparan solo año, mes y día.
7. **Tests reubicados (cerró el 11):** `TestUS11_EstimarItem_BacklogSinArchivo` y `TestUS11_EstimarItem_MismoIDEnOtroProyecto` corren contra el adaptador JSON y contra PostgreSQL en la suite de contrato, sin reducir cobertura.

### Pendientes de decisión (abiertos)

1. **(antes 3) D2 — firmas menores:** `Actualizar`/`ActualizarItem` devuelven solo `error` (el servicio ya tiene el valor modificado); `Agregar` de integrantes devuelve `error` (no hay ID que devolver). Confirmar o ajustar.
2. **(antes 4) D2 — paquete de los contratos:** se decide `internal/service` (D1); confirmar que la dirección de dependencia (adaptadores importan `internal/service`) y el ciclo solo-de-test (tests de servicio importando `jsonrepo`) no presentan problemas en T-09.
3. **(antes 5) D4 — firma exacta de `domain.SiguienteIDItem`/`NumeroDeItem`:** propuesta `SiguienteIDItem(idsExistentes []string) string`. Confirmar (p. ej. ¿recibe IDs o números?).
4. **(antes 6) D6 — CHECK de `proyectos.estado`:** hoy solo se persiste `'Planificado'`. Confirmar el conjunto completo de valores de `domain.Estado` leyendo `internal/domain/estado.go` (fuera de los archivos permitidos de este ADR) y decidir si el CHECK los incluye desde ya o se extiende después con migración nueva.
5. **(antes 7) D6 — materialización del esquema inicial:** primera migración goose (recomendado: la aplica `cmd/server` y la reutilizan tests/CI) vs. script de inicialización del compose de T-04. PENDIENTE DE DECIDIR.
6. **(antes 12) D10 — alcance de los tests de concurrencia:** decidir si el `FOR UPDATE` de D4 se prueba con un test de concurrencia real (dos goroutines) o solo con el test de integración de PK compuesta.
7. **(antes 14) Sección 8/B3 — migración de datos JSON existentes:** script one-shot vs. base vacía. PENDIENTE DE DECIDIR.

## Criterios sin escenario / verificación

Este documento es un ADR, no una historia de usuario: no genera `specs/ID.md` ni `features/ID.feature`, y sus decisiones no se verifican con escenarios BDD sino con los issues T-05+ (tests de contrato e integración, sección 9) y con la tabla de impacto (sección 7). Las correcciones 1–7 de la revisión del 2026-11-11 ya están incorporadas (sección 11). Lo que la persona debe revisar antes de la aceptación final:

1. Las firmas exactas de D2 (interfaces, `ProyectoRegistrado`); los constructores ya están aceptados como excepción documentada (sección 11): son la base de todos los issues T-05–T-09.
2. El DDL de D6 (columnas — incluida `criterios_aceptacion` —, fechas `date` normalizadas a día UTC, CHECKs, índice único de integrantes, PK compuesta de ítems) y sus pendientes abiertos (sección 11).
3. La clasificación de tests de D10 (6 corruptos conservados; 16 conservados + 2 reubicados de los 18 que cuentan Save; nota sobre los 4 de fallo de escritura), con los 2 reubicados corriendo contra ambos adaptadores.
4. Los issues T-05–T-11 y su frontera con T-03/T-04 existentes (frontera decidida en la sección 11: T-10 es el único que modifica el workflow).
5. Las brechas de la sección 8 (crear proyecto, ID de proyecto, migración de datos) y los pendientes abiertos de la sección 11.
