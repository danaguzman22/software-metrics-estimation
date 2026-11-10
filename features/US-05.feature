# language: es
# US-05 pendiente de implementación: solo spec (specs/US-05.md) y estos escenarios.
# Los escenarios describen el comportamiento acordado para CA1–CA4 y las decisiones
# confirmadas: ítem en un proyecto existente identificado por ID; ID automático "HU-01",
# "HU-02", ... único dentro del proyecto y siguiente al mayor existente; título no vacío
# ni solo espacios, guardado recortado; prioridad solo "Alta", "Media" o "Baja" exactas
# tras recortar espacios; estado inicial siempre "Pendiente"; story points 0; descripción
# y criterios opcionales y de texto libre sin validar; orden de errores: proyecto, título,
# prioridad.
# Los errores se describen por su significado; sus identificadores exactos están decididos
# en la spec (D10).
# No se escenifican las decisiones "no se valida" (unicidad de título, límites de longitud,
# cantidad de ítems) ni nada fuera de alcance (editar o reordenar, asignar a Sprint,
# estimar, interfaz web).

@US-05
Característica: Crear ítem del Product Backlog
  Como Product Owner
  Quiero cargar historias en el backlog
  Para organizar el trabajo

  @CA1
  Escenario: Alta exitosa de un ítem del backlog
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se crea en ese proyecto un ítem con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta" y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado en el proyecto "Sistema de Métricas" con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 0 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"

  @CA1
  Escenario: Se puede crear un segundo ítem en el mismo proyecto
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y prioridad "Alta"
    Cuando se crea en ese proyecto un ítem con título "Modificar proyecto", descripción "Permitir cambiar los datos de un proyecto", prioridad "Media" y sin criterios de aceptación
    Entonces el alta se realiza exitosamente
    Y el proyecto "Sistema de Métricas" queda con los ítems "HU-01" con título "Registrar proyecto" y "HU-02" con título "Modificar proyecto"

  @CA1
  Escenario: Primer alta sin archivo de backlog
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene ítems registrados porque todavía no existe el archivo de backlog
    Cuando se crea en ese proyecto un ítem con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta" y sin criterios de aceptación
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado en el proyecto "Sistema de Métricas" con ID "HU-01" y título "Registrar proyecto"

  @CA2
  Escenario: El primer ítem recibe el ID automático HU-01
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene ítems en el backlog
    Cuando se crea en ese proyecto un ítem con título "Registrar proyecto" y prioridad "Alta"
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado con ID "HU-01"

  @CA2
  Escenario: El ID automático es consecutivo al mayor existente
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene ítems con IDs "HU-01" y "HU-03"
    Cuando se crea en ese proyecto un ítem con título "Consultar métricas" y prioridad "Baja"
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado con ID "HU-04"

  @CA2
  Escenario: El ID es único dentro del proyecto
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01"
    Cuando se crea en ese proyecto un ítem con título "Modificar proyecto" y prioridad "Media"
    Entonces el alta se realiza exitosamente
    Y el proyecto "Sistema de Métricas" queda con dos ítems de IDs distintos: "HU-01" y "HU-02"

  @CA2
  Escenario: El mismo ID puede repetirse en proyectos distintos
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01"
    Y un proyecto existente identificado por su ID, con nombre "App Móvil", fecha de inicio "1 de abril de 2026", fecha de fin "31 de diciembre de 2026" y estado "Planificado"
    Y que el proyecto "App Móvil" no tiene ítems en el backlog
    Cuando se crea en el proyecto "App Móvil" un ítem con título "Registrar proyecto" y prioridad "Alta"
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado en el proyecto "App Móvil" con ID "HU-01"
    Y el proyecto "Sistema de Métricas" sigue teniendo su ítem con ID "HU-01"

  @CA4
  Escenario: El estado inicial de un ítem es Pendiente
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se crea en ese proyecto un ítem con título "Registrar proyecto" y prioridad "Alta"
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado con estado "Pendiente"

  @CA1
  Escenario: Los story points de un ítem nuevo quedan en 0
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se crea en ese proyecto un ítem con título "Registrar proyecto" y prioridad "Alta"
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado con story points 0
    Y no se solicita una estimación al crear el ítem

  @CA1
  Escenario: Descripción y criterios de aceptación son opcionales
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se crea en ese proyecto un ítem con título "Revisar métricas", descripción "" y sin criterios de aceptación, con prioridad "Baja"
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado con descripción vacía, sin criterios de aceptación, estado "Pendiente" y story points 0

  @CA1
  Escenario: Descripción y criterios de aceptación admiten texto libre
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se crea en ese proyecto un ítem con título "Analizar métricas", descripción "¿Qué métricas se van a calcular? Aún no está definido.", prioridad "Media" y criterios de aceptación "Se listan las métricas" y "Cada métrica tiene una descripción"
    Entonces el alta se realiza exitosamente
    Y la descripción y los criterios de aceptación quedan registrados tal como se recibieron

  @CA1 @CA3
  Escenario: Título y prioridad con espacios externos se recortan
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se crea en ese proyecto un ítem con título "  Registrar proyecto  ", descripción "Permitir dar de alta proyectos en el sistema", prioridad "  Alta  " y sin criterios de aceptación
    Entonces el alta se realiza exitosamente
    Y el ítem queda registrado con título "Registrar proyecto" y prioridad "Alta", sin espacios al inicio y al final

  @CA3
  Escenario: Título vacío devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta crear en ese proyecto un ítem con título "", descripción "Permitir dar de alta proyectos en el sistema" y prioridad "Alta"
    Entonces se devuelve un error indicando que el título del ítem no puede estar vacío
    Y el alta no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con ítems registrados

  @CA3
  Escenario: Título compuesto únicamente por espacios devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta crear en ese proyecto un ítem con título "   ", descripción "Permitir dar de alta proyectos en el sistema" y prioridad "Alta"
    Entonces se devuelve un error indicando que el título del ítem no puede estar vacío
    Y el alta no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con ítems registrados

  @CA3
  Escenario: Prioridad inválida devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta crear en ese proyecto un ítem con título "Registrar proyecto" y prioridad "Urgente"
    Entonces se devuelve un error indicando que la prioridad del ítem no es válida y debe ser "Alta", "Media" o "Baja"
    Y el alta no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con ítems registrados

  @CA3
  Escenario: Prioridad con mayúsculas distintas devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta crear en ese proyecto un ítem con título "Registrar proyecto" y prioridad "alta"
    Entonces se devuelve un error indicando que la prioridad del ítem no es válida y debe ser "Alta", "Media" o "Baja"
    Y el alta no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con ítems registrados

  @CA3
  Escenario: Prioridad compuesta únicamente por espacios devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta crear en ese proyecto un ítem con título "Registrar proyecto" y prioridad "   "
    Entonces se devuelve un error indicando que la prioridad del ítem no es válida y debe ser "Alta", "Media" o "Baja"
    Y el alta no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con ítems registrados

  @CA1
  Escenario: Proyecto inexistente devuelve error
    Dado que no existe un proyecto identificado por el ID del proyecto al que se quiere crear el ítem
    Cuando se intenta crear en ese proyecto un ítem con título "Registrar proyecto" y prioridad "Alta"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y el alta no se completa exitosamente

  @CA1 @CA3
  Escenario: Proyecto inexistente con título vacío devuelve error de proyecto inexistente
    Dado que no existe un proyecto identificado por el ID del proyecto al que se quiere crear el ítem
    Cuando se intenta crear en ese proyecto un ítem con título "" y prioridad "Alta"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y el alta no se completa exitosamente

  @CA3
  Escenario: Título vacío con prioridad inválida devuelve error de título vacío
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Cuando se intenta crear en ese proyecto un ítem con título "" y prioridad "Urgente"
    Entonces se devuelve un error indicando que el título del ítem no puede estar vacío
    Y el alta no se completa exitosamente
    Y el proyecto "Sistema de Métricas" no queda con ítems registrados
