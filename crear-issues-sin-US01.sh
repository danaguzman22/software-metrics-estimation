#!/usr/bin/env bash
# Crea el Product Backlog inicial como issues de GitHub.
# (US-01 ya fue cargado a mano; este script crea el resto)
# Uso:  gh auth login   (una sola vez)
#       cd <carpeta-del-repo>
#       PROJECT="Nombre de tu GitHub Project" bash crear-issues.sh
# Si no definís PROJECT, los issues se crean igual y los agregás al tablero a mano.
set -u

PROJECT="${PROJECT:-}"

# ---------- Labels ----------
for l in "historia:1D76DB" "tecnica:5319E7" "documentacion:0E8A16" \
         "prioridad:alta:B60205" "prioridad:media:FBCA04" "prioridad:baja:C2E0C6" \
         "sprint-1:BFD4F2" "sprint-2:BFD4F2" "sprint-3:BFD4F2" "sprint-4:BFD4F2"; do
  name="${l%:*}"; color="${l##*:}"
  gh label create "$name" --color "$color" 2>/dev/null || true
done

# us ID "Título" prioridad SP sprint "Como... quiero... para..." "criterios (uno por línea)"
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

# ================= SPRINT 1 — MVP (núcleo de dominio) =================
us US-02 "Modificar proyecto" Media 2 S1 \
"**Como** integrante, **quiero** modificar los datos de un proyecto, **para** corregir o actualizar su información." \
"Se pueden editar nombre y fechas respetando las mismas validaciones que al crear
Modificar un proyecto inexistente devuelve error"

us US-03 "Registrar integrantes del proyecto" Alta 2 S1 \
"**Como** integrante, **quiero** registrar a los miembros del equipo con su rol, **para** asignarles esfuerzo y votos." \
"Se agrega un integrante con nombre y rol
No se permiten integrantes duplicados en un mismo proyecto
Nombre vacío devuelve error"

us US-04 "Consultar estado de un proyecto" Media 2 S1 \
"**Como** Product Owner, **quiero** ver el estado de un proyecto, **para** saber cuánto avanzó." \
"Muestra estado (Planificado / En curso / Finalizado), sprint actual y % de historias completadas
Un proyecto sin historias muestra 0% sin errores (división por cero)"

us US-05 "Crear ítem del Product Backlog" Alta 3 S1 \
"**Como** Product Owner, **quiero** cargar historias en el backlog, **para** organizar el trabajo." \
"Cada ítem tiene ID, título, descripción, prioridad, estado, story points y criterios de aceptación
El ID se genera automáticamente y es único
Título vacío o prioridad inválida devuelve error
Estado inicial: \"Pendiente\""

us US-06 "Editar y priorizar el backlog" Media 3 S1 \
"**Como** Product Owner, **quiero** editar y reordenar historias por prioridad, **para** mantener el backlog actualizado." \
"Se puede editar cualquier campo de un ítem
El backlog se lista ordenado por prioridad
No se puede editar una historia ya completada"

us US-07 "Crear Sprint con Sprint Goal" Alta 3 S1 \
"**Como** equipo, **quiero** crear un Sprint con su objetivo y fechas, **para** planificar una iteración." \
"Se crea con número, Sprint Goal, fecha de inicio y fin
Sprint Goal vacío devuelve error
No puede haber dos Sprints activos a la vez en un proyecto"

us US-08 "Asignar historias a un Sprint" Alta 3 S1 \
"**Como** equipo, **quiero** asignar historias del backlog a un Sprint, **para** definir el Sprint Backlog." \
"Solo se asignan historias estimadas y pendientes
Una historia no puede estar en dos Sprints activos
Los Story Points planificados del Sprint se calculan como suma de sus historias"

us US-09 "Registrar historias completadas y cerrar Sprint" Alta 5 S1 \
"**Como** equipo, **quiero** marcar historias como completadas y cerrar el Sprint, **para** medir lo entregado." \
"Al cerrar, las historias no completadas vuelven al backlog
Se calculan Story Points completados y velocidad del Sprint
Un Sprint cerrado no se puede modificar
Cerrar un Sprint ya cerrado devuelve error"

us US-11 "Estimar historias con Story Points" Alta 2 S1 \
"**Como** equipo, **quiero** asignar Story Points a una historia, **para** dimensionar el trabajo." \
"Solo se aceptan valores de la escala Fibonacci (1, 2, 3, 5, 8, 13, 21)
Un valor fuera de la escala devuelve error"

us US-16 "Métricas de Story Points y velocidad" Alta 5 S1 \
"**Como** Product Owner, **quiero** ver SP planificados, SP completados y velocidad, **para** predecir cuánto puede entregar el equipo." \
"SP planificados = suma de SP de las historias del Sprint
SP completados = suma de SP de historias en estado Completada
Velocidad = promedio de SP completados de los Sprints cerrados
Sin Sprints cerrados la velocidad es 0 (sin error)
% de historias completadas calculado correctamente con 0 historias"

