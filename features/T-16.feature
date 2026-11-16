<<<<<<< HEAD
# language: es
# Estado: Especificación (documental, no ejecutada por la suite Go — no hay godog)
# Fuente: enunciado de la tarea T-16 "Mejora visual de la interfaz" (sin issue asignado;
# el alcance lo define el enunciado). Las decisiones DECIDIDAS y DERIVADAS están en
# specs/T-16.md. La verificación real es la suite en verde + la revisión visual humana
# (specs/T-16.md, secciones 10 y 12).

@T-16
Característica: Mejora visual de la interfaz
  Como persona que usa la aplicación
  Quiero una interfaz más legible y consistente (tipografía, tarjetas, tablas cebradas,
  botones y formularios coherentes, estados vacíos y errores visibles, navegación clara)
  Para trabajar con comodidad en cualquier pantalla y en cualquier ancho, incluido el teléfono

  Antecedente:
    Dada la hoja de estilos de la interfaz actualizada únicamente con cambios de presentación

  @CA1
  Escenario: El cambio toca un único archivo
    Cuando se inspecciona el diff completo del cambio
    Entonces el único archivo con modificaciones es la hoja de estilos de la interfaz
    Y ninguna plantilla, handler, servicio, dominio, persistencia ni test aparece modificado

  @CA2
  Escenario: Las clases existentes se conservan
    Dado que las 16 plantillas usan las clases existentes como contrato
    Cuando se revisa la hoja de estilos después del cambio
    Entonces todas las clases usadas por las plantillas (por ejemplo `acciones`, `boton-imprimir`, `contenido`, `detalle`, `encabezado`, `en-linea`, `error`, `formulario`, `grafico`, `grafico-mensaje`, `grafico-nota`, `grafico-svg`, `leyenda`, `metricas`, `navegacion`, `reporte`, `reporte-encabezado`, `reporte-nota` y `tabla`) siguen definidas o conservadas
    Y la suite de tests que verifica fragmentos del CSS permanece en verde

  @CA3
  Escenario: Sin JavaScript ni recursos externos
    Cuando se analiza la hoja de estilos entregada
    Entonces no hay scripts, importaciones ni URLs externas de ningún tipo
    Y no hay fuentes web ni CDNs: las tipografías son del stack del sistema (`system-ui, -apple-system, "Segoe UI", sans-serif`)

  @CA4
  Escenario: Contraste mínimo AA en todos los pares texto/fondo definidos
    Dado los pares navegación (#f9fafb sobre #1f2937), leyenda (#1e3a8a sobre #e0ecff), notas y mensajes vacíos (#57606a sobre #f7f7f8 y sobre #ffffff), error (#991b1b sobre #fde8e8), celdas y encabezados de tabla (#1f2328 sobre #f0f2f5, #ffffff y el fondo cebra), botón primario (#ffffff sobre #1f4e79) y contorno de foco (#1f4e79)
    Cuando se mide el contraste de cada par con una herramienta WCAG 2.1
    Entonces el texto normal alcanza al menos 4.5:1
    Y el texto grande y los componentes UI alcanzan al menos 3:1
    Y el contorno de foco se distingue tanto sobre fondos claros como sobre el encabezado oscuro

  @CA5
  Escenario: La interfaz funciona en un teléfono de 360 px
    Dado un viewport de 360 px de ancho
    Cuando se recorren todas las pantallas (inicio, proyectos, detalle de proyecto, backlog, formularios, sprints, poker, dashboard, defectos, reporte, error)
    Entonces no hay scroll horizontal de página ni corte de contenido esencial
    Y la navegación envuelve sus enlaces en varias líneas
    Y las tablas se desplazan en horizontal dentro de sí mismas sin romper el diseño
    Y las tarjetas de métricas se muestran en una columna y los formularios ocupan el ancho disponible

  @CA6
  Escenario: Tablas con filas cebradas legibles y print intacto
    Dado una pantalla con una tabla de datos (dashboard, histórico de Sprints, defectos o reporte)
    Cuando se observan las filas en pantalla
    Entonces las filas pares se distinguen de las impares por su fondo
    Y los encabezados de columna se leen con contraste AA
    Cuando se imprime el reporte (Ctrl+P)
    Entonces las filas no se cortan entre páginas y el encabezado de la tabla se repite en cada página

  @CA7
  Escenario: Las métricas del dashboard se presentan como tarjetas
    Dado el dashboard de un proyecto con métricas clave
    Cuando se observa la sección de métricas
    Entonces las métricas se muestran como tarjetas con la etiqueta discreta arriba y el valor destacado abajo
    Y el dashboard se ve igual en contenido que antes, sin tocar la página que lo muestra

  @CA8
  Escenario: Botones y formularios consistentes en toda la aplicación
    Dado un formulario cualquiera de la aplicación (crear proyecto, ítem, Sprint, defectos, backlog, poker)
    Cuando se comparan sus controles entre pantallas
    Y se recorre con el teclado (Tab/Shift+Tab)
    Entonces los botones comparten padding, radio, fondo y peso tipográfico, incluido Imprimir y las acciones en línea
    Y los campos de texto, numéricos, fecha, selección y área de texto comparten borde, radio y padding
    Y el foco visible se muestra en cada control, y el estado hover oscurece el botón primario de forma consistente

  @CA9
  Escenario: Estados vacíos y mensajes de error visibles con AA
    Dado un dashboard sin Sprints cerrados y un formulario que falla
    Cuando se observa el mensaje "Todavía no hay Sprints cerrados", las celdas "Sin datos" y "métricas no disponibles", y la pantalla de error
    Entonces cada uno se lee con contraste AA y señal visual reforzada (borde, fondo o peso), sin depender solo de cursiva

  @CA10
  Escenario: Navegación clara con foco visible
    Dado el encabezado con la navegación principal (Proyectos, Backlog, Sprints, Dashboard, Defectos)
    Cuando se recorre la navegación con teclado
    Entonces cada enlace muestra un contorno de foco distinguible sobre el fondo oscuro del encabezado
    Y los enlaces son legibles (#f9fafb sobre #1f2937) y el estado hover subraya el enlace
    Y el encabezado tiene aire suficiente respecto al contenido

  @CA11
  Escenario: La impresión del reporte conserva sus reglas y su bloque
    Dado el reporte de un proyecto o de un Sprint
    Cuando se imprime (Ctrl+P)
    Entonces la navegación, el encabezado, las acciones y el botón Imprimir no se muestran
    Y las tablas del reporte no cortan sus filas entre páginas y repiten el encabezado en cada página
    Y el bloque de estilos de impresión sigue acotado al reporte y contiene las reglas exigidas por la suite

  @CA12
  Escenario: Los gráficos SVG del dashboard no cambian de color
    Dado el dashboard con sus gráficos de velocidad y de horas
    Cuando se inspeccionan los colores de las barras y de la leyenda
    Y se revisa la hoja de estilos
    Entonces las barras de velocidad y de horas estimadas siguen en #1f4e79 y las de horas reales en #b45309
    Y ninguna regla de estilo altera esos gráficos (sin filtros, opacidad ni color sobre sus elementos)

  @CA13
  Escenario: La suite completa y el build quedan en verde
    Cuando se ejecuta la suite completa y el build de la aplicación
    Entonces todos los tests pasan sin modificar ningún test existente
    Y la construcción termina sin errores
    Y los tests del bloque de impresión del reporte y los de las páginas siguen en verde

  @CA14
  Escenario: Evidencias en verde y de refactor en docs/evidencias
    Dado que la implementación está completa
    Cuando se listan las evidencias producidas por la implementadora
    Entonces existe `docs/evidencias/T-16-green.txt` con la salida de la suite en verde y del build
    Y existe `docs/evidencias/T-16-refactor.txt` con el refactor sin cambios de comportamiento (o "sin cambios")
    Y no existe `docs/evidencias/T-16-red.txt`
=======
Feature: Mejora visual de la interfaz
  Como persona usuaria de Software Metrics & Estimation
  Quiero una interfaz visualmente clara y adaptable
  Para usar las pantallas existentes con comodidad en escritorio y móvil

  Background:
    Given la aplicación está disponible con datos de prueba para recorrer sus pantallas

  @T-16 @CA1
  Scenario Outline: Conservar las respuestas de las pantallas existentes
    Given los mismos datos de prueba antes y después de aplicar T-16
    When solicito "<metodo> <ruta>"
    Then la respuesta conserva el código HTTP que tenía antes de T-16
    And se conservan los enlaces, formularios, acciones y contenido existentes
    And no se agregan controles ni secciones funcionales nuevos

    Examples:
      | pantalla          | metodo | ruta                                                    |
      | Inicio            | GET    | /                                                       |
      | Proyectos         | GET    | /proyectos                                              |
      | Nuevo proyecto    | GET    | /proyectos/nuevo                                        |
      | Detalle proyecto  | GET    | /proyectos/PRJ-001                                      |
      | Backlog           | GET    | /proyectos/PRJ-001/backlog                              |
      | Nuevo ítem        | GET    | /proyectos/PRJ-001/backlog/nuevo                        |
      | Editar ítem       | GET    | /proyectos/PRJ-001/backlog/HU-01/editar                 |
      | Estimar           | POST   | /proyectos/PRJ-001/backlog/HU-01/estimar                |
      | Planning Poker    | GET    | /proyectos/PRJ-001/backlog/HU-01/poker                  |
      | Sprints           | GET    | /proyectos/PRJ-001/sprints                              |
      | Nuevo Sprint      | GET    | /proyectos/PRJ-001/sprints/nuevo                        |
      | Histórico         | GET    | /proyectos/PRJ-001/sprints/historico                    |
      | Defectos proyecto | GET    | /proyectos/PRJ-001/defectos                             |
      | Defectos global   | GET    | /defectos                                                |
      | Integrantes       | GET    | /proyectos/PRJ-001/integrantes                          |
      | Dashboard         | GET    | /proyectos/PRJ-001/dashboard                            |
      | Reporte proyecto  | GET    | /proyectos/PRJ-001/reporte                              |
      | Reporte Sprint    | GET    | /proyectos/PRJ-001/sprints/1/reporte                    |
      | Error             | GET    | /proyectos/PRJ-INEXISTENTE                              |

  @T-16 @CA2
  Scenario: Servir la hoja de estilos actualizada
    When solicito GET "/static/estilos.css"
    Then recibo una respuesta HTTP 200 con contenido CSS
    And las páginas existentes continúan enlazando esa hoja de estilos
    And la hoja contiene las variables visuales y el bloque delimitado de T-16

  @T-16 @CA2
  Scenario: Aplicar la presentación en todas las pantallas a ancho de escritorio
    When recorro cada pantalla accesible con viewport de escritorio
    Then el fondo, la navegación, las superficies, las acciones y la tipografía siguen la dirección visual decidida
    And no aparece una barra lateral ni controles de funciones excluidas

  @T-16 @CA2
  Scenario: Adaptar la interfaz a un viewport móvil
    When recorro cada pantalla accesible con viewport de hasta 640px
    Then la navegación, los formularios y el contenido se adaptan al ancho disponible
    And el contenido de tablas anchas puede desplazarse horizontalmente desde su contenedor
    And no aparece desplazamiento horizontal inesperado de toda la página

  @T-16 @CA3
  Scenario: Mantener tablas tabulares en móvil
    Given una tabla cuyo contenido supera el ancho del viewport
    When se muestra con viewport de hasta 640px
    Then la tabla conserva su formato tabular
    And no se aplica "display: block" a la tabla
    And se puede consultar el contenido restante mediante desplazamiento horizontal del contenedor

  @T-16 @CA4
  Scenario: Conservar los estilos de migas de T-15
    When comparo el bloque delimitado de T-15 con la versión anterior de "internal/web/estaticos/estilos.css"
    Then el contenido del bloque T-15 es idéntico
    And el bloque T-16 aparece separado y delimitado después del bloque T-15

  @T-16 @CA4
  Scenario: Conservar la impresión de reportes
    Given estoy viendo un reporte de proyecto o de Sprint
    When abro la vista previa de impresión
    Then se conservan intactas las reglas "@media print" acotadas a ".reporte"
    And la navegación y los controles de impresión mantienen el comportamiento de impresión existente

  @T-16 @CA5
  Scenario: Limitar la implementación a CSS sin dependencias externas
    When reviso los cambios de implementación de T-16
    Then el único archivo de implementación modificado es "internal/web/estaticos/estilos.css"
    And no hay cambios de plantillas, rutas, dominio, servicios ni JavaScript
    And no se agregan librerías, CDN ni paso de build
    And los gráficos SVG del servidor y el uso existente de "window.print" no se modifican
>>>>>>> 6b230cf (docs: especificación y escenarios BDD de T-16 (mejora visual, solo CSS))
