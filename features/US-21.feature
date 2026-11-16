# language: es
# Estado: Propuesta (pendiente de implementación TDD)
# Los escenarios describen el comportamiento acordado para CA1–CA2 y las
# decisiones D1–D7 del enunciado (dos descargas PDF enlazadas desde las páginas
# de reporte HTML de US-20 sin quitar lo existente; contenido tomado de
# ServicioReportes tal cual, sin recalcular cifras; PDF válido con cabeceras y
# estructura comprobable; errores 404/405/500; reporte sin datos con "No hay
# historias" y "No hay defectos").
# Un solo PR y un solo ciclo RED/GREEN/REFACTOR con GREEN por tandas internas.
# No se escenifican: verificación de texto interno con compresión desactivada en
# todos los casos, normalización detallada de caracteres fuera de cp1252, formato
# tipográfico exacto (márgenes/fuentes/pie), ni la verificación estructural de
# imports del paquete generador: se verifican con tests sin escenario BDD propio.

@US-21
Característica: Exportar informe a PDF
  Como Product Owner
  Quiero exportar el informe de un proyecto o de un Sprint a PDF
  Para compartirlo con el cliente

  @CA1
  Escenario: Descargar el PDF del reporte de un proyecto con el contenido del reporte
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene una historia completada y una historia pendiente
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con velocidad registrada
    Y que el proyecto "PRJ-001" tiene defectos registrados
    Cuando se descarga el reporte en PDF del proyecto "PRJ-001"
    Entonces la respuesta es un archivo PDF descargable del proyecto "PRJ-001"
    Y el PDF es válido
    Y el PDF incluye el contenido del reporte del proyecto

  @CA1
  Escenario: Descargar el PDF del reporte de un Sprint con el contenido del reporte
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado con objetivo "Entregar el núcleo del sistema", fecha de inicio "2026-02-02" y fecha de fin "2026-02-27"
    Y que el Sprint 1 tiene historias asignadas y métricas registradas
    Y que el Sprint 1 tiene defectos detectados y resueltos
    Cuando se descarga el reporte en PDF del Sprint 1 del proyecto "PRJ-001"
    Entonces la respuesta es un archivo PDF descargable del Sprint 1 del proyecto "PRJ-001"
    Y el PDF es válido
    Y el PDF incluye el contenido del reporte del Sprint

  @CA1
  Escenario: Las páginas HTML de reporte ofrecen el enlace para descargar el PDF
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 activo
    Cuando se visita el reporte del proyecto "PRJ-001"
    Entonces se ofrece un enlace para descargar el PDF del reporte del proyecto
    Y se conserva el botón Imprimir
    Cuando se visita el reporte del Sprint 1 del proyecto "PRJ-001"
    Entonces se ofrece un enlace para descargar el PDF del reporte del Sprint
    Y se conserva el botón Imprimir

  @CA1
  Escenario: El PDF de un proyecto inexistente responde 404
    Cuando se descarga el reporte en PDF del proyecto "PRJ-999"
    Entonces se muestra un mensaje de proyecto no encontrado
    Y la respuesta tiene el código de error 404

  @CA1
  Escenario: El PDF de un Sprint inexistente responde 404
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Cuando se descarga el reporte en PDF del Sprint 9 del proyecto "PRJ-001"
    Entonces se muestra un mensaje de Sprint no encontrado
    Y la respuesta tiene el código de error 404

  @CA1
  Escenario: Un número de Sprint inválido al descargar el PDF responde 404 sin invocar el servicio
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Cuando se descarga el reporte en PDF del Sprint "abc" del proyecto "PRJ-001"
    Entonces la respuesta tiene el código de error 404
    Y no se consulta el reporte del Sprint al servicio

  @CA1
  Escenario: Un método no permitido al descargar el PDF responde 405
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Cuando se intenta enviar el formulario de descarga del reporte en PDF del proyecto "PRJ-001"
    Entonces la respuesta tiene el código de error 405

  @CA2
  Escenario: Un proyecto sin datos genera un PDF válido con mensajes vacíos
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" no tiene historias
    Y que el proyecto "PRJ-001" no tiene defectos
    Y que el proyecto "PRJ-001" no tiene Sprints cerrados
    Cuando se descarga el reporte en PDF del proyecto "PRJ-001"
    Entonces la respuesta es un archivo PDF descargable del proyecto "PRJ-001"
    Y el PDF es válido
    Y el PDF muestra el mensaje "No hay historias"
    Y el PDF muestra el mensaje "No hay defectos"

  @CA2
  Escenario: Un Sprint sin datos genera un PDF válido con mensajes vacíos
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 activo sin historias asignadas
    Y que el proyecto "PRJ-001" no tiene defectos
    Cuando se descarga el reporte en PDF del Sprint 1 del proyecto "PRJ-001"
    Entonces la respuesta es un archivo PDF descargable del Sprint 1 del proyecto "PRJ-001"
    Y el PDF es válido
    Y el PDF muestra el mensaje "No hay historias asignadas"
    Y el PDF muestra el mensaje "No hay defectos"

  @CA2
  Escenario: Un Sprint cerrado sin snapshot de esfuerzo genera un PDF válido con métricas no disponibles
    Dado un proyecto existente identificado por "PRJ-001", con nombre "Sistema de Métricas", fecha de inicio "2026-01-05" y fecha de fin "2026-12-31"
    Y que el proyecto "PRJ-001" tiene el Sprint 1 cerrado sin snapshot de esfuerzo
    Y que el Sprint 1 tiene una historia asignada con 5 Story Points
    Y que el proyecto "PRJ-001" no tiene defectos
    Cuando se descarga el reporte en PDF del Sprint 1 del proyecto "PRJ-001"
    Entonces la respuesta es un archivo PDF descargable del Sprint 1 del proyecto "PRJ-001"
    Y el PDF es válido
    Y el PDF muestra el mensaje "métricas no disponibles"