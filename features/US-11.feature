# language: es
# Estado: Implementada
# Los escenarios describen el comportamiento acordado para CA1–CA2 y las decisiones
# confirmadas: la estimación se hace sobre un ítem existente de un proyecto existente,
# identificado por proyectoID e itemID exacto; solo se aceptan los valores de la escala
# Fibonacci 1, 2, 3, 5, 8, 13 y 21; 0 significa "sin estimar" y no es válido; se puede
# re-estimar y el nuevo valor reemplaza al anterior; orden de errores: proyecto, ítem,
# Story Points.
# Los errores se describen por su significado; sus identificadores exactos están decididos
# en la spec (D5, D6).
# No se escenifican las decisiones "no se valida" (estado del ítem, cantidad de ítems) ni
# nada fuera de alcance (Planning Poker, asignar a Sprint, editar otros campos, listar,
# interfaz web).

@US-11
Característica: Estimar historias con Story Points
  Como equipo
  Quiero asignar Story Points a una historia
  Para dimensionar el trabajo

  @CA1
  Esquema del escenario: Estimación válida con cada valor de la escala
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 0 y criterios de aceptación "El nombre del proyecto no puede estar vacío"
    Cuando se estima el ítem "HU-01" del proyecto con <story points> story points
    Entonces la estimación se realiza exitosamente
    Y el ítem "HU-01" queda con <story points> story points

    Ejemplos:
      | story points |
      | 1            |
      | 2            |
      | 3            |
      | 5            |
      | 8            |
      | 13           |
      | 21           |

  @CA1
  Escenario: La estimación conserva el resto de los campos del ítem
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente", story points 0 y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"
    Cuando se estima el ítem "HU-01" del proyecto con 5 story points
    Entonces la estimación se realiza exitosamente
    Y el ítem "HU-01" queda con story points 5
    Y el ítem "HU-01" conserva su ID "HU-01", título "Registrar proyecto", descripción "Permitir dar de alta proyectos en el sistema", prioridad "Alta", estado "Pendiente" y criterios de aceptación "El nombre del proyecto no puede estar vacío" y "El proyecto queda disponible para asignar integrantes"

  @CA1
  Escenario: La re-estimación reemplaza el valor anterior
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 3
    Cuando se estima el ítem "HU-01" del proyecto con 8 story points
    Entonces la estimación se realiza exitosamente
    Y el ítem "HU-01" queda con story points 8
    Y el proyecto sigue teniendo un solo ítem con ID "HU-01"

  @CA1
  Escenario: Un ítem sin estimar pasa a estar estimado
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-02", título "Modificar proyecto" y story points 0
    Cuando se estima el ítem "HU-02" del proyecto con 5 story points
    Entonces la estimación se realiza exitosamente
    Y el ítem "HU-02" queda con story points 5

  @CA1
  Escenario: La estimación conserva los demás ítems y su orden
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 2
    Y que el proyecto tiene un ítem con ID "HU-02", título "Modificar proyecto" y story points 3
    Y que el proyecto tiene un ítem con ID "HU-03", título "Consultar métricas" y story points 5
    Cuando se estima el ítem "HU-02" del proyecto con 13 story points
    Entonces la estimación se realiza exitosamente
    Y el ítem "HU-02" queda con story points 13
    Y el proyecto conserva el orden de ítems "HU-01" con story points 2, "HU-02" con story points 13 y "HU-03" con story points 5

  @CA1
  Escenario: Un ítem con el mismo ID en otro proyecto no se modifica
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 0
    Y un proyecto existente identificado por su ID, con nombre "App Móvil", fecha de inicio "1 de abril de 2026", fecha de fin "31 de diciembre de 2026" y estado "Planificado"
    Y que el proyecto "App Móvil" tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 2
    Cuando se estima el ítem "HU-01" del proyecto "Sistema de Métricas" con 8 story points
    Entonces la estimación se realiza exitosamente
    Y el ítem "HU-01" del proyecto "Sistema de Métricas" queda con story points 8
    Y el ítem "HU-01" del proyecto "App Móvil" conserva sus story points 2

  @CA2
  Escenario: Story Points 0 devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 0
    Cuando se intenta estimar el ítem "HU-01" del proyecto con 0 story points
    Entonces se devuelve un error indicando que los story points no son válidos y deben ser 1, 2, 3, 5, 8, 13 o 21
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 0

  @CA2
  Escenario: Story Points negativos devuelven error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 0
    Cuando se intenta estimar el ítem "HU-01" del proyecto con -1 story points
    Entonces se devuelve un error indicando que los story points no son válidos y deben ser 1, 2, 3, 5, 8, 13 o 21
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 0

  @CA2
  Escenario: Story Points 4 devuelven error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 0
    Cuando se intenta estimar el ítem "HU-01" del proyecto con 4 story points
    Entonces se devuelve un error indicando que los story points no son válidos y deben ser 1, 2, 3, 5, 8, 13 o 21
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 0

  @CA2
  Escenario: Story Points 22 devuelven error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 0
    Cuando se intenta estimar el ítem "HU-01" del proyecto con 22 story points
    Entonces se devuelve un error indicando que los story points no son válidos y deben ser 1, 2, 3, 5, 8, 13 o 21
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 0

  @CA2
  Escenario: Valor inválido sobre un ítem ya estimado conserva el valor anterior
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 3
    Cuando se intenta estimar el ítem "HU-01" del proyecto con 4 story points
    Entonces se devuelve un error indicando que los story points no son válidos y deben ser 1, 2, 3, 5, 8, 13 o 21
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 3

  @CA2
  Escenario: Proyecto inexistente devuelve error
    Dado que no existe un proyecto identificado por el ID del proyecto al que pertenece el ítem que se quiere estimar
    Cuando se intenta estimar el ítem "HU-01" de ese proyecto con 5 story points
    Entonces se devuelve un error indicando que el proyecto no existe
    Y la estimación no se completa exitosamente

  @CA2
  Escenario: Ítem inexistente devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 2
    Cuando se intenta estimar el ítem "HU-02" del proyecto con 5 story points
    Entonces se devuelve un error indicando que el ítem no existe en el proyecto
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 2

  @CA2
  Escenario: Backlog sin archivo devuelve ítem inexistente
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto no tiene ítems registrados porque todavía no existe el archivo de backlog
    Cuando se intenta estimar el ítem "HU-01" del proyecto con 5 story points
    Entonces se devuelve un error indicando que el ítem no existe en el proyecto
    Y la estimación no se completa exitosamente

  @CA2
  Escenario: Un ítem con el mismo ID en otro proyecto no se encuentra
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto "Sistema de Métricas" no tiene ítems en el backlog
    Y un proyecto existente identificado por su ID, con nombre "App Móvil", fecha de inicio "1 de abril de 2026", fecha de fin "31 de diciembre de 2026" y estado "Planificado"
    Y que el proyecto "App Móvil" tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 2
    Cuando se intenta estimar el ítem "HU-01" del proyecto "Sistema de Métricas" con 8 story points
    Entonces se devuelve un error indicando que el ítem no existe en el proyecto
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" del proyecto "App Móvil" conserva sus story points 2

  @CA2
  Escenario: Proyecto inexistente con Story Points inválidos devuelve error de proyecto
    Dado que no existe un proyecto identificado por el ID del proyecto al que pertenece el ítem que se quiere estimar
    Cuando se intenta estimar el ítem "HU-01" de ese proyecto con 4 story points
    Entonces se devuelve un error indicando que el proyecto no existe
    Y la estimación no se completa exitosamente

  @CA2
  Escenario: Ítem inexistente con Story Points inválidos devuelve error de ítem
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado "En curso"
    Y que el proyecto tiene un ítem con ID "HU-01", título "Registrar proyecto" y story points 2
    Cuando se intenta estimar el ítem "HU-02" del proyecto con 4 story points
    Entonces se devuelve un error indicando que el ítem no existe en el proyecto
    Y la estimación no se completa exitosamente
    Y el ítem "HU-01" conserva sus story points 2
