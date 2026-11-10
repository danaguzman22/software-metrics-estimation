#!/usr/bin/env bash
# Crea los issues T-03 (servidor web mínimo) y T-04 (Docker y despliegue).
# Uso:  cd <carpeta-del-repo>
#       PROJECT="Nombre de tu GitHub Project" bash crear-issues-T03-T04.sh
# Si no definís PROJECT, los issues se crean igual y los agregás al tablero a mano.
set -u

PROJECT="${PROJECT:-}"

# Labels (si ya existen, no pasa nada)
for l in "tecnica:5319E7" "prioridad:alta:B60205" "sprint-2:BFD4F2"; do
  name="${l%:*}"; color="${l##*:}"
  gh label create "$name" --color "$color" 2>/dev/null || true
done

# us ID "Título" prioridad SP sprint "descripción" "criterios (uno por línea)" tipo
us() {
  local id="$1" title="$2" prio="$3" sp="$4" sprint="$5" story="$6" crit="$7" kind="${8:-historia}"
  local body="$story

**Prioridad:** ${prio}
**Story Points:** ${sp}
**Sprint planificado:** ${sprint}

### Criterios de aceptación
$(echo "$crit" | sed 's/^/- [ ] /')

### Trazabilidad
- [ ] Especificación SDD (specs/${id}.md)
- [ ] Escenarios BDD (features/${id}.feature)
- [ ] Tests (RED -> GREEN -> REFACTOR)
- [ ] Código Go
"
  local args=(--title "[$id] $title" --label "$kind" --label "prioridad:${prio,,}" --label "sprint-${sprint#S}" --body "$body")
  [ -n "$PROJECT" ] && args+=(--project "$PROJECT")
  gh issue create "${args[@]}"
}

us T-03 "Servidor web mínimo y plantilla base" Alta 5 S2 \
"Dejar la aplicación funcionando en el navegador con las historias ya implementadas (US-01, US-02, US-03 y US-05), para tener un primer MVP visible. Las pantallas usan solo la lógica de los servicios existentes y no duplican reglas de negocio. Si falta una consulta mínima (por ejemplo listar proyectos, integrantes o ítems), se agrega en esta tarea con su test." \
"\`go run ./cmd/server\` arranca el servidor; el puerto se toma de la variable PORT (por defecto 8080)
Existe una plantilla base HTML compartida con navegación y la página de inicio responde 200
Se puede listar y crear proyectos desde un formulario web
Desde el detalle de un proyecto se pueden agregar integrantes y cargar ítems del backlog, mostrando los ya registrados
Un error de validación se muestra en el formulario con un mensaje claro, sin pánico ni página en blanco
Los handlers tienen tests con net/http/httptest
Los datos se guardan en archivos JSON dentro de un directorio configurable (variable DATA_DIR)" tecnica

us T-04 "Dockerfile, Docker Compose y despliegue" Alta 3 S2 \
"Empaquetar la aplicación en una imagen Docker, poder levantarla con Docker Compose y desplegarla en una URL pública (Render o Azure). Depende de T-03." \
"\`docker build\` genera una imagen que arranca la aplicación (build multi-stage con binario Go)
\`docker compose up\` levanta la app y los datos JSON persisten en un volumen entre reinicios
La app lee PORT y DATA_DIR de variables de entorno
Con la variable SEED_DEMO=true, si no hay datos, se cargan datos de demostración (proyecto, integrantes y backlog)
La aplicación queda desplegada en una URL pública y el README explica cómo desplegarla
El README documenta la limitación de persistencia del plan gratuito (los datos se reinician) y los costos a controlar si se usa Azure" tecnica

echo "Listo. Revisá los issues en GitHub y completá Story Points, Prioridad y Sprint en el tablero."
