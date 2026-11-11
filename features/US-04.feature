# language: es
# US-04: los escenarios cubren R1 (estado por fechas), R3 (porcentaje de historias
# completadas), R4 (proyecto sin historias) y R5 (solo lectura del estado guardado),
# según specs/US-04.md y las decisiones D1, D4, D6, D7 y D9.
# No se escenifican R2 (sprint actual: la spec no define una función para mostrarlo)
# ni R6 (proyecto inexistente: depende de P5), ni los pendientes P4 y P5.
# El estado En curso incluye ambos extremos: la fecha de inicio y la fecha de fin.
# La fecha "hoy" se compara solo por día (año, mes y día), ignorando hora y zona
# horaria (D7); el porcentaje se redondea a 2 decimales (D4).
# Los errores se describen por su significado; sus identificadores exactos
# (ErrFechaFinAnteriorAInicio y ErrHistoriasInvalidas) están decididos en la
# spec (D6 y D9).

@US-04
Característica: Consultar estado de un proyecto
  Como Product Owner
  Quiero ver el estado de un proyecto
  Para saber cuánto avanzó

  @CA1
  Esquema del escenario: Estado del proyecto según la fecha de consulta
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026" y fecha de fin "30 de junio de 2026"
    Cuando se consulta el estado del proyecto con la fecha de hoy "<hoy>"
    Entonces el estado mostrado es "<estado>"

    Ejemplos:
      | hoy                   | estado      |
      | 28 de febrero de 2026 | Planificado |
      | 1 de marzo de 2026    | En curso    |
      | 15 de abril de 2026   | En curso    |
      | 30 de junio de 2026   | En curso    |
      | 1 de julio de 2026    | Finalizado  |

  @CA1
  Escenario: La misma fecha con distinta hora no cambia el estado
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026" y fecha de fin "30 de junio de 2026"
    Cuando se consulta el estado del proyecto con la fecha de hoy "1 de marzo de 2026 a las 8:00" y con la fecha de hoy "1 de marzo de 2026 a las 20:00"
    Entonces el estado mostrado es "En curso" en ambas consultas

  @CA1
  Escenario: Fecha de fin anterior a la fecha de inicio devuelve error al consultar el estado
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "30 de junio de 2026" y fecha de fin "1 de marzo de 2026"
    Cuando se consulta el estado del proyecto con la fecha de hoy "15 de abril de 2026"
    Entonces se devuelve un error indicando que la fecha de fin es anterior a la fecha de inicio
    Y la consulta no se completa exitosamente

  @CA1
  Escenario: El estado calculado puede diferir del estado guardado y no lo modifica
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026", fecha de fin "30 de junio de 2026" y estado guardado "Planificado"
    Cuando se consulta el estado del proyecto con la fecha de hoy "1 de julio de 2026"
    Entonces el estado mostrado es "Finalizado"
    Y el proyecto conserva el estado guardado "Planificado"

  @CA1
  Esquema del escenario: Porcentaje de historias completadas redondeado a 2 decimales
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026" y fecha de fin "30 de junio de 2026"
    Cuando se consulta el porcentaje de historias completadas con total de historias <total> y completadas <completadas>
    Entonces el porcentaje mostrado es <porcentaje>%

    Ejemplos:
      | total | completadas | porcentaje |
      | 3     | 0           | 0          |
      | 3     | 1           | 33.33      |
      | 3     | 2           | 66.67      |
      | 3     | 3           | 100        |
      | 8     | 5           | 62.5       |

  @CA1
  Esquema del escenario: Porcentaje con valores inválidos devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026" y fecha de fin "30 de junio de 2026"
    Cuando se consulta el porcentaje de historias completadas con total de historias <total> y completadas <completadas>
    Entonces se devuelve un error indicando que las cantidades de historias son inválidas
    Y la consulta no se completa exitosamente

    Ejemplos:
      | total | completadas |
      | -1    | 0           |
      | 3     | -1          |
      | 3     | 4           |

  @CA2
  Escenario: Proyecto sin historias muestra 0% sin error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026" y fecha de fin "30 de junio de 2026"
    Cuando se consulta el porcentaje de historias completadas con total de historias 0 y completadas 0
    Entonces el porcentaje mostrado es 0%
    Y la consulta se realiza sin errores

  @CA2
  Escenario: Total de historias 0 con completadas mayores que 0 devuelve error
    Dado un proyecto existente identificado por su ID, con nombre "Sistema de Métricas", fecha de inicio "1 de marzo de 2026" y fecha de fin "30 de junio de 2026"
    Cuando se consulta el porcentaje de historias completadas con total de historias 0 y completadas 1
    Entonces se devuelve un error indicando que las cantidades de historias son inválidas
    Y la consulta no se completa exitosamente
