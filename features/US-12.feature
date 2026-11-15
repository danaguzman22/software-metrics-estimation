# language: es
# Estado: Propuesta
# Los escenarios describen US-12 y las decisiones D1–D11, D14, D17, D19–D22
# aprobadas por Dana. D12, D13, D15, D16 y D18 quedan pendientes de
# confirmación. El orden de evaluación al votar es: proyecto, historia,
# integrante del proyecto, ronda, valor. Cuando la ronda está Abierta el
# servicio solo informa quiénes ya votaron (cantidad y nombres), nunca los
# valores; revelar es irreversible. La implementación usa provisionalmente el
# modelo de tablas rondas_poker y votos_poker, el contrato RepositorioRondas
# con las lecturas YaVotaron/Votos y los errores de dominio propuestos en
# specs/US-12.md. specs/US-14.md no existe: la identidad del integrante se
# decide por nombre canónico (D17, DECIDIDA).

@US-12
Característica: Planning Poker: votos ocultos y revelación
  Como equipo
  Quiero votar en secreto y revelar los votos juntos
  Para evitar el sesgo de anclaje

  @CA1
  Escenario: Registrar un voto en una ronda abierta
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces el voto de "Ana" en la ronda número 1 queda registrado con valor 5

  @CA1
  @R5
  Esquema del escenario: Registrar votos con valores de la serie
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Ana" vota <valor> en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces el voto de "Ana" en la ronda número 1 queda registrado con valor <valor>

    Ejemplos:
      | valor |
      | 1     |
      | 2     |
      | 3     |
      | 5     |
      | 8     |
      | 13    |
      | 21    |

  @CA1
  @R5
  Esquema del escenario: No se acepta un voto fuera de la serie
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Ana" vota <valor> en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el voto no es válido
    Y el voto de "Ana" en la ronda número 1 no queda registrado

    Ejemplos:
      | valor |
      | 0     |
      | 4     |
      | 6     |
      | 22    |

  @CA1
  @E2
  Escenario: La historia inexistente rechaza el voto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Cuando el integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-99" del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no existe en el proyecto

  @CA2
  Escenario: Con la ronda abierta no se informan los valores, solo quiénes votaron
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y que la integrante "Vicky" vota 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando se consultan los votos de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la ronda número 1 informa que ya votaron 2 integrantes
    Y la ronda número 1 informa los nombres "Ana" y "Vicky"
    Y la ronda número 1 no informa los valores de los votos
    Y la ronda número 1 informa que está "Abierta"

  @CA2
  Escenario: Revelar muestra todos los votos con sus valores
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y que la integrante "Vicky" vota 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando la integrante "Ana" revela la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y se consultan los votos de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la ronda número 1 informa que está "Revelada"
    Y la ronda número 1 informa los votos de "Ana" con valor 5 y de "Vicky" con valor 8

  @CA2
  Escenario: Revelar es irreversible
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando la integrante "Ana" revela la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y la integrante "Ana" intenta revelar de nuevo la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda ya fue revelada
    Y la ronda número 1 permanece en estado "Revelada"

  @CA2
  Escenario: Revelar una ronda con votos no vuelve a ocultarlos
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando se consultan los votos de la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la ronda número 1 informa los votos de "Ana" con valor 5

  @CA2
  @D9
  Escenario: Revelar una ronda sin votos está permitido
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando la integrante "Ana" revela la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la ronda número 1 queda en estado "Revelada"
    Y la ronda número 1 no informa votos

  @CA3
  Escenario: Cambiar el voto antes de revelar reemplaza el valor
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 3 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando el integrante "Ana" vota 13 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces el voto de "Ana" en la ronda número 1 queda registrado con valor 13
    Y la integrante "Ana" figura una única vez entre quienes votaron en la ronda número 1

  @CA3
  Escenario: No se puede votar dos veces en la misma ronda: el cambio cuenta una vez
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Ana" vota 3 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y el integrante "Ana" vota 3 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la integrante "Ana" figura una única vez entre quienes votaron en la ronda número 1

  @CA3
  Escenario: Tras revelar no se puede votar ni cambiar el voto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando el integrante "Ana" vota 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda ya fue revelada
    Y el voto de "Ana" en la ronda número 1 permanece en valor 5

  @CA3
  @R12
  Escenario: Dos votos concurrentes del mismo integrante dejan un único voto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando dos procesos concurrentes registran para "Ana" los votos 5 y 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la ronda número 1 informa que ya votó 1 integrante
    Y la ronda número 1 conserva uno de los dos valores registrados

  @CA3
  @R12
  Escenario: Votos concurrentes de integrantes distintos se persisten ambos
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando las integrantes "Ana" y "Vicky" votan concurrentemente en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces la ronda número 1 informa que ya votaron 2 integrantes

  @CA4
  Escenario: Un integrante que no es del proyecto no puede votar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Marcos" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el integrante no pertenece al proyecto
    Y la ronda número 1 no informa que "Marcos" haya votado

  @CA4
  Escenario: Un integrante que no es del proyecto no puede revelar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Marcos" revela la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el integrante no pertenece al proyecto
    Y la ronda número 1 permanece en estado "Abierta"

  @CA4
  @E12
  Escenario: Nombre de integrante vacío se rechaza al inicio del paso 3
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el nombre del integrante no puede estar vacío

  @CA4
  @E12
  Escenario: Nombre de integrante de solo espacios se rechaza al inicio del paso 3
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "   " vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el nombre del integrante no puede estar vacío

  @CA4
  @E12
  Escenario: Nombre de integrante vacío o de solo espacios al revelar se rechaza al inicio del paso 3
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "   " revela la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el nombre del integrante no puede estar vacío
    Y la ronda número 1 permanece en estado "Abierta"

  @R1
  Escenario: Proyecto inexistente prevalece sobre ronda y historia inexistentes al votar
    Dado que no existe un proyecto identificado por "p999"
    Cuando el integrante "Ana" vota 5 en la ronda número 9 de la historia "HU-99" del proyecto "p999"
    Entonces se devuelve un error indicando que el proyecto no existe
    Y no se devuelve un error de historia inexistente
    Y no se devuelve un error de ronda inexistente

  @R1
  Escenario: Historia inexistente prevalece sobre ronda inexistente al votar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Cuando el integrante "Ana" vota 5 en la ronda número 9 de la historia "HU-99" del proyecto "p1"
    Entonces se devuelve un error indicando que la historia no existe en el proyecto
    Y no se devuelve un error de ronda inexistente

  @R1
  Escenario: Integrante del proyecto inexistente prevalece sobre ronda inexistente al votar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene la integrante "Marcos"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Cuando el integrante "Marcos" vota 5 en la ronda número 9 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el integrante no pertenece al proyecto
    Y no se devuelve un error de ronda inexistente

  @R16
  Escenario: No se inicia una nueva ronda mientras la anterior está abierta
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando se intenta iniciar una nueva ronda sobre la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la historia ya tiene una ronda abierta
    Y la ronda número 1 permanece en estado "Abierta"

  @R16
  @D7
  Escenario: Tras revelar se puede abrir la siguiente ronda
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Cuando se inicia una nueva ronda sobre la historia "HU-01" del proyecto "p1"
    Entonces se crea la ronda número 2 sobre la historia "HU-01" en estado "Abierta"

  @E5
  Escenario: Ronda inexistente al votar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando el integrante "Ana" vota 5 en la ronda número 9 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda no existe para la historia

  @E5
  Escenario: Ronda inexistente al revelar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando la integrante "Ana" revela la ronda número 9 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda no existe para la historia

  @E5
  Escenario: Ronda inexistente al consultar
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando se consultan los votos de la ronda número 9 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que la ronda no existe para la historia

  @E4
  @D14
  Escenario: Número de ronda inválido devuelve error de dominio en el paso 4
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Cuando el integrante "Ana" vota 5 en la ronda número 0 de la historia "HU-01" del proyecto "p1"
    Entonces se devuelve un error indicando que el número de la ronda debe ser mayor que cero
