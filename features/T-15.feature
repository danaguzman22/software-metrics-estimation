# language: es
# Estado: Decisiones validadas por la persona (2026-11-16)
# Los escenarios describen el comportamiento acordado para CA1-CA6 y las
# decisiones DECIDIDAS 1-8 de specs/T-15.md. Todas las pantallas de la
# jerarquía muestran una migaja de pan con todos los niveles clicables y un
# enlace "← Volver" al nivel padre inmediato, renderizados una única vez en
# base.html a partir de la lista de migas que arma cada handler. Volver desde
# Planning Poker y desde Estimar/Editar lleva al backlog anclado a la historia
# (#item-<ID>, p. ej. #item-HU-01); volver desde el reporte del Sprint
# lleva a Sprints anclado al
# Sprint (#sprint-<numero>). Las páginas de error (404/500) no muestran
# migas; la navegación de US-22 y los enlaces existentes se conservan; sin
# JavaScript y con salida escapada.
# Decisiones del 2026-11-16 (validadas por la persona): C1 resuelto — el
# ancla del ítem usa el ID completo (#item-HU-01); C2 resuelto — Estimar es
# una acción inline (el 422 re-renderiza el backlog), sin pantalla propia;
# C3 resuelto — /proyectos/nuevo entra en la jerarquía como
# Proyectos › Nuevo proyecto (el inicio / y el placeholder global /defectos
# quedan sin migas); C7 resuelto — el ancla del Sprint va en el
# <h2 id="sprint-N">.
# Los textos exactos, el tipo y nombre del helper y las pantallas fuera de la
# jerarquía son DERIVADAS de la spec (secciones 4, 7 y 16), confirmadas por la
# persona (2026-11-16); este archivo las asume tal como quedaron decididas.
# No se escenifican cosas fuera de alcance: nuevas rutas, cambios de
# servicios/dominio/repositorios/migraciones, PDF nuevos, autenticación ni
# cambios en la navegación principal de US-22.