us T-01 "Estructura del proyecto Go, go.mod y CI con go test" Alta 3 S1 \
"Dejar armada la base técnica: módulo Go, carpetas por capa (dominio / servicio / persistencia / web) y workflow de GitHub Actions que corra \`go test ./... -cover\`." \
"El repo compila con \`go build ./...\`
\`go test ./...\` corre en GitHub Actions en cada push
La lógica de negocio vive en un paquete separado de la capa web" tecnica

us T-02 "Persistencia en archivos JSON" Alta 3 S1 \
"Guardar y cargar proyectos, backlog, sprints, votos, esfuerzo y defectos en archivos JSON." \
"Los datos sobreviven al reiniciar la aplicación
Un archivo inexistente inicia con datos vacíos
Un JSON corrupto devuelve error claro, sin pánico" tecnica

# ================= SPRINT 2 — INTERFAZ + ESTIMACIÓN =================
us US-22 "Interfaz web base" Alta 8 S2 \
"**Como** usuario, **quiero** una interfaz web simple, **para** usar la aplicación desde el navegador." \
"Navegación entre proyectos, backlog, sprints y defectos
Formularios con mensajes de error claros
La interfaz usa la lógica del dominio, sin duplicar reglas"

us US-12 "Planning Poker: votos ocultos y revelación" Alta 5 S2 \
"**Como** equipo, **quiero** votar en secreto y revelar los votos juntos, **para** evitar el sesgo de anclaje." \
"Cada integrante registra un voto por ronda
Los votos permanecen ocultos hasta que se revelan
Un integrante no puede votar dos veces en la misma ronda (puede cambiar su voto antes de revelar)
Solo votan integrantes del proyecto"

us US-13 "Planning Poker: diferencias, nuevas rondas y estimación acordada" Alta 5 S2 \
"**Como** equipo, **quiero** detectar diferencias entre votos y repetir rondas hasta acordar, **para** fijar una estimación consensuada." \
"Se detecta diferencia cuando hay votos distintos (se informa mínimo, máximo y quién votó cada uno)
Se puede iniciar una nueva ronda conservando el historial
Se registra la estimación acordada y se guarda como Story Points de la historia
No se puede acordar sin revelar antes"

us US-14 "Registro de esfuerzo" Alta 3 S2 \
"**Como** integrante, **quiero** registrar horas trabajadas en una historia, **para** comparar esfuerzo estimado con real." \
"Se registra integrante, fecha, actividad y horas
Horas menores o iguales a 0 devuelven error
Solo se registra para historias existentes y integrantes del proyecto"

us US-17 "Métricas de horas y desviación" Alta 3 S2 \
"**Como** Product Owner, **quiero** ver horas estimadas, reales y desviación, **para** evaluar la precisión de las estimaciones." \
"Horas reales = suma de los registros de esfuerzo
Desviación = (reales - estimadas) / estimadas, expresada en porcentaje
Con horas estimadas = 0 no hay división por cero
Se calcula por historia, por Sprint y por proyecto"

us US-10 "Consultar Sprints anteriores" Baja 2 S2 \
"**Como** equipo, **quiero** ver el histórico de Sprints cerrados, **para** analizar la evolución." \
"Lista Sprints cerrados con Sprint Goal, SP planificados/completados y velocidad
Sin Sprints cerrados muestra una lista vacía"

# ================= SPRINT 3 — FUNCIONALIDAD Y CALIDAD =================
us US-15 "Gestión de defectos" Media 5 S3 \
"**Como** equipo, **quiero** registrar y seguir defectos, **para** medir la calidad." \
"Se registra descripción, severidad, estado, historia relacionada, Sprint de detección y Sprint de resolución
El Sprint de resolución no puede ser anterior al de detección
Un defecto resuelto requiere Sprint de resolución
Severidad y estado solo admiten valores válidos"

us US-18 "Métricas de defectos" Media 3 S3 \
"**Como** Product Owner, **quiero** ver defectos detectados y resueltos, **para** conocer la calidad del producto." \
"Cantidad de defectos detectados y resueltos por Sprint y por proyecto
Sin defectos devuelve 0"

us US-19 "Dashboard con gráficos" Media 8 S3 \
"**Como** Product Owner, **quiero** un dashboard con el estado del proyecto, **para** decidir rápido." \
"Muestra estado del proyecto y métricas clave
Incluye al menos 2 gráficos (por ejemplo velocidad por Sprint y horas estimadas vs reales)
Funciona con un proyecto sin datos"

us US-20 "Reporte de proyecto o Sprint" Media 5 S3 \
"**Como** Product Owner, **quiero** un reporte de un proyecto o Sprint, **para** presentarlo al cliente." \
"Incluye historias planificadas y completadas, estimaciones, esfuerzo, métricas y defectos
Se puede generar por Sprint o por proyecto completo"

# ================= SPRINT 4 — CIERRE =================
us US-21 "Exportar informe a PDF" Media 5 S4 \
"**Como** Product Owner, **quiero** exportar el reporte a PDF, **para** entregarlo." \
"Se genera un archivo PDF con el contenido del reporte
Un reporte sin datos genera un PDF válido"

us D-01 "Documentación técnica y manual breve de usuario" Media 3 S4 \
"Documentar arquitectura, cómo correr y testear, y un manual breve de uso." \
"README con instrucciones de instalación y ejecución
Manual de usuario con las funciones principales" documentacion

us D-02 "Informe de métricas y cobertura de pruebas" Media 2 S4 \
"Generar el informe de cobertura (\`go test -coverprofile\`) y las métricas del proyecto." \
"Informe con porcentaje de cobertura por paquete
Evidencia de TDD referenciada desde el historial de commits" documentacion

us D-03 "Actas de retrospectivas y evidencia de TDD" Media 2 S4 \
"Registrar las actas de cada retrospectiva y los commits RED/GREEN/REFACTOR." \
"Un acta por Sprint en docs/retros/
Commits con prefijos test:/feat:/refactor: que muestren el ciclo" documentacion

echo "Listo. Revisá los issues en GitHub y armá el tablero."
