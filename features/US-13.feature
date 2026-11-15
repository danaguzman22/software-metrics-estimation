# language: es
# Estado: Propuesta
# Los escenarios describen US-13 y las decisiones D1–D7 aprobadas por Dana.
# D8, D11, D12 y D18–D26 (reúso de
# ErrStoryPointsInvalidos, acuerdo con ronda revelada sin votos permitido,
# semántica provisional de diferencias con ronda abierta) quedan pendientes de
# confirmación. ErrRondaSinVotos se usa únicamente en diferencias con ronda
# revelada sin votos; acordar con ronda revelada sin votos está permitido
# (D18). El mensaje de ErrRondaNoRevelada es "la ronda aún no fue revelada".
# El orden de evaluación al acordar es: proyecto, historia,
# integrante del proyecto, ronda, estado, valor. Con la ronda Abierta no se
# informan diferencias ni valores (privacidad de US-12). La implementación usa
# provisionalmente los servicios DiferenciasRonda, HistorialRondas y
# AcordarEstimacion, el constructor NuevoServicioPlanningPokerConEstimador
# (US-13) y reutiliza el caso de uso de estimar de US-11 (detalle en
# specs/US-13.md).
# FASE 2 (sección 20 de specs/US-13.md): la interfaz web en la página de poker
# por historia; sus escenarios están marcados con @Fase2 y describen las
# secciones "Diferencias", "Historial" y el formulario "Acordar estimación".