@T-15
Característica: Navegación jerárquica en todas las pantallas
  Como persona que usa la aplicación
  Quiero ver una migaja de pan y un enlace "← Volver" en cada pantalla
  Para volver siempre al lugar de origen desde cualquier sección

  Antecedentes:
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas"
    Y que el proyecto "PRJ-001" tiene la historia "HU-01" con título "Iniciar sesión", estado "Pendiente" y 3 Story Points
    Y que el proyecto "PRJ-001" tiene el Sprint número 1 con objetivo "Primera entrega"

  @CA1
  Esquema del escenario: Cada pantalla muestra sus migas con todos los niveles clicables
    Cuando se abre la ruta "<ruta>"
    Entonces la página muestra la migaja de pan con los niveles "<migas>"
    Y cada nivel de la migaja es un enlace clicable

    Ejemplos:
      | ruta                                           | migas                                                          |
      | /proyectos                                      | Proyectos                                                      |
      | /proyectos/nuevo                                | Proyectos › Nuevo proyecto                                     |
      | /proyectos/PRJ-001                              | Proyectos › Sistema de Métricas                                |
      | /proyectos/PRJ-001/backlog                      | Proyectos › Sistema de Métricas › Backlog                      |
      | /proyectos/PRJ-001/backlog/nuevo                | Proyectos › Sistema de Métricas › Backlog › Nuevo ítem         |
      | /proyectos/PRJ-001/backlog/HU-01/editar         | Proyectos › Sistema de Métricas › Backlog › Editar ítem        |
      | /proyectos/PRJ-001/backlog/HU-01/poker          | Proyectos › Sistema de Métricas › Backlog › Planning Poker     |
      | /proyectos/PRJ-001/sprints                      | Proyectos › Sistema de Métricas › Sprints                      |
      | /proyectos/PRJ-001/sprints/nuevo                | Proyectos › Sistema de Métricas › Sprints › Nuevo Sprint       |
      | /proyectos/PRJ-001/sprints/historico            | Proyectos › Sistema de Métricas › Sprints › Histórico          |
      | /proyectos/PRJ-001/sprints/1/reporte            | Proyectos › Sistema de Métricas › Sprints › Reporte del Sprint 1 |
      | /proyectos/PRJ-001/defectos                     | Proyectos › Sistema de Métricas › Defectos                     |
      | /proyectos/PRJ-001/integrantes                  | Proyectos › Sistema de Métricas › Integrantes                  |
      | /proyectos/PRJ-001/dashboard                    | Proyectos › Sistema de Métricas › Dashboard                    |
      | /proyectos/PRJ-001/reporte                      | Proyectos › Sistema de Métricas › Reporte del proyecto         |

  @CA2
  Esquema del escenario: Cada pantalla muestra "← Volver" al nivel padre inmediato
    Cuando se abre la ruta "<ruta>"
    Entonces la página muestra el enlace "← Volver" apuntando a "<volver>"

    Ejemplos:
      | ruta                                         | volver                                      |
      | /proyectos/nuevo                             | /proyectos                                  |
      | /proyectos/PRJ-001                           | /proyectos                                  |
      | /proyectos/PRJ-001/backlog                   | /proyectos/PRJ-001                          |
      | /proyectos/PRJ-001/backlog/nuevo             | /proyectos/PRJ-001/backlog                  |
      | /proyectos/PRJ-001/backlog/HU-01/editar      | /proyectos/PRJ-001/backlog#item-HU-01       |
      | /proyectos/PRJ-001/backlog/HU-01/poker       | /proyectos/PRJ-001/backlog#item-HU-01       |
      | /proyectos/PRJ-001/sprints                   | /proyectos/PRJ-001                          |
      | /proyectos/PRJ-001/sprints/nuevo             | /proyectos/PRJ-001/sprints                  |
      | /proyectos/PRJ-001/sprints/historico         | /proyectos/PRJ-001/sprints                  |
      | /proyectos/PRJ-001/sprints/1/reporte         | /proyectos/PRJ-001/sprints#sprint-1         |
      | /proyectos/PRJ-001/defectos                  | /proyectos/PRJ-001                          |
      | /proyectos/PRJ-001/integrantes               | /proyectos/PRJ-001                          |
      | /proyectos/PRJ-001/dashboard                 | /proyectos/PRJ-001                          |
      | /proyectos/PRJ-001/reporte                   | /proyectos/PRJ-001                          |

  @CA1
  Escenario: La raíz de la jerarquía no muestra "← Volver"
    Cuando se abre la página de proyectos
    Entonces la página muestra la migaja de pan con el nivel "Proyectos"
    Y no muestra el enlace "← Volver"

  @CA1
  Escenario: Las páginas fuera de la jerarquía no muestran migas
    Cuando se abre la página de inicio
    Entonces la página no muestra la migaja de pan
    Cuando se abre la sección global de defectos
    Entonces la página no muestra la migaja de pan
    Y no muestra el enlace "← Volver"

  @CA1
  Escenario: Las migas se renderizan una sola vez por página
    Cuando se abre el backlog del proyecto "PRJ-001"
    Entonces la página contiene un único bloque de migaja de pan

  @CA2
  Escenario: El backlog ancla cada fila de historia para la vuelta anclada
    Cuando se abre el backlog del proyecto "PRJ-001"
    Entonces la fila de la historia "HU-01" lleva el ancla "item-HU-01"

  @CA2
  Escenario: Los Sprints anclan su encabezado para la vuelta anclada
    Cuando se abre la sección de sprints del proyecto "PRJ-001"
    Entonces el encabezado del Sprint número 1 lleva el ancla "sprint-1"

  @CA2
  Escenario: Volver desde Planning Poker lleva al backlog anclado a la historia
    Cuando se abre la página de Planning Poker de la historia "HU-01" del proyecto "PRJ-001"
    Entonces la página muestra el enlace "← Volver" apuntando a "/proyectos/PRJ-001/backlog#item-HU-01"

  @CA2
  Escenario: Volver desde la edición de una historia lleva al backlog anclado a la historia
    Cuando se abre el formulario de edición de la historia "HU-01" del proyecto "PRJ-001"
    Entonces la página muestra el enlace "← Volver" apuntando a "/proyectos/PRJ-001/backlog#item-HU-01"

  @CA2
  Escenario: Volver desde el reporte del Sprint lleva a Sprints anclado al Sprint
    Cuando se abre el reporte del Sprint número 1 del proyecto "PRJ-001"
    Entonces la página muestra el enlace "← Volver" apuntando a "/proyectos/PRJ-001/sprints#sprint-1"

  @CA3
  Escenario: Las migas y "← Volver" se ocultan al imprimir un reporte
    Cuando se abre el reporte del proyecto "PRJ-001"
    Entonces la página muestra la migaja de pan
    Y la impresión de la página no muestra la migaja de pan ni el enlace "← Volver"

  @CA3
  Escenario: El PDF del reporte no incluye las migas ni "← Volver"
    Cuando se descarga el reporte del proyecto "PRJ-001" en PDF
    Entonces el archivo descargado no contiene el bloque de migaja de pan
    Y no contiene el enlace "← Volver"

  @CA4
  Escenario: Un proyecto inexistente muestra el error sin migas ni "← Volver"
    Dado que no existe el proyecto "PRJ-999"
    Cuando se abre el backlog del proyecto "PRJ-999"
    Entonces la página responde con el código 404
    Y muestra el mensaje de proyecto no encontrado
    Y la página no muestra la migaja de pan
    Y no muestra el enlace "← Volver"

  @CA4
  Escenario: Una historia inexistente muestra el error sin migas ni "← Volver"
    Dado que el proyecto "PRJ-001" no tiene la historia "HU-99"
    Cuando se intenta editar la historia "HU-99" del proyecto "PRJ-001"
    Entonces la página responde con el código 404
    Y muestra el mensaje de historia no encontrada
    Y la página no muestra la migaja de pan
    Y no muestra el enlace "← Volver"

  @CA4
  Escenario: Un Sprint inexistente muestra el error sin migas ni "← Volver"
    Cuando se abre el reporte del Sprint número 9 del proyecto "PRJ-001"
    Entonces la página responde con el código 404
    Y muestra el mensaje de Sprint no encontrado
    Y la página no muestra la migaja de pan
    Y no muestra el enlace "← Volver"

  @CA4
  Escenario: Un Sprint con número inválido muestra el error sin migas
    Cuando se abre el reporte del Sprint "abc" del proyecto "PRJ-001"
    Entonces la página responde con el código 404
    Y la página no muestra la migaja de pan
    Y no muestra el enlace "← Volver"

  @CA4
  Escenario: La navegación principal de US-22 se conserva en todas las páginas
    Cuando se abre la página de inicio
    Y se abre la página de proyectos
    Y se abre el detalle del proyecto "PRJ-001"
    Y se abre el backlog del proyecto "PRJ-001"
    Y se abre la sección de sprints del proyecto "PRJ-001"
    Y se abre la sección de defectos
    Entonces todas las páginas muestran los enlaces "Proyectos", "Backlog", "Sprints", "Dashboard" y "Defectos"

  @CA4
  Escenario: Los enlaces propios de cada pantalla se conservan
    Cuando se abre el detalle del proyecto "PRJ-001"
    Entonces la página muestra los enlaces "Ver backlog", "Ver sprints", "Ver defectos", "Ver dashboard", "Ver reporte" e "Ver integrantes"
    Cuando se abre el backlog del proyecto "PRJ-001"
    Entonces la historia "HU-01" muestra los enlaces "Editar" y "Planning Poker"
    Cuando se abre el reporte del proyecto "PRJ-001"
    Entonces la página muestra el botón "Imprimir" y el enlace "Descargar PDF"

  @CA5
  Escenario: Un nombre de proyecto con markup se muestra escapado en las migas
    Dado un proyecto existente identificado por "PRJ-002" con nombre "<script>alert(1)</script>"
    Cuando se abre el detalle del proyecto "PRJ-002"
    Entonces la migaja de pan muestra el texto del nombre escapado
    Y la página no interpreta el nombre como HTML

  @CA5
  Escenario: Un nombre de proyecto vacío muestra el texto genérico en la miga
    Dado un proyecto existente identificado por "PRJ-003" con nombre ""
    Cuando se abre el detalle del proyecto "PRJ-003"
    Entonces la migaja de pan muestra el texto "Proyecto" en lugar del nombre

  @CA5
  Escenario: Un nombre de proyecto muy largo se muestra completo en las migas
    Dado un proyecto existente identificado por "PRJ-004" con nombre "Sistema de métricas y estimación para equipos distribuidos en múltiples zonas horarias con reportes semanales consolidados"
    Cuando se abre el detalle del proyecto "PRJ-004"
    Entonces la migaja de pan muestra el nombre completo

  @CA6
  Escenario: La leyenda de US-22 se conserva donde correspondía
    Cuando se abre la página de proyectos
    Entonces la página muestra la leyenda "Elegí un proyecto para ver su backlog y sus sprints"
    Y muestra la migaja de pan con el nivel "Proyectos"

  @CA6
  Escenario: La estimación sigue funcionando dentro del backlog
    Cuando se abre el backlog del proyecto "PRJ-001"
    Entonces la historia "HU-01" ofrece el campo de Story Points y el botón "Estimar"

  @CA1 @CA6
  Escenario: Un envío inválido del formulario muestra las migas de la pantalla
    Cuando se intenta editar la historia "HU-01" del proyecto "PRJ-001" con título vacío
    Entonces la página responde con el código 422
    Y muestra la migaja de pan con los niveles "Proyectos › Sistema de Métricas › Backlog › Editar ítem"
    Y muestra el enlace "← Volver" apuntando a "/proyectos/PRJ-001/backlog#item-HU-01"
