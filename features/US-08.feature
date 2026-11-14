# language: es
# Estado: Propuesta
# Los escenarios describen el comportamiento acordado para CA1–CA3 y las
# decisiones D1–D10: una historia del backlog se asigna a un Sprint existente
# del mismo proyecto; solo se asignan historias estimadas y pendientes; una
# historia pertenece como máximo a un Sprint del proyecto; los Story Points
# planificados del Sprint son la suma de los story points de sus historias.
# El orden de evaluación es: proyecto, Sprint, ítem, ítem ya asignado, ítem
# estimado e ítem pendiente.
# Los errores se describen por su significado; sus identificadores exactos están
# decididos en la spec (sección 8).
# El escenario de una historia no pendiente usa el estado "Completada", que se
# incorporará con US-09; se verifica con el doble en memoria porque el esquema
# actual de items_backlog solo admite 'Pendiente'.
# No se escenifican la traducción de errores de repositorio (R10), la falta de
# validación del estado del Sprint (R12, cubierta por el escenario de asignación
# válida) ni el formato inverso del ID (D10): se verifican con tests sin
# escenario BDD propio.

@US-08
Característica: Asignar historias a un Sprint
  Como equipo
  Quiero asignar historias del backlog a un Sprint
  Para definir el Sprint Backlog

  @CA1
  Escenario: Asignar una historia estimada y pendiente a un Sprint existente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Cuando se asigna la historia "HU-01" al Sprint número 1 del proyecto "p1"
    Entonces la asignación se realiza exitosamente
    Y el Sprint número 1 del proyecto "p1" queda con la historia "HU-01"

  @CA1
  Escenario: No se asigna una historia sin estimar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-04" con título "Ajustar la configuración", estado "Pendiente" y 0 Story Points
    Cuando se intenta asignar la historia "HU-04" al Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no está estimada
    Y la historia "HU-04" no queda asignada a ningún Sprint
    Y el Sprint número 1 del proyecto "p1" no queda con la historia "HU-04"

  @CA1
  Escenario: No se asigna una historia que no está pendiente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-05" con título "Optimizar consultas", estado "Completada" y 5 Story Points
    Cuando se intenta asignar la historia "HU-05" al Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no está pendiente
    Y la historia "HU-05" no queda asignada a ningún Sprint
    Y el Sprint número 1 del proyecto "p1" no queda con la historia "HU-05"

  @CA2
  Escenario: No se asigna dos veces la misma historia al mismo Sprint
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-06" con título "Registrar horas", estado "Pendiente" y 3 Story Points
    Y que la historia "HU-06" ya está asignada al Sprint número 1 del proyecto "p1"
    Cuando se intenta asignar la historia "HU-06" al Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia ya está asignada a un Sprint
    Y el Sprint número 1 del proyecto "p1" sigue teniendo una sola vez la historia "HU-06"

  @CA2
  Escenario: No se asigna una historia que ya está en otro Sprint
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene el Sprint número 2 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-06" con título "Registrar horas", estado "Pendiente" y 3 Story Points
    Y que la historia "HU-06" ya está asignada al Sprint número 1 del proyecto "p1"
    Cuando se intenta asignar la historia "HU-06" al Sprint número 2 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia ya está asignada a un Sprint
    Y el Sprint número 2 del proyecto "p1" no queda con la historia "HU-06"
    Y el Sprint número 1 del proyecto "p1" conserva la historia "HU-06"

  @CA2
  Escenario: La misma historia de otro proyecto no interfiere
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Y un proyecto existente identificado por "p2", con nombre "Portal de Clientes"
    Y que el proyecto "p2" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p2" tiene la historia "HU-01" con título "Iniciar sesión del portal", estado "Pendiente" y 8 Story Points
    Y que la historia "HU-01" del proyecto "p2" ya está asignada al Sprint número 1 del proyecto "p2"
    Cuando se asigna la historia "HU-01" del proyecto "p1" al Sprint número 1 del proyecto "p1"
    Entonces la asignación se realiza exitosamente
    Y el Sprint número 1 del proyecto "p1" queda con la historia "HU-01" de 3 Story Points
    Y el Sprint número 1 del proyecto "p2" conserva la historia "HU-01" de 8 Story Points

  @CA1
  Escenario: Proyecto inexistente devuelve error
    Dado que no existe un proyecto identificado por "p999"
    Cuando se intenta asignar la historia "HU-01" al Sprint número 1 del proyecto "p999"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y la historia "HU-01" no queda asignada a ningún Sprint

  @CA1
  Escenario: Sprint inexistente devuelve error
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Cuando se intenta asignar la historia "HU-01" al Sprint número 3 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint no existe en el proyecto
    Y la historia "HU-01" no queda asignada a ningún Sprint

  @CA1
  Escenario: Historia inexistente devuelve error
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Cuando se intenta asignar la historia "HU-99" al Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no existe en el proyecto
    Y el Sprint número 1 del proyecto "p1" no queda con la historia "HU-99"

  @CA1
  Escenario: ID de historia malformado se trata como historia inexistente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Cuando se intenta asignar la historia "HU-X" al Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no existe en el proyecto
    Y el Sprint número 1 del proyecto "p1" no queda con ninguna historia

  @CA3
  Escenario: Calcular los Story Points planificados como suma de las historias del Sprint
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Y que el proyecto "p1" tiene la historia "HU-02" con título "Cerrar sesión", estado "Pendiente" y 5 Story Points
    Y que las historias "HU-01" y "HU-02" están asignadas al Sprint número 1 del proyecto "p1"
    Cuando se consultan los Story Points planificados del Sprint número 1 del proyecto "p1"
    Entonces los Story Points planificados del Sprint número 1 son 8

  @CA3
  Escenario: Un Sprint sin historias tiene 0 Story Points planificados
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 2 en estado "Activo"
    Y que el Sprint número 2 del proyecto "p1" no tiene historias asignadas
    Cuando se consultan los Story Points planificados del Sprint número 2 del proyecto "p1"
    Entonces los Story Points planificados del Sprint número 2 son 0

  @CA3
  Escenario: Actualizar los Story Points planificados al asignar una nueva historia
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Y que la historia "HU-01" está asignada al Sprint número 1 del proyecto "p1"
    Y que el proyecto "p1" tiene la historia "HU-03" con título "Recuperar contraseña", estado "Pendiente" y 2 Story Points
    Cuando se asigna la historia "HU-03" al Sprint número 1 del proyecto "p1"
    Entonces la asignación se realiza exitosamente
    Y los Story Points planificados del Sprint número 1 son 5

  @CA1
  Escenario: Proyecto inexistente con Sprint inexistente devuelve primero el error de proyecto
    Dado que no existe un proyecto identificado por "p999"
    Cuando se intenta asignar la historia "HU-01" al Sprint número 9 del proyecto "p999"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de Sprint inexistente
    Y la historia "HU-01" no queda asignada a ningún Sprint

  @CA1
  Escenario: Sprint inexistente con historia inexistente devuelve primero el error de Sprint
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Cuando se intenta asignar la historia "HU-99" al Sprint número 9 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint no existe en el proyecto
    Y no se devuelve un error de historia inexistente
    Y la historia "HU-99" no queda asignada a ningún Sprint
