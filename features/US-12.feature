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
# FASE 2 (sección 20 de specs/US-12.md): la interfaz web de Planning Poker por
# historia; sus escenarios están marcados con @Fase2 y describen una página por
# historia para iniciar la ronda actual, votar en oculto, revelar y abrir la
# siguiente ronda. La exclusión web de D8 se levanta solo para esta fase; el
# consenso y la estimación final siguen siendo US-13.

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

  @CA1 @Fase2
  Escenario: La página de poker ofrece iniciar una ronda cuando la historia no tiene rondas
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Cuando se visita la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página ofrece el botón "Iniciar ronda"
    Y la página no muestra un formulario de voto

  @CA1 @Fase2
  Escenario: Iniciar una ronda desde la página crea la ronda 1
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Cuando se inicia una ronda desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces se redirige a la página de poker de la historia "HU-01" del proyecto "p1"
    Y se crea la ronda número 1 sobre la historia "HU-01" en estado "Abierta"

  @CA1 @D7 @Fase2
  Escenario: Con la ronda anterior revelada, la página ofrece abrir la siguiente ronda
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Cuando se inicia una ronda desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces se redirige a la página de poker de la historia "HU-01" del proyecto "p1"
    Y se crea la ronda número 2 sobre la historia "HU-01" en estado "Abierta"

  @CA2 @Fase2
  Escenario: Con la ronda abierta la página lista quiénes votaron sin valores
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y que la integrante "Vicky" vota 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando se visita la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página informa que ya votaron "Ana" y "Vicky"
    Y la página no informa los valores de los votos

  @CA2 @Fase2
  Escenario: Con la ronda abierta la página no muestra los valores de los votos
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando se visita la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página no informa los valores de los votos
    Y la página muestra el formulario de voto
    Y la página muestra el formulario "Revelar"

  @CA1 @Fase2
  Escenario: Votar desde la página registra el voto
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Ana" vota 5 en la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces se redirige a la página de poker de la historia "HU-01" del proyecto "p1"
    Y la ronda número 1 informa que ya votó 1 integrante

  @CA4 @Fase2
  Escenario: Votar desde la página con un integrante ajeno muestra el error del dominio
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando el integrante "Marcos" vota 5 en la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra el mensaje "el integrante no pertenece al proyecto"
    Y la ronda número 1 no informa que "Marcos" haya votado

  @CA2 @Fase2
  Escenario: Revelar desde la página muestra todos los votos con sus valores
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene las integrantes "Ana" y "Vicky" con roles "Product Owner" y "Product Builder"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Y que la integrante "Ana" vota 5 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Y que la integrante "Vicky" vota 8 en la ronda número 1 de la historia "HU-01" del proyecto "p1"
    Cuando la integrante "Ana" revela la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces se redirige a la página de poker de la historia "HU-01" del proyecto "p1"
    Y la página informa los votos de "Ana" con valor 5 y de "Vicky" con valor 8

  @CA2 @Fase2
  Escenario: Revelar dos veces desde la página muestra el error del dominio
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando la integrante "Ana" revela la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Y la integrante "Ana" revela de nuevo la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra el mensaje "la ronda ya fue revelada"
    Y la ronda número 1 permanece en estado "Revelada"

  @CA1 @CA3 @Fase2
  Escenario: Votar sobre una ronda revelada desde la página muestra el error del dominio
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Revelada"
    Y que la integrante "Ana" votó 5 en la ronda número 1 antes de revelarla
    Cuando el integrante "Ana" vota 8 en la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra el mensaje "la ronda ya fue revelada"
    Y el voto de "Ana" en la ronda número 1 permanece en valor 5

  @CA1 @Fase2
  Escenario: Un valor no numérico en el formulario de voto produce un mensaje de formato
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la ronda número 1 sobre la historia "HU-01" en estado "Abierta"
    Cuando se envía un voto con valor "cinco" en la ronda número 1 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra un mensaje indicando que el voto debe ser un número entero

  @CA1 @E5 @Fase2
  Escenario: Votar sobre una ronda inexistente muestra un mensaje
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la integrante "Ana" con rol "Product Owner"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" no tiene la ronda número 9 sobre la historia "HU-01"
    Cuando el integrante "Ana" vota 5 en la ronda número 9 desde la página de poker de la historia "HU-01" del proyecto "p1"
    Entonces la página muestra el mensaje "la ronda no existe para la historia"

  @CA2 @E2 @Fase2
  Escenario: La historia inexistente muestra un mensaje
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" no tiene la historia "HU-99"
    Cuando se visita la página de poker de la historia "HU-99" del proyecto "p1"
    Entonces la página muestra el mensaje "el ítem no existe en el proyecto"

  @CA4 @R1 @Fase2
  Escenario: El proyecto inexistente muestra un mensaje
    Dado que no existe un proyecto identificado por "p999"
    Cuando se visita la página de poker de la historia "HU-01" del proyecto "p999"
    Entonces la página muestra el mensaje "no se encontró el proyecto"

  @CA1 @Fase2
  Escenario: El backlog muestra el enlace Planning Poker por historia
    Dado un proyecto existente identificado por "p1", con nombre "Sistema de Métricas"
    Y que el proyecto "p1" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 0 Story Points
    Y que el proyecto "p1" tiene la historia "HU-02" con título "Recuperar contraseña", estado "Pendiente" y 0 Story Points
    Cuando se visita el backlog del proyecto "p1"
    Entonces cada historia del backlog muestra un enlace "Planning Poker" hacia su página de poker
