# language: es
# Estado: Propuesta
# Los escenarios describen el comportamiento acordado para CA1–CA3 y las
# decisiones D1–D20: una interfaz web renderizada en el servidor, con
# navegación entre proyectos, backlog, sprints y defectos; formularios que
# muestran errores claros; y una interfaz que delega las reglas en el dominio
# y los servicios.
# FASE 1 incluye inicio, proyectos, backlog y páginas "Disponible próximamente"
# para Sprints y Defectos. FASE 2 incluye el listado, el alta y la asignación de
# Sprints; ya puede implementarse porque US-07, US-08 y US-09 están mergeadas en
# develop.
# Los errores se describen por su significado; sus identificadores y códigos
# HTTP exactos están decididos en la spec (sección 8).
# Los escenarios de FASE 2 están marcados con @Fase2 y ya pueden implementarse.
# No se escenifican el HTTP 500, el HTTP 405, el patrón POST–Redirect–GET, la
# delegación en servicios, la verificación estructural de imports, la
# generación de IDs de proyecto ni el ítem completado: se verifican con tests
# sin escenario BDD propio.

@US-22
Característica: Interfaz web base
  Como equipo
  Quiero una interfaz web para gestionar proyectos y su backlog
  Para operar el sistema sin usar la línea de comandos

  @CA1 @Fase1
  Escenario: La navegación principal está presente en todas las páginas
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Cuando se visita la página de inicio
    Y se visita la página de proyectos
    Y se visita el detalle del proyecto "PRJ-001"
    Y se visita el backlog del proyecto "PRJ-001"
    Y se visita la sección de sprints del proyecto "PRJ-001"
    Y se visita la sección de defectos
    Entonces todas las páginas muestran los enlaces "Proyectos", "Backlog", "Sprints" y "Defectos"
    Y las páginas de inicio y de proyectos muestran la leyenda "Elegí un proyecto para ver su backlog y sus sprints"

  @CA1 @Fase1
  Escenario: La página de inicio ofrece acceder a los proyectos
    Cuando se visita la página de inicio
    Entonces se muestra el título "Software Metrics & Estimation"
    Y se ofrece un enlace para listar los proyectos
    Y se ofrece un enlace para crear un proyecto

  @CA1 @Fase1
  Escenario: La sección de defectos informa que estará disponible próximamente
    Cuando se visita la sección de defectos
    Entonces se muestra el mensaje "Disponible próximamente"
    Y se informa que la gestión de defectos se incorporará más adelante

  @CA1 @Fase1
  Escenario: En la primera fase, la sección de sprints informa que estará disponible próximamente
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Cuando se visita la sección de sprints del proyecto "PRJ-001"
    Entonces se muestra el mensaje "Disponible próximamente"
    Y se informa que la gestión de sprints se incorporará en una próxima entrega

  @CA1 @Fase1
  Escenario: Listar los proyectos registrados
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-03-01" y fecha de fin "2026-06-30"
    Y un proyecto existente identificado por "PRJ-002", con nombre "Portal de Clientes", fecha de inicio "2026-04-15" y fecha de fin "2026-08-31"
    Cuando se visita la página de proyectos
    Entonces se listan los proyectos "Sistema de Métricas" y "Portal de Clientes"
    Y cada proyecto muestra su identificador, sus fechas y su estado "Planificado"
    Y "Sistema de Métricas" aparece antes que "Portal de Clientes"

  @CA1 @Fase1
  Escenario: La lista de proyectos vacía muestra un mensaje
    Dado que no hay proyectos registrados
    Cuando se visita la página de proyectos
    Entonces se muestra el mensaje "No hay proyectos registrados"
    Y no se muestra ningún error

  @CA1 @Fase1
  Escenario: Ver el detalle de un proyecto
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-03-01" y fecha de fin "2026-06-30"
    Cuando se visita el detalle del proyecto "PRJ-001"
    Entonces se muestra el nombre "Sistema de Métricas"
    Y se muestran las fechas "2026-03-01" y "2026-06-30"
    Y se muestra el estado "Planificado"

  @CA1 @Fase1
  Escenario: El estado de un proyecto se calcula por sus fechas
    Dado que hoy es "1 de febrero de 2026"
    Y un proyecto existente identificado por "PRJ-003", con nombre "Sistema legado", fecha de inicio "2025-01-01" y fecha de fin "2025-12-31"
    Cuando se visita el detalle del proyecto "PRJ-003"
    Entonces se muestra el nombre "Sistema legado"
    Y se muestra el estado "Finalizado"
    Cuando se visita la página de proyectos
    Entonces el proyecto "Sistema legado" muestra el estado "Finalizado"

  @CA1 @Fase1
  Escenario: Crear un proyecto con datos válidos
    Dado que hoy es "1 de febrero de 2026"
    Y que no hay proyectos registrados
    Cuando se completa el alta de proyecto con nombre "Sistema de Métricas", fecha de inicio "2026-03-01" y fecha de fin "2026-06-30"
    Y se envía el formulario
    Entonces se muestra el detalle del proyecto creado
    Y el proyecto "Sistema de Métricas" aparece en la lista de proyectos con estado "Planificado"

  @CA1 @CA3 @Fase1
  Escenario: Listar el backlog ordenado por prioridad
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión" y prioridad "Baja"
    Y que el proyecto "PRJ-001" tiene la historia "HU-02" con título "Cerrar sesión" y prioridad "Alta"
    Y que el proyecto "PRJ-001" tiene la historia "HU-03" con título "Recuperar contraseña" y prioridad "Media"
    Cuando se visita el backlog del proyecto "PRJ-001"
    Entonces las historias se muestran en el orden "Cerrar sesión", "Recuperar contraseña", "Iniciar sesión"
    Y cada historia muestra su identificador, su prioridad, su estado y sus story points

  @CA1 @Fase1
  Escenario: Un proyecto sin ítems muestra un mensaje
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" no tiene historias en el backlog
    Cuando se visita el backlog del proyecto "PRJ-001"
    Entonces se muestra el mensaje "El backlog está vacío"
    Y no se muestra ningún error

  @CA1 @Fase1
  Escenario: Crear un ítem del backlog
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Cuando se completa el alta de historia con título "Registrar horas", descripción "Registrar las horas trabajadas", prioridad "Media" y los criterios "La hora debe ser válida" y "El registro queda guardado"
    Y se envía el formulario
    Entonces el backlog del proyecto "PRJ-001" muestra la historia "Registrar horas"
    Y la historia figura con estado "Pendiente" y 0 story points
    Y la historia muestra los dos criterios de aceptación ingresados

  @CA1 @Fase1
  Escenario: Editar un ítem del backlog conservando su estimación
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión" y 3 Story Points
    Cuando se abre el formulario de edición de la historia "HU-01"
    Entonces el formulario muestra el título "Iniciar sesión" y los 3 Story Points actuales
    Cuando se cambia el título a "Iniciar sesión con correo" y se envía el formulario
    Entonces el backlog del proyecto "PRJ-001" muestra la historia "Iniciar sesión con correo"
    Y la historia conserva sus 3 Story Points

  @CA1 @Fase1
  Escenario: Estimar un ítem del backlog
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión" y 0 Story Points
    Cuando se estima la historia "HU-01" con 5 Story Points
    Entonces el backlog del proyecto "PRJ-001" muestra la historia "Iniciar sesión" con 5 Story Points

  @CA2 @Fase1
  Escenario: No se crea un proyecto con nombre vacío
    Dado que no hay proyectos registrados
    Cuando se completa el alta de proyecto con nombre vacío, fecha de inicio "2026-03-01" y fecha de fin "2026-06-30"
    Y se envía el formulario
    Entonces se muestra el mensaje "el nombre del proyecto no puede estar vacío"
    Y el formulario conserva las fechas ingresadas
    Y no se crea ningún proyecto

  @CA2 @Fase1
  Escenario: No se crea un proyecto con fecha de fin anterior a la de inicio
    Dado que no hay proyectos registrados
    Cuando se completa el alta de proyecto con nombre "Sistema de Métricas", fecha de inicio "2026-06-30" y fecha de fin "2026-03-01"
    Y se envía el formulario
    Entonces se muestra el mensaje "la fecha de fin no puede ser anterior a la fecha de inicio"
    Y el formulario conserva el nombre ingresado
    Y no se crea ningún proyecto

  @CA2 @Fase1
  Escenario: No se acepta una fecha con formato inválido
    Dado que no hay proyectos registrados
    Cuando se completa el alta de proyecto con nombre "Sistema de Métricas", fecha de inicio "01/03/2026" y fecha de fin "2026-06-30"
    Y se envía el formulario
    Entonces se muestra el mensaje "La fecha debe tener el formato AAAA-MM-DD"
    Y el formulario conserva el nombre ingresado
    Y no se crea ningún proyecto

  @CA2 @Fase1
  Escenario: No se crea un ítem sin título
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Cuando se completa el alta de historia con título vacío y prioridad "Alta"
    Y se envía el formulario
    Entonces se muestra el mensaje "el título del ítem no puede estar vacío"
    Y el formulario conserva la prioridad "Alta"
    Y el backlog del proyecto "PRJ-001" sigue vacío

  @CA2 @CA3 @Fase1
  Escenario: No se acepta una prioridad fuera de Alta, Media o Baja
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Cuando se envía el alta de historia con título "Registrar horas" y prioridad "Urgente"
    Entonces se muestra el mensaje "la prioridad del ítem no es válida: debe ser Alta, Media o Baja"
    Y el formulario conserva el título "Registrar horas"
    Y el backlog del proyecto "PRJ-001" sigue vacío

  @CA2 @Fase1
  Escenario: No se estima un ítem con story points inválidos
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión" y 0 Story Points
    Cuando se intenta estimar la historia "HU-01" con 4 Story Points
    Entonces se muestra el mensaje "los story points no son válidos: deben ser 1, 2, 3, 5, 8, 13 o 21"
    Y la historia "HU-01" sigue con 0 Story Points

  @CA2 @Fase1
  Escenario: Editar un ítem inexistente muestra un mensaje
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" no tiene la historia "HU-99"
    Cuando se intenta editar la historia "HU-99"
    Entonces se muestra el mensaje "el ítem no existe en el proyecto"
    Y se informa que la historia no existe en el proyecto

  @CA2 @Fase1
  Escenario: Un proyecto inexistente muestra un mensaje
    Dado que no existe un proyecto identificado por "PRJ-999"
    Cuando se visita el detalle del proyecto "PRJ-999"
    Entonces se muestra el mensaje "no se encontró el proyecto"
    Y se informa que el proyecto no existe

  @CA1 @Fase2
  Escenario: Listar los sprints con sus Story Points planificados
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint número 1 con objetivo "Primera entrega"
    Y que el proyecto "PRJ-001" tiene el Sprint número 2 con objetivo "Segunda entrega"
    Y que el Sprint número 1 del proyecto "PRJ-001" tiene historias por 8 Story Points planificados
    Y que el Sprint número 2 del proyecto "PRJ-001" no tiene historias
    Cuando se visita la sección de sprints del proyecto "PRJ-001"
    Entonces se listan los Sprints número 1 y número 2 con sus objetivos, fechas y estado "Activo"
    Y el Sprint número 1 muestra 8 Story Points planificados
    Y el Sprint número 2 muestra 0 Story Points planificados

  @CA1 @Fase2
  Escenario: Crear un sprint
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" no tiene Sprints
    Cuando se completa el alta de sprint con número 1, objetivo "Primera entrega", fecha de inicio "2026-03-01" y fecha de fin "2026-03-14"
    Y se envía el formulario
    Entonces la sección de sprints del proyecto "PRJ-001" muestra el Sprint número 1
    Y el Sprint número 1 muestra el objetivo "Primera entrega" y estado "Activo"

  @CA1 @Fase2
  Escenario: Asignar una historia estimada a un sprint
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint número 1 con objetivo "Primera entrega"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Cuando se asigna la historia "HU-01" al Sprint número 1 del proyecto "PRJ-001"
    Entonces la asignación se realiza exitosamente
    Y la sección de sprints muestra que el Sprint número 1 tiene la historia "HU-01"
    Y el Sprint número 1 muestra 3 Story Points planificados

  @CA2 @Fase2
  Escenario: No se asigna una historia sin estimar
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint número 1 con objetivo "Primera entrega"
    Y que el proyecto "PRJ-001" tiene la historia "HU-04" con título "Ajustar la configuración", estado "Pendiente" y 0 Story Points
    Cuando se intenta asignar la historia "HU-04" al Sprint número 1 del proyecto "PRJ-001"
    Entonces se muestra el mensaje "el ítem no está estimado"
    Y el Sprint número 1 no queda con la historia "HU-04"
    Y el Sprint número 1 sigue con 0 Story Points planificados

  @CA2 @Fase2
  Escenario: Asignar a un sprint inexistente muestra un mensaje
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Cuando se intenta asignar la historia "HU-01" al Sprint número 9 del proyecto "PRJ-001"
    Entonces se muestra el mensaje "el sprint no existe en el proyecto"
    Y la historia "HU-01" no queda asignada a ningún Sprint

  @CA1 @Fase2
  Escenario: Un sprint cerrado muestra sus Story Points planificados
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint número 1 cerrado con objetivo "Primera entrega" y 8 Story Points planificados guardados
    Cuando se visita la sección de sprints del proyecto "PRJ-001"
    Entonces el Sprint número 1 se lista con estado "Cerrado"
    Y el Sprint número 1 muestra 8 Story Points planificados

  @CA2 @Fase2
  Escenario: No se asigna una historia a un sprint cerrado
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint número 1 cerrado con objetivo "Primera entrega"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Cuando se intenta asignar la historia "HU-01" al Sprint número 1 del proyecto "PRJ-001"
    Entonces se muestra el mensaje "el sprint está cerrado"
    Y la historia "HU-01" no queda asignada al Sprint número 1

  @CA3 @Fase2
  Escenario: El formulario de asignación ofrece las historias del proyecto
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene el Sprint número 1 con objetivo "Primera entrega"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión"
    Y que el proyecto "PRJ-001" tiene la historia "HU-02" con título "Cerrar sesión"
    Y un proyecto existente identificado por "PRJ-002", con nombre "Portal de Clientes"
    Y que el proyecto "PRJ-002" tiene la historia "HU-01" con título "Iniciar sesión del portal"
    Cuando se abre el formulario de asignación del Sprint número 1 del proyecto "PRJ-001"
    Entonces el formulario ofrece las historias "HU-01" y "HU-02" del proyecto "PRJ-001"
    Y no ofrece la historia "HU-01" del proyecto "PRJ-002"
