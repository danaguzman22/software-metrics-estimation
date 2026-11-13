# language: es
# Estado: Propuesta
# Los escenarios describen el comportamiento acordado para CA1–CA3 y las decisiones
# D1–D3: se editan título, descripción, prioridad, criterios de aceptación y story
# points; en story points se admite 0 ("sin estimar") o 1, 2, 3, 5, 8, 13, 21; el
# ID y el estado no se editan; el backlog se lista por prioridad Alta, Media, Baja
# y, a igual prioridad, por número ascendente usando comparación numérica; un ítem
# completado no se puede editar.
# Los errores se describen por su significado; sus identificadores exactos están
# decididos en la spec (sección 8).
# No se escenifica la persistencia del estado "Completada" porque queda diferida
# (D3), ni nada fuera de alcance.

@US-06
Característica: Editar y priorizar el backlog
  Como Product Owner
  Quiero editar y reordenar historias por prioridad
  Para mantener el backlog actualizado

  @CA1
  Escenario: Editar el título de un ítem
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"
    Cuando se edita el ítem "HU-01" con título "Registrar y listar proyectos", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con título "Registrar y listar proyectos"
    Y el ítem "HU-01" conserva su descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"

  @CA1
  Escenario: Editar la descripción de un ítem
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se edita el ítem "HU-01" con título "Registrar proyecto", descripción "Permitir dar de alta y consultar proyectos en el sistema", prioridad "Alta", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con descripción "Permitir dar de alta y consultar proyectos en el sistema"
    Y el ítem "HU-01" conserva su título "Registrar proyecto", prioridad "Alta", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA1
  Escenario: Editar la prioridad de un ítem
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se edita el ítem "HU-01" con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Media", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con prioridad "Media"
    Y el ítem "HU-01" conserva su título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA1
  Escenario: Editar los criterios de aceptación de un ítem
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"
    Cuando se edita el ítem "HU-01" con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío", "El proyecto queda disponible para asignar integrantes" y "El proyecto se lista en el backlog"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con criterios de aceptación "El nombre del proyecto no puede estar vacío", "El proyecto queda disponible para asignar integrantes" y "El proyecto se lista en el backlog"
    Y el ítem "HU-01" conserva su título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta" y story points 5

  @CA1
  Escenario: Editar los story points de un ítem
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se edita el ítem "HU-01" con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 8 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con story points 8
    Y el ítem "HU-01" conserva su título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta" y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA1
  Escenario: Editar el título de un ítem sin estimar conserva los story points en 0
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 0 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se edita el ítem "HU-01" con título "Registrar y listar proyectos", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 0 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con título "Registrar y listar proyectos"
    Y el ítem "HU-01" conserva sus story points 0

  @CA1
  Escenario: Dejar un ítem sin estimar con 0 story points
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se edita el ítem "HU-01" con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 0 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con story points 0
    Y el ítem "HU-01" conserva su título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta" y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA1
  Escenario: Editar todos los campos editables a la vez
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se edita el ítem "HU-01" con título "Registrar y listar proyectos", descripción "Permitir dar de alta y consultar proyectos en el sistema", prioridad "Media", story points 8 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto se lista en el backlog"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" queda con título "Registrar y listar proyectos", descripción "Permitir dar de alta y consultar proyectos en el sistema", prioridad "Media", story points 8 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto se lista en el backlog"

  @CA1
  Escenario: El ID y el estado no cambian al editar
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se edita el ítem "HU-01" con título "Registrar y listar proyectos", descripción "Permitir dar de alta y consultar proyectos en el sistema", prioridad "Media", story points 8 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto se lista en el backlog"
    Entonces la edición se realiza exitosamente
    Y el ítem "HU-01" conserva su ID "HU-01" y su estado "Pendiente"
    Y el proyecto sigue teniendo un solo ítem con ID "HU-01"

  @CA1
  Escenario: Título vacío devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se intenta editar el ítem "HU-01" con título "", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces se devuelve un error indicando que el título del ítem no puede estar vacío
    Y la edición no se completa exitosamente
    Y el ítem "HU-01" conserva su título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA1
  Escenario: Título compuesto solo por espacios devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se intenta editar el ítem "HU-01" con título "   ", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces se devuelve un error indicando que el título del ítem no puede estar vacío
    Y la edición no se completa exitosamente
    Y el ítem "HU-01" conserva su título "Registrar proyecto"

  @CA1
  Escenario: Prioridad inválida devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se intenta editar el ítem "HU-01" con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Urgente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces se devuelve un error indicando que la prioridad debe ser Alta, Media o Baja
    Y la edición no se completa exitosamente
    Y el ítem "HU-01" conserva su prioridad "Alta", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA1
  Escenario: Story points inválidos devuelven error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se intenta editar el ítem "HU-01" con título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", story points 4 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces se devuelve un error indicando que los story points no son válidos; se admite 0 (sin estimar) o 1, 2, 3, 5, 8, 13 o 21
    Y la edición no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 5, título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente" y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA1
  Escenario: Proyecto inexistente al editar devuelve error
    Dado que no existe un proyecto identificado por el ID del proyecto al que pertenece el ítem que se quiere editar
    Cuando se intenta editar el ítem "HU-01" de ese proyecto con título "Registrar y listar proyectos", descripción "Permitir dar de alta y consultar proyectos en el sistema", prioridad "Media", story points 8 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y la edición no se completa exitosamente

  @CA1
  Escenario: Ítem inexistente al editar devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se intenta editar el ítem "HU-02" del proyecto con título "Modificar proyecto", descripción "Permitir modificar proyectos existentes", prioridad "Media", story points 3 y criterios de aceptación "El proyecto modificado conserva su ID"
    Entonces se devuelve un error indicando que el ítem no existe en el proyecto
    Y la edición no se completa exitosamente
    Y el ítem "HU-01" conserva su título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA3
  Escenario: No se puede editar un ítem completado
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Completada", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se intenta editar el ítem "HU-01" con título "Registrar y listar proyectos", descripción "Permitir dar de alta y consultar proyectos en el sistema", prioridad "Media", story points 8 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces se devuelve un error indicando que el ítem está completado y no se puede editar
    Y la edición no se completa exitosamente
    Y el ítem "HU-01" conserva su título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Completada", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"

  @CA3
  Escenario: Un ítem completado devuelve error antes que las validaciones de campos
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Completada", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se intenta editar el ítem "HU-01" con título "", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Urgente", story points 4 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Entonces se devuelve un error indicando que el ítem está completado y no se puede editar
    Y no se devuelve un error de título vacío ni de prioridad inválida ni de story points inválidos
    Y la edición no se completa exitosamente
    Y el ítem "HU-01" conserva su estado "Completada" y su título "Registrar proyecto"

  @CA2
  Escenario: El backlog se lista ordenado por prioridad Alta, Media, Baja
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-03", título "Consultar métricas", descripción "Consultar las métricas del proyecto", prioridad "Baja", estado "Pendiente", story points 5 y criterios de aceptación "Se listan las métricas"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Y que el proyecto tiene un ítem con ID "HU-02", título "Modificar proyecto", descripción "Permitir modificar proyectos existentes", prioridad "Media", estado "Pendiente", story points 3 y criterios de aceptación "El proyecto modificado conserva su ID"
    Cuando se lista el backlog del proyecto
    Entonces el backlog se devuelve ordenado por prioridad: "HU-01" con prioridad "Alta", "HU-02" con prioridad "Media" y "HU-03" con prioridad "Baja"

  @CA2
  Escenario: A igual prioridad, el desempate es por número ascendente
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-03", título "Consultar métricas", descripción "Consultar las métricas del proyecto", prioridad "Media", estado "Pendiente", story points 5 y criterios de aceptación "Se listan las métricas"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Media", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Y que el proyecto tiene un ítem con ID "HU-02", título "Modificar proyecto", descripción "Permitir modificar proyectos existentes", prioridad "Media", estado "Pendiente", story points 3 y criterios de aceptación "El proyecto modificado conserva su ID"
    Cuando se lista el backlog del proyecto
    Entonces el backlog se devuelve ordenado por número ascendente: "HU-01", "HU-02" y "HU-03", todos con prioridad "Media"

  @CA2
  Escenario: El desempate numérico ubica HU-2 antes que HU-10
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-10", título "Consultar métricas", descripción "Consultar las métricas del proyecto", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "Se listan las métricas"
    Y que el proyecto tiene un ítem con ID "HU-2", título "Modificar proyecto", descripción "Permitir modificar proyectos existentes", prioridad "Alta", estado "Pendiente", story points 3 y criterios de aceptación "El proyecto modificado conserva su ID"
    Y que el proyecto tiene un ítem con ID "HU-1", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 5 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se lista el backlog del proyecto
    Entonces el backlog se devuelve ordenado numéricamente: "HU-1", "HU-2" y "HU-10", todos con prioridad "Alta"
    Y "HU-2" aparece antes que "HU-10"

  @CA2
  Escenario: Backlog vacío devuelve una lista vacía
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene ítems en el backlog
    Cuando se lista el backlog del proyecto
    Entonces el backlog se devuelve vacío
    Y no se devuelve ningún error

  @CA2
  Escenario: Proyecto inexistente al listar devuelve error
    Dado que no existe un proyecto identificado por el ID del proyecto cuyo backlog se quiere listar
    Cuando se lista el backlog de ese proyecto
    Entonces se devuelve un error indicando que el proyecto no existe
    Y el listado no se completa exitosamente