@US-13
Característica: Planning Poker: diferencias, nuevas rondas y estimación acordada
  Como equipo
  Quiero detectar diferencias entre votos y repetir rondas hasta acordar
  Para fijar una estimación consensuada

  @CA1
  Escenario: Una ronda revelada con votos distintos informa diferencia con mínimo, máximo y quién votó cada uno
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana", "Vicky" y "Marcos" con roles "Product Owner", "Product Builder" y "Ingeniero"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Vicky" votó 8 en la ronda número 1 antes de revelarla
    Y que el integrante "Marcos" votó 13 en la ronda número 1 antes de revelarla
    Cuando se consultan las diferencias de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se informa que hay diferencia en la ronda número 1
    Y se informa que el valor mínimo es 5 y lo votó "Ana"
    Y se informa que el valor máximo es 13 y lo votó "Marcos"

  @CA1
  @R2
  Escenario: Quiénes votaron cada extremo incluye a todos los votantes empatados
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" y el integrante "Marcos" con roles "Product Owner", "Product Builder" e "Ingeniero"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Vicky" votó 8 en la ronda número 1 antes de revelarla
    Y que el integrante "Marcos" votó 5 en la ronda número 1 antes de revelarla
    Cuando se consultan las diferencias de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se informa que hay diferencia en la ronda número 1
    Y se informa que el valor mínimo es 5 y lo votaron "Ana" y "Marcos"
    Y se informa que el valor máximo es 8 y lo votó "Vicky"

  @CA1
  @R2
  Escenario: Una ronda revelada con un solo voto no informa diferencia
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 8 en la ronda número 1 antes de revelarla
    Cuando se consultan las diferencias de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se informa que no hay diferencia en la ronda número 1
    Y se informa que el valor mínimo es 8 y el valor máximo es 8
    Y se informa que lo votó "Ana"

  @CA1
  @R2
  Escenario: Una ronda revelada con todos los votos iguales no informa diferencia
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Vicky" votó 5 en la ronda número 1 antes de revelarla
    Cuando se consultan las diferencias de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se informa que no hay diferencia en la ronda número 1
    Y se informa que el valor mínimo es 5 y el valor máximo es 5
    Y se informa que lo votaron "Ana" y "Vicky"

  @CA1
  @E8
  Escenario: Una ronda revelada sin votos devuelve un error propio
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Cuando se consultan las diferencias de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda no tiene votos

  @CA1
  @E9
  Escenario: Con la ronda abierta no se informan las diferencias ni los valores
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 13 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando se consultan las diferencias de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda aún no fue revelada
    Y la ronda número 1 no informa los valores de los votos

  @CA1
  @E5
  Escenario: Ronda inexistente al consultar diferencias
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando se consultan las diferencias de la ronda número 9 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda no existe para la historia

  @CA1
  @E4
  Escenario: Número de ronda inválido al consultar diferencias
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Cuando se consultan las diferencias de la ronda número 0 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el número de la ronda debe ser mayor que cero

  @CA2
  Escenario: Tras revelar se inicia una nueva ronda y la anterior se conserva en el historial
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando se inicia una nueva ronda sobre la historia "HU-01" del proyecto "p1"
    Entonces se crea la ronda número 2 sobre la historia "HU-01" en estado "Abierta"
    Y el historial de la historia "HU-01" conserva la ronda número 1 con sus votos

  @CA2
  Escenario: El historial informa los votos de las rondas reveladas y nunca los valores de una ronda abierta
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que el proyecto "p1" tiene la ronda número 2 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Vicky" vota 8 en la ronda número 2 de la historia "HU-01" del proyecto "p1"
    Cuando se consulta el historial de rondas de la historia "HU-01" del proyecto "p1"
    Entonces el historial informa la ronda número 1 en estado "Revelada" con el voto de "Ana" con valor 5
    Y el historial informa la ronda número 2 en estado "Abierta" sin informar valores
    Y el historial informa las rondas en orden ascendente por número

  @CA2
  Escenario: El historial de una historia sin rondas es una lista vacía
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Cuando se consulta el historial de rondas de la historia "HU-01" del proyecto "p1"
    Entonces el historial informa cero rondas
    Y no se devuelve un error

  @CA2
  @R16
  Escenario: No se inicia una nueva ronda mientras la anterior está abierta
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando se intenta iniciar una nueva ronda sobre la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la historia ya tiene una ronda abierta

  @CA3
  Escenario: Acordar la estimación guarda los Story Points de la historia
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Vicky" votó 8 en la ronda número 1 antes de revelarla
    Cuando la integrante "Ana" acuerda el valor 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la historia "HU-01" del proyecto "p1" queda con 8 Story Points

  @CA3
  Escenario: Acordar un valor que nadie votó está permitido
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Vicky" votó 8 en la ronda número 1 antes de revelarla
    Cuando la integrante "Ana" acuerda el valor 13 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la historia "HU-01" del proyecto "p1" queda con 13 Story Points

  @CA3
  Escenario: Un integrante del proyecto que no votó puede acordar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" y el integrante "Marcos" con roles "Product Owner", "Product Builder" e "Ingeniero"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando el integrante "Marcos" acuerda el valor 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la historia "HU-01" del proyecto "p1" queda con 5 Story Points

  @CA3
  @R18
  Escenario: Acordar sobre una historia ya estimada reemplaza el valor anterior
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 5 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 8 en la ronda número 1 antes de revelarla
    Cuando la integrante "Ana" acuerda el valor 13 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la historia "HU-01" del proyecto "p1" queda con 13 Story Points

  @CA3
  Escenario: Volver a acordar en una ronda posterior reemplaza el valor acordado en la ronda anterior
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Ana" acuerda el valor 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y que se inicia una nueva ronda sobre la historia "HU-01" del proyecto "p1"
    Y que la integrante "Ana" vota 13 en la ronda número 2 de la historia "HU-01" del proyecto "p1"
    Y que la integrante "Ana" revela la ronda número 2 de la historia "HU-01" del proyecto "p1"
    Cuando la integrante "Ana" acuerda el valor 13 en la ronda número 2 de la historia "HU-01" del proyecto "p1"
    Entonces la historia "HU-01" del proyecto "p1" queda con 13 Story Points

  @CA3
  @R17
  Escenario: Acordar sobre una ronda revelada sin votos está permitido
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Cuando la integrante "Ana" acuerda el valor 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la historia "HU-01" del proyecto "p1" queda con 8 Story Points

  @CA3
  @R6
  Esquema del escenario: Acordar un valor fuera de la serie es rechazado y la historia conserva sus Story Points
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 5 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando la integrante "Ana" acuerda el valor <valor> en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que los story points no son válidos
    Y la historia "HU-01" del proyecto "p1" queda con 5 Story Points

    Ejemplos:
      | valor |
      | 0     |
      | 4     |
      | 22    |

  @CA3
  @E5
  Escenario: Acordar sobre una ronda inexistente
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando la integrante "Ana" acuerda el valor 5 en la ronda número 9 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda no existe para la historia

  @CA3
  @E4
  Escenario: Número de ronda inválido al acordar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Cuando la integrante "Ana" acuerda el valor 5 en la ronda número 0 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el número de la ronda debe ser mayor que cero

  @CA4
  Escenario: Acordar sobre una ronda abierta no está permitido y no modifica los Story Points
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando la integrante "Ana" acuerda el valor 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda aún no fue revelada
    Y la historia "HU-01" del proyecto "p1" queda con 0 Story Points

  @CA4
  Escenario: Un integrante que no es del proyecto no puede acordar ni con la ronda revelada
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando el integrante "Marcos" acuerda el valor 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el integrante no pertenece al proyecto
    Y la historia "HU-01" del proyecto "p1" queda con 0 Story Points

  @CA4
  @E10
  Escenario: Nombre de integrante vacío al acordar se rechaza al inicio del paso 3
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Cuando el integrante "   " acuerda el valor 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el nombre del integrante no puede estar vacío

  @CA4
  @R12
  Escenario: Proyecto inexistente prevalece sobre historia, integrante, ronda y valor al acordar
    Dado que no existe un proyecto identificado por "p999"
    Cuando el integrante "Ana" acuerda el valor 5 en la ronda número 9 de la historia "HU-99" del proyecto "p999"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de historia inexistente

  @CA4
  @R12
  Escenario: Historia inexistente prevalece sobre integrante inexistente al acordar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Y que el proyecto "p1" no tiene la integrante "Marcos"
    Cuando el integrante "Marcos" acuerda el valor 5 en la ronda número 1 de la historia "HU-99" del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no existe en el proyecto

  @CA4
  @R12
  Escenario: Integrante inexistente prevalece sobre ronda inexistente al acordar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene la integrante "Marcos"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando el integrante "Marcos" acuerda el valor 5 en la ronda número 9 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el integrante no pertenece al proyecto
    Y no se devuelve un error de ronda inexistente

  @CA1 @Fase2
  Escenario: Tras revelar, la página muestra las diferencias con mínimo, máximo y quién votó cada uno
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Vicky" votó 13 en la ronda número 1 antes de revelarla
    Cuando se visita la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página informa que hay diferencia
    Y la página informa que el valor mínimo es 5 y lo votó "Ana"
    Y la página informa que el valor máximo es 13 y lo votó "Vicky"

  @CA1 @Fase2
  Escenario: Con la ronda abierta la página no muestra las diferencias ni los valores del historial
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando se visita la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página no informa diferencias
    Y la página no informa los valores de los votos

  @CA2 @Fase2
  Escenario: El historial lista las rondas reveladas con sus votos y las abiertas como votos ocultos
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que el proyecto "p1" tiene la ronda número 2 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Vicky" vota 8 en la ronda número 2 de la historia "HU-01" del proyecto "p1"
    Cuando se visita la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página lista la ronda número 1 con el voto de "Ana" con valor 5
    Y la página lista la ronda número 2 como votos ocultos

  @CA3 @Fase2
  Escenario: Acordar desde la página guarda los Story Points y redirige a la página
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Y que la integrante "Vicky" votó 8 en la ronda número 1 antes de revelarla
    Cuando la integrante "Ana" acuerda el valor 8 en la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces se redirige a la página de poker de la historia "HU-01" del proyecto "p1"
    Y la historia "HU-01" del proyecto "p1" queda con 8 Story Points

  @CA3 @Fase2
  Escenario: Volver a acordar desde la página en una ronda posterior reemplaza el valor anterior
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 8 en la ronda número 1 antes de revelarla
    Y que la integrante "Ana" acuerda el valor 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y que se inicia una nueva ronda sobre la historia "HU-01" del proyecto "p1"
    Y que la integrante "Ana" vota 13 en la ronda número 2 de la historia "HU-01" del proyecto "p1"
    Y que la integrante "Ana" revela la ronda número 2 de la historia "HU-01" del proyecto "p1"
    Cuando la integrante "Ana" acuerda el valor 13 en la ronda número 2 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces se redirige a la página de poker de la historia "HU-01" del proyecto "p1"
    Y la historia "HU-01" del proyecto "p1" queda con 13 Story Points

  @CA4 @Fase2
  Escenario: Acordar desde la página con la ronda abierta muestra el mensaje del dominio
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando se envía un acuerdo del valor 5 en la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra el mensaje "la ronda aún no fue revelada"
    Y la historia "HU-01" del proyecto "p1" queda con 0 Story Points

  @CA3 @Fase2
  Escenario: Un valor fuera de la serie enviado desde el formulario muestra el mensaje del dominio
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando se envía un acuerdo del valor 4 en la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra el mensaje "los story points no son válidos: deben ser 1, 2, 3, 5, 8, 13 o 21"

  @CA3 @E5 @Fase2
  Escenario: Acordar sobre una ronda inexistente muestra el mensaje del dominio
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando se envía un acuerdo del valor 5 en la ronda número 9 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra el mensaje "la ronda no existe para la historia"
