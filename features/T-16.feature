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
