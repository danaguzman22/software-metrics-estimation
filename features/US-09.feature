# language: es
# Estado: Propuesta
# Los escenarios describen US-09 y las decisiones D1–D18 aprobadas para la
# historia. D19–D21 quedan pendientes de confirmación de Dana. El orden de
# evaluación al marcar completada es: proyecto, Sprint, historia; para cerrar y
# consultar métricas, proyecto, Sprint. El cierre conserva las asignaciones
# completadas y elimina las no completadas en la misma transacción; las historias
# devueltas quedan "Pendiente" en el backlog.
# La implementación usa provisionalmente el campo SPPlanificados, el método
# MarcarComoCompletada y el comportamiento del Down propuestos en specs/US-09.md.
# La suite completa debe seguir en verde para US-07. Un test de integración
# verifica que cerrar un Sprint libera el índice parcial de Sprint activo.
# La suite repotest cubre marcar una historia como completada en jsonrepo y
# postgres; son verificaciones de contrato, no escenarios BDD.

@US-09
Característica: Marcar historias como completadas y cerrar el Sprint
  Como equipo
  Quiero marcar historias como completadas y cerrar el Sprint
  Para medir lo entregado

  @D4
  Escenario: Marcar como completada una historia pendiente asignada al Sprint activo
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Y que la historia "HU-01" está asignada al Sprint número 1 del proyecto "p1"
    Cuando se marca como completada la historia "HU-01" del Sprint número 1 del proyecto "p1"
    Entonces la historia "HU-01" queda en estado "Completada"

  @D4
  Escenario: No marcar como completada una historia que no está asignada al Sprint
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Cuando se intenta marcar como completada la historia "HU-01" del Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no está asignada a ese Sprint
    Y la historia "HU-01" permanece en estado "Pendiente"

  @D4
  Escenario: No marcar como completada una historia que no está pendiente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Completada" y 3 Story Points
    Y que la historia "HU-01" está asignada al Sprint número 1 del proyecto "p1"
    Cuando se intenta marcar como completada la historia "HU-01" del Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no está pendiente
    Y la historia "HU-01" permanece en estado "Completada"

  @D4
  Escenario: Proyecto inexistente prevalece sobre Sprint e historia inexistentes al marcar completada
    Dado que no existe un proyecto identificado por "p999"
    Cuando se intenta marcar como completada la historia "HU-99" del Sprint número 9 del proyecto "p999"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de Sprint inexistente
    Y no se devuelve un error de historia inexistente

  @D4
  Escenario: Sprint inexistente prevalece sobre historia inexistente al marcar completada
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene el Sprint número 9
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Cuando se intenta marcar como completada la historia "HU-99" del Sprint número 9 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint no existe en el proyecto
    Y no se devuelve un error de historia inexistente

  @R2
  Escenario: Proyecto inexistente prevalece sobre Sprint inexistente al cerrar
    Dado que no existe un proyecto identificado por "p999"
    Cuando se intenta cerrar el Sprint número 9 del proyecto "p999"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de Sprint inexistente

  @R2
  Escenario: Sprint inexistente al cerrar un Sprint de un proyecto existente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene el Sprint número 9
    Cuando se intenta cerrar el Sprint número 9 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint no existe en el proyecto

  @R2
  Esquema del escenario: Las consultas de métricas validan el proyecto antes que el Sprint
    Dado que no existe un proyecto identificado por "p999"
    Cuando se consultan los "<métrica>" del Sprint número 9 del proyecto "p999"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de Sprint inexistente

    Ejemplos:
      | métrica |
      | Story Points planificados |
      | Story Points completados |
      | velocidad |

  @R2
  Esquema del escenario: Las consultas de métricas requieren un Sprint del proyecto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene el Sprint número 9
    Cuando se consultan los "<métrica>" del Sprint número 9 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint no existe en el proyecto

    Ejemplos:
      | métrica |
      | Story Points planificados |
      | Story Points completados |
      | velocidad |

  @D4
  Escenario: No marcar como completada una historia inexistente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Cuando se intenta marcar como completada la historia "HU-99" del Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no existe en el proyecto

  @D14
  Escenario: Un Sprint cerrado no admite marcar historias como completadas
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Cerrado"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Y que la historia "HU-01" está asignada al Sprint número 1 del proyecto "p1"
    Cuando se intenta marcar como completada la historia "HU-01" del Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint está cerrado
    Y la historia "HU-01" permanece en estado "Pendiente"

  @CA1
  @CA2
  Escenario: Cerrar un Sprint devuelve al backlog las historias no completadas y calcula las métricas
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Completada" y 3 Story Points
    Y que el proyecto "p1" tiene la historia "HU-02" con título "Cerrar sesión", estado "Pendiente" y 5 Story Points
    Y que las historias "HU-01" y "HU-02" están asignadas al Sprint número 1 del proyecto "p1"
    Cuando se cierra el Sprint número 1 del proyecto "p1"
    Entonces el Sprint número 1 del proyecto "p1" queda en estado "Cerrado"
    Y la historia "HU-01" sigue asignada al Sprint número 1 del proyecto "p1"
    Y la historia "HU-02" deja de estar asignada a cualquier Sprint del proyecto "p1"
    Y la historia "HU-02" permanece en estado "Pendiente" en el backlog
    Y los Story Points planificados del Sprint número 1 son 8
    Y los Story Points completados del Sprint número 1 son 3
    Y la velocidad del Sprint número 1 es 3

  @CA2
  Escenario: Calcular los Story Points completados y la velocidad como suma de historias completadas
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que el proyecto "p1" tiene las historias "HU-01" y "HU-02" completadas y asignadas al Sprint número 1, con 3 y 5 Story Points
    Y que el proyecto "p1" tiene la historia "HU-03" pendiente y asignada al Sprint número 1, con 8 Story Points
    Cuando se consultan los Story Points completados y la velocidad del Sprint número 1 del proyecto "p1"
    Entonces los Story Points completados del Sprint número 1 son 8
    Y la velocidad del Sprint número 1 es 8

  @D12
  Escenario: Consultar los Story Points planificados guardados al cerrar el Sprint
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que la historia "HU-01" de 3 Story Points está asignada al Sprint número 1 del proyecto "p1"
    Y que la historia "HU-02" de 5 Story Points está asignada al Sprint número 1 del proyecto "p1"
    Cuando se cierra el Sprint número 1 del proyecto "p1"
    Y se consultan sus Story Points planificados
    Entonces los Story Points planificados del Sprint número 1 son 8

  @D12
  Escenario: Consultar Story Points planificados de un Sprint activo en vivo
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que la historia "HU-01" con 3 Story Points está asignada al Sprint número 1 del proyecto "p1"
    Y que la historia "HU-02" con 5 Story Points está asignada al Sprint número 1 del proyecto "p1"
    Cuando se consultan los Story Points planificados del Sprint número 1 del proyecto "p1"
    Entonces los Story Points planificados del Sprint número 1 son 8

  @D13
  Escenario: Cerrar un Sprint sin historias
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo" y sin historias asignadas
    Cuando se cierra el Sprint número 1 del proyecto "p1"
    Entonces el Sprint número 1 del proyecto "p1" queda en estado "Cerrado"
    Y los Story Points planificados del Sprint número 1 son 0
    Y la velocidad del Sprint número 1 es 0

  @D13
  Escenario: Cerrar un Sprint con todas sus historias completadas
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Y que la historia "HU-01" está completada, asignada al Sprint número 1 y tiene 3 Story Points
    Y que la historia "HU-02" está completada, asignada al Sprint número 1 y tiene 5 Story Points
    Cuando se cierra el Sprint número 1 del proyecto "p1"
    Entonces el Sprint número 1 del proyecto "p1" queda en estado "Cerrado"
    Y las historias "HU-01" y "HU-02" siguen asignadas al Sprint número 1 del proyecto "p1"
    Y los Story Points planificados del Sprint número 1 son 8
    Y la velocidad del Sprint número 1 es 8

  @CA4
  Escenario: Cerrar un Sprint cerrado devuelve error y no modifica sus datos
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Cerrado"
    Cuando se intenta cerrar el Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint está cerrado
    Y el Sprint número 1 del proyecto "p1" permanece en estado "Cerrado"

  @CA3
  @D5
  Escenario: No asignar una historia a un Sprint cerrado
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Cerrado"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Cuando se intenta asignar la historia "HU-01" al Sprint número 1 del proyecto "p1"
    Entonces se devuelve un error indicando que el Sprint está cerrado
    Y la historia "HU-01" no queda asignada a ningún Sprint

  @D9
  Escenario: Crear otro Sprint activo en el proyecto después de cerrar el Sprint actual
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene el Sprint número 1 en estado "Activo"
    Cuando se cierra el Sprint número 1 del proyecto "p1"
    Y se crea el Sprint número 2 en el proyecto "p1" con objetivo "Planificar la segunda iteración", fecha de inicio "23 de marzo de 2026" y fecha de fin "10 de abril de 2026"
    Entonces el Sprint número 2 se crea exitosamente y queda en estado "Activo"

  @D7
  Escenario: Listar los Sprints del proyecto en orden ascendente por número
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene Sprints con números 3, 1 y 2
    Cuando se listan los Sprints del proyecto "p1"
    Entonces se devuelven los Sprints en el orden 1, 2 y 3
